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
	"golang.org/x/term"
)

func updateRangesUp() {
	if From <= 0 {
		return
	}

	if From > 0 {
		From--
	}

	To--
}

func updateRangesReferenceUp() {
	if FromForReferences <= 0 {
		return
	}

	if FromForReferences > 0 {
		FromForReferences--
	}

	ToReferences--
}

func updateRangesDown() {
	if From < len(FileContent) {
		From++
	}

	if To >= len(FileContent) {
		return
	}

	if To < len(FileContent) {
		To++
	}
}

func updateRangesReferenceDown() {
	if FromForReferences < len(References) {
		FromForReferences++
	}

	if ToReferences >= len(References) {
		return
	}

	if ToReferences < len(References) {
		ToReferences++
	}
}

func getSavedStatusInformation(fileName string) string {
	return fmt.Sprintf(`%s <saved "%s">`, getStatusInformation(), fileName)
}

func getStatusInformation() string {
	if !ToggleShowStatus {
		return ""
	}

	percent := getPercentage(To, &FileContent)
	if int(percent) > CurrentPercentage {
		CurrentPercentage = int(percent)
		now := time.Now()
		MinutesToReachNextPercentagePoint[int(percent)] = now.Sub(StartTime)
		StartTime = now
	}

	if PercentagePointStats {
		return fmt.Sprintf(".   %d of %d lines (%.3f%%) [%d lines To next percentage point]                    ",
			To,
			len(FileContent), percent, linesToChangePercentagePoint(To, len(FileContent)))
	}
	return fmt.Sprintf(".   %d of %d lines (%.3f%%)                                            ",
		To, len(FileContent), percent)

}

func init() {
	if err := createDirectories(); err != nil {
		log.Fatal(err)
	}

	MinutesToReachNextPercentagePoint = make(map[int]time.Duration)

	// load words From file
	var err error
	BannedWords, err = loadNonRefsFile(NonRefsFileName)
	if err != nil {
		log.Fatal(err)
	}

	Sidebar.Append(RefsTable)
	// Sidebar.Append(refsStatus)
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
	fileName := *FileToOpen
	if fileName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: missing file To read")
		os.Exit(1)
	}

	var err error

	absoluteFilePath, _ := filepath.Abs(fileName)
	latestFile, err := getFileNameFromLatest(absoluteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	From, To, fileName = latestFile.From, latestFile.To, latestFile.FileName

	file, err := os.Open(fileName)
	check(err)
	FileContent, err = readLines(file)
	check(err)
	defer file.Close()

	//fd := int(os.Stdout.Fd())
	//_, height, err := term.GetSize(fd)
	//if err == nil {
	//	Advance = height - 5 // Subtract for borders, input bar, status, etc. Adjust if needed.
	//	if Advance < 10 {
	//		Advance = 30 // Fallback minimum
	//	}
	//} else {
	//	Advance = 30 // Fallback if detection fails
	//}

	Advance = calculateAdvanceHeight()

	// Adjust To based on new Advance
	To = From + Advance
	if To > len(FileContent) {
		To = len(FileContent)
	}

	StartTime = time.Now()
	CurrentPercentage = int(getPercentage(To, &FileContent))

	txtArea := tui.NewVBox()
	txtAreaScroll := tui.NewScrollArea(txtArea)
	txtAreaScroll.SetAutoscrollToBottom(false)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	txtAreaBox.SetBorder(true)

	inputCommand := newInputCommandEntry()
	inputCommandBox := newInputCommandBox(inputCommand)

	txtReader := tui.NewVBox(txtAreaBox, inputCommandBox)
	txtReader.SetSizePolicy(tui.Expanding, tui.Expanding)

	chunk := getChunk(&FileContent, From, To)
	putText(txtArea, &chunk, txtAreaScroll)

	root := tui.NewHBox(txtReader, Sidebar)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	addUpDownKeyBindings(txtArea, ui, inputCommand, txtAreaScroll)
	addGotoKeyBinding(ui, txtReader)
	addShowStatusKeyBinding(ui, inputCommand)
	addNewNoteKeyBinding(ui, txtArea, inputCommand, fileName, txtAreaScroll)
	addCloseGotoBinding(ui, inputCommand, txtReader, txtArea, txtAreaScroll)
	addSaveStatusKeyBinding(ui, fileName, inputCommand)
	addShowReferencesKeyBinding(ui, txtArea, txtAreaScroll)
	addAnalyzeAndFilterReferencesKeyBinding(ui)
	addPercentageKeyBindings(ui, inputCommand)
	addCloseApplicationKeyBinding(ui, txtArea, txtReader, txtAreaScroll)
	addReferencesNavigationKeyBindings(ui)
	addSaveQuoteKeyBindings(ui, fileName, txtArea, inputCommand, txtAreaScroll)
	addOnSelectedReference()
	addShowMinutesTakenToReachPercentagePointKeyBinding(ui, txtReader)
	addShowHelpKeyBinding(ui, txtReader)
	addOpenRAEWebSite(ui, inputCommand)
	addOpenGoodReadsWebSite(ui, inputCommand)

	inputCommand.SetText(getStatusInformation())

	ClearScreen()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func calculateAdvanceHeight() int {
	advance := 45
	fd := int(os.Stdout.Fd())
	_, height, err := term.GetSize(fd)
	if err == nil {
		advance = height - 5 // Subtract for borders, input bar, status, etc. Adjust if needed.
	}

	return advance
}
