package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func part1(lines []string) int {
	return 0
}

func main() {
	// Get the data into strings
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	str := string(data)
	str = strings.Trim(str, " \n")
	lines := strings.Split(str, "\n")

	fmt.Println("Part 1:", part1(lines))
	// fmt.Println("Part 2:", part2(reports))
}
