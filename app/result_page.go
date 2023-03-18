package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// resultPage is the page model for displaying typing test results.
type resultPage struct {
	app *app

	totalKeysPressed   int
	correctKeysPressed int
	elapsedTime        time.Duration

	grossWPM    float64
	accuracy    float64
	adjustedWPM float64
	cpm         float64
}

func (r *resultPage) init() tea.Cmd {
	// https://support.sunburst.com/hc/en-us/articles/229335208-Type-to-Learn-How-are-Words-Per-Minute-and-Accuracy-Calculated-
	r.grossWPM = (float64(r.totalKeysPressed) / 5) / r.elapsedTime.Minutes()
	r.accuracy = (float64(r.correctKeysPressed) / float64(r.totalKeysPressed)) // range 0 to 1
	r.adjustedWPM = r.grossWPM * r.accuracy
	r.cpm = float64(r.totalKeysPressed) / r.elapsedTime.Minutes()

	return nil
}

func (r *resultPage) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			// exit
			return tea.Quit
		} else if msg.Type == tea.KeyEnter {
			// switch to typing page
			r.app.changePage(newTypingPage(r.app))
			return r.app.Init()
		}
	}
	return nil
}

func (r *resultPage) view() string {
	statStr := fmt.Sprintf("Gross WPM: %.2f\n", r.grossWPM)
	statStr += fmt.Sprintf("Accuracy: %.2f%%\n", r.accuracy*100)
	statStr += fmt.Sprintf("Adjusted WPM: %.2f\n\n", r.adjustedWPM)

	statStr += fmt.Sprintf("Time: %v\n", r.elapsedTime.Round(10*time.Millisecond))
	statStr += fmt.Sprintf("CPM: %.2f\n\n", r.cpm)

	statStr += fmt.Sprintf("Total keys pressed: %d\n", r.totalKeysPressed)
	statStr += fmt.Sprintf("Correct keys: %d", r.correctKeysPressed)

	str := borderStyle.Render(statStr)
	str += "\npress enter to start new game"
	str += "\npress esc or ctrl+c to quit\n"
	return str
}

// newResultPage initialises and returns a new instance of resultPage.
func newResultPage(
	app *app,
	totalKeysPressed int,
	correctKeysPressed int,
	elapsedTime time.Duration,
) *resultPage {

	return &resultPage{
		app:                app,
		totalKeysPressed:   totalKeysPressed,
		correctKeysPressed: correctKeysPressed,
		elapsedTime:        elapsedTime,
	}
}
