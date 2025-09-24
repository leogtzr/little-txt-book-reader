package keybindings

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"textreader/internal/file"
	"textreader/internal/model"
	"textreader/internal/references"
	"textreader/internal/text"
	"textreader/internal/utils"
	"textreader/internal/words"

	"github.com/atotto/clipboard"

	"github.com/marcusolsson/tui-go"
)

// Add word navigation bindings
func AddWordLeftRightKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding("Left", AddWordLeftBinding(txtArea, inputCommand, txtAreaScroll))
	ui.SetKeybinding("Right", AddWordRightBinding(txtArea, inputCommand, txtAreaScroll))
}

func AddWordLeftBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		if model.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveWordLeft(box, txtAreaScroll)
		input.SetText(utils.GetStatusInformation())
	}
}

func AddWordRightBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		if model.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveWordRight(box, txtAreaScroll)
		input.SetText(utils.GetStatusInformation())
	}
}

func AddDownBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		text.MoveTextDown(box, txtAreaScroll)
		input.SetText(utils.GetStatusInformation())
	}
}

func AddUpDownKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.DownKeyBindingAlternative1, AddDownBinding(txtArea, inputCommand, txtAreaScroll))
	//ui.SetKeybinding(model.DownKeyBindingAlternative2, AddDownBinding(txtArea, inputCommand, txtAreaScroll))

	ui.SetKeybinding(model.UpKeyBindingAlternative1, AddUpBinding(txtArea, inputCommand, txtAreaScroll))
	//ui.SetKeybinding(model.UpKeyBindingAlternative2, AddUpBinding(txtArea, inputCommand, txtAreaScroll))
}

func AddHighlightUpDownKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.DownKeyBindingAlternative2, AddHighlightDownBinding(txtArea, inputCommand, txtAreaScroll)) // Down arrow
	ui.SetKeybinding(model.UpKeyBindingAlternative2, AddHighlightUpBinding(txtArea, inputCommand, txtAreaScroll))     // Up arrow
}

func AddHighlightDownBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		if model.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveHighlightDown(box, txtAreaScroll)
		input.SetText(utils.GetStatusInformation())
	}
}

func AddHighlightUpBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		if model.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveHighlightUp(box, txtAreaScroll)
		input.SetText(utils.GetStatusInformation())
	}
}

func AddShowStatusKeyBinding(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(model.ShowStatusKeyBinding, func() {
		model.ToggleShowStatus = !model.ToggleShowStatus
		inputCommand.SetText(utils.GetStatusInformation())
	})
}

func AddUpBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		text.MoveTextUp(box, txtAreaScroll)
		input.SetText(utils.GetStatusInformation())
	}
}

func AddSaveStatusKeyBinding(ui tui.UI, fileName string, inputCommand *tui.Entry) {
	baseFileName := filepath.Base(fileName)
	ui.SetKeybinding(model.SaveStatusKeyBindingAlternative1, func() {
		file.SaveStatus(fileName, model.From, model.To)
		inputCommand.SetText(utils.GetSavedStatusInformation(baseFileName))
	})
}

func AddCloseApplicationKeyBinding(ui tui.UI, txtArea, txtReader *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.CloseApplicationKeyBindingAlternative1, func() {

		switch model.CurrentNavMode {
		case model.ShowReferencesNavigationMode:
			chunk := text.GetChunk(&model.FileContent, model.From, model.To)
			text.PutText(txtArea, &chunk, txtAreaScroll)
			model.CurrentNavMode = model.ReadingNavigationMode
		case model.AnalyzeAndFilterReferencesNavigationMode:
			chunk := text.GetChunk(&model.FileContent, model.From, model.To)
			text.PutText(txtArea, &chunk, txtAreaScroll)
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

func AddPercentageKeyBindings(ui tui.UI, inputCommand *tui.Entry) {
	// Enable percentage tags
	ui.SetKeybinding(model.NextPercentagePointKeyBindingAlternative1, func() {
		model.PercentagePointStats = !model.PercentagePointStats
		inputCommand.SetText(utils.GetStatusInformation())
	})
}

func AddShowReferencesKeyBinding(ui tui.UI, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.ShowReferencesKeyBindingAlternative1, func() {
		model.CurrentNavMode = model.ShowReferencesNavigationMode
		references.LoadReferences()
		chunk := text.GetChunk(&model.References, model.FromForReferences, model.ToReferences)
		text.PutText(txtArea, &chunk, txtAreaScroll)
	})
}

func AddReferencesNavigationKeyBindings(ui tui.UI) {
	// Next References ...
	//ui.SetKeybinding("Right", func() {
	//	if model.PageIndex >= len(model.References) {
	//		return
	//	}
	//	model.PageIndex += model.PageSize
	//	prepareTableForReferences()
	//})
	//
	//// Previous References ...
	//ui.SetKeybinding("Left", func() {
	//	if model.PageIndex < model.PageSize {
	//		return
	//	}
	//	model.PageIndex -= model.PageSize
	//	prepareTableForReferences()
	//})
	// Next References ...
	ui.SetKeybinding("Right", func() {
		if model.CurrentNavMode != model.AnalyzeAndFilterReferencesNavigationMode {
			return
		}
		if model.PageIndex >= len(model.References) {
			return
		}
		model.PageIndex += model.PageSize
		prepareTableForReferences()
	})

	// Previous References ...
	ui.SetKeybinding("Left", func() {
		if model.CurrentNavMode != model.AnalyzeAndFilterReferencesNavigationMode {
			return
		}
		if model.PageIndex < model.PageSize {
			return
		}
		model.PageIndex -= model.PageSize
		prepareTableForReferences()
	})
}

func AddSaveQuoteKeyBindings(ui tui.UI, fileName string, txtArea *tui.Box, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.SaveQuoteKeyBindingAlternative1, func() {
		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		quotesFile := file.GetDirectoryNameForFile("quotes", fileName)

		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}

		clipBoardText = text.RemoveTrailingSpaces(clipBoardText)
		clipBoardText = text.RemoveWhiteSpaces(clipBoardText)
		file.AppendLineToFile(quotesFile, clipBoardText, "\n__________")

		cmd := utils.OpenOSEditor(runtime.GOOS, quotesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr

		// txtReader.SetBorder(true)
		chunk := text.GetChunk(&model.FileContent, model.From, model.To)
		text.PutText(txtArea, &chunk, txtAreaScroll)
		inputCommand.SetText(utils.GetStatusInformation())
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

func AddOnSelectedReference() {
	model.RefsTable.OnItemActivated(func(tui *tui.Table) {

		itemIndexToRemove := tui.Selected()
		itemToAddToNonRefs := model.References[model.PageIndex+itemIndexToRemove]
		// References = remove(References, itemIndexToRemove)
		text.FindAndRemove(&model.References, itemToAddToNonRefs)
		prepareTableForReferences()

		if !words.Contains(model.BannedWords, itemToAddToNonRefs) {
			file.AppendLineToFile(model.NonRefsFileName, itemToAddToNonRefs, "")
		}
	})
}

func AddGotoKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(model.GotoKeyBindingAlternative1, func() {
		utils.AddGotoWidget(txtReader)
	})
}

func AddCloseGotoBinding(ui tui.UI, inputCommand *tui.Entry, txtReader, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.CloseGotoKeyBindingAlternative1, func() {
		// Go To the specified line
		inputCommand.SetText(utils.GetStatusInformation())

		gotoLineNumber := utils.GetNumberLineGoto(model.GotoLine)
		gotoLineNumberDigits, err := strconv.ParseInt(gotoLineNumber, 10, 64)
		if err != nil {
			return
		}
		if int(gotoLineNumberDigits) < (len(model.FileContent) - model.Advance) {
			model.From = int(gotoLineNumberDigits)
			model.To = model.From + model.Advance
			chunk := text.GetChunk(&model.FileContent, model.From, model.To)
			text.PutText(txtArea, &chunk, txtAreaScroll)
			inputCommand.SetText(utils.GetStatusInformation())
		}
		txtReader.Remove(model.GotoWidgetIndex)
		inputCommand.SetText(utils.GetStatusInformation())
		model.CurrentNavMode = model.ReadingNavigationMode
	})
}

func AddNewNoteKeyBinding(ui tui.UI, txtArea *tui.Box, inputCommand *tui.Entry, fileName string, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(model.NewNoteKeyBindingAlternative1, func() {

		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		notesFile := file.GetDirectoryNameForFile("notes", fileName)

		cmd := utils.OpenOSEditor(runtime.GOOS, notesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr
		// txtReader.SetBorder(true)
		chunk := text.GetChunk(&model.FileContent, model.From, model.To)
		text.PutText(txtArea, &chunk, txtAreaScroll)
		inputCommand.SetText(utils.GetStatusInformation())
	})
}

func AddAnalyzeAndFilterReferencesKeyBinding(ui tui.UI) {
	ui.SetKeybinding(model.AnalyzeAndFilterReferencesKeyBinding, func() {
		model.CurrentNavMode = model.AnalyzeAndFilterReferencesNavigationMode
		model.Sidebar.SetTitle("References ... ")
		model.Sidebar.SetBorder(true)
		model.RefsTable.SetColumnStretch(0, 0)
		references.LoadReferences()

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

func AddOpenRAEWebSite(ui tui.UI, inputCommand *tui.Entry) {
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

func AddOpenGoodReadsWebSite(ui tui.UI, inputCommand *tui.Entry) {
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

func AddShowMinutesTakenToReachPercentagePointKeyBinding(ui tui.UI, txtReader *tui.Box) {
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

func AddShowHelpKeyBinding(ui tui.UI, txtReader *tui.Box) {
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
		addKeyBindingDescription(fmt.Sprintf("%10s -> Highlight Down", model.DownKeyBindingAlternative2), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Highlight Up", model.UpKeyBindingAlternative2), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Word Left", "Left"), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Word Right", "Right"), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Copy Word to Clipboard", "c"), &strs)

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

func AddCopyWordKeyBinding(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding("c", func() {
		if model.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		currentLineIndex := model.From + model.CurrentHighlight
		if currentLineIndex >= len(model.FileContent) {
			inputCommand.SetText("No word to copy")
			return
		}
		line := model.FileContent[currentLineIndex]
		wordsList := words.ExtractWords(line)
		if len(wordsList) == 0 || model.CurrentWord >= len(wordsList) {
			inputCommand.SetText("No word to copy")
			return
		}
		word := wordsList[model.CurrentWord]
		err := clipboard.WriteAll(word)
		if err != nil {
			inputCommand.SetText(fmt.Sprintf("Error copying word: %v", err))
			return
		}
		inputCommand.SetText(fmt.Sprintf("Copied '%s' to clipboard", word))
	})
}
