package app

// State is the interface for all concrete states.
type State interface {
	// handleLetter handles alphanumerical keys from user.
	handleLetter(string)

	// handleSpace handles space key from user.
	handleSpace()

	// handleBackspace handles backspace key from user.
	handleBackspace()

	// handleEnter handles enter key from user.
	handleEnter()
}

// correctState controls the behaviour of the typing page,
// when the user has made no mistakes in typing.
type correctState struct {
	typingPage *typingPage
}

func (s *correctState) handleLetter(l string) {
	s.typingPage.pushWordInput(l)

	if l == s.typingPage.text.currentLetter() {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.text.nextLetter()
	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.text.incrementMistypedCount()
		s.typingPage.changeState(s.typingPage.wrongState)
	}
}

func (s *correctState) handleSpace() {
	s.typingPage.pushWordInput(" ")

	if s.typingPage.text.currentLetter() == " " {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.clearWordInput()
		s.typingPage.text.nextLetter()
	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.text.incrementMistypedCount()
		s.typingPage.changeState(s.typingPage.wrongState)
	}
}

func (s *correctState) handleBackspace() {
	poppedLetter := s.typingPage.popWordInput()
	if poppedLetter != "" {
		s.typingPage.text.previousLetter()
	}
}

func (s *correctState) handleEnter() {
	s.typingPage.pushWordInput("⏎")

	if s.typingPage.text.currentLetter() == "\n" {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.clearWordInput()
		s.typingPage.text.nextLetter()

	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.text.incrementMistypedCount()
		s.typingPage.changeState(s.typingPage.wrongState)
	}
}

// newCorrectState initialises and returns a new instance of correctState
func newCorrectState(t *typingPage) *correctState {
	return &correctState{typingPage: t}
}

// wrongState controls the behaviour of the typing page,
// when the user has made any mistakes in typing.
type wrongState struct {
	typingPage *typingPage
}

func (s *wrongState) handleLetter(l string) {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.text.canIncrementMistyped() {
		s.typingPage.pushWordInput(l)
		s.typingPage.text.incrementMistypedCount()
	}
}

func (s *wrongState) handleSpace() {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.text.canIncrementMistyped() {
		s.typingPage.pushWordInput(" ")
		s.typingPage.text.incrementMistypedCount()
	}
}

func (s *wrongState) handleBackspace() {
	poppedLetter := s.typingPage.popWordInput()

	if poppedLetter != "" {
		if s.typingPage.text.anyMistyped() {
			s.typingPage.text.decrementMistypedCount()
		} else {
			s.typingPage.text.previousLetter()
		}
	}

	if !s.typingPage.text.anyMistyped() {
		s.typingPage.changeState(s.typingPage.correctState)
	}
}

func (s *wrongState) handleEnter() {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.text.canIncrementMistyped() {
		s.typingPage.pushWordInput("⏎")
		s.typingPage.text.incrementMistypedCount()
	}
}

// newWrongState initialises and returns a new instance of wrongState
func newWrongState(t *typingPage) *wrongState {
	return &wrongState{typingPage: t}
}
