package main

import(

	."./Elevator"
	."./Driver"
	//"timer"
	//"fmt"
	//"net"
	//."time"
	//"strings"
	//"strconv"
)

const localPort = 20016
const broadcastPort = 20017
const message_size = 1024


//TODO: 
// - Lage kostfunksjon
// - Skal elevators ta i mot seg selv og ha i map?





func main(){	
	
	//Initialiser heis
	Init(localPort, broadcastPort, message_size)
	
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
	
	//TODO: go inni for eller for inni go?
	
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







