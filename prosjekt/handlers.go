package main 

import (
	"fmt"
	"./src/variables"
	"./src/udp"
)




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
					if i == 0{
						to_local_ch <- newStatus
					}
					fmt.Println("----------------STATUSES---------------------")
					for i:=0;i<len(statuses);i++{
						fmt.Println(statuses[i])
					}
					break
				}
			}
			if !(updated){
				statuses = append(statuses[:], newStatus)
			}

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

			
			}else if len(localStatus.Orders) == 1 && localStatus.Orders[0].Floor == localStatus.Floor{
				
				localStatus.Orders = append(localStatus.Orders[1:])
				localStatus.Direction = 0

			}

			from_local_ch <- localStatus

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

