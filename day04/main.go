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

	res := part1(lines)
	fmt.Println("Part 1 (2613):", res)

}

type rc struct {
	r int
	c int
}
type rcdir struct {
	pos rc
	dir rc
}

func part1(grid []string) int {
	res := make(map[rcdir]bool)
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

	// Check all positions
	R := len(grid)
	C := len(grid[0])
	for r := 0; r < R; r++ {
		for c := 0; c < C; c++ {
			checkAllDir(grid, rc{r, c}, dirs, res)
		}
	}
	return len(res)
}

func checkAllDir(grid []string, pos rc, dirs []rc, res map[rcdir]bool) {
	for _, dir := range dirs {
		if checkDir(grid, pos, dir) {
			res[rcdir{pos: pos, dir: dir}] = true
		}
	}
}

func checkDir(grid []string, pos, dir rc) bool {
	R := len(grid)
	C := len(grid[0])
	r := pos.r
	c := pos.c
	// i := 0

	for i := 0; i < 4; i++ {
		// Check boundries
		if r < 0 || r >= R || c < 0 || c >= C {
			return false
		}
		// Check string
		if "XMAS"[i] != grid[r][c] {
			return false
		}
		// Move 1 step, (in the given direction)
		r += dir.r
		c += dir.c
	}
	return true
}

func printMap(grid []string, sol map[rcdir]bool) {
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
