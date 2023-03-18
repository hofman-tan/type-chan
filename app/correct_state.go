package app

type correctState struct {
	typingPage *typingPage
}

func (s *correctState) handleLetter(l string) {
	// update word holder
	s.typingPage.pushWordHolder(l)

	// update textarea
	if l == s.typingPage.text.currentLetter() {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.text.nextLetter()
	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.text.incrementErrorCount()
		s.typingPage.changeState(newWrongState(s.typingPage))
	}
}

func (s *correctState) handleSpace() {
	// update word holder
	s.typingPage.pushWordHolder(" ")

	if s.typingPage.text.currentLetter() == " " {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		// clear word holder
		s.typingPage.clearWordHolder()
		// update textarea
		s.typingPage.text.nextWord()
		s.typingPage.text.nextLetter()

	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.text.incrementErrorCount()
		s.typingPage.changeState(newWrongState(s.typingPage))
	}
}

func (s *correctState) handleBackspace() {
	// update word holder
	poppedLetter := s.typingPage.popWordHolder()

	// update textarea
	if poppedLetter != "" {
		s.typingPage.text.previousLetter()
	}
}

func (s *correctState) handleEnter() { // update word holder
	s.typingPage.pushWordHolder("⏎")

	if s.typingPage.text.currentLetter() == "\n" {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		// clear word holder
		s.typingPage.clearWordHolder()
		// update textarea
		s.typingPage.text.nextWord()
		s.typingPage.text.nextLetter()

	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.text.incrementErrorCount()
		s.typingPage.changeState(newWrongState(s.typingPage))
	}
}

func newCorrectState(t *typingPage) *correctState {
	return &correctState{typingPage: t}
}
