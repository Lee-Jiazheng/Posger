package Posger

import (
	"fmt"
	"testing"
)

func TestTikaServer(t *testing.T) {
	fmt.Print(ExtractDocumentContent("static/articles/BMC.pdf"))
}
