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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STATISTICS (
		TIMESTAMP BIGINT,
		NAME VARCHAR(16) NOT NULL,
		TX INT,
		RX INT);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`DELETE FROM STATISTICS WHERE 1`)
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

func saveStats(data map[string][2]int) {
	t := time.Now()
	timestamp := t.UnixNano() / 1e6
	stmt, err := db.Prepare(`INSERT INTO STATISTICS (TIMESTAMP, NAME, TX, RX) VALUES (?, ?, ?, ?);`)
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

// timestamp, tx, rx
func queryStatsByName(name string) ([][3]int, error) {
	var stats [][3]int
	var rows *sql.Rows

	rows, err = db.Query(`SELECT TIMESTAMP,TX,RX FROM STATISTICS WHERE NAME=?`, name)
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
