package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func init() {
	db, _ = sql.Open("mysql", fmt.Sprintf("%v:%v@(hms_db:3306)/%v", // "hms_db" is the database container's name in the docker-compose.yml
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

	db.Exec(`CREATE TABLE IF NOT EXISTS FDD_TABLE_1(
		ID INT,
		NAME VARCHAR(16), 
		DATA INT);)`)
}

func main() {
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	select {}
}
