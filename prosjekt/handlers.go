package main 

import (
	"fmt"
	"./src/variables"
	"./src/udp"
	"./src/driver"
)

func refreshLights(statuses []variables.Status) {
			//Turn every light off.
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



func statusHandler(to_local_ch, from_local_ch, UDPreceiveStatus, UDPsendStatus chan variables.Status, newOrders chan variables.Order ){
	var statuses []variables.Status
	statuses = append(statuses[:],createStatus())
	var updated bool = false
	for{
		select{
		case newStatus := <- UDPreceiveStatus :
			
			//fmt.Println("statusHandler: UDPreceiveStatus =\n",newStatus)
			for i:=0; i<len(statuses); i++{

				if statuses[i].Addr == newStatus.Addr{
					statuses[i] = newStatus
					updated = true

					fmt.Println("\n----------------STATUSES---------------------")
					for i:=0;i<len(statuses);i++{
						fmt.Println(statuses[i])
					}
					if i == 0{
						to_local_ch <- newStatus
					}
					break
				}
			}
			if !(updated){
				statuses = append(statuses[:], newStatus)
				UDPsendStatus <- statuses[0]
			}
			refreshLights(statuses)
			updated = false


		case newStatus := <- from_local_ch:
			//fmt.Println("statusHandler: from_local_ch =",newStatus)
			
			statuses[0] = newStatus
			UDPsendStatus <- newStatus

		case newOrder := <- newOrders :
			
			newStatus,position := costFunc(statuses,newOrder)

			if position >= 0 {

				if newStatus.Floor > newStatus.Orders[0].Floor{
					newStatus.Direction = -1
				}else if newStatus.Floor < newStatus.Orders[0].Floor{
					newStatus.Direction = 1
				}
				statuses[position] = newStatus

				if position == 0 {
					to_local_ch <- newStatus
				}

				UDPsendStatus <- newStatus
			}
				
		}	
	}
}

func local_handler(to_local_ch ,from_local_ch chan variables.Status,jobDone chan bool,currentFloor chan int,StopCh chan bool,ObsCh chan bool,nextFloor chan int){
	var localStatus variables.Status = createStatus()
	fmt.Println("local_handler: Initializing")
	from_local_ch <- localStatus
	nextFloor <- 0

	for{
		select{
		case localStatus = <- to_local_ch:
			 
			if len(localStatus.Orders) > 0 {
				nextFloor <- localStatus.Orders[0].Floor
				//fmt.Println("local_handler: newStatus sending nextFloor = ",localStatus.Orders[0].Floor)

			}
			
			
			
		case <- jobDone :
			//fmt.Println("local_handler: received jobDone")
			if len(localStatus.Orders) > 1 && localStatus.Orders[0].Floor == localStatus.Floor{
				

				localStatus.Orders = append(localStatus.Orders[1:])
				nextFloor <- localStatus.Orders[0].Floor

				if localStatus.Floor > localStatus.Orders[0].Floor{
					localStatus.Direction = -1
				}else if localStatus.Floor < localStatus.Orders[0].Floor{
					localStatus.Direction = 1
				}

				from_local_ch <- localStatus
			
			}else if len(localStatus.Orders) == 1 && localStatus.Orders[0].Floor == localStatus.Floor{
				
				localStatus.Orders = append(localStatus.Orders[1:])
				localStatus.Direction = 0
				from_local_ch <- localStatus
			}

			

		case localStatus.Floor = <- currentFloor :
			
			if len(localStatus.Orders) > 0{
				if localStatus.Floor > localStatus.Orders[0].Floor{
					localStatus.Direction = -1
				}else if localStatus.Floor < localStatus.Orders[0].Floor{
					localStatus.Direction = 1
				}
			}

			from_local_ch <- localStatus

		case <- StopCh :
			fmt.Println("STOOOP received")
			//send jobs to others?
			//set stop bit in status
		case <- ObsCh  :
			fmt.Println("OOOOBS received")
		}
	}

}


func createStatus() variables.Status{
	status := variables.Status{Floor:1,Direction:-1,Stop:false}
	status.Addr = udp.GetOwnIP()
	status.Orders = append(status.Orders, variables.Order{Floor:0,Dir:0})
	return status
}

func printStatus(status variables.Status){
	//fmt.Println("Addr: ",status.Addr)
	//fmt.Println("Floor: ",status.Floor," Direction: ",status.Direction)
	//fmt.Printf("Orders: %v", statuses.Orders)
	fmt.Printf("%v\n", status)
}

