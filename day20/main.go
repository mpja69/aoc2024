package main

import (
	"bytes"
	"fmt"
	"os"
)

var (
	grid       [][]byte
	R, C       int
	walls      map[P]bool
	edge       map[P]bool
	start, end P
)

type P struct {
	r, c int
}

func main() {
	data, _ := os.ReadFile("data.txt")
	grid = bytes.Split(bytes.TrimSpace(data), []byte{'\n'})
	R = len(grid)
	C = len(grid[0])

	start = getPosOf('S')
	end = getPosOf('E')
	walls = getWalls()
	edge = getEdge()

	p1()
	p2()
}

func p1() {

	dist := getAllDistFromStartBFS()
	cheats := map[int]int{}
	sum := 0
	for p, d := range dist {
		for n := range getAllCheatPosWithRadiusBFS(p, 2) {
			if nd := dist[n] - d - 2; nd >= 100 {
				sum++
				cheats[nd]++
			}
		}
	}
	// fmt.Println(cheats)
	fmt.Printf("P1: (1409) %d\n", sum)
}

func p2() {

	dist := getAllDistFromStartBFS()
	cheats := map[int]int{}
	sum := 0
	for p, d := range dist {
		for n, r := range getAllCheatPosWithRadiusBFS(p, 20) {
			if nd := dist[n] - d - r; nd >= 100 {
				sum++
				cheats[nd]++
			}
		}
	}
	// fmt.Println(cheats)
	fmt.Printf("P2: (1012821) %d\n", sum)
}

type Q struct {
	p P
	i int
}

func getAllDistFromStartBFS() map[P]int {
	q := []Q{{start, 0}}
	seen := map[P]bool{}
	dist := map[P]int{}

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]
		p := curr.p

		// Seen
		if seen[p] {
			continue
		}
		seen[p] = true

		// WORK  - distance from start to p
		dist[p] = curr.i

		// Found - exit condition
		if p == end {
			return dist
		}

		// Neighbours
		for _, d := range []P{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
			n := add(p, d)
			if seen[n] || edge[n] || walls[n] {
				continue
			}
			q = append(q, Q{n, curr.i + 1})
		}

	}

	return map[P]int{}
}

func getAllCheatPosWithRadiusBFS(start P, number int) map[P]int {
	q := []Q{{start, 0}}
	seen := map[P]bool{}

	neighbourAndRadius := map[P]int{}

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]
		p := curr.p

		// Seen
		if seen[p] {
			continue
		}
		seen[p] = true

		// WORK - Collect possible neighbour cheat positions
		if curr.i <= number {
			neighbourAndRadius[p] = curr.i
		}

		// Exit condition: Finish up things in the queue, but stop add adding neighbours
		if curr.i > number {
			continue
		}

		// Neighbours
		for _, d := range []P{{1, 0}, {0, 1}, {-1, 0}, {0, -1}} {
			n := add(p, d)
			if seen[n] {
				continue
			}
			if edge[n] {
				continue
			}

			q = append(q, Q{n, curr.i + 1})
		}

	}

	return neighbourAndRadius
}

func add(a, b P) P { return P{a.r + b.r, a.c + b.c} }

func makeGridWith[T any](init T) [][]T {
	g := make([][]T, R)
	for r := range R {
		g[r] = make([]T, C)
		for c := range C {
			g[r][c] = init
		}
	}
	return g
}

func printGrid(pos P, enabled bool, seen map[P]bool) {
	g := makeGridWith(byte(' '))
	for p := range edge {
		g[p.r][p.c] = '*'
	}
	for p := range seen {
		g[p.r][p.c] = 'o'
	}
	if enabled {
		for p := range walls {
			g[p.r][p.c] = '#'
		}
	} else {
		for p := range walls {
			g[p.r][p.c] = '.'
		}
	}

	g[pos.r][pos.c] = '0'
	g[start.r][start.c] = 'S'
	g[end.r][end.c] = 'E'
	g[pos.r][pos.c] = '0'
	for r := range R {
		fmt.Printf("%s\n", g[r])
	}
	println()
}

func getPosOf(ch byte) P {
	for r := range R {
		for c := range C {
			if grid[r][c] == ch {
				return P{r, c}
			}
		}
	}
	return P{}
}

func getWalls() map[P]bool {
	walls := map[P]bool{}
	for r := 1; r < R-1; r++ {
		for c := 1; c < C-1; c++ {
			if grid[r][c] == '#' {
				walls[P{r, c}] = true
			}
		}
	}
	return walls
}

func getEdge() map[P]bool {
	outside := map[P]bool{}
	for r := range R {
		outside[P{r, 0}] = true
		outside[P{r, C - 1}] = true
	}
	for c := range C {
		outside[P{0, c}] = true
		outside[P{R - 1, c}] = true
	}
	return outside
}
