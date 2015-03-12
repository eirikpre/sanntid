package driver

import (
	"variables"
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

func Init(){
	io_init()
	sensor := make(chan int,1)
	buttons := make(chan int,1)
	currentFloor := make(chan int,50)
	//nextFloor := make(chan int,1)
	
	for i:=0; i<12; i++{
		lightButtons(i, false)
	} 
	
	buttons <- 0
	
	go readSensor(sensor)
	go readButtons(buttons)
	go MoveToFloor(nextFloor,currentFloor,sensor)		//waiting for the polling
	
	
}

func MoveToFloor(nextFloor chan int, currentFloor chan int, sensor chan int){
	tempFloor := 2;
	target := 5;
	for{
		select{
		case tempFloor = <-sensor:
			currentFloor <- tempFloor
			if tempFloor == target{
				motorHandler(0)
				
									
			}
		case target = <- nextFloor:
			fmt.Println("the target is:",target)
			if target >= 0 && target < 4{
				fmt.Println("legal input")
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

func readButtons( buttons chan int){
	lastRead := -1
	
	for{
		
		for i := 0; i<12; i++{
			if io_read_bit(button_channel_matrix[i]) == 1{
				if lastRead != i{
					
					lightButtons(i,true)
					lastRead = i
					buttons <- i
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
			fmt.Println("the sensor is: <---",current)
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
