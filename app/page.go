package app

import tea "github.com/charmbracelet/bubbletea"

// Page is the interface for all concrete page models.
type Page interface {
	// init is called when initialising a page.
	init() tea.Cmd

	// update updates the internal states of the page.
	update(tea.Msg) tea.Cmd

	// view returns the TUI string for the page.
	view() string
}
