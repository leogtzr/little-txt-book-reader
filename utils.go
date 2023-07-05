package main

import (
	"bufio"
	"fmt"
	"io"
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

func prepareReferencesBox(txtArea *tui.TextEdit) {
	txtArea.SetText("")
	txtArea.SetSizePolicy(tui.Expanding, tui.Expanding)
	txtArea.SetFocused(true)
	txtArea.SetWordWrap(true)
}

func saveStatus(fileName string, from, to int) {
	baseFileName := filepath.Base(fileName)
	f, err := os.Create(filepath.Join(home(runtime.GOOS), "ltbr", "progress", baseFileName))
	if err != nil {
		log.Fatal(err)
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
	latestFilePath := filepath.Join(home(runtime.GOOS), "ltbr", "progress", baseFileName)
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
	if len(latestFileFields) != dbFileRequiredNumbermields {
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

func readLines(file io.Reader) ([]string, error) {
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

func countWithoutWhitespaces(words []string) int {
	count := 0
	for _, w := range words {
		count += len(w)
	}
	return count
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
		box.Append(tui.NewVBox(
			tui.NewLabel(txt),
		))
	}
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		// Do not add duplicate.
		if !encountered[elements[v]] {
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

func appendLineToFile(filePath, line, sep string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n%s", sep, line)); err != nil {
		panic(err)
	}
}

func createDirectory(dirPath string) error {
	if !dirExists(dirPath) {
		if err := os.Mkdir(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}

func createDirectories() error {
	ltbrDir := filepath.Join(home(runtime.GOOS), "ltbr")
	if err := createDirectory(ltbrDir); err != nil {
		return err
	}
	return create("notes", "quotes", "progress")
}

func create(dirs ...string) error {
	for _, dir := range dirs {
		if err := createDirectory(filepath.Join(home(runtime.GOOS), "ltbr", dir)); err != nil {
			return err
		}
	}
	return nil
}

func home(opSystem string) string {
	if opSystem == "windows" {
		return os.Getenv("HOMEPATH")
	}
	return os.Getenv("HOME")
}

func getDirectoryNameForFile(dirType, fileName string) string {
	absoluteFilePath, _ := filepath.Abs(fileName)
	baseFileName := path.Base(absoluteFilePath)

	baseFileName = sanitizeFileName(baseFileName)
	notesDir := filepath.Join(home(runtime.GOOS), "ltbr", dirType, baseFileName)

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

func removeTrailingSpaces(s string) string {
	lines := strings.Split(s, "\n")
	var sb strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		sb.WriteString(strings.TrimSpace(line))
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func newInputCommandEntry() *tui.Entry {
	inputCommand := tui.NewEntry()
	inputCommand.SetFocused(true)
	inputCommand.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputCommand.SetEchoMode(tui.EchoModeNormal)

	inputCommandBox := tui.NewHBox(inputCommand)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommand
}

func newInputCommandBox(input *tui.Entry) *tui.Box {
	inputCommandBox := tui.NewHBox(input)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommandBox
}

func getPercentage(currentPosition int, fileContent *[]string) float64 {
	percent := float64(currentPosition) * 100.00
	percent = percent / float64(len(*fileContent))
	return percent
}

func removeWhiteSpaces(input string) string {
	re := regexp.MustCompile(`( |\t){2,}`)
	return re.ReplaceAllString(input, ` `)
}
