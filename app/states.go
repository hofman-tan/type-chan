package app

// State is the interface for all typingPage states.
type State interface {
	handleLetter(string) // handles alphanumerical keys
	handleSpace()
	handleBackspace()
	handleEnter()
}

// correctState handles the 'correct' behaviour of typingPage
// i.e. when there's no mistake in typing.
type correctState struct {
	typingPage *typingPage
}

// newCorrectState returns a new instance of correctState.
func newCorrectState(t *typingPage) *correctState {
	return &correctState{typingPage: t}
}

func (s *correctState) handleLetter(l string) {
	s.typingPage.pushWordInput(l)

	if l == s.typingPage.textarea.currentLetter() {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.textarea.nextLetter()
	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.textarea.incrementMistypedCount()
		s.typingPage.changeState(s.typingPage.wrongState)
	}
}

func (s *correctState) handleSpace() {
	s.typingPage.pushWordInput(" ")

	if s.typingPage.textarea.currentLetter() == " " {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.clearWordInput()
		s.typingPage.textarea.nextLetter()
	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.textarea.incrementMistypedCount()
		s.typingPage.changeState(s.typingPage.wrongState)
	}
}

func (s *correctState) handleBackspace() {
	poppedLetter := s.typingPage.popWordInput()
	if poppedLetter != "" {
		s.typingPage.textarea.previousLetter()
	}
}

func (s *correctState) handleEnter() {
	s.typingPage.pushWordInput("⏎")

	if s.typingPage.textarea.currentLetter() == "\n" {
		// correct letter
		s.typingPage.incrementKeysPressed(true)
		s.typingPage.clearWordInput()
		s.typingPage.textarea.nextLetter()

	} else {
		// wrong letter
		s.typingPage.incrementKeysPressed(false)
		s.typingPage.textarea.incrementMistypedCount()
		s.typingPage.changeState(s.typingPage.wrongState)
	}
}

// wrongState handles the 'wrong' behaviour of typingPage
// i.e. when there's any mistyped letter.
type wrongState struct {
	typingPage *typingPage
}

// newWrongState returns a new instance of wrongState.
func newWrongState(t *typingPage) *wrongState {
	return &wrongState{typingPage: t}
}

func (s *wrongState) handleLetter(l string) {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.textarea.canIncrementMistyped() {
		s.typingPage.pushWordInput(l)
		s.typingPage.textarea.incrementMistypedCount()
	}
}

func (s *wrongState) handleSpace() {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.textarea.canIncrementMistyped() {
		s.typingPage.pushWordInput(" ")
		s.typingPage.textarea.incrementMistypedCount()
	}
}

func (s *wrongState) handleBackspace() {
	poppedLetter := s.typingPage.popWordInput()

	if poppedLetter != "" {
		if s.typingPage.textarea.anyMistyped() {
			s.typingPage.textarea.decrementMistypedCount()
		} else {
			s.typingPage.textarea.previousLetter()
		}
	}

	if !s.typingPage.textarea.anyMistyped() {
		s.typingPage.changeState(s.typingPage.correctState)
	}
}

func (s *wrongState) handleEnter() {
	s.typingPage.incrementKeysPressed(false)

	if s.typingPage.textarea.canIncrementMistyped() {
		s.typingPage.pushWordInput("⏎")
		s.typingPage.textarea.incrementMistypedCount()
	}
}
