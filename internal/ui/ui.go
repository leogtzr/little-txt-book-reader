package ui

import (
	"textreader/internal/model"

	"github.com/marcusolsson/tui-go"
)

func NewInputCommandEntry() *tui.Entry {
	inputCommand := tui.NewEntry()
	inputCommand.SetFocused(true)
	inputCommand.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputCommand.SetEchoMode(tui.EchoModeNormal)

	return inputCommand
}

func NewInputCommandBox(input *tui.Entry) *tui.Box {
	inputCommandBox := tui.NewHBox(input)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommandBox
}

func AddGotoWidget(box *tui.Box, state *model.AppState) {
	gotoInput := tui.NewTextEdit()
	gotoInput.SetText("Go To line: ")
	gotoInput.SetFocused(true)
	gotoInput.OnTextChanged(func(entry *tui.TextEdit) {
		state.GotoLine = entry.Text()
	})
	box.Append(gotoInput)
	state.CurrentNavMode = model.GotoNavigationMode
}
