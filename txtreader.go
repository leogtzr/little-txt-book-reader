package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/marcusolsson/tui-go"
)

const (
	// Advance ...
	Advance int = 30

	// WrapMax ...
	WrapMax = 80

	// GotoWidgetIndex ...
	GotoWidgetIndex = 2
)

var (
	from                 = 0
	to                   = Advance
	fromForReferences    = 0
	toReferences         = Advance
	gotoLine             = ""
	fileToOpen           = flag.String("file", "", "File to open")
	wrapText             = flag.Bool("wrap", false, "Wrap text")
	openLatestFile       = flag.Bool("latest", false, "Open the latest text file")
	percentagePointStats = false
	absoluteFilePath     string
	toggleShowStatus     = true
	showReferencesMode   = false
	references           = []string{}
	fileContent          = []string{}
)

// LatestFile ...
type LatestFile struct {
	FileName string
	From     int
	To       int
}

func updateRangesUp() {
	if from <= 0 {
		return
	}

	if from > 0 {
		from--
	}

	to--
}

func updateRangesReferenceUp() {
	if fromForReferences <= 0 {
		return
	}

	if fromForReferences > 0 {
		fromForReferences--
	}

	toReferences--
}

func updateRangesDown() {
	if from < len(fileContent) {
		from++
	}

	if to >= len(fileContent) {
		return
	}

	if to < len(fileContent) {
		to++
	}
}

func updateRangesReferenceDown() {
	if fromForReferences < len(references) {
		fromForReferences++
	}

	if toReferences >= len(references) {
		return
	}

	if toReferences < len(references) {
		toReferences++
	}
}

func downText(txtArea *tui.Box) {
	if showReferencesMode {
		updateRangesReferenceDown()
	} else {
		updateRangesDown()
	}

	chunk := []string{}
	if showReferencesMode {
		chunk = getChunk(&references, fromForReferences, toReferences)
	} else {
		chunk = getChunk(&fileContent, from, to)
	}

	putText(txtArea, &chunk)
}

func upText(txtArea *tui.Box) {
	if showReferencesMode {
		updateRangesReferenceUp()
	} else {
		updateRangesUp()
	}

	chunk := []string{}
	if showReferencesMode {
		chunk = getChunk(&references, fromForReferences, toReferences)
	} else {
		chunk = getChunk(&fileContent, from, to)
	}

	putText(txtArea, &chunk)
}

func getSavedStatusInformation() string {
	return fmt.Sprintf("%s <saved>", getStatusInformation())
}

func getStatusInformation() string {

	if !toggleShowStatus {
		return ""
	}

	percent := float64(to) * 100.00
	percent = percent / float64(len(fileContent))
	if percentagePointStats {
		return fmt.Sprintf(".   %d of %d lines (%.3f%%) [%d lines to next percentage point]                                                            ",
			to,
			len(fileContent), percent, linesToChangePercentagePoint(to, len(fileContent)))
	}
	return fmt.Sprintf(".   %d of %d lines (%.3f%%)                                                            ",
		to, len(fileContent), percent)

}

func main() {

	flag.Parse()
	fileName := *fileToOpen
	if fileName == "" && !*openLatestFile {
		fmt.Fprintln(os.Stderr, "error: missing file to read")
		os.Exit(1)
	}

	if fileName != "" && *openLatestFile {
		fmt.Fprintln(os.Stderr, "error: conflicting options")
		os.Exit(1)
	}

	var err error

	if *openLatestFile {
		latestFile, err := getFileNameFromLatest()
		from = latestFile.From
		to = latestFile.To

		fileName = latestFile.FileName
		if err != nil {
			log.Fatal(err)
		}
	}

	fileContent, err = readLines(fileName)
	check(err)

	txtArea := tui.NewVBox()
	txtAreaScroll := tui.NewScrollArea(txtArea)
	txtAreaScroll.SetAutoscrollToBottom(true)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	txtAreaBox.SetBorder(true)

	inputCommand := tui.NewEntry()
	inputCommand.SetFocused(true)
	inputCommand.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputCommand.SetEchoMode(tui.EchoModeNormal)

	inputCommandBox := tui.NewHBox(inputCommand)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	txtReader := tui.NewVBox(txtAreaBox, inputCommandBox)
	txtReader.SetSizePolicy(tui.Expanding, tui.Expanding)

	chunk := getChunk(&fileContent, from, to)
	putText(txtArea, &chunk)

	root := tui.NewHBox(txtReader)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	// down ...
	ui.SetKeybinding(downKeyBindingAlternative1, addDownBinding(txtArea, inputCommand))
	ui.SetKeybinding(downKeyBindingAlternative2, addDownBinding(txtArea, inputCommand))

	// Up ...
	ui.SetKeybinding(upKeyBindingAlternative1, addUpBinding(txtArea, inputCommand))
	ui.SetKeybinding(upKeyBindingAlternative2, addUpBinding(txtArea, inputCommand))

	// go to:
	ui.SetKeybinding(gotoKeyBindingAlterntive1, func() {
		addGotoWidget(txtReader)
	})

	// show status key binding:
	ui.SetKeybinding(showStatusKeyBinding, func() {
		toggleShowStatus = !toggleShowStatus
		inputCommand.SetText(getStatusInformation())
	})

	noteBox := tui.NewTextEdit()
	noteBox.SetText("")

	referencesBox := tui.NewTextEdit()
	referencesBox.SetText("")

	// new note key binding:
	ui.SetKeybinding(newNoteKeyBindingAlternative1, func() {
		prepareNewNoteBox(noteBox)
		inputCommand.SetFocused(false)
		inputCommand.SetText("> > > > > Creating note ... ")
		txtReader.SetFocused(false)
		txtArea.SetFocused(false)
		txtAreaScroll.SetFocused(false)

		txtReader.Insert(0, noteBox)
	})

	ui.SetKeybinding(saveNoteKeyBindingAlternative1, func() {
		if !noteBox.IsFocused() {
			return
		}
		saveNote(fileName, noteBox)
		txtReader.Remove(0)
	})

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
			putText(txtArea, &chunk)
			inputCommand.SetText(getStatusInformation())
		}
		txtReader.Remove(GotoWidgetIndex)
		inputCommand.SetText(getStatusInformation())
	})

	ui.SetKeybinding(saveStatusKeyBindingAlternative1, func() {
		absoluteFilePath, _ := filepath.Abs(fileName)
		saveStatus(absoluteFilePath, from, to)

		inputCommand.SetText(getSavedStatusInformation())
	})

	// Enable percentage tags
	ui.SetKeybinding(nextPercentagePointKeyBindingAlternative1, func() {
		percentagePointStats = !percentagePointStats
		inputCommand.SetText(getStatusInformation())
	})

	ui.SetKeybinding(showReferencesKeyBindingAlternative1, func() {
		showReferencesMode = true
		if len(references) == 0 {
			references = extractReferencesFromFileContent(&fileContent)
		}
		chunk := getChunk(&references, fromForReferences, toReferences)
		putText(txtArea, &chunk)
	})

	ui.SetKeybinding(closeApplicationKeyBindingAlternative1, func() {
		if showReferencesMode {
			chunk := getChunk(&fileContent, from, to)
			putText(txtArea, &chunk)
			showReferencesMode = false
		} else {
			ui.Quit()
		}
	})

	inputCommand.SetText(getStatusInformation())

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
