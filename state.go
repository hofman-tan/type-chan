package main

type State interface {
	handleLetter(l string)
	handleSpace()
	handleBackspace()
	view() string
}
