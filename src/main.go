package main

import(
	."./Elevator"
	."./Driver"
)


//TODO: 
// - Lage kostfunksjon
// - Gjøre slik at de kan dra i heisen
// - Send Order i loop med alle bestillinger om heis faller ut av nettverket
// - Fikse buttonlys 





func main() {
	
	//Initialiser heis
	Init()
	
	/*
	
	Lytter hele tiden etter newOrder
	Mottar newOrder på egen heis
		om bestillingen er i samme etg som heisen er, ta bestillingen selv
		ellers:
			sender order, med kost til de andre heisene
			om ikke mottatt svar etter 1 sec, ta bestillingen selv
			om mottar svar fra annen heis, ikke ta bestillingen
	Mottar newOrder fra annen heis
		Sjekker sin egen kost opp mot den andres cost
		Sender svar tilbake dersom lavere kost
	
	*/
	
	
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







