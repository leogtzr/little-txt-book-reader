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
func GetStatusInformation(state *model.AppState) string {
	if !state.ToggleShowStatus {
		return ""
	}

	percent := progress.GetPercentage(state.To, &state.FileContent)
	if int(percent) > state.CurrentPercentage {
		state.CurrentPercentage = int(percent)
		now := time.Now()
		state.MinutesToReachNextPercentagePoint[int(percent)] = now.Sub(state.StartTime)
		state.StartTime = now
	}

	if state.PercentagePointStats {
		return fmt.Sprintf(".   %d of %d lines (%.3f%%) [%d lines To next percentage point]                    ",
			state.To,
			len(state.FileContent), percent, progress.LinesToChangePercentagePoint(state.To, len(state.FileContent)))
	}
	return fmt.Sprintf(".   %d of %d lines (%.3f%%)                                            ",
		state.To, len(state.FileContent), percent)

}

func GetSavedStatusInformation(fileName string, state *model.AppState) string {
	return fmt.Sprintf(`%s <saved "%s">`, GetStatusInformation(state), fileName)
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
