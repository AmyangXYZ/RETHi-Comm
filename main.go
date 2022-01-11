// Purpose: Run the communication network to forward packets from different subsystems
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

import (
	"fmt"
	"time"
)

const (
	GATE_NUM_SUBSYS     = 8
	GATE_NUM_SWITCH     = 8
	QUEUE_NUM_SWITCH    = 8
	QUEUE_LEN_SWITCH    = 4096
	BUF_LEN             = 65536
	SAVE_STATS_PERIOD   = 10 // in seconds
	UPLOAD_STATS_PERIOD = 3  // in seconds
	WSLOG_HEARTBEAT     = -1
	WSLOG_MSG           = 0
	WSLOG_STAT          = 1
	WSLOG_PKT_TX        = 2
)

// switch or subsys
type Node interface {
	Name() string
	OutGate() *Gate // return an idle outcoming gate for connecting
	InGate() *Gate  // return an idle incoming gate for connecting
	Start()
}

type Log struct {
	Type       int               `json:"type"` // -1: heartbeat, 0: log, 1: statistics
	Msg        string            `json:"msg"`
	Statistics map[string][2]int `json:"stats_comm"`
	PktTx      [2]string         `json:"pkt_tx"` // [src, dst]
}

var (
	MODE            = "Simulation"
	CONSOLE_ENABLED = true
	FRER_ENABLED    = true
	boottime        int64
	WSLog                 = make(chan Log, 65536)
	SUBSYS_LIST           = []string{"GCC", "HMS", "STR", "PWR", "ECLSS", "AGT", "INT", "EXT"} // in order
	SequenceNumber  int32 = 0
	Subsystems      []*Subsys
	Switches        []*Switch
	Links           []*Link
	ActiveTopoTag   = ""
)

func main() {
	fmt.Println(`Start Communication Network`)
	boottime = time.Now().Unix()
	initTopology()
	go collectStatistics()
	runHTTPSever()
}
