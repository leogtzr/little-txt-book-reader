package utils

import (
	"os"
	"os/exec"
	"regexp"
	"textreader/internal/model"

	"github.com/marcusolsson/tui-go"
	"golang.org/x/term"
)

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func GetNumberLineGoto(line string) string {
	rgx, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return rgx.ReplaceAllString(line, "")
}

func Percent(currentNumberLine, totalLines int) float64 {
	return float64(currentNumberLine*100.0) / float64(totalLines)
}

func LinesToChangePercentagePoint(currentLine, totalLines int) int {
	start := currentLine
	linesToChangePercentage := -1
	percentageWithCurrentLine := int(Percent(currentLine, totalLines))
	for {
		currentLine++
		nextPercentage := int(Percent(currentLine, totalLines))
		if nextPercentage > percentageWithCurrentLine {
			linesToChangePercentage = currentLine
			break
		}
	}

	return linesToChangePercentage - start
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: move this to the UI package
func AddGotoWidget(box *tui.Box) {
	gotoInput := tui.NewTextEdit()
	gotoInput.SetText("Go To line: ")
	gotoInput.SetFocused(true)
	gotoInput.OnTextChanged(func(entry *tui.TextEdit) {
		model.GotoLine = entry.Text()
	})
	box.Append(gotoInput)
	model.CurrentNavMode = model.GotoNavigationMode
}

func Paginate(x []string, skip, size int) []string {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}

// TODO: move this to the UI package
func NewInputCommandEntry() *tui.Entry {
	inputCommand := tui.NewEntry()
	inputCommand.SetFocused(true)
	inputCommand.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputCommand.SetEchoMode(tui.EchoModeNormal)

	inputCommandBox := tui.NewHBox(inputCommand)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommand
}

// TODO: move this to the UI package
func NewInputCommandBox(input *tui.Entry) *tui.Box {
	inputCommandBox := tui.NewHBox(input)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommandBox
}

func GetPercentage(currentPosition int, fileContent *[]string) float64 {
	percent := float64(currentPosition) * 100.00
	percent = percent / float64(len(*fileContent))
	return percent
}

func CalculateTerminalHeight() int {
	advance := 45
	fd := int(os.Stdout.Fd())
	_, height, err := term.GetSize(fd)
	if err == nil {
		// 5 is for borders, input bar, status, etc. Adjust if needed.
		advance = height - 5
	}

	return advance
}
