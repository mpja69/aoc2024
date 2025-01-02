package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	once := false
	all := false
	ok := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":
			once = true
		case "a":
			m.t = time.Now()
			all = true
		}
	}

	// Execute 1 or All of the moves
	for once || all {
		once = false
		m.pos, m.cost, ok = m.update()
		if !ok {
			log.Printf("%v, %d", time.Since(m.t), m.cost)
			return m, tea.Quit
		}
	}

	// Scroll the canvas if the "y-pos" gets 8 rows from the top or bottom
	cy := m.canvas.Cursor().Y
	if m.pos.Y > (cy + m.canvas.ViewHeight - 8) {
		cy = m.pos.Y - m.canvas.ViewHeight + 8
	} else if m.pos.Y < (cy + 8) {
		cy = m.pos.Y - 8
	}
	cursor := P{0, cy}
	m.canvas.SetCursor(cursor)

	return m, nil
}

func (m *model) View() string {
	s := m.canvas.View() + "\n"
	s += "N: [N]ext\t"
	s += "A: [A]ll\t"
	s += "Q: [Q]uit"
	return s
}
