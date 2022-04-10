package main

import (
	"encoding/binary"
	"math/rand"
)

const (
	PACKET_TYPE_CTRL = 0x00
	PACKET_TYPE_DATA = 0x01
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
	// protocol use
	Src          uint8  `json:"src"`
	Dst          uint8  `json:"dst"`
	MessageType  uint8  `json:"message_type"`
	Priority     uint8  `json:"priority"`
	Version      uint8  `json:"version"`
	Reserved     uint8  `json:"reserved"`
	PhysicalTime uint32 `json:"physical_time"`
	SimulinkTime uint32 `json:"simulink_time"`
	Sequence     uint16 `json:"sequence"`
	Length       uint16 `json:"length"`
	Payload      []byte `json:"-"`

	// internal use
	IsSim       bool
	RawBytes    []byte
	Delay       float64
	Path        []string
	Seq         int32 // for 802.1CB-FRER
	RxTimestamp int64
	TxTimestamp int64
	DupID       int
}

func (pkt *Packet) FromBuf(buf []byte) error {
	pkt.Src = uint8(buf[0])
	pkt.Dst = uint8(buf[1])
	temp := binary.LittleEndian.Uint16(buf[2:4])
	pkt.MessageType = uint8(temp >> 12)
	pkt.Priority = uint8(temp >> 8 & 0x0f)
	pkt.Version = uint8(temp >> 4 & 0x0f)
	pkt.Reserved = uint8(temp & 0x0f)
	pkt.PhysicalTime = uint32(binary.LittleEndian.Uint32(buf[4:8]))
	pkt.SimulinkTime = uint32(binary.LittleEndian.Uint32(buf[8:12]))
	pkt.Sequence = binary.LittleEndian.Uint16(buf[12:14])
	pkt.Length = binary.LittleEndian.Uint16(buf[14:16])
	pkt.Payload = buf[16:]
	pkt.RawBytes = buf
	return nil
}

func (pkt *Packet) ToBuf() []byte {
	var buf [16]byte
	buf[0] = byte(pkt.Src)
	buf[1] = byte(pkt.Dst)
	temp := uint16(pkt.MessageType)<<12 + uint16(pkt.Priority)<<8 + uint16(pkt.Version)<<4 + uint16(pkt.Reserved)
	binary.LittleEndian.PutUint16(buf[2:4], uint16(temp))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(pkt.PhysicalTime))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(pkt.SimulinkTime))
	binary.LittleEndian.PutUint16(buf[12:14], uint16(pkt.Sequence))
	binary.LittleEndian.PutUint16(buf[14:16], uint16(pkt.Length))
	return append(buf[:], pkt.Payload...)
}

// deep copy/duplicate
func (pkt *Packet) Dup() *Packet {
	dup := new(Packet)
	*dup = *pkt
	// copy slice by value
	dup.Path = append([]string{}, pkt.Path...)
	dup.DupID = rand.Intn(1024)
	return dup
}
