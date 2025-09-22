package main

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

	"github.com/atotto/clipboard"

	"github.com/marcusolsson/tui-go"
)

func addUpDownKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(downKeyBindingAlternative1, addDownBinding(txtArea, inputCommand))
	ui.SetKeybinding(downKeyBindingAlternative2, addDownBinding(txtArea, inputCommand))

	ui.SetKeybinding(upKeyBindingAlternative1, addUpBinding(txtArea, inputCommand))
	ui.SetKeybinding(upKeyBindingAlternative2, addUpBinding(txtArea, inputCommand))
}

func addShowStatusKeyBinding(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(showStatusKeyBinding, func() {
		toggleShowStatus = !toggleShowStatus
		inputCommand.SetText(getStatusInformation())
	})
}

func addSaveStatusKeyBinding(ui tui.UI, fileName string, inputCommand *tui.Entry) {
	baseFileName := filepath.Base(fileName)
	ui.SetKeybinding(saveStatusKeyBindingAlternative1, func() {
		saveStatus(fileName, from, to)
		inputCommand.SetText(getSavedStatusInformation(baseFileName))
	})
}

func addCloseApplicationKeyBinding(ui tui.UI, txtArea, txtReader *tui.Box) {
	ui.SetKeybinding(closeApplicationKeyBindingAlternative1, func() {

		switch currentNavMode {
		case showReferencesNavigationMode:
			chunk := getChunk(&fileContent, from, to)
			putText(txtArea, &chunk)
			currentNavMode = readingNavigationMode
		case analyzeAndFilterReferencesNavigationMode:
			chunk := getChunk(&fileContent, from, to)
			putText(txtArea, &chunk)
			currentNavMode = readingNavigationMode
			refsTable.SetFocused(false)
		case gotoNavigationMode, showTimePercentagePointsMode, showHelpMode:
			txtReader.Remove(GotoWidgetIndex)
			currentNavMode = readingNavigationMode
		default:
			clearScreen()
			ui.Quit()
		}
	})
}

func addPercentageKeyBindings(ui tui.UI, inputCommand *tui.Entry) {
	// Enable percentage tags
	ui.SetKeybinding(nextPercentagePointKeyBindingAlternative1, func() {
		percentagePointStats = !percentagePointStats
		inputCommand.SetText(getStatusInformation())
	})
}

func addShowReferencesKeyBinding(ui tui.UI, txtArea *tui.Box) {
	ui.SetKeybinding(showReferencesKeyBindingAlternative1, func() {
		currentNavMode = showReferencesNavigationMode
		loadReferences()
		chunk := getChunk(&references, fromForReferences, toReferences)
		putText(txtArea, &chunk)
	})
}

func addReferencesNavigationKeyBindings(ui tui.UI) {
	// Next references ...
	ui.SetKeybinding("Right", func() {
		if pageIndex >= len(references) {
			return
		}
		pageIndex += pageSize
		prepareTableForReferences()
	})

	// Previous references ...
	ui.SetKeybinding("Left", func() {
		if pageIndex < pageSize {
			return
		}
		pageIndex -= pageSize
		prepareTableForReferences()
	})
}

func addSaveQuoteKeyBindings(ui tui.UI, fileName string, txtArea *tui.Box, inputCommand *tui.Entry) {
	ui.SetKeybinding(saveQuoteKeyBindingAlternative1, func() {
		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		quotesFile := getDirectoryNameForFile("quotes", fileName)

		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}

		clipBoardText = removeTrailingSpaces(clipBoardText)
		clipBoardText = removeWhiteSpaces(clipBoardText)
		appendLineToFile(quotesFile, clipBoardText, "\n__________")

		cmd := openOSEditor(runtime.GOOS, quotesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr

		// txtReader.SetBorder(true)
		chunk := getChunk(&fileContent, from, to)
		putText(txtArea, &chunk)
		inputCommand.SetText(getStatusInformation())
	})
}

func prepareTableForReferences() {
	refsTable.RemoveRows()
	references := paginate(references, pageIndex, pageSize)
	for _, ref := range references {
		refsTable.AppendRow(tui.NewLabel(ref))
	}
	refsTable.SetSelected(0)
}

func addOnSelectedReference() {
	refsTable.OnItemActivated(func(tui *tui.Table) {

		itemIndexToRemove := tui.Selected()
		itemToAddToNonRefs := references[pageIndex+itemIndexToRemove]
		// references = remove(references, itemIndexToRemove)
		findAndRemove(&references, itemToAddToNonRefs)
		prepareTableForReferences()

		if !contains(bannedWords, itemToAddToNonRefs) {
			appendLineToFile(nonRefsFileName, itemToAddToNonRefs, "")
		}
	})
}

func addGotoKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(gotoKeyBindingAlternative1, func() {
		addGotoWidget(txtReader)
	})
}

func addCloseGotoBinding(ui tui.UI, inputCommand *tui.Entry, txtReader, txtArea *tui.Box) {
	ui.SetKeybinding(closeGotoKeyBindingAlternative1, func() {
		// Go to the specified line
		inputCommand.SetText(getStatusInformation())

		gotoLineNumber := getNumberLineGoto(gotoLine)
		gotoLineNumberDigits, err := strconv.ParseInt(gotoLineNumber, 10, 64)
		if err != nil {
			return
		}
		if int(gotoLineNumberDigits) < (len(fileContent) - Advance) {
			from = int(gotoLineNumberDigits)
			to = from + Advance
			chunk := getChunk(&fileContent, from, to)
			putText(txtArea, &chunk)
			inputCommand.SetText(getStatusInformation())
		}
		txtReader.Remove(GotoWidgetIndex)
		inputCommand.SetText(getStatusInformation())
		currentNavMode = readingNavigationMode
	})
}

func addNewNoteKeyBinding(ui tui.UI, txtArea *tui.Box, inputCommand *tui.Entry, fileName string) {
	ui.SetKeybinding(newNoteKeyBindingAlternative1, func() {

		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		notesFile := getDirectoryNameForFile("notes", fileName)

		cmd := openOSEditor(runtime.GOOS, notesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr
		// txtReader.SetBorder(true)
		chunk := getChunk(&fileContent, from, to)
		putText(txtArea, &chunk)
		inputCommand.SetText(getStatusInformation())
	})
}

func addAnalyzeAndFilterReferencesKeyBinding(ui tui.UI) {
	ui.SetKeybinding(analyzeAndFilterReferencesKeyBinding, func() {
		currentNavMode = analyzeAndFilterReferencesNavigationMode
		sidebar.SetTitle("References ... ")
		sidebar.SetBorder(true)
		refsTable.SetColumnStretch(0, 0)
		loadReferences()

		refsTable.RemoveRows()
		prepareTableForReferences()
		refsTable.SetFocused(true)
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
	ui.SetKeybinding(openRAEWebSiteKeyBinging, func() {
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
	ui.SetKeybinding(openGoodReadsWebSiteKeyBinding, func() {
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
	ui.SetKeybinding(showMinutesTakenToReachPercentagePointKeyBinding, func() {

		// Check if we are already in that mode ...
		if currentNavMode == showTimePercentagePointsMode {
			return
		}

		currentNavMode = showTimePercentagePointsMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range minutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		for _, v := range percentages {
			duration := minutesToReachNextPercentagePoint[v]
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
	ui.SetKeybinding(showHelpKeyBinding, func() {

		// Check if we are already in that mode ...
		if currentNavMode == showHelpMode {
			return
		}

		currentNavMode = showHelpMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range minutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		addKeyBindingDescription(fmt.Sprintf("%10s -> Go Down", downKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Go Up", upKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Go To", gotoKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> New Note", newNoteKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Show Status", showStatusKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the Goto Dialog", closeGotoKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Save Progress", saveStatusKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows Next Percentage Point Stats", nextPercentagePointKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows the References Dialog", showReferencesKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the References Dialog", closeReferencesWindowKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the program", closeApplicationKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Analyze and filter references", analyzeAndFilterReferencesKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Add a Quote, gets the text from the clipboard.", saveQuoteKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows Time Stats for each percentage point.", showMinutesTakenToReachPercentagePointKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows this Dialog", showHelpKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens RAE Web site with search from the clipboard.", openRAEWebSiteKeyBinging), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens GoodReads Web site with search from the clipboard.", openGoodReadsWebSiteKeyBinding), &strs)

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
