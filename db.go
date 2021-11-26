package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("sqlite3", "./comm.db")
	if err != nil {
		panic(err)
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
		NODE0 VARCHAR(16) NOT NULL,
		NODE1 VARCHAR(16) NOT NULL);`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STATISTICS (
		ID SMALLINT,
		Name VARCHAR(16) NOT NULL,
		SOURCE INT,
		DELAY DOUBLE);`)
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

	stmtEdges, err := db.Prepare(`INSERT INTO TOPOLOGY_EDGES (TAG, NODE0, NODE1) VALUES (?, ?, ?);`)
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
	rowsEdges, err = db.Query(`SELECT NODE0, NODE1 FROM TOPOLOGY_EDGES where TAG=?`, tag)
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
