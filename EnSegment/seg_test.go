package EnSegment

import (
	"testing"
	"fmt"
)

func TestEnglishText(t *testing.T) {
	res := CutAll("I want to fuck you!are you ok?")
	for _, i := range res {
		fmt.Print(i + ",")
	}
	fmt.Print(res)
}