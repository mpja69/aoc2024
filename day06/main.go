package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

type Pos struct {
	r int
	c int
}

type Path struct {
	Pos
	d int
}

func (p *Path) next() byte {
	p.r += directions[p.d].r
	p.c += directions[p.d].c
	return grid[p.r][p.c]
}
func (p *Path) peek() byte {
	r := p.r + directions[p.d].r
	c := p.c + directions[p.d].c
	if r < 0 || r >= R || c < 0 || c >= R {
		return 0
	}
	return grid[r][c]
}
func (p *Path) turnRight() {
	p.d = (p.d + 1) % len(directions)
}

// ----------------------- Maps/Sets ------------------------
func has[T comparable](set map[T]bool, item T) bool {
	return set[item]
}
func add[T comparable](set map[T]bool, item T) {
	set[item] = true
}

// ------------------------ Globals -------------------------
var (
	R, C       int                                       // Boundries of grid
	grid       [][]byte                                  // the grid
	directions = []Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // 0..3 directions. ONLY turn right
)

// ------------- Instead of writing comments...--------------
func setObstacleOnGridPos(p Pos) {
	grid[p.r][p.c] = '#'
}
func removeObstacleFromGridPos(p Pos) {
	grid[p.r][p.c] = '.'
}

func main() {
	data, err := os.ReadFile("s.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	data = bytes.Trim(data, " \n")
	grid = bytes.Split(data, []byte{'\n'})

	R = len(grid)
	C = len(grid[0])

	fmt.Println("Part 1 (41, 4559):", p1())
	fmt.Println("Part 2 (6, 1604):", p2())
}

func p1() int {
	start, err := findStart()
	if err != nil {
		log.Fatal(err)
	}
	visitedPositions := walkUntilOffTheGrid(start)
	return len(visitedPositions)
}

func findStart() (Path, error) {
	for r := range R {
		for c := range C {
			dir := strings.IndexByte("^>v<", grid[r][c])
			if dir > -1 {
				return Path{Pos{r, c}, dir}, nil
			}
		}
	}
	return Path{}, fmt.Errorf("No Start pos found!")
}

func walkUntilOffTheGrid(start Path) map[Pos]bool {
	path := start
	visited := make(map[Pos]bool)
	add(visited, path.Pos)
	for path.peek() != 0 {
		if path.peek() == '#' {
			path.turnRight()
			continue
		}
		path.next()
		add(visited, path.Pos)
	}
	return visited
}

// --------------------------- PART 2 ------------------------
func p2() int {
	// >> Beginning is same as PART 1
	start, err := findStart()
	if err != nil {
		log.Fatal(err)
	}
	visitedPositions := walkUntilOffTheGrid(start)
	// << End of PART 1

	nbrLoops := 0
	for pos := range visitedPositions {
		setObstacleOnGridPos(pos)
		if checkLoop(start) {
			nbrLoops++
		}
		removeObstacleFromGridPos(pos)
	}
	return nbrLoops
}

func checkLoop(start Path) bool {
	path := start
	traveled := make(map[Path]bool)
	add(traveled, path)

	for path.peek() != 0 {
		if path.peek() == '#' {
			path.turnRight()
			continue
		}
		path.next()

		if has(traveled, path) {
			return true
		}
		add(traveled, path)
	}
	return false
}
