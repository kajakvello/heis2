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
	DoorOpen bool
	Defekt bool
}

var Elevators = make(map[string]ElevStatus)





func calculateCost(loopUp, loopDown bool, highestFloor, lowestFloor, pos int, up, down, inside [N_FLOORS]bool) int {
	
	cost := 0
	highestOrder := -1
	lowestOrder := N_FLOORS
	
	cost += highestFloor - lowestFloor
	
	for i:=lowestFloor; i<highestFloor; i++ {
		if pos == 1 && (up[i] || inside[i]) {
			cost += 4
		} else if pos == 0 && (down[i] || inside[i]) {
			cost += 4
		}
	}
	
	if loopUp {
		for i:=highestFloor; i<N_FLOORS; i++ {
			if up[i] || inside[i] {
				cost += 4
				if i > highestOrder {
					highestOrder = i
				}
			}
			if down[i] {
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
			if up[i] {
				cost += 4
				if i < lowestOrder {
					lowestOrder = i
				}
			}
			if down[i] || inside[i] {
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




func selectCostCase(myFloor, orderFloor, myDirection, orderDirection int, up, down, inside [N_FLOORS]bool) int {
	cost := 0
	
	
	if orderFloor > myFloor {
		if myDirection == -1 {
			cost += calculateCost(false, false, orderFloor, myFloor, 1, up, down, inside)
		} else if orderDirection == 1 {
			if myDirection == 1 {
				cost += calculateCost(false, false, orderFloor, myFloor, 1, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(false, true, orderFloor, myFloor, 1, up, down, inside) 
			}
		} else if orderDirection == 0 {
			if myDirection == 1 {
				cost += calculateCost(true, false, orderFloor, myFloor, 1, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(true, true, orderFloor, myFloor, 1, up, down, inside) 
			}
		}
	
	} else if orderFloor < myFloor {
		if myDirection == -1 {
			cost += calculateCost(false, false, myFloor, orderFloor, 0, up, down, inside)
		} else if orderDirection == 1 {
			if myDirection == 1 {
				cost += calculateCost(true, true, myFloor, orderFloor, 0, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(false, true, myFloor, orderFloor, 0, up, down, inside) 
			}
		} else if orderDirection == 0 {
			if myDirection == 1 {
				cost += calculateCost(true, false, myFloor, orderFloor, 0, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(false, false, myFloor, orderFloor, 0, up, down, inside) 
			}
		}
	}
	return cost
}






//TODO: Gange opp cost med 1000 og legge til IP, slik at cost aldri er lik
//Nå beregnes costen for mange ganger på høyest IP (for some reason), så flere heiser tar samme best.
func GetCost(myFloor, myDirection, orderFloor, orderDirection int, myAddress string) int {
	
	//Find my cost:
	myCost := selectCostCase(myFloor, orderFloor, myDirection, orderDirection, Up, Down, Inside)
	
	myIP, _ := strconv.Atoi(myAddress)
	myCost = (myCost*1000)+myIP
	
	//Check if other elevator got lower cost:
	for IP, elev := range Elevators {
		if elev.Defekt == false {
			elevCost := selectCostCase(elev.LastFloor, orderFloor, elev.Direction, orderDirection, elev.Up, elev.Down, elev.Inside)
		
			elevIP, _ := strconv.Atoi(IP)
			elevCost = (elevCost*1000)+ elevIP
		
			if elevCost < myCost {
				return 0
			}
		}
	}

	return 1

}
