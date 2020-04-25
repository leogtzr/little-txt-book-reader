package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/marcusolsson/tui-go"
	"github.com/muesli/termenv"
)

var (
	from                 = 0
	to                   = Advance
	fromForReferences    = 0
	toReferences         = Advance
	gotoLine             = ""
	fileToOpen           = flag.String("file", "", "File to open")
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
	p                            = termenv.ColorProfile()
)

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
	//return termenv.String(fmt.Sprintf(`%s <saved "%s">`, getStatusInformation(), fileName)).Foreground(p.Color("#E88388")).String()

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

func openOSEditor(os, notesFile string) *exec.Cmd {
	if os == "windows" {
		return exec.Command("notepad", notesFile)
	}
	return exec.Command("/usr/bin/xterm", "-fa", "Monospace", "-fs", "14", "-e", "/usr/bin/vim", "+$", notesFile)
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

	file, err := os.Open(fileName)
	fileContent, err = readLines(file)
	check(err)
	defer file.Close()

	txtArea := tui.NewVBox()
	txtAreaScroll := tui.NewScrollArea(txtArea)
	txtAreaScroll.SetAutoscrollToBottom(true)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	// txtAreaBox.SetBorder(true)

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

	addCloseGotoBinding(ui, inputCommand, txtReader, txtArea)
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
	addSaveQuoteKeyBindings(ui, fileName, txtArea, txtReader, inputCommand)
	addOnSelectedReference()

	inputCommand.SetText(getStatusInformation())

	clearScreen()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}
