package main

import (
	"strings"
	"textreader/internal/model"

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
		updateRangesReferenceDown()
		chunk = GetChunk(&model.References, model.FromForReferences, model.ToReferences)
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		updateRangesDown()
		chunk = GetChunk(&model.FileContent, model.From, model.To)
	}

	PutText(txtArea, &chunk, txtAreaScroll)
}

func MoveTextUp(txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	chunk := []string{}
	switch model.CurrentNavMode {
	case model.ShowReferencesNavigationMode:
		updateRangesReferenceUp()
		chunk = GetChunk(&model.References, model.FromForReferences, model.ToReferences)
	case model.AnalyzeAndFilterReferencesNavigationMode, model.GotoNavigationMode:
		return
	default:
		updateRangesUp()
		chunk = GetChunk(&model.FileContent, model.From, model.To)
	}

	PutText(txtArea, &chunk, txtAreaScroll)
}

func findAndRemove(s *[]string, e string) {
	for i, v := range *s {
		if v == e {
			*s = append((*s)[:i], (*s)[i+1:]...)
			break
		}
	}
}
