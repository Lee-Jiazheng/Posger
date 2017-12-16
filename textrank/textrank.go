package textrank

import (
	"math"
	"github.com/yanyiwu/gojieba"
	"strings"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"sort"
	"os"
	"bufio"
	"io"
	"strconv"
	"io/ioutil"
)

type Sentence struct {
	Words []string
	content string
	Point float64
}

type Article struct {
	Sentences []Sentence
	content   string
}

func (self *Sentence) SegementWords() {
	jieba := gojieba.NewJieba()
	defer jieba.Free()
	self.Words = jieba.Cut(self.content, true)
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

func (self *Article) segementSentence() {
	sentences := strings.FieldsFunc(self.content, split_sentence)
	for _, sentence := range sentences {
		sentence := Sentence{content:sentence}
		sentence.SegementWords()
		self.Sentences = append(self.Sentences, sentence)
	}
}

func (self *Article) summary() {

	n := len(self.Sentences)
	data := make([]float64, n * n)
	var t1, t2 int
	for i1, s1 := range self.Sentences {
		s1 := s1
		for i2, s2 := range self.Sentences {
			s2 := s2
			if i1 == i2 {
				data[i1*n + i2] = 0
			} else {
				data[i1*n + i2] = s1.BM25(&s2)
			}
			t2 = i2
		}
		t1 = i1
	}
	fmt.Print(t1, t2)
	similarityMartix := mat.NewDense(n, n, data)

	PR := mat.NewDense(n, 1,  make([]float64, n))
	for c := 0; c < ITER; c++ {
		var temp, temp2 = new(mat.Dense), new(mat.Dense)
		temp.Mul(generateMatrix(n, n,0.85), similarityMartix)
		temp2.Mul(temp, PR)

		PR.Add(generateMatrix2(n, 1, 0.15), temp2)
		// if converagence, break, the threshold is 0.0001
	}
	// sort sentences points
	result := PR.RawMatrix().Data
	result_map := make(ByPoint, len(result))

	for p, idf := range result {
		result_map[p] = PointMap{p, idf}
	}
	sort.Sort(sort.Reverse(result_map))
	for _, v := range result_map {
		fmt.Print(self.Sentences[v.Key].content + fmt.Sprintf("%.3f", v.Value) + "\n")
	}
}

type PointMap struct {
	Key   int	//Sentence Position(Index)
	Value float64	//Sentence Points
}

type ByPoint []PointMap

func (idf ByPoint) Len() int { return len(idf)}
func (idf ByPoint) Swap(i, j int) { idf[i], idf[j] = idf[j], idf[i]}
func (idf ByPoint) Less(i, j int) bool { return idf[i].Value < idf[j].Value}


//对角线为该值， 相当于行列式乘积
func generateMatrix(raw, col int, num float64) *mat.Dense {
	data := make([]float64, raw*col)
	for r := 0; r < raw; r++ {
		for c := 0; c < col; c++ {
			if r == c {
				data[r*col+c] = num
			}
		}
	}
	return mat.NewDense(raw, col, data)
}

// 全为num
func generateMatrix2(raw, col int, num float64) *mat.Dense {
	data := make([]float64, raw*col)
	for r := 0; r < raw; r++ {
		for c := 0; c < col; c++ {
			data[r*col + c] = num
		}
	}
	return mat.NewDense(raw, col, data)
}

func (self *Sentence) BM25(s *Sentence) float64 {
	var score float64
	for _, word := range self.Words {
		word := word
		var tf int
		for _, wordc := range s.Words {
			if word == wordc {
				tf += 1
			}
		}

		var B float64 = (1 - b) + b * (float64(len(s.Words)) / avsl)
		score += IDF(word) * (float64(tf*(k1+1)) / (float64(tf)+k1*B))
	}
	return score
}

func IDF(word string) float64 {
	len := 6
	idf, ok := WordIDF[word]
	if !ok {
		return math.Log( (float64(float64(len) + 0.5) / 0.5 ))
	}
	return idf
}
const (
	k1 = 2.0
	b = 0.75
	avsl = 5.0
	ITER = 10
)

var wordChan = make(chan string)
var WordDictionary = make(map[string]int)	//一篇文章的词典
var WordDocumentary = make(map[string]int)	//所有词的词典
var WordIDF = make(map[string]float64)

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
	article.segementSentence()
	article.summary()
}