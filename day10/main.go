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
}

type Pos struct{ r, c int }
type Grid [][]byte

func part1(grid Grid) int {
	score := 0

	trailheads := grid.findTrailheads('0')
	for _, start := range trailheads {
		s := grid.findPaths(start, '9')
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

func (g Grid) findPaths(start Pos, end byte) int {
	// Use a queue for the upcoming steps
	q := list.New()
	// Use a set to avoid duplicate visits
	visited := map[Pos]bool{}

	score := 0
	q.PushBack(start)

	// Breadth First search
	for q.Len() > 0 {

		// Pop the oldest item
		curr := q.Front().Value.(Pos)
		q.Remove(q.Front())

		// Track visited positions in a set
		if visited[curr] {
			continue
		}
		visited[curr] = true

		// Have we found the end?
		if g[curr.r][curr.c] == end {
			score++
		}

		// Push all possible neighbours to the queue
		for _, neighbour := range g.possibleNeighbours(curr) {
			q.PushBack(neighbour)
		}
	}
	return score
}

func (g Grid) possibleNeighbours(p Pos) []Pos {
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
