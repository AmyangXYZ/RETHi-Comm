package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AmyangXYZ/sgo"
	"github.com/AmyangXYZ/sgo/middlewares"
	"github.com/gorilla/websocket"
)

const (
	addr = ":8000"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	simStartedFlag = 0
	stopFlowSig    chan bool
)

func runHTTPSever() {
	go func() {
		for {
			l := Log{
				Type: -1,
				Msg:  "heartbeat",
			}
			LogsComm <- l
			time.Sleep(10 * time.Second)
		}
	}()

	app := sgo.New()
	app.SetTemplates("templates", nil)
	app.USE(middlewares.CORS(middlewares.CORSOpt{}))

	app.GET("/", index)

	app.GET("/static/*files", static)

	app.GET("/api/boottime", getBootTime)
	app.GET("/ws/comm", wsComm)

	app.POST("/api/topology", postTopo)
	app.PUT("/api/topology", putTopo)
	app.OPTIONS("/api/topology", sgo.PreflightHandler)

	app.POST("/api/links", postDefaultSetting)
	app.OPTIONS("/api/links", sgo.PreflightHandler)
	app.POST("/api/link/:name", postLink)
	app.OPTIONS("/api/link/:name", sgo.PreflightHandler)

	app.POST("/api/flows", postFlows)
	app.OPTIONS("/api/flows", sgo.PreflightHandler)
	app.GET("/api/flows/start_flag", getStartedFlag)
	app.GET("/api/flows/stop", stopFlows)

	app.POST("/api/fault/switch/:id", postFault)
	app.OPTIONS("/api/fault/switch/:id", sgo.PreflightHandler)

	app.GET("/api/fault/clear", clearFault)

	if err := app.Run(addr); err != nil {
		log.Fatal("Listen error", err)
	}
}

// Index page handler.
func index(ctx *sgo.Context) error {
	return ctx.Render(200, "index")
}

// Static files handler.
func static(ctx *sgo.Context) error {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
	return nil
}

func getBootTime(ctx *sgo.Context) error {
	return ctx.Text(200, fmt.Sprintf("%d", boottime))
}

func wsComm(ctx *sgo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Resp, ctx.Req, nil)
	breakSig := make(chan bool)
	if err != nil {
		return err
	}

	fmt.Println("ws/comm connected")
	defer func() {
		ws.Close()
		fmt.Println("ws/comm client closed")
	}()
	go func() {
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				breakSig <- true
			}
		}
	}()
	for {
		select {
		case l := <-LogsComm:
			ws.WriteJSON(l)
		case <-breakSig:
			return nil
		}
	}
}

type TopologyJSON struct {
	Nodes []string    `json:"nodes"`
	Edges [][2]string `json:"edges"`
}

func postTopo(ctx *sgo.Context) error {
	body, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var topo TopologyJSON
	if err = json.Unmarshal(body, &topo); err != nil {
		fmt.Println(err)
		return err
	}
	if len(Switches) > 0 || len(Subsystems) > 0 {
		return errors.New("topology has already been initialized")
	}

	for _, n := range topo.Nodes {
		if n[:2] == "SW" {
			NewSwitch(n)
		} else {
			NewSubsys(n)
		}
	}
	for _, e := range topo.Edges {
		n0 := findNodeByName(e[0])
		n1 := findNodeByName(e[1])
		if n0.Name() != "" && n1.Name() != "" {
			Connect(n0, n1)
		}
	}
	return ctx.Text(200, "")
}

func putTopo(ctx *sgo.Context) error {
	body, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var topo TopologyJSON
	if err = json.Unmarshal(body, &topo); err != nil {
		fmt.Println(err)
		return err
	}

	for _, n := range topo.Nodes {
		if tmp := findNodeByName(n); tmp.Name() != "" {
			return errors.New("node has already been initialized")
		}
	}

	for _, n := range topo.Nodes {
		if n[:2] == "SW" {
			NewSwitch(n)
		} else {
			NewSubsys(n)
		}
	}
	for _, e := range topo.Edges {
		n0 := findNodeByName(e[0])
		n1 := findNodeByName(e[1])
		if n0.Name() != "" && n1.Name() != "" {
			Connect(n0, n1)
		}
	}
	return ctx.Text(200, "")
}

// set link properties
func postLink(ctx *sgo.Context) error {
	fmt.Println(ctx.Params())

	if linkName := ctx.Param("name"); linkName != "" {
		if nodes := strings.Split(linkName, " > "); len(nodes) == 2 {
			for _, l := range Links {
				if (l.sender1.Owner == nodes[0] && l.sink1.Owner == nodes[1]) ||
					(l.sender2.Owner == nodes[0] && l.sink2.Owner == nodes[1]) {
					l.Bandwidth, _ = strconv.ParseFloat(ctx.Param("bandwidth"), 64)
					l.PacketLossRate, _ = strconv.ParseFloat(ctx.Param("loss"), 64)
					l.Distance, _ = strconv.ParseFloat(ctx.Param("distance"), 64)
				}
			}
		}
	}

	return ctx.Text(200, fmt.Sprintf("%d", boottime))
}

// set link properties
func postDefaultSetting(ctx *sgo.Context) error {
	fmt.Println(ctx.Params())

	if ctx.Param("type") == "wired" {
		for _, l := range Links {
			if l.sender1.Owner != "GCC" && l.sender2.Owner != "GCC" {
				b, _ := strconv.ParseFloat(ctx.Param("bandwidth"), 64)
				l.Bandwidth = b * 1024 * 1024 * 1024 // in Gbps
				d, _ := strconv.ParseFloat(ctx.Param("distance"), 64)
				l.Distance = d
			}
		}
	}
	if ctx.Param("type") == "wireless" {
		for _, l := range Links {
			if l.sender1.Owner == "GCC" || l.sender2.Owner == "GCC" {
				b, _ := strconv.ParseFloat(ctx.Param("bandwidth"), 64)
				l.Bandwidth = b * 1024 // in Kbps
				d, _ := strconv.ParseFloat(ctx.Param("distance"), 64)
				l.Distance = d * 1000 // in km
			}
		}
	}

	return ctx.Text(200, fmt.Sprintf("%d", boottime))
}

type Flow struct {
	Name string   `json:"name"` // subsys name
	ID   int      `json:"id"`   // subsys id
	Dst  []string `json:"dst"`
	Freq string   `json:"freq"`
}

func postFlows(ctx *sgo.Context) error {
	body, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println(string(body))
	var flows []Flow
	err = json.Unmarshal(body, &flows)
	if err != nil {
		fmt.Println(err)
		return err
	}
	simStartedFlag = 1
	stopFlowSig = make(chan bool, 56)
	for _, f := range flows {
		subsys := Subsystems[f.ID]
		for i, flag := range f.Dst {
			if flag == "X" {
				go func(dstID int) {
					for {
						select {
						case <-stopFlowSig:
							return
						default:
							subsys.CreateFlow(dstID)
							if freq, err := strconv.ParseFloat(f.Freq, 64); err == nil {
								time.Sleep(time.Duration(1/freq*1000) * time.Millisecond)
							} else {
								fmt.Println(err)
								return
							}

						}
					}

				}(i)
				// delay between different flows
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
	return ctx.Text(200, "biu")
}

func getStartedFlag(ctx *sgo.Context) error {
	return ctx.Text(200, fmt.Sprintf("%d", simStartedFlag))
}

func stopFlows(ctx *sgo.Context) error {
	simStartedFlag = 0
	for i := 0; i < cap(stopFlowSig); i++ {
		stopFlowSig <- true
	}
	return ctx.Text(200, "stop all flows")
}

func postFault(ctx *sgo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	Switches[id].Failed = true
	duration, _ := strconv.Atoi(ctx.Param("duration"))
	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
		Switches[id].Failed = false
	}()

	fmt.Printf("Inject %s fault on SW%d, duration: %d s\n", ctx.Param("type"), id, duration)
	LogsComm <- Log{
		Type: 0,
		Msg:  fmt.Sprintf("Inject %s fault on SW%d, duration: %d s\n", ctx.Param("type"), id, duration),
	}
	return ctx.Text(200, "")
}

func clearFault(ctx *sgo.Context) error {
	for _, sw := range Switches {
		sw.Failed = false
	}

	fmt.Println("Faults cleared")
	LogsComm <- Log{
		Type: 0,
		Msg:  "Faults cleared",
	}
	return ctx.Text(200, "")
}
