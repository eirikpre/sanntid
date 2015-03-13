package main

import (
	"runtime"
	"time"
	"fmt"
	"./src/variables"
	//"./src/udp"
	"./src/driver"
)

func localStatusHandler(receiveStatus chan variables.Status,sendStatus chan variables.Status,jobDone chan bool,currentFloor chan int,StopCh chan bool,ObsCh chan bool,nextFloor chan int){
	var orderList []variables.Order
	var localStatus variables.Status




	for{
		select{
		case newStatus := <-receiveStatus:

			orderList = newStatus.Orders 
			fmt.Println("localStatusHandler: sending nextFloor=",orderList[0].Floor)
			nextFloor <- orderList[0].Floor
			
		case <- jobDone :

			orderList = append(orderList[1:])
			localStatus.Orders = orderList
			sendStatus <- localStatus
			if len(orderList) != 0{
				nextFloor <- orderList[0].Floor
			}
			
		case <- StopCh :
			fmt.Println("STOOOP received")
		case <- ObsCh  :
			fmt.Println("OOOOBS received")






		}
	}

}
func statusInit(receiveStatus chan variables.Status, sendStatus chan variables.Status){
	status := variables.Status{Floor:0,Direction:0,Stop:false}
	status.Orders = append(status.Orders, variables.Order{Floor:0,Dir:0})
	receiveStatus <- status
	sendStatus <- status

}
func job2status(receiveStatus chan variables.Status, newOrders chan variables.Order,sendStatus chan variables.Status){
	var status variables.Status
	for{
		select{ 
			case newOrder := <- newOrders :
			status.Orders = append(status.Orders,newOrder)
			receiveStatus <- status
			case status = <- sendStatus :
		}
	}

}


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
/*
	var receive_ch chan udp.Udp_message
	     var elevs []variables.Status
	var elev_ch chan []variables.Status = make(chan []variables.Status,1)
	var send_ch chan udp.Udp_message
*/
	jobDone := make(chan bool, 50)
	newOrders := make(chan variables.Order,0)
	nextFloor := make(chan int,1)
	ObsCh := make(chan bool,10)
	StopCh := make(chan bool,10)
	currentFloor := make(chan int,50)
	receiveStatus := make(chan variables.Status,20)
	sendStatus := make(chan variables.Status, 20)
	driver.Init(nextFloor,jobDone,newOrders,StopCh,ObsCh,currentFloor)
	statusInit(receiveStatus,sendStatus)
	go job2status(receiveStatus,newOrders,sendStatus)
	go localStatusHandler(receiveStatus,sendStatus,jobDone,currentFloor,StopCh,ObsCh,nextFloor)
	//errInit := udp.Udp_init(variables.lport, variables.bport, msg_size int, send_ch, receive_ch chan Udp_message)
	//driver.Init()
	for{
		time.Sleep(time.Second*10)
	}
}
