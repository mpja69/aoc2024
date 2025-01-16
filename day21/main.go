package main

import (
	"fmt"
	"github.com/mpja69/aoc2024/day21/numpad"
	"os"
	"strconv"
)

func main() {
	data, _ := os.ReadFile("sample.txt")
	fmt.Printf("Data: %v\n", data)
	numpad := NewNumpad()
	fmt.Printf("%v\n", numpad.move2seq)

}

var (
	dir2sym = map[P]string{
		{0, 1}:  ">",
		{1, 0}:  "v",
		{0, -1}: "<",
		{-1, 0}: "^",
	}
	sym2dir = map[string]P{
		">": {0, 1},
		"v": {1, 0},
		"<": {0, -1},
		"^": {-1, 0},
	}
)

// func dirTodir(key string) string {
var dir2dir = map[string]string{
	"A": "A",
	"^": "<A>A",
	">": "vA^A",
	"v": "v<A>^A",
	"<": "v<<A>>^A",
}

// 	return keyMap[key]
// }

func complexity(code string, sequence string) int {
	c, _ := strconv.Atoi(code[:3])
	s := len(sequence)
	return c * s
}
