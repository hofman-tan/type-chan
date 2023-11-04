package app

// wrongState controls the behaviour of the typing page,
// when the user has made any mistakes in typing.
type wrongState struct {
	typingPage *typingPage
}

func (s *wrongState) handleLetter(l string) {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.text.notErrorCountLimitReached() {
		// update word holder
		s.typingPage.pushWordHolder(l)

		// update textarea
		s.typingPage.text.incrementErrorCount()
	}
}

func (s *wrongState) handleSpace() {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.text.notErrorCountLimitReached() {
		// update word holder
		s.typingPage.pushWordHolder(" ")

		// update textarea
		s.typingPage.text.incrementErrorCount()
	}
}

func (s *wrongState) handleBackspace() {
	// update word holder
	poppedLetter := s.typingPage.popWordHolder()

	// update textarea
	if poppedLetter != "" {
		if s.typingPage.text.hasError() {
			s.typingPage.text.decrementErrorCount()
		} else {
			s.typingPage.text.previousLetter()
		}
	}

	if !s.typingPage.text.hasError() {
		s.typingPage.changeState(s.typingPage.correctState)
	}
}

func (s *wrongState) handleEnter() {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.text.notErrorCountLimitReached() {
		// update word holder
		s.typingPage.pushWordHolder("⏎")

		// update textarea
		s.typingPage.text.incrementErrorCount()
	}
}

// newWrongState initialises and returns a new instance of wrongState
func newWrongState(t *typingPage) *wrongState {
	return &wrongState{typingPage: t}
}
