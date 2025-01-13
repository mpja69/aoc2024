package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	data, _ := os.ReadFile("data.txt")
	wantedPatterns := strings.Split(strings.TrimSpace(string(data)), "\n")

	availablePatterns := strings.Split(wantedPatterns[0], ", ")
	wantedPatterns = wantedPatterns[2:]

	fmt.Printf("P1: %d\n", p1(availablePatterns, wantedPatterns))
}

func p1(available, wanted []string) int {
	res := 0
	for _, pattern := range wanted {
		if canConstruct(pattern, available) {
			res++
		}
	}
	return res
}

// Dynamic programming "canConstruct"
func canConstruct(target string, wordBank []string) bool {
	table := make([]bool, len(target)+1) // Table represent how much of the word that can be built...up to (not including) each position
	table[0] = true                      // The empty string ("") CAN be built -> true in first position

	for i := range table {
		// If I can get to this pos in the table with any word...
		if table[i] {
			// Then loop every word in the bank
			for _, w := range wordBank {
				// Check if the word match the part of the target at the current table position
				if strings.HasPrefix(target[i:], w) {
					// Mark the table
					table[i+len(w)] = true
				}
			}
		}
	}
	return table[len(target)]
}
