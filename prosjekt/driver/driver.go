package driver

import (
	"../variables"
	"bytes"
)

func Init(){
	if !(io_init()) {
	 	io_err = make(error,"io_init() trouble")
	}

}



func Get_buttons(chan error err_chan){
	buttons = [16]byte
	for {
		
		

		if io_read_bit(STOP) || io_read_bit(OBSTRUCTION)  {
			
			err_chan <- make(error, "
		}	
	
	}
}


func motorHandler(nextfloor chan int){

	if (dirn == 0){
        	io_write_analog(MOTOR, 0);
    	} else if (dirn > 0) {
		io_clear_bit(MOTORDIR);
		io_write_analog(MOTOR, 2800);
    	} else if (dirn < 0) {
		io_set_bit(MOTORDIR);
		io_write_analog(MOTOR, 2800);
    	}	



}
func makeNextFloor(buttons int, )


func moveToFloor()













