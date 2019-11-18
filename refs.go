package main

import (
	"strings"
)

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func extractReferences(line string) []string {
	words := extractWords(strings.TrimSpace(line))
	if len(words) == 0 {
		return []string{}
	}

	refs := make([]string, 0)
	f := false
	bag := make([]string, 0)

	for _, word := range words {
		f = isTitle(word)
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

func concat(bag []string) string {
	s := ""
	for _, w := range bag {
		s += w
		s += " "
	}
	return sanitizeWord(strings.TrimSpace(s))
}
