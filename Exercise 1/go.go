// Go 1.2
// go run helloworld_go.go

package main

import (
	. "fmt" // Using '.' to avoid prefixing functions with their package names
		// This is probably not a good idea for large projects...
	"runtime"
	"time"
)

func goroutine_1() {	
	for j := 0; j < 1000000; j++{
		i ++
	}
}

func goroutine_2(){
	for j :=0;j < 1000000;j++{
		i --
	}
}

var i int = 0

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // I guess this is a hint to what GOMAXPROCS does...
						// Try doing the exercise both with and without it!
	go goroutine_1()
	go goroutine_2()
	 // This spawns someGoroutine() as a goroutine

	// We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
	// We'll come back to using channels in Exercise 2. For now: Sleep.
	time.Sleep(100*time.Millisecond)
	Println(i)
}



