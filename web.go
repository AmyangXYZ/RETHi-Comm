package main

import (
	"encoding/json"
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
	stopFlowSig chan bool
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
	app.POST("/api/links", postDefaultSetting)
	app.OPTIONS("/api/links", sgo.PreflightHandler)

	app.POST("/api/link/:name", postLink)
	app.OPTIONS("/api/link/:name", sgo.PreflightHandler)

	app.POST("/api/flows", postFlows)
	app.OPTIONS("/api/flows", sgo.PreflightHandler)
	app.GET("/api/flows/stop", stopFlows)

	app.POST("/api/fault/switch/:id", postFault)
	app.OPTIONS("/api/fault/switch/:id", sgo.PreflightHandler)

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

func stopFlows(ctx *sgo.Context) error {
	for i := 0; i < cap(stopFlowSig); i++ {
		stopFlowSig <- true
	}
	return ctx.Text(200, "stop all flows")
}

func postFault(ctx *sgo.Context) error {
	fmt.Println(ctx.Params())
	id, _ := strconv.Atoi(ctx.Param("id"))
	Switches[id].Failed = true
	return ctx.Text(200, "")
}
