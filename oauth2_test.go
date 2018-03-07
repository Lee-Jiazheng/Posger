package Posger

import (
	"fmt"
	"testing"
)

func TestRedirectUrl(t *testing.T) {
	fmt.Println(oauth2Url("github"))
}