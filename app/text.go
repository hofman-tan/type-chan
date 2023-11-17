package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// text is the model for textarea content, and all relevant position indices.
type text struct {
	lines       []string
	totalLength int
	totalTyped  int

	scroll bool // make textarea scroll (current line appears on top)

	currentLineIndex     int // index position of current line of text
	currentLetterIndex   int // index position of letter, counted from the start of current line
	letterIndexFromStart int // index position of letter, counted from the start of text
	mistypedCount        int // number of wrongly-typed letters, counted from the current letter index
}

// newText initialises and returns a new instance of text.
func newText() *text {
	return &text{
		lines: []string{},
	}
}

// append appends new quote to the text model.
func (t *text) append(q quote) {
	// adds a newline character to the end of text
	if len(t.lines) != 0 {
		t.lines[len(t.lines)-1] += "\n"
		t.totalLength++
	}

	quoteLines := splitTextIntoLines(q.Text)
	t.lines = append(t.lines, quoteLines...)
	t.totalLength += q.length
}

// currentLine returns the current line of text where the cursor lies.
func (t *text) currentLine() string {
	return t.lines[t.currentLineIndex]
}

// nextLetter moves the cursor to the next letter.
func (t *text) nextLetter() {
	t.currentLetterIndex++
	t.letterIndexFromStart++
	if t.currentLetterIndex >= len(t.currentLine()) {
		if t.scroll {
			// remove the previous line from slice, so current line appears on top
			_, trimmed := t.lines[0], t.lines[1:]
			t.lines = trimmed
			t.letterIndexFromStart = 0
		} else {
			// move to next line
			t.currentLineIndex++
		}
		t.currentLetterIndex = 0
	}
	t.totalTyped++
}

// previousLetter moves the cursor to the previous letter.
func (t *text) previousLetter() {
	// ignore if cursor is at the start of the line, or if previous letter is a whitespace
	if t.currentLetterIndex == 0 || string(t.currentLine()[t.currentLetterIndex-1]) == " " {
		return
	}
	t.currentLetterIndex--
	t.letterIndexFromStart--
	t.totalTyped--
}

// incrementMistypedCount increments the number of mistypes made.
func (t *text) incrementMistypedCount() {
	t.mistypedCount++
}

// decrementMistypedCount decrements the number of mistypes made.
func (t *text) decrementMistypedCount() {
	if t.anyMistyped() {
		t.mistypedCount--
	}
}

// canIncrementMistyped tells if mistyped count can still be incremented further.
func (t *text) canIncrementMistyped() bool {
	return t.mistypedCount < t.remainingLettersCount() && t.mistypedCount < maxMistypedCount
}

// anyMistyped tells if there's any mistypes made.
func (t *text) anyMistyped() bool {
	return t.mistypedCount > 0
}

// remainingLetterCount returns the number of letters left to type,
// excluding the current letter.
func (t *text) remainingLettersCount() int {
	return t.totalLength - t.totalTyped
}

// hasReachedEndOfText tells if the cursor has moved beyond the whole text,
// denoting the completion of the typing test.
func (t *text) hasReachedEndOfText() bool {
	return t.remainingLettersCount() <= 0
}

// currentLetter returns the letter currently pointed by the cursor.
func (t *text) currentLetter() string {
	return string(t.lines[t.currentLineIndex][t.currentLetterIndex])
}

// currentProgress returns the current progress of typing.
// The returned value ranges from 0 (just started) to 1 (completed).
func (t *text) currentProgress() float64 {
	return float64(t.totalTyped) / float64(t.totalLength)
}

func (t *text) View() string {
	result := ""
	MistypesToRender := 0
	lineIndex := 0

	for lineIndex < len(t.lines) {
		// ignore lines that are not visible in scroll mode
		if t.scroll && lineIndex >= scrollTextHeight {
			break
		}

		result += strings.Repeat(" ", paddingX)

		for letterIndex, letter := range t.lines[lineIndex] {
			letterStr := string(letter)
			if letter == '\n' {
				letterStr = "⏎"
			}

			if lineIndex < t.currentLineIndex ||
				(lineIndex == t.currentLineIndex && letterIndex < t.currentLetterIndex) {
				// typed letters
				letterStr = lipgloss.NewStyle().Foreground(grey).Render(letterStr)
			}

			if lineIndex == t.currentLineIndex && letterIndex == t.currentLetterIndex {
				// current (untyped) letter
				letterStr = lipgloss.NewStyle().Underline(true).Render(letterStr)
				MistypesToRender = t.mistypedCount
			}

			if MistypesToRender > 0 {
				// mistyped letters
				letterStr = lipgloss.NewStyle().Background(red).Render(letterStr)
				MistypesToRender--
			}

			// no styling applied for untyped letters that come after current letter

			result += letterStr
		}

		result += "\n"
		lineIndex++
	}
	return result
}

func splitTextIntoLines(text string) []string {
	result := []string{}

	if len(text) == 0 {
		return result
	}

	// preserve trailling whitespace or newline after each word
	wordsSlice := []string{}
	buf := []rune{}
	for _, letter := range text {
		buf = append(buf, letter)
		if letter == ' ' || letter == '\n' {
			wordsSlice = append(wordsSlice, string(buf))
			buf = []rune{}
		}
	}
	if len(buf) > 0 {
		wordsSlice = append(wordsSlice, string(buf))
	}

	line := []string{}
	lineLen := 0
	for _, word := range wordsSlice {
		if lineLen != 0 && lineLen+len(word) > windowWidth-paddingX*2 {
			result = append(result, strings.Join(line, ""))
			line = []string{}
			lineLen = 0
		}
		line = append(line, word)
		lineLen += len(word)

		if word[len(word)-1] == '\n' {
			result = append(result, strings.Join(line, ""))
			line = []string{}
			lineLen = 0
		}
	}
	if lineLen != 0 {
		result = append(result, strings.Join(line, ""))
	}
	return result
}

func (t *text) resize() {
	t.lines = splitTextIntoLines(strings.Join(t.lines, ""))

	// determine new values for the letter and line indices
	accLen := 0
	for lineIndex, line := range t.lines {
		if accLen+len(line) >= t.letterIndexFromStart {
			t.currentLetterIndex = t.letterIndexFromStart - accLen
			t.currentLineIndex = lineIndex
			break
		}
		accLen += len(line)
	}
}
