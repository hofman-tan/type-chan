package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// textarea is the model for the typing area.
type textarea struct {
	lines       []string
	totalLength int
	totalTyped  int

	scroll bool // make textarea scroll (current line appears on top)

	currentLineIndex     int // index position of current line in text
	currentLetterIndex   int // index position of letter/cursor, counted from the start of current line
	letterIndexFromStart int // index position of letter/cursor, counted from the start of text
	mistypedCount        int // number of mistyped letters
}

// newTextarea returns a new instance of textarea.
func newTextarea() *textarea {
	return &textarea{
		lines: []string{},
	}
}

// append appends a quote to the textarea model.
func (t *textarea) append(q quote) {
	// adds a newline character to the end of text
	if len(t.lines) != 0 {
		t.lines[len(t.lines)-1] += "\n"
		t.totalLength++
	}

	quoteLines := splitTextIntoLines(q.Text)
	t.lines = append(t.lines, quoteLines...)
	t.totalLength += q.length
}

// currentLine returns the current line in textarea where the cursor lies.
func (t *textarea) currentLine() string {
	return t.lines[t.currentLineIndex]
}

// nextLetter moves the cursor to the next letter.
func (t *textarea) nextLetter() {
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
func (t *textarea) previousLetter() {
	// ignore if cursor is at the start of the line, or if previous letter is a whitespace
	if t.currentLetterIndex == 0 || string(t.currentLine()[t.currentLetterIndex-1]) == " " {
		return
	}
	t.currentLetterIndex--
	t.letterIndexFromStart--
	t.totalTyped--
}

// incrementMistypedCount increments the number of mistypes made.
func (t *textarea) incrementMistypedCount() {
	t.mistypedCount++
}

// decrementMistypedCount decrements the number of mistypes made.
func (t *textarea) decrementMistypedCount() {
	if t.anyMistyped() {
		t.mistypedCount--
	}
}

// canIncrementMistyped tells if mistyped count can still be incremented further.
func (t *textarea) canIncrementMistyped() bool {
	return t.mistypedCount < t.remainingLettersCount() && t.mistypedCount < maxMistypedCount
}

// anyMistyped tells if there's any mistypes made.
func (t *textarea) anyMistyped() bool {
	return t.mistypedCount > 0
}

// remainingLettersCount returns the number of letters left to type.
func (t *textarea) remainingLettersCount() int {
	return t.totalLength - t.totalTyped
}

// hasReachedEndOfText tells if the cursor has moved beyond the whole text,
// denoting the completion of the typing test.
func (t *textarea) hasReachedEndOfText() bool {
	return t.remainingLettersCount() <= 0
}

// currentLetter returns the letter currently pointed by the cursor.
func (t *textarea) currentLetter() string {
	return string(t.lines[t.currentLineIndex][t.currentLetterIndex])
}

// currentProgress returns the current progress of the test in percentage.
func (t *textarea) currentProgress() float64 {
	return float64(t.totalTyped) / float64(t.totalLength)
}

func (t *textarea) View() string {
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

// splitTextIntoLines splits a text string into lines, where the length
// of each line is bounded by the current appWidth.
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
		if lineLen != 0 && lineLen+len(word) > appWidth {
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

// resize resizes the textarea's width by re-splitting the text according
// to the current resized window.
func (t *textarea) resize() {
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
