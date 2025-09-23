package text

import (
	"regexp"
	"strings"
	"textreader/internal/model"
	"textreader/internal/navigation"

	"github.com/marcusolsson/tui-go"
)

func PutText(box *tui.Box, content *[]string, txtAreaScroll *tui.ScrollArea) {
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

func GetChunk(content *[]string, from, to int) []string {
	return (*content)[from:to]
}

func MoveTextDown(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	chunk := []string{}
	switch model.CurrentNavMode {
	case model.ShowReferencesNavigationMode:
		navigation.UpdateRangesReferenceDown()
		chunk = GetChunk(&model.References, model.FromForReferences, model.ToReferences)
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		navigation.UpdateRangesDown()
		chunk = GetChunk(&model.FileContent, model.From, model.To)
	}

	PutText(txtArea, &chunk, txtAreaScroll)
}

func MoveTextUp(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	chunk := []string{}
	switch model.CurrentNavMode {
	case model.ShowReferencesNavigationMode:
		navigation.UpdateRangesReferenceUp()
		chunk = GetChunk(&model.References, model.FromForReferences, model.ToReferences)
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		navigation.UpdateRangesUp()
		chunk = GetChunk(&model.FileContent, model.From, model.To)
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

// Sometimes when we copy in the terminal we get multiple spaces and tabs ...
func RemoveWhiteSpaces(input string) string {
	re := regexp.MustCompile(`( |\t){2,}`)
	return re.ReplaceAllString(input, ` `)
}
