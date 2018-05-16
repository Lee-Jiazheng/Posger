package Posger

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"log"
)

const (
	LOG_PATH = "server.log"
)

var (
	Logger	*log.Logger
)

func init() {
	file, err := os.Create(LOG_PATH)
	if err != nil {
		log.Println("Failed to create " + LOG_PATH + " file.")
	}
	Logger = log.New(file, "", log.LstdFlags|log.Llongfile)
	log.Println("Logger is starting")
}

type Error struct {
	error string
}

func (self Error) Error() string {
	return self.error
}

func Must(data interface{}, err error) interface{} {
	if err != nil {
		Logger.Println(err)
	}
	return data
}

// purifyContent puries a list of string, delete invisible characters.
func purifyContent(contents ...string) (res []string){
	for _, content := range contents {
		//judge type, primary type or
		for _, ch := range []string{"\n", " "} {
			content = strings.Replace(content, ch, "", -1)
		}
		res = append(res, content)
	}
	return
}

func PutRequest(url string, body io.Reader) (string, error) {
	req, _ := http.NewRequest(http.MethodPut, url, body)
	response, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	return string(b), nil
}

// ExtractPdfImages gets JPG images via parse protocol in pdf stream.
func ExtractPdfImages(pdfname string) ([][]byte, error) {
	file, err := os.Open(pdfname)
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	startLoc, resultImageBytes := 0, make([][]byte, 0)
	for {
		streamStart := bytesFind(content, []byte("stream"), startLoc)
		if streamStart == -1 {
			break
		}
		jpgStart := bytesFind(content[:streamStart+20], []byte{0xff, 0xd8}, streamStart)
		if jpgStart == -1 {
			startLoc = streamStart + 20
			continue
		}
		streamEnd := bytesFind(content, []byte("endstream"), jpgStart)
		if streamEnd == -1 {
			return nil, Error{"pdf don't have stream end..."}
		}
		jpgEnd := bytesFind(content, []byte{0xff, 0xd9}, streamEnd-20)
		if jpgEnd == -1 {
			return nil, Error{"pdf don't have jpg end..."}
		}
		resultImageBytes = append(resultImageBytes, content[jpgStart+1:jpgEnd+1])
		if err != nil {
			return nil, err
		}
		startLoc = jpgEnd
	}
	return resultImageBytes, nil
}

func bytesFind(bytes []byte, find []byte, startLoc int) int {
	if len(find) == 0 {
		return 0
	}
	var match int
	for i, bt := range bytes[startLoc:] {
		if bt == find[match] {
			match++
		} else {
			match = 0
		}
		if match == len(find) {
			return i + startLoc - len(find)
		}
	}
	return -1
}

func GenerateJPGImage(image []byte, filename string) error {
	err := ioutil.WriteFile(filename, image, 0666)
	return err
}
