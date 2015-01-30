//port 20018

package main

import ("net"
	"fmt"
)






func main(){
	//p := make([]byte, 1024)
	conn, _ := net.Dial("udp", "129.241.187.255:20008")
	for j := 0; j > -200; j++{
	fmt.Fprintf(conn, "GET RECT")
	}
	conn.Close()
	
	//addr := net.UDPAddr{
	//	Port: 20018,
	//	IP: net.ParseIP("129.241.187.255"),
	//}
	//msg, _ :=net.ListenUDP("udp", &addr)

	//_, remoteaddr,_ := msg.ReadFromUDP(p)
	//fmt.Printf("Message read from %v %s\n", remoteaddr, p)
	

}


//broadcastIP = #.#.#.255. First three bytes are from the local IP
//port = 30000
//addr = new InternetAddress(broadcastIP, port)
//sendSock = new Socket(udp) // UDP, aka SOCK_DGRAM
//sendSock.setOption(broadcast, true)
//sendSock.sendTo(message, addr)
