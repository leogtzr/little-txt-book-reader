package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/marcusolsson/tui-go"
)

func addDownBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		downText(box, txtAreaScroll)
		input.SetText(getStatusInformation())
	}
}

func addUpDownKeyBindings(txtArea *tui.Box, ui tui.UI, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(DownKeyBindingAlternative1, addDownBinding(txtArea, inputCommand, txtAreaScroll))
	ui.SetKeybinding(DownKeyBindingAlternative2, addDownBinding(txtArea, inputCommand, txtAreaScroll))

	ui.SetKeybinding(UpKeyBindingAlternative1, addUpBinding(txtArea, inputCommand, txtAreaScroll))
	ui.SetKeybinding(UpKeyBindingAlternative2, addUpBinding(txtArea, inputCommand, txtAreaScroll))
}

func addShowStatusKeyBinding(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(showStatusKeyBinding, func() {
		ToggleShowStatus = !ToggleShowStatus
		inputCommand.SetText(getStatusInformation())
	})
}

func addUpBinding(box *tui.Box, input *tui.Entry, txtAreaScroll *tui.ScrollArea) func() {
	return func() {
		upText(box, txtAreaScroll)
		input.SetText(getStatusInformation())
	}
}

func addSaveStatusKeyBinding(ui tui.UI, fileName string, inputCommand *tui.Entry) {
	baseFileName := filepath.Base(fileName)
	ui.SetKeybinding(SaveStatusKeyBindingAlternative1, func() {
		saveStatus(fileName, From, To)
		inputCommand.SetText(getSavedStatusInformation(baseFileName))
	})
}

func addCloseApplicationKeyBinding(ui tui.UI, txtArea, txtReader *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(CloseApplicationKeyBindingAlternative1, func() {

		switch CurrentNavMode {
		case ShowReferencesNavigationMode:
			chunk := getChunk(&FileContent, From, To)
			putText(txtArea, &chunk, txtAreaScroll)
			CurrentNavMode = ReadingNavigationMode
		case AnalyzeAndFilterReferencesNavigationMode:
			chunk := getChunk(&FileContent, From, To)
			putText(txtArea, &chunk, txtAreaScroll)
			CurrentNavMode = ReadingNavigationMode
			RefsTable.SetFocused(false)
		case GotoNavigationMode, ShowTimePercentagePointsMode, ShowHelpMode:
			txtReader.Remove(GotoWidgetIndex)
			CurrentNavMode = ReadingNavigationMode
		default:
			ClearScreen()
			ui.Quit()
		}
	})
}

func addPercentageKeyBindings(ui tui.UI, inputCommand *tui.Entry) {
	// Enable percentage tags
	ui.SetKeybinding(NextPercentagePointKeyBindingAlternative1, func() {
		PercentagePointStats = !PercentagePointStats
		inputCommand.SetText(getStatusInformation())
	})
}

func addShowReferencesKeyBinding(ui tui.UI, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(ShowReferencesKeyBindingAlternative1, func() {
		CurrentNavMode = ShowReferencesNavigationMode
		loadReferences()
		chunk := getChunk(&References, FromForReferences, ToReferences)
		putText(txtArea, &chunk, txtAreaScroll)
	})
}

func addReferencesNavigationKeyBindings(ui tui.UI) {
	// Next References ...
	ui.SetKeybinding("Right", func() {
		if PageIndex >= len(References) {
			return
		}
		PageIndex += PageSize
		prepareTableForReferences()
	})

	// Previous References ...
	ui.SetKeybinding("Left", func() {
		if PageIndex < PageSize {
			return
		}
		PageIndex -= PageSize
		prepareTableForReferences()
	})
}

func addSaveQuoteKeyBindings(ui tui.UI, fileName string, txtArea *tui.Box, inputCommand *tui.Entry, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(SaveQuoteKeyBindingAlternative1, func() {
		oldStdout, oldStdin, oldSterr := os.Stdout, os.Stdin, os.Stderr

		quotesFile := getDirectoryNameForFile("quotes", fileName)

		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}

		clipBoardText = removeTrailingSpaces(clipBoardText)
		clipBoardText = removeWhiteSpaces(clipBoardText)
		appendLineToFile(quotesFile, clipBoardText, "\n__________")

		cmd := openOSEditor(runtime.GOOS, quotesFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			panic(cmdErr)
		}

		os.Stdout, os.Stdin, os.Stderr = oldStdout, oldStdin, oldSterr

		// txtReader.SetBorder(true)
		chunk := getChunk(&FileContent, From, To)
		putText(txtArea, &chunk, txtAreaScroll)
		inputCommand.SetText(getStatusInformation())
	})
}

func prepareTableForReferences() {
	RefsTable.RemoveRows()
	references := paginate(References, PageIndex, PageSize)
	for _, ref := range references {
		RefsTable.AppendRow(tui.NewLabel(ref))
	}
	RefsTable.SetSelected(0)
}

func addOnSelectedReference() {
	RefsTable.OnItemActivated(func(tui *tui.Table) {

		itemIndexToRemove := tui.Selected()
		itemToAddToNonRefs := References[PageIndex+itemIndexToRemove]
		// References = remove(References, itemIndexToRemove)
		findAndRemove(&References, itemToAddToNonRefs)
		prepareTableForReferences()

		if !contains(BannedWords, itemToAddToNonRefs) {
			appendLineToFile(NonRefsFileName, itemToAddToNonRefs, "")
		}
	})
}

func addGotoKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(GotoKeyBindingAlternative1, func() {
		addGotoWidget(txtReader)
	})
}

func addCloseGotoBinding(ui tui.UI, inputCommand *tui.Entry, txtReader, txtArea *tui.Box, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(CloseGotoKeyBindingAlternative1, func() {
		// Go To the specified line
		inputCommand.SetText(getStatusInformation())

		gotoLineNumber := getNumberLineGoto(GotoLine)
		gotoLineNumberDigits, err := strconv.ParseInt(gotoLineNumber, 10, 64)
		if err != nil {
			return
		}
		if int(gotoLineNumberDigits) < (len(FileContent) - Advance) {
			From = int(gotoLineNumberDigits)
			To = From + Advance
			chunk := getChunk(&FileContent, From, To)
			putText(txtArea, &chunk, txtAreaScroll)
			inputCommand.SetText(getStatusInformation())
		}
		txtReader.Remove(GotoWidgetIndex)
		inputCommand.SetText(getStatusInformation())
		CurrentNavMode = ReadingNavigationMode
	})
}

func addNewNoteKeyBinding(ui tui.UI, txtArea *tui.Box, inputCommand *tui.Entry, fileName string, txtAreaScroll *tui.ScrollArea) {
	ui.SetKeybinding(NewNoteKeyBindingAlternative1, func() {

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
		chunk := getChunk(&FileContent, From, To)
		putText(txtArea, &chunk, txtAreaScroll)
		inputCommand.SetText(getStatusInformation())
	})
}

func addAnalyzeAndFilterReferencesKeyBinding(ui tui.UI) {
	ui.SetKeybinding(AnalyzeAndFilterReferencesKeyBinding, func() {
		CurrentNavMode = AnalyzeAndFilterReferencesNavigationMode
		Sidebar.SetTitle("References ... ")
		Sidebar.SetBorder(true)
		RefsTable.SetColumnStretch(0, 0)
		loadReferences()

		RefsTable.RemoveRows()
		prepareTableForReferences()
		RefsTable.SetFocused(true)
	})
}

func browserOpenURLCommand(osName, url string) *exec.Cmd {
	switch osName {
	case "linux":
		return exec.Command("xdg-open", url)
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		return exec.Command("open", url)
	default:
		return nil
	}
}

func addOpenRAEWebSite(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(OpenRAEWebSiteKeyBinging, func() {
		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}
		if len(strings.TrimSpace(clipBoardText)) == 0 {
			return
		}
		raeURL := fmt.Sprintf("https://dle.rae.es/%s", clipBoardText)
		if err = browserOpenURLCommand(runtime.GOOS, raeURL).Start(); err != nil {
			inputCommand.SetText(err.Error())
			return
		}
	})
}

func browserOpenGoodReadsURLCommand(osName, url string) *exec.Cmd {
	switch osName {
	case "linux":
		return exec.Command("xdg-open", url)
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		return exec.Command("open", url)
	default:
		return nil
	}
}

func addOpenGoodReadsWebSite(ui tui.UI, inputCommand *tui.Entry) {
	ui.SetKeybinding(OpenGoodReadsWebSiteKeyBinding, func() {
		clipBoardText, err := clipboard.ReadAll()
		if err != nil {
			inputCommand.SetText(err.Error())
			return
		}
		if len(strings.TrimSpace(clipBoardText)) == 0 {
			return
		}
		goodreadsURL := fmt.Sprintf(`https://www.goodreads.com/search?q=%s`, url.QueryEscape(clipBoardText))
		if err = browserOpenGoodReadsURLCommand(runtime.GOOS, goodreadsURL).Start(); err != nil {
			inputCommand.SetText(err.Error())
			return
		}
	})
}

func addShowMinutesTakenToReachPercentagePointKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(ShowMinutesTakenToReachPercentagePointKeyBinding, func() {

		// Check if we are already in that mode ...
		if CurrentNavMode == ShowTimePercentagePointsMode {
			return
		}

		CurrentNavMode = ShowTimePercentagePointsMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range MinutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		for _, v := range percentages {
			duration := MinutesToReachNextPercentagePoint[v]
			strs = append(strs, fmt.Sprintf("%d%% took you %.1f minutes", v, duration.Minutes()))
		}

		l.AddItems(strs...)
		s := tui.NewScrollArea(l)
		s.SetFocused(true)

		txtReader.Append(s)

		ui.SetKeybinding("Alt+Up", func() { s.Scroll(0, -1) })
		ui.SetKeybinding("Alt+Down", func() { s.Scroll(0, 1) })
	})
}

func addShowHelpKeyBinding(ui tui.UI, txtReader *tui.Box) {
	ui.SetKeybinding(ShowHelpKeyBinding, func() {

		// Check if we are already in that mode ...
		if CurrentNavMode == ShowHelpMode {
			return
		}

		CurrentNavMode = ShowHelpMode

		l := tui.NewList()
		var strs []string

		percentages := make([]int, 0)
		for p := range MinutesToReachNextPercentagePoint {
			percentages = append(percentages, p)
		}
		sort.Ints(percentages)

		addKeyBindingDescription(fmt.Sprintf("%10s -> Go Down", DownKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Go Up", UpKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Go To", GotoKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> New Note", NewNoteKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Show Status", showStatusKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the Goto Dialog", CloseGotoKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Save Progress", SaveStatusKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows Next Percentage Point Stats", NextPercentagePointKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows the References Dialog", ShowReferencesKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the References Dialog", CloseReferencesWindowKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Closes the program", CloseApplicationKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Analyze and filter References", AnalyzeAndFilterReferencesKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Add a Quote, gets the text From the clipboard.", SaveQuoteKeyBindingAlternative1), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows Time Stats for each percentage point.", ShowMinutesTakenToReachPercentagePointKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Shows this Dialog", ShowHelpKeyBinding), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens RAE Web site with search From the clipboard.", OpenRAEWebSiteKeyBinging), &strs)
		addKeyBindingDescription(fmt.Sprintf("%10s -> Opens GoodReads Web site with search From the clipboard.", OpenGoodReadsWebSiteKeyBinding), &strs)

		l.AddItems(strs...)
		s := tui.NewScrollArea(l)
		s.SetFocused(true)

		txtReader.Append(s)

		ui.SetKeybinding("Alt+Up", func() { s.Scroll(0, -1) })
		ui.SetKeybinding("Alt+Down", func() { s.Scroll(0, 1) })
	})
}

func addKeyBindingDescription(desc string, keyBindings *[]string) {
	*keyBindings = append(*keyBindings, desc)
}
