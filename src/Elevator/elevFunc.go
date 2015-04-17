package Elevator

import (
	//."./../orderRegister"
	."./../Driver"
	."time"
	."./../Udp"
	"encoding/json"
	"fmt"
	"net"
	."strings"
	"strconv"
)

type ElevStatus struct{
	LastFloor int
	Direction int
}


var myFloor = -1
var lastFloor = 0
var myDirection = -1	// -1 = står i ro, 1 = opp, 0 = ned
var doorOpen = false
var myAddress string
var elevators = make(map[string]ElevStatus)


var receive_ch chan Udp_message
var send_ch chan Udp_message
var openDoor = make(chan bool)


//Elevfunc skal ha initfunksjon, alle elevfunksjoner og de fleste variabler, troooor jeg

func Init(localPort, broadcastPort, message_size int) {
	
	err := Udp_init(localPort, broadcastPort, message_size, Send_ch, receive_ch)
	if err != nil {
		println("Error during udp-init")
		return
	}
	
	Elev_init()		//Fra driver.go
	DeleteAllOrders()
	for Elev_get_floor_sensor_signal() != 0 {
		Elev_set_motor_direction(-300)
	}
	Elev_set_motor_direction(50)
	Sleep(2000*Microsecond)
	Elev_set_motor_direction(0)
	Elev_set_floor_indicator(0)	
	myDirection = -1
	lastFloor = 0
	myFloor = 0
	
	//Henter egen ip-adresse = 147
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				splitip := Split(ip, ".")
				myAddress := splitip[3]
			}
		}
	}
	println("Init completed")
}






func RunElevator() {

	for {
		if doorOpen {
			Sleep(100*Millisecond)
		} else {
			if (EmptyQueue()){
				setDirection()
			}

			if myDirection == 0 {
				Elev_set_motor_direction(-300)
			} else if myDirection == 1 {
				Elev_set_motor_direction(300)
			}

			Sleep(100*Millisecond)
		}
	}
}




func UpdateFloor() {
	for{
		myFloor = Elev_get_floor_sensor_signal()
		
		if lastFloor != myFloor {	
		    if (myFloor != -1) {
		        floorReached(myFloor)
		    } else {
		    	Elev_set_door_open_lamp(0)
		    }
		    lastFloor = myFloor
		}
		Sleep(100*Millisecond)
	}
}




func floorReached(floor int) {
	lastFloor = floor
	Elev_set_floor_indicator(floor)
	
	
	//TODO: Lurer på om GetORder kan sjekke globalliste, så heisen stopper om den når fram før noen andre.
	if (GetOrder(myDirection, floor)) {		//Stops, if orders on floor
		if myDirection == 1 {
			Elev_set_motor_direction(-100)
		} else if (myDirection == 0) {
			Elev_set_motor_direction(100)
		}
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
	
		openDoor <- true
		
	} else if (floor == 0) {			//Stops, so the elevator do not pass 1. floor
		Elev_set_motor_direction(100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		myDirection = 1
		
	} else if (floor == 3) {			//Stops, so the elevator do not pass 4. floor
		Elev_set_motor_direction(-100)
		Sleep(2000*Microsecond)
		Elev_set_motor_direction(0)
		myDirection = 0
	}
	
}




//Registers if any up-buttons is pushed
func CheckButtonCallUp() {
	
	for{
		for i:=0; i<N_FLOORS-1; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_UP, i)) {
				
				if (myDirection == -1 && myFloor == i) || (doorOpen && myFloor == i) {
					openDoor <- true
				} else {
					newOrder := Order{myFloor, myDirection, i, 1, false, true}
					if getCost(i, 1) == 1 {
						if EmptyQueue() {
							UpdateMyOrders(newOrder)
							setDirection()
						} else {
							UpdateMyOrders(newOrder)
						}
					}
					go sendOrder(newOrder)
					UpdateGlobalOrders(newOrder)
				}
			}
		}
		Sleep(100*Millisecond)
	}
}




//Registers if any down-buttons is pushed
func CheckButtonCallDown() {

	for{
		for i:=1; i< N_FLOORS; i++ {
			if (Elev_get_button_signal(BUTTON_CALL_DOWN, i)) {
				
				if (myDirection == -1 && myFloor == i) || (doorOpen && myFloor == i) {
					openDoor <- true
				} else {
					newOrder := Order{myFloor, myDirection, i, 1, false, true}
					if getCost(i, 1) == 1 {
						if EmptyQueue() {
							UpdateMyOrders(newOrder)
							setDirection()
						} else {
							UpdateMyOrders(newOrder)
						}
					}
					go sendOrder(newOrder)
					UpdateGlobalOrders(newOrder)
				}
			}
		}
		Sleep(100*Millisecond)
	}
}




//Registers if any command-buttons is pushed
func CheckButtonCommand() {

	for{
		for i:=0; i<N_FLOORS; i++ {
			if (Elev_get_button_signal(BUTTON_COMMAND, i)) {
			
				if (myDirection == -1 && myFloor == i) || (doorOpen && myFloor == i) {
					openDoor <- true
				} else {
					newOrder := Order{myFloor, myDirection, i, -1, false, true}
					if EmptyQueue() {
						UpdateMyOrders(newOrder)
						setDirection()
					} else {
						UpdateMyOrders(newOrder)
					}
				}
			}
		}
		Sleep(100*Millisecond)
	}
}




func DoorControl() {

	timer := NewTimer(Hour*3)
	for{
	
		select {
			case <- openDoor:
				doorOpen = true
				Elev_set_door_open_lamp(1)
				timer.Reset(Second*3)
				if Elev_get_floor_sensor_signal() == lastFloor {
					deleteOrder := Order{myFloor, myDirection, myFloor, -1, true, false}
					UpdateMyOrders(deleteOrder)
					UpdateGlobalOrders(deleteOrder)
					go sendOrder(deleteOrder)
				}
				
			case <- timer.C:
				println("timer out")
				Elev_set_door_open_lamp(0)
				doorOpen = false
				setDirection()
		}
	}
}




func setDirection(){

	if (EmptyQueue()) {
		myDirection = -1
		
	} else {
		if (myDirection == 0) && !(CheckOrdersUnderFloor(lastFloor)) {
			myDirection = 1
		} else if (myDirection == 1) && !(CheckOrdersAboveFloor(lastFloor)) {
			myDirection = 0
		} else if myDirection == -1 {
			if CheckOrdersAboveFloor(lastFloor) {
				myDirection = 1
			} else if CheckOrdersUnderFloor(lastFloor) {
				myDirection = 0
			}
		}
	}
}






//Calculates cost, returns 1 if myElev got the lowest cost
func getCost(orderFloor int, orderDirection int) int {
	
	myCost := 1 //regner ut egen cost
	lowestCost := myCost
	equalCost := 0
	elevKey := ""
	
	for key, val := range elevators {
		elevCost := 1 			//costfunksjon ut fra val.LastFloor og val.Direction
		if elevCost < lowestCost {
			lowestCost = elevCost
			return 0
		} else if elevCost == myCost {
			equalCost = elevCost
			elevKey = key
		}
	}

	
	if (myCost == lowestCost) && (equalCost == 0) {
		return 1
	} else if (myCost == equalCost) && (myCost == lowestCost) {
		myAddr, _ := strconv.Atoi(myAdress)
		elevAddr, _ := strconv.Atoi(elevators[elevKey])
		if myAddr < elevAddr {
			return 1
		}
		return 0
	} 
	return 0
}



//Receives messages from other elevators continuous
func ReceiveMessage() {
	
	for{
		var receivedMessage Udp_message
		receivedMessage = <- receive_ch
		
		IP := getIP(receivedMessage.Raddr)
		
		if IP == MyAddress {
			break
		}
		
		var receivedOrder Order
		err := json.Unmarshal(receivedMessage.Data[:receivedMessage.Length], &receivedOrder)
		if (err != nil) {
			println("Receive Order Error: ", err)
		}

		
		if receivedOrder.newOrder {
			receiveOrder(receivedOrder)
		}
		
		
		newElevator := true	
		for key,_ := range elevators {
			if key == IP {
				newElevator = false
			}
		}
		
		if newElevator {
			go setMessageTimer(IP)
		} else {
			gotMessage <- IP
		}
		
		elevators[IP] = ElevStatus{LastFloor: receivedMessage.MyFloor, Direction: receivedMessage.MyDirection} 	
		
	}
}



//Returns last three numbers of IP-address
func getIP(address string) string {
	splitaddr := Split(address, ".")
	splitip := Split(splitaddr[3], ":")
	myAddress := splitip[0]
	return myAddress
}





func setMessageTimer(address string) {
	
	timer := NewTimer(3*Hour)
	for {
		select {
		case <- timer.C:
			delete (elevators, address)
			return
		case receivedAddress := <- gotMessage:
			if receivedAddress == address {
				timer.Reset(3*Second)
			}
		}
	}
}







//Receives orders from other elevators
func receiveOrder(receivedOrder Order) {
	
	cost := getCost(receivedOrder.Floor, receivedOrder.Direction)
	
	if cost == 1 {
		UpdateMyOrders(receivedOrder)
	}
	UpdateGlobalOrders(receivedOrder)

}




func sendOrder(order Order) {
	b, err := json.Marshal(order)
	
	if (err != nil) {
		println("Send Order Error: ", err)
	}
	
	var message Udp_message
	message.Raddr = "broadcast"
	message.Data = b
	message.Length = 1024
	
	Send_ch <- message
}



// go fra main. sender hvert sekund oppdatering på floor og direction
func SendUpdateOrder(Send_ch chan Udp_message) {
	for {
		order := Order{myFloor, myDirection, -1, -1, false, false}
		b, err := json.Marshal(order)
		
		if (err != nil) {
			println("Send Order Error: ", err)
		}
		
		var message Udp_message
		message.Raddr = "broadcast"
		message.Data = b
		message.Length = 1024
		
		Send_ch <- message
	}
	Sleep(1*Second)
}




func Stop(ch chan int) {
	for {
		if Elev_get_stop_signal() != 0 {
			ch <- 1
		}
		Sleep(100*Millisecond)
	}
}



