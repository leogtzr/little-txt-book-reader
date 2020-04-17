package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/marcusolsson/tui-go"
)

const (
	downKeyBindingAlternative1                  = "Alt+j"
	downKeyBindingAlternative2                  = "Down"
	upKeyBindingAlternative1                    = "Alt+k"
	upKeyBindingAlternative2                    = "Up"
	gotoKeyBindingAlternative1                  = "Alt+g"
	newNoteKeyBindingAlternative1               = "Alt+n"
	showStatusKeyBinding                        = "Alt+."
	closeGotoKeyBindingAlternative1             = "r"
	saveStatusKeyBindingAlternative1            = "s"
	nextPercentagePointKeyBindingAlternative1   = "Alt+p"
	showReferencesKeyBindingAlternative1        = "Alt+r"
	closeReferencesWindowKeyBindingAlternative1 = "Alt+q"
	closeApplicationKeyBindingAlternative1      = "Esc"
	analyzeAndFilterReferencesKeyBinding        = "Alt+b"
	saveQuoteKeyBindingAlternative1             = "Alt+q"
	maxNumberOfElementsInGUIBox                 = 1000
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func loadNonRefsFile(path string) ([]string, error) {
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

func prepareReferencesBox(guiComponent *tui.TextEdit) {
	guiComponent.SetText("")
	guiComponent.SetSizePolicy(tui.Expanding, tui.Expanding)
	guiComponent.SetFocused(true)
	guiComponent.SetWordWrap(true)
}

func saveStatus(fileName string, from, to int) {
	baseFileName := filepath.Base(fileName)
	f, err := os.Create(filepath.Join(home(), "ltbr", "progress", baseFileName))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s|%d|%d", fileName, from, to))
}

func getNumberLineGoto(line string) string {
	rgx, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return rgx.ReplaceAllString(line, "")
}

func percent(currentNumberLine, totalLines int) float64 {
	return float64(currentNumberLine*100.0) / float64(totalLines)
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

func getFileNameFromLatest(filePath string) (LatestFile, error) {
	baseFileName := filepath.Base(filePath)
	latestFilePath := filepath.Join(home(), "ltbr", "progress", baseFileName)
	latestFile := LatestFile{FileName: filePath, From: 0, To: Advance}
	if !exists(latestFilePath) {
		return latestFile, nil
	}

	f, err := os.Open(latestFilePath)
	if err != nil {
		return LatestFile{}, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	latestFileFields := strings.Split(string(content), "|")
	if len(latestFileFields) != dbFileRequieredNumberFields {
		return LatestFile{}, fmt.Errorf("Wrong format in '%s'", latestFilePath)
	}

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

func check(err error) {
	if err != nil {
		panic(err)
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
	currentNavMode = gotoNavigationMode
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

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	return result
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

	uniqueReferences := removeDuplicates(refs)

	referencesNoBannedWords := make([]string, 0)
	for _, word := range uniqueReferences {
		if !contains(bannedWords, word) {
			referencesNoBannedWords = append(referencesNoBannedWords, word)
		}
	}

	return referencesNoBannedWords
}

func loadReferences() {
	if len(references) == 0 {
		references = extractReferencesFromFileContent(&fileContent)
	}
}

func findAndRemove(s *[]string, e string) {
	for i, v := range *s {
		if v == e {
			*s = append((*s)[:i], (*s)[i+1:]...)
			break
		}
	}
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func paginate(x []string, skip, size int) []string {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}

func appendLineToFile(filePath, line string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString("\n" + line); err != nil {
		panic(err)
	}
}

func createDirectory(dirPath string) error {
	if !dirExists(dirPath) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDirectories() error {
	ltbrDir := filepath.Join(home(), "ltbr")
	if err := createDirectory(ltbrDir); err != nil {
		return err
	}
	return create("notes", "quotes", "progress")
}

func create(dirs ...string) error {
	for _, dir := range dirs {
		if err := createDirectory(filepath.Join(home(), "ltbr", dir)); err != nil {
			return err
		}
	}
	return nil
}

func home() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEPATH")
	}
	return os.Getenv("HOME")
}

func getNotesDirectoryNameForFile(fileName string) string {
	absoluteFilePath, _ := filepath.Abs(fileName)
	baseFileName := path.Base(absoluteFilePath)

	baseFileName = sanitizeFileName(baseFileName)
	notesDir := filepath.Join(home(), "ltbr", "notes", baseFileName)

	return notesDir
}

func getDirectoryNameForFile(dirType, fileName string) string {
	absoluteFilePath, _ := filepath.Abs(fileName)
	baseFileName := path.Base(absoluteFilePath)

	baseFileName = sanitizeFileName(baseFileName)
	notesDir := filepath.Join(home(), "ltbr", dirType, baseFileName)

	return notesDir
}

func longestLineLength(text string) int {
	if len(text) == 0 {
		return 0
	}

	lines := strings.Split(text, "\n")
	longest := len(lines[0])
	for _, line := range lines {
		if len(line) > longest {
			longest = len(line)
		}
	}
	return longest
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
