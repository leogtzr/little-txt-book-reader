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
	"textreader/internal/progress"
	"textreader/internal/references"
	"textreader/internal/terminal"
	"textreader/internal/text"
	"textreader/internal/ui"
	"textreader/internal/utils"
	"textreader/internal/words"

	"github.com/atotto/clipboard"

	"github.com/marcusolsson/tui-go"
)

func AddWordLeftRightKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	ui.SetKeybinding("Left", AddWordLeftBinding(txtArea, inputCommand, txtAreaScroll, state))
	ui.SetKeybinding("Right", AddWordRightBinding(txtArea, inputCommand, txtAreaScroll, state))
}

func AddWordLeftBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) func() {
	return func() {
		if state.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveWordLeft(box, txtAreaScroll, state)
		input.SetText(utils.GetStatusInformation(state))
	}
}

func AddWordRightBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) func() {
	return func() {
		if state.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveWordRight(box, txtAreaScroll, state)
		input.SetText(utils.GetStatusInformation(state))
	}
}

func AddDownBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) func() {
	return func() {
		text.MoveTextDown(box, txtAreaScroll, state)
		input.SetText(utils.GetStatusInformation(state))
	}
}

func AddUpDownKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	ui.SetKeybinding(model.DownKeyBindingAlternative1, AddDownBinding(txtArea, inputCommand, txtAreaScroll, state))
	ui.SetKeybinding(model.UpKeyBindingAlternative1, AddUpBinding(txtArea, inputCommand, txtAreaScroll, state))
}

func AddHighlightUpDownKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	ui.SetKeybinding(model.DownKeyBindingAlternative2, AddHighlightDownBinding(txtArea, inputCommand, txtAreaScroll, state)) // Down arrow
	ui.SetKeybinding(model.UpKeyBindingAlternative2, AddHighlightUpBinding(txtArea, inputCommand, txtAreaScroll, state))     // Up arrow
}

func AddHighlightDownBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) func() {
	return func() {
		if state.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveHighlightDown(box, txtAreaScroll, state)
		input.SetText(utils.GetStatusInformation(state))
	}
}

func AddHighlightUpBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) func() {
	return func() {
		if state.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		text.MoveHighlightUp(box, txtAreaScroll, state)
		input.SetText(utils.GetStatusInformation(state))
	}
}

func AddShowStatusKeyBinding(ui tui.UI, inputCommand *tui.Entry, state *model.AppState) {
	ui.SetKeybinding(model.ShowStatusKeyBinding, func() {
		state.ToggleShowStatus = !state.ToggleShowStatus
		inputCommand.SetText(utils.GetStatusInformation(state))
	})
}

func AddUpBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) func() {
	return func() {
		text.MoveTextUp(box, txtAreaScroll, state)
		input.SetText(utils.GetStatusInformation(state))
	}
}

func AddSaveStatusKeyBinding(ui tui.UI, fileName string, inputCommand *tui.Entry, state *model.AppState) {
	baseFileName := filepath.Base(fileName)
	ui.SetKeybinding(model.SaveStatusKeyBindingAlternative1, func() {
		err := file.SaveStatus(fileName, state.From, state.To, state)
		if err != nil {
			inputCommand.SetText(fmt.Sprintf("Error saving status: %v", err))
			return
		}
		inputCommand.SetText(utils.GetSavedStatusInformation(baseFileName, state))
	})
}

func AddCloseApplicationKeyBinding(ui tui.UI, txtArea, txtReader *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	ui.SetKeybinding(model.CloseApplicationKeyBindingAlternative1, func() {

		switch state.CurrentNavMode {
		case model.ShowReferencesNavigationMode:
			chunk := text.GetChunk(&state.FileContent, state.From, state.To)
			text.PutText(txtArea, &chunk, txtAreaScroll, state)
			state.CurrentNavMode = model.ReadingNavigationMode
		case model.AnalyzeAndFilterReferencesNavigationMode:
			chunk := text.GetChunk(&state.FileContent, state.From, state.To)
			text.PutText(txtArea, &chunk, txtAreaScroll, state)
			state.CurrentNavMode = model.ReadingNavigationMode
			state.RefsTable.SetFocused(false)
			state.RefsTable.RemoveRows()
			state.Sidebar.SetTitle("")
			state.Sidebar.SetBorder(false)
		case model.VocabularyNavigationMode:
			chunk := text.GetChunk(&state.FileContent, state.From, state.To)
			text.PutText(txtArea, &chunk, txtAreaScroll, state)
			state.CurrentNavMode = model.ReadingNavigationMode
			state.VocabTable.SetFocused(false)
			state.VocabTable.RemoveRows()
			state.Sidebar.SetTitle("")
			state.Sidebar.SetBorder(false)
		case model.GotoNavigationMode, model.ShowTimePercentagePointsMode, model.ShowHelpMode:
			txtReader.Remove(model.GotoWidgetIndex)
			state.CurrentNavMode = model.ReadingNavigationMode
		default:
			terminal.ClearScreen()
			ui.Quit()
		}
	})
}

func AddPercentageKeyBindings(ui tui.UI, inputCommand *tui.Entry, state *model.AppState) {
	// Enable percentage tags
	ui.SetKeybinding(model.NextPercentagePointKeyBindingAlternative1, func() {
		state.PercentagePointStats = !state.PercentagePointStats
		inputCommand.SetText(utils.GetStatusInformation(state))
	})
}

func AddShowReferencesKeyBinding(ui tui.UI, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	ui.SetKeybinding(model.ShowReferencesKeyBindingAlternative1, func() {
		state.CurrentNavMode = model.ShowReferencesNavigationMode
		references.LoadReferences(state)
		chunk := text.GetChunk(&state.References, state.FromForReferences, state.ToReferences)
		text.PutText(txtArea, &chunk, txtAreaScroll, state)
	})
}

func AddShowVocabularyKeyBinding(ui tui.UI, txtReader, txtArea *tui.Box, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	ui.SetKeybinding(model.ShowVocabularyKeyBinding, func() {
		if state.CurrentNavMode == model.VocabularyNavigationMode {
			return
		}
		state.CurrentNavMode = model.VocabularyNavigationMode
		state.Sidebar.SetTitle("Vocabulary ... ")
		state.Sidebar.SetBorder(true)
		state.VocabTable.RemoveRows()
		state.PageIndex = 0
		prepareTableForVocabulary(state)
		state.VocabTable.SetFocused(true)
	})
}

func AddReferencesNavigationKeyBindings(ui tui.UI, state *model.AppState) {
	// Next References ...
	ui.SetKeybinding("Right", func() {
		if state.CurrentNavMode != model.AnalyzeAndFilterReferencesNavigationMode {
			return
		}
		if state.PageIndex >= len(state.References) {
			return
		}
		state.PageIndex += model.PageSize
		prepareTableForReferences(state)
	})

	// Previous References ...
	ui.SetKeybinding("Left", func() {
		if state.CurrentNavMode != model.AnalyzeAndFilterReferencesNavigationMode {
			return
		}
		if state.PageIndex < model.PageSize {
			return
		}
		state.PageIndex -= model.PageSize
		prepareTableForReferences(state)
	})
}

func AddSaveQuoteKeyBindings(ui tui.UI, fileName string, txtArea *tui.Box, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
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

		chunk := text.GetChunk(&state.FileContent, state.From, state.To)
		text.PutText(txtArea, &chunk, txtAreaScroll, state)
		inputCommand.SetText(utils.GetStatusInformation(state))
	})
}

func prepareTableForReferences(state *model.AppState) {
	state.RefsTable.RemoveRows()
	paginatedReferences := utils.Paginate(state.References, state.PageIndex, model.PageSize)
	if len(paginatedReferences) == 0 {
		state.RefsTable.AppendRow(tui.NewLabel("No references found"))
	} else {
		for _, ref := range paginatedReferences {
			state.RefsTable.AppendRow(tui.NewLabel(ref))
		}
	}
	state.RefsTable.SetSelected(0)
}

func AddOnSelectedReference(state *model.AppState) {
	state.RefsTable.OnItemActivated(func(t *tui.Table) {
		itemIndexToRemove := t.Selected()
		itemToAddToNonRefs := state.References[state.PageIndex+itemIndexToRemove]
		text.FindAndRemove(&state.References, itemToAddToNonRefs)
		prepareTableForReferences(state)

		if !words.Contains(state.BannedWords, itemToAddToNonRefs) {
			file.AppendLineToFile(model.NonRefsFileName, itemToAddToNonRefs, "")
		}
	})
}

func AddGotoKeyBinding(tuiUI tui.UI, txtReader *tui.Box, state *model.AppState) {
	tuiUI.SetKeybinding(model.GotoKeyBindingAlternative1, func() {
		ui.AddGotoWidget(txtReader, state)
	})
}

func AddCloseGotoBinding(ui tui.UI, inputCommand *tui.Entry, txtReader, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
	ui.SetKeybinding(model.CloseGotoKeyBindingAlternative1, func() {
		// Go To the specified line
		inputCommand.SetText(utils.GetStatusInformation(state))

		gotoLineNumber := progress.GetNumberLineGoto(state.GotoLine)
		gotoLineNumberDigits, err := strconv.ParseInt(gotoLineNumber, 10, 64)
		if err != nil {
			return
		}
		if int(gotoLineNumberDigits) < (len(state.FileContent) - state.Advance) {
			state.From = int(gotoLineNumberDigits)
			state.To = state.From + state.Advance
			chunk := text.GetChunk(&state.FileContent, state.From, state.To)
			text.PutText(txtArea, &chunk, txtAreaScroll, state)
			inputCommand.SetText(utils.GetStatusInformation(state))
		}
		txtReader.Remove(model.GotoWidgetIndex)
		inputCommand.SetText(utils.GetStatusInformation(state))
		state.CurrentNavMode = model.ReadingNavigationMode
	})
}

func AddNewNoteKeyBinding(ui tui.UI, txtArea *tui.Box, inputCommand *tui.Entry, fileName string, txtAreaScroll *tui.ScrollArea, state *model.AppState) {
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
		chunk := text.GetChunk(&state.FileContent, state.From, state.To)
		text.PutText(txtArea, &chunk, txtAreaScroll, state)
		inputCommand.SetText(utils.GetStatusInformation(state))
	})
}

func AddAnalyzeAndFilterReferencesKeyBinding(ui tui.UI, state *model.AppState) {
	ui.SetKeybinding(model.AnalyzeAndFilterReferencesKeyBinding, func() {
		state.CurrentNavMode = model.AnalyzeAndFilterReferencesNavigationMode
		state.Sidebar.SetTitle("References ... ")
		state.Sidebar.SetBorder(true)
		state.RefsTable.SetColumnStretch(0, 0)
		references.LoadReferences(state)

		state.RefsTable.RemoveRows()
		prepareTableForReferences(state)
		state.RefsTable.SetFocused(true)
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

func AddShowMinutesTakenToReachPercentagePointKeyBinding(ui tui.UI, txtReader *tui.Box, state *model.AppState) {
	ui.SetKeybinding(model.ShowMinutesTakenToReachPercentagePointKeyBinding, func() {

		// Check if we are already in that mode ...
		if state.CurrentNavMode == model.ShowTimePercentagePointsMode {
			return
		}

		state.CurrentNavMode = model.ShowTimePercentagePointsMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range state.MinutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		for _, v := range percentages {
			duration := state.MinutesToReachNextPercentagePoint[v]
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

func AddShowHelpKeyBinding(ui tui.UI, txtReader *tui.Box, state *model.AppState) {
	ui.SetKeybinding(model.ShowHelpKeyBinding, func() {

		// Check if we are already in that mode ...
		if state.CurrentNavMode == model.ShowHelpMode {
			return
		}

		state.CurrentNavMode = model.ShowHelpMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range state.MinutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		addKeyBindingDescription(fmt.Sprintf("%10s -> Go Down / Go Up",
			model.DownKeyBindingAlternative1+"/"+model.UpKeyBindingAlternative1), &strs)
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
		addKeyBindingDescription(fmt.Sprintf("%10s -> Add a Quote, gets the text from the clipboard.", model.SaveQuoteKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows Time Stats for each percentage point.", model.ShowMinutesTakenToReachPercentagePointKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows this Dialog", model.ShowHelpKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens RAE Web site search with the clipboard content", model.OpenRAEWebSiteKeyBinging), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens GoodReads Web site with the clipboard content", model.OpenGoodReadsWebSiteKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Highlight Down", model.DownKeyBindingAlternative2), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Highlight Up", model.UpKeyBindingAlternative2), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Word Left / Word Right", "Left/Right"), &strs)
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

func AddCopyWordKeyBinding(ui tui.UI, inputCommand *tui.Entry, state *model.AppState) {
	ui.SetKeybinding("c", func() {
		if state.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		currentLineIndex := state.From + state.CurrentHighlight
		if currentLineIndex >= len(state.FileContent) {
			inputCommand.SetText("No word to copy")
			return
		}
		line := state.FileContent[currentLineIndex]
		wordsList := words.ExtractWords(line)
		if len(wordsList) == 0 || state.CurrentWord >= len(wordsList) {
			inputCommand.SetText("No word to copy")
			return
		}
		word := wordsList[state.CurrentWord]
		word = words.SanitizeWord(word)
		err := clipboard.WriteAll(word)
		if err != nil {
			inputCommand.SetText(fmt.Sprintf("Error copying word: %v", err))
			return
		}
		inputCommand.SetText(fmt.Sprintf("Copied '%s' to clipboard", word))
	})
}

func AddSaveVocabularyKeyBinding(ui tui.UI, fileName string, inputCommand *tui.Entry, state *model.AppState) {
	ui.SetKeybinding(model.SaveVocabularyKeyBinding, func() {
		if state.CurrentNavMode != model.ReadingNavigationMode {
			return
		}
		currentLineIndex := state.From + state.CurrentHighlight
		if currentLineIndex >= len(state.FileContent) {
			inputCommand.SetText("No word to save")
			return
		}
		line := state.FileContent[currentLineIndex]
		wordsList := words.ExtractWords(line)
		if len(wordsList) == 0 || state.CurrentWord >= len(wordsList) {
			inputCommand.SetText("No word to save")
			return
		}
		word := wordsList[state.CurrentWord]
		word = words.SanitizeWord(word)
		if words.Contains(state.Vocabulary, word) {
			inputCommand.SetText(fmt.Sprintf("Word '%s' already in vocabulary", word))
			return
		}
		state.Vocabulary = append(state.Vocabulary, word)
		err := file.SaveStatus(fileName, state.From, state.To, state)
		if err != nil {
			inputCommand.SetText(fmt.Sprintf("Error saving vocabulary: %v", err))
			return
		}
		inputCommand.SetText(fmt.Sprintf("Saved '%s' to vocabulary", word))
	})
}

func prepareTableForVocabulary(state *model.AppState) {
	state.VocabTable.RemoveRows()
	paginatedVocabulary := utils.Paginate(state.Vocabulary, state.PageIndex, model.PageSize)
	if len(paginatedVocabulary) == 0 {
		state.VocabTable.AppendRow(tui.NewLabel("No vocabulary words saved"))
	} else {
		for _, word := range paginatedVocabulary {
			state.VocabTable.AppendRow(tui.NewLabel(word))
		}
	}
	state.VocabTable.SetSelected(0)
}

func AddVocabularyNavigationKeyBindings(ui tui.UI, state *model.AppState) {
	ui.SetKeybinding("Right", func() {
		if state.CurrentNavMode != model.VocabularyNavigationMode {
			return
		}
		if state.PageIndex >= len(state.Vocabulary)-model.PageSize {
			return
		}
		state.PageIndex += model.PageSize
		prepareTableForVocabulary(state)
	})
	ui.SetKeybinding("Left", func() {
		if state.CurrentNavMode != model.VocabularyNavigationMode {
			return
		}
		if state.PageIndex < model.PageSize {
			return
		}
		state.PageIndex -= model.PageSize
		prepareTableForVocabulary(state)
	})
	ui.SetKeybinding(model.UpKeyBindingAlternative2, func() {
		if state.CurrentNavMode != model.VocabularyNavigationMode {
			return
		}
		selected := state.VocabTable.Selected()
		if selected > 0 {
			state.VocabTable.SetSelected(selected - 1)
		}
	})
	ui.SetKeybinding(model.DownKeyBindingAlternative2, func() {
		if state.CurrentNavMode != model.VocabularyNavigationMode {
			return
		}
		selected := state.VocabTable.Selected()
		paginatedVocabulary := utils.Paginate(state.Vocabulary, state.PageIndex, model.PageSize)
		if selected < len(paginatedVocabulary)-1 {
			state.VocabTable.SetSelected(selected + 1)
		}
	})
	ui.SetKeybinding("k", func() {
		if state.CurrentNavMode != model.VocabularyNavigationMode {
			return
		}
		selected := state.VocabTable.Selected()
		if selected > 0 {
			state.VocabTable.SetSelected(selected - 1)
		}
	})
	ui.SetKeybinding("j", func() {
		if state.CurrentNavMode != model.VocabularyNavigationMode {
			return
		}
		selected := state.VocabTable.Selected()
		paginatedVocabulary := utils.Paginate(state.Vocabulary, state.PageIndex, model.PageSize)
		if selected < len(paginatedVocabulary)-1 {
			state.VocabTable.SetSelected(selected + 1)
		}
	})
}

func AddOnSelectedVocabulary(state *model.AppState) {
	state.VocabTable.OnItemActivated(func(t *tui.Table) {
		itemIndexToRemove := t.Selected()
		itemToAddToNonRefs := state.Vocabulary[state.PageIndex+itemIndexToRemove]
		text.FindAndRemove(&state.Vocabulary, itemToAddToNonRefs)
		prepareTableForVocabulary(state)
		file.SaveStatus(state.FileToOpen, state.From, state.To, state)
		if !words.Contains(state.BannedWords, itemToAddToNonRefs) {
			file.AppendLineToFile(model.NonRefsFileName, itemToAddToNonRefs, "")
		}
	})
}
