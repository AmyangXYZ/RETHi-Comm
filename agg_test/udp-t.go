package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

//                     1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      SRC      |       DST     | TYPE  | PRIO  |  VER  |  RES  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                     PHYSICAL_TIMESTAMP                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                     SIMULINK_TIMESTAMP                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |            SEQUENCE           |              LEN              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      DATA…
// +-+-+-+-+-+-+-+-+

type Packet struct {
	Src uint8
	Dst uint8
	// and other headers
	SubFramesNum uint8
}

func (pkt *Packet) FromBuf(buf []byte) error {
	pkt.Src = uint8(buf[0])
	pkt.Dst = uint8(buf[1])
	// other headers

	// start reading subframes
	pkt.SubFramesNum = uint8(buf[2])
	bufCursor := uint8(3)
	for i := 0; i < int(pkt.SubFramesNum); i++ {
		row := uint8(buf[bufCursor])
		col := uint8(buf[bufCursor+1])
		length := uint8(buf[bufCursor+2])
		bufCursor += 3

		payload := buf[bufCursor : bufCursor+length*8]
		bufCursor += length * 8
		data64 := []float64{}
		for j := 0; j < int(length); j++ {
			data64 = append(data64, math.Float64frombits(binary.LittleEndian.Uint64(payload[j*8:j*8+8])))
		}

		fmt.Println("subframe", i, ":", row, col, length, data64)
	}
	return nil

}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":10001")
	if err != nil {
		fmt.Println("invalid address")
		return
	}

	inConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		var buf [1024]byte
		n, _ := inConn.Read(buf[0:])
		fmt.Println("buffer", buf[0:n])

		pkt := new(Packet)
		if err = pkt.FromBuf(buf[0:n]); err != nil {
			fmt.Println(err)
		}
	}
}
