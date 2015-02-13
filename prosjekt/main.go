package main

import (
	"fmt"
	"net"
	"runtime"
	"time"
	"./udp"
	"./variables"
	"./driver"
)


 
func Init(elev_ch chan []Status) {

	//driver.Init()
	//mySelf = driver.getStatus()
	append(elevs, mySelf)
	elev_ch <- elevs
	errUDP := udp.Udp_init(variables.bport, variables.bport, msg_size int, send_ch, receive_ch chan Udp_message)
}


func StatusUpdate(status Status){
	// Update heiser with received statuses
	
	
}



func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	var receive_ch chan Udp_message
	var elevs []Status
	var elev_ch chan []Status = make(chan []Status,1)
	var send_ch chan Udp_message
	var err_ch
	Init(elev_ch)
	errInit := udp.Udp_init(variables.lport, variables.bport, msg_size int, send_ch, receive_ch chan Udp_message)
		
	
	time.Sleep(2*time.Second)
	

}
