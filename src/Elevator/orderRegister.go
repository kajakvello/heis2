package Elevator

import (
	."./../Driver"
	."./../Udp"
)



//Lager en liste for antall heiser, men adresse, posisjon og retning.
//Om en heis faller ut (vi ikke får meldinger etter viss tid), settes floor og direction til -2.
//Om en bestilling ikke har blitt tatt, regner gjenværende heiser ut ny cost seg i mellom.



//ElevatorPositions:
//Sets from senders adress where the other elevators are located




// My inside orders
var Inside [N_FLOORS]bool


// My up orders
var Up [N_FLOORS]bool


// My down orders
var Down [N_FLOORS]bool




type ElevStatus struct{
	LastFloor int
	Direction int
	Up [N_FLOORS]bool
	Down [N_FLOORS]bool
	Inside [N_FLOORS]bool
}

var Elevators = make(map[string]ElevStatus)
var MyAddress string





// Update my orders
func UpdateMyOrders(receivedOrder Order, IP string) {


	if receivedOrder.OrderHandled {
		
		if IP == MyAddress {
			Inside[receivedOrder.Floor] = false
			Elev_set_button_lamp(BUTTON_COMMAND, receivedOrder.Floor, 0)
		}
		
		Up[receivedOrder.Floor] = false
		Down[receivedOrder.Floor] = false
		
	} else if receivedOrder.UpdateOrder && IP != MyAddress {
		
		if receivedOrder.Direction == 0 {
			Down[receivedOrder.Floor] = false
			
		} else if receivedOrder.Direction == 1 {
			Up[receivedOrder.Floor] = false
		}
		
	} else if receivedOrder.NewOrder {
		
		if receivedOrder.Direction == 0 {
			Down[receivedOrder.Floor] = true
			
		} else if receivedOrder.Direction == 1 {
			Up[receivedOrder.Floor] = true
			
		} else if receivedOrder.Direction == -1 {
			Inside[receivedOrder.Floor] = true
			Elev_set_button_lamp(BUTTON_COMMAND, receivedOrder.Floor, 1)
			
		} else {
			println("Unvalid direction, or unvalid floor")
		}	
		
	} else {
		println("Error in update my order")
	}

}


func SetButtonLight(order Order) {
	
	if order.NewOrder && order.Direction == 0 {
		Elev_set_button_lamp(BUTTON_CALL_DOWN, order.Floor, 1)
		
	} else if order.NewOrder && order.Direction == 1 {
		Elev_set_button_lamp(BUTTON_CALL_UP, order.Floor, 1)
		
	} else if order.OrderHandled {
		if order.Floor < N_FLOORS-1 {
			Elev_set_button_lamp(BUTTON_CALL_UP, order.Floor, 0)
		}
		if order.Floor > 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, order.Floor, 0)
		}
	}
}





//Funker fra init
func DeleteAllOrders() {

	for j:=0; j<N_FLOORS; j++ {
		Inside[j] = false
	}

	for j:=0; j<N_FLOORS; j++ {
		Up[j] = false
	}

	for j:=0; j<N_FLOORS; j++ {
		Down[j] = false
	}
}




// Returns true if the elevator should take an order from "floor". If it exists an order in the same direction as the elevator is headed.
func GetOrder(direction int, floor int) (bool, string) {
	
	if Inside[floor] {
		return true, MyAddress
	}
	if Up[floor] && (direction == 1 || direction == -1 || floor == 0 || !CheckOrdersUnderFloor(floor)) {
		return true, MyAddress
	}
	if Down[floor] && (direction == 0 || direction == -1 || floor == N_FLOORS-1 || !CheckOrdersAboveFloor(floor)) {
		return true, MyAddress
	}

	for IP, val := range Elevators {
		if ((val.Up)[floor] && direction == 1) || ((val.Down)[floor] && direction == 0) {
			return true, IP
		}
	}

	return false, ""
}




func CheckOrdersUnderFloor(floor int) bool {
	for i:=0; i<floor; i++ {
		if (Up[i] || Down[i] || Inside[i]) {
			return true
		}
	}
	return false
}




func CheckOrdersAboveFloor(floor int) bool {
	for i:=floor+1; i<N_FLOORS; i++ {
		if (Up[i] || Down[i] || Inside[i]) {
			return true
		}
	}
	return false
}




func EmptyQueue() bool {
	for i:=0; i<N_FLOORS; i++ {
		if (Up[i] || Down[i] || Inside[i]) {
			return false
		}
	}
	return true
}












