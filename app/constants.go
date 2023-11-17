package app

import "github.com/charmbracelet/lipgloss"

// TODO: move under app.go
// Mode is the user setting that determines the type of typing test.

type Mode int

const (
	Sprint Mode = iota
	Timed
)

const maxMistypedCount int = 10
const quoteBufferSize int = 2
const scrollTextHeight int = 3

const paddingX int = 10
const paddingY int = 2
const minWindowWidth int = 50

const red = lipgloss.Color("#cc001b")
const green = lipgloss.Color("#5ac700")
const grey = lipgloss.Color("#595959")
