package references

import (
	"bufio"
	"os"
	"strings"
	"textreader/internal/model"
	"textreader/internal/terminal"
	"textreader/internal/words"
)

// TODO: fix this, we can pass state only.
func ExtractReferencesFromFileContent(fileContent *[]string, state *model.AppState) []string {
	refs := make([]string, 0)
	i := 0
	for _, lineInFile := range *fileContent {
		i++
		lineInFile = strings.TrimSpace(lineInFile)
		if lineInFile == "" || len(lineInFile) == 0 {
			continue
		}

		r := ExtractReferences(lineInFile)
		if len(r) > 0 {
			for _, ref := range r {
				refs = append(refs, ref)
			}
		}
	}

	uniqueReferences := removeDuplicates(refs)

	referencesNoBannedWords := make([]string, 0)
	for _, word := range uniqueReferences {
		if !words.Contains(state.BannedWords, word) {
			referencesNoBannedWords = append(referencesNoBannedWords, word)
		}
	}

	return referencesNoBannedWords
}

func LoadReferences(state *model.AppState) {
	if len(state.References) == 0 {
		state.References = ExtractReferencesFromFileContent(&state.FileContent, state)
		state.ToReferences = terminal.CalculateTerminalHeight()
	}
}

func LoadNonRefsFile(path string) ([]string, error) {
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

func ExtractReferences(line string) []string {
	referenceWords := words.ExtractWords(strings.TrimSpace(line))
	if len(referenceWords) == 0 {
		return []string{}
	}

	refs := make([]string, 0)
	f := false
	bag := make([]string, 0)

	for _, word := range referenceWords {
		if word == "" {
			continue
		}
		f = words.IsTitle(word)
		if f {
			bag = append(bag, word)
			if strings.Contains(word, ",") || strings.Contains(word, ".") {
				f = false
			}
		}
		if !f {
			if len(bag) > 0 {
				refs = append(refs, concat(bag))
				bag = nil
			}
		}
	}

	if len(bag) > 0 {
		refs = append(refs, concat(bag))
	}

	return refs
}

// TODO: change the name of the params/variables.
func concat(bag []string) string {
	s := ""
	for _, w := range bag {
		s += w
		s += " "
	}
	return words.SanitizeWord(strings.TrimSpace(s))
}
