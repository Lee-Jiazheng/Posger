package Posger

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	//"gopkg.in/mgo.v2"
	"fmt"
	"log"
	"gonum.org/v1/gonum/mat"
	"math"
	"sort"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	articleCount   int64
	wordDictionary = make(map[string]int)
)

func init() {
	idfFile, _ := os.Open("tools/idf.dict")
	defer idfFile.Close()
	buffer := bufio.NewReader(idfFile)
	line, _ := buffer.ReadString('\n')
	articleCount, _ = strconv.ParseInt(line[:len(line)-1], 10, 32)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		res := strings.Split(line[:len(line)-1], " ")
		idf, _ := strconv.ParseInt(res[1], 10, 32)
		wordDictionary[res[0]] = int(idf)
	}
	log.Println("wordDict loading end, word count is " + fmt.Sprintf("%d", articleCount))
}

// PathPrefix is api, uploading files api and so on.
func registeSummaryApi(router *mux.Router) {
	router.HandleFunc("/paper", getPaperApi).Methods("GET")
	//router.HandleFunc("/paper/{paperId}", ).Methods("GET")
	//router.HandleFunc("/paper", addPaper).Methods("POST")
	router.HandleFunc("/paper", del).Methods("DELETE")
}


func getPaper(username string) (string){
	return "login have been"
}

func getPaperApi(w http.ResponseWriter, r *http.Request) {
	RequireLoginApi(w, r, getPaper)
}

func del(w http.ResponseWriter, r *http.Request) {

}





type baseArticle struct {
	Abstract	string		// summary by author, abstract in paper
	Keywords	[]string	// keywords in paper
	References	[]string	// multi references
	Title		string		// article's title
}

type jsonArticle struct {
	baseArticle
	content		[]string	// summary Content
}

func NewJsonArticle(filepath string) (*jsonArticle, error) {
	article, err := NewArticle(filepath)
	if err != nil {
		return nil, err
	}
	return &jsonArticle{baseArticle: article.baseArticle, content: article.Summary()}, nil
}

// NewArticle use filepath parameter to construct a
// segmented article
func NewArticle(filepath string) (*Article, error) {
	article := &Article{Sentences: make([]Sentence, 0), filepath: filepath}
	content, err := ExtractDocumentContent(filepath)
	if err != nil {
		return nil, err
	}
	article.setMeta(content)
	fmt.Print(content)
	article.segmentation(content)
	return article, nil
}

// segmentation will segment an article to
// an array of multi sentences
// Warning: Now Only Chinese
func (self *Article) segmentation(content string) {
	jiebaSentenceSegmentation(self, content)
}

// By analysising the rules, we can get the infomation
// i.e. author, title, abstract, reference, acknowledge
func (self *Article) setMeta(content string) (c string){
	// Consider the condition that we can't find below "keywords"

	sum_s, sum_e := self.getSummaryIndex(content)
	self.setSummaryIndex(content[sum_s:sum_e])
	content = content[sum_e:]

	key_s, key_e := self.getKeywordsIndex(content[:])
	self.setKeywordsIndex(content[key_s:key_e])

	ref_s, ref_e := self.getReferenceIndex(content)
	self.setReferenceIndex(content[ref_s:ref_e])

	return content[key_e: ref_s]
}

func (self *Article) getSummaryIndex(content string) (s, e int) {
	s = strings.LastIndex(content[:strings.Index(content, "摘 要")], "\n")
	e = strings.LastIndex(content[:strings.Index(content, "关键词")], "\n") + 1
	return
}

func (self *Article) setSummaryIndex(abstract string) {
	self.Abstract = strings.Replace(abstract, "\n", "", -1)
}

func (self *Article) setTitleAndAuthor(content string) {

}

func (self *Article) getKeywordsIndex(content string) (s, e int) {
	e = strings.Index(content[:], "\n")
	return
}

func (self *Article) setKeywordsIndex(keywords string) {
	self.Keywords = strings.Split(keywords, " ")
}

func (self *Article) getReferenceIndex(content string) (s, e int) {
	s = strings.LastIndex(content[:strings.LastIndex(content, "参考文献")], "\n") + 1
	e = len(content)
	return
}

func (self *Article) setReferenceIndex(reference string) {
	// First is 参考文献
	// By [1] [2] ... rule consists an array.
	var pos []int; res := strings.Split(reference, "\n")
	for i, r := range res[1:] {
		if strings.HasPrefix(r, fmt.Sprintf("[%d]", len(pos) + 1)) {
			pos = append(pos, i)
		}
	}
	pos = append(pos, len(res))
	for i, p := range pos[:len(pos)-1] {
		self.References = append(self.References, strings.Join(purifyContent(res[p:pos[i+1]]...), ""))
	}

}

type Sentence struct {
	Words []string
	Point float64
}

type Article struct {
	baseArticle
	filepath  string
	Sentences []Sentence
}

func getWordIdf(word string) float64 {
	idf := math.Log((float64(articleCount) + 0.5) / (float64(wordDictionary[word]) + 0.5))
	return idf
}

const SUMMARY_ADJUST = 0.85

// Summarizer belongs to Article
// through similarity matrix, we
// can get the score matrix eventually.
func (self *Article) Summary() []string {
	const ITER = 10
	n := len(self.Sentences)
	similarityMatrix := self.similarityMatrix()
	PR := mat.NewDense(n, 1, make([]float64, n))

	for c := 0; c < ITER; c++ {
		// PR = 0.15 + 0.85 * M * PR
		// and similarityMatrix is 0.85*M
		MmulPR := new(mat.Dense)
		MmulPR.Mul(similarityMatrix, PR)
		PR.Add(wholeMatrix(n, 1, 1-SUMMARY_ADJUST), MmulPR)
		// if converagence, break, the threshold is 0.0001
	}

	resultMap := sortSentences(PR)
	result := make([]string, n)
	for i, v := range *resultMap {
		// assemble all the words, append to the result.
		result[i] = connectSliceWords(self.Sentences[v.Key].Words...)
	}
	return purifyContent(result...)
}

func sortSentences(PR *mat.Dense) *ByPoint {
	// sort sentences points
	result_map := make(ByPoint, len(PR.RawMatrix().Data))
	for p, idf := range PR.RawMatrix().Data {
		result_map[p] = PointMap{p, idf}
	}
	sort.Sort(sort.Reverse(result_map))
	return &result_map
}

func connectSliceWords(words ...string) string {
	var result string
	for _, word := range words {
		result += word + " "
	}
	return result
}

// simlarityMatrix is a score matrix by BM25 alorgithm
// it will record every dual sentences similarity scores.
func (self *Article) similarityMatrix() *mat.Dense {
	n := len(self.Sentences)
	data := make([]float64, n*n)
	for i1, s1 := range self.Sentences {
		s1 := s1
		for i2, s2 := range self.Sentences {
			s2 := s2
			if i1 == i2 {
				// diagonal Set Zero
				data[i1*n+i2] = 0
			} else {
				data[i1*n+i2] = s1.BM25(&s2) * SUMMARY_ADJUST
			}
		}
	}
	return mat.NewDense(n, n, data)
}

// BM25 caculates the similarity score of two sentences.
func (self *Sentence) BM25(s *Sentence) float64 {
	const k1, b, avsl = 2.0, 0.75, 5.0
	var score float64
	for _, word := range self.Words {
		var tf int
		for _, wordc := range s.Words {
			if word == wordc {
				tf += 1
			}
		}
		var B float64 = (1 - b) + b*(float64(len(s.Words))/avsl)
		score += getWordIdf(word) * (float64(tf*(k1+1)) / (float64(tf) + k1*B))
	}
	return score
}

// wholeMatrix will return a matrix(mat.Dense)
// that all member is "num" parameter.
// Its size is designited by raw and col.
func wholeMatrix(raw int, col int, num float64) *mat.Dense {
	data := make([]float64, raw*col)
	for i, _ := range data {
		data[i] = num
	}
	return mat.NewDense(raw, col, data)
}

type PointMap struct {
	Key   int     //Sentence Position(Index)
	Value float64 //Sentence Points
}

type ByPoint []PointMap

func (idf ByPoint) Len() int           { return len(idf) }
func (idf ByPoint) Swap(i, j int)      { idf[i], idf[j] = idf[j], idf[i] }
func (idf ByPoint) Less(i, j int) bool { return idf[i].Value < idf[j].Value }

// We also need keyword.
// keyword
// keyword
// keyword

func (self *Article) GetKeyWords() []string {
	var resultMap = make(ByWordIdf, 0)
	tokenNum, statistics := self.TokenFrequencyStat()
	for _, sentence := range self.Sentences {
		for _, word := range sentence.Words {
			if isKeywordMapExist(&resultMap, word) {
				continue
			} // Later judge stop words
			resultMap = append(resultMap, WordIdfMap{word, getWordIdf(word) * float64(statistics[word]) / float64(tokenNum)})
		}
	}
	sort.Sort(sort.Reverse(resultMap))

	result := make([]string, 0)
	for _, v := range resultMap {
		result = append(result, v.Word)
	}
	return result
}

func (self *Article) TokenFrequencyStat() (wordCount int, countMap map[string]int) {
	countMap = make(map[string]int)
	for _, sentence := range self.Sentences {
		for _, word := range sentence.Words {
			countMap[word]++
			wordCount++
		}
	}
	return
}

func isKeywordMapExist(m *ByWordIdf, word string) bool {
	for _, idfs := range *m {
		if word == idfs.Word {
			return true
		}
	}
	return false
}

type WordIdfMap struct {
	Word  string  // Word Content
	Value float64 // Word's IDF Point
}

type ByWordIdf []WordIdfMap

func (idf ByWordIdf) Len() int           { return len(idf) }
func (idf ByWordIdf) Swap(i, j int)      { idf[i], idf[j] = idf[j], idf[i] }
func (idf ByWordIdf) Less(i, j int) bool { return idf[i].Value < idf[j].Value }
