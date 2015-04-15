package Elevator

import (
	//."./../orderRegister"
	."./../Driver"
	."time"
	."./../Udp"
	"encoding/json"
	//"fmt"
)


var myFloor = -1
var lastFloor = 0
var myDirection = -1	// -1 = står i ro, 1 = opp, 0 = ned
var doorOpen = false


var receive_ch chan Udp_message
var openDoor = make(chan bool)


//Elevfunc skal ha initfunksjon, alle elevfunksjoner og de fleste variabler, troooor jeg

func Init(localPort, broadcastPort, message_size int) {
	/*
	err := Udp_init(localPort, broadcastPort, message_size, Send_ch, receive_ch)
	if err != nil {
		println("Error during udp-init")
		return
	} */
	
	Elev_init()		//Fra driver.go
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
	println("ferdig init")
}






func RunElevator() {

	for {
		if doorOpen {
			Sleep(100*Millisecond)
		} else {
			if (EmptyQueue()){
				setDirection()
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
		    lastFloor = myFloor
		}
		Sleep(100*Millisecond)
	}
}




func floorReached(floor int) {
	lastFloor = floor
	Elev_set_floor_indicator(floor)		//set light on floor
	
	if (GetOrder(myDirection, floor)) {	//Stops, if orders on floor
		if myDirection == 1 {
			Elev_set_motor_direction(-100)
		} else if (myDirection == 0) {
			Elev_set_motor_direction(100)
		}
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
	
		openDoor <- true
		
	} else if (floor == 0) {			//Stops, so the elevator do not pass 1. floor
		Elev_set_motor_direction(100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		myDirection = 1
		
	} else if (floor == 3) {			//Stops, so the elevator do not pass 4. floor
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
					//Regn ut egen cost og send newOrder
					//getCost(i, 1)
					newOrder := Order{myFloor, myDirection, i, 1, false}
					if EmptyQueue() {
						UpdateMyOrders(newOrder)
						setDirection()
					}
					UpdateGlobalOrders(newOrder)
					UpdateMyOrders(newOrder)		//for testing
					//go SendOrder(newOrder)
					//Set en timer som hører etter svar, ta bestillingen selv om ingen svar etter timer går ut.
				}
			}
		}
		Sleep(100*Millisecond)
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
					//Regn ut egen cost og send newOrder
					//getCost(i, 0)
					newOrder := Order{myFloor, myDirection, i, 0, false}
					if EmptyQueue() {
						UpdateMyOrders(newOrder)
						setDirection()
					}
					UpdateMyOrders(newOrder)
					UpdateGlobalOrders(newOrder)
					//go SendOrder(newOrder)
					//Set en timer som hører etter svar, ta bestillingen selv om ingen svar etter timer går ut.
				}
			}
		}
		Sleep(100*Millisecond)
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
					newOrder := Order{myFloor, myDirection, i, -1, false}
					if EmptyQueue() {
						UpdateMyOrders(newOrder)
						setDirection()
					}
					UpdateMyOrders(newOrder)
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
					deleteOrder := Order{myFloor, myDirection, myFloor, -1, true}
					UpdateMyOrders(deleteOrder)
					UpdateGlobalOrders(deleteOrder)
				}
				
			case <- timer.C:
				println("timer out")
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







//Calculates cost, returns 1 if myElev got the lowest cost
func getCost(orderFloor int, orderDirection int) int {

	myCost := 1 //regner ut egen cost
	elevOneCost := 1 //reger ut elev1 sin cost ut fra ElevatorPositions.ElevOneFloor og ElevatorPositions.ElevOneDirection
	elevTwoCost := 1 //reger ut elev2 sin cost ut fra ElevatorPositions.ElevTwoFloor og ElevatorPositions.ElevTwoDirection
	if (myCost < elevOneCost) && (myCost < elevTwoCost) {
		return 1
	} 
	return 0
}








//Receives orders from other elevators
func ReceiveOrder() {

	var receivedMessage Udp_message
	receivedMessage = <- receive_ch
	
	var receivedOrder Order
	
	err := json.Unmarshal(receivedMessage.Data, &receivedOrder)
	
	if (err != nil) {
		println("Receive Order Error: ", err)
	}
	
	//Init messages from the other elevators
	if ElevOneAdress == nil {
		ElevOneAdress = receivedMessage.Raddr
	} else if ElevTwoAdress == nil {
		ElevTwoAdress = receivedMessage.Raddr
	}
	
	//Set other elevators positon
	if receivedMessage.Raddr == "???" {
		ElevOneFloor = receivedOrder.MyFloor
		ElevOneDirection = receivedOrder.MyDirection
	} else if receivedMessage.Raddr == "???2" {
		ElevTwoFloor = receivedOrder.MyFloor
		ElevTwoDirection = receivedOrder.MyDirection
	}
	
	//regn ut kost
	//legg til i globalOrders
		//evt også i myOrders
		
	
}





func Stop(ch chan int) {
	for {
		if Elev_get_stop_signal() != 0 {
			ch <- 1
		}
		Sleep(100*Millisecond)
	}
}



