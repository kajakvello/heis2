package Elevator

import (
	//."./../orderRegister"
	."./../Driver"
	."time"
	."./../Udp"
	"encoding/json"
	//"fmt"
)


var floor = -1
var last_floor = 0
var direction = -1	// -1 = står i ro, 1 = opp, 0 = ned
var doorOpen = false


var receive_ch chan Udp_message
var exit chan bool


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
	direction = 1
	last_floor = -1
	floor = 0
	println("ferdig init")
}




func floorReached(floor int) {
	last_floor = floor
	Elev_set_floor_indicator(floor)		//set light on floor
	
	if (GetOrder(direction, floor)) {
		if direction == 1 {
			Elev_set_motor_direction(-100)
		} else if (direction == 0) {
			Elev_set_motor_direction(100)
		}
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		if doorOpen {
			exit <- true
		}
		go openDoorAndDeleteOrder()
		
	} else if (floor == 0) {			//Stops, so the elevator do not pass 1. floor
		Elev_set_motor_direction(100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		direction = 1
		
	} else if (floor == 3) {			//Stops, so the elevator do not pass 4. floor
		Elev_set_motor_direction(-100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		direction = 0
	}
	
}



func RunElevator() {

	for {
		if doorOpen {
			Sleep(100*Millisecond)
		} else {
		
			setDirection()

			if direction == 0 {
				Elev_set_motor_direction(-300)
			} else if direction == 1 {
				Elev_set_motor_direction(300)
			}

			Sleep(100*Millisecond)
		}
	}
}



//Calculates cost
func getCost(orderFloor int, direction int) int {
	//Hente ut info om alle andre heiser og regne ut alle tre koster!
	cost := 1//KOSTFUNKSJON
	return cost
}






//Registers if any up-buttons is pushed
func CheckButtonCallUp() {
	for{
		for i:=0; i<N_FLOORS-1; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_UP, i)) {
				println("er i første if")
				if (direction == -1 && floor == i) {
					if doorOpen {
						exit <- true
					}
					go openDoorAndDeleteOrder()
				
				} else {
					//Regn ut egen cost og send newOrder
					//getCost(i, 1)
					newOrder := Order{floor, direction, i, 1, false}
					println(i)
					UpdateGlobalOrders(newOrder)
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
			
				if (direction == -1 && floor == i) {
					if doorOpen {
						exit <- true
					}
					go openDoorAndDeleteOrder()

				} else {
					//Regn ut egen cost og send newOrder
					//getCost(i, 0)
					newOrder := Order{floor, direction, i, 0, false}
					UpdateGlobalOrders(newOrder)
					go SendOrder(newOrder)
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
				if (direction == -1 && floor == i) {
					if doorOpen {
						exit <- true
					}
					go openDoorAndDeleteOrder()

				} else {
					newOrder := Order{floor, direction, i, -1, false}
					UpdateMyOrders(newOrder)
				}
			}
		}
		Sleep(100*Millisecond)
	}
}


func UpdateFloor() {
	for{
		floor = Elev_get_floor_sensor_signal()
		
		if last_floor != floor {	
		    if (floor != -1) {
		        floorReached(floor)
		    } else {
		    	Elev_set_door_open_lamp(0)
		    }
		    last_floor = floor
		}
		Sleep(100*Millisecond)
	}
}


func openDoorAndDeleteOrder() {
	
	select {
	case ok := <- exit:
    	if ok {
    		println("kill thread")
    		return
    	}
	default:
		doorOpen = true
		Elev_set_door_open_lamp(1)
		Sleep(3*Second)
		
		
		if Elev_get_floor_sensor_signal() == last_floor {
			deleteOrder := Order{floor, direction, floor, -1, true}
			UpdateMyOrders(deleteOrder)
			UpdateGlobalOrders(deleteOrder)
		}
		
		Elev_set_door_open_lamp(0)
		doorOpen = false
	}
}


func setDirection(){

	if (EmptyQueue()) {
		direction = -1
	} else {

		if (direction == 0) && !(CheckOrdersUnderFloor(last_floor)) {
			direction = 1
		} else if (direction == 1) && !(CheckOrdersAboveFloor(last_floor)) {
			direction = 0
		} else if direction == -1 {
			if CheckOrdersAboveFloor(last_floor) {
				direction = 1
			} else if CheckOrdersUnderFloor(last_floor) {
				direction = 0
			}
		}
	}
}



//Receives orders from other elevators
func ReceiveOrder() Order {

	var receivedMessage Udp_message
	receivedMessage = <- receive_ch
	
	var receivedOrder Order
	
	err := json.Unmarshal(receivedMessage.Data, &receivedOrder)
	
	if (err != nil) {
		println("Receive Order Error: ", err)
	}
	return receivedOrder
}





func Stop(ch chan int) int {
	for {
		if Elev_get_stop_signal() != 0 {
			ch <- 1
		}
	}
}



