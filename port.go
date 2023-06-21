package comm

// Send/Receive interface of switch/subsys
// based on Go channel
type Port struct {
	ID       int
	Owner    string
	Neighbor string
	Channel  chan *Packet
	Failed   bool
}

func NewPort(id int, owner string) *Port {
	p := new(Port)
	p.ID = id
	p.Owner = owner
	p.Channel = make(chan *Packet, 1024)
	return p
}
