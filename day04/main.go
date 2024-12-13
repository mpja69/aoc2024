package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Get the data into strings
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	str := string(data)
	str = strings.Trim(str, " \n")
	lines := strings.Split(str, "\n")

	fmt.Println("Part 1 (18, 2613):", part1(lines))
	fmt.Println("Part 2 (9, 1905):", part2(lines))

}

type rc struct {
	r int
	c int
}
type key struct {
	pos rc
	dir rc
}

func part2(grid []string) int {
	res := 0
	// Check all VALID positions, i.e the ones where the squares overlap
	R := len(grid) - 1
	C := len(grid[0]) - 1
	for r := 1; r < R; r++ {
		for c := 1; c < C; c++ {
			if grid[r][c] == 'A' {
				res += checkAllPatterns(grid, rc{r, c})
			}
		}
	}
	return res
}

func checkAllPatterns(grid []string, pos rc) int {
	patterns := []string{"MMSS", "SMMS", "SSMM", "MSSM"}
	for _, pattern := range patterns {
		if checkCorners(grid, pos, pattern) {
			return 1
		}
	}
	return 0
}

func checkCorners(grid []string, pos rc, pattern string) bool {
	offsets := []rc{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}}
	for i, off := range offsets {
		if grid[pos.r+off.r][pos.c+off.c] != pattern[i] {
			return false
		}
	}
	return true
}

func part1(grid []string) int {
	res := make(map[key]bool)

	// Check all positions
	R := len(grid)
	C := len(grid[0])
	for r := 0; r < R; r++ {
		for c := 0; c < C; c++ {
			checkAllDirections(grid, rc{r, c}, res)
		}
	}
	return len(res)
}

func checkAllDirections(grid []string, pos rc, res map[key]bool) {
	dirs := []rc{
		{-1, 0},
		{1, 0},
		{0, 1},
		{0, -1},
		{-1, 1},
		{-1, -1},
		{1, 1},
		{1, -1},
	}
	for _, dir := range dirs {
		if checkPattern(grid, pos, dir) {
			res[key{pos: pos, dir: dir}] = true
		}
	}
}

func checkPattern(grid []string, pos, dir rc) bool {
	R := len(grid)
	C := len(grid[0])
	r := pos.r
	c := pos.c
	pattern := "XMAS"

	for i := 0; i < 4; i++ {
		// Check boundries
		if r < 0 || r >= R || c < 0 || c >= C {
			return false
		}
		// Check string
		if pattern[i] != grid[r][c] {
			return false
		}
		// Move 1 step, (in the given direction)
		r += dir.r
		c += dir.c
	}
	return true
}

// --------------------- utility func ----------------------
func printMap(grid []string, sol map[key]bool) {
	R := len(grid)
	C := len(grid[0])
	res := make([][]byte, R)
	for r := range grid {
		res[r] = make([]byte, C)
	}

	for s := range sol {
		println("Sol: ")
		for r := range res {
			for c := range res[r] {
				res[r][c] = '.'
			}
		}

		r := s.pos.r
		c := s.pos.c
		res[r][c] = 'X'
		for i := 0; i < 3; i++ {
			r += s.dir.r
			c += s.dir.c
			res[r][c] = "MAS"[i]
		}

		for r := range res {
			fmt.Printf("%s\n", res[r])
		}
		println()
	}
}
