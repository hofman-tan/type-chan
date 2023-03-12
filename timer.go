package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type timer struct {
	seconds int
	countUp bool
}

type TickMsg time.Time
type TimesUpMsg time.Time

func (t *timer) tick() tea.Cmd {
	return tea.Tick(time.Second, func(time time.Time) tea.Msg {
		if t.countUp {
			// increment one second
			t.seconds++
			return TickMsg(time)

		} else {
			if t.seconds > 0 {
				// decrement one second
				t.seconds--
				return TickMsg(time)
			} else {
				return TimesUpMsg(time)
			}
		}
	})
}

func (t *timer) string() string {
	duration := t.getTime()
	return duration.String()
}

func (t *timer) getTime() time.Duration {
	return time.Duration(float64(t.seconds) * float64(time.Second))
}

func newCountUpTimer() *timer {
	return &timer{countUp: true}
}

func newCountDownTimer(seconds int) *timer {
	return &timer{
		countUp: false,
		seconds: seconds,
	}
}
