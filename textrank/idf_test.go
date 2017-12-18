package textrank

import (
	"testing"
)

func TestPdfDirectory(t *testing.T) {
	errors := ArticlesBatchIDF([]string{"../static/articles/test.pdf"})
	if len(errors) != 0{
		t.Error(errors)
	}
}

func TestSaveIDF(t *testing.T) {
	errors := ArticlesBatchIDF([]string{"../static/articles/test.pdf"})
	if len(errors) != 0{
		t.Error(errors)
	}
	err := SaveWordDictionary()
	if err != nil {
		t.Error(err)
	}
}