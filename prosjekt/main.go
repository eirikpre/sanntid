package main

import (
	"runtime"
	"time"
	"./src/variables"
	"./src/udp"
	"./src/driver"
)

/*
TODO

Door Open
lastRead / receiveOrder / multiple equal orders 94- 98

Testing
UDP modul



*/




func main(){
	runtime.GOMAXPROCS(runtime.NumCPU()) 
	jobDone := make(chan bool)
	newOrders := make(chan variables.Order)
	nextFloor := make(chan int,10)
	ObsCh := make(chan bool,10)
	StopCh := make(chan bool,10)
	currentFloor := make(chan int)
	to_local_ch := make(chan variables.Status)
	from_local_ch := make(chan variables.Status)
	UDPreceiveStatus := make(chan variables.Status)
	UDPsendStatus := make(chan variables.Status)
	

	driver.Init(newOrders,jobDone,StopCh,ObsCh,nextFloor,currentFloor)
	udp.Udp_Init(UDPsendStatus, UDPreceiveStatus)

	go statusHandler(to_local_ch,from_local_ch,UDPreceiveStatus,UDPsendStatus, newOrders, jobDone)
	go local_handler(to_local_ch,from_local_ch,jobDone,currentFloor,StopCh,nextFloor)
	

	for{
		time.Sleep(time.Second*10)
	}
	
}

