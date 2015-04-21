package main

import(
	."./Elevator"
	."./Driver"
)


//TODO: 
// - Lage kostfunksjon
// - Gjøre slik at de kan dra i heisen
// - Send Order i loop med alle bestillinger om heis faller ut av nettverket





func main() {
	
	//Initialiser heis
	Init()
	
	
	
	go CheckButtonCallUp()
	go CheckButtonCallDown()
	go CheckButtonCommand()
	go RunElevator()
	go UpdateFloor()
	go DoorControl()

	
	go ReceiveMessage()
	go SendUpdateMessage()
	

	s := make(chan int)
	go Stop(s)
	
	select {
	case <- s:
		Elev_set_motor_direction(0)
		break
	} 
}







