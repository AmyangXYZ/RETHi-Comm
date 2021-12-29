package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type RoutingEntry struct {
	Dst     int
	NextHop string
	Cost    int
}

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
	presetTopos := []TopologyData{}
	j, err := ioutil.ReadFile("./topos.json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(j, &presetTopos); err != nil {
		panic(err)
	}

	for _, topo := range presetTopos {
		insertTopo(topo)
		if topo.Tag == "frer" {
			loadTopo(topo)
		}
	}
}

// load a topo
func loadTopo(topo TopologyData) error {
	fmt.Println("load topology - " + topo.Tag)
	genRoutingTable(topo)

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

func genRoutingTable(topo TopologyData) {
	// fmt.Println(topo.Edges)
}
