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




*/
func distributeOrder(statuses []variables.Status, newOrder variables.Order) variables.Status{
	for j:=0; j<40;j++{
		for i:=0; i<len(statuses); i++{
			if len(statuses[i].Orders[j:]) == 0{
				statuses[i].Orders = append(statuses[i].Orders, newOrder)
				statuses[i].Direction = statuses[i].Orders[0].Floor - statuses[i].Floor
				 
				return statuses[i]
			}
		}
	}
	statuses[0].Orders = append(statuses[0].Orders, newOrder)
	statuses[0].Direction = statuses[0].Orders[0].Floor - statuses[0].Floor
	return statuses[0]

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
			
			newStatus := distributeOrder(statuses,newOrder)
			//UDPsendStatus <- newStatus
			fmt.Println("statusHandler: sending newStatus with newOrder", newStatus)
			
			LOCALreceiveStatus <- newStatus

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
			fmt.Println("localStatusHandler: received jobDone")

			
			//fmt.Println("localStatusHandler: length = ",len(localStatus.Orders))
			if len(localStatus.Orders) > 1 {
				//fmt.Println("localStatusHandler: More Orders to be done, starting next: " ,localStatus.Orders)
				localStatus.Orders = append(localStatus.Orders[1:])
				localStatus.Direction = localStatus.Orders[0].Floor - localStatus.Floor 
				nextFloor <- localStatus.Orders[0].Floor
				fmt.Println("localStatusHandler: jobDone sending nextFloor = ",localStatus.Orders[0].Floor)
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