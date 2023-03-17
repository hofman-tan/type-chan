package app

import tea "github.com/charmbracelet/bubbletea"

type Page interface {
	init() tea.Cmd
	update(tea.Msg) tea.Cmd
	view() string
}
