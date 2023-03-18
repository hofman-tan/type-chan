package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Timer interface {
	tick() tea.Cmd
	getTimeElapsed() time.Duration
	string() string
}

type countUpTimer struct {
	startTime time.Time
}

type countDownTimer struct {
	seconds     int
	secondsLeft int
}

type TickMsg time.Time
type TimesUpMsg time.Time

func newCountUpTimer() *countUpTimer {
	return &countUpTimer{}
}

func (t *countUpTimer) tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(curTime time.Time) tea.Msg {
		if t.startTime.IsZero() {
			t.startTime = time.Now()
		}
		return TickMsg(curTime)
	})
}

func (t *countUpTimer) string() string {
	d := t.getTimeElapsed()

	ms := d.Milliseconds() % 1000
	msStr := fmt.Sprintf("%03d", ms)[:2]
	s := int(d.Seconds()) % 60
	m := int(d.Minutes())

	return fmt.Sprintf("%02d:%02d:%s", m, s, msStr)
}

func (t *countUpTimer) getTimeElapsed() time.Duration {
	if t.startTime.IsZero() {
		return time.Duration(0)
	}
	return time.Since(t.startTime)
}

func newCountDownTimer(seconds int) *countDownTimer {
	return &countDownTimer{
		seconds:     seconds,
		secondsLeft: seconds,
	}
}

func (t *countDownTimer) tick() tea.Cmd {
	return tea.Tick(time.Second, func(curTime time.Time) tea.Msg {
		t.secondsLeft--

		if t.secondsLeft > 0 {
			return TickMsg(curTime)
		} else {
			return TimesUpMsg(curTime)
		}
	})
}

func (t *countDownTimer) getTimeElapsed() time.Duration {
	return time.Duration(float64(t.seconds-t.secondsLeft) * float64(time.Second))
}

func (t *countDownTimer) string() string {
	d := time.Duration(float64(t.secondsLeft) * float64(time.Second))

	s := int(d.Seconds()) % 60
	m := int(d.Minutes())

	return fmt.Sprintf("%02d:%02d", m, s)
}
