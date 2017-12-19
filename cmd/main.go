package main

import (
	"../textrank"
	"fmt"
)

func main() {

	baseDir := `E:\Posger\static\articles\`
	files := []string{"BMC.pdf", "Modeling.pdf", "test.pdf"}
	for i, f := range files {
		files[i] = baseDir + f
	}
	errors := textrank.ArticlesBatchIDF(files)
	if len(errors) != 0{
		fmt.Println("error")
	}
	article, err :=  textrank.NewArticle(baseDir + "jetbrains.txt", "en")
	if err != nil {
		fmt.Println("error")
	}
	rank := article.Summarizer()
	for _, sentence := range rank {
		fmt.Println(sentence)
	}
}