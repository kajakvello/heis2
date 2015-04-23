package main

import(
	."./Elevator"
	."./Driver"
)


//TODO: 

// - Heisene regner av og til ut forskjellig cost på samme heis. Fml.

// - Fiks så man kan nappe strømmen til motor
// - Fiks så de kan drepe prossessen din, men interne ordre ikke går tapt





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
	
	//go PrintStatus()
	

	s := make(chan int)
	go Stop(s)
	
	select {
	case <- s:
		Elev_set_motor_direction(0)
		break
	} 
	
}







