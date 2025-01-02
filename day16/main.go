package main

import (
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
	input := strings.Trim(string(data), " \n")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Printf("P1: (91464) %d\n", p1(input))
	//P1: 354040 to high DFS/Stack
	//P1:  95436 to high BFS/Queue
}

var (
	hlStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Background(lipgloss.Color("2"))
	defStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
)

type P = canvas.Point // Use type ALIAS, instead of type DEFINITION (to get all the methods)

type model struct {
	canvas canvas.Model
	cost   int
	pos    P
	update func() ([]int, bool)
}

func p1(input string) int {
	m := &model{}

	// Init the canvas
	lines := strings.Split(string(input), "\n")
	m.canvas = canvas.New(len(lines[0]), len(lines), canvas.WithViewWidth(150), canvas.WithViewHeight(50))
	m.canvas.SetLinesWithStyle(lines, defStyle)

	// Get start pos and the input stream
	m.pos = P{1, m.canvas.Height() - 2}
	m.canvas.SetCellStyle(m.pos, hlStyle)
	m.update = getBFSFunc(m.canvas, m.pos, P{1, 0})

	return BFS(m.canvas, m.pos, P{1, 0})

	// fmt.Println(m.canvas.View())
	// Run Bubble Tea
	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}

	// Calculate and return the result
	return m.cost
}

// Create a closure for a BFS function
func getBFSFunc(c canvas.Model, pos, dir P) func() ([]int, bool) {
	// NOTE: seen can NOT be just a set. It has to map a postion to the cost of getting there!
	seen := map[P]int{}
	res := []int{}
	q := Queue{}
	q.push(Item{pos, dir, 0})
	return func() ([]int, bool) {

		// "While" the Queue is not empty
		if q.empty() {
			return res, false
		}

		// Get current item
		curr := q.pop()

		// NOTE: If we reach a place to lower cost, then continue searching
		if cost, ok := seen[curr.pos]; ok {
			if curr.cost > cost {
				return res, true
			}
		}

		// Goal condition and calculation
		if c.Cell(curr.pos).Rune == 'E' {
			// Found a solution!
			res = append(res, curr.cost)
			return res, true
		}

		// NOTE: Save the cost for getting here
		c.SetCellStyle(curr.pos, hlStyle)
		seen[curr.pos] = curr.cost

		// Find/Add new neighbours to Queue
		for _, neighbour := range neighbours(c, curr.pos, curr.dir) {
			neighbour.cost += curr.cost
			q.push(neighbour)
		}

		// Return from this "iteration"
		return res, true
	}
}

type Item struct {
	pos  P
	dir  P
	cost int
}

// Get possible neighbours (and their new dir and cost)
func neighbours(c canvas.Model, pos, dir P) []Item {
	res := []Item{}
	// Straight
	d := dir
	p := pos.Add(d)
	if c.Cell(p).Rune != '#' {
		res = append(res, Item{p, d, 1})
	}

	// Right
	d.X, d.Y = dir.Y, -dir.X // Turn right
	p = pos.Add(d)
	if c.Cell(p).Rune != '#' {
		res = append(res, Item{p, d, 1001})
	}
	// Left
	d.X, d.Y = -dir.Y, dir.X // Turn left
	p = pos.Add(d)
	if c.Cell(p).Rune != '#' {
		res = append(res, Item{p, d, 1001})
	}
	return res
}

// ------------------------------- Utility functions --------------------------------

type Set map[P]bool

func (s Set) add(pos P) {
	s[pos] = true
}
func (s Set) delete(pos P) {
	s[pos] = false
}
func (s Set) has(pos P) bool {
	return s[pos]
}

type Stack struct {
	items []Item
}

func (q *Stack) pop() Item {
	item := q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	return item
}

func (q *Stack) push(item Item) {
	q.items = append(q.items, item)
}

func (q *Stack) empty() bool {
	return len(q.items) == 0
}

type Queue struct {
	items []Item
}

func (q *Queue) pop() Item {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) push(item Item) {
	q.items = append(q.items, item)
}

func (q *Queue) empty() bool {
	return len(q.items) == 0
}

func BFS(c canvas.Model, pos, dir P) int {
	// seen := Set{}
	seen := map[P]int{}
	res := []int{}
	q := Queue{}
	q.push(Item{pos, dir, 0})

	// While the Queue is not empty
	for !q.empty() {

		// Get current item
		curr := q.pop()

		// Guard againt already seen items with map
		// if seen.has(curr.pos) {
		// 	continue
		// }
		if cost, ok := seen[curr.pos]; ok {
			if curr.cost > cost {
				continue
			}
		}

		// Goal condition and calculation
		if c.Cell(curr.pos).Rune == 'E' {
			// Found a solution! Add to list, and stop this path.
			res = append(res, curr.cost)
			continue
		}

		// seen.add(curr.pos)
		seen[curr.pos] = curr.cost

		// Find/Add new neighbours to Queue
		for _, neighbour := range neighbours(c, curr.pos, curr.dir) {
			neighbour.cost += curr.cost
			q.push(neighbour)
		}
	}

	fmt.Printf("Results: %v,\n", res)
	min := res[0]
	for _, i := range res[1:] {
		if i < min {
			min = i
		}
	}
	return min
}
