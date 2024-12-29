package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := bytes.Trim(data, " \n")
	parts := bytes.Split(input, []byte{'\n', '\n'})
	grid = bytes.Split(parts[0], []byte{'\n'})
	moves := parts[1]

	R = len(grid)
	C = len(grid[0])

	fmt.Printf("P1: %d\n", p1(moves))

	// small: 2028
	// large: 10092
}

var (
	R, C int
	grid Grid
)

type P struct{ r, c int }
type Grid [][]byte

func p1(moves []byte) int {
	p := start()
	for _, m := range moves {
		dir, ch := peek(p, m)
		switch ch {
		case '.':
			p = move(p, dir)
		case 'O':
			boxes, ok := checkBoxes(p, dir)
			if ok {
				pushBoxes(p, dir, boxes)
				p = move(p, dir)
			}
		}
	}
	return sumGPS()
}

func peek(p P, m byte) (d P, c byte) {
	move := map[byte]P{'^': {-1, 0}, '>': {0, 1}, 'v': {1, 0}, '<': {0, -1}}
	d = move[m]
	n := p
	n.r, n.c = p.r+d.r, p.c+d.c
	c = grid[n.r][n.c]
	return
}

func move(p, d P) P {
	n := p
	n.r, n.c = p.r+d.r, p.c+d.c
	grid[n.r][n.c] = '@'
	grid[p.r][p.c] = '.'
	return n
}

func checkBoxes(p, d P) (boxes int, ok bool) {
	boxes = 0
	for {
		p.r, p.c = p.r+d.r, p.c+d.c
		if grid[p.r][p.c] != 'O' {
			break
		}
		boxes++
	}
	ok = grid[p.r][p.c] == '.'
	return
}

func pushBoxes(p, d P, boxes int) {
	r, c := p.r, p.c
	r, c = r+d.r, c+d.c
	grid[r][c] = '.'
	for i := 0; i < boxes; i++ {
		r, c = r+d.r, c+d.c
		grid[r][c] = 'O'
	}
}

func start() P {
	for r := range R {
		for c := range C {
			if grid[r][c] == '@' {
				return P{r, c}
			}
		}
	}
	return P{}
}

func sumGPS() int {
	sum := 0
	for r := range R {
		for c := range C {
			if grid[r][c] == 'O' {
				sum += 100*r + c
			}
		}
	}
	return sum
}

func printGrid() {
	for r := range R {
		for c := range C {
			fmt.Printf("%c", grid[r][c])
		}
		fmt.Println()
	}
	fmt.Println()
}
