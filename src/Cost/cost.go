package Cost

import (
	"strconv"
	."./../Driver"
	."./../OrderRegister"
)

type ElevStatus struct{
	LastFloor int
	Direction int
	Up [N_FLOORS]bool
	Down [N_FLOORS]bool
	Inside [N_FLOORS]bool
}


var Elevators = make(map[string]ElevStatus)
var MyAddress string


func calculateCost(loopUp, loopDown bool, highestFloor, lowestFloor, pos int) int {

	cost := 0
	highestOrder := -1
	lowestOrder := N_FLOORS
	
	cost += highestFloor - lowestFloor
	
	for i:=lowestFloor; i<highestFloor; i++ {
		if pos == 1 && (Up[i] || Inside[i]) {
			cost += 4
		} else if pos == 0 && (Down[i] || Inside[i]) {
			cost += 4
		}
	}
	
	if loopUp {
		for i:=highestFloor; i<N_FLOORS; i++ {
			if Up[i] || Inside[i] {
				cost += 4
				if i > highestOrder {
					highestOrder = i
				}
			}
			if Down[i] {
				cost += 4
				if i > highestOrder {
					highestOrder = i
				}
			}
		}
		if highestOrder != -1 {
			cost += (highestOrder - highestFloor)*2
		}
	}
	if loopDown {
		for i:=0; i<lowestFloor; i++ {
			if Up[i] {
				cost += 4
				if i < lowestOrder {
					lowestOrder = i
				}
			}
			if Down[i] || Inside[i] {
				cost += 4
				if i < lowestOrder {
					lowestOrder = i
				}
			}
		}
		if lowestOrder != N_FLOORS {
			cost += (lowestFloor - lowestOrder)*2
		}
	}
	return cost
}




func selectCostCase(myFloor, orderFloor, myDirection, orderDirection int) int {
	cost := 0
	
	if orderFloor > myFloor {
		if myDirection == -1 {
			cost += calculateCost(false, false, orderFloor, myFloor, 1)
		} else if orderDirection == 1 {
			if myDirection == 1 {
				cost += calculateCost(false, false, orderFloor, myFloor, 1)
			} else if myDirection == 0 {
				cost += calculateCost(false, true, orderFloor, myFloor, 1) 
			}
		} else if orderDirection == 0 {
			if myDirection == 1 {
				cost += calculateCost(true, false, orderFloor, myFloor, 1)
			} else if myDirection == 0 {
				cost += calculateCost(true, true, orderFloor, myFloor, 1) 
			}
		}
	
	} else if orderFloor < myFloor {
		if myDirection == -1 {
			cost += calculateCost(false, false, myFloor, orderFloor, 0)
		} else if orderDirection == 1 {
			if myDirection == 1 {
				cost += calculateCost(true, true, myFloor, orderFloor, 0)
			} else if myDirection == 0 {
				cost += calculateCost(false, true, myFloor, orderFloor, 0) 
			}
		} else if orderDirection == 0 {
			if myDirection == 1 {
				cost += calculateCost(true, false, myFloor, orderFloor, 0)
			} else if myDirection == 0 {
				cost += calculateCost(false, false, myFloor, orderFloor, 0) 
			}
		}
	}
	return cost
}






//TODO: Gange opp cost med 1000 og legge til IP, slik at cost aldri er lik
//Nå beregnes costen for mange ganger på høyest IP (for some reason), så flere heiser tar samme best.
func GetCost(myFloor, myDirection, orderFloor, orderDirection int) int {
	
	//Find my cost:
	myCost := selectCostCase(myFloor, orderFloor, myDirection, orderDirection)
	
	myIP, _ := strconv.Atoi(MyAddress)
	myCost = (myCost*1000)+myIP
	
	//Check if other elevator got lower cost:
	for IP, val := range Elevators {
	
		elevCost := selectCostCase(val.LastFloor, orderFloor, val.Direction, orderDirection)
		
		elevIP, _ := strconv.Atoi(IP)
		elevCost = (elevCost*1000)+ elevIP
		
		println("MyCost = ", myCost, " ElevCost = ", elevCost, " Direction = ", orderDirection)
		
		if elevCost < myCost {
			return 0
		}
	}

	return 1

}
