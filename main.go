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

// one-way delay (OWD) ~= transmission delay + propogation delay = packet_size/bit_rate + distance/speed
// wireless speed: c, wirespeed: .59c to .77c

const (
	GATE_NUM_SERVER  = 2 // 0 for in-habitat network, 1 for ground-habitat link
	GATE_NUM_SWITCH  = 8
	QUEUE_NUM_SWITCH = 8
	QUEUE_LEN_SWITCH = 4096
	BUF_LEN          = 65536
	CONFIG_LOC       = "./flex_config.json"
)

var (
	boottime int64
	LogsComm = make(chan Log, 10000)
)

type Log struct {
	Type       int               `json:"type"` // -1: heartbeat, 0: log, 1: statistics
	Msg        string            `json:"msg"`
	Statistics map[string][2]int `json:"stats_comm"`
}

type Subsys struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
}

var (
	SUBSYS_LIST   []Subsys              // access by id
	SUBSYS_TABLE  = map[string]Subsys{} // access by name
	ROUTING_TABLE = map[int]string{     // subsys: switch
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
	Servers     []*Server
	Switches    []*Switch
)

func main() {
	boottime = time.Now().Unix()

	go runHTTPSever()

	// go runDataRepo()

	select {}
}

// configure topo
func init() {
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

	// initialize subsystem
	gcc := NewServer("GCC")
	// start go-routine for this subsystem
	go gcc.Start()
	// add to server list
	Servers = append(Servers, gcc)

	hms := NewServer("HMS")
	go hms.Start()
	Servers = append(Servers, hms)

	agt := NewServer("AGT")
	go agt.Start()
	Servers = append(Servers, agt)

	str := NewServer("STR")
	go str.Start()
	Servers = append(Servers, str)

	inv := NewServer("INV")
	go inv.Start()
	Servers = append(Servers, inv)

	pwr := NewServer("PWR")
	go pwr.Start()
	Servers = append(Servers, pwr)

	eclss := NewServer("ECLSS")
	go eclss.Start()
	Servers = append(Servers, eclss)

	it := NewServer("INT")
	go it.Start()
	Servers = append(Servers, it)

	ext := NewServer("EXT")
	go ext.Start()
	Servers = append(Servers, ext)

	SW0 := NewSwitch("SW0")
	go SW0.Start()
	Switches = append(Switches, SW0)

	SW1 := NewSwitch("SW1")
	go SW1.Start()
	Switches = append(Switches, SW1)

	SW2 := NewSwitch("SW2")
	go SW2.Start()
	Switches = append(Switches, SW2)

	SW3 := NewSwitch("SW3")
	go SW3.Start()
	Switches = append(Switches, SW3)

	SW4 := NewSwitch("SW4")
	go SW4.Start()
	Switches = append(Switches, SW4)

	SW5 := NewSwitch("SW5")
	go SW5.Start()
	Switches = append(Switches, SW5)

	SW6 := NewSwitch("SW6")
	go SW6.Start()
	Switches = append(Switches, SW6)

	SW7 := NewSwitch("SW7")
	go SW7.Start()
	Switches = append(Switches, SW7)

	SW8 := NewSwitch("SW8")
	go SW8.Start()
	Switches = append(Switches, SW8)

	go collectStatistics()

	// switches
	gateIdxSW0 := 0
	go connect(SW0.GateIn, SW1.GateOut[0])
	go connect(SW1.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	go connect(SW0.GateIn, SW2.GateOut[0])
	go connect(SW2.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	go connect(SW0.GateIn, SW3.GateOut[0])
	go connect(SW3.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	go connect(SW0.GateIn, SW4.GateOut[0])
	go connect(SW4.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	go connect(SW0.GateIn, SW5.GateOut[0])
	go connect(SW5.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	go connect(SW0.GateIn, SW6.GateOut[0])
	go connect(SW6.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	go connect(SW0.GateIn, SW7.GateOut[0])
	go connect(SW7.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	go connect(SW0.GateIn, SW8.GateOut[0])
	go connect(SW8.GateIn, SW0.GateOut[gateIdxSW0])
	gateIdxSW0++

	// subsystems
	go connect(hms.GateIn, gcc.GateOut[1])
	go connect(gcc.GateIn, hms.GateOut[1])

	go connect(SW1.GateIn, hms.GateOut[0])
	go connect(hms.GateIn, SW1.GateOut[1])

	go connect(SW2.GateIn, agt.GateOut[0])
	go connect(agt.GateIn, SW2.GateOut[1])

	go connect(SW3.GateIn, str.GateOut[0])
	go connect(str.GateIn, SW3.GateOut[1])

	go connect(SW4.GateIn, inv.GateOut[0])
	go connect(inv.GateIn, SW4.GateOut[1])

	go connect(SW5.GateIn, pwr.GateOut[0])
	go connect(pwr.GateIn, SW5.GateOut[1])

	go connect(SW6.GateIn, eclss.GateOut[0])
	go connect(eclss.GateIn, SW6.GateOut[1])

	go connect(SW7.GateIn, it.GateOut[0])
	go connect(it.GateIn, SW7.GateOut[1])

	go connect(SW8.GateIn, ext.GateOut[0])
	go connect(ext.GateIn, SW8.GateOut[1])
}

func collectStatistics() {
	for {
		var tmp = map[string][2]int{}

		for _, sr := range Servers {
			tmp[sr.Name] = [2]int{sr.fwdCnt, sr.recvCnt}
		}

		for _, sw := range Switches {
			tmp[sw.Name] = [2]int{sw.fwdCnt, sw.recvCnt}
		}
		LogsComm <- Log{
			Type:       1,
			Msg:        "",
			Statistics: tmp,
		}
		time.Sleep(5 * time.Second)
	}
}
