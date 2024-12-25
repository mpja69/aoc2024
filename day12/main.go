package main

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"maps"
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
			Region, area, perimeter := g.findOneRegion(Pos{r, c})
			maps.Insert(visited, maps.All(Region))
			price += area * perimeter
		}
	}
	return price
}

// Breadth First search - Find all positions in a one region
func (g Grid) findOneRegion(start Pos) (visited map[Pos]bool, area int, perimeter int) {
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
			id := g[r][c]
			horizontalSides := g.scanRows(region, id)
			verticalSides := g.scanCols(region, id)
			perimeter := horizontalSides + verticalSides
			price += area * perimeter
		}
	}
	return price
}

func (g Grid) scanRows(region map[Pos]bool, id byte) int {
	length := 0
	// Scan each row...
	for r := range R {
		top := false
		bottom := false
		// ...from left to right
		for c := range C {
			if region[Pos{r, c}] == false {
				top = false
				bottom = false
				continue
			}
			// Have one above
			if r > 0 && g[r-1][c] == id {
				top = false
			}
			// Have one below
			if r < R-1 && g[r+1][c] == id {
				bottom = false
			}
			// Have none above
			if r == 0 || r > 0 && g[r-1][c] != id {
				if top == false {
					length++
					top = true
				}
			}
			// have none below
			if r == R-1 || r < R-1 && g[r+1][c] != id {
				if bottom == false {
					length++
					bottom = true
				}
			}
		}
	}
	return length
}

func (g Grid) scanCols(region map[Pos]bool, id byte) int {
	length := 0
	// Scan each col...
	for c := range C {
		left := false
		right := false
		// ...from top to bottom
		for r := range R {
			if region[Pos{r, c}] == false {
				left = false
				right = false
				continue
			}
			// have one left
			if c > 0 && g[r][c-1] == id {
				left = false
			}
			// have one right
			if c < C-1 && g[r][c+1] == id {
				right = false
			}
			// have none left
			if c == 0 || c > 0 && g[r][c-1] != id {
				if left == false {
					length++
					left = true
				}
			}
			// have none right
			if c == C-1 || c < C-1 && g[r][c+1] != id {
				if right == false {
					length++
					right = true
				}
			}
		}
	}
	return length
}

// ---------------------- Helper functions --------------------
func printGrid(g [][]byte) {
	for r := range R {
		for c := range C {
			fmt.Printf("%c", g[r][c])
		}
		fmt.Println()
	}
}

func (p Pos) add(d Pos) Pos {
	return Pos{p.r + d.r, p.c + d.c}
}
