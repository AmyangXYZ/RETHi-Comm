package main

import (
	"time"
)

// default values
var (
	SpeedWireless     float64 = 300000000   // m/s
	DistanceWireless  float64 = 57968000000 // meter 54500000000-401300000000
	BandwidthWireless float64 = 2048        // bps, 500~32000, ref: https://mars.nasa.gov/msl/mission/communications/
	SpeedWire         float64 = 231000000   // .77c
	DistanceWire      float64 = 30          // meter
	BandwidthWire     float64 = 1073741824  // 1Gbps
	PacketLossRate    float64 = 0           // percenttge
)

// Emulate ethernet cables, connect two gates, no direction
type Link struct {
	sender1 *Gate
	sink1   *Gate
	sender2 *Gate
	sink2   *Gate

	PacketLossRate float64 // percentage
	Bandwidth      float64 // in Mbps
	Speed          float64 // m/s
	Distance       float64 // in meter
	delay          float64
	// for mars-earth
	HardcodedDelay float64
	Failed         bool
	stopSig        chan bool
}

// Connect the gates of two nodes.
// Each node has a used gate counter,
// so no need to specify the gate index here
func Connect(n1, n2 Node) {
	l := new(Link)
	l.PacketLossRate = 0
	if n1.Name() == "GCC" || n2.Name() == "GCC" {
		l.Bandwidth = BandwidthWireless
		l.Speed = SpeedWireless
		l.Distance = DistanceWireless
		l.HardcodedDelay = 308
	} else {
		l.Bandwidth = BandwidthWire
		l.Speed = SpeedWire
		l.Distance = DistanceWire
	}

	l.sender1 = n1.OutGate()
	l.sink1 = n2.InGate()
	l.sender1.Neighbor = l.sink1.Owner
	l.sink1.Neighbor = l.sender1.Owner

	l.sender2 = n2.OutGate()
	l.sink2 = n1.InGate()
	l.sender2.Neighbor = l.sink2.Owner
	l.sink2.Neighbor = l.sender2.Owner

	l.stopSig = make(chan bool)

	Links = append(Links, l)

	go l.forward()
}

// size in bytes
func (l *Link) computeDelay(pktSize int) {
	// var jitter float64
	// if l.Distance < 1000 { // in-habitat
	// 	jitter = rand.ExpFloat64() * float64(JITTER_BASE) / 1000000 // us
	// }
	if l.HardcodedDelay > 0 {
		l.delay = l.HardcodedDelay
		return
	}
	l.delay = float64(pktSize)*8/l.Bandwidth + l.Distance/l.Speed
	// + jitter
	// fmt.Println("delay", l.delay)
}

func (l *Link) Stop() {
	l.stopSig <- true
}

// forward packet from one gate to the other
func (l *Link) forward() {
	for {
		select {
		case <-l.stopSig:
			return
		case pkt := <-l.sender1.Channel:
			go func() {
				if l.Failed {
					return
				}
				l.computeDelay(len(pkt.RawBytes))

				if ANIMATION_ENABLED {
					time.Sleep(1000 * time.Millisecond)
				}
				if DELAY_ENABLED {
					time.Sleep(time.Duration(l.delay) * time.Second)
				}
				pkt.Delay += l.delay
				l.sink1.Channel <- pkt
			}()
		case pkt := <-l.sender2.Channel:
			go func() {
				if l.Failed {
					return
				}
				l.computeDelay(len(pkt.RawBytes))
				// fmt.Println(l.delay)
				pkt.Delay += l.delay
				if ANIMATION_ENABLED {
					time.Sleep(1000 * time.Millisecond)
				}
				if DELAY_ENABLED {
					time.Sleep(time.Duration(l.delay) * time.Second)
				}
				l.sink2.Channel <- pkt
			}()
		}
	}
}
