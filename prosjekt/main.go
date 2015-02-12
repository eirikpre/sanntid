package main

import (
	"fmt"
	"net"
	"runtime"
	"time"
	"./udp"
	"./variables"
)


 
func Init() {
	mySelf := new(Status)
	mySelf.laddr = "129.241.187.255"
	mySelf.floor = 1
	mySelf.direction = 2
	mySelf.destination = 2
	
	heiser.PushBack(mySelf)
	//Get floor, direction,destination
	//Set up WorkList
	// Store in heiser.Front()
}


func StatusUpdate(status Status){
	// Update heiser with received statuses
	
	
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	Init()
	
	//go UDPStatusUpdate()		//Updates local status & puts it in heiser
	go UDPSendStatus()		//Sending current status to all the other heiser
	go UDPReceiveStatus(baddr)	//Receiving statuses and stores them in heiser
	//go UDPReceiveWork()		//Receiving work and puts them in workList
	
	time.Sleep(2*time.Second)
	

}
