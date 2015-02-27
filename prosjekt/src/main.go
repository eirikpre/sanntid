package main

import (
	"runtime"
	"time"
	//"variables"
	//"udp"
	"driver"
)

/*
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
*/


func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
/*
	var receive_ch chan udp.Udp_message
	     var elevs []variables.Status
	var elev_ch chan []variables.Status = make(chan []variables.Status,1)
	var send_ch chan udp.Udp_message
*/
	

	driver.Init()
	
	//errInit := udp.Udp_init(variables.lport, variables.bport, msg_size int, send_ch, receive_ch chan Udp_message)
	//driver.Init()
	for{
	time.Sleep(time.Second*10)
	}
}

func eliajgiadg
