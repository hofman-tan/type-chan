package app

import "github.com/charmbracelet/lipgloss"

const red = lipgloss.Color("#cc001b")
const green = lipgloss.Color("#5ac700")
const grey = lipgloss.Color("#595959")
const white = lipgloss.Color("#ffffff")

var containerStyle = lipgloss.NewStyle().
	Padding(2, 2)

var borderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Padding(1, 2)

var greenBorderStyle = borderStyle.Copy().
	BorderForeground(green)

var redBorderStyle = borderStyle.Copy().
	BorderForeground(red)

var whiteTextStyle = lipgloss.NewStyle().
	Foreground(white)

var greyTextStyle = lipgloss.NewStyle().
	Foreground(grey)

var underlinedStyle = lipgloss.NewStyle().
	Underline(true)

var redTextStyle = lipgloss.NewStyle().
	Background(red)

var wordHolderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(white).
	Padding(0, 1).
	Width(30)

var sidebarStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(grey).
	Padding(1, 1).
	Width(15)
