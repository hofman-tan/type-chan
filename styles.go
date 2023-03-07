package main

import "github.com/charmbracelet/lipgloss"

const red = lipgloss.Color("#cc001b")
const green = lipgloss.Color("#5ac700")
const grey = lipgloss.Color("#595959")
const white = lipgloss.Color("#ffffff")

const textareaWidth = 60

var ContainerStyle = lipgloss.NewStyle().
	Padding(2, 2)

var textAreaStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Padding(1, 2).
	Width(textareaWidth)

var greenTextAreaStyle = textAreaStyle.Copy().
	BorderForeground(green)

var redTextAreaStyle = textAreaStyle.Copy().
	BorderForeground(red)

var pastTextStyle = lipgloss.NewStyle().
	Foreground(grey)

var currentLetterStyle = lipgloss.NewStyle().
	Underline(true)

var errorOffsetStyle = lipgloss.NewStyle().
	Background(red)

var wordHolderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(white).
	Padding(0, 1).
	Width(30)

var progressBarBlankStyle = lipgloss.NewStyle().
	Foreground(grey)

var progressBarContentStyle = lipgloss.NewStyle().
	Foreground(white)

var sidebarStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(grey).
	Padding(0, 1).
	Width(17).
	Height(2)
