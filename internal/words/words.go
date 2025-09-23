package words

import (
	"strings"
	"unicode"
)

func ExtractWords(line string) []string {
	words := strings.Split(line, " ")
	return words
}

func SanitizeWord(line string) string {
	line = strings.ReplaceAll(line, ".", "")
	line = strings.ReplaceAll(line, ",", "")
	line = strings.ReplaceAll(line, "\"", "")
	line = strings.ReplaceAll(line, ")", "")
	line = strings.ReplaceAll(line, "(", "")
	line = strings.ReplaceAll(line, ":", "")
	line = strings.ReplaceAll(line, ";", "")
	return line
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
