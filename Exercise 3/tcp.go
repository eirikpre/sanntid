package main

import(	"fmt"
	"net"
)

//port 34933 fixed
//PORT 33546
func main(){
	//str := "Melding tilbake"
	//s := make([]byte,1024)
	//copy(s[:],str)
	
	p := make([]byte, 1024)
	conn, err := net.Dial("tcp", "129.241.187.136:34933")	
	if(err != nil){
		fmt.Println("err")
	}
	
	conn.Read(p)
	fmt.Printf("%s\n",p)
	for j := 0; j < 10; j++{
		str := "Melding tilbake"
		s := make([]byte,1024)
		copy(s[:],str)
		conn.Write(s)
		conn.Read(p)
		fmt.Printf("%s\n\n",p)
	}
	conn.Close()
}
