package driver

import (
	"../variables"
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

func Init(newOrders chan variables.Order, jobDone, StopCh, ObsCh chan bool, nextFloor, currentFloor chan int){
	io_init()
	sensor := make(chan int,0)
	for i:=0; i<12; i++{		//Turn every light off.
		LightButtons(i, false)
	}



	go readSensor(sensor)
	go readButtons(newOrders,ObsCh,StopCh)
	go moveToFloor(nextFloor,currentFloor,sensor,jobDone,ObsCh)
/*
	var sleepStart chan int
	var sleepDone chan bool
	go sleeper(sleepStart,sleepDone)	
*/	
}

func moveToFloor(nextFloor chan int, currentFloor, sensor chan int, jobDone,ObsCh chan bool ){
	tempFloor := 1;
	target := 0;
	fmt.Println("moveToFloor: Initializing")

	for{
		select{
		case tempFloor = <-sensor:
			//fmt.Println("senor is registering: ", tempFloor)
			currentFloor <- tempFloor
			//fmt.Println("MoveToFloor: target=",target, "tempFloor=",tempFloor)
			if tempFloor == target {
				time.Sleep(time.Millisecond*150)
				motorHandler(0)
				// Åpner dører og venter 4 sek. DÅRLIG IMPLEMENTASJON
				LightButtons(9,true)
				time.Sleep(time.Second*2)
				LightButtons(9,false)
				jobDone <- true
			}

			

		case target = <- nextFloor:
			//fmt.Println("MoveToFloor: target=",target, "tempFloor=",tempFloor)
				//fmt.Println("nextfloor = target: ", target)
				if target < tempFloor{
				 	motorHandler(-1)	 
				}else if target > tempFloor { 
					motorHandler(1) 
				}else if target == tempFloor{ 
					//fmt.Println("tempFloor = target: ", target)
					motorHandler(0)
					// Åpner dører og venter 4 sek. DÅRLIG IMPLEMENTASJON
					LightButtons(9,true)

					time.Sleep(time.Second*2)
					LightButtons(9,false)
					jobDone <- true
				}

		}

		
	}
}

func LightButtons(light int, on bool){
	if on {	
		io_set_bit(lamp_channel_matrix[light])
	}else{ 
		io_clear_bit(lamp_channel_matrix[light]) 
	}
}

func readButtons( newOrders chan variables.Order, ObsCh chan bool, StopCh chan bool){ 
	
	var order variables.Order
	fmt.Println("readButtons: Initializing")
	for{
		
		for i := 0; i<12; i++{
			if io_read_bit(button_channel_matrix[i]) == 1{
				
				if i == 9 {

					ObsCh <- true			
				}else if i == 1{
					LightButtons(i,true)
					fmt.Println("EMERGENCY STOP!!")
					motorHandler(0)
					StopCh <- true
					return				
				}else{
					switch i {
						case 0 , 2:
							order.Floor = 0
							order.Dir = 1 - i/2
							newOrders <- order
							LightButtons(i,true)
						case 3 , 4 , 5 :
							order.Floor = 1
							if i == 3{
								order.Dir = 1
							}else if i == 4	{
								order.Dir = -1
							}else {
								order.Dir = 0
							}
							newOrders <- order
							LightButtons(i,true)
						case 6 , 7 , 8 :
							order.Floor = 2 
							if i == 6{
								order.Dir = 1
							}else if i == 7	{
								order.Dir = -1
							}else {
								order.Dir = 0
							}
							newOrders <- order
							LightButtons(i,true)
						case 10 , 11:
							order.Floor = 3
							order.Dir =  i - 11
							newOrders <- order
							LightButtons(i,true)
					}
						

						


					
					
					//fmt.Println("readButtons: Sending order" , order)
					
				}
					

				
			}
		}
		time.Sleep(time.Millisecond*80)
	}
}

func readSensor(sensor chan int){
	lastRead := -1
	current := -1
	fmt.Println("readSensor: Initializing")
	for {		
		if io_read_bit(variables.SENSOR_FLOOR1) == 1 	   {
			current = 0
			io_clear_bit(variables.LIGHT_FLOOR_IND1)//00
			io_clear_bit(variables.LIGHT_FLOOR_IND2)
			io_clear_bit(variables.LIGHT_COMMAND1)

		}else if io_read_bit(variables.SENSOR_FLOOR2) == 1 {
			current = 1
			io_clear_bit(variables.LIGHT_FLOOR_IND1)//01
			io_set_bit(variables.LIGHT_FLOOR_IND2)
			io_clear_bit(variables.LIGHT_COMMAND2)

		}else if io_read_bit(variables.SENSOR_FLOOR3) == 1 {
			current = 2
			io_set_bit(variables.LIGHT_FLOOR_IND1)//10
			io_clear_bit(variables.LIGHT_FLOOR_IND2)
			io_clear_bit(variables.LIGHT_COMMAND3)

		}else if io_read_bit(variables.SENSOR_FLOOR4) == 1 {
			current = 3
			io_set_bit(variables.LIGHT_FLOOR_IND1)//11
			io_set_bit(variables.LIGHT_FLOOR_IND2)
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

/*
func sleeper(sleepStart chan int,sleepDone chan bool){
	var time.Timer

}
*/