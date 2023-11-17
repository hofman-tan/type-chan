package app

import "time"

// currentMode keeps track of the current mode setting.
var currentMode Mode

// Timeout is the time limit (in seconds) for Timed mode.
var Timeout time.Duration = time.Second * 5 * 60
