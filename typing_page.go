package main

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

const maxErrorOffset int = 10
const quoteBufferSize int = 3

var countdown = 5 * 60 // 5 minutes (for timed test) TODO: convert to argument

type Mode int

const (
	Sprint Mode = iota
	Timed
)

type typingPage struct {
	app *app

	textLines  []string
	textLength int
	words      []string
	wordHolder string

	currentTextIndex   int // index position of current un-typed letter, counted from the very beginning of the whole text
	currentLineIndex   int // index position of current line of text
	currentLetterIndex int // index position of letter, counted from the start of current line
	currentWordIndex   int // index position of current word from the word slices
	errorCount         int // number of wrongly-typed letters, counted from the current letter index

	totalKeysPressed   int
	correctKeysPressed int

	currentState State
	mode         Mode

	quoteFetcher *QuoteFetcher
	timer        Timer
	started      bool
	viewBuilder  *typingPageViewBuilder
}

func (t *typingPage) Init() tea.Cmd {
	//text := "test"
	//text := "hello there how are you my friend?"
	//text := "During the first part of your life, you only become aware of happiness once you have lost it. Then an age comes, a second one, in which you already know, at the moment when you begin to experience true happiness, that you are, at the end of the day, going to lose it. When I met Belle, I understood that I had just entered this second age. I also understood that I hadn't reached the third age, in which anticipation of the loss of happiness prevents you from living."
	//text := "‘Margareta! I’m surprised at you! We both know there’s no such thing as love!’"
	//text := "hey»\nthere"

	quotes := []Quote{}
	if t.mode == Timed {
		// fill up the buffer
		for i := 0; i < quoteBufferSize; i++ {
			quotes = append(quotes, getRandomQuote())
		}
		// begin endless fetching
		t.quoteFetcher.Start(quoteBufferSize)

	} else {
		quotes = append(quotes, getRandomQuote())
	}

	// TODO: create a separate struct for text, which will initialize given a slice of quotes
	for _, quote := range quotes {
		if len(t.textLines) != 0 {
			t.textLines[len(t.textLines)-1] += "\n"
			t.textLength++
		}

		t.textLines = append(t.textLines, quote.lines...)
		t.textLength += quote.length

		t.words = append(t.words, quote.words...)
	}

	return nil
}

func (t *typingPage) pushWordHolder(l string) {
	if t.errorCount < t.remainingLettersCount() && t.errorCount < maxErrorOffset {
		t.wordHolder += l
	}
}

func (t *typingPage) popWordHolder() string {
	word := []rune(t.wordHolder)

	if len(word) > 0 {
		lastLetter := word[len(word)-1]
		word = word[:len(word)-1] // remove the last letter
		t.wordHolder = string(word)
		return string(lastLetter)
	}
	return ""
}

func (t *typingPage) clearWordHolder() {
	t.wordHolder = ""
}

func (t *typingPage) currentLine() string {
	return t.textLines[t.currentLineIndex]
}

func (t *typingPage) nextLetter() {
	t.currentTextIndex++

	t.currentLetterIndex++
	if t.currentLetterIndex >= len(t.currentLine()) {
		t.currentLineIndex++
		t.currentLetterIndex = 0
	}
}

func (t *typingPage) previousLetter() {
	t.currentTextIndex--

	t.currentLetterIndex--
	if t.currentLetterIndex < 0 {
		t.currentLineIndex--
		t.currentLetterIndex = len(t.currentLine()) - 1
	}
}

func (t *typingPage) nextWord() {
	t.currentWordIndex++
}

func (t *typingPage) incrementErrorOffset() {
	if t.errorCount < t.remainingLettersCount() && t.errorCount < maxErrorOffset {
		t.errorCount++
	}
}

func (t *typingPage) decrementErrorOffset() {
	if t.errorCount != 0 {
		t.errorCount--
	}
}

func (t *typingPage) changeState(s State) {
	t.currentState = s
}

func (t *typingPage) remainingLettersCount() int {
	return t.textLength - t.currentTextIndex
}

func (t *typingPage) isEndOfTextReached() bool {
	return t.remainingLettersCount() <= 0
}

func (t *typingPage) currentLetter() string {
	return string(t.textLines[t.currentLineIndex][t.currentLetterIndex])
}

// range within 0 to 1
func (t *typingPage) currentProgress() float64 {
	return float64(t.currentTextIndex) / float64(t.textLength)
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
		} else if msg.Type == tea.KeyEnter {
			t.currentState.handleEnter()
		} else if msg.Type == tea.KeyTab || msg.Type == tea.KeyUp || msg.Type == tea.KeyDown || msg.Type == tea.KeyLeft || msg.Type == tea.KeyRight {
			// do nothing
		} else {
			t.currentState.handleLetter(msg.String())
		}

		if !t.started {
			t.started = true
			return t.timer.tick()
		}

		if t.mode == Timed && t.remainingLettersCount() < 10 {
			t.textLines[len(t.textLines)-1] += "\n"
			t.textLength++

			quote := <-t.quoteFetcher.quotes
			t.textLines = append(t.textLines, quote.lines...)
			t.textLength += quote.length
		}

		// done typing the whole text
		if t.isEndOfTextReached() {
			// switch to result page
			t.quoteFetcher.Stop()
			resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, t.timer.getTimeElapsed())
			t.app.changePage(resultPage)
			return t.app.Init()
		}

	case TickMsg:
		return t.timer.tick()

	case TimesUpMsg:
		// time's up!
		t.quoteFetcher.cancel()
		resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, t.timer.getTimeElapsed())
		t.app.changePage(resultPage)
		return t.app.Init()
	}

	return nil
}

func (t *typingPage) View() string {
	if t.isEndOfTextReached() {
		return ""
	}

	if t.mode == Timed {
		timeProgress := t.timer.getTimeElapsed().Seconds() / float64(countdown)
		t.viewBuilder.addProgressBar(timeProgress)

		// make current line the first line in textarea
		lineCount := int(math.Min(textareaHeight, float64(len(t.textLines)-t.currentLineIndex)))
		lines := t.textLines[t.currentLineIndex : t.currentLineIndex+lineCount]
		t.viewBuilder.addTextarea(lines, 0, t.currentLetterIndex, t.errorCount)

	} else {

		// render full text, with progress bar
		t.viewBuilder.addProgressBar(t.currentProgress())
		t.viewBuilder.addTextarea(t.textLines, t.currentLineIndex, t.currentLetterIndex, t.errorCount)
	}

	t.viewBuilder.addSidebar(t.started, t.timer.string())
	t.viewBuilder.addWordHolder(t.wordHolder)
	return t.viewBuilder.render()
}

func newTypingPage(app *app) *typingPage {
	typingPage := &typingPage{app: app}
	// initially at correct state
	typingPage.currentState = &CorrectState{typingPage: typingPage}
	typingPage.mode = Timed //Sprint

	typingPage.textLines = []string{}
	typingPage.textLength = 0
	typingPage.words = []string{}

	if typingPage.mode == Timed {
		typingPage.timer = newCountDownTimer(countdown)
	} else {
		typingPage.timer = newCountUpTimer()
	}

	typingPage.quoteFetcher = newQuoteFetcher()
	typingPage.viewBuilder = &typingPageViewBuilder{}
	return typingPage
}
