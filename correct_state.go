package main

import "github.com/muesli/reflow/wordwrap"

type CorrectState struct {
	typingPage *typingPage
}

func (s *CorrectState) handleLetter(l string) {
	// update word holder
	s.typingPage.pushWordHolder(l)

	// update textarea
	if l == s.typingPage.currentLetter() {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.nextLetter()
	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.incrementErrorOffset()
		s.typingPage.changeState(newWrongState(s.typingPage))
	}
}

func (s *CorrectState) handleSpace() {
	// update word holder
	s.typingPage.pushWordHolder(" ")

	if s.typingPage.currentLetter() == " " {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		// clear word holder
		s.typingPage.clearWordHolder()
		// update textarea
		s.typingPage.nextWord()
		s.typingPage.nextLetter()

	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.incrementErrorOffset()
		s.typingPage.changeState(newWrongState(s.typingPage))
	}
}

func (s *CorrectState) handleBackspace() {
	// update word holder
	poppedLetter := s.typingPage.popWordHolder()

	// update textarea
	if poppedLetter != "" {
		s.typingPage.previousLetter()
	}
}

func (s *CorrectState) textareaView() string {
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
			str += currentLetterStyle.Render(letter)
		} else {
			str += letter
		}

		if rune == '\n' {
			str += "\n"
		}
	}

	return greenTextAreaStyle.Render(str)
}

func newCorrectState(t *typingPage) *CorrectState {
	return &CorrectState{typingPage: t}
}
