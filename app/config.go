package app

import "time"

// currentMode keeps track of the current mode setting.
var currentMode Mode

// Countdown is the time limit (in seconds) for Timed mode.
var Countdown time.Duration = time.Second * 5 * 60
