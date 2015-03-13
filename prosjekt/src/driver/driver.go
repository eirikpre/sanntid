package driver

import (
	"../variables"
	"time"
	"fmt"
)
/*
TODO

Door Open





*/

var button_channel_matrix = []int{ 
	variables.BUTTON_UP1, variables.STOP, variables.BUTTON_COMMAND1,
	variables.BUTTON_UP2, variables.BUTTON_DOWN2, variables.BUTTON_COMMAND2,
	variables.BUTTON_UP3, variables.BUTTON_DOWN3, variables.BUTTON_COMMAND3,
	variables.OBSTRUCTION, variables.BUTTON_DOWN4, variables.BUTTON_COMMAND4,
}

var lamp_channel_matrix  = []int{
    variables.LIGHT_UP1, variables.LIGHT_STOP, variables.LIGHT_COMMAND1,
    variables.LIGHT_UP2, variables.LIGHT_DOWN2, variables.LIGHT_COMMAND2,
    variables.LIGHT_UP3, variables.LIGHT_DOWN3, variables.LIGHT_COMMAND3,
    variables.LIGHT_DOOR_OPEN, variables.LIGHT_DOWN4, variables.LIGHT_COMMAND4,
}

func Init(nextFloor chan int, jobDone chan bool, newOrders chan variables.Order, StopCh chan bool, ObsCh chan bool,	currentFloor chan int){
	io_init()

	sensor := make(chan int,0)

	//nextFloor := make(chan int,1)

	
	for i:=0; i<12; i++{
		lightButtons(i, false)
	} 
	
	
	go readSensor(sensor)
	go readButtons(newOrders,ObsCh,StopCh)
	go moveToFloor(nextFloor,currentFloor,sensor,jobDone)		//waiting for the polling
	
}



func moveToFloor(nextFloor chan int, currentFloor chan int, sensor chan int, arrivedCh chan bool ){
	tempFloor := 0;
	target := 0;
	for{
		select{
		case tempFloor = <-sensor:
			currentFloor <- tempFloor
			fmt.Println("MoveToFloor: target=",target, "tempFloor=",tempFloor)
			if tempFloor == target{
				time.Sleep(time.Millisecond*150)
				motorHandler(0)
				arrivedCh <- true
				// Åpner dører og venter 10 sek. DÅRLIG IMPLEMENTASJON

			}

		case target = <- nextFloor:
			fmt.Println("MoveToFloor: target=",target, "tempFloor=",tempFloor)
			if target >= 0 && target < 4{
				
				if target < tempFloor{
				 	motorHandler(-1)	 
				}else if target > tempFloor { 
					motorHandler(1) 
				}else{ 	
				motorHandler(0) /*Error!!*/	
				}
			}else {
				fmt.Println("MoveToFloor :: illegal input")
				motorHandler(0)
				/*Error!!*/
			}
			
		}
		
	}
}




func lightButtons(light int, on bool){
	if on {	
		io_set_bit(lamp_channel_matrix[light])
	}else{ 
		io_clear_bit(lamp_channel_matrix[light]) 
	}
}

func readButtons( newOrders chan variables.Order, ObsCh chan bool, StopCh chan bool){ 
	lastRead := -1
	var order variables.Order
	for{
		
		for i := 0; i<12; i++{
			if io_read_bit(button_channel_matrix[i]) == 1{
				
				if i == 9 {
					ObsCh <- true				
				}else if i == 1{
					lightButtons(i,true)
					StopCh <- true				
				}else if lastRead != i{
					switch i {
						case 0 , 2:
							order.Floor = 0
							order.Dir = 1 - i/2
						case 3 , 4 , 5 :
							order.Floor = 1
							if i == 3{
								order.Dir = 1
							}else if i == 4	{
								order.Dir = -1
							}else {
								order.Dir = 0
							}
						case 6 , 7 , 8 :
							order.Floor = 2 
							if i == 6{
								order.Dir = 1
							}else if i == 7	{
								order.Dir = -1
							}else {
								order.Dir = 0
							}
						case 10 , 11:
							order.Floor = 3
							order.Dir = 11 - i
						


					}	
					lightButtons(i,true)
					lastRead = i
					fmt.Println("readButtons: Sending order" , order)
					newOrders <- order
				}
			}
		}
		time.Sleep(time.Millisecond*80)
	}
}


func readSensor(sensor chan int){
	lastRead := -1
	current := -1
	for {		
		if io_read_bit(variables.SENSOR_FLOOR1) == 1 	   {
			current = 0
			io_clear_bit(variables.LIGHT_FLOOR_IND1)//00
			io_clear_bit(variables.LIGHT_FLOOR_IND2)
			io_clear_bit(variables.LIGHT_UP1)
			io_clear_bit(variables.LIGHT_COMMAND1)

		}else if io_read_bit(variables.SENSOR_FLOOR2) == 1 {
			current = 1
			io_clear_bit(variables.LIGHT_FLOOR_IND1)//01
			io_set_bit(variables.LIGHT_FLOOR_IND2)
			io_clear_bit(variables.LIGHT_DOWN2)
			io_clear_bit(variables.LIGHT_UP2)
			io_clear_bit(variables.LIGHT_COMMAND2)

		}else if io_read_bit(variables.SENSOR_FLOOR3) == 1 {
			current = 2
			io_set_bit(variables.LIGHT_FLOOR_IND1)//10
			io_clear_bit(variables.LIGHT_FLOOR_IND2)
			io_clear_bit(variables.LIGHT_DOWN3)
			io_clear_bit(variables.LIGHT_UP3)
			io_clear_bit(variables.LIGHT_COMMAND3)

		}else if io_read_bit(variables.SENSOR_FLOOR4) == 1 {
			current = 3
			io_set_bit(variables.LIGHT_FLOOR_IND1)//11
			io_set_bit(variables.LIGHT_FLOOR_IND2)
			io_clear_bit(variables.LIGHT_DOWN4)
			io_clear_bit(variables.LIGHT_UP4)
			io_clear_bit(variables.LIGHT_COMMAND4)
		}
		
		if lastRead != current {

			lastRead = current
			sensor <- current
			
		}
		time.Sleep(time.Millisecond*50)
	}
}


func motorHandler(motorDir int ) {
	if (motorDir == 0){
			io_write_analog(variables.MOTOR, 0);
	    } else if (motorDir > 0) {
			io_clear_bit(variables.MOTORDIR);
			io_write_analog(variables.MOTOR, 2800);
	    } else if (motorDir < 0) {
			io_set_bit(variables.MOTORDIR);
			io_write_analog(variables.MOTOR, 2800);
		}
  	
}
