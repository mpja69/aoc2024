package main

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"maps"
	"math"
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

	fmt.Printf("Part 1: (1930, 1550156): %d\n", part1(grid))
	fmt.Printf("Part 2: (1206, 946084): %d\n", part2(grid))
}

type Pos struct{ r, c int }
type Grid [][]byte

type Region struct {
	plant byte
	plots []Pos
}

func part1(g Grid) int {
	visited := make(map[Pos]bool)
	price := 0
	for r := range R {
		for c := range C {
			if visited[Pos{r, c}] {
				continue
			}
			region, area, perimeter := g.findOneRegion(Pos{r, c})
			maps.Insert(visited, maps.All(region))
			price += area * perimeter
		}
	}
	return price
}

// Breadth First search - Find all positions in a one region
func (g Grid) findOneRegion(start Pos) (visited map[Pos]bool, area, perimeter int) {
	q := list.New()          // Use a queue for the upcoming steps
	visited = map[Pos]bool{} // Use a set to avoid duplicate visits

	// Initial values
	q.PushBack(start)

	// Process the queue
	for q.Len() > 0 {

		// Pop the oldest item
		pos := q.Front().Value.(Pos)
		q.Remove(q.Front())

		// Track visited positions in a set
		if visited[pos] {
			continue
		}
		visited[pos] = true

		// Count the things we want
		neighbours := g.neighbours(pos)
		perimeter += 4 - len(neighbours)
		area += 1

		// Push all possible neighbours to the queue
		for _, neighbour := range neighbours {
			q.PushBack(neighbour)
		}
	}
	return visited, area, perimeter
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
		// Check valid plant
		if g[r][c] != val {
			continue
		}
		neighbours = append(neighbours, Pos{r, c})
	}
	return neighbours
}

// ---------------------------------------------- Part 2 ------------------------------------------------

func part2(g Grid) int {
	visited := make(map[Pos]bool)
	price := 0
	for r := range R {
		for c := range C {
			if visited[Pos{r, c}] {
				continue
			}
			region, area, _ := g.findOneRegion(Pos{r, c})
			maps.Insert(visited, maps.All(region))

			// --- Solution where we count the corners ---
			sides := corners(region)
			price += area * sides

			// --- Solution where we scan the lines ---
			// horizontalSides := scanRowsRegion(region)
			// verticalSides := scanColsRegion(region)
			// perimeter := horizontalSides + verticalSides
			// price += area * perimeter
		}
	}
	return price
}

type F2 struct {
	r, c float64
}

func corners(region map[Pos]bool) int {
	virtualGrid := map[F2]bool{}
	offsets := []F2{{-0.5, 0.5}, {0.5, 0.5}, {0.5, -0.5}, {-0.5, -0.5}}

	// Create a grid for the virtual corner points
	for p := range region {
		for _, o := range offsets {
			r := float64(p.r) + o.r
			c := float64(p.c) + o.c
			virtualGrid[F2{r, c}] = true
		}
	}

	corners := 0

	for c := range virtualGrid {
		b := [4]bool{}
		for i, o := range offsets {
			r := int(math.Round(c.r + o.r))
			c := int(math.Round(c.c + o.c))
			b[i] = region[Pos{r, c}]
		}

		if b[0] != b[1] != b[2] != b[3] {
			corners += 1
		} else if b == [4]bool{true, false, true, false} || b == [4]bool{false, true, false, true} {
			corners += 2
		}
	}
	return corners
}

// ------------------- Scan funcions ---------------------
func scanRowsRegion(region map[Pos]bool) int {
	length := 0
	// Scan each row...
	for r := range R {
		above := false
		below := false
		// ...from left to right
		for c := range C {
			pos := Pos{r, c}

			// First check if we are in region
			if region[pos] == false {
				above = false
				below = false
				continue
			}

			// Check for edge above
			if !region[pos.above()] && !above {
				length++
			}

			// Check for edge below
			if !region[pos.below()] && !below {
				length++
			}

			// Set flags
			above = !region[pos.above()]
			below = !region[pos.below()]
		}
	}
	return length
}

func scanColsRegion(region map[Pos]bool) int {
	length := 0
	// Scan each col...
	for c := range C {
		left := false
		right := false
		// ...from top to bottom
		for r := range R {
			pos := Pos{r, c}

			// First check if we are in region
			if region[pos] == false {
				left = false
				right = false
				continue
			}

			// Check for edge above: I.e. NOT a_square_above AND NOT in_a_ongoing_line_above
			if !region[pos.left()] && !left {
				length++
			}

			// Check for edge below
			if !region[pos.right()] && !right {
				length++
			}

			// Set flags
			left = !region[pos.left()]
			right = !region[pos.right()]
		}
	}
	return length
}

// ---------------------- Helper functions --------------------
// --- Grid ---
func (g Grid) inside(r, c int) bool {
	return r >= 0 && r <= R-1 && c >= 0 && r <= C-1
}

func (g Grid) set(pos Pos, ch byte) {
	if g.inside(pos.r, pos.c) {
		g[pos.r][pos.c] = ch
	}
}
func (g Grid) printGrid() {
	for r := range R {
		for c := range C {
			fmt.Printf("%c", g[r][c])
		}
		fmt.Println()
	}
	fmt.Println()
}
func regionToGrid(region map[Pos]bool) Grid {
	grid := [][]byte{}
	for r := range R {
		grid = append(grid, []byte{})
		for c := range C {
			if region[Pos{r, c}] {
				grid[r] = append(grid[r], 'x')
			} else {
				grid[r] = append(grid[r], '.')
			}
		}
	}
	return grid
}

// --- Pos ---
// --- methods ---
func (p Pos) above() Pos {
	return above(p)
}
func (p Pos) below() Pos {
	return below(p)
}
func (p Pos) left() Pos {
	return left(p)
}
func (p Pos) right() Pos {
	return right(p)
}

// --- Functions ---
type Move func(Pos) Pos

func above(p Pos) Pos {
	return Pos{p.r - 1, p.c}
}
func below(p Pos) Pos {
	return Pos{p.r + 1, p.c}
}
func left(p Pos) Pos {
	return Pos{p.r, p.c - 1}
}
func right(p Pos) Pos {
	return Pos{p.r, p.c + 1}
}
