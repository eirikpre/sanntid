package udp

import (
	"net"
	"encoding/json"
	"../variables"
	"fmt"
	//"time"
	
)

const port = 30000

func Udp_Init(UDPsendStatus, UDPreceiveStatus chan variables.Status){
	fmt.Println("Udp_Init: Initialzing")
	var msg_rcv []byte = make([]byte,200)
	
	var status variables.Status

	addr, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:30000")
	conn, _ := net.ListenUDP("udp4",addr)


	go func() {
		for{	
				length,_,_ := conn.ReadFromUDP(msg_rcv)
				msg_rcv = msg_rcv[:length]
				json.Unmarshal(msg_rcv, &status)
				fmt.Println("udp_received: ", msg_rcv)
				UDPreceiveStatus <- status
		}
	}()

	newAddr := new(net.UDPAddr)
	*newAddr = *addr
	newAddr.IP = make(net.IP, len(addr.IP))
	copy(newAddr.IP,addr.IP)	

	go func() {
		for{
				status = <- UDPsendStatus
				message,_ := json.Marshal(status)
				fmt.Println("udp_send: ", message)
				conn.WriteToUDP(message,addr)
			}
	}()
}


func GetOwnIP() string {
     interfaces, _ := net.InterfaceAddrs()
     for _, address := range interfaces {

 		ipnet, ok := address.(*net.IPNet)

       	if  ok && !ipnet.IP.IsLoopback(){

          	if ipnet.IP.To4() != nil {

             	return ipnet.IP.String()

          	}
      	}
     }
     panic("GetOwnIP")
}


