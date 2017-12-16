package tika

import (
	"testing"
	"fmt"
)

func TestPdfDirectory(t *testing.T) {
	content, err := GetPdfContent("../static/articles/test.pdf")
	if err != nil {
		t.Error(err)
	}
	fmt.Print(content)
}
