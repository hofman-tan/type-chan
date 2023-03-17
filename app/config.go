package app

type Mode int

const (
	Sprint Mode = iota
	Timed
)

var currentMode Mode

// Countdown (in seconds) for timed mode
var Countdown = 5 * 60
