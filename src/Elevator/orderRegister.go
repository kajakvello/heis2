package Elevator

import (
	."./../Driver"
	"encoding/json"
	."./../Udp"
	."time"
)



//Lager en liste for antall heiser, men adresse, posisjon og retning.
//Om en heis faller ut (vi ikke får meldinger etter viss tid), settes floor og direction til -2.
//Om en bestilling ikke har blitt tatt, regner gjenværende heiser ut ny cost seg i mellom.



//ElevatorPositions:
//Sets from senders adress where the other elevators are located


var orderHandled chan int		//sletter timer dersom order er tatt








//0 = 1.etg opp, 1 = 2.etg opp, 2 = 3.etg opp, 3 = 2.etg ned, 4 = 3.etg ned, 5 = 4.etg ned
var globalOrders [N_FLOORS*2-2]bool
//Alle bestillinger for alle heiser. Må sjekkes jevnlig at bestillinger blir ekspedert. Aner ikke hvordan, shalalalala, aner ikke hvordan shaaa-la-la-la-la-la!


// My inside orders
var inside [N_FLOORS]bool


// My up orders
var up [N_FLOORS]bool


// My down orders
var down [N_FLOORS]bool






// Update my orders
func UpdateMyOrders(receivedOrder Order) {

	if receivedOrder.OrderHandledAtFloor {
		
		inside[receivedOrder.Floor] = false
		up[receivedOrder.Floor] = false
		down[receivedOrder.Floor] = false
		
		Elev_set_button_lamp(BUTTON_COMMAND, receivedOrder.Floor, 0)
		if (receivedOrder.Floor < N_FLOORS-1) {
			Elev_set_button_lamp(BUTTON_CALL_UP, receivedOrder.Floor, 0)
		}
		if (receivedOrder.Floor > 0) {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, receivedOrder.Floor, 0)
		}

		
	} else if receviedOrder.NewOrder {
	
		if receivedOrder.Direction == 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, receivedOrder.Floor, 1)
			down[receivedOrder.Floor] = true
		} else if receivedOrder.Direction == 1 {
			up[receivedOrder.Floor] = true
			Elev_set_button_lamp(BUTTON_CALL_UP, receivedOrder.Floor, 1)
		} else if receivedOrder.Direction == -1 {
			inside[receivedOrder.Floor] = true
			Elev_set_button_lamp(BUTTON_COMMAND, receivedOrder.Floor, 1)
		} else {
			println("Unvalid direction, or unvalid floor")
		}	
		
	} else {
		println("Error in update my order")
	}

}




// Runs everytime the program receives a new order
func UpdateGlobalOrders(receivedOrder Order) {

	if globalOrders[receivedOrder.Floor] == false {
		setGlobalOrderTimer(receivedOrder)
	}
	
	if receivedOrder.OrderHandledAtFloor {
		
		globalOrders[receivedOrder.Floor] = false
		globalOrders[N_FLOORS-2 + receivedOrder.Floor] = false
		orderHandled <- receivedOrder.Floor
		
		if (receivedOrder.Floor < N_FLOORS-1) {
			Elev_set_button_lamp(BUTTON_CALL_UP, receivedOrder.Floor, 0)
		}
		if (receivedOrder.Floor > 0) {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, receivedOrder.Floor, 0)
		}
		
	} else if receivedOrder.NewOrder {
	
		if receivedOrder.Direction == 1 {
			globalOrders[receivedOrder.Floor] = true
			Elev_set_button_lamp(BUTTON_CALL_UP, receivedOrder.Floor, 1)
			setOrder <- receivedOrder.Floor
		} else if receivedOrder.Direction == 0 {
			globalOrders[N_FLOORS-2 + receivedOrder.Floor] = true
			Elev_set_button_lamp(BUTTON_CALL_DOWN, receivedOrder.Floor, 1)
			setOrder <- N_FLOORS-2 + receivedOrder.Floor
		} else {
			println("Not valid direction, or unvalid floor")
		}
		
	} else {
		println("Error in update global order")
	}

}





//Funker fra init
func DeleteAllOrders() {
	for j:=0; j<N_FLOORS*2-2; j++ {
		globalOrders[j] = false
		orderHandled <- j
	}

	for j:=0; j<N_FLOORS; j++ {
		inside[j] = false
	}

	for j:=0; j<N_FLOORS; j++ {
		up[j] = false
	}

	for j:=0; j<N_FLOORS; j++ {
		down[j] = false
	}
}




//TODO: Lurer på om denne kan sjekke global liste i stede, så den tar bestillingen om den plutselig når fram først.
// Returns true if the elevator should take an order from "floor". If it exists an order in the same direction as the elevator is headed.
func GetOrder(direction int, floor int) bool {
	/*
	if inside[floor] {
		return true	
	}
	if (globalOrders[floor]) && (direction == 1 || direction == -1 || floor == 0 || !CheckOrdersUnderFloor(floor)) {
		return true	
	}
	if (globalOrders[N_FLOORS-2+floor] && (direction == 0 || direction == -1 || floor == 3 || !CheckOrdersAboveFloor(floor)) {
		return true
	}
	return false
	*/

	if inside[floor] {
		return true
	}
	if up[floor] && (direction == 1 || direction == -1 || floor == 0 || !CheckOrdersUnderFloor(floor)) {
		return true
	}
	if down[floor] && (direction == 0 || direction == -1 || floor == 3 || !CheckOrdersAboveFloor(floor)) {
		return true
	}
	return false
}




func CheckOrdersUnderFloor(floor int) bool {
	for i:=0; i<floor; i++ {
		if (up[i] || down[i] || inside[i]) {
			return true
		}
	}
	return false
}




func CheckOrdersAboveFloor(floor int) bool {
	for i:=floor+1; i<N_FLOORS; i++ {
		if (up[i] || down[i] || inside[i]) {
			return true
		}
	}
	return false
}




func EmptyQueue() bool {
	for i:=0; i<N_FLOORS; i++ {
		if (up[i] || down[i] || inside[i]) {
			return false
		}
	}
	return true
}






func setGlobalOrderTimer(order Order) {
	
	timer := NewTimer(10*Second)
	for {
		select {
		case <- timer.C:
			println("Order not handled")
			UpdateMyOrders(order)
			go setGlobalOrderTimer(order)
			return
		case i := <- orderHandled:
			if i == order.Floor {
				timer.Stop()				//Endret fra Reset() til Stop()
				return
			}
		}
	}
}













