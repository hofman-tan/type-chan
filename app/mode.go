package app

type Mode int

const (
	Sprint Mode = iota
	Timed
)

var currentMode Mode
