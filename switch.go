// Purpose: Simulate TSN switches, include scheduling, queueing and routing functions.
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
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
	position       [2]int
	fwdCnt         int
	recvCnt        int
	GCL            [GATE_NUM_SWITCH][]TimeWindow // gateid:schedule
	gatesIn        [GATE_NUM_SWITCH]*Gate
	gatesOut       [GATE_NUM_SWITCH]*Gate
	gatesInIdx     int
	gatesOutIdx    int
	Neighbors      []string
	queue          [GATE_NUM_SWITCH][QUEUE_NUM_SWITCH][]*Packet // priority queue
	queueLocker    [GATE_NUM_SWITCH][QUEUE_NUM_SWITCH]sync.Mutex
	pktWaitlistNum [GATE_NUM_SWITCH]chan int8
	Faults         map[string]Fault

	logFwdCntMutex  sync.Mutex
	logRecvCntMutex sync.Mutex

	SeqRecoverHistory      map[int32]bool
	SeqRecoverHistoryMutex sync.Mutex
	RoutingTable           map[string][]RoutingEntry

	stopSig chan bool
}

type Fault struct {
	Type      string
	Happening bool
	Durtaion  int
}

func NewSwitch(name string, position [2]int) *Switch {
	var gatesIn [GATE_NUM_SWITCH]*Gate
	var gatesOut [GATE_NUM_SWITCH]*Gate
	var queue [GATE_NUM_SWITCH][QUEUE_NUM_SWITCH][]*Packet
	var pktWaitlistNum [GATE_NUM_SWITCH]chan int8
	var schedule [GATE_NUM_SWITCH][]TimeWindow

	for i := 0; i < GATE_NUM_SWITCH; i++ {
		gatesIn[i] = NewGate(i, name)
		gatesOut[i] = NewGate(i, name)
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
		gatesIn:           gatesIn,
		gatesOut:          gatesOut,
		gatesInIdx:        -1,
		gatesOutIdx:       -1,
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
func (sw *Switch) OutGate() *Gate {
	sw.gatesOutIdx++
	return sw.gatesOut[sw.gatesOutIdx]
}

// implement Node interface
func (sw *Switch) InGate() *Gate {
	sw.gatesInIdx++
	return sw.gatesIn[sw.gatesInIdx]
}

// starts the switch routine
func (sw *Switch) Start() {
	// fmt.Println("Start Switch", sw.ID)

	// handle incoming packets
	for _, inGate := range sw.gatesIn {
		go func(g *Gate) {
			for {
				select {
				case <-sw.stopSig:
					return
				case pkt := <-g.Channel:
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
		}(inGate)
	}

	// watch queues of each gate and send out according to the schedule
	for _, outGate := range sw.gatesOut {
		go func(g *Gate) {
			for {
				select {
				case <-sw.stopSig:
					return
				case asn := <-NEW_SLOT_SIGNAL:
					window := sw.GCL[g.ID][asn%HYPER_PERIOD]

					if len(sw.queue[g.ID][window.Queue]) > 0 {
						sw.queueLocker[g.ID][window.Queue].Lock()
						pkt := sw.queue[g.ID][window.Queue][0]
						sw.queue[g.ID][window.Queue] = sw.queue[g.ID][window.Queue][1:]
						sw.queueLocker[g.ID][window.Queue].Unlock()
						sw.send(pkt, g)
						if sw.Faults[FAULT_FLOODING].Happening {
							go func() {
								for i := 0; i < 10; i++ {
									dup := pkt.Dup()
									sw.send(dup, g)
									time.Sleep(200 * time.Millisecond)
								}
							}()
						}
					}
				}

			}
		}(outGate)
	}

}

func (sw *Switch) Stop() {
	for range sw.gatesIn {
		sw.stopSig <- true // stop ingates
	}
	for range sw.gatesOut {
		sw.stopSig <- true // stop outgates
	}
	// fmt.Println(sw.name, "stopped")
}

// classify packet belongs to which out-gate
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
		if gates, err := sw.routingFRER(pkt); err == nil {
			for _, g := range gates {
				dup := pkt.Dup()
				// enqueue
				sw.queueLocker[g.ID][pkt.Priority].Lock()
				sw.queue[g.ID][dup.Priority] = append(sw.queue[g.ID][pkt.Priority], dup)
				// sw.pktWaitlistNum[g.ID] <- 1
				sw.queueLocker[g.ID][pkt.Priority].Unlock()
				time.Sleep(100 * time.Millisecond)
			}
		}
	} else {
		if g, err := sw.routing(pkt); err == nil {
			// enqueue
			sw.queueLocker[g.ID][pkt.Priority].Lock()
			sw.queue[g.ID][pkt.Priority] = append(sw.queue[g.ID][pkt.Priority], pkt)
			sw.queueLocker[g.ID][pkt.Priority].Unlock()
			// sw.pktWaitlistNum[g.ID] <- 1
		} else {
			fmt.Println(err)
		}
	}
}

// find out-gate
func (sw *Switch) routing(pkt *Packet) (*Gate, error) {
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
			for _, g := range sw.gatesOut {
				if g.Neighbor == entry.NextHop {
					return g, nil
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
			for _, g := range sw.gatesOut {
				if g.Neighbor == entry.NextHop {
					return g, nil
				}
			}
		}
	}

	return nil, errors.New("[" + sw.name + "] cannot found next hop")
}

func (sw *Switch) routingFRER(pkt *Packet) ([]*Gate, error) {
	gates := []*Gate{}
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
		for _, g := range sw.gatesOut {
			if g.Neighbor == entry.NextHop {
				gates = append(gates, g)
			}
		}
	}
	if len(gates) == 0 {
		return nil, errors.New("[" + sw.name + "] cannot found next hop")
	}
	return gates, nil
}

func (sw *Switch) send(pkt *Packet, gate *Gate) {
	// fmt.Println("sent to", gate.Neighbor)
	pkt.Path = append(pkt.Path, sw.name)
	gate.Channel <- pkt
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
