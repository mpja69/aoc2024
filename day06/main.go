package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// ------------------------ Globals -------------------------
var (
	R, C       int                                       // Boundries of grid
	directions = []Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // 0..3 directions. ONLY turn right
)

type Pos struct {
	r int
	c int
}

func (p Pos) pos() Pos {
	return p
}

type Path struct {
	Pos
	d int
}

func (p Path) pos() Pos {
	return p.Pos
}
func (p *Path) moveForward() {
	p.r += directions[p.d].r
	p.c += directions[p.d].c
}
func (p *Path) peek() Pos {
	r := p.r + directions[p.d].r
	c := p.c + directions[p.d].c
	return Pos{r, c}
}
func (p *Path) inside() bool {
	r := p.r + directions[p.d].r
	c := p.c + directions[p.d].c
	return !(r < 0 || r >= R || c < 0 || c >= R)
}
func (p *Path) turnRight() {
	p.d = (p.d + 1) % len(directions)
}

// ----------------------- Maps/Sets ------------------------
type Item interface {
	pos() Pos
}

type Set map[Item]bool

func (s Set) add(i Item) {
	s[i] = true
}
func (s Set) has(i Item) bool {
	return s[i]
}
func (s Set) delete(i Item) {
	s[i] = false
}

// ----------------------- Generics ------------------------
// func has[T comparable](set map[T]bool, item T) bool {
// 	return set[item]
// }
// func add[T comparable](set map[T]bool, item T) {
// 	set[item] = true
// }

func main() {
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	data = bytes.Trim(data, " \n")
	grid := bytes.Split(data, []byte{'\n'})

	R = len(grid)
	C = len(grid[0])

	start, obstacles := findStart(grid)
	t := time.Now()
	fmt.Println("Part 1 (41, 4559):", p1(start, obstacles))
	fmt.Println("Part 2 (6, 1604):", p2(start, obstacles))
	fmt.Printf("Millis: %d\n", time.Since(t).Milliseconds())
}

func p1(start Path, obstacles Set) int {
	visitedPositions := walkUntilOffTheGrid(start, obstacles)
	return len(visitedPositions)
}

func findStart(grid [][]byte) (Path, Set) {
	start := Path{}
	obstacles := make(Set)
	for r := range R {
		for c := range C {
			if grid[r][c] == '#' {
				obstacles.add(Pos{r, c})
			}
			dir := strings.IndexByte("^>v<", grid[r][c])
			if dir > -1 {
				start = Path{Pos{r, c}, dir}
			}
		}
	}
	return start, obstacles
}

func walkUntilOffTheGrid(start Path, obstacles Set) Set {
	path := start
	visited := make(Set)
	visited.add(path.Pos)
	for path.inside() {
		if obstacles.has(path.peek()) {
			path.turnRight()
			continue
		}
		path.moveForward()
		visited.add(path.Pos)
	}
	return visited
}

// --------------------------- PART 2 ------------------------
func p2(start Path, obstacles Set) int {
	// Beginning is same as PART 1
	visitedCells := walkUntilOffTheGrid(start, obstacles)

	nbrLoops := 0
	for cell := range visitedCells {
		obstacles.add(cell)
		if checkLoop(start, obstacles) {
			nbrLoops++
		}
		obstacles.delete(cell)
	}
	return nbrLoops
}

func checkLoop(start Path, obstacles Set) bool {
	path := start
	traveled := make(Set)
	traveled.add(path)

	for path.inside() {
		if obstacles.has(path.peek()) {
			path.turnRight()
			continue
		}
		path.moveForward()

		if traveled.has(path) {
			return true
		}
		traveled.add(path)
	}
	return false
}
