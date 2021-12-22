// Purpose: Simulate TSN switches, include scheduling, queueing and routing functions.
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

import (
	"errors"
	"fmt"
)

// Switch simulates MQMO TSN switch
type Switch struct {
	name           string
	position       [2]int
	fwdCnt         int
	recvCnt        int
	gatesIn        [GATE_NUM_SWITCH]*Gate
	gatesOut       [GATE_NUM_SWITCH]*Gate
	gatesInIdx     int
	gatesOutIdx    int
	queue          [QUEUE_NUM_SWITCH]chan *Packet // priority queue
	Failed         bool
	FailedDuration int
	stopSig        chan bool
}

func NewSwitch(name string, position [2]int) *Switch {
	var gatesIn [GATE_NUM_SWITCH]*Gate
	var gatesOut [GATE_NUM_SWITCH]*Gate
	for i := 0; i < GATE_NUM_SWITCH; i++ {
		gatesIn[i] = NewGate(i, name)
		gatesOut[i] = NewGate(i, name)
	}

	var queue [QUEUE_NUM_SWITCH]chan *Packet
	for j := 0; j < QUEUE_NUM_SWITCH; j++ {
		queue[j] = make(chan *Packet, QUEUE_LEN_SWITCH)
	}

	sw := &Switch{
		name:        name,
		position:    position,
		gatesIn:     gatesIn,
		gatesOut:    gatesOut,
		gatesInIdx:  -1,
		gatesOutIdx: -1,
		queue:       queue,
		stopSig:     make(chan bool),
	}
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

	// enqueue
	for _, inGate := range sw.gatesIn {
		go func(g *Gate) {
			for {
				select {
				case <-sw.stopSig:
					return
				case pkt := <-g.Channel:
					sw.queue[pkt.Priority] <- pkt
					sw.recvCnt++
					// fmt.Println(sw.name, "enqueue packet to queue", pkt.Priority)
				}
			}
		}(inGate)
	}

	// queue with higher priority is more likely to be accessed
	// "https://medium.com/a-journey-with-go/go-ordering-in-select-statements-fd0ff80fd8d6"
	for {
		select {
		case <-sw.stopSig:
			return
		case pkt := <-sw.queue[0]:
			sw.handle(pkt)

		case pkt := <-sw.queue[1]:
			sw.handle(pkt)
		case pkt := <-sw.queue[1]:
			sw.handle(pkt)

		case pkt := <-sw.queue[2]:
			sw.handle(pkt)
		case pkt := <-sw.queue[2]:
			sw.handle(pkt)
		case pkt := <-sw.queue[2]:
			sw.handle(pkt)

		case pkt := <-sw.queue[3]:
			sw.handle(pkt)
		case pkt := <-sw.queue[3]:
			sw.handle(pkt)
		case pkt := <-sw.queue[3]:
			sw.handle(pkt)
		case pkt := <-sw.queue[3]:
			sw.handle(pkt)

		case pkt := <-sw.queue[4]:
			sw.handle(pkt)
		case pkt := <-sw.queue[4]:
			sw.handle(pkt)
		case pkt := <-sw.queue[4]:
			sw.handle(pkt)
		case pkt := <-sw.queue[4]:
			sw.handle(pkt)
		case pkt := <-sw.queue[4]:
			sw.handle(pkt)

		case pkt := <-sw.queue[5]:
			sw.handle(pkt)
		case pkt := <-sw.queue[5]:
			sw.handle(pkt)
		case pkt := <-sw.queue[5]:
			sw.handle(pkt)
		case pkt := <-sw.queue[5]:
			sw.handle(pkt)
		case pkt := <-sw.queue[5]:
			sw.handle(pkt)
		case pkt := <-sw.queue[5]:
			sw.handle(pkt)

		case pkt := <-sw.queue[6]:
			sw.handle(pkt)
		case pkt := <-sw.queue[6]:
			sw.handle(pkt)
		case pkt := <-sw.queue[6]:
			sw.handle(pkt)
		case pkt := <-sw.queue[6]:
			sw.handle(pkt)
		case pkt := <-sw.queue[6]:
			sw.handle(pkt)
		case pkt := <-sw.queue[6]:
			sw.handle(pkt)
		case pkt := <-sw.queue[6]:
			sw.handle(pkt)

		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		case pkt := <-sw.queue[7]:
			sw.handle(pkt)
		}
	}
}

func (sw *Switch) Stop() {
	sw.stopSig <- true // stop queue handler
	for range sw.gatesIn {
		sw.stopSig <- true // stop ingates
	}
	// fmt.Println(sw.name, "stopped")
}

// handle incoming packets, schedule based on its priority and forward to the destination switch/subsys
func (sw *Switch) handle(pkt *Packet) {
	if g, err := sw.routing(pkt); err == nil {
		// fmt.Println("sent to", g.Neighbor)
		pkt.Path = append(pkt.Path, sw.name)
		g.Channel <- pkt
		sw.fwdCnt++
	} else {
		fmt.Println(err)
	}
}

func (sw *Switch) routing(pkt *Packet) (*Gate, error) {
	for _, g := range sw.gatesOut {
		if g.Neighbor == SUBSYS_LIST[pkt.Dst] {
			return g, nil
		}
	}

	for _, g := range sw.gatesOut {
		for _, otherSw := range ROUTING_TABLE[int(pkt.Dst)] {
			if g.Neighbor == otherSw {
				for _, swww := range Switches {
					if swww.name == otherSw {
						if !swww.Failed {
							return g, nil
						}
					}
				}
			}
		}
	}

	// not found, send to SW0
	for _, g := range sw.gatesOut {
		if g.Neighbor == "SW0" {
			return g, nil
		}
	}

	return nil, errors.New("[" + sw.name + "] cannot found next hop")
}
