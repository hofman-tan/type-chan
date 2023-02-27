package main

type WrongState struct {
	model *model
}

func (s *WrongState) handleLetter(l string) {
	// update word holder
	s.model.pushWordHolder(l)

	// update textarea
	s.model.incrementErrorOffset()
}

func (s *WrongState) handleSpace() {
	// update word holder
	s.model.pushWordHolder(" ")

	// update textarea
	s.model.incrementErrorOffset()
}

func (s *WrongState) handleBackspace() {
	// update word holder
	poppedLetter := s.model.popWordHolder()

	// update textarea
	if poppedLetter != "" {
		if s.model.errorOffset != 0 {
			s.model.decrementErrorOffset()
		} else {
			s.model.previousLetter()
		}
	}

	if s.model.errorOffset == 0 {
		s.model.changeState(s.model.correctState)
	}
}

func (s *WrongState) view() string {
	str := ""

	// textarea
	past := pastTextStyle.Render(s.model.pastText())
	errorOffset := errorOffsetStyle.Render(s.model.text[s.model.currentTextIndex : s.model.currentTextIndex+s.model.errorOffset])
	future := s.model.text[s.model.currentTextIndex+s.model.errorOffset:]
	str += redTextAreaStyle.Render(past + errorOffset + future)
	str += "\n"

	// wordholder
	str += wordHolderStyle.Render(s.model.wordHolder)
	str += "\npress esc or ctrl+c to quit\n"

	return str
}
