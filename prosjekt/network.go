package udp

import (
	"fmt"
	"net"
	"strconv"
)

//Only IP
var laddr IP		//local address
var baddr IP		//Broadcast address

type status struct{
	laddr IP		//Local IP
	floor int 		//Current floor
	direction [2]byte	//Up = 2, Still = 1, Down = 0
	destination int		//Destination floor
	Length int		//Length of received data, in #bytes
}

func UDPInit(){
	baddr, err = net.ResolveUDPAddr()
	if err != nil {
		return err
	}
	
}

func UDPSendWork(){		//Sending work

}

func UDPSendStatus(){		//Sending status

}

func UDPReceiveStatus(){	//Receive status

}

func UDPReceiveWork(){		//Receive work

}

func UDPConnection(){		
}
