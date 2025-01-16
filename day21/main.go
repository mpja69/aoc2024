package main

import (
	"aoc2024/day21/types"
	"fmt"
	"os"
	"strconv"

	"github.com/mpja69/aoc2024/day21/numpad"
)

type P = types.P

func main() {
	data, _ := os.ReadFile("sample.txt")
	fmt.Printf("Data: %v\n", data)
	np := numpad.New()
	fmt.Printf("%v\n", np.PeekTo("5"))

}

var dir2dir = map[string]string{
	"A": "A",
	"^": "<A>A",
	">": "vA^A",
	"v": "v<A>^A",
	"<": "v<<A>>^A",
}

func complexity(code string, sequence string) int {
	c, _ := strconv.Atoi(code[:3])
	s := len(sequence)
	return c * s
}
