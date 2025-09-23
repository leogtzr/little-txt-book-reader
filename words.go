package main

import (
	"strings"
	"textreader/internal/model"
	"unicode"
)

func extractWords(line string) []string {
	words := strings.Split(line, " ")
	return words
}

func sanitizeWord(line string) string {
	line = strings.ReplaceAll(line, ".", "")
	line = strings.ReplaceAll(line, ",", "")
	line = strings.ReplaceAll(line, "\"", "")
	line = strings.ReplaceAll(line, ")", "")
	line = strings.ReplaceAll(line, "(", "")
	line = strings.ReplaceAll(line, ":", "")
	line = strings.ReplaceAll(line, ";", "")
	return line
}

func sanitizeFileName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, " ", "")
	return fileName
}

func isTitle(word string) bool {
	chars := []rune(word)
	return unicode.IsUpper(chars[0])
}

func contains(words []string, word string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

func shouldIgnoreWord(word string) bool {
	return contains(model.BannedWords, word)
}
