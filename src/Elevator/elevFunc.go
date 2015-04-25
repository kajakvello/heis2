package Elevator

import (
	."./../Driver"
	."./../OrderRegister"
	."./../Timer"
	."./../Udp"
	."./../Cost"
	"encoding/json"
	"net"
	."strings"
	."fmt"
	."time"
)



const localPort = 20016
const broadcastPort = 20017
const message_size = 1024






//Elevfunc skal ha initfunksjon, alle elevfunksjoner og de fleste variabler, troooor jeg

func Init() {

	err := Udp_init(localPort, broadcastPort, message_size, Send_ch, Receive_ch)
	if err != nil {
		println("Error during udp-init")
		return
	}
	Elev_init()		
	DeleteAllOrders()
	
	//Runs elevator to 1. floor
	for Elev_get_floor_sensor_signal() != 0 {
		Elev_set_motor_direction(-300)
	}
	breakElevator(0)
	Elev_set_floor_indicator(0)	
	MyDirection = -1
	LastFloor = 0
	MyFloor = 0
	Defect = false
	
	//Gets own IP-address (last three numbers)
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				splitip := Split(ip, ".")
				MyAddress = splitip[3]
			}
		}
	}
	go SelfAliveTimer()
	println("Init completed")
}



/////////////////////////////////////////////////////////////////////////////////////////////



func PrintStatus() {

	for{
		//Println("UP: ", GlobalUp)
		//Println("DOWN: ", GlobalDown)	
		
		
		println("Direction: ", MyDirection)
		Println("UP: ",Up)
		Println("DOWN: ", Down)
		Println("INSIDE: ", Inside)
		Sleep(3*Second)
		
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



func RunElevator() {

	for {
		if DoorOpen {
			Sleep(100*Millisecond)
			
		} else if Defect {
			for Elev_get_floor_sensor_signal() != 0 {
				Elev_set_motor_direction(-300)
			}
			breakElevator(0)
			
		} else {
			if EmptyQueue() {
				MyDirection = -1
				Alive <- MyAddress
			}
			
			if MyDirection == 0 {
				Elev_set_motor_direction(-300)
			} else if MyDirection == 1 {
				Elev_set_motor_direction(300)
			}

			Sleep(100*Millisecond)
		}
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



func UpdateFloor() {
	for{
		MyFloor = Elev_get_floor_sensor_signal()
		
		if LastFloor != MyFloor {	
		    if (MyFloor != -1) {
		        floorReached(MyFloor)
		    } else {
		    	Elev_set_door_open_lamp(0)		//trenger vi denne dobbelsjekkingen?
		    }
		}
		Sleep(100*Millisecond)
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



func floorReached(floor int) {
	LastFloor = floor
	Elev_set_floor_indicator(floor)
	
	orderDir := GetOrder(MyDirection, floor) 
	
	if orderDir != 2 {						//Stops if order on floor
		breakElevator(MyDirection)
		OpenDoor <- orderDir
		
	} else if (floor == 0) {				//Stops, so the elevator do not pass 1. floor
		breakElevator(0)
		MyDirection = 1
		
	} else if (floor == N_FLOORS-1) {		//Stops, so the elevator do not pass N. floor
		breakElevator(1)
		MyDirection = 0
	}
	
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Registers if any up-buttons is pushed and sends order
func CheckButtonCallUp() {
	
	for{
		for i:=0; i<N_FLOORS-1; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_UP, i)) {
				
				if (MyDirection == -1 && MyFloor == i) || (DoorOpen && MyFloor == i) {
					OpenDoor <- 1
				} else {
					newOrder := Order{LastFloor, MyDirection, i, 1, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(newOrder)
				}
			}
		}
		Sleep(100*Millisecond)
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Registers if any down-buttons is pushed and sends order
func CheckButtonCallDown() {

	for{
		for i:=1; i< N_FLOORS; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_DOWN, i)) {
				
				if (MyDirection == -1 && MyFloor == i) || (DoorOpen && MyFloor == i) {
					OpenDoor <- 0
				} else {
					newOrder := Order{LastFloor, MyDirection, i, 0, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(newOrder)
				}
			}
		}
		Sleep(100*Millisecond)
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Registers if any command-buttons is pushed
func CheckButtonCommand() {

	for{
		for i:=0; i<N_FLOORS; i++ {
			if (Elev_get_button_signal(BUTTON_COMMAND, i)) {
			
				if (MyDirection == -1 && MyFloor == i) || (DoorOpen && MyFloor == i) {
					OpenDoor <- -1
				} else {
					newOrder := Order{MyFloor, MyDirection, i, -1, false, true, DoorOpen, Up, Down, Inside}
					if EmptyQueue() {
						UpdateMyOrders(newOrder)
						SetDirectionToOrder(-1)
					} else {
						UpdateMyOrders(newOrder)
					}
				}
			}
		}
		Sleep(100*Millisecond)
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Receives messages from other elevators continuous
func ReceiveMessage() {
	
	for{
		var receivedMessage Udp_message
		receivedMessage = <- Receive_ch
		
		IP := getIP(receivedMessage.Raddr)
		
		var receivedOrder Order
		err := json.Unmarshal(receivedMessage.Data[:receivedMessage.Length], &receivedOrder)
		if (err != nil) {
			Println("Receive Order Error: ", err)
			Println("when decoding: ", string(receivedMessage.Data))
		}

	
		if receivedOrder.NewOrder || receivedOrder.OrderHandled {
			go receiveOrder(receivedOrder, IP)
			
		} else if IP != MyAddress {
		
			newElevator := true	
			for key,_ := range Elevators {
				if key == IP {
					newElevator = false
				}
			}
	
			if newElevator {
				go SetMessageTimer(IP)
				go AliveTimer(IP)
				Elevators[IP] = ElevStatus{LastFloor: receivedOrder.MyFloor, 
					Direction: receivedOrder.MyDirection, DoorOpen: receivedOrder.DoorOpen, 
					Up: receivedOrder.Up, Down: receivedOrder.Down, 
					Inside: receivedOrder.Inside, Defect: false} 
			} else {
			
				GotMessage <- IP
				Alive <- IP
				def := Elevators[IP].Defect
				Elevators[IP] = ElevStatus{LastFloor: receivedOrder.MyFloor, 
					Direction: receivedOrder.MyDirection, DoorOpen: receivedOrder.DoorOpen, 
					Up: receivedOrder.Up, Down: receivedOrder.Down, 
					Inside: receivedOrder.Inside, Defect: def} 
			}
				
		} else {
			Alive <- IP
		}
		
		//Sleep(10*Millisecond) Merker ingen forskjell med denne fra 0-100 millisekund
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Receives orders from all elevators
func receiveOrder(order Order, IP string) {
	
	go SetButtonLight(order, IP)
	
	if order.OrderHandled {		//sletter ordre
		UpdateMyOrders(order)
		UpdateGlobalOrders(order)
		return
	}
	
	if order.NewOrder {
		
		//Check if the order already exists
		//Tror vi trenger denne sjekken, om ikke vil nytt knappetrykk overskrive gammel bestilling. Men da vil ikke heisene ta over bestillinger fra hverandre...
		
		if ((order.Direction == 1 && Up[order.Floor]) || (order.Direction == 0 && Down[order.Floor])) && !Defect {
			return
		}
		for _, val := range Elevators {
			if ((order.Direction == 1 && val.Up[order.Floor]) || (order.Direction == 0 && val.Down[order.Floor])) && !val.Defect {
				return
			}
		}
		
		
		//Mulig disse gjør at heisene tror en annen heis er i etg, selv om den akkurat har begynt å gå.
		
		//Ckeck if i am on the right floor
		if (MyFloor == order.Floor) && (DoorOpen || MyDirection == -1 ) && (!Defect) {
			OpenDoor <- order.Direction
			return
		}
		
		//Check if any other elevator is on the right floor
		for _, val := range Elevators {
			if (order.Floor == val.LastFloor) && (val.DoorOpen || val.Direction == -1) {
				if !val.Defect {
					return
				}
			}
		}
	}
	
	UpdateGlobalOrders(order)
	
	if !Defect && GetCost(LastFloor, MyDirection, order.Floor, order.Direction, MyAddress) == 1 {
		if EmptyQueue() {
			UpdateMyOrders(order)
			SetDirectionToOrder(order.Direction)
		} else {
			UpdateMyOrders(order)
		}
	}
	
}



/////////////////////////////////////////////////////////////////////////////////////////////



// go fra main. sender hvert sekund oppdatering på floor og direction
func SendUpdateMessage() {
	for {
		order := Order{LastFloor, MyDirection, -1, -1, false, false, DoorOpen, Up, Down, Inside}
		b, err := json.Marshal(order)
		
		if (err != nil) {
			println("Send Order Error: ", err)
		}
		
		var message Udp_message
		message.Raddr = "broadcast"
		message.Data = b
		message.Length = 1024
		
		
		Send_ch <- message
		Sleep(10*Millisecond)
	}
	
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Mulig denne må justeres på forskjellige heiser
func breakElevator(direction int) {
	if direction == 1 {
		//Elev_set_motor_direction(-100)
		//Sleep(2000*Millisecond)
		Elev_set_motor_direction(0)
	} else if direction == 0 {
		//Elev_set_motor_direction(100)
		//Sleep(2000*Millisecond)
		Elev_set_motor_direction(0)
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Returns last three numbers of IP-address
func getIP(address string) string {
	splitaddr := Split(address, ".")
	splitip := Split(splitaddr[3], ":")
	myAddress := splitip[0]
	return myAddress
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Used to stop the program and elevator from running
func Stop(ch chan int) {
	for {
		if Elev_get_stop_signal() != 0 {
			ch <- 1
		}
		Sleep(100*Millisecond)
	}
}



