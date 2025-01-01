package main

import (
	"log"
	"os"
	"strings"

	"github.com/NimbleMarkets/ntcharts/canvas"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func p2(grid string, moves string) int {
	m := &model{}

	// Styles
	m.hlStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Background(lipgloss.Color("2"))
	m.defStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	// Init the canvas
	lines := strings.Split(string(grid), "\n")
	m.canvas = canvas.New(len(lines[0])*2, len(lines), canvas.WithViewWidth(100), canvas.WithViewHeight(30))
	for y, line := range lines {
		for x, r := range line {
			if r == 'O' {
				m.canvas.SetRuneWithStyle(P{x * 2, y}, '[', m.defStyle)
				m.canvas.SetRuneWithStyle(P{x*2 + 1, y}, ']', m.defStyle)
			} else if r == '@' {
				m.canvas.SetRuneWithStyle(P{x * 2, y}, '@', m.defStyle)
				m.canvas.SetRuneWithStyle(P{x*2 + 1, y}, '.', m.defStyle)
			} else {
				m.canvas.SetRuneWithStyle(P{x * 2, y}, r, m.defStyle)
				m.canvas.SetRuneWithStyle(P{x*2 + 1, y}, r, m.defStyle)
			}
		}
	}

	// Get start pos and the input stream
	m.currPos = m.start()
	m.canvas.SetCellStyle(m.currPos, m.hlStyle)
	m.moves = strings.NewReader(moves)

	// Run Bubble Tea
	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}

	// Calculate and return the result
	return m.p2SumGPS()
}

func (m model) p2Update(p P, chDir byte) P {
	dir, ch := m.peek(p, chDir)
	switch ch {
	case '.':
		p = m.move(p, dir)
	case '[', ']':
		log.Printf("Update...found: []")
		boxes, ok := m.p2FindBoxes(p, dir)
		if ok {
			m.p2PushBoxes(dir, boxes)
			p = m.move(p, dir)
		}
	}
	return p
}

func (m model) p2FindBoxes(pos, dir P) (boxes map[P]rune, ok bool) {
	if dir.X != 0 {
		boxes, ok = m.p2FindBoxesHorizontal(pos, dir)
	} else {

		boxes, ok = m.p2FindBoxesVertical(pos, dir)
	}
	for pp := range boxes {
		m.canvas.SetCellStyle(pp, m.hlStyle)
	}
	return
}

// BSF returns a set with possible pushable boxes, and a flag if they are pushable
func (m model) p2FindBoxesVertical(pos, dir P) (boxes map[P]rune, ok bool) {
	log.Printf("Vertical BSF")
	ok = true
	seen := map[P]bool{}
	boxes = map[P]rune{}
	q := []P{}
	q = append(q, pos.Add(dir))

	for len(q) > 0 {
		log.Printf("Queue not empty: %d", len(q))
		pos, q = q[0], q[1:] // Pop Queue
		log.Printf("Popped. Now have length: %d", len(q))
		if seen[pos] { // Skip seen
			continue
		}

		symbol := m.canvas.Cell(pos).Rune
		log.Printf("Passed not-seen: %v, %c", pos, symbol)

		// Check if to stop. (Stop searching, or stop adding neighbours)
		if symbol == '#' {
			log.Printf("# -> Return!")
			ok = false
			break
		} else if symbol == '.' {
			log.Printf(". -> Continue")
			continue
		}

		seen[pos] = true    // Add to seen
		boxes[pos] = symbol // Add to the result

		// Add neighbours
		// RIGHT
		right := m.canvas.Cell(pos.Add(P{1, 0})).Rune
		if symbol == '[' || right == ']' {
			q = append(q, pos.Add(P{1, 0}))
		}
		// LEFT
		left := m.canvas.Cell(pos.Add(P{-1, 0})).Rune
		if left == '[' || symbol == ']' {
			q = append(q, pos.Add(P{-1, 0}))
		}
		// UP/DOWN: Could be any of: ., #, [, ]
		q = append(q, pos.Add(dir))
	}
	return
}

func (m model) p2FindBoxesHorizontal(pos, dir P) (boxes map[P]rune, ok bool) {
	var symbol rune
	boxes = map[P]rune{}
	for {
		pos = pos.Add(dir)
		symbol = m.canvas.Cell(pos).Rune
		if symbol != '[' && symbol != ']' {
			break
		}
		boxes[pos] = symbol
	}
	ok = symbol == '.'
	return boxes, ok
}

func (m model) p2PushBoxes(dir P, boxes map[P]rune) {
	// "Clear" all boxes on current positions
	for pos := range boxes {
		m.canvas.SetCell(pos, canvas.NewCell('.'))
	}
	// Then "draw" the boxes on new posotions
	for pos, symbol := range boxes {
		m.canvas.SetCell(pos.Add(dir), canvas.NewCell(symbol))
	}
}
func (m model) p2SumGPS() int {
	sum := 0
	for y := range m.canvas.Height() {
		for x := range m.canvas.Width() {
			if m.canvas.Cell(P{x, y}).Rune == '[' {
				sum += 100*y + x
			}
		}
	}
	return sum
}
