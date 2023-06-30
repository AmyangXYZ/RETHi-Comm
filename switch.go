package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type TimeWindow struct {
	Queue int `json:"queue"`
}

// Switch simulates MIMOMQ TSN switch
type Switch struct {
	name           string
	position       [2]int // position in the grid (frontend) (x, y)
	fwdCnt         int
	recvCnt        int
	GCL            [PORT_NUM_SWITCH][]TimeWindow // portid:schedule
	portsIn        [PORT_NUM_SWITCH]*Port
	portsOut       [PORT_NUM_SWITCH]*Port
	portsInIdx     int
	portsOutIdx    int
	Neighbors      []string
	queue          [PORT_NUM_SWITCH][QUEUE_NUM_SWITCH][]*Packet // priority queue
	queueLocker    [PORT_NUM_SWITCH][QUEUE_NUM_SWITCH]sync.Mutex
	pktWaitlistNum [PORT_NUM_SWITCH]chan int8
	Faults         map[string]Fault

	logFwdCntMutex  sync.Mutex
	logRecvCntMutex sync.Mutex

	SeqRecoverHistory      map[int32]bool
	SeqRecoverHistoryMutex sync.Mutex
	RoutingTable           map[string][]RoutingEntry

	stopSig chan bool
}

// Fault injected to the switch
type Fault struct {
	Type      string
	Happening bool
	Durtaion  int
}

// New switch
func NewSwitch(name string, position [2]int) *Switch {
	var portsIn [PORT_NUM_SWITCH]*Port
	var portsOut [PORT_NUM_SWITCH]*Port
	var queue [PORT_NUM_SWITCH][QUEUE_NUM_SWITCH][]*Packet
	var pktWaitlistNum [PORT_NUM_SWITCH]chan int8
	var schedule [PORT_NUM_SWITCH][]TimeWindow

	for i := 0; i < PORT_NUM_SWITCH; i++ {
		portsIn[i] = NewPort(i, name)
		portsOut[i] = NewPort(i, name)
		pktWaitlistNum[i] = make(chan int8, 2048)
		schedule[i] = make([]TimeWindow, HYPER_PERIOD)
		// factors := []int{}
		// for f := 2; f < HYPER_PERIOD; f++ {
		// 	if HYPER_PERIOD%f == 0 {
		// 		factors = append(factors, f)
		// 	}
		// }
		// interval := factors[rand.Intn(len(factors))]
		offset := rand.Int()
		for k := 0; k < len(schedule[i]); k++ {
			q := rand.Int() % 8
			if (k+offset)%3 == 0 {
				q = 0
			} else if (k+offset)%7 == 0 {
				q = 1
			} else if (k+offset)%8 == 0 {
				q = 2
			} else if (k+offset)%11 == 0 {
				q = 3
			} else if (k+offset)%20 == 0 {
				q = 4
			} else if (k+offset)%23 == 0 {
				q = 5
			} else if (k+offset)%39 == 0 {
				q = 6
			} else if (k+offset)%47 == 0 {
				q = 7
			}
			schedule[i][k] = TimeWindow{
				Queue: q,
			}
		}
	}

	sw := &Switch{
		name:              name,
		position:          position,
		portsIn:           portsIn,
		portsOut:          portsOut,
		portsInIdx:        -1,
		portsOutIdx:       -1,
		queue:             queue,
		pktWaitlistNum:    pktWaitlistNum,
		GCL:               schedule,
		RoutingTable:      make(map[string][]RoutingEntry),
		SeqRecoverHistory: make(map[int32]bool),
		stopSig:           make(chan bool),
		Faults:            make(map[string]Fault),
	}

	for _, subsys := range Subsystems {
		paths := Graph.FindAllPaths(name, subsys.Name())
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
		sw.RoutingTable[subsys.Name()] = table
	}
	// fmt.Println(name)
	// for dst, p := range sw.RoutingTable {
	// 	fmt.Println("    ", dst, p)
	// }

	Switches = append(Switches, sw)
	return sw
}

// implement Node interface
func (sw *Switch) Name() string {
	return sw.name
}

// implement Node interface
func (sw *Switch) OutPort() *Port {
	sw.portsOutIdx++
	return sw.portsOut[sw.portsOutIdx]
}

// implement Node interface
func (sw *Switch) InPort() *Port {
	sw.portsInIdx++
	return sw.portsIn[sw.portsInIdx]
}

// starts the switch routine
func (sw *Switch) Start() {
	// fmt.Println("Start Switch", sw.ID)

	// handle incoming packets
	for _, inPort := range sw.portsIn {
		go func(p *Port) {
			for {
				select {
				case <-sw.stopSig:
					return
				case pkt := <-p.Channel:
					// sw.queue[pkt.Priority] <- pkt
					sw.logRecvCntMutex.Lock()
					sw.recvCnt++
					sw.logRecvCntMutex.Unlock()
					if ANIMATION_ENABLED {
						WSLog <- Log{
							Type:  WSLOG_PKT_TX,
							PktTx: PktTx{Node: sw.name, UID: pkt.UID},
						}
					}
					if !sw.Faults[FAULT_FAILURE].Happening {
						if sw.Faults[FAULT_SLOW].Happening {
							time.Sleep(2 * time.Second)
						}
						if sw.Faults[FAULT_OVERFLOW].Happening {
							if sw.recvCnt%5 == 1 {
								continue
							}
						}
						go sw.Classify(pkt)
					}
				}
			}
		}(inPort)
	}

	// watch queues of each port and send out according to the schedule
	for _, outPort := range sw.portsOut {
		go func(p *Port) {
			for {
				select {
				case <-sw.stopSig:
					return
				case asn := <-NEW_SLOT_SIGNAL:
					window := sw.GCL[p.ID][asn%HYPER_PERIOD]

					if len(sw.queue[p.ID][window.Queue]) > 0 {
						sw.queueLocker[p.ID][window.Queue].Lock()
						pkt := sw.queue[p.ID][window.Queue][0]
						sw.queue[p.ID][window.Queue] = sw.queue[p.ID][window.Queue][1:]
						sw.queueLocker[p.ID][window.Queue].Unlock()
						sw.send(pkt, p)
						if sw.Faults[FAULT_FLOODING].Happening {
							go func() {
								for i := 0; i < 10; i++ {
									dup := pkt.Dup()
									sw.send(dup, p)
									time.Sleep(200 * time.Millisecond)
								}
							}()
						}
					}
				}
			}
		}(outPort)
	}
}

// Stop stops the switch
func (sw *Switch) Stop() {
	for range sw.portsIn {
		sw.stopSig <- true // stop inports
	}
	for range sw.portsOut {
		sw.stopSig <- true // stop outports
	}
	// fmt.Println(sw.name, "stopped")
}

// classify packet belongs to which out-port
func (sw *Switch) Classify(pkt *Packet) {
	if len(pkt.Path) > 20 {
		return
	}

	if FRER_ENABLED || DUP_ELI_ENABLED {
		// eliminate dup
		sw.SeqRecoverHistoryMutex.Lock()
		if _, ok := sw.SeqRecoverHistory[pkt.Seq]; ok {
			// fmt.Println(sw.name, "eliminate dup from", pkt.Path[len(pkt.Path)-1])
			sw.SeqRecoverHistoryMutex.Unlock()
			return
		}
		sw.SeqRecoverHistory[pkt.Seq] = true
		sw.SeqRecoverHistoryMutex.Unlock()
	}
	if FRER_ENABLED {
		// send dup
		if ports, err := sw.routingFRER(pkt); err == nil {
			for _, p := range ports {
				dup := pkt.Dup()
				if TAS_ENABLED {
					// enqueue
					sw.queueLocker[p.ID][pkt.Priority].Lock()
					sw.queue[p.ID][dup.Priority] = append(sw.queue[p.ID][pkt.Priority], dup)
					// sw.pktWaitlistNum[p.ID] <- 1
					sw.queueLocker[p.ID][pkt.Priority].Unlock()
					time.Sleep(100 * time.Millisecond)
				} else {
					sw.send(pkt, p)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	} else {
		if p, err := sw.routing(pkt); err == nil {
			if TAS_ENABLED {
				// enqueue
				sw.queueLocker[p.ID][pkt.Priority].Lock()
				sw.queue[p.ID][pkt.Priority] = append(sw.queue[p.ID][pkt.Priority], pkt)
				sw.queueLocker[p.ID][pkt.Priority].Unlock()
				// sw.pktWaitlistNum[p.ID] <- 1
			} else {
				sw.send(pkt, p)
			}
		} else {
			fmt.Println(err)
		}
	}
}

// find out-port
func (sw *Switch) routing(pkt *Packet) (*Port, error) {
	if sw.Faults[FAULT_MIS_ROUTING].Happening {
	L:
		for i := len(sw.RoutingTable[subsysID2Name(pkt.Dst)]) - 1; i >= 0; i-- {
			entry := sw.RoutingTable[subsysID2Name(pkt.Dst)][i]
			if entry.NextHop[:2] == "SW" {
				for _, swww := range Switches {
					if REROUTE_ENABLED {
						if swww.name == entry.NextHop && swww.Faults[FAULT_FAILURE].Happening {
							continue L
						}
					}
				}
			}
			for _, p := range sw.portsOut {
				if p.Neighbor == entry.NextHop {
					return p, nil
				}
			}
		}

	} else {
	L1:
		for _, entry := range sw.RoutingTable[subsysID2Name(pkt.Dst)] {
			if entry.NextHop[:2] == "SW" {
				for _, swww := range Switches {
					if REROUTE_ENABLED {
						if swww.name == entry.NextHop && swww.Faults[FAULT_FAILURE].Happening {
							continue L1
						}
					}
				}
			}
			for _, p := range sw.portsOut {
				if p.Neighbor == entry.NextHop {
					return p, nil
				}
			}
		}
	}

	return nil, errors.New("[" + sw.name + "] cannot found next hop")
}

// routing when FRER enabled
func (sw *Switch) routingFRER(pkt *Packet) ([]*Port, error) {
	ports := []*Port{}
L1:
	for _, entry := range sw.RoutingTable[subsysID2Name(pkt.Dst)] {
		if entry.NextHop[:2] == "SW" {
			for _, swww := range Switches {
				if REROUTE_ENABLED {
					if swww.name == entry.NextHop && swww.Faults[FAULT_FAILURE].Happening {
						continue L1
					}
				}
			}
		}
		if entry.NextHop == pkt.Path[len(pkt.Path)-1] {
			continue
		}
		for _, p := range sw.portsOut {
			if p.Neighbor == entry.NextHop {
				ports = append(ports, p)
			}
		}
	}
	if len(ports) == 0 {
		return nil, errors.New("[" + sw.name + "] cannot found next hop")
	}
	return ports, nil
}

// Send packet to port
func (sw *Switch) send(pkt *Packet, port *Port) {
	// fmt.Println("sent to", port.Neighbor)
	pkt.Path = append(pkt.Path, sw.name)
	port.Channel <- pkt
	if ANIMATION_ENABLED {
		WSLog <- Log{
			Type:  WSLOG_PKT_TX,
			PktTx: PktTx{Node: sw.name, UID: pkt.UID},
		}
	}
	sw.logFwdCntMutex.Lock()
	sw.fwdCnt++
	sw.logFwdCntMutex.Unlock()
}
