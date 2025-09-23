package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	files "textreader/internal/file"
	"textreader/internal/model"
	"textreader/internal/utils"
	"time"

	"github.com/marcusolsson/tui-go"
)

func getSavedStatusInformation(fileName string) string {
	return fmt.Sprintf(`%s <saved "%s">`, getStatusInformation(), fileName)
}

func getStatusInformation() string {
	if !model.ToggleShowStatus {
		return ""
	}

	percent := utils.GetPercentage(model.To, &model.FileContent)
	if int(percent) > model.CurrentPercentage {
		model.CurrentPercentage = int(percent)
		now := time.Now()
		model.MinutesToReachNextPercentagePoint[int(percent)] = now.Sub(model.StartTime)
		model.StartTime = now
	}

	if model.PercentagePointStats {
		return fmt.Sprintf(".   %d of %d lines (%.3f%%) [%d lines To next percentage point]                    ",
			model.To,
			len(model.FileContent), percent, utils.LinesToChangePercentagePoint(model.To, len(model.FileContent)))
	}
	return fmt.Sprintf(".   %d of %d lines (%.3f%%)                                            ",
		model.To, len(model.FileContent), percent)

}

func init() {
	if err := files.CreateDirectories(); err != nil {
		log.Fatal(err)
	}

	model.MinutesToReachNextPercentagePoint = make(map[int]time.Duration)

	// load words From file
	var err error
	model.BannedWords, err = loadNonRefsFile(model.NonRefsFileName)
	if err != nil {
		log.Fatal(err)
	}

	model.Sidebar.Append(model.RefsTable)
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
	fileName := *model.FileToOpen
	if fileName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: missing file To read")
		os.Exit(1)
	}

	var err error

	absoluteFilePath, _ := filepath.Abs(fileName)
	latestFile, err := files.GetFileNameFromLatest(absoluteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	model.From, model.To, fileName = latestFile.From, latestFile.To, latestFile.FileName

	file, err := os.Open(fileName)
	utils.Check(err)
	model.FileContent, err = files.ReadLines(file)
	utils.Check(err)
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

	model.Advance = utils.CalculateTerminalHeight()

	// Adjust To based on new Advance
	model.To = model.From + model.Advance
	if model.To > len(model.FileContent) {
		model.To = len(model.FileContent)
	}

	model.StartTime = time.Now()
	model.CurrentPercentage = int(utils.GetPercentage(model.To, &model.FileContent))

	txtArea := tui.NewVBox()
	txtAreaScroll := tui.NewScrollArea(txtArea)
	txtAreaScroll.SetAutoscrollToBottom(false)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	txtAreaBox.SetBorder(true)

	inputCommand := utils.NewInputCommandEntry()
	inputCommandBox := utils.NewInputCommandBox(inputCommand)

	txtReader := tui.NewVBox(txtAreaBox, inputCommandBox)
	txtReader.SetSizePolicy(tui.Expanding, tui.Expanding)

	chunk := GetChunk(&model.FileContent, model.From, model.To)
	PutText(txtArea, &chunk, txtAreaScroll)

	root := tui.NewHBox(txtReader, model.Sidebar)

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

	utils.ClearScreen()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
