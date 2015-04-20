package main

import (
	"runtime"
	"time"
	"fmt"
	"./src/variables"
	//"./src/udp"
	"./src/driver"
)

/*
TODO

Door Open
lastRead / receiveOrder / multiple equal orders 94- 98

Testing
UDP modul



*/
func sort(status variables.Status) variables.Status { 	// Sorts the Orders
	var zeros []variables.Order
	var currentDir []variables.Order
	var wrongDir []variables.Order
	
	for i:=0;i<len(status.Orders);i++{					
		if status.Orders[i].Dir*status.Direction > 0{		// Adding a sorting the currentDir orders
			for j:=0;j<len(currentDir);j++{
				if currentDir[j].Floor*status.Direction < currentDir[j+1].Floor {
					currentDir = append(currentDir[:j], status.Orders[i], currentDir[(j+1):])
					break
				}
			}
		}if else status.Orders[i].Dir*status.Direction < 0{	// Adding a sorting the wrongDir orders
			for j:=0;j<len(wrongDir);j++{
				if wrongDir[j].Floor*status.Direction > wrongDir[j+1].Floor {
					wrongDir = append(wrongDir[:j], status.Orders[i], wrongDir[(j+1):])
					break
				}
			}
		}else{
			for j:=0;j<len(zeros);j++{						// Adding a sorting the zeros orders
				if zeros[j].Floor*status.Direction < zeros[j+1].Floor {
					zeros = append(zeros[:j], status.Orders[i], zeros[(j+1):])
					break
				}
			}
		}		
	}
	status.Orders = append(currentDir, wrongDir)
	for i:=0; i<len(zeros);i++{								// Adding the zeros in correct place
		for j:=0;j<len(status.Orders);j++{
			if status.Orders[j] >= zeros[i]*status.Direction{
				status.Orders = append(status.Orders[:j], zeros[i], status.Orders[(j+1):])
				break
			}
		}		
	}
	for i:=0; i<len(status.Orders); i++{
		if status.Orders[0].Floor > status.Floor*status.Direction{
			status.Orders = append(status.Orders[1:],status.Orders[0])
		}	
	}
	return status
}
func costFunc(statuses []variables.Status, newOrder variables.Order) variables.Status, int{
	var costArray []int = make(int[],len(statuses));
	if newOrder.Dir == 0{ 					// Buttons inside the elevator => job has to be done by self
		statuses[0].Orders = append(statuses[0].Orders[:],newOrder)
		statuses[i] = sort(statuses[i]) // Sorting
		return statuses[0],0
	}
	for i:=0; i<len(statuses); i++ {     	// Checking for similar orders
		for j:=0; j<statuses[i].Orders;j++{
			if statuses[i].Orders[j] == newOrder{
				return statuses[i],-1
			}
		}
	}

	
	for i:=0;i<len(statuses);i++{
		if statuses[i].Floor == newOrder.Floor && (statuses[i].Direction == 0 || statuses[i].Direction/abs(statuses[i].Direction) == newOrder.Dir)
		{ 									// Elevator at the current floor && same direction
			statuses[i].Orders = append(statuses[i].Orders[:],newOrder)
			statuses[i] = sort(statuses[i]) // Sorting
			return statuses[i],i
		}else{
			for j:=0;j<len(statuses[i].Orders);j++{		// Add costs in the array
				costArray[i] += abs(statuses[i].Floor - statuses[i].Orders[j].Floor)
				costArray[i] += abs(statuses[i].Floor - newOrder.Floor)
			}		
			if statuses[i].Direction/abs(statuses[i].Direction) != newOrder.Dir && statuses[i].Direction != 0{
				costArray[i] += 10
			}			
		}
	}
	minimum:=256;
	position:=0;
	for i:=0;i<len(statuses);i++{ 			// Find the cheapest elevator
		if minimum > costArray[i]{
			minimum = costArray[i]
			position = i
		}
	}
	statuses[position].Orders = append(statuses[position].Orders[:],newOrder)
	statuses[i] = sort(statuses[i]) // Sorting
	return statuses[position],position
}


func statusHandler(LOCALreceiveStatus chan variables.Status, LOCALsendStatus chan variables.Status, UDPreceiveStatus chan variables.Status,	UDPsendStatus chan variables.Status, newOrders chan variables.Order ){
	
	var statuses []variables.Status
	statuses = append(statuses, createStatus())
	var added bool = false
	for{
		select{
		case newStatus := <- UDPreceiveStatus :
			fmt.Println("statusHandler: received from UDPreceiveStatus =",newStatus)
			for i:=0; i<len(statuses); i++{
				if statuses[i].Addr == newStatus.Addr{
					statuses[i] = newStatus
					added = true
					if i == 0{
						LOCALreceiveStatus <- newStatus
					}
					break
				}
			}
			if !(added){
				statuses = append(statuses, newStatus)
			}
			added = false

		case newStatus := <- LOCALsendStatus:
			fmt.Println("statusHandler: received from LOCALsendStatus =",newStatus)
			statuses[0] = newStatus
			//UDPsendStatus <- newStatus

		case newOrder := <- newOrders :
			
			//cost function!
			
			newStatus,position := costFunc(statuses,newOrder)
			if position >= 0 {
				if position == 0 {
					LOCALreceiveStatus <- newStatus
				}
				//UDPsendStatus <- newStatus
				
			}	
				

		}
	}

}


func localStatusHandler(receiveStatus chan variables.Status,sendStatus chan variables.Status,jobDone chan bool,currentFloor chan int,StopCh chan bool,ObsCh chan bool,nextFloor chan int){
	var localStatus variables.Status
	

	for{
		select{
		case localStatus = <-receiveStatus:
			fmt.Println("localStatusHandler: received new Status")
			
			if len(localStatus.Orders) > 0 {
				nextFloor <- localStatus.Orders[0].Floor
				fmt.Println("localStatusHandler: newStatus sending nextFloor = ",localStatus.Orders[0].Floor)

			}
			
			
			
		case <- jobDone :
			//fmt.Println("localStatusHandler: received jobDone")
			if len(localStatus.Orders) > 1 {
				//fmt.Println("localStatusHandler: More Orders to be done, starting next: " ,localStatus.Orders)
				localStatus.Orders = append(localStatus.Orders[1:])
				localStatus.Direction = localStatus.Orders[0].Floor - localStatus.Floor 
				nextFloor <- localStatus.Orders[0].Floor
				//fmt.Println("localStatusHandler: jobDone sending nextFloor = ",localStatus.Orders[0].Floor)
			}else if len(localStatus.Orders) == 1{
				//fmt.Println("localStatusHandler: Last order done")
				localStatus.Orders = append(localStatus.Orders[1:])
				localStatus.Direction = 0
			}else{
				localStatus.Direction = 0
			}
			sendStatus <- localStatus

		case localStatus.Floor = <- currentFloor :
			//fmt.Println("localStatusHandler: recieved currentFloor=",localStatus.Floor)
			if len(localStatus.Orders) != 0{
				localStatus.Direction = localStatus.Orders[0].Floor - localStatus.Floor
			}
			sendStatus <- localStatus
		case <- StopCh :
			fmt.Println("STOOOP received")
			//send jobs to others?
			//set stop bit in status
		case <- ObsCh  :
			fmt.Println("OOOOBS received")
		}
	}

}





func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())


	jobDone := make(chan bool)
	newOrders := make(chan variables.Order)
	nextFloor := make(chan int,10)
	ObsCh := make(chan bool,10)
	StopCh := make(chan bool,10)
	currentFloor := make(chan int)
	LOCALreceiveStatus := make(chan variables.Status)
	LOCALsendStatus := make(chan variables.Status)
	UDPreceiveStatus := make(chan variables.Status, 50)
	UDPsendStatus := make(chan variables.Status, 50)
	driver.Init(nextFloor,jobDone,newOrders,StopCh,ObsCh,currentFloor)

	//statusInit(LOCALreceiveStatus,LOCALsendStatus)
	go statusHandler(LOCALreceiveStatus,LOCALsendStatus,UDPreceiveStatus,UDPsendStatus, newOrders,)
	go localStatusHandler(LOCALreceiveStatus,LOCALsendStatus,jobDone,currentFloor,StopCh,ObsCh,nextFloor)
	//errInit := udp.Udp_init(variables.lport, variables.bport, msg_size int, send_ch, receive_ch chan Udp_message)
	//driver.Init()
	for{
		time.Sleep(time.Second*10)
	}
}


func statusInit(receiveStatus chan variables.Status, sendStatus chan variables.Status){
	status := variables.Status{Floor:0,Direction:0,Stop:false}
	status.Orders = append(status.Orders, variables.Order{Floor:0,Dir:0})
	receiveStatus <- status
	sendStatus <- status

}

func createStatus() variables.Status {
		status := variables.Status{Floor:0,Direction:0,Stop:false}
		// Sette riktig ip-adresse.
	status.Orders = append(status.Orders, variables.Order{Floor:0,Dir:0})

	return status

}