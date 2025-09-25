package text

import (
	"regexp"
	"strings"
	"textreader/internal/model"
	"textreader/internal/navigation"
	"textreader/internal/words"

	"github.com/marcusolsson/tui-go"
)

func PutText(box *tui.Box, content *[]string, txtAreaScroll *tui.ScrollArea) {
	for box.Length() > 0 {
		box.Remove(0)
	}

	for i, txt := range *content {
		txt = strings.Replace(txt, " ", " ", -1)
		txt = strings.Replace(txt, "\t", "    ", -1)

		if i != model.CurrentHighlight {
			label := tui.NewLabel(txt)
			box.Append(label)
		} else {
			wordsList := words.ExtractWords(txt)
			lineBox := tui.NewHBox()
			for j, word := range wordsList {
				wordLabel := tui.NewLabel(word)
				if j == model.CurrentWord {
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

func MoveTextDown(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	var chunk []string
	switch model.CurrentNavMode {
	case model.ShowReferencesNavigationMode:
		navigation.UpdateRangesReferenceDown()
		chunk = GetChunk(&model.References, model.FromForReferences, model.ToReferences)
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		navigation.UpdateRangesDown()
		chunk = GetChunk(&model.FileContent, model.From, model.To)
		model.CurrentHighlight = 0 // Reset highlight on scroll
		model.CurrentWord = 0
	}

	PutText(txtArea, &chunk, txtAreaScroll)
}

func MoveTextUp(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	var chunk []string
	switch model.CurrentNavMode {
	case model.ShowReferencesNavigationMode:
		navigation.UpdateRangesReferenceUp()
		chunk = GetChunk(&model.References, model.FromForReferences, model.ToReferences)
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		navigation.UpdateRangesUp()
		chunk = GetChunk(&model.FileContent, model.From, model.To)
		model.CurrentHighlight = 0 // Reset highlight on scroll
		model.CurrentWord = 0
	}

	PutText(txtArea, &chunk, txtAreaScroll)
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

func MoveHighlightDown(box *tui.Box, txtAreaScroll *tui.ScrollArea) {
	model.CurrentHighlight++
	visibleLines := model.To - model.From
	if model.CurrentHighlight >= visibleLines {
		model.CurrentHighlight = visibleLines - 1
	}
	model.CurrentWord = 0
	chunk := GetChunk(&model.FileContent, model.From, model.To)
	PutText(box, &chunk, txtAreaScroll)
}

func MoveHighlightUp(box *tui.Box, txtAreaScroll *tui.ScrollArea) {
	model.CurrentHighlight--
	if model.CurrentHighlight < 0 {
		model.CurrentHighlight = 0
	}
	model.CurrentWord = 0
	chunk := GetChunk(&model.FileContent, model.From, model.To)
	PutText(box, &chunk, txtAreaScroll)
}

func MoveWordLeft(box *tui.Box, txtAreaScroll *tui.ScrollArea) {
	currentLineIndex := model.From + model.CurrentHighlight
	if currentLineIndex >= len(model.FileContent) {
		return
	}
	line := model.FileContent[currentLineIndex]
	wordsList := words.ExtractWords(line)
	if len(wordsList) == 0 {
		return
	}
	model.CurrentWord--
	if model.CurrentWord < 0 {
		model.CurrentWord = len(wordsList) - 1
	}
	chunk := GetChunk(&model.FileContent, model.From, model.To)
	PutText(box, &chunk, txtAreaScroll)
}

func MoveWordRight(box *tui.Box, txtAreaScroll *tui.ScrollArea) {
	currentLineIndex := model.From + model.CurrentHighlight
	if currentLineIndex >= len(model.FileContent) {
		return
	}
	line := model.FileContent[currentLineIndex]
	wordsList := words.ExtractWords(line)
	if len(wordsList) == 0 {
		return
	}
	model.CurrentWord++
	if model.CurrentWord >= len(wordsList) {
		model.CurrentWord = 0
	}
	chunk := GetChunk(&model.FileContent, model.From, model.To)
	PutText(box, &chunk, txtAreaScroll)
}
