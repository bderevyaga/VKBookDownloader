package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"./file"
	"./vk"
)

type DB struct {
	Id     int
	Domain string
	Query  string
	LasId  int
}

type Doc struct {
	URL   string
	Title string
}

type Message struct {
	Response struct {
		Count int
		Items []struct {
			Id          int
			Attachments []struct {
				Doc Doc
			}
		}
	}
}

var nilDoc = Doc{}

const access_token = "463cb4750555a2d9609f220bc056c9da82b8645572e9b6b6a766f63551307f5a74cbfc689b12a08b6bcf9"

func getBook(offset int, db DB) int {
	counter := 10

	parameters := make(map[string]string)
	parameters["domain"] = db.Domain
	parameters["query"] = db.Query
	parameters["count"] = strconv.Itoa(counter)
	parameters["offset"] = strconv.Itoa(offset)
	parameters["v"] = "5.68"
	parameters["access_token"] = access_token

	fmt.Print("Get 'wall.search' request \n")
	resp, errVk := vk.Request("wall.search", parameters)
	if errVk != nil {
		fmt.Println("VK error:", errVk)
	}

	var message Message

	fmt.Print("Decode request \n\n")

	errJson := json.Unmarshal(resp, &message)
	if errJson != nil {
		fmt.Println("JSON error:", errJson)
	}

	count := message.Response.Count

	if offset == 0 {
		statement, _ := database.Prepare("UPDATE groups SET las_id = ? WHERE id = ?")
		statement.Exec(message.Response.Items[0].Id, db.Id)

		fmt.Print("Elements count:", count, "\n\n")
	}

	for _, element := range message.Response.Items {
		if element.Id == db.LasId {
			return offset
		}

		for _, element := range element.Attachments {
			if element.Doc != nilDoc {
				fmt.Println("Download:", element.Doc.Title)

				errDownload := file.Download(element.Doc.Title, element.Doc.URL)
				if errDownload != nil {
					fmt.Println("Download error:", errDownload)
				}

				fmt.Print("Done download \n\n")
			}
		}
	}

	fmt.Print("Donloads", offset, "for", count, "elements... \n\n")

	if count > offset {
		return getBook(offset+counter, db)
	}

	return offset
}

var database *sql.DB

func main() {
	fmt.Print("Set https proxy \n\n")
	os.Setenv("HTTPS_PROXY", "https://193.37.152.6:3128")

	fmt.Print("Open DB \n\n")
	database, _ = sql.Open("sqlite3", "./parse.db")

	var datasDB []DB

	rows, _ := database.Query("SELECT id, domain, query, las_id FROM groups")
	for rows.Next() {
		datas := new(DB)
		rows.Scan(&datas.Id, &datas.Domain, &datas.Query, &datas.LasId)

		datasDB = append(datasDB, *datas)
	}

	for _, db := range datasDB {
		fmt.Println("Domain: ", db.Domain, " Query: ", db.Query)
		getBook(0, db)
	}

	fmt.Println("Done ALL...")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
