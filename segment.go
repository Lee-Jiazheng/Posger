package Posger

import (
	"github.com/yanyiwu/gojieba"
)

func jiebaSentenceSegmentation(article *Article, content string) {
	words := jiebaWordSegmentation(content)
	s := 0
	for e, word := range words {
		if canSegment(word) {
			article.Sentences = append(article.Sentences, Sentence{words[s: e+1], 0})
			s = e + 1
		}
	}
	article.Sentences = append(article.Sentences, Sentence{words[s:], 0})
}

func jiebaWordSegmentation(sentence string) []string {
	jieba := gojieba.NewJieba()
	defer jieba.Free()
	return jieba.Cut(sentence, true)
}

func canSegment(word string) (res bool) {
	if word == "。" {
		return true
	}
	return false
}