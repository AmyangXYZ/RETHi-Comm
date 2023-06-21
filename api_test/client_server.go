package comm

import (
	"fmt"
	"net"
	"time"
)

var (
	inConn     *net.UDPConn
	outConn    net.Conn
	txTime     time.Time
	delay      time.Duration
	totalDelay int64
)

func send(outConn net.Conn) {
	pkt := &Packet{
		Src: uint8(8),
		Dst: uint8(1),
	}
	var buf [64]byte
	buf[0] = pkt.Src
	buf[1] = pkt.Dst
	pkt.RawBytes = buf[:]

	_, err := outConn.Write(pkt.RawBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	txTime = time.Now()
	// fmt.Println("sent", n)
}

func listen(inConn *net.UDPConn) {
	cnt := 0
	for {
		var buf [128]byte
		n, err := inConn.Read(buf[0:])
		if err != nil {
			// fmt.Println(err)
			return
		}
		delay = time.Since(txTime)
		pkt := new(Packet)
		err = pkt.FromBuf(buf[0:n])
		if err != nil {
			fmt.Println(err)
			return
		}
		if cnt > 0 {
			totalDelay += int64(delay)
		}
		// fmt.Println("recv", len(pkt.RawBytes), "delay", delay)
		cnt++
	}
}

func init() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":20001")
	if err != nil {
		fmt.Println("invalid address")
		return
	}

	inConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	outConn, err = net.Dial("udp", "localhost:20001")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	go listen(inConn)
	pktNum := 1000
	time.Sleep(1 * time.Second)
	for i := 0; i < pktNum; i++ {
		send(outConn)
		time.Sleep(20 * time.Millisecond)
	}
	// select {}
	fmt.Println(float64(totalDelay) / float64(pktNum-1) / 1000)
}
