package main

import (
	"encoding/json"
	"fmt"
	"io"
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
				Type: WSLOG_HEARTBEAT,
				Msg:  "heartbeat",
			}
			WSLog <- l
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

	app.GET("/api/topologies", getTopoTags)
	app.GET("/api/topology", getTopo)
	app.POST("/api/topology", postTopo)
	app.OPTIONS("/api/topology", sgo.PreflightHandler)

	app.GET("/api/stats/:name/io", getStatsIOByName)
	app.GET("/api/stats/:name/delay", getStatsDelayByName)

	app.POST("/api/links", postDefaultSetting)
	app.OPTIONS("/api/links", sgo.PreflightHandler)
	app.POST("/api/link/:name", postLink)
	app.OPTIONS("/api/link/:name", sgo.PreflightHandler)

	app.POST("/api/flows", postFlows)
	app.OPTIONS("/api/flows", sgo.PreflightHandler)
	app.GET("/api/flows/start_flag", getStartedFlag)
	app.GET("/api/flows/stop", stopFlows)

	app.GET("/api/switch/:id", getSwitchSchedule)

	app.POST("/api/fault/switch/:id", postFault)
	app.OPTIONS("/api/fault/switch/:id", sgo.PreflightHandler)
	app.GET("/api/fault/clear", clearFault)

	app.GET("/api/animation", getCurrentAnimationFlag)
	app.GET("/api/animation/:flag", getAnimationFlag)

	app.GET("/api/reroute", getCurrentRerouteFlag)
	app.GET("/api/reroute/:flag", getRerouteFlag)

	app.GET("/api/frer", getCurrentFRERFlag)
	app.GET("/api/frer/:flag", getFRERFlag)
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
		case l := <-WSLog:
			ws.WriteJSON(l)
		case <-breakSig:
			return nil
		}
	}
}

func getCurrentAnimationFlag(ctx *sgo.Context) error {
	return ctx.Text(200, strconv.FormatBool(ANIMATION_ENABLED))
}

func getAnimationFlag(ctx *sgo.Context) error {
	flag := ctx.Param("flag")
	if flag == "true" {
		ANIMATION_ENABLED = true
		fmt.Println("packet animation enabled")
	} else {
		ANIMATION_ENABLED = false
		fmt.Println("packet animation disabled")
	}
	return ctx.Text(200, strconv.FormatBool(ANIMATION_ENABLED))
}

func getCurrentRerouteFlag(ctx *sgo.Context) error {
	return ctx.Text(200, strconv.FormatBool(REROUTE_ENABLED))
}

func getRerouteFlag(ctx *sgo.Context) error {
	flag := ctx.Param("flag")
	if flag == "true" {
		REROUTE_ENABLED = true
		fmt.Println("FRER enabled")
	} else {
		REROUTE_ENABLED = false
		fmt.Println("FRER disabled")
	}
	return ctx.Text(200, strconv.FormatBool(REROUTE_ENABLED))
}

func getCurrentFRERFlag(ctx *sgo.Context) error {
	return ctx.Text(200, strconv.FormatBool(FRER_ENABLED))
}

func getFRERFlag(ctx *sgo.Context) error {
	flag := ctx.Param("flag")
	if flag == "true" {
		FRER_ENABLED = true
		fmt.Println("FRER enabled")
	} else {
		FRER_ENABLED = false
		fmt.Println("FRER disabled")
	}
	return ctx.Text(200, strconv.FormatBool(FRER_ENABLED))
}

func getSwitchSchedule(ctx *sgo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	Switches[id].Neighbors = []string{}
	for i := 0; i < Switches[id].gatesOutIdx; i++ {
		Switches[id].Neighbors = append(Switches[id].Neighbors, Switches[id].gatesOut[i].Neighbor)
	}
	return ctx.JSON(200, 1, "success", Switches[id])
}

func getTopoTags(ctx *sgo.Context) error {
	tags, err := queryTopoTags()
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(tags) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}
	return ctx.JSON(200, 1, "success", tags)
}

func getTopo(ctx *sgo.Context) error {
	topo, err := queryTopo(ctx.Param("tag"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(topo.Nodes) == 0 || len(topo.Edges) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}

	if ctx.Param("tag") != ActiveTopoTag {
		ActiveTopoTag = ctx.Param("tag")
		loadTopo(topo)
	}

	return ctx.JSON(200, 1, "success", topo)
}

func postTopo(ctx *sgo.Context) error {
	body, err := io.ReadAll(ctx.Req.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var topo TopologyData
	fmt.Println(string(body))
	if err = json.Unmarshal(body, &topo); err != nil {
		fmt.Println(err)
		return err
	}
	loadTopo(topo)
	insertTopo(topo)
	return ctx.Text(200, "")
}

func getStatsIOByName(ctx *sgo.Context) error {
	stats, err := queryStatsIOByName(ctx.Param("name"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(stats) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}

	return ctx.JSON(200, 1, "success", stats)
}

func getStatsDelayByName(ctx *sgo.Context) error {
	stats, err := queryStatsDelayByName(ctx.Param("name"))
	if err != nil {
		return ctx.JSON(500, 0, err.Error(), nil)
	}
	if len(stats) == 0 {
		return ctx.JSON(200, 0, "no result found", nil)
	}

	return ctx.JSON(200, 1, "success", stats)
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
					l.Failed, _ = strconv.ParseBool(ctx.Param("failed"))
				}
			}
		}
	}

	return ctx.Text(200, fmt.Sprintf("%d", boottime))
}

// set link properties
func postDefaultSetting(ctx *sgo.Context) error {
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
				l.Distance = d // in km
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
	body, err := io.ReadAll(ctx.Req.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	var flows []Flow
	err = json.Unmarshal(body, &flows)
	if err != nil {
		fmt.Println(err)
		return err
	}
	simStartedFlag = 1
	stopFlowSig = make(chan bool, 56)
	subsysList := []string{"GCC", "HMS", "STR", "SPL", "PWR", "ECLSS", "AGT", "IE", "EXT", "DTB", "COORD"}
	for _, f := range flows {
		for _, subsys := range Subsystems {
			if subsys.Name() == f.Name {
				for i, flag := range f.Dst {
					if flag == "X" {
						go func(dstID int, f Flow) {
							for {
								select {
								case <-stopFlowSig:
									return
								default:
									subsys.CreateFlow(dstID)
									if freq, err := strconv.ParseFloat(f.Freq, 64); err == nil {
										time.Sleep(time.Duration(1/freq*1000*1000*1000) * time.Nanosecond)
									} else {
										fmt.Println(err)
										return
									}

								}
							}

						}(subsysName2ID(subsysList[i]), f)
						// delay between different flows
						time.Sleep(200 * time.Millisecond)
					}
				}
				break
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
	WSLog <- Log{
		Type: WSLOG_MSG,
		Msg:  fmt.Sprintf("Inject %s fault on SW%d, duration: %d s\n", ctx.Param("type"), id, duration),
	}
	return ctx.Text(200, "")
}

func clearFault(ctx *sgo.Context) error {
	for _, sw := range Switches {
		sw.Failed = false
	}

	fmt.Println("Faults cleared")
	WSLog <- Log{
		Type: WSLOG_MSG,
		Msg:  "Faults cleared",
	}
	return ctx.Text(200, "")
}
