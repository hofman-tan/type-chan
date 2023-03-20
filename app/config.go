package app

// Mode is the user setting that determines the type of typing test.
type Mode int

const (
	Sprint Mode = iota
	Timed
)

// currentMode keeps track of the current mode setting.
var currentMode Mode

// Countdown is the time limit (in seconds) for Timed mode.
var Countdown = 5 * 60
