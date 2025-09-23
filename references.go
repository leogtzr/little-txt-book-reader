package main

import (
	"bufio"
	"os"
	"strings"
)

func extractReferencesFromFileContent(fileContent *[]string) []string {
	refs := make([]string, 0)
	i := 0
	for _, lineInFile := range *fileContent {
		i++
		lineInFile = strings.TrimSpace(lineInFile)
		if lineInFile == "" || len(lineInFile) == 0 {
			continue
		}

		r := extractReferences(lineInFile)
		if len(r) > 0 {
			for _, ref := range r {
				refs = append(refs, ref)
			}
		}
	}

	uniqueReferences := removeDuplicates(refs)

	referencesNoBannedWords := make([]string, 0)
	for _, word := range uniqueReferences {
		if !contains(BannedWords, word) {
			referencesNoBannedWords = append(referencesNoBannedWords, word)
		}
	}

	return referencesNoBannedWords
}

func loadReferences() {
	if len(References) == 0 {
		References = extractReferencesFromFileContent(&FileContent)
		ToReferences = calculateAdvanceHeight()
	}
}

func findAndRemove(s *[]string, e string) {
	for i, v := range *s {
		if v == e {
			*s = append((*s)[:i], (*s)[i+1:]...)
			break
		}
	}
}

func loadNonRefsFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
