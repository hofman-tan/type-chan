package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// typingPage is the page model for the typing test program.
type typingPage struct {
	app *app

	timer        Timer
	text         *text
	quoteFetcher *quoteFetcher
	viewBuilder  *typingPageViewBuilder

	wordHolder string
	started    bool

	totalKeysPressed   int
	correctKeysPressed int

	currentState State
}

func (t *typingPage) init() tea.Cmd {
	//text := "test"
	//text := "hello there how are you my friend?"
	//text := "During the first part of your life, you only become aware of happiness once you have lost it. Then an age comes, a second one, in which you already know, at the moment when you begin to experience true happiness, that you are, at the end of the day, going to lose it. When I met Belle, I understood that I had just entered this second age. I also understood that I hadn't reached the third age, in which anticipation of the loss of happiness prevents you from living."
	//text := "‘Margareta! I’m surprised at you! We both know there’s no such thing as love!’"
	//text := "hey»\nthere"

	quotes := []quote{}
	if currentMode == Timed {
		// fill up the buffer
		for i := 0; i < quoteBufferSize; i++ {
			quotes = append(quotes, getRandomQuote())
		}
		// begin endless fetching
		t.quoteFetcher.start(quoteBufferSize)

	} else {
		quotes = append(quotes, getRandomQuote())
	}

	for _, quote := range quotes {
		t.text.append(quote)
	}

	return nil
}

// pushWordHolder appends the letter to the word holder.
func (t *typingPage) pushWordHolder(l string) {
	t.wordHolder += l
}

// popWordHolder removes the last letter from the word holder.
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

		if currentMode == Timed && t.text.remainingLettersCount() < (textCountThreshold) {
			t.text.append(<-t.quoteFetcher.quotes)
		}

		// done typing the whole text
		if t.text.isEndOfTextReached() {
			return t.toResultPage()
		}

	case TickMsg:
		return t.timer.tick()

	case TimesUpMsg:
		// time's up!
		return t.toResultPage()
	}

	return nil
}

func (t *typingPage) view() string {
	if t.text.isEndOfTextReached() {
		return ""
	}

	t.viewBuilder.addTextarea(t.text.textLines, t.text.currentLineIndex, t.text.currentLetterIndex, t.text.errorCount)

	if currentMode == Timed {
		// show elapsed time as current progress
		timeProgress := t.timer.getTimeElapsed().Seconds() / float64(Countdown)
		t.viewBuilder.addProgressBar(timeProgress)
		// make textarea scrolls (current line appears on top)
		t.viewBuilder.setTextareaScroll(true)
	} else {
		t.viewBuilder.addProgressBar(t.text.currentProgress())
	}

	t.viewBuilder.addSidebar(t.started, t.timer.string())
	t.viewBuilder.addWordHolder(t.wordHolder)
	return t.viewBuilder.render()
}

// toResultPage initialises and transitions to result page.
func (t *typingPage) toResultPage() tea.Cmd {
	t.quoteFetcher.stop()
	resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, t.timer.getTimeElapsed())
	t.app.changePage(resultPage)
	return t.app.Init()
}

// newTypingPage initialises and returns a new instance of typingPage.
func newTypingPage(app *app) *typingPage {
	typingPage := &typingPage{app: app}
	// initially at correct state
	typingPage.currentState = newCorrectState(typingPage)
	typingPage.text = newText()

	if currentMode == Timed {
		typingPage.timer = newCountDownTimer(Countdown)
	} else {
		typingPage.timer = newCountUpTimer()
	}

	typingPage.quoteFetcher = newQuoteFetcher()
	typingPage.viewBuilder = newViewBuilder()
	return typingPage
}
