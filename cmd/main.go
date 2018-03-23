package main

import (
	. ".."
	"fmt"
)

func main() {
	article, _ := NewArticle("static/articles/大数据时代我国企业财务共享中心的优化.pdf")
	sum := article.Summary()

	fmt.Print(sum)
	keys := article.GetKeyWords()
	fmt.Print(keys)
	RunServer()
}
