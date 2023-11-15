package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
)

type progressBar struct {
	progress.Model
}

func (p *progressBar) View() string {
	p.Width = windowWidth - paddingX*2 - 4
	return strings.Repeat(" ", paddingX) + p.Model.View()
}

func newProgressBar() *progressBar {
	pb := progress.New(progress.WithSolidFill(string(green)))
	return &progressBar{pb}
}
