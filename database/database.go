package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type StructureGroups struct {
	Id     int
	Domain string
	LasId  int
}

var database *sql.DB

func OpenDB() *sql.DB {
	if database == nil {

		fmt.Println("\x1b[30;1m[DB] Open data base\x1b[0m")

		database, _ = sql.Open("sqlite3", "./parse.db")
	}

	return database
}

func GetGroups() []StructureGroups {
	database = OpenDB()

	fmt.Println("\x1b[30;1m[DB] Get get rows from groups\x1b[0m")

	rows, _ := database.Query("SELECT id, domain, las_id FROM groups")

	var elements []StructureGroups

	for rows.Next() {
		element := new(StructureGroups)
		rows.Scan(&element.Id, &element.Domain, &element.LasId)
		elements = append(elements, *element)
	}

	return elements
}

func SetLasIdToGroups (dataBaseId int, lasId int) {
	database = OpenDB()

	fmt.Println("\x1b[30;1m[DB] Set las id to group: ", lasId, "\x1b[0m")

	statement, _ := database.Prepare("UPDATE groups SET las_id = ? WHERE id = ?")
	statement.Exec(lasId, dataBaseId)
}