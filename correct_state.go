package main

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

func (s *CorrectState) view() string {
	str := ""

	// textarea
	past := pastTextStyle.Render(s.typingPage.pastText())
	current := currentLetterStyle.Render(s.typingPage.currentLetter())
	future := s.typingPage.futureText()
	str += greenTextAreaStyle.Render(past + current + future)
	str += "\n"

	// word holder
	str += wordHolderStyle.Render(s.typingPage.wordHolder)
	str += "\npress esc or ctrl+c to quit\n"

	return str
}

func newCorrectState(t *typingPage) *CorrectState {
	return &CorrectState{typingPage: t}
}
