package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	currentPage Page
}

func (p *app) Init() tea.Cmd {
	return p.currentPage.Init()
}

func (p *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, p.currentPage.Update(msg)
}

func (p *app) View() string {
	return p.currentPage.View()
}

func (p *app) changePage(page Page) {
	p.currentPage = page
}
