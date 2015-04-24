package Elevator

import (
	."./../Driver"
	."./../OrderRegister"
	."time"
	."./../Udp"
	."./../Cost"
	"encoding/json"
	"net"
	."strings"
	."fmt"
)



const localPort = 20016
const broadcastPort = 20017
const message_size = 1024
var receive_ch = make(chan Udp_message)
var send_ch = make(chan Udp_message)



var myFloor = -1
var lastFloor = 0
var myDirection = -1	// -1 = står i ro, 1 = opp, 0 = ned 
var doorOpen = false



var openDoor = make(chan bool)
var gotMessage = make(chan string)


//Elevfunc skal ha initfunksjon, alle elevfunksjoner og de fleste variabler, troooor jeg

func Init() {

	err := Udp_init(localPort, broadcastPort, message_size, send_ch, receive_ch)
	if err != nil {
		println("Error during udp-init")
		return
	}
	Elev_init()		
	DeleteAllOrders()
	for Elev_get_floor_sensor_signal() != 0 {
		Elev_set_motor_direction(-300)
	}
	Elev_set_motor_direction(50)
	Sleep(2000*Microsecond)
	Elev_set_motor_direction(0)
	Elev_set_floor_indicator(0)	
	myDirection = -1
	lastFloor = 0
	myFloor = 0
	
	//Henter egen ip-adresse = 147
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
	println("Init completed")
}



func PrintStatus() {

	for{
		println("Direction: ", myDirection)
		Println("UP: ",Up)
		Println("DOWN: ", Down)
		Println("INSIDE: ", Inside)
		Sleep(1*Second)
	}
}


func RunElevator() {

	for {
		if doorOpen {
			Sleep(100*Millisecond)
		} else {
			if EmptyQueue() {
				myDirection = -1
			}
			
			if myDirection == 0 {
				Elev_set_motor_direction(-300)
			} else if myDirection == 1 {
				Elev_set_motor_direction(300)
			}

			Sleep(100*Millisecond)
		}
	}
}




func UpdateFloor() {
	for{
		myFloor = Elev_get_floor_sensor_signal()
		
		if lastFloor != myFloor {	
		    if (myFloor != -1) {
		        floorReached(myFloor)
		    } else {
		    	Elev_set_door_open_lamp(0)
		    }
		}
		Sleep(100*Millisecond)
	}
}




func floorReached(floor int) {
	lastFloor = floor
	Elev_set_floor_indicator(floor)
	
	orderDir := GetOrder(myDirection, floor) 
	
	if orderDir != 2 {				
		if myDirection == 1 {
			Elev_set_motor_direction(-100)
		} else if (myDirection == 0) {
			Elev_set_motor_direction(100)
		}
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
	
		openDoor <- true
		
		deleteOrder := Order{lastFloor, myDirection, floor, orderDir, true, false, Up, Down, Inside}
		go sendOrder(deleteOrder)
		
	} else if (floor == 0) {				//Stops, so the elevator do not pass 1. floor
		Elev_set_motor_direction(100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		myDirection = 1
		
	} else if (floor == N_FLOORS-1) {		//Stops, so the elevator do not pass 4. floor
		Elev_set_motor_direction(-100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		myDirection = 0
	}
	
}




//Registers if any up-buttons is pushed
func CheckButtonCallUp() {
	
	for{
		for i:=0; i<N_FLOORS-1; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_UP, i)) {
				
				if (myDirection == -1 && myFloor == i) || (doorOpen && myFloor == i) {
					openDoor <- true
				} else {
					newOrder := Order{lastFloor, myDirection, i, 1, false, true, Up, Down, Inside}
					go sendOrder(newOrder)
				}
			}
		}
		Sleep(50*Millisecond)
	}
}




//Registers if any down-buttons is pushed
func CheckButtonCallDown() {

	for{
		for i:=1; i< N_FLOORS; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_DOWN, i)) {
				
				if (myDirection == -1 && myFloor == i) || (doorOpen && myFloor == i) {
					openDoor <- true
				} else {
					newOrder := Order{lastFloor, myDirection, i, 0, false, true, Up, Down, Inside}
					go sendOrder(newOrder)
				}
			}
		}
		Sleep(50*Millisecond)
	}
}




//Registers if any command-buttons is pushed
func CheckButtonCommand() {

	for{
		for i:=0; i<N_FLOORS; i++ {
			if (Elev_get_button_signal(BUTTON_COMMAND, i)) {
			
				if (myDirection == -1 && myFloor == i) || (doorOpen && myFloor == i) {
					openDoor <- true
				} else {
					newOrder := Order{myFloor, myDirection, i, -1, false, true, Up, Down, Inside}
					if EmptyQueue() {
						UpdateMyOrders(newOrder)
						setDirection()
					} else {
						UpdateMyOrders(newOrder)
					}
				}
			}
		}
		Sleep(100*Millisecond)
	}
}




func DoorControl() {

	timer := NewTimer(Hour*3)
	for{
	
		select {
			case <- openDoor:
				doorOpen = true
				Elev_set_door_open_lamp(1)
				timer.Reset(Second*3)
				if Elev_get_floor_sensor_signal() == lastFloor {
					deleteOrder := Order{lastFloor, myDirection, lastFloor, myDirection, true, false, Up, Down, Inside}
					go sendOrder(deleteOrder)
				}
				
			case <- timer.C:
				Elev_set_door_open_lamp(0)
				doorOpen = false
				go setDirection()
		}
	}
}




func setDirection(){
	
	if (EmptyQueue()) {
		myDirection = -1
		
	} else if GetOrder(2, lastFloor) != 2 {
		if myDirection == 1 && !CheckOrdersAboveFloor(lastFloor) {
			myDirection = 0
			println("åpner dør fra setDirection. myDir=", myDirection)
			openDoor <- true
			
		} else if myDirection == 0 && !CheckOrdersUnderFloor(lastFloor) {
			myDirection = 1
			println("Åpner dør fra SetDirection. myDir=", myDirection)
			openDoor <- true
		}
	} else {
		if (myDirection == 0) && !(CheckOrdersUnderFloor(lastFloor)) {
			myDirection = 1
		} else if (myDirection == 1) && !(CheckOrdersAboveFloor(lastFloor)) {
			myDirection = 0
		} else if myDirection == -1 {
			if CheckOrdersAboveFloor(lastFloor) {
				myDirection = 1
			} else if CheckOrdersUnderFloor(lastFloor) {
				myDirection = 0
			}
		}
	}
}






//Receives messages from other elevators continuous
func ReceiveMessage() {
	
	for{
		var receivedMessage Udp_message
		receivedMessage = <- receive_ch
		
		IP := getIP(receivedMessage.Raddr)
		
		
		var receivedOrder Order
		err := json.Unmarshal(receivedMessage.Data[:receivedMessage.Length], &receivedOrder)
		if (err != nil) {
			Println("Receive Order Error: ", err)
			Println("when decoding: ", string(receivedMessage.Data))
		}

	
		if receivedOrder.NewOrder || receivedOrder.OrderHandled {
			receiveOrder(receivedOrder)
			
		} else if IP != MyAddress {
		
			newElevator := true	
			for key,_ := range Elevators {
				if key == IP {
					newElevator = false
				}
			}
	
			if newElevator {
				go setMessageTimer(IP)
			} else {
				gotMessage <- IP
			}
			
			Elevators[IP] = ElevStatus{LastFloor: receivedOrder.MyFloor, Direction: receivedOrder.MyDirection, Up: receivedOrder.Up, Down: receivedOrder.Down, Inside: receivedOrder.Inside} 			
		}
		Sleep(1*Millisecond)
	}
}






//Receives orders from other elevators
func receiveOrder(receivedOrder Order) {
	
	go SetButtonLight(receivedOrder)
	
	if receivedOrder.OrderHandled {
		UpdateMyOrders(receivedOrder)
		return
	}
	
	if receivedOrder.NewOrder {
		for _, val := range Elevators {
			if (val.Up)[receivedOrder.Floor] && receivedOrder.Direction == 1 {
				return
			} else if (val.Down)[receivedOrder.Floor] && receivedOrder.Direction == 0 {
				return
			}
		}
	}
	
	if (myDirection == -1 && myFloor == receivedOrder.Floor) || (doorOpen && myFloor == receivedOrder.Floor) {
		openDoor <- true
	
	} else if GetCost(lastFloor, myDirection, receivedOrder.Floor, receivedOrder.Direction) == 1 {
		if EmptyQueue() {
			UpdateMyOrders(receivedOrder)
			setDirection()
		} else {
			UpdateMyOrders(receivedOrder)
		}
	}
}





func sendOrder(order Order) {
	b, err := json.Marshal(order)
	
	if (err != nil) {
		println("Send Order Error: ", err)
	}
	
	var message Udp_message
	message.Raddr = "broadcast"
	message.Data = b
	message.Length = 1024
	
	send_ch <- message
}






// go fra main. sender hvert sekund oppdatering på floor og direction
func SendUpdateMessage() {
	for {
		order := Order{lastFloor, myDirection, -1, -1, false, false, Up, Down, Inside}
		b, err := json.Marshal(order)
		
		if (err != nil) {
			println("Send Order Error: ", err)
		}
		
		var message Udp_message
		message.Raddr = "broadcast"
		message.Data = b
		message.Length = 1024
		
		send_ch <- message
		Sleep(100*Millisecond)
	}
	
}






func setMessageTimer(address string) {
	
	timer := NewTimer(3*Hour)
	for {
		select {
		case <- timer.C:
			for i:=0; i<N_FLOORS; i++ {
				if (Elevators[address].Up)[i] {
					order := Order{lastFloor, myDirection, i, 1, false, true, Up, Down, Inside}
					go sendOrder(order)
				}
				if (Elevators[address].Down)[i] {
					order := Order{lastFloor, myDirection, i, 0, false, true, Up, Down, Inside}
					go sendOrder(order)
				}	
			}
			delete (Elevators, address)
			return
			
		case receivedAddress := <- gotMessage:
			if receivedAddress == address {
				timer.Reset(3*Second)
			}
		}
	}
}






//Returns last three numbers of IP-address
func getIP(address string) string {
	splitaddr := Split(address, ".")
	splitip := Split(splitaddr[3], ":")
	myAddress := splitip[0]
	return myAddress
}



//Used to stop the program and elevator from running
func Stop(ch chan int) {
	for {
		if Elev_get_stop_signal() != 0 {
			ch <- 1
		}
		Sleep(100*Millisecond)
	}
}



