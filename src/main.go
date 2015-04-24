package main

import(
	."./Elevator"
	."./Driver"
	."./Timer"
)


//TODO: 


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
	
	go PrintStatus()
	

	s := make(chan int)
	go Stop(s)
	
	select {
	case <- s:
		Elev_set_motor_direction(0)
		break
	} 
	
	
}







