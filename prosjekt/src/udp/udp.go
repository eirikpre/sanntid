package udp

import (
	"net"
	"encoding/json"
	"../variables"
	"fmt"
	"errors"
	
)


const port = 30000

func Udp_Init(UDPsendStatus, UDPreceiveStatus chan variables.Status,error_ch chan error){
	var conn *net.UDPConn

	defer func(){
		if r:=recover();r!=nil{
			error_ch <- errors.New("UDPRESET")
			conn.Close()
			return
		}
	}()

	fmt.Println("Udp_Init: Initialzing")
	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:30020")
	
	error_ch <- err

	conn, err = net.ListenUDP("udp4",addr)
	error_ch <- err


	go func() {
		for{	
				
				var msg_rcv []byte = make([]byte,2000)
				length,_,err := conn.ReadFromUDP(msg_rcv)
				if(err != nil){

					panic("CONNECTION DOWN!")
				}
				 
				msg_rcv = msg_rcv[:length]
				var status_rcv variables.Status	
				err = json.Unmarshal(msg_rcv,&status_rcv)
				error_ch <- err
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
				bytes,err := json.Marshal(&status_snd)
				error_ch <- err
				_,err = conn.WriteToUDP(bytes,addr)
				error_ch <- err
				
				
			}
	}()
}


func GetOwnIP(error_ch chan error) string {
     interfaces, err := net.InterfaceAddrs()
     error_ch <- err
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


