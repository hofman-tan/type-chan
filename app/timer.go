package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

type stopWatch struct {
	startTime time.Time
}

func (s *stopWatch) start() tea.Cmd {
	return s.tick()
}

func (s *stopWatch) tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(curTime time.Time) tea.Msg {
		if s.startTime.IsZero() {
			s.startTime = time.Now()
		}
		return TickMsg(curTime)
	})
}

func (s *stopWatch) elapsed() time.Duration {
	if s.startTime.IsZero() {
		return 0
	}
	return time.Since(s.startTime)
}

func (s *stopWatch) view() string {
	return s.elapsed().Round(time.Millisecond * 100).String()
}

func newStopwatch() stopWatch {
	return stopWatch{}
}
