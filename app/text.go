package app

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

func newText() *text {
	return &text{
		textLines: []string{},
		words:     []string{},
	}
}

func (t *text) append(q quote) {
	if len(t.textLines) != 0 {
		t.textLines[len(t.textLines)-1] += "\n"
		t.textLength++
	}

	t.textLines = append(t.textLines, q.lines...)
	t.textLength += q.length
	t.words = append(t.words, q.words...)
}

func (t *text) currentLine() string {
	return t.textLines[t.currentLineIndex]
}

func (t *text) nextLetter() {
	t.currentTextIndex++

	t.currentLetterIndex++
	if t.currentLetterIndex >= len(t.currentLine()) {
		t.currentLineIndex++
		t.currentLetterIndex = 0
	}
}

func (t *text) previousLetter() {
	t.currentTextIndex--

	t.currentLetterIndex--
	if t.currentLetterIndex < 0 {
		t.currentLineIndex--
		t.currentLetterIndex = len(t.currentLine()) - 1
	}
}

func (t *text) nextWord() {
	t.currentWordIndex++
}

func (t *text) incrementErrorOffset() {
	if t.errorCount < t.remainingLettersCount() && t.errorCount < maxErrorOffset {
		t.errorCount++
	}
}

func (t *text) hasError() bool {
	return t.errorCount > 0
}

func (t *text) decrementErrorOffset() {
	if t.hasError() {
		t.errorCount--
	}
}

func (t *text) remainingLettersCount() int {
	return t.textLength - t.currentTextIndex
}

func (t *text) isEndOfTextReached() bool {
	return t.remainingLettersCount() <= 0
}

func (t *text) currentLetter() string {
	return string(t.textLines[t.currentLineIndex][t.currentLetterIndex])
}

// range within 0 to 1
func (t *text) currentProgress() float64 {
	return float64(t.currentTextIndex) / float64(t.textLength)
}
