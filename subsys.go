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
	"os"
	"sort"
	"sync"
	"time"
)

// Subsys listens and forward UDP packets from each subsystem
type Subsys struct {
	name     string
	position [2]int
	id       int
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

	RoutingTable map[string][]RoutingEntry

	SeqRecoverHistory      map[int32]bool // for frer
	SeqRecoverHistoryMutex sync.Mutex

	logMutex sync.Mutex

	stopSig chan bool
}

// returns a Subsys pointer
func NewSubsys(name string, position [2]int) *Subsys {

	var gatesIn [GATE_NUM_SUBSYS]*Gate
	var gatesOut [GATE_NUM_SUBSYS]*Gate
	for i := 0; i < GATE_NUM_SUBSYS; i++ {
		gatesIn[i] = NewGate(i, name)
		gatesOut[i] = NewGate(i, name)
	}

	s := &Subsys{
		name:              name,
		position:          position,
		id:                subsysName2ID(name),
		gatesIn:           gatesIn,
		gatesOut:          gatesOut,
		gatesInIdx:        -1,
		gatesOutIdx:       -1,
		RoutingTable:      make(map[string][]RoutingEntry),
		SeqRecoverHistory: make(map[int32]bool),
		stopSig:           make(chan bool),
	}

	for subsys, _ := range SUBSYS_MAP {
		if subsys == name {
			continue
		}
		paths := Graph.FindAllPaths(name, subsys)
		table := []RoutingEntry{}
		for _, p := range paths {
			entry := RoutingEntry{NextHop: p[1], HopCount: len(p) - 1}
			found := false
			for _, e := range table {
				if e.NextHop == entry.NextHop {
					found = true
					if e.HopCount > entry.HopCount {
						e.HopCount = entry.HopCount
					}
				}
			}
			if !found {
				table = append(table, entry)
			}
		}
		sort.SliceStable(table, func(i, j int) bool {
			return table[i].HopCount < table[j].HopCount
		})
		s.RoutingTable[subsys] = table
	}
	// fmt.Println(name)
	// for dst, p := range s.RoutingTable {
	// 	fmt.Println("    ", dst, p)
	// }
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
	fmt.Printf("Start virtual subsys (id: %d) - %s: local_addr: %s, remote_addr: %s\n",
		s.id, s.name, os.Getenv("ADDR_LOCAL_"+s.name), os.Getenv("ADDR_REMOTE_"+s.name))
	udpAddr, err := net.ResolveUDPAddr("udp", os.Getenv("ADDR_LOCAL_"+s.name))
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

	s.outConn, err = net.Dial("udp", os.Getenv("ADDR_REMOTE_"+s.name))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, g := range s.gatesIn {
		go s.handleMessage(g)
	}
}

func (s *Subsys) Stop() {
	for range s.gatesIn {
		s.stopSig <- true // stop ingates
	}
	s.inConn.Close() // stop udp server
	// fmt.Println(s.name, "stopped")
}

// receive external UDP packets
func (s *Subsys) handlePacket() {
	for {
		var buf [PKT_BUF_LEN]byte
		n, err := s.inConn.Read(buf[0:])
		if err != nil {
			// fmt.Println(err)
			return
		}

		go func(buffer []byte) {
			pkt := new(Packet)
			err = pkt.FromBuf(buf[0:n])
			if err != nil {
				fmt.Println(err)
				return
			}
			pkt.Seq = getSeqNum()
			pkt.RxTimestamp = time.Now().UnixNano()
			// fmt.Println("packet recv #", pkt.Seq, pkt.RxTimestamp)
			if pkt.Src != uint8(s.id) {
				fmt.Printf("[%s]WARNING! SRC doesn't match\n", s.name)
			}

			if g, err := s.routing(pkt); err == nil {
				s.send(pkt, g)
			} else {
				fmt.Println(err)
			}
		}(buf[0:n])

	}
}

// send internal messages from switches to outside
func (s *Subsys) handleMessage(inGate *Gate) {
	// fmt.Println(s.name, "waiting msg from switches")
	for {
		select {
		case <-s.stopSig:
			fmt.Println(s.name, "terminate an ingate goroutine")
			return
		case pkt := <-inGate.Channel:
			s.recvCnt++
			if FRER_ENABLED {
				// eliminate dup
				s.SeqRecoverHistoryMutex.Lock()
				if _, ok := s.SeqRecoverHistory[pkt.Seq]; ok {
					// fmt.Println(sw.name, "eliminate dup from", pkt.Path[len(pkt.Path)-1])
					s.SeqRecoverHistoryMutex.Unlock()
					continue
				}
				s.SeqRecoverHistory[pkt.Seq] = true
				s.SeqRecoverHistoryMutex.Unlock()
			}
			pkt.Path = append(pkt.Path, s.name)

			pkt.TxTimestamp = time.Now().UnixNano()
			fmt.Println(s.name, "packet send out #", pkt.Seq, pkt.TxTimestamp)

			// fmt.Println(pkt.RawBytes)
			if !pkt.IsSim {
				go func() {
					_, err := s.outConn.Write(pkt.RawBytes)
					if err != nil {
						fmt.Printf("[%s] sending UDP to remote error %v\n", s.name, err)
					}
				}()
			}

			// go saveStatsDelay(s.name, subsysID2Name(pkt.Src), pkt.Seq, pkt.Delay)

			if pkt.Delay < 1 {
				pkt.Delay *= 1000000
				// fmt.Printf("Pkt #%d: %d bytes, %v, delay: %.3f us\n", pkt.Seq, len(pkt.RawBytes), pkt.Path, pkt.Delay)
				// fmt.Printf("Pkt #%d: %d bytes, %v, delay: %v us\n", pkt.Seq, len(pkt.RawBytes), pkt.Path, (pkt.TxTimestamp-pkt.RxTimestamp)/1000)
				if CONSOLE_ENABLED {
					WSLog <- Log{
						Type: WSLOG_MSG,
						Msg:  fmt.Sprintf("Pkt #%d: %d bytes, %v, delay: %v us", pkt.Seq, len(pkt.RawBytes), pkt.Path, (pkt.TxTimestamp-pkt.RxTimestamp)/1000),
						// Msg:  fmt.Sprintf("Pkt #%d: %d bytes, %v, delay: %.3f us", pkt.Seq, len(pkt.RawBytes), pkt.Path, pkt.Delay),
					}
				}
			} else {
				fmt.Printf("Pkt #%d: %d bytes, %v, delay: %.3f us\n", pkt.Seq, len(pkt.RawBytes), pkt.Path, pkt.Delay)
				if CONSOLE_ENABLED {
					WSLog <- Log{
						Type: WSLOG_MSG,
						Msg:  fmt.Sprintf("Pkt #%d: %d bytes, %v, delay: %.2f s", pkt.Seq, len(pkt.RawBytes), pkt.Path, pkt.Delay),
					}
				}
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
	var buf [64]byte
	buf[0] = pkt.Src
	buf[1] = pkt.Dst
	pkt.RawBytes = buf[:]
	pkt.Seq = getSeqNum()
	if g, err := s.routing(pkt); err == nil {
		s.send(pkt, g)
	} else {
		fmt.Println(err)
	}
}

func (s *Subsys) routing(pkt *Packet) (*Gate, error) {
L1:
	for _, entry := range s.RoutingTable[subsysID2Name(pkt.Dst)] {
		if entry.NextHop[:2] == "SW" {
			for _, swww := range Switches {
				if swww.name == entry.NextHop && swww.Failed {
					continue L1
				}
			}
		}
		for _, g := range s.gatesOut {
			if g.Neighbor == entry.NextHop {
				return g, nil
			}
		}
	}

	return nil, errors.New("[" + s.name + "] cannot found next hop")
}

func (s *Subsys) send(pkt *Packet, gate *Gate) {
	// fmt.Println("sent to", gate.Neighbor)
	pkt.Path = append(pkt.Path, s.name)
	gate.Channel <- pkt
	if ANIMATION_ENABLED {
		WSLog <- Log{
			Type:  WSLOG_PKT_TX,
			PktTx: [2]string{s.name, gate.Neighbor},
		}
	}
	s.logMutex.Lock()
	s.fwdCnt++
	s.logMutex.Unlock()
}
