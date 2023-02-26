package main

import "github.com/charmbracelet/lipgloss"

const red = "#cc001b"

var greenTextAreaStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#5ac700")).
	Padding(1, 2).
	Width(50)

var redTextAreaStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(red)).
	Padding(1, 2).
	Width(50)

var pastTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#333333"))

var currentLetterStyle = lipgloss.NewStyle().
	Underline(true)

var errorOffsetStyle = lipgloss.NewStyle().
	Background(lipgloss.Color(red))

var wordHolderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#ffffff")).
	Padding(0, 1).
	Width(30)

	// var wrongCurrentLetterStyle = lipgloss.NewStyle().
	// 	Underline(true).
	// 	Foreground(lipgloss.Color("#c40000"))
