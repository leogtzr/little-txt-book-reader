package ui

import "github.com/marcusolsson/tui-go"

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

func NewInputCommandBox(input *tui.Entry) *tui.Box {
	inputCommandBox := tui.NewHBox(input)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	return inputCommandBox
}
