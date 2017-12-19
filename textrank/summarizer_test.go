package textrank

import (
	"testing"
	"fmt"
)

func TestSummarizer(t *testing.T) {
	baseDir := "../static/articles/"
	files := []string{"BMC.pdf", "Modeling.pdf", "test.pdf"}
	for i, f := range files {
		files[i] = baseDir + f
	}
	errors := ArticlesBatchIDF(files)
	if len(errors) != 0{
		t.Error(errors)
	}
	article, err :=  NewArticle(baseDir + "Modeling.pdf", "en")
	if err != nil {
		t.Error(err)
	}
	rank := article.Summarizer()
	fmt.Println("\n\n")
	for _, sentence := range rank {
		fmt.Println(sentence)
	}
}