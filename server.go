// Purpose: The virtual subsys in this model, running UDP server to receive packets from subsystems, and
//          send back to their UDP servers.
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

import (
	"fmt"
	"net"
)

// Server listens and forward UDP packets from each subsystem
type Server struct {
	Name string
	// recv conn from other subsys
	inConn *net.UDPConn
	// send conn to other subsys
	outConn net.Conn
	recvCnt int
	fwdCnt  int

	GateIn  *Gate
	GateOut [GATE_NUM_SERVER]*Gate
}

// returns a Server pointer
func NewServer(name string) *Server {
	var gateOut [GATE_NUM_SERVER]*Gate
	for i := 0; i < GATE_NUM_SERVER; i++ {
		gateOut[i] = NewGate(i, name)
	}
	return &Server{
		Name:    name,
		GateIn:  NewGate(0, name),
		GateOut: gateOut,
	}
}

func (s *Server) Start() {
	fmt.Printf("Start virtual subsys - %s: local_addr: %s, remote_addr: %s\n",
		s.Name, SUBSYS_TABLE[s.Name].LocalAddr, SUBSYS_TABLE[s.Name].RemoteAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", SUBSYS_TABLE[s.Name].LocalAddr)
	if err != nil {
		fmt.Println("invalid address")
		return
	}

	s.inConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.handlePacket()

	s.outConn, err = net.Dial("udp", SUBSYS_TABLE[s.Name].RemoteAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.handleMessage()

}

// receive external UDP packets
func (s *Server) handlePacket() {
	for {
		var buf [BUF_LEN]byte
		n, err := s.inConn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}

		// if n < 14 {
		// 	fmt.Printf("[%s] invalid packet\n", s.Name)
		// 	return
		// }

		s.recvCnt++
		fmt.Printf("[%s] Received packet #%d\n", s.Name, s.recvCnt)
		pkt := new(Packet)

		err = pkt.FromBuf(buf[0:n])
		if err != nil {
			fmt.Println(err)
			// invalid packet, drop
			continue
		}

		if pkt.Src != uint8(SUBSYS_TABLE[s.Name].ID) {
			fmt.Printf("[%s]WARNING! SRC doesn't match\n", s.Name)
		}

		// the ground-habitat link
		if (s.Name == "HMS" && pkt.Dst == 0) || (s.Name == "GCC" && pkt.Dst == 1) {
			pkt.Path = append(pkt.Path, s.Name)
			s.GateOut[1].Channel <- pkt
			var dst = "HMS"
			if s.Name == "HMS" {
				dst = "GCC"
			}
			LogsComm <- Log{
				Type: 0,
				Msg:  fmt.Sprintf("[%s] Received and forwarding to %s", s.Name, dst),
			}
		} else {
			pkt.Path = append(pkt.Path, s.Name)
			s.GateOut[0].Channel <- pkt
		}

	}
}

// send internal messages from switches to outside
func (s *Server) handleMessage() {
	// fmt.Println("waiting msg from switches")
	for {
		pkt := <-s.GateIn.Channel
		pkt.Path = append(pkt.Path, s.Name)
		s.fwdCnt++
		fmt.Println(pkt.RawBytes)
		_, err := s.outConn.Write(pkt.RawBytes)
		if err != nil {
			fmt.Printf("[%s] sending UDP to remote error %v\n", s.Name, err)
		}
		fwdCntTotal++
		// fmt.Printf("Forward packet #%d {", fwdCntTotal)
		// for i, v := range pkt.Path {

		// 	fmt.Print(v)
		// 	if i != len(pkt.Path)-1 {
		// 		fmt.Print("->")
		// 	}
		// }

		// fmt.Printf("}\n")

		LogsComm <- Log{
			Type: 0,
			Msg:  fmt.Sprintf("[%s] Forwarded packet", s.Name),
		}
	}
}
