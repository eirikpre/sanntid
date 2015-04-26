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
	
	
	

	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:30020")
	if(err != nil){
		fmt.Println("Error in resovle: ", err)
	}
	
	conn, err := net.ListenUDP("udp4",addr)
	if(err != nil){
		fmt.Println("Error in ListenUDP: ", err)
	}


	go func() {
		for{	
				
				var msg_rcv []byte = make([]byte,2000)
				length,_,_ := conn.ReadFromUDP(msg_rcv)
				msg_rcv = msg_rcv[:length]
				var status_rcv variables.Status	
				err := json.Unmarshal(msg_rcv,&status_rcv)
				if(err != nil){
					fmt.Println("error i unmarshal: ", err)
				}
				//fmt.Println("udp_received: ",status_rcv)
				UDPreceiveStatus <- status_rcv
		}
	}()

	newAddr := new(net.UDPAddr)
	*newAddr = *addr
	newAddr.IP = make(net.IP, len(addr.IP))
	copy(newAddr.IP,addr.IP)	

	go func() {
		for{	
				status_snd := <- UDPsendStatus
				//fmt.Println("udp_send: ",status_snd)
				
				bytes,err := json.Marshal(&status_snd)
				if(err != nil){
					fmt.Println("error i Marshal: ", err)
				}

				_,err = conn.WriteToUDP(bytes,addr)
				if(err != nil){
					fmt.Println("Error in write: ", err)
				}
				//time.Sleep(time.Millisecond*300)
				
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


