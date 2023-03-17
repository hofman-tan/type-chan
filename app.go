package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	currentPage Page
}

func (p *app) Init() tea.Cmd {
	return p.currentPage.init()
}

func (p *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, p.currentPage.update(msg)
}

func (p *app) View() string {
	return containerStyle.Render(p.currentPage.view())
}

func (p *app) changePage(page Page) {
	p.currentPage = page
}
