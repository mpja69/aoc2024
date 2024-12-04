package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getNumbers(line string) []int {
	parts := strings.Fields(line)
	var numbers []int
	for i := range parts {
		number, _ := strconv.Atoi(parts[i])
		numbers = append(numbers, number)
	}
	return numbers
}

func getReports(lines []string) [][]int {
	var reports [][]int
	for _, line := range lines {
		numbers := getNumbers(line)
		reports = append(reports, numbers)
	}
	return reports
}

func isIncreasing(report []int) bool {
	ok := true
	for i := 1; i < len(report); i++ {
		distance := report[i] - report[i-1]
		if distance < 1 || distance > 3 {
			return false
		}
	}
	return ok
}

func isDecreasing(report []int) bool {
	ok := true
	for i := 1; i < len(report); i++ {
		distance := report[i-1] - report[i]
		if distance < 1 || distance > 3 {
			return false
		}
	}
	return ok
}
func problemDampener(report []int, i int) []int {
	cpy := make([]int, len(report))
	copy(cpy, report)
	return append(cpy[:i], cpy[i+1:]...)
}

func isSafe(report []int) bool {
	for i := range report {
		r := problemDampener(report, i)
		if isIncreasing(r) || isDecreasing(r) {
			return true
		}
	}
	return false
}

func part1(reports [][]int) int {
	nbrSafeReports := 0
	for _, report := range reports {
		if isIncreasing(report) || isDecreasing(report) {
			nbrSafeReports++
		}
	}
	return nbrSafeReports
}

func part2(reports [][]int) int {
	nbrSafeReports := 0
	for _, report := range reports {
		if isSafe(report) {
			nbrSafeReports++
		}
	}
	return nbrSafeReports
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

	reports := getReports(lines)

	fmt.Println("Part 1:", part1(reports))
	fmt.Println("Part 2:", part2(reports))
}
