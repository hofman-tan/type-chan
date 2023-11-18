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

// typingPage is the model for the typing test page.
type typingPage struct {
	app          *app
	quoteFetcher *quoteFetcher
	started      bool

	totalKeysPressed   int
	correctKeysPressed int

	progressBar progress.Model
	textarea    *textarea
	wordInput   string
	stopWatch   stopwatch
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
		// fill up the buffer first
		for i := 0; i < quoteBufferSize; i++ {
			q, err := getRandomQuote()
			if err != nil {
				return err
			}
			quotes = append(quotes, q)
		}
		t.quoteFetcher.start(quoteBufferSize)
	}

	for _, quote := range quotes {
		t.textarea.append(quote)
	}
	return nil
}

// pushWordInput appends a letter to the word input.
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

// changeState changes the current state to the given value.
func (t *typingPage) changeState(s State) {
	t.currentState = s
}

// incrementKeysPressed increments the total number of keys pressed.
// 'correct' tells whether the current keypress is correct.
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

		if currentMode == Timed && len(t.textarea.lines) < scrollTextHeight {
			t.textarea.append(<-t.quoteFetcher.quotes)
		}

		if t.textarea.hasReachedEndOfText() {
			t.toResultPage()
		}

	case TickMsg:
		cmds = append(cmds, t.stopWatch.tick())

	case timer.TimeoutMsg:
		t.toResultPage()

	case tea.WindowSizeMsg:
		t.progressBar.Width = appWidth
		t.textarea.resize()
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
		progressPercent = t.textarea.currentProgress()
	case Timed:
		progressPercent = float64(Timeout-t.timer.Timeout) / float64(Timeout)
	}
	if t.textarea.anyMistyped() {
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
		t.textarea.View() + "\n\n" +
		strings.Repeat(" ", paddingX) + lipgloss.JoinHorizontal(lipgloss.Top, wordInput, timeStr) + "\n" +
		strings.Repeat(" ", paddingX) + lipgloss.NewStyle().Foreground(grey).Render("esc or ctrl+c to quit")

}

// toResultPage initialises and directs user to the result page.
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

// newTypingPage returns a new instance of typingPage.
func newTypingPage(app *app) *typingPage {
	t := &typingPage{app: app}
	t.correctState = newCorrectState(t)
	t.wrongState = newWrongState(t)
	t.currentState = t.correctState // initially at correct state

	t.textarea = newTextarea()
	t.progressBar = progress.New(progress.WithWidth(appWidth), progress.WithoutPercentage())

	switch currentMode {
	case Sprint:
		t.stopWatch = newStopwatch()
	case Timed:
		t.textarea.scroll = true
		t.timer = timer.NewWithInterval(Timeout, time.Second)
	}

	t.quoteFetcher = newQuoteFetcher(context.Background())
	return t
}
