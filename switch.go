// Purpose: Simulate TSN switches, include scheduling, queueing and routing functions.
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

// Switch simulates MQMO TSN switch
type Switch struct {
	name       string
	fwdCnt     int
	recvCnt    int
	gateIn     *Gate
	gateOut    [GATE_NUM_SWITCH]*Gate
	gateOutIdx int
	queue      [QUEUE_NUM_SWITCH]chan *Packet // priority queue
}

func NewSwitch(name string) *Switch {
	var gateOut [GATE_NUM_SWITCH]*Gate
	for i := 0; i < GATE_NUM_SWITCH; i++ {
		gateOut[i] = NewGate(i, name)
	}

	var queue [QUEUE_NUM_SWITCH]chan *Packet
	for j := 0; j < QUEUE_NUM_SWITCH; j++ {
		queue[j] = make(chan *Packet, QUEUE_LEN_SWITCH)
	}

	sw := &Switch{
		name:       name,
		gateIn:     NewGate(0, name),
		gateOut:    gateOut,
		gateOutIdx: -1,
		queue:      queue,
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
	sw.gateOutIdx++
	return sw.gateOut[sw.gateOutIdx]
}

// implement Node interface
func (sw *Switch) InGate() *Gate {
	return sw.gateIn
}

// starts the switch routine
func (sw *Switch) Start() {
	// fmt.Println("Start Switch", sw.ID)

	// enqueue
	go func() {
		for {
			pkt := <-sw.gateIn.Channel
			sw.queue[pkt.Priority] <- pkt
			sw.recvCnt++
			// fmt.Println("enqueue packet to queue", pkt.Priority)
		}
	}()

	// queue with higher priority is more likely to be accessed
	// "https://medium.com/a-journey-with-go/go-ordering-in-select-statements-fd0ff80fd8d6"
	for {
		select {
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

// handle incoming packets, schedule based on its priority and forward to the destination switch/subsys
func (sw *Switch) handle(pkt *Packet) {
	sent := false
	if sw.name != "SW0" {
		for _, g := range sw.gateOut {
			if g.Neighbor == SUBSYS_LIST[pkt.Dst].Name {
				sw.send(g, pkt)
				sent = true
				break
			}
		}
		if !sent {
			// to switch 0
			for _, g := range sw.gateOut {
				if g.Neighbor == "SW0" {
					sw.send(g, pkt)
					break
				}
			}
		}
	} else {
		for _, g := range sw.gateOut {
			if g.Neighbor == ROUTING_TABLE[int(pkt.Dst)] {
				sw.send(g, pkt)
			}
		}
	}

}

// send to next hop
func (sw *Switch) send(gate *Gate, pkt *Packet) {
	pkt.Path = append(pkt.Path, sw.name)
	gate.Channel <- pkt
	sw.fwdCnt++
	// fmt.Printf("[%s] Forward packet src=%d, dst=%d to gate %s\n",
	// 	sw.name, pkt.Src, pkt.Dst, gate.Neighbor)
}
