package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type timer struct {
	seconds int
}

type TickMsg time.Time
type TimesUpMsg time.Time

func (t *timer) tick() tea.Cmd {
	return tea.Tick(time.Second, func(time time.Time) tea.Msg {
		if t.seconds > 0 {
			t.seconds -= 1 // decrement one second
			return TickMsg(time)
		} else {
			return TimesUpMsg(time)
		}
	})
}

func (t *timer) string() string {
	duration := time.Duration(float64(t.seconds) * float64(time.Second))
	return duration.String()
}

func newTimer() *timer {
	return &timer{
		seconds: 5 * 60, // 5 minutes
	}
}
