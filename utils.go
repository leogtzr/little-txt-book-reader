package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/marcusolsson/tui-go"
)

const (
	downKeyBindingAlternative1                  = "Alt+j"
	downKeyBindingAlternative2                  = "Down"
	upKeyBindingAlternative1                    = "Alt+k"
	upKeyBindingAlternative2                    = "Up"
	gotoKeyBindingAlterntive1                   = "Alt+g"
	newNoteKeyBindingAlternative1               = "Alt+n"
	saveNoteKeyBindingAlternative1              = "Alt+s"
	showStatusKeyBinding                        = "Alt+."
	closeGotoKeyBindingAlternative1             = "r"
	saveStatusKeyBindingAlternative1            = "s"
	nextPercentagePointKeyBindingAlternative1   = "Alt+p"
	showReferencesKeyBindingAlternative1        = "Alt+r"
	closeReferencesWindowKeyBindingAlternative1 = "Alt+q"
	closeApplicationKeyBindingAlternative1      = "Esc"
	maxNumberOfElementsInGUIBox                 = 1000
)

func prepareNewNoteBox(noteBox *tui.TextEdit) {
	noteBox.SetText("")
	noteBox.SetSizePolicy(tui.Expanding, tui.Expanding)
	noteBox.SetFocused(true)
	noteBox.SetWordWrap(true)
}

func prepareReferencesBox(guiComponent *tui.TextEdit) {
	guiComponent.SetText("")
	guiComponent.SetSizePolicy(tui.Expanding, tui.Expanding)
	guiComponent.SetFocused(true)
	guiComponent.SetWordWrap(true)
}

func saveStatus(fileName string, from, to int) {
	// write from, to y el nombre del archivo ...
	homeDir := os.Getenv("HOME")
	f, err := os.Create(filepath.Join(homeDir, "txtread"))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s|%d|%d", fileName, from, to))
}

func saveReadingStatus(fileName string, from, to int) func() {
	return func() {
		absoluteFilePath, _ := filepath.Abs(fileName)
		saveStatus(absoluteFilePath, from, to)
	}
}

func getNumberLineGoto(line string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(line, "")
}

func percent(i, totalLines int) float64 {
	return float64(i*100.0) / float64(totalLines)
}

func linesToChangePercentagePoint(currentLine, totalLines int) int {
	start := currentLine
	linesToChangePercentage := -1
	percentageWithCurrentLine := int(percent(currentLine, totalLines))
	for {
		currentLine++
		nextPercentage := int(percent(currentLine, totalLines))
		if nextPercentage > percentageWithCurrentLine {
			linesToChangePercentage = currentLine
			break
		}
	}

	return linesToChangePercentage - start
}

func getFileNameFromLatest() (LatestFile, error) {
	homeDir := os.Getenv("HOME")
	latestFilePath := filepath.Join(homeDir, "txtread")

	if !exists(latestFilePath) {
		return LatestFile{}, fmt.Errorf("'%s' does not exist", latestFilePath)
	}

	f, err := os.Open(latestFilePath)
	if err != nil {
		return LatestFile{}, err
	}
	defer f.Close()
	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		return LatestFile{}, err
	}
	latestFileFields := strings.Split(string(fileContent), "|")
	if len(latestFileFields) != 3 {
		return LatestFile{}, fmt.Errorf("Wrong format in '%s'", latestFilePath)
	}
	latestFile := LatestFile{}
	latestFile.FileName = latestFileFields[0]
	fromInt, _ := strconv.ParseInt(latestFileFields[1], 10, 32)
	toInt, _ := strconv.ParseInt(latestFileFields[2], 10, 32)

	latestFile.From = int(fromInt)
	latestFile.To = int(toInt)
	return latestFile, nil
}

func dirExists(dirPath string) bool {
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getChunk(content *[]string, from, to int) []string {
	return (*content)[from:to]
}

func clearBox(box *tui.Box, contentLen int) {
	for i := 0; i < contentLen; i++ {
		box.Append(tui.NewHBox(
			tui.NewLabel(""),
			tui.NewSpacer(),
		))
	}
}

func needsSemiWrap(line string) bool {
	len := len(line)
	if len < (WrapMax / 2) {
		return false
	}
	return (len > (WrapMax / 2)) && (len < WrapMax)
}

func countWithoutWhitespaces(words []string) int {
	count := 0
	for _, w := range words {
		count += len(w)
	}
	return count
}

func wrap(line string) string {
	if !needsSemiWrap(line) {
		return line
	}
	fields := strings.Fields(line)
	numberOfWords := len(fields)
	countWithoutSpaces := countWithoutWhitespaces(fields)
	wrapLength := WrapMax - countWithoutSpaces
	if numberOfWords == 1 || numberOfWords == 0 {
		return line
	}
	return fmt.Sprintf("%s", strings.Join(fields, strings.Repeat(" ", wrapLength/(numberOfWords-1))))
}

func addGotoWidget(box *tui.Box) {
	gotoInput := tui.NewTextEdit()
	gotoInput.SetText("Go to line: ")
	gotoInput.SetFocused(true)
	gotoInput.OnTextChanged(func(entry *tui.TextEdit) {
		gotoLine = entry.Text()
	})
	box.Append(gotoInput)
}

func exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

func addUpBinding(box *tui.Box, input *tui.Entry) func() {
	return func() {
		upText(box)
		input.SetText(getStatusInformation())
	}
}

func addDownBinding(box *tui.Box, input *tui.Entry) func() {
	return func() {
		downText(box)
		input.SetText(getStatusInformation())
	}
}

func putText(box *tui.Box, content *[]string) {
	clearBox(box, len(*content))

	/*
		Had to introduce this code to reduce the number of elements added to the
		GUI, otherwise the memory would be increasing all the time ...
	*/
	if box.Length() > 0 {
		for i := 0; i < box.Length(); i++ {
			box.Remove(i)
		}
	}

	for _, txt := range *content {
		txt = strings.Replace(txt, " ", " ", -1)
		txt = strings.Replace(txt, "\t", "    ", -1)
		if *wrapText {
			txt = wrap(txt)
		}
		box.Append(tui.NewVBox(
			tui.NewLabel(txt),
		))
	}
}

func extractReferencesFromFileContent(fileContent *[]string) []string {
	refs := make([]string, 0)
	i := 0
	for _, lineInFile := range *fileContent {
		i++
		lineInFile = strings.TrimSpace(lineInFile)
		if lineInFile == "" || len(lineInFile) == 0 {
			continue
		}

		r := extractReferences(lineInFile)
		if len(r) > 0 {
			for _, ref := range r {
				refs = append(refs, ref)
			}
		}
	}

	return refs
}
