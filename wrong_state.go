package main

type WrongState struct {
	typingPage *typingPage
}

func (s *WrongState) handleLetter(l string) {
	s.typingPage.incrementKeysPressed(false)

	// update word holder
	s.typingPage.pushWordHolder(l)

	// update textarea
	s.typingPage.text.incrementErrorOffset()
}

func (s *WrongState) handleSpace() {
	s.typingPage.incrementKeysPressed(false)

	// update word holder
	s.typingPage.pushWordHolder(" ")

	// update textarea
	s.typingPage.text.incrementErrorOffset()
}

func (s *WrongState) handleBackspace() {
	// update word holder
	poppedLetter := s.typingPage.popWordHolder()

	// update textarea
	if poppedLetter != "" {
		if s.typingPage.text.hasError() {
			s.typingPage.text.decrementErrorOffset()
		} else {
			s.typingPage.text.previousLetter()
		}
	}

	if !s.typingPage.text.hasError() {
		s.typingPage.changeState(newCorrectState(s.typingPage))
	}
}

func (s *WrongState) handleEnter() {
	s.typingPage.incrementKeysPressed(false)

	// update word holder
	s.typingPage.pushWordHolder("⏎")

	// update textarea
	s.typingPage.text.incrementErrorOffset()
}

func newWrongState(t *typingPage) *WrongState {
	return &WrongState{typingPage: t}
}
