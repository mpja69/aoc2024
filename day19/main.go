package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	data, _ := os.ReadFile("data.txt")
	wantedPatterns := strings.Split(strings.TrimSpace(string(data)), "\n")

	availablePatterns := strings.Split(wantedPatterns[0], ", ")
	wantedPatterns = wantedPatterns[2:]

	fmt.Printf("P1: (340) \n")
	p1(availablePatterns, wantedPatterns)
	fmt.Printf("P2:(717561822679428)\n")
	p2(availablePatterns, wantedPatterns)
}

func p1(available, wanted []string) {
	memo := map[string]bool{}
	resRec := 0
	resDP := 0
	stopTimer := startTimer()
	for _, pattern := range wanted {
		if canConstructRec(pattern, available, memo) {
			resRec++
		}
	}
	fmt.Printf("P1: Recusrsive solution: %d, %v\n", resRec, stopTimer())
	stopTimer = startTimer()
	for _, pattern := range wanted {
		if canConstructDP(pattern, available) {
			resDP++
		}
	}
	fmt.Printf("P1: Table solution:      %d, %v\n", resDP, stopTimer())
}

func p2(available, wanted []string) {
	memo := map[string]int{}
	resRec := 0
	resDP := 0
	stopTimer := startTimer()
	for _, pattern := range wanted {
		resRec += countConstructRec(pattern, available, memo)
	}
	fmt.Printf("P2: Recusrsive solution: %d, %v\n", resRec, stopTimer())

	stopTimer = startTimer()
	for _, pattern := range wanted {
		resDP += countConstructDP(pattern, available)
	}
	fmt.Printf("P2: Table solution:      %d, %v\n", resDP, stopTimer())
}

// Dynamic programming with a table
// Time complexity: O(m^2*n), where m is target length, and n is the length of the word bank
// Space complexity: O(m)
func canConstructDP(target string, wordBank []string) bool {
	table := make([]bool, len(target)+1) // Table represent how much of the word that can be built...up to (not including) each position
	table[0] = true                      // The empty string ("") CAN be built -> true in first position

	for i := range table {
		if table[i] {
			for _, word := range wordBank {
				if strings.HasPrefix(target[i:], word) {
					table[i+len(word)] = true
				}
			}
		}
	}
	return table[len(target)]
}

// Like canConstruct, but counting and adding up all ways the target can be constructed
func countConstructDP(target string, wordBank []string) int {
	table := make([]int, len(target)+1)
	table[0] = 1

	for i := range table {
		if table[i] > 0 {
			for _, word := range wordBank {
				if strings.HasPrefix(target[i:], word) {
					table[i+len(word)] += table[i]
				}
			}
		}
	}
	return table[len(target)]
}

func canConstructRec(target string, wordBank []string, memo map[string]bool) bool {
	if len(target) == 0 {
		return true
	}
	if memo[target] {
		return true
	}
	for _, word := range wordBank {
		if strings.HasPrefix(target, word) {
			if canConstructRec(target[len(word):], wordBank, memo) {
				memo[target] = true
				return true
			}
		}
	}
	return false
}

func countConstructRec(target string, wordBank []string, memo map[string]int) int {
	if len(target) == 0 {
		return 1
	}
	if val, ok := memo[target]; ok {
		return val
	}
	val := 0
	for _, word := range wordBank {
		if strings.HasPrefix(target, word) {
			val += countConstructRec(target[len(word):], wordBank, memo)
		}
	}
	memo[target] = val
	return val
}

// Starts a timer and return a function that returns the elapsed time.
func startTimer() func() time.Duration {
	t0 := time.Now()
	return func() time.Duration {
		return time.Duration(time.Since(t0).Milliseconds())
	}
}
