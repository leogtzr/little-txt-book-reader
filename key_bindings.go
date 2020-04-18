package main

import (
	"os"
	"path/filepath"
	"runtime"

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

func addcloseApplicationKeyBinding(ui tui.UI, txtArea, txtReader *tui.Box) {
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
		case gotoNavigationMode:
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
		appendLineToFile(quotesFile, clipBoardText, "\n__________")

		cmd := openOSEditor(runtime.GOOS, quotesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr

		txtReader.SetBorder(true)

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
			appendLineToFile("", nonRefsFileName, itemToAddToNonRefs)
		}
	})
}

func addGotoKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(gotoKeyBindingAlternative1, func() {
		addGotoWidget(txtReader)
	})
}
