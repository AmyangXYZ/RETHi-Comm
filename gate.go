// Purpose: Define the gate of TSN switches, connecting switches and servers.
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

type Gate struct {
	ID       int
	Owner    string
	Neighbor string
	Channel  chan *Packet
}

func NewGate(id int, owner string) *Gate {
	g := new(Gate)
	g.ID = id
	g.Owner = owner
	g.Channel = make(chan *Packet, 1024)
	return g
}
