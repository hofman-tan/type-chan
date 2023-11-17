package app

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// typingPage is the page model for the typing test program.
type typingPage struct {
	app *app

	text         *text
	quoteFetcher *quoteFetcher

	wordInput string
	started   bool

	totalKeysPressed   int
	correctKeysPressed int

	progressBar progress.Model
	stopWatch   stopWatch
	timer       timer.Model

	currentState State
	correctState *correctState
	wrongState   *wrongState
}

func (t *typingPage) init() error {
	quotes := []quote{}
	switch currentMode {
	case Sprint:
		q, err := getRandomQuote()
		if err != nil {
			return err
		}
		quotes = append(quotes, q)

	case Timed:
		// fill up the buffer
		for i := 0; i < quoteBufferSize; i++ {
			q, err := getRandomQuote()
			if err != nil {
				return err
			}
			quotes = append(quotes, q)
		}
		// begin continuous querying
		t.quoteFetcher.start(quoteBufferSize)
	}

	for _, quote := range quotes {
		t.text.append(quote)
	}
	return nil
}

// pushWordInput appends the letter to the word input.
func (t *typingPage) pushWordInput(l string) {
	t.wordInput += l
}

// popWordInput removes the last letter from the word input.
func (t *typingPage) popWordInput() string {
	word := []rune(t.wordInput)
	if len(word) == 0 {
		return ""
	}
	lastLetter := word[len(word)-1]
	word = word[:len(word)-1] // remove the last letter
	t.wordInput = string(word)
	return string(lastLetter)
}

// clearWordInput clears the word input.
func (t *typingPage) clearWordInput() {
	t.wordInput = ""
}

// changeState sets the current state to the given value.
func (t *typingPage) changeState(s State) {
	t.currentState = s
}

// incrementKeysPressed increments the total number of keys pressed.
// 'correct' param denotes whether the current key pressed is correct.
func (t *typingPage) incrementKeysPressed(correct bool) {
	t.totalKeysPressed++
	if correct {
		t.correctKeysPressed++
	}
}

func (t *typingPage) update(msg tea.Msg) (tea.Cmd, error) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			// exit
			return tea.Quit, nil
		case tea.KeyBackspace:
			t.currentState.handleBackspace()
		case tea.KeySpace:
			t.currentState.handleSpace()
		case tea.KeyEnter:
			t.currentState.handleEnter()
		case tea.KeyTab, tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
			// do nothing
		default:
			t.currentState.handleLetter(msg.String())
		}

		if !t.started {
			t.started = true
			if currentMode == Sprint {
				cmds = append(cmds, t.stopWatch.start())
			}
			if currentMode == Timed {
				cmds = append(cmds, t.timer.Start())
			}
		}

		if currentMode == Timed && len(t.text.lines) < scrollTextHeight {
			t.text.append(<-t.quoteFetcher.quotes)
		}

		if t.text.hasReachedEndOfText() {
			t.toResultPage()
		}

	case TickMsg:
		cmds = append(cmds, t.stopWatch.tick())

	// Time's up!
	case timer.TimeoutMsg:
		t.toResultPage()

	// Terminal window is resized
	case tea.WindowSizeMsg:
		t.progressBar.Width = appWidth
		t.text.resize()

	}

	if currentMode == Timed {
		var cmd tea.Cmd
		t.timer, cmd = t.timer.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...), nil
}

func (t *typingPage) view() string {
	var progressPercent float64
	switch currentMode {
	case Sprint:
		progressPercent = t.text.currentProgress()
	case Timed:
		progressPercent = float64(Timeout-t.timer.Timeout) / float64(Timeout)
	}
	if t.text.anyMistyped() {
		t.progressBar.FullColor = string(red)
	} else {
		t.progressBar.FullColor = string(green)
	}

	progressBar := t.progressBar.ViewAs(progressPercent)
	wordInput := "> " + t.wordInput
	wordInput = lipgloss.NewStyle().Width(appWidth / 2).Align(lipgloss.Left).Render(wordInput)

	var timeStr string
	if currentMode == Sprint {
		timeStr = t.stopWatch.view()
	} else {
		timeStr = t.timer.View()
	}
	timeStr = lipgloss.NewStyle().Width(appWidth / 2).Align(lipgloss.Right).Render(timeStr)

	return strings.Repeat(" ", paddingX) + progressBar + "\n\n" +
		t.text.View() + "\n\n" +
		strings.Repeat(" ", paddingX) + lipgloss.JoinHorizontal(lipgloss.Top, wordInput, timeStr) + "\n" +
		strings.Repeat(" ", paddingX) + lipgloss.NewStyle().Foreground(grey).Render("esc or ctrl+c to quit")

}

// toResultPage initialises and transitions to result page.
func (t *typingPage) toResultPage() error {
	t.quoteFetcher.stop()

	var elapsed time.Duration
	if currentMode == Sprint {
		elapsed = t.stopWatch.elapsed()
	} else {
		elapsed = Timeout
	}
	resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, elapsed)
	return t.app.changePage(resultPage)
}

// newTypingPage initialises and returns a new instance of typingPage.
func newTypingPage(app *app) *typingPage {
	t := &typingPage{app: app}
	t.correctState = newCorrectState(t)
	t.wrongState = newWrongState(t)
	t.currentState = t.correctState // initially at correct state

	t.text = newText()
	t.progressBar = progress.New(progress.WithWidth(windowWidth-paddingX*2), progress.WithoutPercentage())

	switch currentMode {
	case Sprint:
		t.stopWatch = newStopwatch()
	case Timed:
		t.text.scroll = true
		t.timer = timer.NewWithInterval(Timeout, time.Second)
	}

	t.quoteFetcher = newQuoteFetcher(context.Background())
	return t
}
