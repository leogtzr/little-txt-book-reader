package model

import (
	"time"

	"github.com/marcusolsson/tui-go"
)

// NavMode represents the navigation mode of the application.
type NavMode int

// AppState holds the application state.
type AppState struct {
	From, To, FromForReferences, ToReferences, FromForVocabulary, ToForVocabulary int
	GotoLine                                                                      string
	FileToOpen                                                                    string
	Vocabulary                                                                    []string
	PercentagePointStats, ToggleShowStatus                                        bool
	References, FileContent, BannedWords                                          []string
	CurrentNavMode                                                                NavMode
	Sidebar                                                                       *tui.Box
	RefsTable, VocabTable                                                         *tui.Table
	PageIndex, CurrentPercentage, Advance                                         int
	MinutesToReachNextPercentagePoint                                             map[int]time.Duration
	StartTime                                                                     time.Time
	CurrentHighlight, CurrentWord                                                 int
}

// NewAppState initializes a new AppState instance.
func NewAppState() *AppState {
	return &AppState{
		From:                              0,
		To:                                0, // Will be set based on Advance
		FromForReferences:                 0,
		ToReferences:                      10,
		GotoLine:                          "",
		FileToOpen:                        "", // Initialize as empty string
		PercentagePointStats:              false,
		ToggleShowStatus:                  true,
		References:                        []string{},
		FileContent:                       []string{},
		BannedWords:                       []string{},
		CurrentNavMode:                    ReadingNavigationMode,
		Vocabulary:                        []string{},
		Sidebar:                           tui.NewVBox(),
		RefsTable:                         tui.NewTable(0, 0),
		VocabTable:                        tui.NewTable(0, 0),
		PageIndex:                         0,
		MinutesToReachNextPercentagePoint: make(map[int]time.Duration),
		CurrentPercentage:                 0,
		Advance:                           0,
		CurrentHighlight:                  0,
		CurrentWord:                       0,
	}
}

// LatestFile represents the latest file state.
type LatestFile struct {
	FileName string
	From     int
	To       int
}

const (
	DownKeyBindingAlternative1                       = "j"
	DownKeyBindingAlternative2                       = "Down"
	UpKeyBindingAlternative1                         = "k"
	UpKeyBindingAlternative2                         = "Up"
	GotoKeyBindingAlternative1                       = "g"
	NewNoteKeyBindingAlternative1                    = "n"
	ShowStatusKeyBinding                             = "."
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
	SaveVocabularyKeyBinding                         = "w"
	ShowVocabularyKeyBinding                         = "W" // New key binding for vocabulary dialog

)

const (
	ReadingNavigationMode                    NavMode = 1
	ShowReferencesNavigationMode             NavMode = 2
	AnalyzeAndFilterReferencesNavigationMode NavMode = 3
	GotoNavigationMode                       NavMode = 4
	ShowTimePercentagePointsMode             NavMode = 5
	ShowHelpMode                             NavMode = 6
	VocabularyNavigationMode                 NavMode = 7 // New navigation mode

	GotoWidgetIndex = 2

	NonRefsFileName = "non-refs.txt"

	PageSize = 20

	DBFileRequiredNumberFields = 3
)
