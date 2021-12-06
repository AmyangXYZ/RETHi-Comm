// Purpose: Run the communication network to forward packets from different subsystems
// Date Created: 15 Apr 2021
// Date Last Modified: 17 Apr 2021
// Modeler Name: Jiachen Wang (UConn)
// Funding Acknowledgement: Funded by the NASA RETHi Project
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	GATE_NUM_SUBSYS     = 8
	GATE_NUM_SWITCH     = 8
	QUEUE_NUM_SWITCH    = 8
	QUEUE_LEN_SWITCH    = 4096
	BUF_LEN             = 65536
	CONFIG_LOC          = "./flex_config.json"
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

// for reading config from json
type SubsysConfig struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
}

var (
	CONSOLE_ENABLED = false
	boottime        int64
	LogsComm        = make(chan Log, 6553600)
	SUBSYS_LIST     []SubsysConfig              // access by id
	SUBSYS_TABLE    = map[string]SubsysConfig{} // access by name
	ROUTING_TABLE   = map[int][]string{         // subsysID: switches
		1: {"SW1", "SW7", "SW2"},
		2: {"SW2", "SW1", "SW3"},
		3: {"SW3", "SW2", "SW4"},
		4: {"SW4", "SW3", "SW5"},
		5: {"SW5", "SW4", "SW6"},
		6: {"SW6", "SW5", "SW7"},
		7: {"SW7", "SW6", "SW1"},
	}

	FwdCntTotal   = 0
	Subsystems    []*Subsys
	Switches      []*Switch
	Links         []*Link
	ActiveTopoTag = ""
)

func main() {
	boottime = time.Now().Unix()
	config, err := ioutil.ReadFile(CONFIG_LOC)
	if err != nil {
		fmt.Println("reading configuration error")
		return
	}
	err = json.Unmarshal(config, &SUBSYS_LIST)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(SUBSYS_LIST)
	for _, v := range SUBSYS_LIST {
		SUBSYS_TABLE[v.Name] = v
	}

	fmt.Println(`Start Communication Network`)

	go collectStatistics()
	runHTTPSever()
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

// load a topo
func loadTopo(topo TopologyData) error {
	fmt.Println("load topology-" + topo.Tag)
	if len(Switches) > 0 || len(Subsystems) > 0 || len(Links) > 0 {
		// fmt.Println("stop all running nodes and links")
		for _, l := range Links {
			l.Stop()
		}
		for _, s := range Subsystems {
			s.Stop()
		}
		for _, sw := range Switches {
			sw.Stop()
		}
		Links = []*Link{}
		Subsystems = []*Subsys{}
		Switches = []*Switch{}
	}

	for _, n := range topo.Nodes {
		if n.Name[:2] == "SW" {
			NewSwitch(n.Name)
		} else {
			NewSubsys(n.Name)
		}
	}
	for _, e := range topo.Edges {
		n0 := findNodeByName(e[0])
		n1 := findNodeByName(e[1])
		if n0.Name() != "" && n1.Name() != "" {
			Connect(n0, n1)
		}
	}
	return nil
}
