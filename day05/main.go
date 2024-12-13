package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type pair [2]int
type pairMap map[pair]bool

// Utility functions for the input data
func convertToUpdateSlice(lines []string, delimiter string) [][]int {
	items := make([][]int, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, delimiter)
		item := make([]int, 0, len(parts))
		for i := range parts {
			num, _ := strconv.Atoi(parts[i])
			item = append(item, num)
		}
		items = append(items, item)
	}
	return items
}

func convertToRuleMap(lines []string, delimiter string) pairMap {
	ruleMap := make(pairMap)
	for _, line := range lines {
		rule := pair{}
		parts := strings.Split(line, delimiter)
		rule[0], _ = strconv.Atoi(parts[0])
		rule[1], _ = strconv.Atoi(parts[1])
		ruleMap[rule] = true
	}
	return ruleMap
}

func main() {
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	str := string(data)
	str = strings.Trim(str, " \n")
	parts := strings.Split(str, "\n\n")
	part1 := strings.Split(parts[0], "\n")
	rules := convertToRuleMap(part1, "|")

	part2 := strings.Split(parts[1], "\n")
	updates := convertToUpdateSlice(part2, ",")

	res := p1(updates, rules)
	fmt.Println("Part 1 (5374):", res)

	res = p2(updates, rules)
	fmt.Println("Part 2 (4260):", res)
}

// ------------------------- PART 1 ----------------------------
func p1(updates [][]int, rules pairMap) int {
	res := 0
	for _, update := range updates {
		if rules.followedBy(update) == true {
			res += getMiddle(update)
		}
	}
	return res
}

func (rules pairMap) followedBy(update []int) bool {
	for i := 0; i < len(update)-1; i++ {
		rule := pair{update[i], update[i+1]}
		if rules.missingRule(rule) {
			return false
		}
	}
	return true
}

func (rules pairMap) missingRule(p pair) bool {
	return rules[p] == false
}

func getMiddle(update []int) int {
	return update[len(update)/2]
}

// ------------------------- PART 2 ----------------------------
func p2(updates [][]int, rules pairMap) int {
	// Get the invalid updates
	invalidUpdates := [][]int{}
	for _, update := range updates {
		if rules.followedBy(update) == false {
			invalidUpdates = append(invalidUpdates, update)
		}
	}

	// Make them valid...to find correct middle
	for _, update := range invalidUpdates {
		makeUpdateValid(update, rules)
	}

	// Calculate sum
	res := 0
	for _, update := range invalidUpdates {
		res += getMiddle(update)
	}

	return res
}

// Like a "bubble sort"
// This could probably be optimized...
// So that we don't start from the beginning every time
func makeUpdateValid(update []int, rules pairMap) {
	// for !isValidUpdate(update, rules) {
	for rules.followedBy(update) == false {
		for i, j := 0, 1; j < len(update); i, j = i+1, j+1 {
			rule := pair{update[i], update[j]}
			if rules.missingRule(rule) {
				update[i], update[j] = update[j], update[i]
			}
		}
	}
}
