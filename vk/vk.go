package vk

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"fmt"
	"encoding/json"

	"../file"
	"../database"
	"time"
)

type Doc struct {
	URL   string
	Title string
}

type Message struct {
	Response struct {
		Count int
		Items []struct {
			Id int
			Attachments []struct {
				Doc Doc
			}
		}
	}
}

const counter = 100
const accessToken = "463cb4750555a2d9609f220bc056c9da82b8645572e9b6b6a766f63551307f5a74cbfc689b12a08b6bcf9"
const apiMethodURL = "https://api.vk.com/method/"

var nilDoc = Doc{}

func Request(methodName string, parameters map[string]string) ([]byte, error) {
	requestURL, err := url.Parse(apiMethodURL + methodName)
	if err != nil {
		return nil, err
	}

	requestQuery := requestURL.Query()

	for key, value := range parameters {
		requestQuery.Set(key, value)
	}

	requestURL.RawQuery = requestQuery.Encode()

	response, err := http.Get(requestURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func GetWall(domain string, counter int, offset int) (int, Message) {
	parameters := make(map[string]string)
	parameters["domain"] = domain
	parameters["count"] = strconv.Itoa(counter)
	parameters["offset"] = strconv.Itoa(offset)
	parameters["v"] = "5.68"
	parameters["access_token"] = accessToken

	fmt.Println("\x1b[32;1m[VK] Get wall request\x1b[0m")

	resp, errVk := Request("wall.get", parameters)
	if errVk != nil {
		fmt.Println("VK error:", errVk)
	}

	var message Message

	fmt.Println("\x1b[32;1m[VK] Decode request\x1b[0m")

	errJson := json.Unmarshal(resp, &message)
	if errJson != nil {
		fmt.Println("JSON error:", errJson)
	}

	count := message.Response.Count

	fmt.Println("\x1b[32;1m[VK] Donloads ", offset, " for ", count, " elements\x1b[0m")

	return count, message
}

func GetBook(offset int, db database.StructureGroups) int {
	count, message := GetWall(db.Domain, counter, offset)

	if offset == 0 {
		fmt.Println("\x1b[32;1m[VK] Elements count:", count, "\x1b[0m")
		database.SetLasIdToGroups(db.Id, message.Response.Items[0].Id)
	}

	for _, element := range message.Response.Items {
		if element.Id == db.LasId {
			return offset
		}

		for _, element := range element.Attachments {
			if element.Doc != nilDoc {
				fmt.Println("\x1b[32;1m[VK] Download:", element.Doc.Title, "\x1b[0m")
				file.Download(element.Doc.Title, element.Doc.URL)
			}
		}
	}

	if count > offset {
		return GetBook(offset+counter, db)
	}

	return offset
}

func Parse()  {
	rows := database.GetGroups()

	fmt.Println("\x1b[32;1m[VK] Start parse vk.com\x1b[0m")

	for _, db := range rows {
		fmt.Println("\x1b[32;1m[VK] Domain: ", db.Domain, "\x1b[0m")
		amt := time.Duration(20)
		time.Sleep(time.Second * amt)

		GetBook(0, db)
	}

	fmt.Println("\x1b[32;1m[VK] Done parse vk.com\x1b[0m")
}
