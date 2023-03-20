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
