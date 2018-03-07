package Posger

import (
	"os"
	"log"
)

var (
	tikaServerUrl string
)

// TIKA init function should
// 1. command start tika, e.g.: javac **
// 2. read tika config file.
func init() {
	//exec.Command("java -jar /home/lee/workspace/Posger2/tools/tika-server-1.17.jar")
	tikaServerUrl = "http://localhost:9998"
	log.Println("TIKA started at " + tikaServerUrl)
}

func ExtractDocumentContent(filepath string) (string, error) {
	body, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	return PutRequest(tikaServerUrl+"/tika", body)
}
