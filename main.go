package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"./file"
	"./vk"
	"./proxy"
	"./database"
)

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

const accessToken = "463cb4750555a2d9609f220bc056c9da82b8645572e9b6b6a766f63551307f5a74cbfc689b12a08b6bcf9"

func getBook(offset int, db database.StructureGroups) int {
	counter := 10

	parameters := make(map[string]string)
	parameters["domain"] = db.Domain
	parameters["count"] = strconv.Itoa(counter)
	parameters["offset"] = strconv.Itoa(offset)
	parameters["v"] = "5.68"
	parameters["access_token"] = accessToken

	fmt.Print("Get 'wall.search' request \n")
	resp, errVk := vk.Request("wall.get", parameters)
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
		fmt.Print("Elements count:", count, "\n\n")
		database.SetLasIdToGroups(db.Id, message.Response.Items[0].Id)
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

func main() {
	proxy.Set()
	rows := database.GetGroups()

	for _, db := range rows {
		fmt.Println("Domain: ", db.Domain, " Query: ", db.Query)
		getBook(0, db)
	}

	fmt.Println("Done ALL...")
}
