package main

import (
	"bytes"
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

	antennas := scan(grid)
	fmt.Printf("Part 1: %d\n", len(part1(antennas)))
	fmt.Printf("Part 2: %d\n", len(part2(antennas)))
}

type Vec2 struct{ r, c int }

func scan(grid [][]byte) map[byte][]Vec2 {
	antennas := map[byte][]Vec2{}
	for r := 0; r < R; r++ {
		for c := 0; c < C; c++ {
			ch := grid[r][c]
			if ch != '.' {
				antennas[ch] = append(antennas[ch], Vec2{r, c})
			}
		}
	}
	return antennas
}

func part1(antennas map[byte][]Vec2) map[Vec2]bool {
	// Use a map to avoid duplicates
	antinodes := map[Vec2]bool{}

	// Loop over all different antenna types
	for _, val := range antennas {
		// Double loop within the same antenna type
		for _, a := range val {
			for _, b := range val {
				// Skip the same antenna position
				if a == b {
					continue
				}
				// calculate antinode position with regards to the delta
				r := a.r - (b.r - a.r)
				c := a.c - (b.c - a.c)
				if 0 <= r && r < R && 0 <= c && c < C {
					antinodes[Vec2{r, c}] = true
				}
			}
		}
	}
	return antinodes
}

func part2(antennas map[byte][]Vec2) map[Vec2]bool {
	// Use a map to avoid duplicates
	antinodes := map[Vec2]bool{}

	// Loop over all different antenna types
	for _, val := range antennas {
		// Double loop within the same antenna type
		for _, a := range val {
			for _, b := range val {
				// Skip the same antenna position
				if a == b {
					continue
				}
				// calculate the deltas
				dr := (b.r - a.r)
				dc := (b.c - a.c)
				// Loop over all the possible positions
				for r, c := a.r, a.c; r < R && c < C && r >= 0 && c >= 0; r, c = r+dr, c+dc {
					antinodes[Vec2{r, c}] = true
				}
			}
		}
	}
	return antinodes
}

// ---------------------- Utility functions ----------------
func printGrid(grid [][]byte) {
	for r := 0; r < R; r++ {
		for c := 0; c < C; c++ {
			fmt.Printf("%c", grid[r][c])
		}
		fmt.Println()
	}
}

func printMap(m map[byte][]Vec2) {
	for key, val := range m {
		fmt.Printf("%c: %v\n", key, val)
	}
	fmt.Println()
}

func printAntinodes(antinodes map[Vec2]bool, grid [][]byte) {
	cpy := make([][]byte, len(grid))
	for r := 0; r < len(cpy); r++ {
		cpy[r] = make([]byte, len(grid[0]))
	}
	for r := 0; r < len(cpy); r++ {
		for c := 0; c < len(cpy[r]); c++ {
			cpy[r][c] = grid[r][c]
		}
	}
	for p := range antinodes {
		cpy[p.r][p.c] = '#'
	}
	printGrid(cpy)
}
