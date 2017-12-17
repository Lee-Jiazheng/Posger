package textrank

var symbols = []string{".", "!", "?", "..."}

func isSegmentSymbol(word string) bool {
	for _, w := range symbols {
		if w == word {
			return true
		}
	}
	return false
}
