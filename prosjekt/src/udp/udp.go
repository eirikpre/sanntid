package udp

import (
	"fmt"
	"net"
	"strconv"
	"encoding/json"
	"../variables"
)
port := 30000;

func Udp_Init(UDPsendStatus, UDPreceiveStatus chan variables.Status){
	
	go udp_send(UDPsendStatus)
	go udp_listen(UDPreceiveStatus)

}

func udp_listen(UDPreceiveStatus chan variables.Status){
	addr, _ := net.ResolveUDPAddr("udp", ":port")
	conn, _ := net.ListenUDP("udp",addr)
	
	var status variables.Status
	var message []byte
	
	for{
		_, _,_ := conn.ReadFromUDP(message)
		_ := json.Unmarshal(message, &status)
		UDPreceiveStatus <- status
	}
}

func udp_send(UDPsendStatus chan variables.Status){
	broadcast := net,IPv4(255.255.255.255)
	broadcast = net.UDPAddr{ IP: broadcast, Port: port }
	conn, _ := net.DialUDP("udp",nil,&broadcast)
	
	for{
		select{
		case status := <- UDPsendStatus 
			b,_ := json.Marshal(status)
			_,_ := conn.Write(b)
		}
	}
}

func GetOwnIP() string{
         interfaces, _ := net.InterfaceAddrs()
         for _, address := range interfaces {
               if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                  if ipnet.IP.To4() != nil {
                     return (ipnet.IP.String())
                  }

               }
         }
}


