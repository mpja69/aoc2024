package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	data, err := os.ReadFile("s.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := strings.Trim(string(data), " \n")
	lines := strings.Split(input, "\n")

	fmt.Println("Part 1 (, ):", p1(lines))
}

func p1(lines []string) int {
	res := 0

	return res
}
