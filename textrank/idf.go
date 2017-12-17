package textrank

import (
	"fmt"
	"math"
	"../tika"
	"os"
	"bufio"
	"sync"
)

var (
	// trained IDF article count
	articleCount int
	WordDictionary map[string]int
	dictMutex sync.Mutex
)

var WordIDF = make(map[string]float64)


func init() {
	WordDictionary = make(map[string]int)
	articleCount = 0
}

func SaveWordDictionary() error{
	file, err := os.Create("IDF.dict")
	if err != nil { return err}
	defer file.Close()
	w := bufio.NewWriter(file)
	for word, count := range WordDictionary {
		_, err := w.WriteString(fmt.Sprintf(`%s %.6f\n`, word, count))
		if err != nil { return err}
	}
	w.Flush()
	return nil
}

func GetWordIDF(word string) float64 {
	idf := math.Log( (float64(articleCount) + 0.5) / (float64(WordDictionary[word]) + 0.5) )
	return idf
}

func ArticlesBatchIDF(filenames []string) []error {
	errCh := make(chan error)
	errors := make([]error, 0)
	for _, filename := range filenames {
		go articleSingleIDF(filename, errCh)
	}

	for range filenames{
		err := <-errCh
		if err != nil { errors = append(errors, err)}
	}
	fmt.Print(WordDictionary)
	return errors
}

func articleSingleIDF(filename string, errCh chan<- error) {
	content, err := tika.GetPdfContent(filename)
	if err != nil { errCh <- err; return }
	article := Article{Sentences:make([]Sentence, 0), content:content}
	article.segementSentence("en")
	words := article.countWord()
	processWordDict(words)
	articleCount++
	errCh <- nil
}

func processWordDict(words []string) {
	dictMutex.Lock()
	for _, word := range words {
		WordDictionary[word] += 1
	}
	dictMutex.Unlock()
}
