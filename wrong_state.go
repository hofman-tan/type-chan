package main

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
		if s.typingPage.errorCount != 0 {
			s.typingPage.decrementErrorOffset()
		} else {
			s.typingPage.previousLetter()
		}
	}

	if s.typingPage.errorCount == 0 {
		s.typingPage.changeState(newCorrectState(s.typingPage))
	}
}

func (s *WrongState) handleEnter() {
	s.typingPage.incrementKeysPressed(false)

	// update word holder
	s.typingPage.pushWordHolder("⏎")

	// update textarea
	s.typingPage.incrementErrorOffset()
}

func newWrongState(t *typingPage) *WrongState {
	return &WrongState{typingPage: t}
}
