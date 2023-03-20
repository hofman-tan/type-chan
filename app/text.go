package app

// text is the model for textarea content, and all relevant position indices.
type text struct {
	textLines  []string
	textLength int
	words      []string

	currentTextIndex   int // index position of current un-typed letter, counted from the very beginning of the whole text
	currentLineIndex   int // index position of current line of text
	currentLetterIndex int // index position of letter, counted from the start of current line
	currentWordIndex   int // index position of current word from the word slices
	errorCount         int // number of wrongly-typed letters, counted from the current letter index
}

// newText initialises and returns a new instance of text.
func newText() *text {
	return &text{
		textLines: []string{},
		words:     []string{},
	}
}

// append appends new quote to the text model.
func (t *text) append(q quote) {
	if len(t.textLines) != 0 {
		t.textLines[len(t.textLines)-1] += "\n"
		t.textLength++
	}

	t.textLines = append(t.textLines, q.lines...)
	t.textLength += q.length
	t.words = append(t.words, q.words...)
}

// currentLine returns the current line of text where the cursor lies.
func (t *text) currentLine() string {
	return t.textLines[t.currentLineIndex]
}

// nextLetter moves the cursor to the next letter.
func (t *text) nextLetter() {
	t.currentTextIndex++

	t.currentLetterIndex++
	if t.currentLetterIndex >= len(t.currentLine()) {
		t.currentLineIndex++
		t.currentLetterIndex = 0
	}
}

// previousLetter moves the cursor to the previous letter.
func (t *text) previousLetter() {
	t.currentTextIndex--

	t.currentLetterIndex--
	if t.currentLetterIndex < 0 {
		t.currentLineIndex--
		t.currentLetterIndex = len(t.currentLine()) - 1
	}
}

// nextWord points the word cursor to the next word.
func (t *text) nextWord() {
	t.currentWordIndex++
}

// incrementErrorCount increments the number of errors made.
func (t *text) incrementErrorCount() {
	t.errorCount++
}

// decrementErrorCount decrements the number of errors made.
func (t *text) decrementErrorCount() {
	if t.hasError() {
		t.errorCount--
	}
}

// notErrorCountLimitReached tells if error count can still be incremented further.
func (t *text) notErrorCountLimitReached() bool {
	return t.errorCount < t.remainingLettersCount() && t.errorCount < maxErrorCount
}

// hasError tells if there's any errors made.
func (t *text) hasError() bool {
	return t.errorCount > 0
}

// remainingLetterCount returns the number of letters left to type,
// excluding the current letter.
func (t *text) remainingLettersCount() int {
	return t.textLength - t.currentTextIndex
}

// isEndOfTextReached tells if the cursor has moved beyond the whole text,
// denoting the completion of the typing test.
func (t *text) isEndOfTextReached() bool {
	return t.remainingLettersCount() <= 0
}

// currentLetter returns the letter currently pointed by the cursor.
func (t *text) currentLetter() string {
	return string(t.textLines[t.currentLineIndex][t.currentLetterIndex])
}

// currentProgress returns the current progress of typing.
// The returned value ranges from 0 (just started) to 1 (completed).
func (t *text) currentProgress() float64 {
	return float64(t.currentTextIndex) / float64(t.textLength)
}
