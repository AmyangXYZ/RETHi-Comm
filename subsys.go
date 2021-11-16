// Purpose: The virtual subsys in this model, running UDP Subsys to receive packets from subsystems, and
//          send back to their UDP servers.
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

import (
	"errors"
	"fmt"
	"net"
)

// Subsys listens and forward UDP packets from each subsystem
type Subsys struct {
	name string
	id   int
	// recv conn from other node
	inConn *net.UDPConn
	// send conn to other node
	outConn net.Conn
	recvCnt int
	fwdCnt  int

	gatesIn     [GATE_NUM_SUBSYS]*Gate
	gatesOut    [GATE_NUM_SUBSYS]*Gate
	gatesInIdx  int
	gatesOutIdx int
}

// returns a Subsys pointer
func NewSubsys(name string) *Subsys {
	id := 0

	var gatesIn [GATE_NUM_SUBSYS]*Gate
	var gatesOut [GATE_NUM_SUBSYS]*Gate
	for i := 0; i < GATE_NUM_SUBSYS; i++ {
		gatesIn[i] = NewGate(i, name)
		gatesOut[i] = NewGate(i, name)
	}
	for _, v := range SUBSYS_LIST {
		if name == v.Name {
			id = v.ID
			break
		}
	}

	s := &Subsys{
		name:        name,
		id:          id,
		gatesIn:     gatesIn,
		gatesOut:    gatesOut,
		gatesInIdx:  -1,
		gatesOutIdx: -1,
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
	s.gatesOutIdx++
	return s.gatesOut[s.gatesOutIdx]
}

// implement Node interface
func (s *Subsys) InGate() *Gate {
	s.gatesInIdx++
	return s.gatesIn[s.gatesInIdx]
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

	for _, g := range s.gatesIn {
		go s.handleMessage(g)
	}

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

		// fmt.Printf("[%s] Received packet #%d\n", s.name, s.recvCnt)
		pkt := new(Packet)
		err = pkt.FromBuf(buf[0:n])
		if err != nil {
			fmt.Println(err)
			continue
		}
		if pkt.Src != uint8(SUBSYS_TABLE[s.name].ID) {
			fmt.Printf("[%s]WARNING! SRC doesn't match\n", s.name)
		}

		if g, err := s.routing(pkt); err == nil {
			// fmt.Println("sent to", g.Neighbor)
			g.Channel <- pkt
			s.fwdCnt++
		} else {
			fmt.Println(err)
		}
	}
}

// send internal messages from switches to outside
func (s *Subsys) handleMessage(inGate *Gate) {
	// fmt.Println("waiting msg from switches")
	for {
		pkt := <-inGate.Channel
		pkt.Path = append(pkt.Path, s.name)

		s.recvCnt++
		// fmt.Println(pkt.RawBytes)
		if !pkt.IsSim {
			_, err := s.outConn.Write(pkt.RawBytes)
			if err != nil {
				fmt.Printf("[%s] sending UDP to remote error %v\n", s.name, err)
			}
		}

		FwdCntTotal++
		if pkt.Delay < 1 {
			pkt.Delay *= 1000000
			fmt.Printf("Pkt #%d: %d bytes, %v, delay: %.3f us\n", FwdCntTotal, len(pkt.RawBytes), pkt.Path, pkt.Delay)
			LogsComm <- Log{
				Type: 0,
				Msg:  fmt.Sprintf("Pkt #%d: %d bytes, %v, delay: %.3f us", FwdCntTotal, len(pkt.RawBytes), pkt.Path, pkt.Delay),
			}
		} else {
			fmt.Printf("Pkt #%d: %d bytes, %v, delay: %.2f s\n", FwdCntTotal, len(pkt.RawBytes), pkt.Path, pkt.Delay)
			LogsComm <- Log{
				Type: 0,
				Msg:  fmt.Sprintf("Pkt #%d: %d bytes, %v, delay: %.2f s", FwdCntTotal, len(pkt.RawBytes), pkt.Path, pkt.Delay),
			}
		}
	}
}

func (s *Subsys) CreateFlow(dst int) {
	pkt := &Packet{
		Src:   uint8(s.id),
		Dst:   uint8(dst),
		IsSim: true,
	}
	var buf [2]byte
	buf[0] = pkt.Src
	buf[1] = pkt.Dst
	pkt.RawBytes = buf[:]

	if g, err := s.routing(pkt); err == nil {
		// fmt.Println("sent to", g.Neighbor)
		pkt.Path = append(pkt.Path, s.name)
		g.Channel <- pkt
		s.fwdCnt++
	} else {
		fmt.Println(err)
	}
}

func (s *Subsys) routing(pkt *Packet) (*Gate, error) {
	if pkt.Dst == 0 || s.name == "GCC" { // to/from GCC
		for _, g := range s.gatesOut {
			if g.Neighbor == SUBSYS_LIST[pkt.Dst].Name {
				return g, nil
			}
		}
	}

	for _, g := range s.gatesOut {
		for _, sw := range ROUTING_TABLE[int(pkt.Dst)] {
			if g.Neighbor == sw {
				return g, nil
			}
		}
	}

	// not found, sent to an arbitrary switch to reach sw0
	// prefer switch with consistent id
	for _, g := range s.gatesOut {
		if g.Neighbor == ROUTING_TABLE[s.id][0] {
			return g, nil
		}
	}

	return nil, errors.New("[" + s.name + "] cannot found next hop")
}
