package main

import (
	"encoding/json"
	"fmt"
)

var (
	presetTopos = []string{
		`{ "tag": "default", "nodes": [ { "name": "GCC", "value": [ 150, 1300 ] }, { "name": "HMS", "value": [ 500, 1150 ] }, { "name": "STR", "value": [ 650, 300 ] }, { "name": "PWR", "value": [ 1200, 75 ] }, { "name": "ECLSS", "value": [ 1750, 300 ] }, { "name": "AGT", "value": [ 1900, 850 ] }, { "name": "INT", "value": [ 1750, 1425 ] }, { "name": "EXT", "value": [ 1200, 1600 ] }, { "name": "SW0", "value": [ 1200, 850 ] }, { "name": "SW1", "value": [ 850, 1000 ] }, { "name": "SW2", "value": [ 875, 525 ] }, { "name": "SW3", "value": [ 1200, 375 ] }, { "name": "SW4", "value": [ 1525, 525 ] }, { "name": "SW5", "value": [ 1650, 850 ] }, { "name": "SW6", "value": [ 1525, 1175 ] }, { "name": "SW7", "value": [ 1200, 1300 ] } ], "edges": [ [ "HMS", "GCC" ], [ "HMS", "SW1" ], [ "HMS", "SW2" ], [ "HMS", "SW7" ], [ "STR", "SW2" ], [ "STR", "SW1" ], [ "STR", "SW3" ], [ "PWR", "SW3" ], [ "PWR", "SW2" ], [ "PWR", "SW4" ], [ "ECLSS", "SW4" ], [ "ECLSS", "SW3" ], [ "ECLSS", "SW5" ], [ "AGT", "SW5" ], [ "AGT", "SW4" ], [ "AGT", "SW6" ], [ "INT", "SW6" ], [ "INT", "SW5" ], [ "INT", "SW7" ], [ "EXT", "SW7" ], [ "EXT", "SW1" ], [ "EXT", "SW6" ], [ "SW1", "SW2" ], [ "SW2", "SW3" ], [ "SW3", "SW4" ], [ "SW4", "SW5" ], [ "SW5", "SW6" ], [ "SW6", "SW7" ], [ "SW7", "SW1" ], [ "SW1", "SW0" ], [ "SW2", "SW0" ], [ "SW3", "SW0" ], [ "SW4", "SW0" ], [ "SW5", "SW0" ], [ "SW6", "SW0" ], [ "SW7", "SW0" ] ] }`,
		`{ "tag": "noRedundancy", "nodes": [ { "name": "GCC", "value": [ 150, 1300 ] }, { "name": "HMS", "value": [ 500, 1150 ] }, { "name": "STR", "value": [ 650, 300 ] }, { "name": "PWR", "value": [ 1200, 75 ] }, { "name": "ECLSS", "value": [ 1750, 300 ] }, { "name": "AGT", "value": [ 1900, 850 ] }, { "name": "INT", "value": [ 1750, 1425 ] }, { "name": "EXT", "value": [ 1200, 1600 ] }, { "name": "SW0", "value": [ 1200, 850 ] }, { "name": "SW1", "value": [ 850, 1000 ] }, { "name": "SW2", "value": [ 875, 525 ] }, { "name": "SW3", "value": [ 1200, 375 ] }, { "name": "SW4", "value": [ 1525, 525 ] }, { "name": "SW5", "value": [ 1650, 850 ] }, { "name": "SW6", "value": [ 1525, 1175 ] }, { "name": "SW7", "value": [ 1200, 1300 ] } ], "edges": [ [ "HMS", "GCC" ], [ "HMS", "SW1" ], [ "STR", "SW2" ], [ "PWR", "SW3" ], [ "ECLSS", "SW4" ], [ "AGT", "SW5" ], [ "INT", "SW6" ], [ "EXT", "SW7" ], [ "SW1", "SW2" ], [ "SW2", "SW3" ], [ "SW3", "SW4" ], [ "SW4", "SW5" ], [ "SW5", "SW6" ], [ "SW6", "SW7" ], [ "SW7", "SW1" ], [ "SW1", "SW0" ], [ "SW2", "SW0" ], [ "SW3", "SW0" ], [ "SW4", "SW0" ], [ "SW5", "SW0" ], [ "SW6", "SW0" ], [ "SW7", "SW0" ] ] }`,
	}
)

type TopologyData struct {
	Tag   string         `json:"tag"`
	Nodes []TopologyNode `json:"nodes"`
	Edges [][2]string    `json:"edges"`
}

type TopologyNode struct {
	Name     string `json:"name"`
	Position [2]int `json:"value"`
}

func init() {
	for _, topoJSON := range presetTopos {
		var topo TopologyData
		if err = json.Unmarshal([]byte(topoJSON), &topo); err != nil {
			panic(err)
		}
		insertTopo(topo)
		if topo.Tag == "default" {
			loadTopo(topo)
		}
	}
}

// load a topo
func loadTopo(topo TopologyData) error {
	fmt.Println("load topology - " + topo.Tag)
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
			NewSwitch(n.Name, n.Position)
		} else {
			NewSubsys(n.Name, n.Position)
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
