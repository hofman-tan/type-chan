package app

import tea "github.com/charmbracelet/bubbletea"

// Page is the interface for page models.
type Page interface {
	// init handles the initialisation of a page
	// i.e. when the page first starts.
	init() error

	// update updates the underlying model of the page.
	update(tea.Msg) (tea.Cmd, error)

	// view renders the page UI.
	view() string
}
