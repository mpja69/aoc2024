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
	fmt.Printf("%s\n", data)

	// Create the Numpad, with the desired output data we want
	np := numpad.New(data)
	buf := make([]byte, 1000)
	// Playing around with the Reader interface: Read each line, returning the sequence into the buffer
	np.Read(buf)
	// Print the buffer as a string to visualize the sequences
	fmt.Printf("%s\n", buf)

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
