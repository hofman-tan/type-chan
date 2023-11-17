package app

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	windowWidth  int
	appWidth     int
	initializing bool
)

// app represents the main typing test program.
// It keeps track of the page the user is currently on.
type app struct {
	currentPage Page
}

func (a *app) Init() tea.Cmd {
	return tea.ClearScreen
}

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		if msg.Width < minWindowWidth {
			windowWidth = minWindowWidth
		} else {
			windowWidth = msg.Width
		}
		appWidth = windowWidth - paddingX*2
		initializing = false
	}
	if a.currentPage == nil {
		return a, nil
	}
	return a, a.currentPage.update(msg)
}

func (a *app) View() string {
	if initializing {
		return "Initializing..."
	}

	return strings.Repeat("\n", paddingY) +
		a.currentPage.view() +
		strings.Repeat("\n", paddingY)
}

// changePage sets the current page to the given value.
func (a *app) changePage(page Page) {
	a.currentPage = page
	a.currentPage.init()
}

// New returns a new app instance.
func New() *app {
	return &app{}
}

// Start launches the program with the given mode.
func (a *app) Start(m Mode) {
	currentMode = m
	initializing = true

	// switch to typing page
	a.changePage(newTypingPage(a))

	p := tea.NewProgram(a)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting the program: %v", err)
		os.Exit(1)
	}
}
