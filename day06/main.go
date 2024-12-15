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
	R, C int // Boundries of grid
)

type Vec2 struct {
	r int
	c int
}

type Step struct {
	pos Vec2
	dir Vec2
}

func (p *Step) peek() Vec2 {
	r := p.pos.r + p.dir.r
	c := p.pos.c + p.dir.c
	return Vec2{r, c}
}
func (p *Step) moveForward() {
	p.pos.r += p.dir.r
	p.pos.c += p.dir.c
}
func (p *Step) inside() bool {
	pp := p.peek()
	return !(pp.r < 0 || pp.r >= R || pp.c < 0 || pp.c >= R)
}

func (p *Step) turnRight() {
	p.dir.r, p.dir.c = p.dir.c, -p.dir.r
}

// ----------------------- Maps/Sets ------------------------
type Item interface{}
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
func has[T comparable](set map[T]bool, item T) bool {
	return set[item]
}
func add[T comparable](set map[T]bool, item T) {
	set[item] = true
}
func del[T comparable](set map[T]bool, item T) {
	set[item] = false
}

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
	visited := walkUntilOffTheGrid(start, obstacles)
	fmt.Println("Part 1 (41, 4559):", len(visited))
	fmt.Println("Part 2 (6, 1604):", p2(obstacles, visited))
	fmt.Printf("Millis: %d\n", time.Since(t).Milliseconds())
}

// --------------------------- PART 2 ------------------------
func p2(obstacles Set, visitedList []Step) int {
	nbrLoops := 0
	isLoop := false
	start := visitedList[0]
	for _, newStart := range visitedList[1:] {
		newObstacle := newStart.pos
		obstacles.add(newObstacle)
		if isLoop = checkLoop(start, obstacles); isLoop {
			nbrLoops++
		}
		obstacles.delete(newObstacle)
		start = newStart
	}
	return nbrLoops
}

func findStart(grid [][]byte) (Step, Set) {
	start := Step{}
	obstacles := make(Set)
	for r := range R {
		for c := range C {
			if grid[r][c] == '#' {
				obstacles.add(Vec2{r, c})
			}
			dir := strings.IndexByte("^>v<", grid[r][c])
			if dir > -1 {
				start = Step{Vec2{r, c}, []Vec2{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}[dir]}
			}
		}
	}
	return start, obstacles
}

func walkUntilOffTheGrid(start Step, obstacles Set) []Step {
	path := start
	list := []Step{start}
	visited := make(Set)
	visited.add(start.pos)
	for path.inside() {
		if obstacles.has(path.peek()) {
			path.turnRight()
			continue
		}
		path.moveForward()
		if visited.has(path.pos) == false {
			visited.add(path.pos)
			list = append(list, path)
		}
	}
	return list
}

func checkLoop(start Step, obstacles Set) bool {
	path := start
	visited := make(Set)
	visited.add(path) // 23 ggr:w

	for path.inside() {
		if obstacles.has(path.peek()) {
			path.turnRight()
			continue
		}
		path.moveForward()

		if visited.has(path) {
			return true
		}
		visited.add(path)
	}
	return false
}
