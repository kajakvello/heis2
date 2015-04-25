package Timer


import (
	."./../Driver"
	."./../Cost"
	."./../OrderRegister"
	."time"
	//."fmt"
	
)




//Opens door for three seconds. Deletes lights and order when doors open.
func DoorControl() {

	timer := NewTimer(Hour*3)
	orderDir := 2
	for {
	
		select {
		case orderDir = <- OpenDoor:
			DoorOpen = true
			Elev_set_door_open_lamp(1)
			timer.Reset(Second*3)
			if Elev_get_floor_sensor_signal() == LastFloor {
				deleteOrder := Order{LastFloor, MyDirection, LastFloor, orderDir, true, false, DoorOpen, Up, Down, Inside}
				SendOrder(deleteOrder)
			}
			
		case <- timer.C:
			Elev_set_door_open_lamp(0)
			DoorOpen = false
			if orderDir == 1{
				orderDir = 0
			} else if orderDir == 0 {
				orderDir = 1
			}
			go SetDirectionToOrder(orderDir)
		}
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



func SetMessageTimer(address string) {
	
	timer := NewTimer(3*Hour)
	for {
		select {
		
		case IP := <- GotMessage:
			if IP == address {
				timer.Reset(3*Second)
			}
		
		case <- timer.C:
			temp := Elevators[address]
			temp.Defect = true
			Elevators[address] = temp
			
			for i:=0; i<N_FLOORS; i++ {
				if (Elevators[address].Up)[i] {
					order := Order{LastFloor, MyDirection, i, 1, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(order)
				}
				if (Elevators[address].Down)[i] {
					order := Order{LastFloor, MyDirection, i, 0, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(order)
				}	
			}
			delete (Elevators, address)
			return
			
		}	
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////


/*
func OrderTimer(order int) {

	timer := NewTimer(10*Second)
	
	if order.Direction == 1 {
		floor = order.Floor
	} else if order.Direction == 0 {
		floor = N_FLOORS-2+order.Floor
	}
	
	for {
		select {
		
		case handledFloor := <- orderHandled:
			if handledFloor == floor {
				return
			}
		}
		
		case <- timer.C:
			if order.Direction == 1 {
				Up[order.Floor] = false
			} else if order.Direction == 0 {
				Down[order.Floor] = false
			}
			newOrder := Order{floor, direction, floor, direction, false, true, DoorOpen, Up, Down, Inside}
			go SendOrder(newOrder)
			timer.Reset(10*Second)
			return
			
	}


}
*/


/////////////////////////////////////////////////////////////////////////////////////////////



//Resets timer if order has been handled or if the elevator has no orders (dir == -1).
//Sends orders to other elevators if timer runs out. Deletes all outside orders and sets one order true to check if its running again
func AliveTimer(address string) {

	timer := NewTimer(3*Hour)
	oldUp := [N_FLOORS]bool{}
	oldDown := [N_FLOORS]bool{}
	oldInside := [N_FLOORS]bool{}
	
	for {	
		select {
		case IP := <- Alive: 

			if IP == address {
				temp := Elevators[IP]
				
				if Elevators[IP].Direction == -1 {
					timer.Reset(10*Second)
					break
				}
				for i:=0; i<N_FLOORS; i++ {
					if (oldUp[i] && !(Elevators[IP].Up)[i]) || (oldDown[i] && !(Elevators[IP].Down)[i]) || (oldInside[i] && !(Elevators[IP].Inside)[i]) {
						timer.Reset(10*Second)
						temp.Defect = false
					}
					oldUp[i] = temp.Up[i]
					oldDown[i] = temp.Down[i]
					oldInside[i] = temp.Inside[i]
				}
				Elevators[IP] = temp
			}
		case <- timer.C:
		
			temp := Elevators[address]
			temp.Defect = true
			Elevators[address] = temp
			println("Elevator nr ", address, " is defect")
			
			for i:=0; i<N_FLOORS; i++ {
				if (Elevators[address].Up)[i] {
					order := Order{0, -1, i, 1, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}
				if (Elevators[address].Down)[i] {
					order := Order{0, -1, i, 0, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}
				temp.Up[i] = false
				temp.Down[i] = false
			}
			temp.Up[0] = true
			Elevators[address] = temp
			timer.Reset(3*Hour)
		}
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



func SelfAliveTimer() {	//For å vite om seg selv er Defekt. Kalles i Init()

	timer := NewTimer(3*Hour)
	oldUp := [N_FLOORS]bool{}
	oldDown := [N_FLOORS]bool{}
	oldInside := [N_FLOORS]bool{}
	
	for {
		select {
		case IP := <- Alive: 

			if IP == MyAddress {
				
				if MyDirection == -1 {
					timer.Reset(10*Second)
					break
				}
				for i:=0; i<N_FLOORS; i++ {
					if (oldUp[i] && !Up[i]) || (oldDown[i] && !Down[i]) || (oldInside[i] && !Inside[i]) {
						timer.Reset(10*Second)
						Defect = false
					}
					oldUp[i] = Up[i]
					oldDown[i] = Down[i]
					oldInside[i] = Inside[i]
				}
			}
		case <- timer.C:
		
			Defect = true
			println("I am defect")
			
			for i:=0; i<N_FLOORS; i++ {
				if Up[i] {
					order := Order{0, -1, i, 1, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}
				if Down[i] {
					order := Order{0, -1, i, 0, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}
				Up[i] = false
				Down[i] = false
			}
			Up[0] = true
			timer.Reset(3*Hour)
		}
	}
}



