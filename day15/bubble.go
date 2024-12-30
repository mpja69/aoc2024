package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	once := false
	all := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":
			once = true
		case "a":
			all = true
		}
	}

	// Execute 1 or All of the moves
	for once || all {
		once = false
		move, err := m.moves.ReadByte()
		log.Printf("move: %c, err: %v", move, err)
		if err != nil {
			return m, tea.Quit
		}
		m.currPos = m.p1Update(m.currPos, move)
	}

	// Scroll the canvas if the "y-pos" gets 8 rows from the top or bottom
	cy := m.canvas.Cursor().Y
	if m.currPos.Y > (cy + m.canvas.ViewHeight - 8) {
		cy = m.currPos.Y - m.canvas.ViewHeight + 8
	} else if m.currPos.Y < (cy + 8) {
		cy = m.currPos.Y - 8
	}
	cursor := P{0, cy}
	m.canvas.SetCursor(cursor)

	return m, nil
}

func (m *model) View() string {
	s := fmt.Sprintf("Moves left: %d\n", m.moves.Len())
	s += m.canvas.View() + "\n"
	s += "N: [N]ext\t"
	s += "A: [A]ll\t"
	s += "Q: [Q]uittn"
	return s
}
