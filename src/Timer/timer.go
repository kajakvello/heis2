package Timer


import (
	."./../Udp"
	."./../Driver"
	."./../Cost"
	."./../OrderRegister"
	."time"
	
)




//Opens door for three seconds. Deletes lights when doors opens, deletes order when doors closing.
func DoorControl() {

	timer := NewTimer(Hour*3)
	for {
	
		select {
		case <- OpenDoor:
			println("åpner dør")
			DoorOpen = true
			Elev_set_door_open_lamp(1)
			timer.Reset(Second*3)
			//delete light, but not order
			if Elev_get_floor_sensor_signal() == LastFloor {
				deleteLight := Order{LastFloor, MyDirection, LastFloor, MyDirection, true, false, DoorOpen, Up, Down, Inside}
				go SendOrder(deleteLight)
			}
			
		case <- timer.C:
			println("slukker lys")
			Elev_set_door_open_lamp(0)
			DoorOpen = false
			//delete order
			if Elev_get_floor_sensor_signal() == LastFloor {
				deleteOrder := Order{LastFloor, MyDirection, LastFloor, MyDirection, true, false, DoorOpen, Up, Down, Inside}
				SendOrder(deleteOrder)	
			}
			go SetDirectionToOrder()
		}
	}
}






func SetMessageTimer(address string) {
	
	timer := NewTimer(3*Hour)
	for {
		select {
		case <- timer.C:
			temp := Elevators[address]
			temp.Defekt = true
			Elevators[address] = temp
			
			for i:=0; i<N_FLOORS; i++ {
				if (Elevators[address].Up)[i] {
					order := Order{LastFloor, MyDirection, i, 1, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}
				if (Elevators[address].Down)[i] {
					order := Order{LastFloor, MyDirection, i, 0, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}	
			}
			delete (Elevators, address)
			return
			
		case receivedAddress := <- GotMessage:
			if receivedAddress == address {
				timer.Reset(3*Second)
			}
		}
	}
}




//Resets timer if order has been handled or if the elevator has no orders (dir == -1).
//Sends orders to other elevators if timer runs out. Deletes all outside orders and sets one order true to check if its running again
func AliveTimer(address string) {

	timer := NewTimer(3*Hour)
	oldUp := [N_FLOORS]bool{}
	oldDown := [N_FLOORS]bool{}
	
	for {	
		select {
		case IP := <- Alive: 
		
			if IP == address {
				temp := Elevators[IP]
				
				if Elevators[IP].Direction == -1 {
					timer.Reset(5*Second)
					break
				}
				for i:=0; i<N_FLOORS; i++ {
					if (oldUp[i] && !(Elevators[IP].Up)[i] || oldDown[i] && !(Elevators[IP].Down)[i]) {
						timer.Reset(5*Second)
					}
					oldUp[i] = temp.Up[i]
					oldDown[i] = temp.Down[i]
				}
				temp.Defekt = false
				Elevators[IP] = temp
			}
		case <- timer.C:
		
			temp := Elevators[address]
			temp.Defekt = true
			Elevators[address] = temp
			
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
		}
	}
}







