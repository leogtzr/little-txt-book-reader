package main

import (
	"fmt"
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
	// down ...
	ui.SetKeybinding(downKeyBindingAlternative1, addDownBinding(txtArea, inputCommand))
	ui.SetKeybinding(downKeyBindingAlternative2, addDownBinding(txtArea, inputCommand))

	// Up ...
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

func addSaveQuoteKeyBindings(ui tui.UI, fileName string, txtArea, txtReader *tui.Box, inputCommand *tui.Entry) {
	ui.SetKeybinding(saveQuoteKeyBindingAlternative1, func() {
		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		quotesFile := getDirectoryNameForFile("quotes", fileName)

		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}

		clipBoardText = removeTrailingSpaces(clipBoardText)
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

func addOpenRAEWebSite(ui tui.UI, inputCommand *tui.Entry) {
	// openRAEWebSiteKeyBinging
	ui.SetKeybinding(openRAEWebSiteKeyBinging, func() {
		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}
		if len(strings.TrimSpace(clipBoardText)) == 0 {
			return
		}
		url := fmt.Sprintf("https://dle.rae.es/%s", clipBoardText)
		if err = exec.Command("xdg-open", url).Start(); err != nil {
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

/*
showHelpKeyBinding                               = "Alt+h"
*/

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

		strs = append(strs, "Alt+j    -> Go Down")
		strs = append(strs, "Down     -> Go Down")
		strs = append(strs, "Alt+k    -> Go Up")
		strs = append(strs, "Up       -> Go Up")
		strs = append(strs, "Go To    -> Go To")
		strs = append(strs, "Alt+n    -> New Note")
		strs = append(strs, "Alt+.    -> Show Status")
		strs = append(strs, "r        -> Closes the Goto Dialog")
		strs = append(strs, "s        -> Save Progress")
		strs = append(strs, "Alt+p    -> Shows Next Percentage Point Stats")
		strs = append(strs, "Alt+r    -> Shows the References Dialog")
		strs = append(strs, "Alt+q    -> Closes the References Dialog")
		strs = append(strs, "Esc      -> Closes the program")
		strs = append(strs, "Alt+b    -> Analyze and filter references")
		strs = append(strs, "Alt+q    -> Add a Quote, gets the text from the clipboard.")
		strs = append(strs, "Alt+m    -> Shows Time Stats for each percentage point.")
		strs = append(strs, "Alt+y    -> Shows this Dialog")

		l.AddItems(strs...)
		s := tui.NewScrollArea(l)
		s.SetFocused(true)

		txtReader.Append(s)

		ui.SetKeybinding("Alt+Up", func() { s.Scroll(0, -1) })
		ui.SetKeybinding("Alt+Down", func() { s.Scroll(0, 1) })
	})
}
