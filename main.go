package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"textreader/internal/file"
	"textreader/internal/keybindings"
	"textreader/internal/model"
	"textreader/internal/progress"
	"textreader/internal/references"
	"textreader/internal/terminal"
	"textreader/internal/text"
	"textreader/internal/ui"
	"textreader/internal/utils"
	"time"

	"github.com/marcusolsson/tui-go"
)

func main() {
	fileFlag := flag.String("file", "", "File to open")
	flag.Parse()
	state := model.NewAppState()
	state.FileToOpen = *fileFlag

	if err := run(state); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(state *model.AppState) error {
	fileName := state.FileToOpen
	if fileName == "" {
		return fmt.Errorf("missing file to read")
	}

	if err := file.CreateDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	var err error
	state.BannedWords, err = references.LoadNonRefsFile(model.NonRefsFileName)
	if err != nil {
		return fmt.Errorf("failed to load banned words: %w", err)
	}

	state.Sidebar.Append(state.RefsTable)

	absoluteFilePath, err := filepath.Abs(fileName)
	if err != nil {
		return fmt.Errorf("failed to resolve file path: %w", err)
	}

	latestFile, err := file.GetFileNameFromLatest(absoluteFilePath, state)
	if err != nil {
		return fmt.Errorf("failed to load latest file: %w", err)
	}

	state.From, state.To, fileName = latestFile.From, latestFile.To, latestFile.FileName

	txtFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer txtFile.Close()

	state.FileContent, err = file.ReadLines(txtFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	state.Advance = terminal.CalculateTerminalHeight()
	state.To = state.From + state.Advance
	if state.To > len(state.FileContent) {
		state.To = len(state.FileContent)
	}

	state.StartTime = time.Now()
	state.CurrentPercentage = int(progress.GetPercentage(state.To, &state.FileContent))

	txtArea := tui.NewVBox()
	txtAreaScroll := tui.NewScrollArea(txtArea)
	txtAreaScroll.SetAutoscrollToBottom(false)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	txtAreaBox.SetBorder(true)

	inputCommand := ui.NewInputCommandEntry()
	inputCommandBox := ui.NewInputCommandBox(inputCommand)

	txtReader := tui.NewVBox(txtAreaBox, inputCommandBox)
	txtReader.SetSizePolicy(tui.Expanding, tui.Expanding)

	chunk := text.GetChunk(&state.FileContent, state.From, state.To)
	text.PutText(txtArea, &chunk, txtAreaScroll, state)

	root := tui.NewHBox(txtReader, state.Sidebar)

	tuiUI, err := tui.New(root)
	if err != nil {
		return fmt.Errorf("failed to initialize UI: %w", err)
	}

	theme := tui.NewTheme()
	theme.SetStyle("label.highlight", tui.Style{
		Fg: tui.ColorBlack,
		Bg: tui.ColorGreen,
	})
	theme.SetStyle("label.wordhighlight", tui.Style{
		Fg: tui.ColorBlack,
		Bg: tui.ColorCyan,
	})
	tuiUI.SetTheme(theme)

	keybindings.AddUpDownKeyBindings(txtArea, tuiUI, inputCommand, txtAreaScroll, state)
	keybindings.AddHighlightUpDownKeyBindings(txtArea, tuiUI, inputCommand, txtAreaScroll, state)
	keybindings.AddWordLeftRightKeyBindings(txtArea, tuiUI, inputCommand, txtAreaScroll, state)
	keybindings.AddCopyWordKeyBinding(tuiUI, inputCommand, state)
	keybindings.AddGotoKeyBinding(tuiUI, txtReader, state)
	keybindings.AddShowStatusKeyBinding(tuiUI, inputCommand, state)
	keybindings.AddNewNoteKeyBinding(tuiUI, txtArea, inputCommand, fileName, txtAreaScroll, state)
	keybindings.AddCloseGotoBinding(tuiUI, inputCommand, txtReader, txtArea, txtAreaScroll, state)
	keybindings.AddSaveStatusKeyBinding(tuiUI, fileName, inputCommand, state)
	keybindings.AddShowReferencesKeyBinding(tuiUI, txtArea, txtAreaScroll, state)
	keybindings.AddAnalyzeAndFilterReferencesKeyBinding(tuiUI, state)
	keybindings.AddPercentageKeyBindings(tuiUI, inputCommand, state)
	keybindings.AddCloseApplicationKeyBinding(tuiUI, txtArea, txtReader, txtAreaScroll, state)
	keybindings.AddReferencesNavigationKeyBindings(tuiUI, state)
	keybindings.AddSaveQuoteKeyBindings(tuiUI, fileName, txtArea, inputCommand, txtAreaScroll, state)
	keybindings.AddOnSelectedReference(state)
	keybindings.AddShowMinutesTakenToReachPercentagePointKeyBinding(tuiUI, txtReader, state)
	keybindings.AddShowHelpKeyBinding(tuiUI, txtReader, state)
	keybindings.AddOpenRAEWebSite(tuiUI, inputCommand)
	keybindings.AddOpenGoodReadsWebSite(tuiUI, inputCommand)

	inputCommand.SetText(utils.GetStatusInformation(state))

	terminal.ClearScreen()

	if err := tuiUI.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}
	return nil
}
