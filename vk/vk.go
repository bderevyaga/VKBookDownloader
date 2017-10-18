package vk

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiMethodURL = "https://api.vk.com/method/"

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
