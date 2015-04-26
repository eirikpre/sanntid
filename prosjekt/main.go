package main

import (
	"runtime"
	
	"./src/variables"
	"./src/udp"
	"./src/driver"
)


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU()) 
	jobDone := make(chan bool,5)
	newOrders := make(chan variables.Order)
	nextFloor := make(chan int,100)
	ObsCh := make(chan bool,10)
	StopCh := make(chan bool,200)
	currentFloor := make(chan int,5)
	to_local_ch := make(chan variables.Status,100)
	from_local_ch := make(chan variables.Status)
	UDPreceiveStatus := make(chan variables.Status)
	UDPsendStatus := make(chan variables.Status)
	error_ch := make(chan error,200)
	

	driver.Init(newOrders,jobDone,StopCh,ObsCh,nextFloor,currentFloor)
	udp.Udp_Init(UDPsendStatus, UDPreceiveStatus,error_ch)

	go statusHandler(to_local_ch,from_local_ch,UDPreceiveStatus,UDPsendStatus, newOrders, jobDone, error_ch)
	go local_handler(to_local_ch,from_local_ch,jobDone,currentFloor,StopCh,nextFloor,error_ch,newOrders)
	
	errorHandler(UDPsendStatus, UDPreceiveStatus, error_ch)
	
}

