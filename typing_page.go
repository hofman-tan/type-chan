package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxErrorOffset int = 10
const progressBarLength int = 80

type typingPage struct {
	app *app

	text       string
	words      []string
	wordHolder string

	currentTextIndex int // points to the currently un-typed letter
	currentWordIndex int // points to the currently un-typed word
	errorOffset      int // the index offset starting from currentTextIndex, marks the number of wrongly-typed letters

	startTime   time.Time
	elapsedTime time.Duration

	totalKeysPressed   int
	correctKeysPressed int

	currentState State

	timer *timer
}

func (t *typingPage) Init() tea.Cmd {
	//text := "test"
	//text := "hello there how are you my friend?"
	//text := "During the first part of your life, you only become aware of happiness once you have lost it. Then an age comes, a second one, in which you already know, at the moment when you begin to experience true happiness, that you are, at the end of the day, going to lose it. When I met Belle, I understood that I had just entered this second age. I also understood that I hadn't reached the third age, in which anticipation of the loss of happiness prevents you from living."
	//text := "‘Margareta! I’m surprised at you! We both know there’s no such thing as love!’"
	//text := "hey»\nthere"

	text := getRandomQuote()
	words := strings.Split(text, " ")

	t.text = text
	t.words = words

	return nil
}

func (t *typingPage) pushWordHolder(l string) {
	if t.errorOffset < t.remainingLettersCount() && t.errorOffset < maxErrorOffset {
		t.wordHolder += l
	}
}

func (t *typingPage) popWordHolder() string {
	if len(t.wordHolder) > 0 {
		lastLetter := t.wordHolder[len(t.wordHolder)-1]
		t.wordHolder = t.wordHolder[:len(t.wordHolder)-1] // remove last letter from word holder
		return string(lastLetter)
	}
	return ""
}

func (t *typingPage) clearWordHolder() {
	t.wordHolder = ""
}

func (t *typingPage) nextLetter() {
	t.currentTextIndex++
}

func (t *typingPage) previousLetter() {
	t.currentTextIndex--
}

func (t *typingPage) nextWord() {
	t.currentWordIndex++
}

func (t *typingPage) incrementErrorOffset() {
	if t.errorOffset < t.remainingLettersCount() && t.errorOffset < maxErrorOffset {
		t.errorOffset++
	}
}

func (t *typingPage) decrementErrorOffset() {
	if t.errorOffset != 0 {
		t.errorOffset--
	}
}

func (t *typingPage) changeState(s State) {
	t.currentState = s
}

func (t *typingPage) remainingLettersCount() int {
	return len(t.text) - t.currentTextIndex
}

func (t *typingPage) isEndOfTextReached() bool {
	return t.currentTextIndex >= len(t.text)
}

func (t *typingPage) currentLetter() string {
	return string(t.text[t.currentTextIndex])
}

// range within 0 to 1
func (t *typingPage) currentProgress() float64 {
	return float64(t.currentTextIndex) / float64(len(t.text))
}

func (t *typingPage) markStartTime() time.Time {
	time := time.Now()
	t.startTime = time
	return time
}

func (t *typingPage) markElapsedTime() time.Duration {
	duration := time.Since(t.startTime)
	t.elapsedTime = duration
	return duration
}

func (t *typingPage) incrementKeysPressed(correct bool) {
	t.totalKeysPressed++
	if correct {
		t.correctKeysPressed++
	}
}

func (t *typingPage) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			// exit
			return tea.Quit
		} else if msg.Type == tea.KeyBackspace {
			t.currentState.handleBackspace()
		} else if msg.Type == tea.KeySpace {
			t.currentState.handleSpace()
		} else if msg.Type == tea.KeyTab || msg.Type == tea.KeyEnter || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown || msg.Type == tea.KeyLeft || msg.Type == tea.KeyRight {
			// do nothing
		} else {
			t.currentState.handleLetter(msg.String())
		}

		if t.startTime.IsZero() {
			t.markStartTime()
			return t.timer.tick()
		}

		// done typing the whole text
		if t.isEndOfTextReached() {
			t.markElapsedTime()
			// switch to result page
			resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, t.elapsedTime)
			t.app.changePage(resultPage)
			return t.app.Init()
		}

	case TickMsg:
		return t.timer.tick()

	case TimesUpMsg:
		// time's up!
		t.markElapsedTime()
		resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, t.elapsedTime)
		t.app.changePage(resultPage)
		return t.app.Init()
	}

	return nil
}

func (t *typingPage) progressBarView() string {
	times := int(math.Floor(t.currentProgress() * float64(progressBarLength)))
	bar := progressBarContentStyle.Render(strings.Repeat("_", times))
	blank := progressBarBlankStyle.Render(strings.Repeat("_", progressBarLength-times))
	return bar + blank
}

func (t *typingPage) sidebarView() string {
	if t.startTime.IsZero() {
		return sidebarStyle.Render("Start typing!")
	} else {
		return sidebarStyle.Render(fmt.Sprintf("Time left:\n%s", t.timer.string()))
	}
}

func (t *typingPage) wordHolderView() string {
	str := wordHolderStyle.Render(t.wordHolder)
	str += "\npress esc or ctrl+c to quit\n"
	return str
}

func (t *typingPage) View() string {
	if t.isEndOfTextReached() {
		return ""
	}

	str := ""
	str += t.progressBarView() + "\n"
	str += lipgloss.JoinHorizontal(lipgloss.Top, t.currentState.textareaView(), t.sidebarView()) + "\n"
	str += t.wordHolderView()
	return ContainerStyle.Render(str)
}

func newTypingPage(app *app) *typingPage {
	typingPage := &typingPage{app: app}
	// initially at correct state
	typingPage.currentState = &CorrectState{typingPage: typingPage}
	typingPage.timer = newTimer()
	return typingPage
}
