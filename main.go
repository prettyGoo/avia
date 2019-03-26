// proj project main.go
package main

import (
	"database/sql"
	"fmt"
	"strings"

	"proj/src/csv_parser"
	"proj/src/init_db"
	"proj/src/json_parser"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = ""
	DB_PASSWORD = ""
	DB_NAME     = "avia"
	FILE_NAME   = "data.csv"
)

func main() {
	var connParams string = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", connParams)
	if err != nil {
		panic(err)
	}

	init_db.InitDB(db, false)
	if strings.Contains(FILE_NAME, ".csv") {
		csv_parser.Parse(FILE_NAME, db)
	} else if strings.Contains(FILE_NAME, ".json") {
		json_parser.Parse(FILE_NAME, db)
	}

	defer db.Close()
}
