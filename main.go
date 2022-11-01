// Purpose: Run the communication network to forward packets from different subsystems
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	PKT_BUF_LEN         = 65535
	GATE_NUM_SUBSYS     = 8
	GATE_NUM_SWITCH     = 8
	QUEUE_NUM_SWITCH    = 8
	QUEUE_LEN_SWITCH    = 1024
	SAVE_STATS_PERIOD   = 10 // in seconds
	UPLOAD_STATS_PERIOD = 3  // in seconds
	WSLOG_HEARTBEAT     = -1
	WSLOG_MSG           = 0
	WSLOG_STAT          = 1
	WSLOG_PKT_TX        = 2
	FAULT_FAILURE       = "Failure"
	FAULT_SLOW          = "Slow"
	FAULT_OVERFLOW      = "Overflow"
	FAULT_FLOODING      = "Flooding"
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
	PktTx      PktTx             `json:"pkt_tx"`
}

type PktTx struct {
	Node     string `json:"node"`
	UID      int    `json:"uid"`
	Finished bool   `json:"finished"`
}

var (
	ASN               = 0
	NEW_SLOT_SIGNAL   chan int
	HYPER_PERIOD                    = 100
	SLOT_DURATION     time.Duration = 100 // us, interval of ASN incremental
	ANIMATION_ENABLED               = false
	CONSOLE_ENABLED                 = false
	DELAY_ENABLED                   = false
	FRER_ENABLED                    = false
	REROUTE_ENABLED                 = false
	DUP_ELI_ENABLED                 = false
	JITTER_BASE                     = 0
	boottime          int64
	WSLog             = make(chan Log, 65536)
	SUBSYS_MAP        = map[string]uint8{
		"GCC":   0,
		"HMS":   1,
		"STR":   2,
		"SPL":   11,
		"ECLSS": 5,
		"PWR":   3,
		"AGT":   6,
		"IE":    8,
		"DTB":   9,
		"EXT":   7,
		"COORD": 10,
	}
	SequenceNumber int32 = 0
	UID                  = 0
	Subsystems     []*Subsys
	Switches       []*Switch
	Links          []*Link
	ActiveTopoTag  = ""
	resetASN       chan bool
)

func init() {
	// read env configs
	if os.Getenv("CONSOLE_ENABLED") == "true" {
		CONSOLE_ENABLED = true
	}
	if os.Getenv("DELAY_ENABLED") == "true" {
		DELAY_ENABLED = true
	}
	if len(os.Getenv("RAND_SEED")) > 0 {
		seed, err := strconv.Atoi(os.Getenv("RAND_SEED"))
		if err != nil {
			panic(err)
		}
		rand.Seed(int64(seed))
	} else {
		rand.Seed(time.Now().UnixNano())
	}
	if j, err := strconv.Atoi(os.Getenv("JITTER")); err == nil {
		JITTER_BASE = j
	}
	resetASN = make(chan bool)
}

func main() {
	fmt.Println(`Start Communication Network`)
	boottime = time.Now().Unix()
	initTopology()
	go collectStatistics()
	runHTTPSever()
}
