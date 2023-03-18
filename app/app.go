package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// app represents the main typing test program.
// It keeps track of the page the user is currently on.
type app struct {
	currentPage Page
}

func (a *app) Init() tea.Cmd {
	return a.currentPage.init()
}

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return a, a.currentPage.update(msg)
}

func (a *app) View() string {
	return containerStyle.Render(a.currentPage.view())
}

// changePage sets the current page to the given value.
func (a *app) changePage(page Page) {
	a.currentPage = page
}

// New returns a new app instance.
func New() *app {
	return &app{}
}

// Start launches the program with the given mode.
func (a *app) Start(m Mode) {
	currentMode = m

	// switch to typing page
	a.changePage(newTypingPage(a))

	p := tea.NewProgram(a)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting the program: %v", err)
		os.Exit(1)
	}
}
