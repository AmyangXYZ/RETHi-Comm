package main

import "time"

func getSeqNum() int32 {
	tmp := SequenceNumber
	SequenceNumber++
	return tmp
}

func subsysID2Name(id uint8) string {
	return SUBSYS_LIST[int(id)]
}

func subsysName2ID(name string) int {
	for i, n := range SUBSYS_LIST {
		if n == name {
			return i
		}
	}
	return -1
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
