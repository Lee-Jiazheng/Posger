package textrank

import (
	"github.com/yanyiwu/gojieba"
	"strings"
	"os"
	"bufio"
	"io"
	"strconv"
	"io/ioutil"
	"log"
	"../EnSegment"
	"../tika"
)

type Sentence struct {
	Words []string
	Point float64
}

type Article struct {
	Sentences []Sentence
	content   string
}


func split_sentence(char rune) bool {
	split_chars := ".!?" + "。！？" + "\n"
	for _, s := range split_chars {
		if s == char {
			return true
		}
	}
	return false
}

func NewArticle(filepath, lang string) (*Article, error){
	article := Article{}
	content, err :=tika.GetPdfContent(filepath)
	if err != nil { return nil, err}
	article.content = content
	article.segementSentence(lang)
	return &article, nil
}

func (self *Article) segementSentence(lang string) {
	var words []string
	if lang == "en" {
		words = EnSegment.CutAll(self.content)
	} else {
		jieba := gojieba.NewJieba()
		defer jieba.Free()
		log.Print("jieba cutting article...")
		words = jieba.CutAll(self.content)
	}
	self.Sentences = self.cutSentences(&words)
}

func (self *Article) cutSentences(words *[]string) []Sentence{
	var start, end int
	var sentences = make([]Sentence, 0)
	for i, word := range *words {
		if isSegmentSymbol(word) {
			end = i
			sentences = append(sentences, Sentence{Words: (*words)[start:end]})
			start = i
		}
	}
	return sentences
}

func (self *Article)countWord() []string{
	recordDictionary := make(map[string]bool)
	words := make([]string, 0)
	for _, sentence := range self.Sentences {
		for _, word := range sentence.Words {
			if recordDictionary[word] == false {
				recordDictionary[word] = true
				words = append(words, word)
			}
		}
	}
	return words
}

type PointMap struct {
	Key   int	//Sentence Position(Index)
	Value float64	//Sentence Points
}

type ByPoint []PointMap

func (idf ByPoint) Len() int { return len(idf)}
func (idf ByPoint) Swap(i, j int) { idf[i], idf[j] = idf[j], idf[i]}
func (idf ByPoint) Less(i, j int) bool { return idf[i].Value < idf[j].Value}

func run() {
	idfFile, _ := os.Open("/home/lee/articles/IDF.txt")
	defer idfFile.Close()
	buffer := bufio.NewReader(idfFile)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		res := strings.Split(line, " ")
		idf, _ := strconv.ParseFloat(res[1], 64)
		WordIDF[res[0]] = idf
	}
	testFile, err := os.Open("/home/lee/articles/6.txt")
	if err != nil{panic(err)}
	defer testFile.Close()
	c, err := ioutil.ReadAll(testFile)
	// fmt.Println(string(fd))
	content := string(c[:])

	article := Article{Sentences:make([]Sentence, 0), content:string(content)}
	//article.segementSentence()
	article.Summarizer()
}