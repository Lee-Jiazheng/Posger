package Posger

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestWordDictionaryIdf(t *testing.T) {
	//assertEqual(t, wordDictionary["会话"], 393, "会话-"+string(wordDictionary["会话"]))
	//assertEqual(t, wordDictionary["网页"], 6609, "网页-"+string(wordDictionary["网页"]))
	//assertEqual(t, articleCount, 91208, "article count is "+ fmt.Sprintf("%d", articleCount))
}

func TestSummary(t *testing.T) {
	article, err := NewArticle("static/articles/大数据时代我国企业财务共享中心的优化.pdf")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(article.Summary()[1])
}

// about 0.7s/op
func BenchmarkNewJsonArticle(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewJsonArticle("static/articles/大数据时代我国企业财务共享中心的优化.pdf")
		if err != nil {
			b.Error(err)
		}
	}

}