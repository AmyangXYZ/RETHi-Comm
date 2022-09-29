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

type TimeWindows struct {
	Queue int `json:"queue"`
	Start int `json:"start"`
	End   int `json:"end"`
}

// Switch simulates MIMOMQ TSN switch
type Switch struct {
	name           string
	position       [2]int
	fwdCnt         int
	recvCnt        int
	GCL            [GATE_NUM_SWITCH][]TimeWindows // gateid:schedule
	gatesIn        [GATE_NUM_SWITCH]*Gate
	gatesOut       [GATE_NUM_SWITCH]*Gate
	gatesInIdx     int
	gatesOutIdx    int
	queue          [GATE_NUM_SWITCH][QUEUE_NUM_SWITCH][]*Packet // priority queue
	pktWaitlistNum [GATE_NUM_SWITCH]chan int8
	Failed         bool
	FailedDuration int

	logFwdCntMutex  sync.Mutex
	logRecvCntMutex sync.Mutex

	SeqRecoverHistory      map[int32]bool
	SeqRecoverHistoryMutex sync.Mutex
	RoutingTable           map[string][]RoutingEntry

	stopSig chan bool
}

func NewSwitch(name string, position [2]int) *Switch {
	var gatesIn [GATE_NUM_SWITCH]*Gate
	var gatesOut [GATE_NUM_SWITCH]*Gate
	var queue [GATE_NUM_SWITCH][QUEUE_NUM_SWITCH][]*Packet
	var pktWaitlistNum [GATE_NUM_SWITCH]chan int8
	var schedule [GATE_NUM_SWITCH][]TimeWindows
	// var utilization [GATE_NUM_SWITCH][]float64

	for i := 0; i < GATE_NUM_SWITCH; i++ {
		gatesIn[i] = NewGate(i, name)
		gatesOut[i] = NewGate(i, name)
		pktWaitlistNum[i] = make(chan int8, 2048)
		schedule[i] = make([]TimeWindows, HYPER_PERIOD)
		// utilization[i] = make([]float64, QUEUE_NUM_SWITCH)
		// for k := 0; k < QUEUE_NUM_SWITCH; k++ {
		// 	utilization[i][k] = math.Pow(1/3.0, float64(k)+1)
		// }

		schedule[i] = make([]TimeWindows, 0)
		for k, _ := range schedule[i] {
			schedule[i][k] = TimeWindows{
				int(rand.ExpFloat64()*9) % QUEUE_NUM_SWITCH,
				k * int(SLOT_DURATION),
				(k + 1) * int(SLOT_DURATION),
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

	go sw.Start()
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
					go sw.Classify(pkt)
					sw.logRecvCntMutex.Lock()
					sw.recvCnt++
					sw.logRecvCntMutex.Unlock()
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
				case <-sw.pktWaitlistNum[g.ID]:
					slot := ASN % HYPER_PERIOD
					window := TimeWindows{Queue: -1}
					for _, w := range sw.GCL[g.ID] {
						if w.Start <= slot && slot < w.End {
							window = w
							break
						}
					}
					if window.Queue != -1 && len(sw.queue[g.ID][window.Queue]) > 0 {
						pkt := sw.queue[g.ID][window.Queue][0]
						sw.queue[g.ID][window.Queue] = sw.queue[g.ID][window.Queue][1:]
						sw.send(pkt, g)
					} else {
						// continue loop
						sw.pktWaitlistNum[g.ID] <- 1
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
	if ANIMATION_ENABLED {
		WSLog <- Log{
			Type:  WSLOG_PKT_TX,
			PktTx: PktTx{Node: sw.name, UID: pkt.UID},
		}
	}
	if FRER_ENABLED {
		// eliminate dup
		sw.SeqRecoverHistoryMutex.Lock()
		if _, ok := sw.SeqRecoverHistory[pkt.Seq]; ok {
			// fmt.Println(sw.name, "eliminate dup from", pkt.Path[len(pkt.Path)-1])
			sw.SeqRecoverHistoryMutex.Unlock()
			return
		}
		sw.SeqRecoverHistory[pkt.Seq] = true
		sw.SeqRecoverHistoryMutex.Unlock()
		// send dup
		if gates, err := sw.routingFRER(pkt); err == nil {
			for _, g := range gates {
				dup := pkt.Dup()
				// enqueue
				sw.queue[g.ID][dup.Priority] = append(sw.queue[g.ID][pkt.Priority], dup)
				sw.pktWaitlistNum[g.ID] <- 1
				time.Sleep(100 * time.Millisecond)
			}
		}
	} else {
		if g, err := sw.routing(pkt); err == nil {
			// enqueue
			sw.queue[g.ID][pkt.Priority] = append(sw.queue[g.ID][pkt.Priority], pkt)
			sw.pktWaitlistNum[g.ID] <- 1
		} else {
			fmt.Println(err)
		}
	}
}

// find out-gate
func (sw *Switch) routing(pkt *Packet) (*Gate, error) {
L1:
	for _, entry := range sw.RoutingTable[subsysID2Name(pkt.Dst)] {
		if entry.NextHop[:2] == "SW" {
			for _, swww := range Switches {
				if swww.name == entry.NextHop && swww.Failed {
					continue L1
				}
			}
		}
		for _, g := range sw.gatesOut {
			if g.Neighbor == entry.NextHop {
				return g, nil
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
				if swww.name == entry.NextHop && swww.Failed {
					continue L1
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
