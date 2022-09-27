package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, _ = sql.Open("mysql", fmt.Sprintf("%v:%v@(comm_db:3306)/%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DB")))

	// wait database container to start
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(2 * time.Second)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS TOPOLOGY_NODES (
		TAG VARCHAR(64) NOT NULL,
		NAME VARCHAR(16) NOT NULL,
		POS_X INT NOT NULL,
		POS_Y INT NOT NULL);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS TOPOLOGY_EDGES (
		TAG VARCHAR(64) NOT NULL,
		SOURCE VARCHAR(16) NOT NULL,
		TARGET VARCHAR(16) NOT NULL);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STATISTICS_IO (
		TIMESTAMP BIGINT,
		NAME VARCHAR(16) NOT NULL,
		TX INT,
		RX INT);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STATISTICS_DELAY (
		TIMESTAMP BIGINT,
		NAME VARCHAR(16) NOT NULL,
		SOURCE VARCHAR(16) NOT NULL,
		SEQ INT,
		DELAY DOUBLE);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`DELETE FROM STATISTICS_IO WHERE 1`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`DELETE FROM STATISTICS_DELAY WHERE 1`)
	if err != nil {
		panic(err)
	}
}

func insertTopo(topo TopologyData) {
	db.Exec(fmt.Sprintf("DELETE FROM TOPOLOGY_NODES WHERE TAG=\"%s\"", topo.Tag))
	db.Exec(fmt.Sprintf("DELETE FROM TOPOLOGY_EDGES WHERE TAG=\"%s\"", topo.Tag))

	stmtNodes, err := db.Prepare(`INSERT INTO TOPOLOGY_NODES (TAG, NAME, POS_X, POS_Y) VALUES (?, ?, ?, ?);`)
	if err != nil {
		fmt.Println(err)
		return
	}

	stmtEdges, err := db.Prepare(`INSERT INTO TOPOLOGY_EDGES (TAG, SOURCE, TARGET) VALUES (?, ?, ?);`)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmtNodes.Close()
	defer stmtEdges.Close()

	for _, n := range topo.Nodes {

		_, err = stmtNodes.Exec(topo.Tag, n.Name, n.Position[0], n.Position[1])
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, e := range topo.Edges {
		_, err = stmtEdges.Exec(topo.Tag, e[0], e[1])
		if err != nil {
			fmt.Println(err)
		}
	}
}

func queryTopo(tag string) (TopologyData, error) {
	var topo TopologyData
	var node TopologyNode
	var edge = [2]string{}
	var rowsNodes *sql.Rows
	var rowsEdges *sql.Rows

	topo.Tag = tag
	rowsNodes, err := db.Query(`SELECT NAME,POS_X, POS_Y FROM TOPOLOGY_NODES where TAG=?`, tag)
	if err != nil {
		return topo, err
	}
	defer rowsNodes.Close()
	for rowsNodes.Next() {
		rowsNodes.Scan(&node.Name, &node.Position[0], &node.Position[1])
		topo.Nodes = append(topo.Nodes, node)
	}
	rowsEdges, err = db.Query(`SELECT SOURCE, TARGET FROM TOPOLOGY_EDGES where TAG=?`, tag)
	if err != nil {
		return topo, err
	}
	defer rowsEdges.Close()
	for rowsEdges.Next() {
		rowsEdges.Scan(&edge[0], &edge[1])
		topo.Edges = append(topo.Edges, edge)
	}
	return topo, nil
}

func queryTopoTags() ([]string, error) {
	var tags []string
	var rows *sql.Rows

	rows, err = db.Query(`SELECT TAG FROM TOPOLOGY_NODES GROUP BY TAG`)
	if err != nil {
		return tags, err
	}

	defer rows.Close()
	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	return tags, nil
}

func saveStatsIO(data map[string][2]int) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6
	stmt, err := db.Prepare(`INSERT INTO STATISTICS_IO (TIMESTAMP, NAME, TX, RX) VALUES (?, ?, ?, ?);`)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()

	for name, value := range data {
		_, err = stmt.Exec(timestamp, name, value[0], value[1])
		if err != nil {
			fmt.Println(err)
		}
	}
}

func saveStatsDelay(name, source string, seq int32, delay float64) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6
	stmt, err := db.Prepare(`INSERT INTO STATISTICS_DELAY (TIMESTAMP, NAME, SOURCE, SEQ, DELAY) VALUES (?, ?, ?, ?, ?);`)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(timestamp, name, source, seq, delay)
	if err != nil {
		fmt.Println(err)
	}
}

// timestamp, tx, rx
func queryStatsIOByName(name string) ([][3]int, error) {
	var stats [][3]int
	var rows *sql.Rows

	rows, err = db.Query(`SELECT TIMESTAMP,TX,RX FROM STATISTICS_IO WHERE NAME=?`, name)
	if err != nil {
		return stats, err
	}

	defer rows.Close()
	for rows.Next() {
		var data [3]int
		rows.Scan(&data[0], &data[1], &data[2])
		stats = append(stats, data)
	}
	return stats, nil
}

type StatsDelay struct {
	Source    string    `json:"source"`
	Timestamp []int     `json:"timestamp"`
	Delay     []float64 `json:"delay"`
}

func queryStatsDelayByName(name string) ([]StatsDelay, error) {
	var stats []StatsDelay

	for subsys := range SUBSYS_MAP {
		if subsys == name {
			continue
		}
		var rows *sql.Rows

		rows, err = db.Query(`SELECT TIMESTAMP, DELAY FROM STATISTICS_DELAY WHERE NAME=? and SOURCE=?`, name, subsys)
		if err != nil {
			return stats, err
		}
		var entry StatsDelay
		entry.Source = subsys
		defer rows.Close()
		for rows.Next() {
			var ts int
			var delay float64
			rows.Scan(&ts, &delay)
			entry.Timestamp = append(entry.Timestamp, ts)
			entry.Delay = append(entry.Delay, delay)
		}
		if len(entry.Timestamp) > 0 {
			stats = append(stats, entry)
		}
	}

	return stats, nil
}
