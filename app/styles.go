package app

import "github.com/charmbracelet/lipgloss"

const paddingX = 8
const paddingY = 2

const red = lipgloss.Color("#cc001b")
const green = lipgloss.Color("#5ac700")
const grey = lipgloss.Color("#595959")

var greyTextStyle = lipgloss.NewStyle().
	Foreground(grey)

var underlinedStyle = lipgloss.NewStyle().
	Underline(true)

var redTextStyle = lipgloss.NewStyle().
	Background(red)
