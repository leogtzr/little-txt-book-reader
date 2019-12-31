package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/marcusolsson/tui-go"
)

type navMode int

const (
	readingNavigationMode                    navMode = 1
	showReferencesNavigationMode             navMode = 2
	analyzeAndFilterReferencesNavigationMode navMode = 3
	gotoNavigationMode                       navMode = 4

	// Advance ...
	Advance int = 30

	// WrapMax ...
	WrapMax = 80

	// GotoWidgetIndex ...
	GotoWidgetIndex = 2

	nonRefsFileName = "non-refs.txt"

	pageSize = 10

	dbFileRequieredNumberFields = 3

	txtDBFile = "txtread"
)

var (
	from                 = 0
	to                   = Advance
	fromForReferences    = 0
	toReferences         = Advance
	gotoLine             = ""
	fileToOpen           = flag.String("file", "", "File to open")
	wrapText             = flag.Bool("wrap", false, "Wrap text")
	percentagePointStats = false
	absoluteFilePath     string
	toggleShowStatus             = true
	references                   = []string{}
	fileContent                  = []string{}
	currentNavMode       navMode = readingNavigationMode
	bannedWords                  = []string{}
	sidebar                      = tui.NewVBox()
	refsTable                    = tui.NewTable(0, 0)
	refsStatus                   = tui.NewStatusBar("_")
	pageIndex                    = 0
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
	case analyzeAndFilterReferencesNavigationMode, gotoNavigationMode:
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
	case analyzeAndFilterReferencesNavigationMode, gotoNavigationMode:
		return
	default:
		updateRangesUp()
		chunk = getChunk(&fileContent, from, to)
	}

	putText(txtArea, &chunk)
}

func getSavedStatusInformation(fileName string) string {
	return fmt.Sprintf(`%s <saved "%s">`, getStatusInformation(), fileName)
}

func getStatusInformation() string {

	if !toggleShowStatus {
		return ""
	}

	percent := float64(to) * 100.00
	percent = percent / float64(len(fileContent))
	if percentagePointStats {
		return fmt.Sprintf(".   %d of %d lines (%.3f%%) [%d lines to next percentage point]                    ",
			to,
			len(fileContent), percent, linesToChangePercentagePoint(to, len(fileContent)))
	}
	return fmt.Sprintf(".   %d of %d lines (%.3f%%)                                            ",
		to, len(fileContent), percent)

}

func init() {
	if err := createDirectories(); err != nil {
		log.Fatal(err)
	}

	// load words from file
	var err error
	bannedWords, err = loadNonRefsFile(nonRefsFileName)
	if err != nil {
		log.Fatal(err)
	}

	sidebar.Append(refsTable)
	// sidebar.Append(refsStatus)
}

func main() {

	flag.Parse()
	fileName := *fileToOpen
	if fileName == "" {
		fmt.Fprintln(os.Stderr, "error: missing file to read")
		os.Exit(1)
	}

	var err error

	absoluteFilePath, _ := filepath.Abs(fileName)
	latestFile, err := getFileNameFromLatest(absoluteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	from, to, fileName = latestFile.From, latestFile.To, latestFile.FileName

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
	// inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	txtReader := tui.NewVBox(txtAreaBox, inputCommandBox)
	txtReader.SetSizePolicy(tui.Expanding, tui.Expanding)

	chunk := getChunk(&fileContent, from, to)
	putText(txtArea, &chunk)

	root := tui.NewHBox(txtReader, sidebar)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	addUpDownKeyBindings(txtArea, ui, inputCommand)
	addGotoKeyBinding(ui, txtReader)

	// show status key binding:
	addShowStatusKeyBinding(ui, inputCommand)

	ui.SetKeybinding(newNoteKeyBindingAlternative1, func() {

		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		notesFile := getNotesDirectoryNameForFile(fileName)

		cmd := exec.Command("/usr/bin/xterm", "-fa", "Monospace", "-fs", "14", "-e", "/usr/bin/vim", "+$", notesFile)
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
		currentNavMode = readingNavigationMode
	})

	addSaveStatusKeyBinding(ui, fileName, inputCommand)
	addShowReferencesKeyBinding(ui, txtArea)

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

	addPercentageKeyBindings(ui, inputCommand)
	addcloseApplicationKeyBinding(ui, txtArea, txtReader)
	addReferencesNavigationKeyBindings(ui)
	addOnSelectedReference()

	inputCommand.SetText(getStatusInformation())

	clearScreen()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
