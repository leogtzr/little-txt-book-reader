package text

import (
	"regexp"
	"strings"
	"textreader/internal/model"
	"textreader/internal/navigation"
	"textreader/internal/words"

	"github.com/marcusolsson/tui-go"
)

func PutText(box *tui.Box, content *[]string, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	for box.Length() > 0 {
		box.Remove(0)
	}

	for i, txt := range *content {
		txt = strings.Replace(txt, " ", " ", -1)
		txt = strings.Replace(txt, "\t", "    ", -1)

		if i != state.CurrentHighlight {
			label := tui.NewLabel(txt)
			box.Append(label)
		} else {
			wordsList := words.ExtractWords(txt)
			lineBox := tui.NewHBox()
			for j, word := range wordsList {
				wordLabel := tui.NewLabel(word)
				if j == state.CurrentWord {
					wordLabel.SetStyleName("wordhighlight")
				}
				lineBox.Append(wordLabel)
				if j < len(wordsList)-1 {
					lineBox.Append(tui.NewLabel(" "))
				}
			}
			box.Append(lineBox)
		}
	}

	txtAreaScroll.ScrollToTop()
}

func GetChunk(content *[]string, from, to int) []string {
	return (*content)[from:to]
}

func MoveTextDown(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	var chunk []string
	switch state.CurrentNavMode {
	case model.ShowReferencesNavigationMode:
		navigation.UpdateRangesReferenceDown(state)
		chunk = GetChunk(&state.References, state.FromForReferences, state.ToReferences)
	case model.VocabularyNavigationMode:
		return // Disable scrolling in table mode
	// In these modes we don't want to scroll the text area
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		navigation.UpdateRangesDown(state)
		chunk = GetChunk(&state.FileContent, state.From, state.To)
		state.CurrentHighlight = 0 // Reset highlight on scroll
		state.CurrentWord = 0
	}

	PutText(txtArea, &chunk, txtAreaScroll, state)
}

func MoveTextUp(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	var chunk []string
	switch state.CurrentNavMode {
	case model.ShowReferencesNavigationMode:
		navigation.UpdateRangesReferenceUp(state)
		chunk = GetChunk(&state.References, state.FromForReferences, state.ToReferences)
	case model.VocabularyNavigationMode:
		return // Disable scrolling in table mode
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		navigation.UpdateRangesUp(state)
		chunk = GetChunk(&state.FileContent, state.From, state.To)
		state.CurrentHighlight = 0 // Reset highlight on scroll
		state.CurrentWord = 0
	}

	PutText(txtArea, &chunk, txtAreaScroll, state)
}

func FindAndRemove(s *[]string, e string) {
	for i, v := range *s {
		if v == e {
			*s = append((*s)[:i], (*s)[i+1:]...)
			break
		}
	}
}

func RemoveTrailingSpaces(s string) string {
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

// Note: Sometimes when we copy in the terminal we get multiple spaces and tabs ...
func RemoveWhiteSpaces(input string) string {
	re := regexp.MustCompile(`([ \t]){2,}`)
	return re.ReplaceAllString(input, ` `)
}

func MoveHighlightDown(box *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	state.CurrentHighlight++
	visibleLines := state.To - state.From
	if state.CurrentHighlight >= visibleLines {
		state.CurrentHighlight = visibleLines - 1
	}
	state.CurrentWord = 0
	chunk := GetChunk(&state.FileContent, state.From, state.To)
	PutText(box, &chunk, txtAreaScroll, state)
}

func MoveHighlightUp(box *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	state.CurrentHighlight--
	if state.CurrentHighlight < 0 {
		state.CurrentHighlight = 0
	}
	state.CurrentWord = 0
	chunk := GetChunk(&state.FileContent, state.From, state.To)
	PutText(box, &chunk, txtAreaScroll, state)
}

func MoveWordLeft(box *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	currentLineIndex := state.From + state.CurrentHighlight
	if currentLineIndex >= len(state.FileContent) {
		return
	}
	line := state.FileContent[currentLineIndex]
	wordsList := words.ExtractWords(line)
	if len(wordsList) == 0 {
		return
	}
	state.CurrentWord--
	if state.CurrentWord < 0 {
		state.CurrentWord = len(wordsList) - 1
	}
	chunk := GetChunk(&state.FileContent, state.From, state.To)
	PutText(box, &chunk, txtAreaScroll, state)
}

func MoveWordRight(box *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	currentLineIndex := state.From + state.CurrentHighlight
	if currentLineIndex >= len(state.FileContent) {
		return
	}
	line := state.FileContent[currentLineIndex]
	wordsList := words.ExtractWords(line)
	if len(wordsList) == 0 {
		return
	}
	state.CurrentWord++
	if state.CurrentWord >= len(wordsList) {
		state.CurrentWord = 0
	}
	chunk := GetChunk(&state.FileContent, state.From, state.To)
	PutText(box, &chunk, txtAreaScroll, state)
}
