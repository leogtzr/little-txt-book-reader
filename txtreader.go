package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/marcusolsson/tui-go"
)

var (
	from                              = 0
	to                                = Advance
	fromForReferences                 = 0
	toReferences                      = Advance
	gotoLine                          = ""
	fileToOpen                        = flag.String("file", "", "File to open")
	percentagePointStats              = false
	toggleShowStatus                  = true
	references                        = []string{}
	fileContent                       = []string{}
	currentNavMode                    = readingNavigationMode
	bannedWords                       = []string{}
	sidebar                           = tui.NewVBox()
	refsTable                         = tui.NewTable(0, 0)
	pageIndex                         = 0
	minutesToReachNextPercentagePoint map[int]time.Duration
	startTime                         time.Time
	currentPercentage                 int
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
}

func getStatusInformation() string {

	if !toggleShowStatus {
		return ""
	}

	percent := getPercentage(to, &fileContent)
	if int(percent) > currentPercentage {
		currentPercentage = int(percent)
		now := time.Now()
		minutesToReachNextPercentagePoint[int(percent)] = now.Sub(startTime)
		startTime = now
	}

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

	minutesToReachNextPercentagePoint = make(map[int]time.Duration)

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
	if os == "darwin" {
		script := fmt.Sprintf(`tell application "Terminal"
	activate
	do script "vim +$ %s; exit"
end tell`, notesFile)
		return exec.Command("osascript", "-e", script)
	}
	return exec.Command("/usr/bin/xterm", "-fa", "Monospace", "-fs", "14", "-e", "/usr/bin/vim", "+$", notesFile)
}

func main() {

	flag.Parse()
	fileName := *fileToOpen
	if fileName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: missing file to read")
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
	check(err)
	fileContent, err = readLines(file)
	check(err)
	defer file.Close()

	startTime = time.Now()
	currentPercentage = int(getPercentage(to, &fileContent))

	txtArea := tui.NewVBox()
	txtAreaScroll := tui.NewScrollArea(txtArea)
	txtAreaScroll.SetAutoscrollToBottom(true)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	txtAreaBox.SetBorder(true)

	inputCommand := newInputCommandEntry()
	inputCommandBox := newInputCommandBox(inputCommand)

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
	addShowStatusKeyBinding(ui, inputCommand)
	addNewNoteKeyBinding(ui, txtArea, inputCommand, fileName)
	addCloseGotoBinding(ui, inputCommand, txtReader, txtArea)
	addSaveStatusKeyBinding(ui, fileName, inputCommand)
	addShowReferencesKeyBinding(ui, txtArea)
	addAnalyzeAndFilterReferencesKeyBinding(ui)
	addPercentageKeyBindings(ui, inputCommand)
	addCloseApplicationKeyBinding(ui, txtArea, txtReader)
	addReferencesNavigationKeyBindings(ui)
	addSaveQuoteKeyBindings(ui, fileName, txtArea, inputCommand)
	addOnSelectedReference()
	addShowMinutesTakenToReachPercentagePointKeyBinding(ui, txtReader)
	addShowHelpKeyBinding(ui, txtReader)
	addOpenRAEWebSite(ui, inputCommand)
	addOpenGoodReadsWebSite(ui, inputCommand)

	inputCommand.SetText(getStatusInformation())

	clearScreen()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
