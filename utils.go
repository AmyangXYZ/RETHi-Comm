package main

import (
	"sync"
	"time"
)

var seqNumIncMutex sync.Mutex
var UIDIncMutex sync.Mutex

func getSeqNum() int32 {
	seqNumIncMutex.Lock()
	tmp := SequenceNumber
	SequenceNumber++
	seqNumIncMutex.Unlock()
	return tmp
}

func getUID() int {
	UIDIncMutex.Lock()
	tmp := UID
	UID++
	UIDIncMutex.Unlock()
	return tmp
}

func subsysID2Name(id uint8) string {
	for n, i := range SUBSYS_MAP {
		if id == i {
			return n
		}
	}
	return ""
}

func subsysName2ID(name string) int {
	return int(SUBSYS_MAP[name])
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
