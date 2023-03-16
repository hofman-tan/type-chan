package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const textareaWidth = 60
const textareaHeight = 10

type typingPageViewBuilder struct {
	withProgressBar            bool
	progressBarCurrentProgress float64

	withTextarea               bool
	textareaLines              []string
	textareaCurrentLetterIndex int
	textareaCurrentLineIndex   int
	textareaErrorCount         int

	withSidebar bool
	sidebarStr  string

	withWordHolder bool
	wordHolder     string
}

func (t *typingPageViewBuilder) addProgressBar(currentProgress float64) {
	t.withProgressBar = true
	t.progressBarCurrentProgress = currentProgress
}

func (t *typingPageViewBuilder) renderProgressBar(totalProgress int) string {
	times := int(math.Floor(t.progressBarCurrentProgress * float64(totalProgress)))
	bar := progressBarContentStyle.Render(strings.Repeat("_", times))
	blank := progressBarBlankStyle.Render(strings.Repeat("_", totalProgress-times))

	return bar + blank
}

func (t *typingPageViewBuilder) addTextarea(lines []string, currentLineIndex int, currentLetterIndex int, errorCount int) {
	t.withTextarea = true
	t.textareaLines = lines
	t.textareaCurrentLineIndex = currentLineIndex
	t.textareaCurrentLetterIndex = currentLetterIndex
	t.textareaErrorCount = errorCount
}

func (t *typingPageViewBuilder) renderTextarea() string {
	str := ""
	errorsToRender := 0

	for lineIndex, line := range t.textareaLines {
		for letterIndex, rune := range line {
			letter := string(rune)
			if rune == '\n' {
				letter = "⏎"
			}

			if lineIndex < t.textareaCurrentLineIndex ||
				(lineIndex == t.textareaCurrentLineIndex && letterIndex < t.textareaCurrentLetterIndex) {
				// past letters
				letter = pastTextStyle.Render(letter)
			}

			if lineIndex == t.textareaCurrentLineIndex && letterIndex == t.textareaCurrentLetterIndex {
				// current letter
				letter = currentLetterStyle.Render(letter)
				errorsToRender = t.textareaErrorCount
			}

			if errorsToRender > 0 {
				// wrong letters
				letter = errorOffsetStyle.Render(letter)
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
	}

	str = lipgloss.NewStyle().Width(textareaWidth).Render(str)
	if t.textareaErrorCount > 0 {
		return redTextAreaStyle.Render(str)
	} else {
		return greenTextAreaStyle.Render(str)
	}
}

func (t *typingPageViewBuilder) addSidebar(started bool, time string) {
	t.withSidebar = true
	if !started {
		t.sidebarStr = "Start typing!"
	} else {
		t.sidebarStr = fmt.Sprintf("Time:\n%s", time)
	}
}

func (t *typingPageViewBuilder) renderSidebar(height int) string {
	return sidebarStyle.Height(height).Render(t.sidebarStr)
}

func (t *typingPageViewBuilder) addWordHolder(word string) {
	t.withWordHolder = true
	t.wordHolder = wordHolderStyle.Render(word)
}

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
		wordHolder = t.wordHolder
	}

	str := ""
	str += progressBar + "\n"
	str += textareaSidebar + "\n"
	str += wordHolder + "\n"
	str += "press esc or ctrl+c to quit\n"
	return str
}
