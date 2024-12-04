package main

import (
	"aoc2024/day03/lexer"
	"fmt"
	"log"
	"os"
	"strings"
)

func part1(lines []string) int {
	for _, line := range lines {
		l := lexer.New(line)
	}
	return 0
}

// "mul(XXX,XXX)", where XXX is 1-3 digits
func main() {
	// Get the data into strings
	data, err := os.ReadFile("s.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	str := string(data)
	str = strings.Trim(str, " \n")
	lines := strings.Split(str, "\n")

	fmt.Println("Part 1:", part1(lines), " => mul(2,4), mul(5,5), mul(11,8), mul(8,4) ", " => 161 ?")
	// fmt.Println("Part 2:", part2(reports))
}
