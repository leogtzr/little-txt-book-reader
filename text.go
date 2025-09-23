package main

import (
	"regexp"
	"strings"

	"github.com/marcusolsson/tui-go"
)

func putText(box *tui.Box, content *[]string, txtAreaScroll *tui.ScrollArea) {
	for box.Length() > 0 {
		box.Remove(0)
	}

	for _, txt := range *content {
		txt = strings.Replace(txt, " ", " ", -1)
		txt = strings.Replace(txt, "\t", "    ", -1)
		box.Append(tui.NewVBox(
			tui.NewLabel(txt),
		))
	}

	txtAreaScroll.ScrollToTop()
}

func getChunk(content *[]string, from, to int) []string {
	return (*content)[from:to]
}

func downText(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	chunk := []string{}
	switch CurrentNavMode {
	case ShowReferencesNavigationMode:
		updateRangesReferenceDown()
		chunk = getChunk(&References, FromForReferences, ToReferences)
	case AnalyzeAndFilterReferencesNavigationMode, GotoNavigationMode:
		return
	default:
		updateRangesDown()
		chunk = getChunk(&FileContent, From, To)
	}

	putText(txtArea, &chunk, txtAreaScroll)
}

func upText(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	chunk := []string{}
	switch CurrentNavMode {
	case ShowReferencesNavigationMode:
		updateRangesReferenceUp()
		chunk = getChunk(&References, FromForReferences, ToReferences)
	case AnalyzeAndFilterReferencesNavigationMode, GotoNavigationMode:
		return
	default:
		updateRangesUp()
		chunk = getChunk(&FileContent, From, To)
	}

	putText(txtArea, &chunk, txtAreaScroll)
}

func longestLineLength(text string) int {
	if len(text) == 0 {
		return 0
	}

	lines := strings.Split(text, "\n")
	longest := len(lines[0])
	for _, line := range lines {
		if len(line) > longest {
			longest = len(line)
		}
	}
	return longest
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

func removeTrailingSpaces(s string) string {
	lines := strings.Split(s, "\n")
	var sb strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		sb.WriteString(strings.TrimSpace(line))
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func removeWhiteSpaces(input string) string {
	re := regexp.MustCompile(`( |\t){2,}`)
	return re.ReplaceAllString(input, ` `)
}

func listsAreEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
