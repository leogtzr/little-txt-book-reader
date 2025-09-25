package utils

import (
	"fmt"
	"os/exec"
	"textreader/internal/model"
	"textreader/internal/progress"
	"time"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
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

// TODO: check if we can move these two to a different package
func GetStatusInformation() string {
	if !model.ToggleShowStatus {
		return ""
	}

	percent := progress.GetPercentage(model.To, &model.FileContent)
	if int(percent) > model.CurrentPercentage {
		model.CurrentPercentage = int(percent)
		now := time.Now()
		model.MinutesToReachNextPercentagePoint[int(percent)] = now.Sub(model.StartTime)
		model.StartTime = now
	}

	if model.PercentagePointStats {
		return fmt.Sprintf(".   %d of %d lines (%.3f%%) [%d lines To next percentage point]                    ",
			model.To,
			len(model.FileContent), percent, progress.LinesToChangePercentagePoint(model.To, len(model.FileContent)))
	}
	return fmt.Sprintf(".   %d of %d lines (%.3f%%)                                            ",
		model.To, len(model.FileContent), percent)

}

func GetSavedStatusInformation(fileName string) string {
	return fmt.Sprintf(`%s <saved "%s">`, GetStatusInformation(), fileName)
}

func OpenOSEditor(os, notesFile string) *exec.Cmd {
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
