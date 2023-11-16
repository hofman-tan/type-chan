package app

import (
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

	wordHolder string
	started    bool

	totalKeysPressed   int
	correctKeysPressed int

	progressBar progress.Model
	stopWatch   stopWatch
	timer       timer.Model

	currentState State
	correctState *correctState
	wrongState   *wrongState
}

func (t *typingPage) init() {
	//text := "test"
	//text := "hello there how are you my friend?"
	//text := "During the first part of your life, you only become aware of happiness once you have lost it. Then an age comes, a second one, in which you already know, at the moment when you begin to experience true happiness, that you are, at the end of the day, going to lose it. When I met Belle, I understood that I had just entered this second age. I also understood that I hadn't reached the third age, in which anticipation of the loss of happiness prevents you from living."
	//text := "‘Margareta! I’m surprised at you! We both know there’s no such thing as love!’"
	//text := "hey»\nthere"

	quotes := []quote{}
	switch currentMode {
	case Sprint:
		quotes = append(quotes, getRandomQuote())
	case Timed:
		// fill up the buffer
		for i := 0; i < quoteBufferSize; i++ {
			quotes = append(quotes, getRandomQuote())
		}
		// begin continuous querying
		t.quoteFetcher.start(quoteBufferSize)
	}

	for _, quote := range quotes {
		t.text.append(quote)
	}
}

// pushWordHolder appends the letter to the word holder.
func (t *typingPage) pushWordHolder(l string) {
	t.wordHolder += l
}

// popWordHolder removes the last letter from the word holder.
func (t *typingPage) popWordHolder() string {
	word := []rune(t.wordHolder)
	if len(word) == 0 {
		return ""
	}
	lastLetter := word[len(word)-1]
	word = word[:len(word)-1] // remove the last letter
	t.wordHolder = string(word)
	return string(lastLetter)
}

// clearWordHolder clears the word holder.
func (t *typingPage) clearWordHolder() {
	t.wordHolder = ""
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

func (t *typingPage) update(msg tea.Msg) tea.Cmd {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			// exit
			return tea.Quit
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

		if currentMode == Timed && len(t.text.lines) < scrollWindowHeight {
			t.text.append(<-t.quoteFetcher.quotes)
		}

		if t.text.isEndOfTextReached() {
			t.toResultPage()
		}

	case TickMsg:
		cmds = append(cmds, t.stopWatch.tick())

	// Time's up!
	case timer.TimeoutMsg:
		t.toResultPage()

	// Terminal window is resized
	case tea.WindowSizeMsg:
		t.progressBar.Width = windowWidth - paddingX*2
		t.text.resize()

	}

	if currentMode == Timed {
		var cmd tea.Cmd
		t.timer, cmd = t.timer.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (t *typingPage) view() string {
	if t.text.isEndOfTextReached() {
		return ""
	}

	appWidth := windowWidth - paddingX*2

	var progressPercent float64
	switch currentMode {
	case Sprint:
		progressPercent = t.text.currentProgress()
	case Timed:
		progressPercent = float64(Countdown-t.timer.Timeout) / float64(Countdown)
	}

	progressBar := t.progressBar.ViewAs(progressPercent)
	wordHolder := "> " + t.wordHolder
	wordHolder = lipgloss.NewStyle().Width(appWidth / 2).Align(lipgloss.Left).Render(wordHolder)

	var timeStr string
	if currentMode == Sprint {
		timeStr = t.stopWatch.view()
	} else {
		timeStr = t.timer.View()
	}
	timeStr = lipgloss.NewStyle().Width(appWidth / 2).Align(lipgloss.Right).Render(timeStr)

	return strings.Repeat(" ", paddingX) + progressBar + "\n\n" +
		t.text.View() + "\n\n" +
		strings.Repeat(" ", paddingX) + lipgloss.JoinHorizontal(lipgloss.Top, wordHolder, timeStr) + "\n" +
		strings.Repeat(" ", paddingX) + lipgloss.NewStyle().Foreground(grey).Render("esc or ctrl+c to quit")

}

// toResultPage initialises and transitions to result page.
func (t *typingPage) toResultPage() {
	t.quoteFetcher.stop()

	var elapsed time.Duration
	if currentMode == Sprint {
		elapsed = t.stopWatch.elapsed()
	} else {
		elapsed = Countdown
	}
	resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, elapsed)
	t.app.changePage(resultPage)
}

// newTypingPage initialises and returns a new instance of typingPage.
func newTypingPage(app *app) *typingPage {
	t := &typingPage{app: app}
	t.correctState = newCorrectState(t)
	t.wrongState = newWrongState(t)
	t.currentState = t.correctState // initially at correct state

	t.text = newText()
	t.progressBar = progress.New(progress.WithSolidFill(string(green)), progress.WithWidth(windowWidth-paddingX*2), progress.WithoutPercentage())

	switch currentMode {
	case Sprint:
		t.stopWatch = newStopwatch()
	case Timed:
		t.text.scroll = true
		t.timer = timer.NewWithInterval(Countdown, time.Second)
	}

	t.quoteFetcher = newQuoteFetcher()
	return t
}
