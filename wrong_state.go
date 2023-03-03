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

func (s *WrongState) view() string {
	str := ""

	// textarea
	past := pastTextStyle.Render(s.typingPage.pastText())
	errorOffset := errorOffsetStyle.Render(s.typingPage.text[s.typingPage.currentTextIndex : s.typingPage.currentTextIndex+s.typingPage.errorOffset])
	future := s.typingPage.text[s.typingPage.currentTextIndex+s.typingPage.errorOffset:]
	str += redTextAreaStyle.Render(past + errorOffset + future)
	str += "\n"

	// wordholder
	str += wordHolderStyle.Render(s.typingPage.wordHolder)
	str += "\npress esc or ctrl+c to quit\n"

	return str
}

func newWrongState(t *typingPage) *WrongState {
	return &WrongState{typingPage: t}
}
