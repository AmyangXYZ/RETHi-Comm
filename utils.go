package main

import (
	"sync"
	"time"
)

var seqNumIncMutex sync.Mutex // for sequence number increment
var UIDIncMutex sync.Mutex    // for packet UID increment

// get the current sequence number and increment it
func getSeqNum() int32 {
	seqNumIncMutex.Lock()
	tmp := SequenceNumber
	SequenceNumber++
	seqNumIncMutex.Unlock()
	return tmp
}

// get the current packet UID and increment it
func getUID() int {
	UIDIncMutex.Lock()
	tmp := UID
	UID++
	UIDIncMutex.Unlock()
	return tmp
}

// convert subsys id to name
func subsysID2Name(id uint8) string {
	for n, i := range SUBSYS_MAP {
		if id == i {
			return n
		}
	}
	return ""
}

// convert subsys name to id
func subsysName2ID(name string) int {
	return int(SUBSYS_MAP[name])
}

// find node by name in Subsystems and Switches
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
	return &Switch{name: ""}
}

// Collect statistics and send to frontend via websocket
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
			saveStatsIO(tmp)
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

		WSLog <- Log{
			Type:       WSLOG_STAT,
			Msg:        "",
			Statistics: tmp,
		}
		time.Sleep(UPLOAD_STATS_PERIOD * time.Second)
	}
}
