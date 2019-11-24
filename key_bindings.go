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
	ui.SetKeybinding(saveStatusKeyBindingAlternative1, func() {
		absoluteFilePath, _ := filepath.Abs(fileName)
		saveStatus(absoluteFilePath, from, to)

		inputCommand.SetText(getSavedStatusInformation())
	})
}

func addcloseApplicationKeyBinding(ui tui.UI, txtArea *tui.Box) {
	ui.SetKeybinding(closeApplicationKeyBindingAlternative1, func() {
		switch currentNavMode {
		case showReferencesNavigationMode:
			chunk := getChunk(&fileContent, from, to)
			putText(txtArea, &chunk)
			currentNavMode = readingNavigationMode
		default:
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
		if len(references) == 0 {
			references = extractReferencesFromFileContent(&fileContent)
		}
		chunk := getChunk(&references, fromForReferences, toReferences)
		putText(txtArea, &chunk)
	})
}
