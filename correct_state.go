package main

type CorrectState struct {
	model *model
}

func (s *CorrectState) handleLetter(l string) {
	// update word holder
	s.model.pushWordHolder(l)

	// update textarea
	if l == s.model.currentLetter() {
		// correct letter
		s.model.nextLetter()
	} else {
		// wrong letter
		s.model.incrementErrorOffset()
		s.model.changeState(s.model.wrongState)
	}
}

func (s *CorrectState) handleSpace() {
	// update word holder
	s.model.pushWordHolder(" ")

	if s.model.currentLetter() == " " {
		// correct letter
		// clear word holder
		s.model.clearWordHolder()
		// update textarea
		s.model.nextWord()
		s.model.nextLetter()

	} else {
		// wrong letter
		s.model.incrementErrorOffset()
		s.model.changeState(s.model.wrongState)
	}
}

func (s *CorrectState) handleBackspace() {
	// update word holder
	poppedLetter := s.model.popWordHolder()

	// update textarea
	if poppedLetter != "" {
		s.model.previousLetter()
	}
}

func (s *CorrectState) view() string {
	str := ""

	// textarea
	past := pastTextStyle.Render(s.model.pastText())
	current := currentLetterStyle.Render(s.model.currentLetter())
	future := s.model.futureText()
	str += greenTextAreaStyle.Render(past + current + future)
	str += "\n"

	// word holder
	str += wordHolderStyle.Render(s.model.wordHolder)
	str += "\npress esc or ctrl+c to quit\n"

	return str
}
