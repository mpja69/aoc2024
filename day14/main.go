package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type P struct {
	x, y int
}
type Robot struct {
	p, v P
}

func part2(robots []Robot) int {
	res := 0
	W := 101 // 11
	H := 103 // 7

	g := robotsToGrid(robots, 0, 0, W, H)

	// Since the speed is constant, it's just a linear equation
	for {
		for i := range robots {
			// Move each robot
			g.set(robots[i].p, ' ')
			robots[i].p.x += robots[i].v.x
			robots[i].p.y += robots[i].v.y

			// Modulo, both positive and negative
			for robots[i].p.x < 0 {
				robots[i].p.x += W
			}
			for robots[i].p.y < 0 {
				robots[i].p.y += H
			}
			for robots[i].p.x >= W {
				robots[i].p.x -= W
			}
			for robots[i].p.y >= H {
				robots[i].p.y -= H
			}

			g.set(robots[i].p, '#')
		}
		res++
		if g.line(10, '#') {
			break
		}
		if res > W*H {
			return 0
		}
	}
	g.printGrid()
	return res
}

func part1(robots []Robot) int {
	res := 0
	W := 101 // 11
	H := 103 // 7
	moves := 100

	// Since the speed is constant, it's just a linear equation
	for i := range robots {
		// Move each robot 100 times
		robots[i].p.x += robots[i].v.x * moves
		robots[i].p.y += robots[i].v.y * moves

		// Wrap around
		robots[i].p.x %= W
		robots[i].p.y %= H

		// (Modulo can be negative...so adjust)
		if robots[i].p.x < 0 {
			robots[i].p.x += W
		}
		if robots[i].p.y < 0 {
			robots[i].p.y += H
		}
	}

	// Count nuber of robots in each quadrant, (exclude the middle lines)
	xc := W / 2
	yc := H / 2
	var topleft, topright, bottomleft, bottomright int

	for _, r := range robots {
		if r.p.x < xc && r.p.y < yc {
			topleft++
		}
		if r.p.x > xc && r.p.y < yc {
			topright++
		}
		if r.p.x < xc && r.p.y > yc {
			bottomleft++
		}
		if r.p.x > xc && r.p.y > yc {
			bottomright++
		}
	}
	// and multiply together
	res = topleft * topright * bottomleft * bottomright
	return res
}
func main() {
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	robots := []Robot{}

	input := strings.Trim(string(data), " \n")
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		c := parseLine(line)
		robots = append(robots, c)
	}

	// fmt.Printf("Part 1: (12, 229069152) %d\n", part1(robots))
	fmt.Printf("Part 2: (, 7383) %d\n", part2(robots))
}

func parseLine(str string) Robot {
	re := regexp.MustCompile(`(-?\d+)`)
	matches := re.FindAllString(str, -1)
	r := Robot{}
	r.p.x, _ = strconv.Atoi(matches[0])
	r.p.y, _ = strconv.Atoi(matches[1])
	r.v.x, _ = strconv.Atoi(matches[2])
	r.v.y, _ = strconv.Atoi(matches[3])
	return r
}

type Grid [][]byte

func (g Grid) set(pos P, ch byte) {
	g[pos.y][pos.x] = ch
}

func (g Grid) printGrid() {
	for y := range len(g) {
		for x := range len(g[0]) {
			fmt.Printf("%c", g[y][x])
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g Grid) line(n int, ch byte) bool {
	for y := range len(g) {
		length := 0
		for x := range len(g[0]) {
			if g[y][x] == ch {
				length++
			} else {
				length = 0
			}
			if length == n {
				return true
			}
		}
	}
	return false
}
func robotsToGrid(robots []Robot, x0, y0, x1, y1 int) Grid {
	W := x1 - x0
	H := y1 - y0
	g := make([][]byte, H)
	for y := range H {
		g[y] = make([]byte, W)
		for x := range W {
			g[y][x] = ' '
		}
	}

	for _, r := range robots {
		if r.p.x < x0 || r.p.x >= x1 || r.p.y < y0 || r.p.y >= y1 {
			continue
		}
		x := r.p.x - x0
		y := r.p.y - y0
		g[y][x] = '#'
	}
	return g
}
