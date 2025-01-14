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

	fmt.Printf("P1: Finding how many of the patterns that can be constructed (340) \n")
	p1(availablePatterns, wantedPatterns)
	fmt.Printf("P2: Finding how many ways the patterns can be constructed (717561822679428)\n")
	p2(availablePatterns, wantedPatterns)
}

func p1(available, wanted []string) {
	// Top-Down approach
	memo := map[string]bool{}
	result := 0
	stopTimer := startTimer()
	for _, pattern := range wanted {
		if canConstructTD(pattern, available, memo) {
			result++
		}
	}
	fmt.Printf("P1: Top-Down solution:  %d, %v\n", result, stopTimer())
	// Bottom-Up approach
	result = 0
	stopTimer = startTimer()
	for _, pattern := range wanted {
		if canConstructBU(pattern, available) {
			result++
		}
	}
	fmt.Printf("P1: Bottom-Up solution: %d, %v\n", result, stopTimer())
}

func p2(available, wanted []string) {
	// Top-Down approach
	memo := map[string]int{}
	result := 0
	stopTimer := startTimer()
	for _, pattern := range wanted {
		result += countConstructTD(pattern, available, memo)
	}
	fmt.Printf("P2: Top-Down solution:  %d, %v\n", result, stopTimer())
	// Bottom-Up approach
	result = 0
	stopTimer = startTimer()
	for _, pattern := range wanted {
		result += countConstructBU(pattern, available)
	}
	fmt.Printf("P2: Bottom-Up solution: %d, %v\n", result, stopTimer())
}

// Bottom-Up for canConstruct
func canConstructBU(target string, wordBank []string) bool {
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

// Bottom-Up for countConstruct
func countConstructBU(target string, wordBank []string) int {
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

// Top-Down approach to canConstruct
func canConstructTD(target string, wordBank []string, memo map[string]bool) bool {
	if len(target) == 0 {
		return true
	}
	if memo[target] {
		return true
	}
	for _, word := range wordBank {
		if strings.HasPrefix(target, word) {
			if canConstructTD(target[len(word):], wordBank, memo) {
				memo[target] = true
				return true
			}
		}
	}
	return false
}

// Top-Down approach to countConstruct
func countConstructTD(target string, wordBank []string, memo map[string]int) int {
	if len(target) == 0 {
		return 1
	}
	if val, ok := memo[target]; ok {
		return val
	}
	val := 0
	for _, word := range wordBank {
		if strings.HasPrefix(target, word) {
			val += countConstructTD(target[len(word):], wordBank, memo)
		}
	}
	memo[target] = val
	return val
}

// Starts a timer. And return a function that returns the elapsed time.
func startTimer() func() time.Duration {
	t0 := time.Now()
	return func() time.Duration {
		return time.Duration(time.Since(t0).Milliseconds())
	}
}
