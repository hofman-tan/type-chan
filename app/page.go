package app

import tea "github.com/charmbracelet/bubbletea"

// Page is the interface for all concrete page models.
type Page interface {
	// init handles the initialisation of a page i.e. when page transition occurs.
	init()

	// update updates the internal states of the page.
	update(tea.Msg) tea.Cmd

	// view returns the TUI string for the page.
	view() string
}
