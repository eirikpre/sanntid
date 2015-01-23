// Go 1.2
// go run helloworld_go.go

package main

import (
	. "fmt" // Using '.' to avoid prefixing functions with their package names
		// This is probably not a good idea for large projects...
	"runtime"
	//"time"
)

func goroutine_1(channel chan int, channelDone chan int) {		
	for j := 0; j < 1000000; j++{
		i = <- channel
		//i = k
		i ++
		channel <- i
	}
	channelDone <- 1
}

func goroutine_2(channel chan int, channelDone chan int){
	for j :=0;j < 1000001;j++{
		i = <- channel
		i --
		channel <- i
	}
	channelDone <- 1
}


var i int = 0
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
	var channel chan int = make(chan int,1)
	var channelDone chan int = make(chan int,1)
	
	channel <- 0
		// Try doing the exercise both with and without it!
	go goroutine_1(channel, channelDone)
	go goroutine_2(channel, channelDone)
	 // This spawns someGoroutine() as a goroutine

	// We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
	// We'll come back to using channels in Exercise 2. For now: Sleep.
	//time.Sleep(900*time.Millisecond)

	_ = <- channelDone 
	_ = <- channelDone
	i = <- channel		
	Println(i,"\n")
	
}



