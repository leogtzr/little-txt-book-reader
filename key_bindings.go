package main

import (
	"path/filepath"

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
			appendLineToFile(nonRefsFileName, itemToAddToNonRefs)
		}
	})
}

func addGotoKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(gotoKeyBindingAlternative1, func() {
		addGotoWidget(txtReader)
	})
}
