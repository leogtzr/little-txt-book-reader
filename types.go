package main

import (
	"flag"
	"time"

	"github.com/marcusolsson/tui-go"
)

type navMode int

// LatestFile ...
type LatestFile struct {
	FileName string
	From     int
	To       int
}

var (
	From                              = 0
	To                                = Advance
	FromForReferences                 = 0
	ToReferences                      = 10
	GotoLine                          = ""
	FileToOpen                        = flag.String("file", "", "File To open")
	PercentagePointStats              = false
	ToggleShowStatus                  = true
	References                        = []string{}
	FileContent                       = []string{}
	CurrentNavMode                    = ReadingNavigationMode
	BannedWords                       = []string{}
	Sidebar                           = tui.NewVBox()
	RefsTable                         = tui.NewTable(0, 0)
	PageIndex                         = 0
	MinutesToReachNextPercentagePoint map[int]time.Duration
	StartTime                         time.Time
	CurrentPercentage                 int
	Advance                           int
)

const (
	DownKeyBindingAlternative1                       = "j"
	DownKeyBindingAlternative2                       = "Down"
	UpKeyBindingAlternative1                         = "k"
	UpKeyBindingAlternative2                         = "Up"
	GotoKeyBindingAlternative1                       = "g"
	NewNoteKeyBindingAlternative1                    = "n"
	showStatusKeyBinding                             = "."
	CloseGotoKeyBindingAlternative1                  = "r"
	SaveStatusKeyBindingAlternative1                 = "s"
	NextPercentagePointKeyBindingAlternative1        = "p"
	ShowReferencesKeyBindingAlternative1             = "f"
	CloseReferencesWindowKeyBindingAlternative1      = "q"
	CloseApplicationKeyBindingAlternative1           = "Esc"
	AnalyzeAndFilterReferencesKeyBinding             = "Alt+b"
	SaveQuoteKeyBindingAlternative1                  = "Alt+q"
	ShowMinutesTakenToReachPercentagePointKeyBinding = "m"
	OpenRAEWebSiteKeyBinging                         = "o"
	OpenGoodReadsWebSiteKeyBinding                   = "d"
	ShowHelpKeyBinding                               = "h"
	maxNumberOfElementsInGUIBox                      = 200
)

const (
	ReadingNavigationMode                    navMode = 1
	ShowReferencesNavigationMode             navMode = 2
	AnalyzeAndFilterReferencesNavigationMode navMode = 3
	GotoNavigationMode                       navMode = 4
	ShowTimePercentagePointsMode             navMode = 5
	ShowHelpMode                             navMode = 6

	GotoWidgetIndex = 2

	NonRefsFileName = "non-refs.txt"

	PageSize = 10

	DBFileRequiredNumbermields = 3
)
