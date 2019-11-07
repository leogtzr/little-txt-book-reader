package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/marcusolsson/tui-go"
)

// Advance ...
const Advance int = 30

// WrapMax ...
const WrapMax = 80

// GotoWidgetIndex ...
const GotoWidgetIndex = 2

const exampleBody = ``

var from = 0
var to = Advance
var gotoLine = ""
var fileToOpen = flag.String("file", "", "File to open")
var openLatestFile = flag.Bool("latest", false, "Open the latest text file")
var percentagePointStats = false
var absoluteFilePath string
var toggleShowStatus = true

// LatestFile ...
type LatestFile struct {
	FileName string
	From     int
	To       int
}

func downText(fileContent *[]string, txtArea *tui.Box) {
	if from < len(*fileContent) {
		from++
	}
	if to >= len(*fileContent) {
		return
	}

	if to < len(*fileContent) {
		to++
	}
	chunk := getChunk(fileContent, from, to)
	putText(txtArea, &chunk)
}

func upText(fileContent *[]string, txtArea *tui.Box) {
	if from <= 0 {
		return
	}

	if from > 0 {
		from--
	}

	to--

	chunk := getChunk(fileContent, from, to)
	putText(txtArea, &chunk)
}

func getSavedStatusInformation(fileContent *[]string) string {
	return fmt.Sprintf("%s <saved>", getStatusInformation(fileContent))
}

func getStatusInformation(fileContent *[]string) string {

	if !toggleShowStatus {
		return ""
	}

	percent := float64(to) * 100.00
	percent = percent / float64(len(*fileContent))
	if percentagePointStats {
		return fmt.Sprintf(".   %d of %d lines (%.3f%%) [%d lines to next percentage point]                                                            ",
			to,
			len(*fileContent), percent, linesToChangePercentagePoint(to, len(*fileContent)))
	}
	return fmt.Sprintf(".   %d of %d lines (%.3f%%)                                                            ",
		to, len(*fileContent), percent)

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

	fileContent, err := readLines(fileName)
	check(err)

	// sidebar := tui.NewVBox(
	// 	tui.NewLabel("CHANNELS"),
	// 	tui.NewLabel("general"),
	// 	tui.NewLabel("random"),
	// 	tui.NewLabel(""),
	// 	tui.NewLabel("DIRECT MESSAGES"),
	// 	tui.NewLabel("slackbot"),
	// 	tui.NewSpacer(),
	// )
	// sidebar.SetBorder(true)

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
	ui.SetKeybinding(downKeyBindingAlternative1, addDownBinding(&fileContent, txtArea, inputCommand))
	ui.SetKeybinding(downKeyBindingAlternative2, addDownBinding(&fileContent, txtArea, inputCommand))

	// Up ...
	ui.SetKeybinding(upKeyBindingAlternative1, addUpBinding(&fileContent, txtArea, inputCommand))
	ui.SetKeybinding(upKeyBindingAlternative2, addUpBinding(&fileContent, txtArea, inputCommand))

	// go to:
	ui.SetKeybinding(gotoKeyBindingAlterntive1, func() {
		addGotoWidget(txtReader)
	})

	// show status key binding:
	ui.SetKeybinding(showStatusKeyBinding, func() {
		toggleShowStatus = !toggleShowStatus
		inputCommand.SetText(getStatusInformation(&fileContent))
	})

	noteBox := tui.NewTextEdit()
	noteBox.SetText("")

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
		// Save note ...
		notesDir := filepath.Join(os.Getenv("HOME"), "txtnotes")
		if !dirExists(notesDir) {
			err := os.Mkdir(notesDir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: creating notes dir: %s", notesDir)
			}
		}
		rand.Seed(time.Now().UnixNano())
		absoluteFilePath, _ := filepath.Abs(fileName)
		baseFileName := path.Base(absoluteFilePath)
		noteFileName := fmt.Sprintf("%d-%s", rand.Intn(150), baseFileName)

		noteContent := noteBox.Text()

		ioutil.WriteFile(filepath.Join(notesDir, noteFileName), []byte(noteContent), 0666)
		txtReader.Remove(0)
	})

	ui.SetKeybinding(closeGotoKeyBindingAlternative1, func() {
		// Go to the specified line
		inputCommand.SetText(getStatusInformation(&fileContent))

		gotoLineNumber := getNumberLineGoto(gotoLine)
		gotoLineNumberDigits, err := strconv.ParseInt(gotoLineNumber, 10, 64)
		if err != nil {
			return
		}
		if int(gotoLineNumberDigits) < (len(fileContent) - Advance) {
			from = int(gotoLineNumberDigits)
			to = from + Advance
			putText(txtArea, &chunk)
			inputCommand.SetText(getStatusInformation(&fileContent))
		}
		txtReader.Remove(GotoWidgetIndex)
		inputCommand.SetText(getStatusInformation(&fileContent))
	})

	ui.SetKeybinding(saveStatusKeyBindingAlternative1, func() {
		absoluteFilePath, _ := filepath.Abs(fileName)
		saveStatus(absoluteFilePath, from, to)

		inputCommand.SetText(getSavedStatusInformation(&fileContent))
	})

	// Enable percentage stags
	ui.SetKeybinding(nextPercentagePointKeyBindingAlternative1, func() {
		percentagePointStats = !percentagePointStats
		inputCommand.SetText(getStatusInformation(&fileContent))
	})

	ui.SetKeybinding(closeApplicationKeyBindingAlterntive1, func() {
		ui.Quit()
	})

	inputCommand.SetText(getStatusInformation(&fileContent))

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
