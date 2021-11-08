// Purpose: Define the gate of TSN switches, connecting switches and servers.
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

import (
	"fmt"
	"time"
)

var GROUND_LINK_DELAY = 3 // in second

type Gate struct {
	ID        int
	Owner     string
	Neighbor  string
	Bandwidth float64
	Delay     float64
	Channel   chan *Packet
}

func NewGate(id int, owner string) *Gate {
	return &Gate{
		ID:      id,
		Owner:   owner,
		Channel: make(chan *Packet, 1),
	}
}

// connect two gates
func connect(dst, src *Gate) {
	dst.Neighbor = src.Owner
	src.Neighbor = dst.Owner
	for {
		pkt := <-src.Channel
		fmt.Println("latency", GROUND_LINK_DELAY)
		if dst.Owner == "GCC" || src.Owner == "GCC" {
			time.Sleep(3 * time.Second)
		}
		dst.Channel <- pkt
	}
}
