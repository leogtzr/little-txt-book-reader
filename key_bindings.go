package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"textreader/internal/file"
	"textreader/internal/model"
	"textreader/internal/utils"

	"github.com/atotto/clipboard"

	"github.com/marcusolsson/tui-go"
)

func addDownBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		MoveTextDown(box, txtAreaScroll)
		input.SetText(getStatusInformation())
	}
}

func addUpDownKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.DownKeyBindingAlternative1, addDownBinding(txtArea, inputCommand, txtAreaScroll))
	ui.SetKeybinding(model.DownKeyBindingAlternative2, addDownBinding(txtArea, inputCommand, txtAreaScroll))

	ui.SetKeybinding(model.UpKeyBindingAlternative1, addUpBinding(txtArea, inputCommand, txtAreaScroll))
	ui.SetKeybinding(model.UpKeyBindingAlternative2, addUpBinding(txtArea, inputCommand, txtAreaScroll))
}

func addShowStatusKeyBinding(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(model.ShowStatusKeyBinding, func() {
		model.ToggleShowStatus = !model.ToggleShowStatus
		inputCommand.SetText(getStatusInformation())
	})
}

func addUpBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		MoveTextUp(box, txtAreaScroll)
		input.SetText(getStatusInformation())
	}
}

func addSaveStatusKeyBinding(ui tui.UI, fileName string, inputCommand *tui.Entry) {
	baseFileName := filepath.Base(fileName)
	ui.SetKeybinding(model.SaveStatusKeyBindingAlternative1, func() {
		file.SaveStatus(fileName, model.From, model.To)
		inputCommand.SetText(getSavedStatusInformation(baseFileName))
	})
}

func addCloseApplicationKeyBinding(ui tui.UI, txtArea, txtReader *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.CloseApplicationKeyBindingAlternative1, func() {

		switch model.CurrentNavMode {
		case model.ShowReferencesNavigationMode:
			chunk := GetChunk(&model.FileContent, model.From, model.To)
			PutText(txtArea, &chunk, txtAreaScroll)
			model.CurrentNavMode = model.ReadingNavigationMode
		case model.AnalyzeAndFilterReferencesNavigationMode:
			chunk := GetChunk(&model.FileContent, model.From, model.To)
			PutText(txtArea, &chunk, txtAreaScroll)
			model.CurrentNavMode = model.ReadingNavigationMode
			model.RefsTable.SetFocused(false)
		case model.GotoNavigationMode, model.ShowTimePercentagePointsMode, model.ShowHelpMode:
			txtReader.Remove(model.GotoWidgetIndex)
			model.CurrentNavMode = model.ReadingNavigationMode
		default:
			utils.ClearScreen()
			ui.Quit()
		}
	})
}

func addPercentageKeyBindings(ui tui.UI, inputCommand *tui.Entry) {
	// Enable percentage tags
	ui.SetKeybinding(model.NextPercentagePointKeyBindingAlternative1, func() {
		model.PercentagePointStats = !model.PercentagePointStats
		inputCommand.SetText(getStatusInformation())
	})
}

func addShowReferencesKeyBinding(ui tui.UI, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.ShowReferencesKeyBindingAlternative1, func() {
		model.CurrentNavMode = model.ShowReferencesNavigationMode
		loadReferences()
		chunk := GetChunk(&model.References, model.FromForReferences, model.ToReferences)
		PutText(txtArea, &chunk, txtAreaScroll)
	})
}

func addReferencesNavigationKeyBindings(ui tui.UI) {
	// Next References ...
	ui.SetKeybinding("Right", func() {
		if model.PageIndex >= len(model.References) {
			return
		}
		model.PageIndex += model.PageSize
		prepareTableForReferences()
	})

	// Previous References ...
	ui.SetKeybinding("Left", func() {
		if model.PageIndex < model.PageSize {
			return
		}
		model.PageIndex -= model.PageSize
		prepareTableForReferences()
	})
}

func addSaveQuoteKeyBindings(ui tui.UI, fileName string, txtArea *tui.Box, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.SaveQuoteKeyBindingAlternative1, func() {
		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		quotesFile := file.GetDirectoryNameForFile("quotes", fileName)

		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}

		clipBoardText = removeTrailingSpaces(clipBoardText)
		clipBoardText = removeWhiteSpaces(clipBoardText)
		file.AppendLineToFile(quotesFile, clipBoardText, "\n__________")

		cmd := openOSEditor(runtime.GOOS, quotesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr

		// txtReader.SetBorder(true)
		chunk := GetChunk(&model.FileContent, model.From, model.To)
		PutText(txtArea, &chunk, txtAreaScroll)
		inputCommand.SetText(getStatusInformation())
	})
}

func prepareTableForReferences() {
	model.RefsTable.RemoveRows()
	references := utils.Paginate(model.References, model.PageIndex, model.PageSize)
	for _, ref := range references {
		model.RefsTable.AppendRow(tui.NewLabel(ref))
	}
	model.RefsTable.SetSelected(0)
}

func addOnSelectedReference() {
	model.RefsTable.OnItemActivated(func(tui *tui.Table) {

		itemIndexToRemove := tui.Selected()
		itemToAddToNonRefs := model.References[model.PageIndex+itemIndexToRemove]
		// References = remove(References, itemIndexToRemove)
		findAndRemove(&model.References, itemToAddToNonRefs)
		prepareTableForReferences()

		if !contains(model.BannedWords, itemToAddToNonRefs) {
			file.AppendLineToFile(model.NonRefsFileName, itemToAddToNonRefs, "")
		}
	})
}

func addGotoKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(model.GotoKeyBindingAlternative1, func() {
		utils.AddGotoWidget(txtReader)
	})
}

func addCloseGotoBinding(ui tui.UI, inputCommand *tui.Entry, txtReader, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.CloseGotoKeyBindingAlternative1, func() {
		// Go To the specified line
		inputCommand.SetText(getStatusInformation())

		gotoLineNumber := utils.GetNumberLineGoto(model.GotoLine)
		gotoLineNumberDigits, err := strconv.ParseInt(gotoLineNumber, 10, 64)
		if err != nil {
			return
		}
		if int(gotoLineNumberDigits) < (len(model.FileContent) - model.Advance) {
			model.From = int(gotoLineNumberDigits)
			model.To = model.From + model.Advance
			chunk := GetChunk(&model.FileContent, model.From, model.To)
			PutText(txtArea, &chunk, txtAreaScroll)
			inputCommand.SetText(getStatusInformation())
		}
		txtReader.Remove(model.GotoWidgetIndex)
		inputCommand.SetText(getStatusInformation())
		model.CurrentNavMode = model.ReadingNavigationMode
	})
}

func addNewNoteKeyBinding(ui tui.UI, txtArea *tui.Box, inputCommand *tui.Entry, fileName string, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.NewNoteKeyBindingAlternative1, func() {

		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		notesFile := file.GetDirectoryNameForFile("notes", fileName)

		cmd := openOSEditor(runtime.GOOS, notesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr
		// txtReader.SetBorder(true)
		chunk := GetChunk(&model.FileContent, model.From, model.To)
		PutText(txtArea, &chunk, txtAreaScroll)
		inputCommand.SetText(getStatusInformation())
	})
}

func addAnalyzeAndFilterReferencesKeyBinding(ui tui.UI) {
	ui.SetKeybinding(model.AnalyzeAndFilterReferencesKeyBinding, func() {
		model.CurrentNavMode = model.AnalyzeAndFilterReferencesNavigationMode
		model.Sidebar.SetTitle("References ... ")
		model.Sidebar.SetBorder(true)
		model.RefsTable.SetColumnStretch(0, 0)
		loadReferences()

		model.RefsTable.RemoveRows()
		prepareTableForReferences()
		model.RefsTable.SetFocused(true)
	})
}

func browserOpenURLCommand(osName, url string) *exec.Cmd {
	switch osName {
	case "linux":
		return exec.Command("xdg-open", url)
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		return exec.Command("open", url)
	default:
		return nil
	}
}

func addOpenRAEWebSite(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(model.OpenRAEWebSiteKeyBinging, func() {
		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}
		if len(strings.TrimSpace(clipBoardText)) == 0 {
			return
		}
		raeURL := fmt.Sprintf("https://dle.rae.es/%s", clipBoardText)
		if err = browserOpenURLCommand(runtime.GOOS, raeURL).Start(); err != nil {
			inputCommand.SetText(err.Error())
			return
		}
	})
}

func browserOpenGoodReadsURLCommand(osName, url string) *exec.Cmd {
	switch osName {
	case "linux":
		return exec.Command("xdg-open", url)
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		return exec.Command("open", url)
	default:
		return nil
	}
}

func addOpenGoodReadsWebSite(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(model.OpenGoodReadsWebSiteKeyBinding, func() {
		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}
		if len(strings.TrimSpace(clipBoardText)) == 0 {
			return
		}
		goodreadsURL := fmt.Sprintf(`https://www.goodreads.com/search?q=%s`, url.QueryEscape(clipBoardText))
		if err = browserOpenGoodReadsURLCommand(runtime.GOOS, goodreadsURL).Start(); err != nil {
			inputCommand.SetText(err.Error())
			return
		}
	})
}

func addShowMinutesTakenToReachPercentagePointKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(model.ShowMinutesTakenToReachPercentagePointKeyBinding, func() {

		// Check if we are already in that mode ...
		if model.CurrentNavMode == model.ShowTimePercentagePointsMode {
			return
		}

		model.CurrentNavMode = model.ShowTimePercentagePointsMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range model.MinutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		for _, v := range percentages {
			duration := model.MinutesToReachNextPercentagePoint[v]
			strs = append(strs, fmt.Sprintf("%d%% took you %.1f minutes", v, duration.Minutes()))
		}

		l.AddItems(strs...)
		s := tui.NewScrollArea(l)
		s.SetFocused(true)

		txtReader.Append(s)

		ui.SetKeybinding("Alt+Up", func() { s.Scroll(0, -1) })
		ui.SetKeybinding("Alt+Down", func() { s.Scroll(0, 1) })
	})
}

func addShowHelpKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(model.ShowHelpKeyBinding, func() {

		// Check if we are already in that mode ...
		if model.CurrentNavMode == model.ShowHelpMode {
			return
		}

		model.CurrentNavMode = model.ShowHelpMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range model.MinutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		addKeyBindingDescription(fmt.Sprintf("%10s -> Go Down", model.DownKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Go Up", model.UpKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Go To", model.GotoKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> New Note", model.NewNoteKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Show Status", model.ShowStatusKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the Goto Dialog", model.CloseGotoKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Save Progress", model.SaveStatusKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows Next Percentage Point Stats", model.NextPercentagePointKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows the References Dialog", model.ShowReferencesKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the References Dialog", model.CloseReferencesWindowKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the program", model.CloseApplicationKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Analyze and filter References", model.AnalyzeAndFilterReferencesKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Add a Quote, gets the text From the clipboard.", model.SaveQuoteKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows Time Stats for each percentage point.", model.ShowMinutesTakenToReachPercentagePointKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows this Dialog", model.ShowHelpKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens RAE Web site with search From the clipboard.", model.OpenRAEWebSiteKeyBinging), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens GoodReads Web site with search From the clipboard.", model.OpenGoodReadsWebSiteKeyBinding), &strs)

		l.AddItems(strs...)
		s := tui.NewScrollArea(l)
		s.SetFocused(true)

		txtReader.Append(s)

		ui.SetKeybinding("Alt+Up", func() { s.Scroll(0, -1) })
		ui.SetKeybinding("Alt+Down", func() { s.Scroll(0, 1) })
	})
}

func addKeyBindingDescription(desc string, keyBindings *[]string) {
	*keyBindings = append(*keyBindings, desc)
}

// Sometimes when we copy in the terminal we get multiple spaces and tabs ...
func removeWhiteSpaces(input string) string {
	re := regexp.MustCompile(`( |\t){2,}`)
	return re.ReplaceAllString(input, ` `)
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
