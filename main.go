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
	GATE_NUM_SUBSYS  = 2 // 0 for in-habitat network, 1 for ground-habitat link
	GATE_NUM_SWITCH  = 8
	QUEUE_NUM_SWITCH = 8
	QUEUE_LEN_SWITCH = 4096
	BUF_LEN          = 65536
	CONFIG_LOC       = "./flex_config.json"
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
	boottime      int64
	LogsComm      = make(chan Log, 10000)
	SUBSYS_LIST   []SubsysConfig              // access by id
	SUBSYS_TABLE  = map[string]SubsysConfig{} // access by name
	ROUTING_TABLE = map[int]string{           // subsys: switch
		1: "SW1",
		2: "SW2",
		3: "SW3",
		4: "SW4",
		5: "SW5",
		6: "SW6",
		7: "SW7",
		8: "SW8",
	}
	fwdCntTotal = 0
	Subsystems  []*Subsys
	Switches    []*Switch
	Links       []*Link
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

	gcc := NewSubsys("GCC")
	hms := NewSubsys("HMS")
	agt := NewSubsys("AGT")
	str := NewSubsys("STR")
	inv := NewSubsys("INV")
	pwr := NewSubsys("PWR")
	eclss := NewSubsys("ECLSS")
	it := NewSubsys("INT")
	ext := NewSubsys("EXT")

	SW0 := NewSwitch("SW0")
	SW1 := NewSwitch("SW1")
	SW2 := NewSwitch("SW2")
	SW3 := NewSwitch("SW3")
	SW4 := NewSwitch("SW4")
	SW5 := NewSwitch("SW5")
	SW6 := NewSwitch("SW6")
	SW7 := NewSwitch("SW7")
	SW8 := NewSwitch("SW8")

	Connect(gcc, hms)
	Connect(hms, SW1)
	Connect(agt, SW2)
	Connect(str, SW3)
	Connect(inv, SW4)
	Connect(pwr, SW5)
	Connect(eclss, SW6)
	Connect(it, SW7)
	Connect(ext, SW8)

	Connect(SW1, SW0)
	Connect(SW1, SW2)

	Connect(SW2, SW0)
	Connect(SW2, SW3)

	Connect(SW3, SW0)
	Connect(SW3, SW4)

	Connect(SW4, SW0)
	Connect(SW4, SW5)

	Connect(SW5, SW0)
	Connect(SW5, SW6)

	Connect(SW6, SW0)
	Connect(SW6, SW7)

	Connect(SW7, SW0)
	Connect(SW7, SW8)

	Connect(SW8, SW0)
	Connect(SW8, SW1)

	go collectStatistics()

	go runHTTPSever()
	select {}
}

func collectStatistics() {
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
		time.Sleep(5 * time.Second)
	}
}
