package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type typingPageViewBuilder struct {
	progressBarCurrentProgress float64
	textarea                   string
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

func (t *typingPageViewBuilder) addTextarea(ta string) {
	t.textarea = ta
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
	// sidebar follows the same height as textarea
	sidebar := t.renderSidebar(lipgloss.Height(t.textarea) - 2)
	textareaSidebar := lipgloss.JoinHorizontal(lipgloss.Top, t.textarea, sidebar)

	// progress bar follows the same width as textarea + sidebar
	progressBar := t.renderProgressBar(lipgloss.Width(textareaSidebar))

	str := ""
	str += progressBar + "\n"
	str += textareaSidebar + "\n"
	str += t.wordHolder + "\n"
	str += "press esc or ctrl+c to quit\n"
	return str
}
