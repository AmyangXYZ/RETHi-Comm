// Purpose: The virtual subsys in this model, running UDP Subsys to receive packets from subsystems, and
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

// Subsys listens and forward UDP packets from each subsystem
type Subsys struct {
	name string
	// recv conn from other node
	inConn *net.UDPConn
	// send conn to other node
	outConn net.Conn
	recvCnt int
	fwdCnt  int

	gateIn     *Gate
	gateOut    [GATE_NUM_SUBSYS]*Gate
	gateOutIdx int
}

// returns a Subsys pointer
func NewSubsys(name string) *Subsys {
	var gateOut [GATE_NUM_SUBSYS]*Gate
	for i := 0; i < GATE_NUM_SUBSYS; i++ {
		gateOut[i] = NewGate(i, name)
	}
	s := &Subsys{
		name:       name,
		gateIn:     NewGate(0, name),
		gateOut:    gateOut,
		gateOutIdx: -1,
	}
	go s.Start()
	Subsystems = append(Subsystems, s)
	return s
}

// implement Node interface
func (s *Subsys) Name() string {
	return s.name
}

// implement Node interface
func (s *Subsys) OutGate() *Gate {
	s.gateOutIdx++
	return s.gateOut[s.gateOutIdx]
}

// implement Node interface
func (s *Subsys) InGate() *Gate {
	return s.gateIn
}

func (s *Subsys) Start() {
	fmt.Printf("Start virtual subsys - %s: local_addr: %s, remote_addr: %s\n",
		s.name, SUBSYS_TABLE[s.name].LocalAddr, SUBSYS_TABLE[s.name].RemoteAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", SUBSYS_TABLE[s.name].LocalAddr)
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

	s.outConn, err = net.Dial("udp", SUBSYS_TABLE[s.name].RemoteAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.handleMessage()

}

// receive external UDP packets
func (s *Subsys) handlePacket() {
	for {
		var buf [BUF_LEN]byte
		n, err := s.inConn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}

		// if n < 14 {
		// 	fmt.Printf("[%s] invalid packet\n", s.name)
		// 	return
		// }

		s.recvCnt++
		fmt.Printf("[%s] Received packet #%d\n", s.name, s.recvCnt)
		pkt := new(Packet)

		err = pkt.FromBuf(buf[0:n])
		if err != nil {
			fmt.Println(err)
			// invalid packet, drop
			continue
		}

		if pkt.Src != uint8(SUBSYS_TABLE[s.name].ID) {
			fmt.Printf("[%s]WARNING! SRC doesn't match\n", s.name)
		}
		pkt.Path = append(pkt.Path, s.name)
		fmt.Println("biu")
		// routing
		foundGate := false
		for _, g := range s.gateOut {
			if g.Neighbor == SUBSYS_LIST[pkt.Dst].Name {
				fmt.Println("sent to gate", g)
				g.Channel <- pkt
				foundGate = true
				break
			}
		}

		if !foundGate { // not hms or gcc
			for _, g := range s.gateOut {
				if g.Neighbor[:2] == "SW" {
					fmt.Println("sent to gate", g)
					g.Channel <- pkt
					break
				}
			}
		}
	}
}

// send internal messages from switches to outside
func (s *Subsys) handleMessage() {
	// fmt.Println("waiting msg from switches")
	for {
		pkt := <-s.gateIn.Channel
		pkt.Path = append(pkt.Path, s.name)
		s.fwdCnt++
		// fmt.Println(pkt.RawBytes)
		_, err := s.outConn.Write(pkt.RawBytes)
		if err != nil {
			fmt.Printf("[%s] sending UDP to remote error %v\n", s.name, err)
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
		fmt.Printf("[%s] Forwarded packet %v\n", s.name, pkt.Path)
		LogsComm <- Log{
			Type: 0,
			Msg:  fmt.Sprintf("[%s] Forwarded packet %v", s.name, pkt.Path),
		}
	}
}
