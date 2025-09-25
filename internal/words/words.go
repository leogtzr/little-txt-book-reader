package words

import (
	"regexp"
	"strings"
	"unicode"
)

func ExtractWords(line string) []string {
	return strings.Fields(line)
}

func SanitizeWord(line string) string {
	//line = strings.ReplaceAll(line, ".", "")
	//line = strings.ReplaceAll(line, ",", "")
	//line = strings.ReplaceAll(line, "\"", "")
	//line = strings.ReplaceAll(line, ")", "")
	//line = strings.ReplaceAll(line, "(", "")
	//line = strings.ReplaceAll(line, ":", "")
	//line = strings.ReplaceAll(line, ";", "")
	//return line
	re := regexp.MustCompile(`[.,"()?:;\\]`)
	return re.ReplaceAllString(line, "")
}

func IsTitle(word string) bool {
	chars := []rune(word)
	return unicode.IsUpper(chars[0])
}

func Contains(words []string, word string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}
