package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func getNumbers(line string) (int, int) {
	parts := strings.Fields(line)
	left, _ := strconv.Atoi(parts[0])
	right, _ := strconv.Atoi(parts[1])
	return left, right
}

func getLists(lines []string) ([]int, []int) {
	// get the lines into 2 int arrays
	var ll []int
	var rr []int
	for _, line := range lines {
		left, right := getNumbers(line)
		ll = append(ll, left)
		rr = append(rr, right)
	}
	return ll, rr
}

func getDistance(a, b int) int {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d
}

// Similarity score
func getScore(r int, ll []int) int {
	s := 0
	for _, l := range ll {
		if r == l {
			s += r
		}
	}
	return s
}

func part1(ll, rr []int) int {
	sort.Ints(ll)
	sort.Ints(rr)

	dd := 0
	for i := range ll {
		dd += getDistance(ll[i], rr[i])
	}
	return dd
}

func part2(ll, rr []int) int {
	s := 0
	for i := range ll {
		s += getScore(ll[i], rr)
	}
	return s
}

func main() {
	// Get the data into strings
	data, _ := os.ReadFile("1.txt")
	str := string(data)
	str = strings.Trim(str, " \n")
	lines := strings.Split(str, "\n")

	ll, rr := getLists(lines)

	fmt.Println("Part 1:", part1(ll, rr))
	fmt.Println("Part 2:", part2(ll, rr))

}
