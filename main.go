package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	// "os/exec"
	"path/filepath"
	files "textreader/internal/file"
	"textreader/internal/keybindings"
	"textreader/internal/model"
	"textreader/internal/references"
	"textreader/internal/text"
	"textreader/internal/utils"
	"time"

	"github.com/marcusolsson/tui-go"
)

func init() {
	if err := files.CreateDirectories(); err != nil {
		log.Fatal(err)
	}

	model.MinutesToReachNextPercentagePoint = make(map[int]time.Duration)

	// load words From file
	var err error
	model.BannedWords, err = references.LoadNonRefsFile(model.NonRefsFileName)
	if err != nil {
		log.Fatal(err)
	}

	model.Sidebar.Append(model.RefsTable)
	// Sidebar.Append(refsStatus)
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

	chunk := text.GetChunk(&model.FileContent, model.From, model.To)
	text.PutText(txtArea, &chunk, txtAreaScroll)

	root := tui.NewHBox(txtReader, model.Sidebar)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	keybindings.AddUpDownKeyBindings(txtArea, ui, inputCommand, txtAreaScroll)
	keybindings.AddGotoKeyBinding(ui, txtReader)
	keybindings.AddShowStatusKeyBinding(ui, inputCommand)
	keybindings.AddNewNoteKeyBinding(ui, txtArea, inputCommand, fileName, txtAreaScroll)
	keybindings.AddCloseGotoBinding(ui, inputCommand, txtReader, txtArea, txtAreaScroll)
	keybindings.AddSaveStatusKeyBinding(ui, fileName, inputCommand)
	keybindings.AddShowReferencesKeyBinding(ui, txtArea, txtAreaScroll)
	keybindings.AddAnalyzeAndFilterReferencesKeyBinding(ui)
	keybindings.AddPercentageKeyBindings(ui, inputCommand)
	keybindings.AddCloseApplicationKeyBinding(ui, txtArea, txtReader, txtAreaScroll)
	keybindings.AddReferencesNavigationKeyBindings(ui)
	keybindings.AddSaveQuoteKeyBindings(ui, fileName, txtArea, inputCommand, txtAreaScroll)
	keybindings.AddOnSelectedReference()
	keybindings.AddShowMinutesTakenToReachPercentagePointKeyBinding(ui, txtReader)
	keybindings.AddShowHelpKeyBinding(ui, txtReader)
	keybindings.AddOpenRAEWebSite(ui, inputCommand)
	keybindings.AddOpenGoodReadsWebSite(ui, inputCommand)

	inputCommand.SetText(utils.GetStatusInformation())

	utils.ClearScreen()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
