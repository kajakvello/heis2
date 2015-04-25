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

<<<<<<< HEAD
var myFloor = -1
var lastFloor = 0
var myDirection = -1	// -1 = står i ro, 1 = opp, 0 = ned 
var doorOpen = false
=======


>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f



//Elevfunc skal ha initfunksjon, alle elevfunksjoner og de fleste variabler, troooor jeg

func Init() {

	err := Udp_init(localPort, broadcastPort, message_size, Send_ch, Receive_ch)
	if err != nil {
		println("Error during udp-init")
		return
	}
	Elev_init()		
	DeleteAllOrders()
	for Elev_get_floor_sensor_signal() != 0 {
		Elev_set_motor_direction(-300)
	}
	Elev_set_motor_direction(100)
	Sleep(2000*Microsecond)
	Elev_set_motor_direction(0)
	Elev_set_floor_indicator(0)	
	MyDirection = -1
	LastFloor = 0
	MyFloor = 0
	Defekt = false
	
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
	go SelfAliveTimer(MyAddress)
	println("Init completed")
}




func PrintStatus() {

	for{
		//Println("UP: ", GlobalUp)
		//Println("DOWN: ", GlobalDown)	
		
		Println("Defekt = ", Defekt)
		
		Sleep(2*Second)
		/*
		println("Direction: ", MyDirection)
		Println("UP: ",Up)
		Println("DOWN: ", Down)
		Println("INSIDE: ", Inside)
		Sleep(2*Second)
		*/
	}
}




func RunElevator() {

	for {
		if DoorOpen {
			Sleep(100*Millisecond)
<<<<<<< HEAD
		} else {
			if (EmptyQueue()){
			
				for _, val := range Elevators {
					for i:=0; i<N_FLOORS; i++ {
						if val.Up[i] {
							if getCost(i, 1) == 1 {
								newOrder := Order{myFloor, myDirection, i, 1, false, true, true, Up, Down, Inside}
								sendOrder(newOrder)
							}
						}
						if val.Down[i] {
							if getCost(i, 0) == 1 {
								newOrder := Order{myFloor, myDirection, i, 0, false, true, true, Up, Down, Inside}
								sendOrder(newOrder)
							}
						}
					}
				}
				Sleep(10*Millisecond)
			}
			
			if EmptyQueue() {
				myDirection = -1
=======
		/*} else if Defekt {
			for Elev_get_floor_sensor_signal() != 0 {
				Elev_set_motor_direction(-300)
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
			}
			Elev_set_motor_direction(50)
			Sleep(2000*Microsecond)
			Elev_set_motor_direction(0)
		}*/ } else {
			if EmptyQueue() {
				MyDirection = -1
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




func UpdateFloor() {
	for{
		MyFloor = Elev_get_floor_sensor_signal()
		
		if LastFloor != MyFloor {	
		    if (MyFloor != -1) {
		        floorReached(MyFloor)
		    } else {
		    	Elev_set_door_open_lamp(0)
		    }
		}
		Sleep(100*Millisecond)
	}
}




func floorReached(floor int) {
	LastFloor = floor
	Elev_set_floor_indicator(floor)
	orderOnFloor, IP := GetOrder(myDirection, floor)
	
<<<<<<< HEAD
	if (orderOnFloor) {				//Breaks and stops, if orders on floor
		if myDirection == 1 {
=======
	orderDir := GetOrder(MyDirection, floor) 
	
	if orderDir != 2 {						//Stops if order on floor
		if MyDirection == 1 {
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
			Elev_set_motor_direction(-100)
		} else if (MyDirection == 0) {
			Elev_set_motor_direction(100)
		}
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
<<<<<<< HEAD
		
		if IP != MyAddress {
			updateOrder := Order{myFloor, myDirection, floor, myDirection, false, false, true, Up, Down, Inside}
			sendOrder(updateOrder)
		}
		openDoor <- true
=======
		OpenDoor <- orderDir
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
		
	} else if (floor == 0) {				//Stops, so the elevator do not pass 1. floor
		Elev_set_motor_direction(100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		MyDirection = 1
		
	} else if (floor == N_FLOORS-1) {		//Stops, so the elevator do not pass 4. floor
		Elev_set_motor_direction(-100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		MyDirection = 0
	}
	
}




//Registers if any up-buttons is pushed
func CheckButtonCallUp() {
	
	for{
		for i:=0; i<N_FLOORS-1; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_UP, i)) {
				
				if (MyDirection == -1 && MyFloor == i) || (DoorOpen && MyFloor == i) {
					OpenDoor <- 1
				} else {
<<<<<<< HEAD
					newOrder := Order{myFloor, myDirection, i, 1, false, true, false, Up, Down, Inside}
					go sendOrder(newOrder)
=======
					newOrder := Order{LastFloor, MyDirection, i, 1, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(newOrder)
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
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
				
				if (MyDirection == -1 && MyFloor == i) || (DoorOpen && MyFloor == i) {
					OpenDoor <- 0
				} else {
<<<<<<< HEAD
					newOrder := Order{myFloor, myDirection, i, 0, false, true, false, Up, Down, Inside}
					go sendOrder(newOrder)
=======
					newOrder := Order{LastFloor, MyDirection, i, 0, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(newOrder)
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
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
			
				if (MyDirection == -1 && MyFloor == i) || (DoorOpen && MyFloor == i) {
					OpenDoor <- -1
				} else {
<<<<<<< HEAD
					newOrder := Order{myFloor, myDirection, i, -1, false, true, false, Up, Down, Inside}
					if EmptyQueue() {
						UpdateMyOrders(newOrder, "")
						setDirection()
=======
					newOrder := Order{MyFloor, MyDirection, i, -1, false, true, DoorOpen, Up, Down, Inside}
					if EmptyQueue() {
						UpdateMyOrders(newOrder)
						SetDirectionToOrder(-1)
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
					} else {
						UpdateMyOrders(newOrder, "")
					}
				}
			}
		}
		Sleep(100*Millisecond)
	}
}




<<<<<<< HEAD
func DoorControl() {

	timer := NewTimer(3*Hour)
	for{
	
		select {
			case <- openDoor:
				doorOpen = true
				Elev_set_door_open_lamp(1)
				timer.Reset(Second*3)
				if Elev_get_floor_sensor_signal() == lastFloor {
					deleteOrder := Order{myFloor, myDirection, myFloor, -1, true, false, false, Up, Down, Inside}
					go sendOrder(deleteOrder)
				}
				
			case <- timer.C:
				Elev_set_door_open_lamp(0)
				doorOpen = false
				setDirection()
		}
	}
}




func setDirection(){

	if (EmptyQueue()) {
		myDirection = -1
		
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





//TODO: Gange opp cost med 1000 og legge til IP, slik at cost aldri er lik
//Nå beregnes costen for mange ganger på høyest IP (for some reason), så flere heiser tar samme best.
func getCost(orderFloor int, orderDirection int) int {
	equalCost := []string{}
	
	//Find my cost:
	myCost := int(Abs(float64(orderFloor - myFloor))*3)
	
	for i:=0; i<N_FLOORS; i++ {
		if Up[i] || Down[i] || Inside[i] {
			myCost += 4
		} 
	}
	if orderDirection != myDirection {
		myCost += 1
	}
	
	//Check if other elevator got lower cost:
	for IP, val := range Elevators {
		
		elevCost := int(Abs(float64(orderFloor - val.LastFloor))*3)
		
		for i:=0; i<N_FLOORS; i++ {
			if val.Up[i] || val.Down[i] || val.Inside[i] {
				elevCost += 4
			}
		}
		if orderDirection != val.Direction {
			elevCost += 1
		}
		
		if elevCost < myCost {
			return 0
		} else if elevCost == myCost {
			equalCost = append(equalCost, IP)
		}
	}
	
	if len(equalCost) != 0 {
		myAddr, _ := strconv.Atoi(MyAddress)
		for i:=0; i<len(equalCost); i++ {
			elevAddr, _ := strconv.Atoi(equalCost[i])
			if elevAddr < myAddr {
				println("Lik cost. Myaddr =", myAddr, " mens elevaddr =", elevAddr, "så jeg chiller")
				return 0
			}
		}
	}
	Sleep(Millisecond*5)
	println("Got it!! Fra IP: ", MyAddress)
	return 1

}


=======
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f


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

	
<<<<<<< HEAD
		if receivedOrder.NewOrder || receivedOrder.OrderHandled || receivedOrder.UpdateOrder {
			receiveOrder(receivedOrder, IP)
		}
	
		if IP != MyAddress {
		
			newElevator := true	
			for IP,_ := range Elevators {
				if IP == IP {
=======
		if receivedOrder.NewOrder || receivedOrder.OrderHandled {
			go receiveOrder(receivedOrder)
			
		} else if IP != MyAddress {
		
			newElevator := true	
			for key,_ := range Elevators {
				if key == IP {
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
					newElevator = false
				}
			}
	
			if newElevator {
				go SetMessageTimer(IP)
				go AliveTimer(IP)
				Elevators[IP] = ElevStatus{LastFloor: receivedOrder.MyFloor, 
					Direction: receivedOrder.MyDirection, 
					Up: receivedOrder.Up, Down: receivedOrder.Down, Inside: receivedOrder.Inside, 
					Defekt: false} 
			} else {
<<<<<<< HEAD
				gotMessage <- IP
			}
	
			Elevators[IP] = ElevStatus{LastFloor: receivedOrder.MyFloor, Direction: receivedOrder.MyDirection, Up: receivedOrder.Up, Down: receivedOrder.Down, Inside: receivedOrder.Inside} 			
		}
		Sleep(Millisecond*1)
	}
}



//Returns last three numbers of IP-address
func getIP(address string) string {
	splitaddr := Split(address, ".")
	splitip := Split(splitaddr[3], ":")
	MyAddress := splitip[0]
	return MyAddress
}





func setMessageTimer(IP string) {
	
	timer := NewTimer(3*Hour)			//TODO: fikse timer?
	for {
		select {
		case <- timer.C:
			for i:=0; i<N_FLOORS; i++ {
				if (Elevators[IP].Up)[i] {
					order := Order{myFloor, myDirection, i, 1, false, true, false, Up, Down, Inside}
					go sendOrder(order)
				}
				if (Elevators[IP].Down)[i] {
					order := Order{myFloor, myDirection, i, 0, false, true, false, Up, Down, Inside}
					go sendOrder(order)
				}	
			}
			delete (Elevators, IP)
			return
			
		case receivedAddress := <- gotMessage:
			if receivedAddress == IP {
				timer.Reset(3*Second)
=======
			
				GotMessage <- IP
				Alive <- IP
				def := Elevators[IP].Defekt
				Elevators[IP] = ElevStatus{LastFloor: receivedOrder.MyFloor, 
					Direction: receivedOrder.MyDirection, DoorOpen: receivedOrder.DoorOpen, 
					Up: receivedOrder.Up, Down: receivedOrder.Down, Inside: receivedOrder.Inside, 
					Defekt: def} 
					
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
			}
			
					
		} else {
			Alive <- IP
		}
			
		Sleep(1*Millisecond)
	}
}






<<<<<<< HEAD

//Receives orders from other elevators
func receiveOrder(receivedOrder Order, IP string) {
=======
//Receives orders from all elevators
func receiveOrder(receivedOrder Order) {
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
	
	go SetButtonLight(receivedOrder)
	
<<<<<<< HEAD
	if receivedOrder.OrderHandled {
		UpdateMyOrders(receivedOrder, IP)
		return
	}

	if receivedOrder.UpdateOrder && IP != MyAddress {
		UpdateMyOrders(receivedOrder, IP)
=======
	if receivedOrder.OrderHandled {		//sletter ordre
		UpdateMyOrders(receivedOrder)
		UpdateGlobalOrders(receivedOrder)
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
		return
	}
	
	if receivedOrder.NewOrder {
		//sjekker om jeg er i rett etg
		if (MyFloor == receivedOrder.Floor) && (DoorOpen || MyDirection == -1 ) {
			OpenDoor <- receivedOrder.Direction
			return
		} 

		//sjekker om noen av de andre heisene er i rett etg
		for _, val := range Elevators {
			if (receivedOrder.Floor == val.LastFloor) && (val.DoorOpen || val.Direction == -1) {
				if !val.Defekt {
					return
				}
			}
		}
	}
	
	UpdateGlobalOrders(receivedOrder)
	if !Defekt && GetCost(LastFloor, MyDirection, receivedOrder.Floor, receivedOrder.Direction, MyAddress) == 1 {
		if EmptyQueue() {
<<<<<<< HEAD
			UpdateMyOrders(receivedOrder, "")
			setDirection()
=======
			UpdateMyOrders(receivedOrder)
			SetDirectionToOrder(receivedOrder.Direction)
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
		} else {
			UpdateMyOrders(receivedOrder, "")
		}
	} 
}






// go fra main. sender hvert sekund oppdatering på floor og direction
func SendUpdateMessage() {
	for {
<<<<<<< HEAD
		order := Order{myFloor, myDirection, -1, -1, false, false, false, Up, Down, Inside}
=======
		order := Order{LastFloor, MyDirection, -1, -1, false, false, DoorOpen, Up, Down, Inside}
>>>>>>> 0beabc0fd3bf2b289636ed18cfb04b7c1b1ec22f
		b, err := json.Marshal(order)
		
		if (err != nil) {
			println("Send Order Error: ", err)
		}
		
		var message Udp_message
		message.Raddr = "broadcast"
		message.Data = b
		message.Length = 1024
		
		
		Send_ch <- message
		Sleep(200*Millisecond)
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




