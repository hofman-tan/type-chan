package app

type State interface {
	handleLetter(string)
	handleSpace()
	handleBackspace()
	handleEnter()
}
