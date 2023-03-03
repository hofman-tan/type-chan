package main

import tea "github.com/charmbracelet/bubbletea"

type Page interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
}
