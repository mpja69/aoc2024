package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NimbleMarkets/ntcharts/canvas"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	data, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := bytes.Trim(data, " \n")
	parts := strings.Split(string(input), "\n\n")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	fmt.Printf("P1: %d\n", p1(parts[0], parts[1]))

	// small: 2028
	// large:	10092
	// data:	1514333
}

type model struct {
	canvas   canvas.Model
	defStyle lipgloss.Style
	hlStyle  lipgloss.Style
	currPos  P
	moves    *strings.Reader
}

type P = canvas.Point // Use type ALIAS, instead of type DEFINITION

func p1(grid string, moves string) int {
	m := &model{}

	// Styles
	m.hlStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Background(lipgloss.Color("2"))
	m.defStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	// Init the canvas
	lines := strings.Split(string(grid), "\n")
	m.canvas = canvas.New(len(lines[0]), len(lines), canvas.WithViewWidth(60), canvas.WithViewHeight(30))
	m.canvas.SetLinesWithStyle(lines, m.defStyle)
	m.canvas.Focus()

	// Get start pos and the input stream
	m.currPos = m.start()
	m.canvas.SetCellStyle(m.currPos, m.hlStyle)
	m.moves = strings.NewReader(moves)

	// Run Bubble Tea
	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}

	// Calculate and return the result
	return m.sumGPS()
}

func (m model) p1Update(p P, chDir byte) P {
	dir, ch := m.peek(p, chDir)
	switch ch {
	case '.':
		p = m.move(p, dir)
	case 'O':
		boxes, ok := m.checkBoxes(p, dir)
		if ok {
			m.pushBoxes(p, dir, boxes)
			p = m.move(p, dir)
		}
	}
	return p
}

func (m model) peek(p P, move byte) (dir P, ch byte) {
	movements := map[byte]P{'^': {0, -1}, '>': {1, 0}, 'v': {0, 1}, '<': {-1, 0}}
	dir = movements[move]
	n := p.Add(dir)
	ch = byte(m.canvas.Cell(n).Rune)
	return
}

func (m model) move(pos, dir P) P {
	m.canvas.SetCell(pos, canvas.NewCellWithStyle('.', m.defStyle))
	pos = pos.Add(dir)
	m.canvas.SetCell(pos, canvas.NewCellWithStyle('@', m.hlStyle))
	return pos
}

func (m model) checkBoxes(pos, dir P) (boxes int, ok bool) {
	for {
		pos = pos.Add(dir)
		if m.canvas.Cell(pos).Rune != 'O' {
			break
		}
		boxes++
	}
	ok = m.canvas.Cell(pos).Rune == '.'
	return
}

func (m model) pushBoxes(p, d P, boxes int) {
	p = p.Add(d)
	m.canvas.SetCell(p, canvas.NewCellWithStyle('.', m.defStyle))
	for i := 0; i < boxes; i++ {
		p = p.Add(d)
		m.canvas.SetCell(p, canvas.NewCellWithStyle('O', m.defStyle))
	}
}

func (m model) start() P {
	for y := range m.canvas.Height() {
		for x := range m.canvas.Width() {
			if m.canvas.Cell(P{x, y}).Rune == '@' {
				return P{x, y}
			}
		}
	}
	return P{}
}

func (m model) sumGPS() int {
	sum := 0
	for y := range m.canvas.Height() {
		for x := range m.canvas.Width() {
			if m.canvas.Cell(P{x, y}).Rune == 'O' {
				sum += 100*y + x
			}
		}
	}
	return sum
}
