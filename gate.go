package main

// Send/Receive interface of switch/subsys
// based on Go channel
type Gate struct {
	ID       int
	Owner    string
	Neighbor string
	Channel  chan *Packet
	Failed   bool
}

func NewGate(id int, owner string) *Gate {
	g := new(Gate)
	g.ID = id
	g.Owner = owner
	g.Channel = make(chan *Packet, 1024)
	return g
}
