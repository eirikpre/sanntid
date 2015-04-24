package main

import (
	"runtime"
	"time"
	"fmt"
	"./src/variables"
	"./src/udp"
	"./src/driver"
	"math"
)

/*
TODO

Door Open
lastRead / receiveOrder / multiple equal orders 94- 98

Testing
UDP modul



*/



func sort(status variables.Status) variables.Status { 	// Sorts the Orders
	if len(status.Orders) < 2{
		if status.Floor > status.Orders[0].Floor{
			status.Direction = -1
		}else{
			status.Direction = 1
		}
		return status
	}
	fmt.Println("sort: I have to sort")
	var zeros [] variables.Order
	var currentDir [] variables.Order
	var wrongDir [] variables.Order

	for i:=0;i<len(status.Orders);i++{

			if status.Orders[i].Dir*status.Direction > 0{		
			// Adding the currentDir orders
				currentDir = append(currentDir[:], status.Orders[i])
			}else if (status.Orders[i].Dir*status.Direction < 0){	
			// Adding the wrongDir orders
				wrongDir = append(wrongDir[:], status.Orders[i])
			}else{	
			// Adding  the zero orders
				zeros = append(zeros[:], status.Orders[i])
			}		
	}
	status.Orders = nil
	status.Orders = bubbleSort(currentDir, status.Direction)
	wrongDir = bubbleSort(wrongDir, status.Direction*-1)
	zeros = bubbleSort(zeros, status.Direction)
	

	for i := range wrongDir{
		status.Orders = append( status.Orders[:], wrongDir[i])
	}


	
	// Adding the zeros in correct place
	if len(zeros) != 0{
		for i:=0; i<len(zeros);i++{	
			for j:=0;j<len(status.Orders);j++{
				if status.Orders[j].Floor >= zeros[i].Floor*status.Direction{

					wrongDir = status.Orders[j:]
					status.Orders = append(status.Orders[:], variables.Order{0,0} )
					copy( status.Orders[j+1:], wrongDir[:] )
					status.Orders[j] = zeros[i]

					break
				}
			}		
		}
	}

	printStatus(status)

	for i:=0; i<len(status.Orders); i++{
		if status.Orders[i].Floor*status.Direction <= status.Floor*status.Direction{
			status.Orders = append(status.Orders[0:],status.Orders[0])
		}else{
			break
		}
	}

	if status.Floor > status.Orders[0].Floor{
		status.Direction = -1
	}else{
		status.Direction = 1
	}


	fmt.Println("sort: Exiting")
	printStatus(status)
	return status
}


func costFunc(statuses []variables.Status, newOrder variables.Order) (variables.Status, int){
	fmt.Println("costFunc: Running")

	costArray := make([]int,len(statuses))
	// Buttons inside the elevator => job has to be done by self
	if newOrder.Dir == 0{ 	
		statuses[0].Orders = append(statuses[0].Orders[:],newOrder)
		statuses[0] = sort(statuses[0]) 
		return statuses[0],0
	}

	// Checking for identical orders
	for i:=0; i<len(statuses); i++ {     				
		for j:=0; j<len(statuses[i].Orders);j++{
			if statuses[i].Orders[j] == newOrder{
				fmt.Println("costFunc: identical order: ",newOrder)
				return statuses[i],-1
			}
		}
	}

	
	for i:=0;i<len(statuses);i++{						// Creating a costArray
		if (statuses[i].Floor == newOrder.Floor) && (statuses[i].Direction*newOrder.Dir >= 0){ 												// Elevator at the current floor && same direction
			
			statuses[i].Orders = append(statuses[i].Orders[:],newOrder)
			statuses[i] = sort(statuses[i]) 
			return statuses[i],i

		}else{
			for j:=0;j<len(statuses[i].Orders);j++{		// Add costs in the array
				costArray[i] += int(math.Abs(float64(statuses[i].Floor - statuses[i].Orders[j].Floor)))
				costArray[i] += int(math.Abs(float64(statuses[i].Floor - newOrder.Floor)))
			}		
			if statuses[i].Direction != newOrder.Dir && statuses[i].Direction != 0{
				costArray[i] += 10
			}			
		}
	}
	minimum:=256;
	position:=0;
	for i:=0;i<len(statuses);i++{ 			// Find the cheapest elevator
		if minimum > costArray[i]{
			minimum = costArray[i]
			position = i
		}
	}
	statuses[position].Orders = append(statuses[position].Orders[:],newOrder)
	statuses[position] = sort(statuses[position])
	fmt.Println("costFunc: Exiting")
	return statuses[position],position
}


func statusHandler(to_local_ch, from_local_ch, UDPreceiveStatus, UDPsendStatus chan variables.Status, newOrders chan variables.Order ){
	var statuses []variables.Status
	statuses = append(statuses[:],createStatus())
	var updated bool = false
	for{
		select{
		case newStatus := <- UDPreceiveStatus :
			fmt.Println("statusHandler: UDPreceiveStatus =",newStatus)
			for i:=0; i<len(statuses); i++{
				if statuses[i].Addr == newStatus.Addr{
					statuses[i] = newStatus
					updated = true
					if i == 0{
						to_local_ch <- newStatus
					}
					break
				}
			}
			if !(updated){
				statuses = append(statuses, newStatus)
			}
			updated = false

		case newStatus := <- from_local_ch:
			fmt.Println("statusHandler: from_local_ch =",newStatus)
			statuses[0] = newStatus
			UDPsendStatus <- newStatus

		case newOrder := <- newOrders :
			
			newStatus,position := costFunc(statuses,newOrder)
			if position >= 0 {
				if position == 0 {
					to_local_ch <- newStatus
				}
				UDPsendStatus <- newStatus
			}
				
		}	
	}
}

func local_handler(receiveStatus ,sendStatus chan variables.Status,jobDone chan bool,currentFloor chan int,StopCh chan bool,ObsCh chan bool,nextFloor chan int){
	var localStatus variables.Status
	

	for{
		select{
		case localStatus = <-receiveStatus:
			fmt.Println("local_handler: received new Status")
			
			if len(localStatus.Orders) > 0 {
				nextFloor <- localStatus.Orders[0].Floor
				fmt.Println("local_handler: newStatus sending nextFloor = ",localStatus.Orders[0].Floor)

			}
			
			
			
		case <- jobDone :
			//fmt.Println("local_handler: received jobDone")
			if len(localStatus.Orders) > 1 {
				

				localStatus.Orders = append(localStatus.Orders[1:])
				nextFloor <- localStatus.Orders[0].Floor

				if localStatus.Floor > localStatus.Orders[0].Floor{
					localStatus.Direction = -1
				}else if localStatus.Floor < localStatus.Orders[0].Floor{
					localStatus.Direction = 1
				}

			
			}else if len(localStatus.Orders) == 1{
				
				localStatus.Orders = append(localStatus.Orders[1:])
				

			}else{
				localStatus.Direction = 0
			}
			sendStatus <- localStatus

		case localStatus.Floor = <- currentFloor :
			
			if len(localStatus.Orders) > 0{
				if localStatus.Floor > localStatus.Orders[0].Floor{
					localStatus.Direction = -1
				}else if localStatus.Floor < localStatus.Orders[0].Floor{
					localStatus.Direction = 1
				}
			}

			
			sendStatus <- localStatus
		case <- StopCh :
			fmt.Println("STOOOP received")
			//send jobs to others?
			//set stop bit in status
		case <- ObsCh  :
			fmt.Println("OOOOBS received")
		}
	}

}



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
	UDPreceiveStatus := make(chan variables.Status, 50)
	UDPsendStatus := make(chan variables.Status, 50)
	
	driver.Init(nextFloor,jobDone,newOrders,StopCh,ObsCh,currentFloor)
	
	go statusHandler(to_local_ch,from_local_ch,UDPreceiveStatus,UDPsendStatus, newOrders,)
	go local_handler(to_local_ch,from_local_ch,jobDone,currentFloor,StopCh,ObsCh,nextFloor)
	//errInit := udp.Udp_init(variables.lport, variables.bport, msg_size int, send_ch, receive_ch chan Udp_message)
	//driver.Init()
	for{
		time.Sleep(time.Second*10)
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

func bubbleSort(orders []variables.Order,direction int) []variables.Order { 	
	var temp variables.Order
	for i:=1;i<len(orders);i++{
		for j:=1;j<len(orders);j++{
			if orders[j-1].Floor*direction > orders[j].Floor*direction {
				temp = orders[j-1]
				orders[j-1] = orders[j]
				orders[j] = temp
			}
		}
	}
	//fmt.Printf("bubbleSort: Orders: %v\n",orders)
	return orders
}