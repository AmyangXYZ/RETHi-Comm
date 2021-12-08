package main

import (
	"encoding/binary"
	"time"
)

type Packet struct {
	// protocol use
	Src          uint8  `json:"src"`
	Dst          uint8  `json:"dst"`
	MessageType  uint8  `json:"msg_type"`
	DataType     uint8  `json:"data_type"`
	Priority     uint8  `json:"priority"`
	PhysicalTime uint32 `json:"phy_time"`
	SimulinkTime uint32 `json:"sim_time"`
	Row          uint8  `json:"row"`
	Col          uint8  `json:"col"`
	Length       uint16 `json:"len"`
	Payload      []byte `json:"payload"`

	// internal use
	IsSim        bool
	RawBytes     []byte
	Delay        float64
	Path         []string
	TimeCreated  time.Time
	TimeReceived time.Time
}

func (pkt *Packet) FromBuf(buf []byte) error {
	pkt.Src = uint8(buf[0])
	pkt.Dst = uint8(buf[1])
	// temp := binary.LittleEndian.Uint16(buf[2:4])
	// pkt.MessageType = uint8(temp >> 12 & 0x0f)
	// pkt.DataType = uint8(temp >> 4 & 0xff)
	// pkt.Priority = uint8(temp & 0x0f)
	// pkt.PhysicalTime = uint32(binary.LittleEndian.Uint16(buf[4:8]))
	// pkt.SimulinkTime = uint32(binary.LittleEndian.Uint16(buf[8:12]))
	// pkt.Row = uint8(buf[12])
	// pkt.Col = uint8(buf[13])
	// pkt.Length = binary.LittleEndian.Uint16(buf[14:16])
	pkt.RawBytes = buf
	return nil
}

func (pkt *Packet) ToBuf() []byte {
	var buf [16]byte
	buf[0] = byte(pkt.Src)
	buf[1] = byte(pkt.Dst)
	temp := uint16(pkt.MessageType)<<12 + uint16(pkt.DataType)<<4 + uint16(pkt.Priority)
	binary.LittleEndian.PutUint16(buf[2:4], uint16(temp))
	binary.LittleEndian.PutUint16(buf[4:8], uint16(pkt.PhysicalTime))
	binary.LittleEndian.PutUint16(buf[8:12], uint16(pkt.SimulinkTime))
	buf[12] = byte(pkt.Row)
	buf[13] = byte(pkt.Col)
	binary.LittleEndian.PutUint16(buf[14:16], uint16(pkt.Length))
	return append(buf[:], pkt.Payload...)
}
