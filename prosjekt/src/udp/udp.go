package udp

import (
	"fmt"
	"net"
	"strconv"
	"json"
	"../variables"
)


func Udp_Init(UDPsendStatus, UDPreceiveStatus chan variables.Status){
	var receive chan []byte
	var status variables.Status
	
	go udp_send(UDPsendStatus)
	go udp_listen(UDPreceiveStatus)

			

		case b := <- receive :							// Receive from the router.

		}





}

func udp_listen(UDPreceiveStatus chan variables.Status){
var message []byte

	for{
		_, remoteaddr,_ := conn.ReadFromUDP(message)
		_ := json.Unmarshal(b, &status)
		UDPreceiveStatus <- status
		message -> receive
	}
}

func udp_send(UDPsendStatus chan variables.Status){
	
	
	for{
		  				// Send the status
		b,_ := json.Marshal(status)
		_,_ := conn.Write(b)

	}





}



// Fra Sindre Langeveld
/* func getOwnIP() net.IPAddr {
	ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
	    	// handle err
	    	for _, addr := range addrs {
			switch v := addr.(type) 
			{
				case *net.IPAddr:
			}
		    	// process IP address
		}
	}
}
*/




