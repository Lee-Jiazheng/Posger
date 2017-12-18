package textrank

import "testing"

func TestSummarizer(t *testing.T) {
	baseDir := "../static/articles/"
	files := []string{""}
	errors := ArticlesBatchIDF([]string{"../static/articles/test.pdf"})
	if len(errors) != 0{
		t.Error(errors)
	}

}