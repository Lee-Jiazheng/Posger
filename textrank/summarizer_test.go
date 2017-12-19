package textrank

import "testing"

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
	article, err :=  NewArticle(baseDir + "jetbrains.txt", "en")
	if err != nil {
		t.Error(err)
	}
	article.Summarizer()
}