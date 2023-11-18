package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

// model for stopwatch
type stopwatch struct {
	startTime time.Time
}

// start starts the stopwatch
func (s *stopwatch) start() tea.Cmd {
	return s.tick()
}

// tick ticks the stopwatch at every 100ms interval.
func (s *stopwatch) tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(curTime time.Time) tea.Msg {
		if s.startTime.IsZero() {
			s.startTime = time.Now()
		}
		return TickMsg(curTime)
	})
}

// elapsed returns the elapsed duration.
func (s *stopwatch) elapsed() time.Duration {
	if s.startTime.IsZero() {
		return 0
	}
	return time.Since(s.startTime)
}

// view returns the UI string of stopwatch.
func (s *stopwatch) view() string {
	return s.elapsed().Round(time.Millisecond * 100).String()
}

// newStopwatch returns a new instance of stopwatch.
func newStopwatch() stopwatch {
	return stopwatch{}
}
