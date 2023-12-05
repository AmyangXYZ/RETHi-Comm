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

// Subsys is the virtual node that represents a subsystem,
// it communicates with outside real subsystem and pass packets to TSN switches
type Subsys struct {
	name     string
	position [2]int
	id       int
	// recv conn from other node
	inConn  *net.UDPConn
	recvCnt int
	fwdCnt  int

	portsIn     [PORT_NUM_SUBSYS]*Port
	portsOut    [PORT_NUM_SUBSYS]*Port
	portsInIdx  int
	portsOutIdx int

	Priority int // will overwrite priority in packets

	RoutingTable map[string][]RoutingEntry

	SeqRecoverHistory      map[int32]bool // for frer
	SeqRecoverHistoryMutex sync.Mutex

	logMutex sync.Mutex

	stopSig chan bool
}

// returns a Subsys pointer
func NewSubsys(name string, position [2]int) *Subsys {

	var portsIn [PORT_NUM_SUBSYS]*Port
	var portsOut [PORT_NUM_SUBSYS]*Port
	for i := 0; i < PORT_NUM_SUBSYS; i++ {
		portsIn[i] = NewPort(i, name)
		portsOut[i] = NewPort(i, name)
	}

	s := &Subsys{
		name:              name,
		position:          position,
		id:                subsysName2ID(name),
		Priority:          1,
		portsIn:           portsIn,
		portsOut:          portsOut,
		portsInIdx:        -1,
		portsOutIdx:       -1,
		RoutingTable:      make(map[string][]RoutingEntry),
		SeqRecoverHistory: make(map[int32]bool),
		stopSig:           make(chan bool),
	}

	for subsys := range SUBSYS_MAP {
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
	go s.Start()
	Subsystems = append(Subsystems, s)
	return s
}

// implement Node interface
func (s *Subsys) Name() string {
	return s.name
}

// implement Node interface
func (s *Subsys) OutPort() *Port {
	s.portsOutIdx++
	return s.portsOut[s.portsOutIdx]
}

// implement Node interface
func (s *Subsys) InPort() *Port {
	s.portsInIdx++
	return s.portsIn[s.portsInIdx]
}

// Start the Subsys
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

	for _, p := range s.portsIn {
		go s.handleMessage(p)
	}
}

// Stop the Subsys
func (s *Subsys) Stop() {
	for range s.portsIn {
		s.stopSig <- true // stop inports
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
			pkt.Priority = uint8(s.Priority)
			pkt.Seq = getSeqNum()
			pkt.UID = getUID()
			pkt.RxTimestamp = time.Now().UnixNano()
			// fmt.Println("packet recv #", pkt.Seq, pkt.RxTimestamp)
			if pkt.Src != uint8(s.id) {
				fmt.Printf("[%s]WARNING! SRC doesn't match\n", s.name)
			}
			if ANIMATION_ENABLED {
				WSLog <- Log{
					Type:  WSLOG_PKT_TX,
					PktTx: PktTx{Node: s.name, UID: pkt.UID},
				}
			}
			if p, err := s.routing(pkt); err == nil {
				s.send(pkt, p)
			} else {
				fmt.Println(err)
			}
		}(buf[0:n])

	}
}

// send internal messages from switches to outside
func (s *Subsys) handleMessage(inPort *Port) {
	// fmt.Println(s.name, "waiting msg from switches")
	for {
		select {
		case <-s.stopSig:
			// fmt.Println(s.name, "terminate an inport goroutine")
			return
		case pkt := <-inPort.Channel:
			s.recvCnt++
			if FRER_ENABLED || DUP_ELI_ENABLED {
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

			if s.name == "HMS" &&
				(pkt.Dst == uint8(subsysName2ID("GCC")) ||
					(pkt.Src == uint8(subsysName2ID("GCC")) && pkt.Dst != uint8(subsysName2ID("HMS")))) {
				if p, err := s.routing(pkt); err == nil {
					s.send(pkt, p)
				} else {
					fmt.Println(err)
				}

			} else {
				// fmt.Println(s.name, "packet send out #", pkt.Seq, pkt.TxTimestamp)

				// fmt.Println(pkt.RawBytes)
				// if !pkt.IsSim {
				go func() {
					outConn, err := net.Dial("udp", os.Getenv("ADDR_REMOTE_"+s.name))
					if err != nil {
						fmt.Println(err)
						return
					}
					_, err = outConn.Write(pkt.RawBytes)
					if err != nil {
						fmt.Printf("[%s] sending UDP to remote error %v\n", s.name, err)
					}
				}()
				// }
				if ANIMATION_ENABLED {
					WSLog <- Log{
						Type:  WSLOG_PKT_TX,
						PktTx: PktTx{Node: s.name, UID: pkt.UID},
					}
					go func() {
						time.Sleep(1 * time.Second)
						WSLog <- Log{
							Type:  WSLOG_PKT_TX,
							PktTx: PktTx{Node: s.name, UID: pkt.UID, Finished: true},
						}
					}()
				}
				if SAVE_STATS {
					go saveStatsDelay(s.name, subsysID2Name(pkt.Src), pkt.Seq, float64(pkt.TxTimestamp-pkt.RxTimestamp)/1000) // pkt.Delay)
				}
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
					// fmt.Printf("Pkt #%d: %d bytes, %v, delay: %.3f us\n", pkt.Seq, len(pkt.RawBytes), pkt.Path, pkt.Delay)
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
}

// create a simulated internal packet flow
func (s *Subsys) CreateFlow(dst int) {
	pkt := &Packet{
		Src:      uint8(s.id),
		Dst:      uint8(dst),
		Priority: uint8(s.Priority),
		IsSim:    true,
	}
	var buf [64]byte
	buf[0] = pkt.Src
	buf[1] = pkt.Dst
	pkt.RawBytes = buf[:]
	pkt.RxTimestamp = time.Now().UnixNano()
	pkt.Seq = getSeqNum()
	pkt.UID = getUID()
	if ANIMATION_ENABLED {
		WSLog <- Log{
			Type:  WSLOG_PKT_TX,
			PktTx: PktTx{Node: s.name, UID: pkt.UID},
		}
	}
	if p, err := s.routing(pkt); err == nil {
		s.send(pkt, p)
	} else {
		fmt.Println(err)
	}
}

// find the right port to send this packet
func (s *Subsys) routing(pkt *Packet) (*Port, error) {
L1:
	for _, entry := range s.RoutingTable[subsysID2Name(pkt.Dst)] {
		if entry.NextHop[:2] == "SW" {
			for _, swww := range Switches {
				if REROUTE_ENABLED {
					if swww.name == entry.NextHop &&
						swww.Faults[FAULT_FAILURE].Happening {
						continue L1
					}
				}
			}
		}
		for _, p := range s.portsOut {
			if p.Neighbor == entry.NextHop {
				return p, nil
			}
		}
	}

	return nil, errors.New("[" + s.name + "] cannot found next hop to " + string(pkt.Dst))
}

// send a packet out from a port
func (s *Subsys) send(pkt *Packet, port *Port) {
	// fmt.Println("sent to", port.Neighbor)
	pkt.Path = append(pkt.Path, s.name)
	port.Channel <- pkt
	s.logMutex.Lock()
	s.fwdCnt++
	s.logMutex.Unlock()
}
