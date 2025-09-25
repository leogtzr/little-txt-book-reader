package terminal

import (
	"os"
	"os/exec"

	"golang.org/x/term"
)

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
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
