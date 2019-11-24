package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/marcusolsson/tui-go"
)

type navMode int

const (
	readingNavigationMode                    navMode = 1
	showReferencesNavigationMode             navMode = 2
	analyzeAndFilterReferencesNavigationMode navMode = 3
	// Advance ...
	Advance int = 30

	// WrapMax ...
	WrapMax = 80

	// GotoWidgetIndex ...
	GotoWidgetIndex = 2

	nonRefsFileName = "non-refs.txt"
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
	toggleShowStatus             = true
	references                   = []string{}
	fileContent                  = []string{}
	currentNavMode       navMode = readingNavigationMode
	bannedWords                  = []string{}
	sidebar                      = tui.NewVBox()
	refsTable                    = tui.NewTable(0, 0)
	refsStatus                   = tui.NewStatusBar("__________")
	// refsTableScroll              = tui.NewScrollArea(sidebar)
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
	chunk := []string{}
	switch currentNavMode {
	case showReferencesNavigationMode:
		updateRangesReferenceDown()
		chunk = getChunk(&references, fromForReferences, toReferences)
	case analyzeAndFilterReferencesNavigationMode:
		return
	default:
		updateRangesDown()
		chunk = getChunk(&fileContent, from, to)
	}

	putText(txtArea, &chunk)
}

func upText(txtArea *tui.Box) {
	chunk := []string{}
	switch currentNavMode {
	case showReferencesNavigationMode:
		updateRangesReferenceUp()
		chunk = getChunk(&references, fromForReferences, toReferences)
	case analyzeAndFilterReferencesNavigationMode:
		return
	default:
		updateRangesUp()
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

// load words from file
func init() {
	var err error
	bannedWords, err = loadNonRefsFile(nonRefsFileName)
	if err != nil {
		log.Fatal(err)
	}

	sidebar.Append(refsTable)
	sidebar.Append(refsStatus)
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

	// <<<<<<<
	//
	//sidebar.SetBorder(true)
	// >>>>>>>

	root := tui.NewHBox(txtReader, sidebar)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	addUpDownKeyBindings(txtArea, ui, inputCommand)

	// go to:
	ui.SetKeybinding(gotoKeyBindingAlterntive1, func() {
		addGotoWidget(txtReader)
	})

	// show status key binding:
	addShowStatusKeyBinding(ui, inputCommand)

	noteBox := tui.NewTextEdit()
	noteBox.SetText("")

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

	addSaveStatusKeyBinding(ui, fileName, inputCommand)
	addShowReferencesKeyBinding(ui, txtArea)

	ui.SetKeybinding(analyzeAndFilterReferencesKeyBinding, func() {
		currentNavMode = analyzeAndFilterReferencesNavigationMode
		sidebar.SetTitle("References ... ")
		sidebar.SetBorder(true)

		refsTable.SetColumnStretch(0, 0)

		loadReferences()
		// TODO: Need a clever way of getting this shit ...
		// TODO: a status bar ...
		for _, ref := range references[0:10] {
			refsTable.AppendRow(
				tui.NewLabel(ref),
			)
		}

		refsTable.SetFocused(true)
		// refsTableScroll.SetFocused(true)
	})

	addPercentageKeyBindings(ui, inputCommand)
	addcloseApplicationKeyBinding(ui, txtArea)

	refsTable.OnItemActivated(func(table *tui.Table) {
		// TODO: remove item ...
		// fmt.Println(table.Selected())
	})

	inputCommand.SetText(getStatusInformation())

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
