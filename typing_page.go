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

	timer        Timer
	text         *text
	quoteFetcher *QuoteFetcher
	viewBuilder  *typingPageViewBuilder

	wordHolder string
	started    bool

	totalKeysPressed   int
	correctKeysPressed int

	currentState State
	mode         Mode
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

	for _, quote := range quotes {
		t.text.append(quote)
	}

	return nil
}

func (t *typingPage) pushWordHolder(l string) {
	if t.text.errorCount < t.text.remainingLettersCount() &&
		t.text.errorCount < maxErrorOffset {
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

func (t *typingPage) changeState(s State) {
	t.currentState = s
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

		if t.mode == Timed && t.text.remainingLettersCount() < 10 {
			t.text.append(<-t.quoteFetcher.quotes)
		}

		// done typing the whole text
		if t.text.isEndOfTextReached() {
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
	if t.text.isEndOfTextReached() {
		return ""
	}

	if t.mode == Timed {
		timeProgress := t.timer.getTimeElapsed().Seconds() / float64(countdown)
		t.viewBuilder.addProgressBar(timeProgress)

		// make current line the first line in textarea
		lineCount := int(math.Min(textareaHeight, float64(len(t.text.textLines)-t.text.currentLineIndex)))
		lines := t.text.textLines[t.text.currentLineIndex : t.text.currentLineIndex+lineCount]
		t.viewBuilder.addTextarea(lines, 0, t.text.currentLetterIndex, t.text.errorCount)

	} else {

		// render full text, with progress bar
		t.viewBuilder.addProgressBar(t.text.currentProgress())
		t.viewBuilder.addTextarea(t.text.textLines, t.text.currentLineIndex, t.text.currentLetterIndex, t.text.errorCount)
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

	typingPage.text = newText()

	if typingPage.mode == Timed {
		typingPage.timer = newCountDownTimer(countdown)
	} else {
		typingPage.timer = newCountUpTimer()
	}

	typingPage.quoteFetcher = newQuoteFetcher()
	typingPage.viewBuilder = &typingPageViewBuilder{}
	return typingPage
}
