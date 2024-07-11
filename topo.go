package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

var (
	Graph *TopoGraph
)

type RoutingEntry struct {
	NextHop  string
	HopCount int
}

// for en/decode
type TopologyData struct {
	Tag   string         `json:"tag"`
	Nodes []TopologyNode `json:"nodes"`
	Edges [][2]string    `json:"edges"`
}

// for en/decode
type TopologyNode struct {
	Name     string `json:"name"`
	Position [2]int `json:"value"`
}

type LinkDelay struct {
	Delay map[string]map[string]int `json:"delays"`
}

func initTopology() {
	presetTopos := []TopologyData{}
	j, err := os.ReadFile("./topos.json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(j, &presetTopos); err != nil {
		panic(err)
	}

	for _, topo := range presetTopos {
		insertTopo(topo)
		if topo.Tag == "default" {
			ActiveTopoTag = topo.Tag
			loadTopo(topo)
		}
	}

	if DELAY_ENABLED {
		j, err = os.ReadFile("/delays.json")
		if err != nil {
			fmt.Println("No preset link delays")
			return
		}
		presetDelays := LinkDelay{}
		if err = json.Unmarshal(j, &presetDelays); err != nil {
			panic(err)
		}

		for srcName, dsts := range presetDelays.Delay {
			for dstName, delay := range dsts {
				l := findLinkByNodes(findNodeByName(srcName), findNodeByName(dstName))
				if l != nil {
					l.HardcodedDelay = float64(delay)
				}
			}
		}

		HARDCODE_DELAY_ENABLED = true
		fmt.Println("Preset link delays enabled")
	}
}

// load a topo
func loadTopo(topo TopologyData) error {
	fmt.Println("load topology - " + topo.Tag)
	Graph = new(TopoGraph)
	Graph.construct(topo)
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
		resetASN <- true
		close(NEW_SLOT_SIGNAL)
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

	NEW_SLOT_SIGNAL = make(chan int, len(Switches)*PORT_NUM_SWITCH)
	go func() {
		// time.Sleep(200 * time.Millisecond)
		for {
			select {
			case <-resetASN:
				ASN = 0
				return
			case <-time.After(SLOT_DURATION * time.Microsecond):
				ASN++
				for i := 0; i < cap(NEW_SLOT_SIGNAL); i++ {
					NEW_SLOT_SIGNAL <- ASN
				}
			}
		}
	}()

	for _, sw := range Switches {
		go sw.Start()
	}

	return nil
}

// for routing table generation
type TopoGraph struct {
	Nodes []*TopologyGraphNode
}

type TopologyGraphNode struct {
	name      string
	neighbors []*TopologyGraphNode
}

func (g *TopoGraph) construct(topo TopologyData) {
	// construct graph
	for _, n := range topo.Nodes {
		node := new(TopologyGraphNode)
		node.name = n.Name
		g.Nodes = append(g.Nodes, node)
	}
	for _, e := range topo.Edges {
		for _, n := range g.Nodes {
			if n.name == e[0] {
				for _, nn := range g.Nodes {
					if nn.name == e[1] {
						n.neighbors = append(n.neighbors, nn)
						nn.neighbors = append(nn.neighbors, n)
					}
				}
			}
		}
	}
}

// FindAllPaths returns all paths from src to dst
func (g *TopoGraph) FindAllPaths(src, dst string) [][]string {
	visited := make(map[string]bool)
	path := []string{}
	res := [][]string{}
	g.findPath(src, dst, visited, path, &res)
	sort.SliceStable(res, func(i, j int) bool {
		return len(res[i]) < len(res[j])
	})
	return res
}

// helper function for FindAllPaths, DFS
func (g *TopoGraph) findPath(cur, dst string, visited map[string]bool, path []string, res *[][]string) {
	visited[cur] = true
	path = append(path, cur)

	if cur == dst {
		tmp := make([]string, len(path))
		copy(tmp, path)
		*res = append(*res, tmp)
	} else {
		for _, n := range g.Nodes {
			if n.name == cur {
				for _, nn := range n.neighbors {
					if !visited[nn.name] && (nn.name[:2] == "SW" || nn.name == dst ||
						((n.name == "GCC" || dst == "GCC") && nn.name == "HMS")) {
						g.findPath(nn.name, dst, visited, path, res)
					}
				}
			}
		}
	}
	visited[cur] = false
}
