package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type StructureGroups struct {
	Id     int
	Domain string
	Query  string
	LasId  int
}

var database *sql.DB

func GetDB() *sql.DB {
	if database == nil {
		fmt.Print("Get data base \n\n")

		database, _ = sql.Open("sqlite3", "./parse.db")
	}

	return database
}

func GetGroups() []StructureGroups {
	database = GetDB()

	fmt.Print("Get get rows from groups \n\n")

	rows, _ := database.Query("SELECT id, domain, query, las_id FROM groups")

	var elements []StructureGroups

	for rows.Next() {
		element := new(StructureGroups)
		rows.Scan(&element.Id, &element.Domain, &element.Query, &element.LasId)
		elements = append(elements, *element)
	}

	return elements
}

func SetLasIdToGroups (dataBaseId int, lasId int) {
	database = GetDB()

	fmt.Print("Set las id to group: ", lasId, "\n\n")

	statement, _ := database.Prepare("UPDATE groups SET las_id = ? WHERE id = ?")
	statement.Exec(lasId, dataBaseId)
}