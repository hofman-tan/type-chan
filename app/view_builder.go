package app

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// typingPageViewBuilder is a view builder that constructs the typing page TUI.
type typingPageViewBuilder struct {
	withProgressBar            bool
	progressBarCurrentProgress float64

	withTextarea               bool
	textareaLines              []string
	textareaCurrentLetterIndex int
	textareaCurrentLineIndex   int
	textareaErrorCount         int
	textareaScroll             bool

	withSidebar bool
	sidebarStr  string

	withWordHolder bool
	wordHolder     string
}

// newViewBuilder returns a new instance of typingPageViewBuilder.
func newViewBuilder() *typingPageViewBuilder {
	return &typingPageViewBuilder{}
}

// addProgressBar configures the builder to render the progress bar in the final TUI.
func (t *typingPageViewBuilder) addProgressBar(currentProgress float64) {
	t.withProgressBar = true
	t.progressBarCurrentProgress = currentProgress
}

// renderProgressBar returns the TUI string representation of the progress bar.
func (t *typingPageViewBuilder) renderProgressBar(totalProgress int) string {
	times := int(math.Floor(t.progressBarCurrentProgress * float64(totalProgress)))
	bar := whiteTextStyle.Render(strings.Repeat("_", times))
	blank := greyTextStyle.Render(strings.Repeat("_", totalProgress-times))

	return bar + blank
}

// addTextarea configures the builder to render the text area in the final TUI.
func (t *typingPageViewBuilder) addTextarea(lines []string, currentLineIndex int, currentLetterIndex int, errorCount int) {
	t.withTextarea = true
	t.textareaLines = lines
	t.textareaCurrentLineIndex = currentLineIndex
	t.textareaCurrentLetterIndex = currentLetterIndex
	t.textareaErrorCount = errorCount
}

// setTextareaScroll sets the text area to scroll mode, which fixes the
// current line to the top of text area.
func (t *typingPageViewBuilder) setTextareaScroll(scroll bool) {
	t.textareaScroll = scroll
}

// renderTextarea returns the TUI string representation of the textarea.
func (t *typingPageViewBuilder) renderTextarea() string {
	str := ""
	errorsToRender := 0
	linesRendered := 0

	lineIndex := 0
	if t.textareaScroll {
		// make current line the first line in textarea
		lineIndex = t.textareaCurrentLineIndex
	}

	for lineIndex < len(t.textareaLines) {
		if t.textareaScroll && linesRendered >= textareaMaxHeight {
			break
		}

		line := t.textareaLines[lineIndex]
		for letterIndex, rune := range line {
			letter := string(rune)
			if rune == '\n' {
				letter = "⏎"
			}

			if lineIndex < t.textareaCurrentLineIndex ||
				(lineIndex == t.textareaCurrentLineIndex && letterIndex < t.textareaCurrentLetterIndex) {
				// past letters
				letter = greyTextStyle.Render(letter)
			}

			if lineIndex == t.textareaCurrentLineIndex && letterIndex == t.textareaCurrentLetterIndex {
				// current letter
				letter = underlinedStyle.Render(letter)
				errorsToRender = t.textareaErrorCount
			}

			if errorsToRender > 0 {
				// wrong letters
				letter = redTextStyle.Render(letter)
				errorsToRender--
			}

			// no styling applied for future letters.

			str += letter

			// add more space between quotes
			if rune == '\n' {
				str += "\n"
			}
		}

		str += "\n"
		lineIndex++
		linesRendered++
	}

	textBoxStyle := lipgloss.NewStyle().
		Width(textareaWidth).
		MaxWidth(textareaWidth)

	if t.textareaScroll {
		// fixed height
		textBoxStyle = textBoxStyle.Height(textareaMaxHeight).MaxHeight(textareaMaxHeight)
	} else {
		// variable height
		textBoxStyle = textBoxStyle.Height(textareaMinHeight)
	}

	str = textBoxStyle.Render(str)

	if t.textareaErrorCount > 0 {
		return redBorderStyle.Render(str)
	} else {
		return greenBorderStyle.Render(str)
	}
}

// addSidebar configures the builder to render the sidebar in the final TUI.
func (t *typingPageViewBuilder) addSidebar(started bool, time string) {
	t.withSidebar = true
	if !started {
		t.sidebarStr = "Start typing!"
	} else {
		t.sidebarStr = fmt.Sprintf("Time:\n%s", time)
	}
}

// renderSidebar returns the TUI string representation of the sidebar.
func (t *typingPageViewBuilder) renderSidebar(height int) string {
	return sidebarStyle.Height(height).Render(t.sidebarStr)
}

// addWordHolder configures the builder to render the word holder in the final TUI.
func (t *typingPageViewBuilder) addWordHolder(word string) {
	t.withWordHolder = true
	t.wordHolder = wordHolderStyle.Render(word)
}

// renderWordHolder returns the TUI string representation of the word holder.
func (t *typingPageViewBuilder) renderWordHolder() string {
	return t.wordHolder
}

// render constructs and view and returns the final TUI string.
func (t *typingPageViewBuilder) render() string {
	textarea := ""
	if t.withTextarea {
		textarea = t.renderTextarea()
	}

	sidebar := ""
	if t.withSidebar {
		// sidebar follows the same height as textarea
		sidebar = t.renderSidebar(lipgloss.Height(textarea) - 2) // minus off the padding
	}

	textareaSidebar := lipgloss.JoinHorizontal(lipgloss.Top, textarea, sidebar)

	progressBar := ""
	if t.withProgressBar {
		// progress bar follows the same width as textarea + sidebar
		progressBar = t.renderProgressBar(lipgloss.Width(textareaSidebar))
	}

	wordHolder := ""
	if t.withWordHolder {
		wordHolder = t.renderWordHolder()
	}

	str := ""
	str += progressBar + "\n"
	str += textareaSidebar + "\n"
	str += wordHolder + "\n"
	str += "press esc or ctrl+c to quit\n"
	return str
}
