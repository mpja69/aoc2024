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
	directions = []Vec{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // 0..3 directions. ONLY turn right
)

type Vec struct {
	r int
	c int
}

func (p Vec) getPos() Vec {
	return p
}

type Path struct {
	pos Vec
	dir Vec
}

func (p Path) getPos() Vec {
	return p.pos
}
func (p *Path) peek() Vec {
	r := p.pos.r + p.dir.r
	c := p.pos.c + p.dir.c
	return Vec{r, c}
}
func (p *Path) moveForward() {
	p.pos.r += p.dir.r
	p.pos.c += p.dir.c
}
func (p *Path) inside() bool {
	pp := p.peek()
	return !(pp.r < 0 || pp.r >= R || pp.c < 0 || pp.c >= R)
}

func (p *Path) turnRight() {
	p.dir.r, p.dir.c = p.dir.c, -p.dir.r
}

// ----------------------- Maps/Sets ------------------------
type Item interface {
	getPos() Vec
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
	visited := walkUntilOffTheGrid(start, obstacles)
	fmt.Println("Part 1 (41, 4559):", p1(visited))
	fmt.Println("Part 2 (6, 1604):", p2(start, obstacles, visited))
	fmt.Printf("Millis: %d\n", time.Since(t).Milliseconds())
}

func p1(visited Set) int {
	return len(visited)
}

func findStart(grid [][]byte) (Path, Set) {
	start := Path{}
	obstacles := make(Set)
	for r := range R {
		for c := range C {
			if grid[r][c] == '#' {
				obstacles.add(Vec{r, c})
			}
			dir := strings.IndexByte("^>v<", grid[r][c])
			if dir > -1 {
				start = Path{Vec{r, c}, directions[dir]}
			}
		}
	}
	return start, obstacles
}

func walkUntilOffTheGrid(start Path, obstacles Set) Set {
	path := start
	visited := make(Set)
	visited.add(path.pos)
	for path.inside() {
		if obstacles.has(path.peek()) {
			path.turnRight()
			continue
		}
		path.moveForward()
		visited.add(path.pos)
	}
	return visited
}

// --------------------------- PART 2 ------------------------
func p2(start Path, obstacles Set, visited Set) int {
	nbrLoops := 0
	for cell := range visited {
		obstacles.add(cell)
		// NOTE:	start kan få vara "sista/senaste path:en innan hinder"...
		//			För stegen innan är ju redan testade!
		//				=> Behöver fixa detta i checkloop + att ändra start här!!
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
