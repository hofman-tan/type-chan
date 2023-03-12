package main

import (
	"github.com/muesli/reflow/wordwrap"
)

type WrongState struct {
	typingPage *typingPage
}

func (s *WrongState) handleLetter(l string) {
	s.typingPage.incrementKeysPressed(false)

	// update word holder
	s.typingPage.pushWordHolder(l)

	// update textarea
	s.typingPage.incrementErrorOffset()
}

func (s *WrongState) handleSpace() {
	s.typingPage.incrementKeysPressed(false)

	// update word holder
	s.typingPage.pushWordHolder(" ")

	// update textarea
	s.typingPage.incrementErrorOffset()
}

func (s *WrongState) handleBackspace() {
	// update word holder
	poppedLetter := s.typingPage.popWordHolder()

	// update textarea
	if poppedLetter != "" {
		if s.typingPage.errorOffset != 0 {
			s.typingPage.decrementErrorOffset()
		} else {
			s.typingPage.previousLetter()
		}
	}

	if s.typingPage.errorOffset == 0 {
		s.typingPage.changeState(newCorrectState(s.typingPage))
	}
}

func (s *WrongState) textareaView() string {
	str := ""
	test := wordwrap.String(s.typingPage.text, textareaWidth)

	for index, rune := range test {
		letter := string(rune)
		if rune == '\n' {
			letter = " "
		}

		if index < s.typingPage.currentTextIndex {
			str += pastTextStyle.Render(letter)
		} else if index == s.typingPage.currentTextIndex {
			str += errorOffsetStyle.Render(currentLetterStyle.Render(letter))
		} else if index > s.typingPage.currentTextIndex && index < s.typingPage.currentTextIndex+s.typingPage.errorOffset {
			str += errorOffsetStyle.Render(letter)
		} else {
			str += letter
		}

		if rune == '\n' {
			str += "\n"
		}
	}

	return redTextAreaStyle.Render(str)
}

func newWrongState(t *typingPage) *WrongState {
	return &WrongState{typingPage: t}
}
