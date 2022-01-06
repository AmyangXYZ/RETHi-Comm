package main

import "time"

// default values
var (
	SpeedWireless     float64 = 300000000    // m/s
	DistanceWireless  float64 = 331940000000 // meter
	BandwidthWireless float64 = 2048         // bps, 500~32000, ref: https://mars.nasa.gov/msl/mission/communications/
	SpeedWire         float64 = 231000000    // .77c
	DistanceWire      float64 = 30           // meter
	BandwidthWire     float64 = 1073741824   // 1Gbps
	PacketLossRate    float64 = 0            // percenttge
)

// no direction
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

	Failed  bool
	stopSig chan bool
}

func Connect(n1, n2 Node) {
	l := new(Link)
	l.PacketLossRate = 0
	if n1.Name() == "GCC" || n2.Name() == "GCC" {
		l.Bandwidth = BandwidthWireless
		l.Speed = SpeedWireless
		l.Distance = DistanceWireless
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
	l.delay = float64(pktSize)*8/l.Bandwidth + l.Distance/l.Speed
	// fmt.Println("delay", l.delay)
}

func (l *Link) Stop() {
	l.stopSig <- true
}

func (l *Link) forward() {
	for {
		select {
		case <-l.stopSig:
			return
		case pkt := <-l.sender1.Channel:
			l.computeDelay(len(pkt.RawBytes))
			// fmt.Println(l.delay)
			time.Sleep(time.Duration(l.delay) * time.Second)
			pkt.Delay += l.delay
			l.sink1.Channel <- pkt
		case pkt := <-l.sender2.Channel:
			l.computeDelay(len(pkt.RawBytes))
			// fmt.Println(l.delay)
			pkt.Delay += l.delay
			// time.Sleep(time.Duration(l.delay) * time.Second)
			l.sink2.Channel <- pkt
		}
	}
}
