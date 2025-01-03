package main

import (
	"container/heap"
	"fmt"
	"log"
	"maps"
	"math"
	"os"
	"slices"
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

	fmt.Printf("P2: (494) %d\n", p2(input))
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

func p2(input string) int {
	m := &model{}

	// Init the canvas
	lines := strings.Split(string(input), "\n")
	m.canvas = canvas.New(len(lines[0]), len(lines), canvas.WithViewWidth(150), canvas.WithViewHeight(50))
	m.canvas.SetLinesWithStyle(lines, defStyle)
	m.pos = findStart(m.canvas)
	m.canvas.SetCellStyle(m.pos, hlStyle)

	endStates, backtrackMap := dijkstraWithBacktrack(m.canvas, m.pos, P{1, 0})
	tracks := backtrackMap.backTrack(endStates)
	return countPositions(tracks)
}

func countPositions(tracks map[State]bool) int {
	set := map[P]bool{}
	for state := range tracks {
		set[state.pos] = true
	}
	return len(set)
}
func findStart(c canvas.Model) P {
	for x := range c.Width() {
		for y := range c.Height() {
			pos := P{x, y}
			if c.Cell(pos).Rune == 'S' {
				return pos
			}
		}
	}
	return P{}
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
	return cost

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
		for _, neighbour := range neighbours(curr.value.pos, curr.value.dir, cost) {
			if c.Cell(neighbour.value.pos).Rune == '#' {
				continue
			}
			if seen[neighbour.value] {
				continue
			}
			heap.Push(&pq, &neighbour)
			// log.Printf("Added to queue: %v", neighbour)
		}

		// Return from this "iteration"
		// log.Printf("End of BFS: %v", res)
		return state.pos, cost, true
	}
}

// Could either be a step forward or a step to the right/left.
func neighbours(pos, dir P, cost int) []Item {
	res := []Item{}
	// Straight
	d := dir
	p := pos.Add(d)
	res = append(res, Item{value: State{p, d}, priority: cost + 1})

	// Only do the turn. Stay in same place.
	// Right
	d.X, d.Y = dir.Y, -dir.X // Turn right
	p = pos.Add(d)
	res = append(res, Item{value: State{p, d}, priority: cost + 1001})
	// Left
	d.X, d.Y = -dir.Y, dir.X // Turn left
	p = pos.Add(d)
	res = append(res, Item{value: State{p, d}, priority: cost + 1001})
	return res
}

type CostMap map[State]int
type BacktrackMap map[State]map[State]bool // A Map of a Set

func dijkstraWithBacktrack(c canvas.Model, pos, dir P) (map[State]bool, BacktrackMap) {
	lowestCostAt := CostMap{} // Like a seen-set, but mapped to a cost
	bestCost := math.MaxInt   // Init with max cost

	pq := PriorityQueue{}
	heap.Push(&pq, &Item{State{pos, dir}, 0, 0})

	backtrackMap := BacktrackMap{} // A Map of a Set
	endStates := map[State]bool{}  // If we can reach the end from different tiles

	// While the Queue is not empty
	for pq.Len() > 0 {
		// Get current item
		curr := heap.Pop(&pq).(*Item)
		state := curr.value
		cost := curr.priority

		// Avoid more expensive paths
		if cost > lowestCostAt.getOrMax(state) {
			continue
		}
		// Add to seen-set
		// lowestCostAt[state] = curr.priority // Unnecessary?!

		// Goal condition
		if c.Cell(curr.value.pos).Rune == 'E' {
			if cost > bestCost {
				// Allready found all solutions that are "better/best"
				break
			}
			// Found a new solution!
			bestCost = cost
			endStates[state] = true
		}

		// Find/Add new neighbours to Queue
		for _, neighbour := range neighbours(curr.value.pos, curr.value.dir, cost) {
			newCost, newPos := neighbour.priority, neighbour.value.pos
			if c.Cell(newPos).Rune == '#' {
				continue
			}
			lowest := lowestCostAt.getOrMax(neighbour.value)
			if newCost > lowest {
				continue
			} else if newCost < lowest {
				backtrackMap[neighbour.value] = map[State]bool{}
				lowestCostAt[neighbour.value] = newCost
			}
			backtrackMap[neighbour.value][state] = true
			heap.Push(&pq, &neighbour)
		}
	}
	return endStates, backtrackMap
}

// NOTE: Get max instead of 0, (zero), if missing in the map
func (cm CostMap) getOrMax(state State) int {
	if lowest, ok := cm[state]; ok {
		return lowest
	}
	return math.MaxInt
}

// Start at the end(s) and flodfill backwards
func (btm BacktrackMap) backTrack(endStates map[State]bool) map[State]bool {
	seen := map[State]bool{}                       // A set (map) with only the positions
	states := slices.Collect(maps.Keys(endStates)) // Init a the queue (slice) with what we have in seen
	maps.Insert(seen, maps.All(endStates))         // Init the seen map with the end states
	for len(states) > 0 {
		key := states[0]
		states = states[1:]
		for last := range btm[key] {
			if seen[last] {
				continue
			}
			seen[last] = true
			states = append(states, last)
		}
	}
	return seen
}

func dijkstra(c canvas.Model, pos, dir P) int {
	seen := map[State]bool{}
	bestCost := 0

	pq := PriorityQueue{}
	heap.Push(&pq, &Item{State{pos, dir}, 0, 0})

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
			bestCost = curr.priority
			break
		}

		// Find/Add new neighbours to Queue
		for _, neighbour := range neighbours(curr.value.pos, curr.value.dir, curr.priority) {
			if c.Cell(neighbour.value.pos).Rune == '#' {
				continue
			}
			if seen[neighbour.value] {
				continue
			}
			heap.Push(&pq, &neighbour)
		}
	}
	return bestCost
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
