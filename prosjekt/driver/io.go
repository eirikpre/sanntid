package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func io_init() int {
	return C.io_init()
}

func  io_set_bit(channel int) { 
	C.io_set_bit(channel)
}

func io_clear_bit(channel int) { 
	C.io_clear_bit(channel)
}

func io_write_analog(channel, value int) {
	C.io_write_analog(channel,value)
}

func io_read_bit(channel int) int {
	return C.io_read_bit(channel)
}

func io_read_analog(channel int) int {
	return io_read_analog(channel)
}

