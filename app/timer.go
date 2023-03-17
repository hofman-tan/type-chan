package app

import (
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

// func fmtDuration(d time.Duration) string {
// 	d = d.Round(10 * time.Millisecond)
// 	m := d / time.Minute
// 	d -= m * time.Minute

// 	s := d / time.Second
// 	d -= s * time.Second

// 	ms := d / time.Millisecond
// 	return fmt.Sprintf("%02d:%02d:%02d", m, s, ms)
// }

func (t *countUpTimer) string() string {
	return t.getTimeElapsed().String()
}

func (t *countUpTimer) getTimeElapsed() time.Duration {
	return time.Since(t.startTime).Round(10 * time.Millisecond)
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
	duration := time.Duration(float64(t.secondsLeft) * float64(time.Second))
	return duration.String()
}
