package app

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	windowWidth int
	appWidth    int
)

// app represents the main typing test program.
// It keeps track of the page the user is currently on.
type app struct {
	currentPage Page
	error       error
}

func (a *app) Init() tea.Cmd {
	// starts on typing page
	if err := a.changePage(newTypingPage(a)); err != nil {
		a.error = err
		return tea.Quit
	}
	return tea.ClearScreen
}

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// window is resized
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		if msg.Width < minWindowWidth {
			windowWidth = minWindowWidth
		} else {
			windowWidth = msg.Width
		}
		appWidth = windowWidth - paddingX*2
	}

	if a.currentPage == nil {
		return a, nil
	}
	cmd, err := a.currentPage.update(msg)
	if err != nil {
		a.error = err
		return a, tea.Quit
	}
	return a, cmd
}

func (a *app) View() string {
	if a.error != nil {
		return fmt.Sprintf("Something went wrong!\nError: %s\n", a.error)
	}

	return strings.Repeat("\n", paddingY) +
		a.currentPage.view() +
		strings.Repeat("\n", paddingY)
}

// changePage sets the current page to the given value.
func (a *app) changePage(page Page) error {
	a.currentPage = page
	return a.currentPage.init()
}

// New returns a new app instance.
func New() *app {
	return &app{}
}

// Start launches the program with the given mode.
func (a *app) Start(m Mode) {
	currentMode = m

	p := tea.NewProgram(a)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting the program: %v", err)
		os.Exit(1)
	}
}
