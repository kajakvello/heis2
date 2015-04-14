package Elevator

import (
	."./../Driver"
	"encoding/json"
	."./../Udp"
)





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
		
		SendOrder(receivedOrder)
		Elev_set_button_lamp(BUTTON_CALL_UP, receivedOrder.Floor, 0)
		Elev_set_button_lamp(BUTTON_CALL_DOWN, receivedOrder.Floor, 0)
		Elev_set_button_lamp(BUTTON_COMMAND, receivedOrder.Floor, 0)
		
	} else {
	
		if receivedOrder.Direction == 0 {
			down[receivedOrder.Floor] = true
		} else if receivedOrder.Direction == 1 {
			up[receivedOrder.Floor] = true
		} else if receivedOrder.Direction == -1 {
			inside[receivedOrder.Floor] = true
		} else {
			println("Unvalid direction, or unvalid floor")
		}
		
	}
}




// Runs everytime the program receives a new order
func UpdateGlobalOrders(receivedOrder Order) {

	if receivedOrder.OrderHandledAtFloor {
	
		globalOrders[receivedOrder.Floor] = false
		globalOrders[N_FLOORS-2 + receivedOrder.Floor] = false
		
		Elev_set_button_lamp(BUTTON_CALL_UP, receivedOrder.Floor, 0)
		Elev_set_button_lamp(BUTTON_CALL_DOWN, receivedOrder.Floor, 0)
		
	} else {
	
		if receivedOrder.Direction == 1 {
			globalOrders[receivedOrder.Floor] = true
		} else if receivedOrder.Direction == 0 {
			globalOrders[N_FLOORS-2 + receivedOrder.Floor] = true
		} else {
			println("Not valid direction, or unvalid floor")
		}
		
	}
}






func DeleteAllOrders() {
	for j:=0; j<N_FLOORS*2-2; j++ {
		globalOrders[j] = false
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





// Returns true if the elevator should take an order from "floor". If it exists an order in the same direction as the elevator is headed.
func GetOrder(direction int, floor int) bool {

	if (inside[floor] == true) {
		return true
	}
	if ( up[floor] == true && (direction == 0 || direction == -1 || floor == 0 || !checkOrdersUnderFloor(floor)) ) {
		return true
	}
	if ( down[floor] == true && (direction == 1 || direction == -1 || floor == 3 || !checkOrdersAboveFloor(floor)) ) {
		return true
	}
	return false
}




func checkOrdersUnderFloor(floor int) bool {
	for i:=0; i<floor; i++ {
		if (up[i] || down[i] || inside[i]) {
			return true
		}
	}
	return false
}




func checkOrdersAboveFloor(floor int) bool {
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




func SendOrder(order Order) {
	b, err := json.Marshal(order)
	
	var message Udp_message
	message.Data = b
	message.Raddr = "broadcast"
	
	
	if (err != nil) {
		println("Send Order Error: ", err)
	}
	
	Send_ch <- message
	
}






