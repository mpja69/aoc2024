package main

import (
	"bytes"
	"container/list"
	"fmt"
	"image"
	"iter"
	"os"
)

var (
	GRID_SIZE = 70
)

func main() {
	data, _ := os.ReadFile("data.txt")

	fmt.Printf("P1: %d\n", p1(data, 1024))
	nbr, pos := p2(data)
	fmt.Printf("P2: Nbr: %d, Pos: %v\n", nbr, pos)
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

// Binary Search point "between ok/false"
func p2(data []byte) (int, Pos) {
	lo, hi := 0, 3450
	p := Pos{}
	for lo < hi {
		mid := (lo + hi) / 2
		fmt.Printf("lo: %d, mid: %d, hi: %d\n", lo, mid, hi)
		r := bytes.NewReader(data)
		mem := map[Pos]bool{}
		for range mid + 1 {
			fmt.Fscanf(r, "%d,%d\n", &p.X, &p.Y)
			mem[p] = true
		}

		_, ok := bsf(mem)

		if ok {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo, p
}

type Pos = image.Point
type Item struct {
	Pos
	i int
}

func bsf(obstacles map[Pos]bool) (int, bool) {
	// BFS to find the shorters path
	seen := map[Pos]bool{}
	q := list.List{}

	q.PushBack(Item{Pos{0, 0}, 0})

	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		curr := e.Value.(Item)

		if seen[curr.Pos] {
			continue
		}
		seen[curr.Pos] = true

		if curr.Pos.Eq(Pos{GRID_SIZE, GRID_SIZE}) {
			return curr.i, true
		}

		for n := range neighbours(curr.Pos) {
			if !n.In(image.Rect(0, 0, GRID_SIZE+1, GRID_SIZE+1)) || obstacles[n] || seen[n] {
				continue
			}
			q.PushBack(Item{n, curr.i + 1})
		}
	}
	return 0, false
}

// Testing an iterator
//
// A closure that holds:
//
//	the "start pos p"
//	and the state of its inner iteration over directions
//	and "yields" the current value via the callback yield(...)
func neighbours(p Pos) iter.Seq[Pos] {
	return func(yield func(p Pos) bool) {
		for _, d := range []Pos{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
			n := p.Add(d)
			yield(n)
		}
	}
}
