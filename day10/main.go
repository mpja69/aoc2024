package main

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"os"
)

var (
	R, C int
)

func main() {
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := bytes.Trim(data, " \n")
	grid := bytes.Split(input, []byte{'\n'})

	R = len(grid)
	C = len(grid[0])

	fmt.Printf("Part 1: (36, 744): %d\n", part1(grid))
	fmt.Printf("Part 2: (81, 1651): %d\n", part2(grid))
}

type Pos struct{ r, c int }
type Grid [][]byte

func part1(grid Grid) int {
	score := 0

	trailheads := grid.findTrailheads('0')
	for _, start := range trailheads {
		s := grid.calcScore(start, '9')
		score += s
	}
	return score
}

func (g Grid) findTrailheads(start byte) []Pos {
	trailheads := []Pos{}
	for r := range R {
		for c := range C {
			if g[r][c] == start {
				trailheads = append(trailheads, Pos{r, c})
			}
		}
	}
	return trailheads
}

// Breadth First search
func (g Grid) calcScore(start Pos, end byte) int {
	q := list.New()           // Use a queue for the upcoming steps
	visited := map[Pos]bool{} // Use a set to avoid duplicate visits

	// Initial values
	score := 0
	q.PushBack(start)

	// Process the queue
	for q.Len() > 0 {

		// Pop the oldest item
		curr := q.Front().Value.(Pos)
		q.Remove(q.Front())

		// Track visited positions in a set
		if visited[curr] {
			continue
		}
		visited[curr] = true

		// Have we found an end position?
		if g[curr.r][curr.c] == end {
			score++
		}

		// Push all possible neighbours to the queue
		for _, neighbour := range g.neighbours(curr) {
			q.PushBack(neighbour)
		}
	}
	return score
}

func (g Grid) neighbours(p Pos) []Pos {
	neighbours := []Pos{}
	val := g[p.r][p.c]

	// There are 4 potential steps, Up, Left and Right
	potentialNeighbours := []Pos{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, s := range potentialNeighbours {
		r, c := p.r+s.r, p.c+s.c
		// Check inside grid
		if r < 0 || r >= R || c < 0 || c >= C {
			continue
		}
		// Check valid height
		if g[r][c] != val+1 {
			continue
		}

		neighbours = append(neighbours, Pos{r, c})
	}
	return neighbours
}

func part2(grid Grid) int {
	trails := 0
	trailheads := grid.findTrailheads('0')
	for _, start := range trailheads {
		s := grid.calcRating(start, '9')
		trails += s
	}
	return trails
}

type Item struct {
	step int
	path [10]Pos
}

// Breadth First search - Part 2
func (g Grid) calcRating(start Pos, end byte) int {
	q := list.New()                     // Use a queue for the upcoming steps
	distinctPaths := map[[10]Pos]bool{} // Save distinct paths

	// Initial values
	curr := Item{}
	curr.path[0] = start
	q.PushBack(curr)

	// Process the queue
	for q.Len() > 0 {

		// Pop the oldest item
		curr := q.Front().Value.(Item)
		q.Remove(q.Front())

		pos := curr.path[curr.step]

		// Have we found a new trail to same end?
		if g[pos.r][pos.c] == end {
			distinctPaths[curr.path] = true
			continue // Not needed? Because we will NOT get any possible neighbours
		}

		// Push all possible neighbours to the queue
		for _, neighbour := range g.neighbours(pos) {
			next := Item{}
			next.step = curr.step + 1
			copy(next.path[:], curr.path[:])
			next.path[next.step] = neighbour
			q.PushBack(next)
		}
	}
	return len(distinctPaths)
}
