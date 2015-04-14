package main

import(
	
	//. "./udp"
	//."./Driver"
	. "./Elevator"
	//. "./orderRegister"
	//"timer"
	//."fmt"
	//"net"
	//"time"
	//"strings"
	//"strconv"
)

const localPort = 20016
const broadcastPort = 20017
const message_size = 1024

var exit chan bool

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
	go UpdateFloor()
	
	<- exit
	//go ReceiveOrder()
	//go RunElevator()
	
	//go checklastfloor()????
	//go checktimer()????
	
	//Elev_set_motor_direction(0)
}

