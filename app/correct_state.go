package app

// correctState controls the behaviour of the typing page,
// when the user has made no mistakes in typing.
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

// newCorrectState initialises and returns a new instance of correctState
func newCorrectState(t *typingPage) *correctState {
	return &correctState{typingPage: t}
}
