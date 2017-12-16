package textrank

import (
	"fmt"
	"time"
	"math/big"
	"io/ioutil"
	"math"
)

func trainIDF() {
	articles := make([]Article, 0)
	base_dir := "/home/lee/articles/"
	paths := []string{"1.txt", "2.txt", "3.txt", "4.txt", "5.txt", "6.txt"}
	for _, path := range paths {
		content, _ := ioutil.ReadFile(base_dir + path)
		article := Article{Sentences:make([]Sentence, 0), content:string(content)}
		article.segementSentence()
		articles = append(articles, article)
	}


	for _, article := range articles {
		recordDictionary := make(map[string]bool)
		for _, sentence := range article.Sentences {
			// sentence := sentence
			for _, word := range sentence.Words {

				if recordDictionary[word] == false {
					recordDictionary[word] = true
					WordDocumentary[word] += 1
				}
			}
		}
	}

	content, _ := ioutil.ReadFile(base_dir + "test.txt")
	article := Article{content:string(content), Sentences:nil}
	article.segementSentence()

	fmt.Print(WordDocumentary)

	start := time.Now()

	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))

	var temp []byte

	for word, count := range WordDocumentary {
		idf := math.Log( (float64(float64(len(articles)) + 0.5) / (float64(count) + 0.5) ))
		WordIDF[word] = idf
		temp = append(temp, []byte(word + " " + fmt.Sprintf("%.6f", idf) + "\n")...)
	}

	ioutil.WriteFile(base_dir + "IDF.txt", temp,0666)
	elapsed := time.Since(start)
	fmt.Printf("\nBinomial took %s", elapsed)
}