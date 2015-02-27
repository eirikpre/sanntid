package driver

import (
	"variables"
	"time"
	"fmt"
)
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
	driveMotor := make(chan int,0)
	sensor := make(chan int,0)
	buttons := make(chan int,0)
	for i:=0; i<12; i++{
		lightButtons(i, false)
	}

	motorHandler(driveMotor) 
	go readSensor(sensor)
	go readButtons(buttons)
	go elevatorHandler(status)
	
	
	
	//test case
	for{
		select {
			case  lol := <- sensor:
				fmt.Println("Now at floor:",lol)
			case  button := <- buttons:
				fmt.Println("Button:",button)
				lightButtons(button, true)
		}
		time.Sleep(time.Millisecond*200)
	}	
/*	driveMotor <- 1
	time.Sleep(2*time.Second)
	driveMotor <- 0
	time.Sleep(2*time.Second)*/
}

func MoveToFloor(nextFloor chan int, done chan bool ){
	currentFloor := 0
	if( buttonLight = true && heis != currentfloor){
		nextFloor = buttonLightFloor
		if(sensor != currentFloor){
			if(nextFloor > currentFloor){
				motorhandler(1)
			}else if{ nextFloor < currentFloor){
				motorhandler(-1)
			}
		}
		else{
			motorhandler(0)
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
	for{
		for i := 0; i<12; i++{
			if io_read_bit(button_channel_matrix[i]) == 1{
				buttons <- i
			}
		}/*
		if io_read_bit(STOP) || io_read_bit(OBSTRUCTION)  {
			
			err_chan <- make(error, "obstruction/stop")
		}*/
		time.Sleep(time.Millisecond*80)
	}
}


func readSensor(sensor chan int){
	for {
		if io_read_bit(variables.SENSOR_FLOOR1) == 1 	   {
			sensor <- 0
			io_clear_bit(variables.LIGHT_FLOOR_IND1)//00
			io_clear_bit(variables.LIGHT_FLOOR_IND2)
		}else if io_read_bit(variables.SENSOR_FLOOR2) == 1 {
			sensor <- 1
			io_clear_bit(variables.LIGHT_FLOOR_IND1)//01
			io_set_bit(variables.LIGHT_FLOOR_IND2)
		}else if io_read_bit(variables.SENSOR_FLOOR3) == 1 {
			sensor <- 2
			io_set_bit(variables.LIGHT_FLOOR_IND1)//10
			io_clear_bit(variables.LIGHT_FLOOR_IND2)
		}else if io_read_bit(variables.SENSOR_FLOOR4) == 1 {
			sensor <- 3
			io_set_bit(variables.LIGHT_FLOOR_IND1)//11
			io_set_bit(variables.LIGHT_FLOOR_IND2)
		}
		time.Sleep(time.Millisecond*200)
	}
}



func motorHandler(motorDir int ) {
	for {
		dirn := <- motorDir
		if (dirn == 0){
			io_write_analog(variables.MOTOR, 0);
	    	} else if (dirn > 0) {
			io_clear_bit(variables.MOTORDIR);
			io_write_analog(variables.MOTOR, 2800);
	    	} else if (dirn < 0) {
			io_set_bit(variables.MOTORDIR);
			io_write_analog(variables.MOTOR, 2800);
			}
    	}	
}
