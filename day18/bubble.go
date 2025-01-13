package main

import (
	"github.com/NimbleMarkets/ntcharts/canvas"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":
			m.next()
			m.updateCanvas()
		}
	}
	return m, nil
}

func (m *model) updateCanvas() {
	path, obstacles, _ := m.val()
	for x := range GRID_SIZE {
		for y := range GRID_SIZE {
			if obstacles[Pos{x, y}] {
				m.c.SetCell(Pos{x, y}, canvas.NewCell('#'))
			} else {
				m.c.SetCell(Pos{x, y}, canvas.NewCell('.'))
			}
		}
	}
	for _, p := range path {
		m.c.SetCell(p, canvas.NewCell('0'))
	}
}

func (m *model) View() string {
	s := m.c.View() + "\n"
	s += "N: [N]ext\t"
	s += "Q: [Q]uit"
	return s
}
