package app

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// typingPage is the page model for the typing test program.
type typingPage struct {
	app *app

	timer        Timer
	text         *text
	quoteFetcher *quoteFetcher

	wordHolder *wordHolder
	started    bool

	totalKeysPressed   int
	correctKeysPressed int

	progressBar *progressBar

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
}

// pushWordHolder appends the letter to the word holder.
func (t *typingPage) pushWordHolder(l string) {
	t.wordHolder.add(l)
}

// popWordHolder removes the last letter from the word holder.
func (t *typingPage) popWordHolder() string {
	return t.wordHolder.pop()
}

// clearWordHolder clears the word holder.
func (t *typingPage) clearWordHolder() {
	t.wordHolder.clear()
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
			cmds = append(cmds, t.timer.tick())
		}

		if currentMode == Timed && len(t.text.lines) < scrollWindowHeight {
			t.text.append(<-t.quoteFetcher.quotes)
		}

		// done typing the whole text
		if t.text.isEndOfTextReached() {
			t.toResultPage()
		}

		if currentMode == Timed {
			// show elapsed time as current progress
			timeProgress := t.timer.getTimeElapsed().Seconds() / float64(Countdown)
			cmds = append(cmds, t.progressBar.SetPercent(timeProgress))
		} else {
			cmds = append(cmds, t.progressBar.SetPercent(t.text.currentProgress()))
		}

	case TickMsg:
		cmds = append(cmds, t.timer.tick())

	// Time's up!
	case TimesUpMsg:
		t.toResultPage()

	// Terminal window is resized
	case tea.WindowSizeMsg:
		t.text.resize()

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := t.progressBar.Update(msg)
		t.progressBar.Model = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (t *typingPage) view() string {
	if t.text.isEndOfTextReached() {
		return ""
	}

	return t.progressBar.View() + "\n\n" +
		t.text.View() + "\n\n" +
		t.wordHolder.View()
}

// toResultPage initialises and transitions to result page.
func (t *typingPage) toResultPage() {
	t.quoteFetcher.stop()
	resultPage := newResultPage(t.app, t.totalKeysPressed, t.correctKeysPressed, t.timer.getTimeElapsed())
	t.app.changePage(resultPage)
}

// newTypingPage initialises and returns a new instance of typingPage.
func newTypingPage(app *app) *typingPage {
	typingPage := &typingPage{app: app}
	typingPage.correctState = newCorrectState(typingPage)
	typingPage.wrongState = newWrongState(typingPage)
	// initially at correct state
	typingPage.currentState = typingPage.correctState
	typingPage.text = newText()
	typingPage.progressBar = newProgressBar()
	typingPage.wordHolder = newWordHolder()

	if currentMode == Timed {
		typingPage.text.scroll = true
		typingPage.timer = newCountDownTimer(Countdown)
	} else {
		typingPage.text.scroll = false
		typingPage.timer = newCountUpTimer()
	}

	typingPage.quoteFetcher = newQuoteFetcher()
	return typingPage
}
