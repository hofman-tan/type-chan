package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app := &app{}
	// switch to typing page
	app.changePage(newTypingPage(app))

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting the program: %v", err)
		os.Exit(1)
	}
}
