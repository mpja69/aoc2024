package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/mpja69/aoc2024/day21/keypad"
)

func main() {
	data, _ := os.ReadFile("sample.txt")
	data = bytes.TrimSpace(data)
	first := bytes.Split(data, []byte("\n"))
	// fmt.Println("029A")
	fmt.Printf("%s\n\n", data)

	// Create the Numpad, with the desired output data we want
	np := keypad.NewNumpad(data)
	numBuf := make([]byte, 1000)
	n, _ := np.Read(numBuf)

	// fmt.Println("<A^A>^^AvvvA")
	fmt.Printf("%s\n\n", numBuf[:n])

	dp := keypad.NewDirpad(numBuf[:n])
	dirBuf := make([]byte, 1000)
	n, _ = dp.Read(dirBuf)

	// fmt.Println("v<<A>>^A<A>AvA<^AA>A<vAAA>^A")
	fmt.Printf("%s\n\n", dirBuf[:n])

	dp = keypad.NewDirpad(dirBuf[:n])
	dirBuf = make([]byte, 1000)
	n, _ = dp.Read(dirBuf)

	// fmt.Println("<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A")
	fmt.Printf("%s\n\n", dirBuf[:n])

	sum := 0
	for i, line := range bytes.Split(dirBuf, []byte("\n")) {
		fmt.Printf("%s: %s\n", first[i], line)
		sum += complexity(first[i], line)
		// num, _ := strconv.Atoi(string(first[0][:3]))
		// fmt.Printf("%d * %d = %d\n", n, num, n*num)
	}
	fmt.Printf("Complexity: %d\n", sum)
}

func complexity(code []byte, sequence []byte) int {
	s := len(sequence)
	c, _ := strconv.Atoi(string(code[:3]))
	fmt.Printf("%d * %d = %d\n", s, c, s*c)
	return c * s
}
