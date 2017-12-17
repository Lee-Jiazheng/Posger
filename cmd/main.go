package main

import (
	"../textrank"
	"fmt"
)

func main() {
	errors := textrank.ArticlesBatchIDF([]string{`E:\Posger\static\articles\test.pdf`})
	fmt.Print(errors)
}