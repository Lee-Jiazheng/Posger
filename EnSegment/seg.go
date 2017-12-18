package EnSegment

var cutSymbolsEn = []rune{' ', '?', '!', '.', '\n', ',', ';', '\'', '"', '/', '(', ')'}
//var cutSymbolsCn = []rune{'，', '？', '！', '。', '；', '”', '“',}
var cutSymbolsCn = []rune{}
var CutSymbols = append(cutSymbolsEn, cutSymbolsCn...)

func CutAll(content string) []string {
	var start, end int
	strings := make([]string, 0)
	for i, s := range content {
		if canSegment(s) {
			end = i
			if start == end { start = i + 1; continue}
			strings = append(strings, content[start:end])
			if s != ' ' && s != '\n'{
				strings = append(strings, content[i:i+1])
				end += 1
			}
			start = i + 1
		}
	}
	return strings
}

func canSegment(char rune) bool {
	for _, ch := range CutSymbols {
		if ch == char {
			return true
		}
	}
	return false
}