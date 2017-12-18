package EnSegment

import (
	"testing"
	"fmt"
)

func TestEnglishText(t *testing.T) {
	res := CutAll("I want to fuck you!are you, mother fucker, ok?")
	for _, i := range res {
		fmt.Print(i + "\n")
	}
	fmt.Print(res)
}