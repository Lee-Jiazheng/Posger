package tika

import (
	"io/ioutil"
	"net/http"
	"os"
)

var (
	tikaServerUrl string
	staticDir     string
)

func init() {
	tikaServerUrl = "http://localhost:9998/"
}

func putRequest(url, filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPut, url, file)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	return string(b), nil
}

func GetPdfContent(filename string) (string, error) {
	url := tikaServerUrl + "tika"
	return putRequest(url, filename)
}

// Waiting, maybe json or xml
func getPdfMeta(filename string) (string, error) {
	url := tikaServerUrl + "meta"
	return putRequest(url, filename)
}
