package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

const textareaWidth = 60

type typingPageViewBuilder struct {
	progressBarCurrentProgress float64
	textareaStr                string
	textareaCurrentIndex       int
	textareaErrorEndIndex      int
	sidebarStr                 string
	wordHolder                 string
}

func (t *typingPageViewBuilder) addProgressBar(currentProgress float64) {
	t.progressBarCurrentProgress = currentProgress
}

func (t *typingPageViewBuilder) renderProgressBar(totalProgress int) string {
	times := int(math.Floor(t.progressBarCurrentProgress * float64(totalProgress)))
	bar := progressBarContentStyle.Render(strings.Repeat("_", times))
	blank := progressBarBlankStyle.Render(strings.Repeat("_", totalProgress-times))

	return bar + blank
}

func (t *typingPageViewBuilder) addTextarea(ta string, currentIndex int, errorEndIndex int) {
	t.textareaStr = ta
	t.textareaCurrentIndex = currentIndex
	t.textareaErrorEndIndex = errorEndIndex
}

func (t *typingPageViewBuilder) renderTextarea() string {
	str := ""
	test := wordwrap.String(t.textareaStr, textareaWidth)

	for index, rune := range test {
		letter := string(rune)
		if rune == '\n' {
			letter = " "
		}

		if index < t.textareaCurrentIndex {
			// past letters
			str += pastTextStyle.Render(letter)

		} else if index == t.textareaCurrentIndex {
			// current letter
			l := currentLetterStyle.Render(letter)
			if t.textareaErrorEndIndex > t.textareaCurrentIndex {
				l = errorOffsetStyle.Render(l)
			}
			str += l

		} else if index > t.textareaCurrentIndex && index < t.textareaErrorEndIndex {
			// wrong letters
			str += errorOffsetStyle.Render(letter)

		} else {
			// future letters
			str += letter
		}

		if rune == '\n' {
			str += "\n"
		}
	}

	if t.textareaErrorEndIndex > t.textareaCurrentIndex {
		return redTextAreaStyle.Render(str)
	} else {
		return greenTextAreaStyle.Render(str)
	}
}

func (t *typingPageViewBuilder) addSidebar(started bool, time string) {
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
	t.wordHolder = wordHolderStyle.Render(word)
}

func (t *typingPageViewBuilder) render() string {
	textarea := t.renderTextarea()
	// sidebar follows the same height as textarea
	sidebar := t.renderSidebar(lipgloss.Height(textarea) - 2)
	textareaSidebar := lipgloss.JoinHorizontal(lipgloss.Top, textarea, sidebar)

	// progress bar follows the same width as textarea + sidebar
	progressBar := t.renderProgressBar(lipgloss.Width(textareaSidebar))

	str := ""
	str += progressBar + "\n"
	str += textareaSidebar + "\n"
	str += t.wordHolder + "\n"
	str += "press esc or ctrl+c to quit\n"
	return str
}
