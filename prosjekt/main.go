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
Fix jobDone, proccing two times because sending on nextFloor line 73




*/


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
				}
			}
			if !(added){
				statuses = append(statuses, newStatus)
			}
			added = false

		case newStatus := <- LOCALsendStatus:
			fmt.Println("statusHandler: received from LOCALsendStatus =",newStatus)
			statuses[0] = newStatus
			UDPsendStatus <- newStatus

		case newOrder := <- newOrders :
			fmt.Println("statusHandler: received newOrder")
			//cost function!
			statuses[0].Orders = append(statuses[0].Orders,newOrder)
			statuses[0].Direction = statuses[0].Orders[0].Floor - statuses[0].Floor 
			fmt.Println("statusHandler: sending newStatus with newOrder", statuses[0])
			LOCALreceiveStatus <- statuses[0]

		}
	}

}


func localStatusHandler(receiveStatus chan variables.Status,sendStatus chan variables.Status,jobDone chan bool,currentFloor chan int,StopCh chan bool,ObsCh chan bool,nextFloor chan int){
	var localStatus variables.Status
	for{
		select{
		case localStatus = <-receiveStatus:
			fmt.Println("localStatusHandler: received new Status")
			//fmt.Println("localStatusHandler: Current OrderList:" ,localStatus.Orders)
			nextFloor <- localStatus.Orders[0].Floor
			
			
		case <- jobDone :
			fmt.Println("localStatusHandler: received jobDone")
			localStatus.Orders = append(localStatus.Orders[1:])

			if len(localStatus.Orders) != 0 {
				fmt.Println("localStatusHandler: More Orders to be done, starting next: " ,localStatus.Orders)
				localStatus.Direction = localStatus.Orders[0].Floor - localStatus.Floor 
				nextFloor <- localStatus.Orders[0].Floor
				
			}else{
				localStatus.Direction = 0
				
			}
			sendStatus <- localStatus

		case localStatus.Floor = <- currentFloor :
			fmt.Println("localStatusHandler: recieved currentFloor=",localStatus.Floor)
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


	jobDone := make(chan bool,1)
	newOrders := make(chan variables.Order)
	nextFloor := make(chan int,1)
	ObsCh := make(chan bool,10)
	StopCh := make(chan bool,10)
	currentFloor := make(chan int,50)
	LOCALreceiveStatus := make(chan variables.Status)
	LOCALsendStatus := make(chan variables.Status)
	UDPreceiveStatus := make(chan variables.Status, 20)
	UDPsendStatus := make(chan variables.Status, 20)
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
	status.Orders = append(status.Orders, variables.Order{Floor:0,Dir:0})

	return status

}