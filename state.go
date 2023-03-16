package main

type State interface {
	handleLetter(string)
	handleSpace()
	handleBackspace()
	handleEnter()
}
