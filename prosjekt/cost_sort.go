package main 
 
import (
	"fmt"
	"./src/variables"
	"math"
)



func costFunc(statuses []variables.Status, newOrder variables.Order) ([]variables.Status, int){

	costArray := make([]int,len(statuses))

	// Checking for identical orders
	for i:=0; i<len(statuses); i++ {     				
		for j:=0; j<len(statuses[i].Orders);j++{
			if statuses[i].Orders[j] == newOrder {
				fmt.Println("costFunc: identical order: ",newOrder)
				return statuses,-1
			}
		}
		if newOrder.Dir == 0 && i > 0{
			break
		}
	}


	// Buttons inside the elevator => job has to be done by self
	if newOrder.Dir == 0{ 	
		statuses[0].Orders = append(statuses[0].Orders[:],newOrder)
		statuses[0] = sort(statuses[0]) 
		return statuses,0
	}


	// Generating a costArray
	for i:=0;i<len(statuses);i++{						
		for j:=0;j<len(statuses[i].Orders);j++{	
			costArray[i] += int(math.Abs(float64(statuses[i].Floor - statuses[i].Orders[j].Floor)))
			costArray[i] += int(math.Abs(float64(statuses[i].Floor - newOrder.Floor)))
		}		
		if statuses[i].Direction != newOrder.Dir && statuses[i].Direction != 0{
			costArray[i] += 10
		}			
	}

	minimum:=256;
	position:=-1;

	// Find the cheapest elevator
	for i:=0;i<len(statuses);i++{ 			
		if minimum > costArray[i]{
			minimum = costArray[i]
			position = i
		}
	}

	// Add to the cheapest & sort
	fmt.Printf("costFunc: updating%v",statuses[position])
	statuses[position].Orders = append(statuses[position].Orders[:],newOrder)
	statuses[position] = sort(statuses[position])
	fmt.Println("into ",statuses[position])
	return statuses,position
}

func sort(status variables.Status) variables.Status { 	// Sorts the Orders
	if len(status.Orders) < 2{
		return status
	}

	if status.Floor > status.Orders[0].Floor{
		status.Direction = -1
	}else{
		status.Direction = 1
	}

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
			// Adding the zero orders
				zeros = append(zeros[:], status.Orders[i])
			}		
	}
	

	// Sorting individual lists
	status.Orders = bubbleSort(currentDir, status.Direction)
	wrongDir = bubbleSort(wrongDir, status.Direction*-1)
	zeros = bubbleSort(zeros, status.Direction)

	for i:=0;i<len(wrongDir);i++{
		status.Orders = append( status.Orders[:], wrongDir[i])
	}


	
	// Adding the zeros in correct place
	if len(status.Orders) != 0{
		for i:=0; i<len(zeros);i++{	
			for j:=0;j<len(status.Orders);j++{
			
				if status.Orders[j].Floor*status.Direction >= zeros[i].Floor*status.Direction{

					wrongDir = status.Orders[j:]
					status.Orders = append(status.Orders[:], variables.Order{0,0} )
					copy( status.Orders[j+1:], wrongDir[:] )
					status.Orders[j] = zeros[i]
					break

				}else if j == len(status.Orders)-1{
					status.Orders = append(status.Orders[:],zeros[i])
				}
			}		
		}
	}else{
		status.Orders = zeros
	}



	// Find the 'cheapest' elevator
	for i:=0; i<len(status.Orders); i++{
		if status.Orders[0].Floor*status.Direction <= status.Floor*status.Direction{

			wrongDir = status.Orders[1:]
			status.Orders = append(wrongDir[:], status.Orders[0])

		}else{
			break
		}
	}


	return status
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




