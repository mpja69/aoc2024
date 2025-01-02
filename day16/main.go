package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	update func() (P, int, bool)
	t      time.Time
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
	m.update = dijkstraFunc(m.canvas, m.pos, P{1, 0})
	m.canvas.SetCursor(m.pos)

	t := time.Now()
	cost := dijkstra(m.canvas, m.pos, P{1, 0})
	log.Printf("%v, %d\n", time.Since(t), cost)

	// fmt.Println(m.canvas.View())
	// Run Bubble Tea
	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}

	// Calculate and return the result
	return m.cost
}

type State struct {
	pos P
	dir P
}

type Item struct {
	value    State // The actual value we are storing
	priority int   // I.e. the cost
	index    int   // Used by the heap
}

// Define the list that holds the items in the Queue
type PriorityQueue []*Item

// Create a closure for a BFS function
func dijkstraFunc(c canvas.Model, pos, dir P) func() (P, int, bool) {
	// NOTE: seen can NOT be just a set. It has to map a postion to the cost of getting there!
	seen := map[State]bool{}
	pq := PriorityQueue{}
	startState := State{pos, dir}
	seen[startState] = true
	heap.Push(&pq, &Item{value: startState, priority: 0})

	log.Printf("Created queue: %v", pq)

	return func() (P, int, bool) {
		// log.Printf("Enter BFS...")

		// "While" the Queue is not empty
		if pq.Len() == 0 {
			// log.Printf("Empty queue: %v", res)
			return P{}, -1, false
		}

		// Pop current item
		curr := heap.Pop(&pq).(*Item)
		cost, state := curr.priority, curr.value

		c.SetCellStyle(state.pos, hlStyle)
		seen[state] = true

		// Goal condition and calculation
		if c.Cell(state.pos).Rune == 'E' {
			// Found the solution!
			return state.pos, cost, false
		}

		// Find/Add new neighbours to Queue
		for _, neighbour := range neighbours(curr.value.pos, curr.value.dir) {
			if c.Cell(neighbour.value.pos).Rune == '#' {
				continue
			}
			if seen[neighbour.value] {
				continue
			}
			neighbour.priority += curr.priority
			heap.Push(&pq, &neighbour)
			// log.Printf("Added to queue: %v", neighbour)
		}

		// Return from this "iteration"
		// log.Printf("End of BFS: %v", res)
		return state.pos, cost, true
	}
}

// Could either be a step forward or a step to the right/left.
func neighbours(pos, dir P) []Item {
	res := []Item{}
	// Straight
	d := dir
	p := pos.Add(d)
	res = append(res, Item{value: State{p, d}, priority: 1})

	// Only do the turn. Stay in same place.
	// Right
	d.X, d.Y = dir.Y, -dir.X // Turn right
	p = pos.Add(d)
	res = append(res, Item{value: State{p, d}, priority: 1001})
	// Left
	d.X, d.Y = -dir.Y, dir.X // Turn left
	p = pos.Add(d)
	res = append(res, Item{value: State{p, d}, priority: 1001})
	return res
}

func dijkstra(c canvas.Model, pos, dir P) int {
	seen := map[State]bool{}
	res := 0
	pq := PriorityQueue{}
	startState := State{pos, dir}
	heap.Push(&pq, &Item{startState, 0, 0})

	// While the Queue is not empty
	for pq.Len() > 0 {

		// Get current item
		curr := heap.Pop(&pq).(*Item)
		state := curr.value

		// Add to seen-set
		seen[state] = true

		// Goal condition
		if c.Cell(curr.value.pos).Rune == 'E' {
			// Found the solution!
			res = curr.priority
			break
		}

		// Find/Add new neighbours to Queue
		for _, neighbour := range neighbours(curr.value.pos, curr.value.dir) {
			if c.Cell(neighbour.value.pos).Rune == '#' {
				continue
			}
			if seen[neighbour.value] {
				continue
			}
			neighbour.priority += curr.priority
			heap.Push(&pq, &neighbour)
		}
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
