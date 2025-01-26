package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var cache map[string]int
var expr map[string]string

func main() {
	data, _ := os.ReadFile("data.txt")
	parts := strings.Split(strings.TrimSpace(string(data)), "\n\n")
	cache = map[string]int{}
	// The inputs
	for _, line := range strings.Split(parts[0], "\n") {
		l := strings.Split(line, ": ")
		val, _ := strconv.Atoi(l[1])
		cache[l[0]] = val
	}

	// The expressions
	expr = map[string]string{}
	for _, line := range strings.Split(parts[1], "\n") {
		l := strings.Split(line, " -> ")
		expr[l[1]] = l[0]
	}
	p1()
}

func p1() {
	for k := range expr {
		solve(k)
	}
	fmt.Println("P1: ", output())
}

func solve(ans string) int {
	if val, ok := cache[ans]; ok {
		return val
	}

	a, op, b := parse(expr[ans])
	if _, ok := cache[a]; !ok {
		cache[a] = solve(a)
	}

	if _, ok := cache[b]; !ok {
		cache[b] = solve(b)
	}

	switch op {
	case "AND":
		cache[ans] = cache[a] & cache[b]
	case "OR":
		cache[ans] = cache[a] | cache[b]
	case "XOR":
		cache[ans] = cache[a] ^ cache[b]
	}

	return cache[ans]
}

func parse(e string) (a, op, b string) {
	s := strings.Split(e, " ")
	return s[0], s[1], s[2]

}
func output() int {
	z := []string{}

	// Find all starting with z
	for k := range cache {
		if k[0] != 'z' {
			continue
		}
		z = append(z, k)
	}

	// order them
	slices.Sort(z)

	// Calculate value
	res := 0
	for _, zz := range slices.Backward(z) {
		res = res | cache[zz]
		res = res << 1
	}
	return res >> 1
}
