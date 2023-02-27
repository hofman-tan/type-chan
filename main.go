package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const maxErrorOffset int = 10
const progressBarLength int = 60

type model struct {
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

	correctState State
	wrongState   State
	state        State
}

func initialModel() *model {
	text := "hello there how are you my friend?"
	//text := "During the first part of your life, you only become aware of happiness once you have lost it. Then an age comes, a second one, in which you already know, at the moment when you begin to experience true happiness, that you are, at the end of the day, going to lose it. When I met Belle, I understood that I had just entered this second age. I also understood that I hadn't reached the third age, in which anticipation of the loss of happiness prevents you from living."
	words := strings.Split(text, " ")

	m := &model{
		text:  text,
		words: words,
	}
	m.correctState = &CorrectState{model: m}
	m.wrongState = &WrongState{model: m}
	m.state = m.correctState

	return m
}

func (m *model) pushWordHolder(l string) {
	if m.errorOffset < m.remainingLettersCount() && m.errorOffset < maxErrorOffset {
		m.wordHolder += l
	}
}

func (m *model) popWordHolder() string {
	if len(m.wordHolder) > 0 {
		lastLetter := m.wordHolder[len(m.wordHolder)-1]
		m.wordHolder = m.wordHolder[:len(m.wordHolder)-1] // remove last letter from word holder
		return string(lastLetter)
	}
	return ""
}

func (m *model) clearWordHolder() {
	m.wordHolder = ""
}

func (m *model) nextLetter() {
	m.currentTextIndex++
}

func (m *model) previousLetter() {
	m.currentTextIndex--
}

func (m *model) nextWord() {
	m.currentWordIndex++
}

func (m *model) incrementErrorOffset() {
	if m.errorOffset < m.remainingLettersCount() && m.errorOffset < maxErrorOffset {
		m.errorOffset++
	}
}

func (m *model) decrementErrorOffset() {
	if m.errorOffset != 0 {
		m.errorOffset--
	}
}

func (m *model) changeState(s State) {
	m.state = s
}

func (m *model) remainingLettersCount() int {
	return len(m.text) - m.currentTextIndex
}

func (m *model) currentLetter() string {
	return string(m.text[m.currentTextIndex])
}

func (m *model) isEndOfTextReached() bool {
	return m.currentTextIndex >= len(m.text)
}

func (m *model) pastText() string {
	return m.text[0:m.currentTextIndex]
}

func (m *model) futureText() string {
	if m.currentTextIndex+1 == len(m.text) {
		return ""
	}
	return m.text[m.currentTextIndex+1:]
}

// range within 0 to 1
func (m *model) currentProgress() float64 {
	return float64(m.currentTextIndex) / float64(len(m.text))
}

func (m *model) markStartTime() time.Time {
	time := time.Now()
	m.startTime = time
	return time
}

func (m *model) markElapsedTime() time.Duration {
	duration := time.Since(m.startTime)
	m.elapsedTime = duration
	return duration
}

func (m *model) incrementKeysPressed(correct bool) {
	m.totalKeysPressed++
	if correct {
		m.correctKeysPressed++
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.startTime.IsZero() {
			m.markStartTime()
		}

		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			// exit
			return m, tea.Quit
		} else if msg.Type == tea.KeyBackspace {
			m.state.handleBackspace()
		} else if msg.Type == tea.KeySpace {
			m.state.handleSpace()
		} else {
			m.state.handleLetter(msg.String())
		}

		// done typing the whole text
		if m.isEndOfTextReached() {
			m.markElapsedTime()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) progressBar() string {
	times := int(math.Floor(m.currentProgress() * float64(progressBarLength)))
	bar := progressBarContentStyle.Render(strings.Repeat(" ", times))
	blank := progressBarBlankStyle.Render(strings.Repeat(" ", progressBarLength-times))
	return bar + blank + "\n"
}

func (m *model) View() string {
	if m.isEndOfTextReached() {
		return ""
	}

	str := ""
	str += m.progressBar()
	str += m.state.view()
	return ContainerStyle.Render(str)
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting the program: %v", err)
		os.Exit(1)
	}

	fmt.Printf("start: %s, duration: %v\n", m.startTime, m.elapsedTime.Minutes())
	fmt.Printf("total keys: %v, correct: %v\n", m.totalKeysPressed, m.correctKeysPressed)
	fmt.Printf("gross: %.2f, acc: %.2f%%, awpm: %.2f\n", m.grossWPM(), m.accuracy()*100, m.adjustedWPM())
}
