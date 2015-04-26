package main 

import (
	"fmt"
	"./src/variables"
	"./src/udp"
	"./src/driver"
)

func statusHandler(to_local_ch, from_local_ch, UDPreceiveStatus, UDPsendStatus chan variables.Status, newOrders chan variables.Order , jobDone chan bool){
	fmt.Println("statusHandler: Initializing")
	statuses := make([]variables.Status,0)

	statuses = append(statuses[:],createStatus())

	var newOrder variables.Order
	var tempStatus variables.Status
	var position int

	var updated bool = false
	for{
		select{
		case tempStatus = <- UDPreceiveStatus :

			for i:=0; i<len(statuses); i++{

				if (statuses[i].Addr == tempStatus.Addr)	{
					//fmt.Println("udp: updating",statuses[i],"into ",tempStatus," pos:",i)
					statuses[i] = tempStatus
					updated = true
					if i == 0{
						//fmt.Println("Sending to_local_ch: ", statuses[0])
						to_local_ch <- statuses[0]
						
					}
					
				}
			}

			if !(updated){
				//fmt.Println("Adding a new entry to statuses!!")
				
				statuses = append(statuses[:], tempStatus)
				UDPsendStatus <- statuses[0]
			}
			updated = false


			fmt.Println("\n----------------STATUSES---------------------")
			for i:=0;i<len(statuses);i++{
				fmt.Println(statuses[i])
			}
			
			refreshLights(statuses)


		case tempStatus = <- from_local_ch:
				//fmt.Println("status_from_local: updating",statuses[0],"into ",tempStatus)
				statuses[0] = tempStatus
				UDPsendStatus <- tempStatus
			

		case newOrder = <- newOrders :
			
			statuses,position = costFunc(statuses,newOrder)
			

			if position >= 0 {
				//fmt.Println("Posistion after cost: ", position)
				statuses[position].Direction = getDir(statuses[position])
				//fmt.Println("Dir after getDir: ", statuses[position].Direction)
				UDPsendStatus <- statuses[position]

			}		
		}	
	}
}

func local_handler(to_local_ch ,from_local_ch chan variables.Status,jobDone chan bool,currentFloor chan int,StopCh chan bool,nextFloor chan int){
	var localStatus variables.Status = createStatus()
	fmt.Println("local_handler: Initializing")
	//from_local_ch <- localStatus
	nextFloor <- 0

	for{
		select{
		case localStatus = <- to_local_ch:
			 //fmt.Println("localStatus Order is: ", localStatus.Orders)
			if len(localStatus.Orders) > 0 {
				//fmt.Println("localStatus case: Not empty")
				nextFloor <- localStatus.Orders[0].Floor

			}else{
				//fmt.Println("localStatus case: Orders are empty")
			}
			
			
		// Some job has been done, target == tempFloor
		case joDone := <- jobDone :
			//fmt.Println("jobDone : ", joDone)
			if(joDone == true){
				//fmt.Println("In jobDone case: ", localStatus.Floor)
				if len(localStatus.Orders) > 1 && localStatus.Orders[0].Floor == localStatus.Floor{
					//fmt.Println("I'm here if orderlist is larger than 1")
					localStatus.Orders = append(localStatus.Orders[1:])
					//fmt.Println("localStatus.Orders : ", localStatus.Orders[:])
					localStatus = sort(localStatus)
					nextFloor <- localStatus.Orders[0].Floor

					localStatus.Direction = getDir(localStatus)

				}else if len(localStatus.Orders) == 1 && localStatus.Orders[0].Floor == localStatus.Floor{
					//fmt.Println("I'm here if orderlist equals 1, and order.floor = Floor")
					localStatus.Orders = nil
					localStatus.Direction = 0
					
				}
				//fmt.Println("I'm at the end of jobDone")
				from_local_ch <- localStatus
			}

		case localStatus.Floor = <- currentFloor :
			//fmt.Println("In localStatus case: ", localStatus.Floor)
			
			localStatus.Direction = getDir(localStatus)
			from_local_ch <- localStatus

		case <- StopCh :
			fmt.Println("STOOOP received")
			//send jobs to others?
			//set stop bit in status


		}
	}

}


func createStatus() variables.Status{
	status := variables.Status{Floor:1,Direction:-1,Stop:false}
	status.Addr = udp.GetOwnIP()
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
