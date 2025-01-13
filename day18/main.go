package main

import (
	"bytes"
	"container/list"
	"fmt"
	"github.com/NimbleMarkets/ntcharts/canvas"
	tea "github.com/charmbracelet/bubbletea"
	"image"
	"iter"
	"os"
)

var (
	GRID_SIZE = 6 //70
)

func main() {
	data, _ := os.ReadFile("sample.txt")

	// fmt.Printf("P1: %d\n", p1(data, 1024))
	p2Visual(data)
}

func p1(data []byte, nbrBytes int) int {
	r := bytes.NewReader(data)
	mem := map[Pos]bool{}

	for range nbrBytes {
		p := Pos{}
		fmt.Fscanf(r, "%d,%d\n", &p.X, &p.Y)
		mem[p] = true
	}
	res, _ := bsf(mem)
	return res
}

type model struct {
	c    canvas.Model
	next func() bool
	val  func() ([]Pos, map[Pos]bool, Pos)
}

// Binary Search point "between ok/false"
func p2Visual(data []byte) {
	m := model{}
	m.next, m.val = getUpdateFunc(data)

	m.c = canvas.New(GRID_SIZE, GRID_SIZE)
	if _, err := tea.NewProgram(&m).Run(); err != nil {
		os.Exit(1)
	}

}

func getUpdateFunc(data []byte) (func() bool, func() ([]Pos, map[Pos]bool, Pos)) {

	lo, hi := 0, 25 //3450
	p := Pos{}
	path := []Pos{}
	mem := map[Pos]bool{}
	ok := false

	next := func() bool {
		if lo < hi {
			mid := (lo + hi) / 2
			// fmt.Printf("lo: %d, mid: %d, hi: %d\n", lo, mid, hi)
			r := bytes.NewReader(data)
			mem = map[Pos]bool{}
			for range mid + 1 {
				fmt.Fscanf(r, "%d,%d\n", &p.X, &p.Y)
				mem[p] = true
			}

			path, ok = bsfPath(mem)

			if ok {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		return lo < hi
	}
	val := func() ([]Pos, map[Pos]bool, Pos) {
		return path, mem, p
	}
	return next, val
}

type Pos = image.Point
type Path struct {
	Pos
	path []Pos
}

func bsfPath(obstacles map[Pos]bool) ([]Pos, bool) {
	// Init the BFS
	seen := map[Pos]bool{}
	q := list.List{}
	q.PushBack(Path{Pos{0, 0}, []Pos{}})

	for q.Len() > 0 {
		// Queue
		e := q.Front()
		q.Remove(e)
		curr := e.Value.(Path)

		// Seen
		if seen[curr.Pos] {
			continue
		}
		seen[curr.Pos] = true

		// Found
		if curr.Pos.Eq(Pos{GRID_SIZE, GRID_SIZE}) {
			return curr.path, true
		}

		// Neighbours
		for n := range neighbours(curr.Pos) {
			if obstacles[n] || seen[n] {
				continue
			}
			q.PushBack(Path{n, append(curr.path, n)})
		}
	}
	// No path
	return nil, false
}

type Item struct {
	Pos
	i int
}

func bsf(obstacles map[Pos]bool) (int, bool) {
	// Init the BFS
	seen := map[Pos]bool{}
	q := list.List{}
	q.PushBack(Item{Pos{0, 0}, 0})

	for q.Len() > 0 {
		// Queue
		e := q.Front()
		q.Remove(e)
		curr := e.Value.(Item)

		// Seen
		if seen[curr.Pos] {
			continue
		}
		seen[curr.Pos] = true

		// Found
		if curr.Pos.Eq(Pos{GRID_SIZE, GRID_SIZE}) {
			return curr.i, true
		}

		// Neighbours
		for n := range neighbours(curr.Pos) {
			if obstacles[n] || seen[n] {
				continue
			}
			q.PushBack(Item{n, curr.i + 1})
		}
	}
	// No path
	return 0, false
}

// Testing a "pointless iterator, since it could just as wll been a slice
//
// A closure that holds...
//   - the "start pos p"
//   - and the state of its inner iteration over directions
//   - and "yields" the current value via the callback yield(...)
func neighbours(p Pos) iter.Seq[Pos] {
	return func(yield func(p Pos) bool) {
		for _, d := range []Pos{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
			n := p.Add(d)
			if n.In(image.Rect(0, 0, GRID_SIZE+1, GRID_SIZE+1)) {
				yield(n)
			}
		}
	}
}
