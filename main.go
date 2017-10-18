package main

import (
	"encoding/json"
	"fmt"
	"os"

	"./file"
	"./vk"
)

type Doc struct {
	URL   string
	Title string
}

type Message struct {
	Response struct {
		Count int
		Items []struct {
			Attachments []struct {
				Doc Doc
			}
		}
	}
}

var nilDoc = Doc{}

func main() {
	fmt.Println("Set https proxy")
	os.Setenv("HTTPS_PROXY", "https://193.37.152.6:3128")

	parameters := make(map[string]string)
	parameters["domain"] = "proglib"
	parameters["query"] = "#book@proglib"
	parameters["count"] = "100"
	parameters["v"] = "5.68"
	parameters["access_token"] = "463cb4750555a2d9609f220bc056c9da82b8645572e9b6b6a766f63551307f5a74cbfc689b12a08b6bcf9"

	fmt.Println("Get 'wall.search' request")
	resp, errVk := vk.Request("wall.search", parameters)
	if errVk != nil {
		fmt.Println("VK error:", errVk)
	}

	var message Message

	fmt.Println("Decode request")
	errJson := json.Unmarshal(resp, &message)
	if errJson != nil {
		fmt.Println("JSON error:", errJson)
	}

	for _, element := range message.Response.Items {
		for _, element := range element.Attachments {
			if element.Doc != nilDoc {
				fmt.Println("Download " + element.Doc.Title + "...")

				errDownload := file.Download(element.Doc.Title, element.Doc.URL)
				if errDownload != nil {
					fmt.Println("Download error:", errDownload)
				}

				fmt.Println("Done download \n")
			}
		}
	}

	fmt.Println("Done ALL...")
}
