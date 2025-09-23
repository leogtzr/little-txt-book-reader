package main

import (
	"bufio"
	"os"
	"strings"
	"textreader/internal/model"
	"textreader/internal/utils"
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
		if !contains(model.BannedWords, word) {
			referencesNoBannedWords = append(referencesNoBannedWords, word)
		}
	}

	return referencesNoBannedWords
}

func loadReferences() {
	if len(model.References) == 0 {
		model.References = extractReferencesFromFileContent(&model.FileContent)
		model.ToReferences = utils.CalculateTerminalHeight()
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

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		// Do not add duplicate.
		if !encountered[elements[v]] {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append To result slice.
			result = append(result, elements[v])
		}
	}
	return result
}
