package main 

import (
	"fmt"
	"./src/variables"
	"./src/udp"
	"./src/driver"
	"time"
	"errors"
	"math/rand"
)

func errorHandler(UDPsendStatus, UDPreceiveStatus chan variables.Status, error_ch chan error){

	for{

		error_msg := <- error_ch

		switch{
		case error_msg == nil:
			break
		case error_msg.Error() == "STOP":
			//????
		case error_msg.Error() == "UDPRESET":
				time.Sleep(2*time.Second)
				udp.Udp_Init(UDPsendStatus, UDPreceiveStatus,error_ch)
		default:
			fmt.Println(error_msg)

		}
	}
}







func statusHandler(to_local_ch, from_local_ch, UDPreceiveStatus, UDPsendStatus chan variables.Status, newOrders chan variables.Order , jobDone chan bool, error_ch chan error){
	fmt.Println("statusHandler: Initializing")
	statuses := make([]variables.Status,0)
	timers := make([]*time.Timer,0)
	statuses = append(statuses[:], createStatus(error_ch))
	timers = append(timers[:], time.NewTimer(20*time.Second))
	var newOrder variables.Order
	var tempStatus variables.Status


	// watchDog
	go func() {
		for{

			for i:=0;i<len(timers);i++{
				
				select{

				case <- timers[i].C:
					fmt.Println("Timeout on ", statuses[i])
					if len(statuses[i].Orders) > 0{
						statuses[i].Stop = true
						dummy := statuses[i]
						statuses[i].Orders = nil
						UDPsendStatus <- statuses[i]
						error_ch <- errors.New("TIMEOUT")
						time.Sleep(200*time.Millisecond)
						for j:=0; j<len(dummy.Orders); j++{
							if dummy.Orders[j].Dir != 0{
								newOrders <- dummy.Orders[j]
							}
						}
					} 
				default:
					time.Sleep(time.Duration(rand.Intn(1000)*4)*time.Millisecond)
				}
			}
		}
	}()



	go func (){
		var updated bool = false
		for{

			select{
			case tempStatus = <- UDPreceiveStatus :


				for i:=0; i<len(statuses); i++{

					if (statuses[i].Addr == tempStatus.Addr)	{
						statuses[i] = tempStatus
						updated = true
						timers[i].Reset(20*time.Second)
						if i == 0{

							to_local_ch <- statuses[0]
							
						}
						
					}
				}

				if !(updated){

					statuses = append(statuses[:], tempStatus)
					timers = append(timers[:],time.NewTimer(20*time.Second))
					UDPsendStatus <- statuses[0]
				}
				updated = false
				refreshLights(statuses)

				fmt.Println("\n-------------------------------STATUSES-----------------------------------")
				for i:=0;i<len(statuses);i++{
					fmt.Println(statuses[i])
				}
				
				


			case tempStatus = <- from_local_ch:

					statuses[0] = tempStatus
					UDPsendStatus <- tempStatus
				

			case newOrder = <- newOrders :
				
				status,position := costFunc(statuses,newOrder)
				
				if position >= 0 {
					
					statuses[position] = status
					statuses[position].Direction = getDir(statuses[position])
					UDPsendStatus <- statuses[position]

				}		
			}	
		}
	}()
}

func local_handler(to_local_ch ,from_local_ch chan variables.Status,jobDone chan bool,currentFloor chan int,StopCh chan bool,nextFloor chan int, error_ch chan error,newOrders chan variables.Order){
	var localStatus variables.Status = createStatus(error_ch)
	fmt.Println("local_handler: Initializing")
	nextFloor <- 0

	for{
		select{
		case localStatus = <- to_local_ch:

			if len(localStatus.Orders) > 0 {
				nextFloor <- localStatus.Orders[0].Floor
			}

		case <- jobDone :
			
			if len(localStatus.Orders) > 1 && localStatus.Orders[0].Floor == localStatus.Floor {

				localStatus.Orders = append(localStatus.Orders[1:])
				localStatus.Direction = getDir(localStatus)
				localStatus = sort(localStatus)
				nextFloor <- localStatus.Orders[0].Floor

			}else if len(localStatus.Orders) == 1 && localStatus.Orders[0].Floor == localStatus.Floor{

				localStatus.Orders = nil
				localStatus.Direction = 0
				
			}
			localStatus.Stop = false
	
			from_local_ch <- localStatus
		

		case localStatus.Floor = <- currentFloor :
			
			localStatus.Direction = getDir(localStatus)
			from_local_ch <- localStatus

		case <- StopCh :
			localStatus.Stop = true
			from_local_ch <- localStatus
			time.Sleep(200*time.Millisecond)
			for i:=0; i<len(localStatus.Orders); i++{
				if localStatus.Orders[i].Dir != 0{
					newOrders <- localStatus.Orders[i]
				}
			}
			error_ch <- errors.New("STOP")
			return

		}
	}

}


func createStatus(error_ch chan error) variables.Status{
	status := variables.Status{Floor:1,Direction:-1,Stop:false}
	status.Addr = udp.GetOwnIP(error_ch)
	status.Orders = append(status.Orders, variables.Order{Floor:0,Dir:0})
	return status
}

func getDir(status variables.Status) int {

	if len(status.Orders) > 0{
		if status.Floor > status.Orders[0].Floor{
			return -1
		}else if status.Floor < status.Orders[0].Floor{
			return 1
		}
	}
	return status.Direction
}


func refreshLights(statuses []variables.Status) {

	driver.LightButtons(0, false)
	driver.LightButtons(3, false)
	driver.LightButtons(4, false)
	driver.LightButtons(6, false)
	driver.LightButtons(7, false)
	driver.LightButtons(10, false)


	for i:=0 ; i<len(statuses); i++ {
		for j:=0 ; j<len(statuses[i].Orders) ; j++{
			if statuses[i].Stop{
				break
			}else{
				switch(statuses[i].Orders[j]){
	
				case variables.Order{0,1}:
					driver.LightButtons(0,true)
					
				case variables.Order{1,1}:
					driver.LightButtons(3, true)
					
				case variables.Order{1,-1}:
					driver.LightButtons(4, true)
					
				case variables.Order{2,1}:
					driver.LightButtons(6, true)
					
				case variables.Order{2,-1}:
					driver.LightButtons(7, true)
					
				case variables.Order{3,-1}:
					driver.LightButtons(10, true)
					
			}
			}
		}
	}
}
