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
}

var (
	CONSOLE_ENABLED = true
	boottime        int64
	LogsComm        = make(chan Log, 6553600)
	SUBSYS_LIST     = []string{"GCC", "HMS", "STR", "PWR", "ECLSS", "AGT", "INT", "EXT"} // in order
	ROUTING_TABLE   = map[int][]string{                                                  // subsysID: switches
		1: {"SW1", "SW2"},
		2: {"SW3", "SW4"},

		// 1: {"SW1", "SW7", "SW2"},
		// 2: {"SW2", "SW1", "SW3"},
		// 3: {"SW3", "SW2", "SW4"},
		// 4: {"SW4", "SW3", "SW5"},
		// 5: {"SW5", "SW4", "SW6"},
		// 6: {"SW6", "SW5", "SW7"},
		// 7: {"SW7", "SW6", "SW1"},
	}

	SequenceNumber = 0
	Subsystems     []*Subsys
	Switches       []*Switch
	Links          []*Link
	ActiveTopoTag  = ""
)

func main() {
	boottime = time.Now().Unix()

	fmt.Println(`Start Communication Network`)
	go collectStatistics()
	runHTTPSever()
}

func getSeqNum() int {
	tmp := SequenceNumber
	SequenceNumber++
	return tmp
}

func subsysID2Name(id int) string {
	return SUBSYS_LIST[id]
}

func subsysName2ID(name string) int {
	for i, n := range SUBSYS_LIST {
		if n == name {
			return i
		}
	}
	return -1
}

func findNodeByName(name string) Node {
	for _, subsys := range Subsystems {
		if subsys.Name() == name {
			return subsys
		}
	}
	for _, sw := range Switches {
		if sw.Name() == name {
			return sw
		}
	}
	// nil
	return &Switch{name: ""}
}

func collectStatistics() {
	// save to db
	go func() {
		for {
			var tmp = map[string][2]int{}
			for _, s := range Subsystems {
				tmp[s.Name()] = [2]int{s.fwdCnt, s.recvCnt}
			}

			for _, sw := range Switches {
				tmp[sw.Name()] = [2]int{sw.fwdCnt, sw.recvCnt}
			}
			saveStats(tmp)
			time.Sleep(SAVE_STATS_PERIOD * time.Second)
		}
	}()

	// send to frontend
	for {
		var tmp = map[string][2]int{}

		for _, s := range Subsystems {
			tmp[s.Name()] = [2]int{s.fwdCnt, s.recvCnt}
		}

		for _, sw := range Switches {
			tmp[sw.Name()] = [2]int{sw.fwdCnt, sw.recvCnt}
		}

		LogsComm <- Log{
			Type:       1,
			Msg:        "",
			Statistics: tmp,
		}
		time.Sleep(UPLOAD_STATS_PERIOD * time.Second)
	}
}
