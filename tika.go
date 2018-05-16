package Posger

import (
	"os"
	"log"
	"os/exec"
	"bytes"
	"fmt"
)

var (
	tikaServerUrl string
)

// TIKA init function should
// 1. command start tika, e.g.: java **
// 2. read tika config file.
func init() {
	//exec.Command("java -jar /home/lee/workspace/Posger2/tools/tika-server-1.17.jar")

	go func() {
		//cmd := exec.Command("/home/lee/Programs/jdk1.8.0_161/bin/java", "-jar", "./tools/tika-server-1.17.jar")
		cmd := exec.Command("/home/lee/Install_Program/jdk1.8.0_151/bin/java", "-jar", "./tools/tika-server-1.17.jar")

		var out bytes.Buffer

		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s", out.String())
	}()

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
