package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := strings.Trim(string(data), " \n")
	parts := strings.Split(input, " ")

	numbers := make([]int, len(parts))
	for i, nbr := range parts {
		num, _ := strconv.Atoi(nbr)
		numbers[i] = int(num)
	}

	fmt.Printf("Part 1: (55321, 200446): %d\n", part1(numbers, 25))
	fmt.Printf("Part 2: (65601038650482, 238317474993392): %d\n", part2(numbers, 75))
}

// ----------------------------- Part 1 - Linear solution using slices ----------------------------------------
func part1(values []int, steps int) int {
	for range steps {
		values = applyRules(values)
	}
	return len(values)
}

func applyRules(values []int) []int {
	res := []int{}
	for _, val := range values {
		if val == 0 {
			res = append(res, 1)
		} else if length := lenInt(val); length%2 == 0 {
			a, b := divideInt(val)
			res = append(res, a, b)
		} else {
			res = append(res, val*2024)
		}
	}
	return res
}

func lenInt(n int) int {
	x, i := int(10), 1
	for x <= n {
		x *= 10
		i++
	}
	return i
}

func divideInt(n int) (a, b int) {
	half := lenInt(n) / 2
	a = n / int(math.Pow10(half))
	b = n - a*int(math.Pow10(half))
	return
}

// ----------------------- Part 2 - Recursive solution using length AND Memoisation ---------------------
type Item struct {
	val   int
	steps int
}

var memo map[Item]int

func part2(values []int, steps int) int {
	memo = make(map[Item]int)
	res := 0
	for _, val := range values {
		res += applyRulesRec(val, steps)
	}
	return res
}

func applyRulesRec(val, steps int) int {
	item := Item{val, steps}

	if res, ok := memo[item]; ok {
		return res
	}

	if steps == 0 {
		memo[item] = 1
	} else if val == 0 {
		memo[item] = applyRulesRec(1, steps-1)
	} else if length := lenInt(val); length%2 == 0 {
		a, b := divideInt(val)
		memo[item] = applyRulesRec(a, steps-1) + applyRulesRec(b, steps-1)
	} else {
		memo[item] = applyRulesRec(2024*val, steps-1)
	}
	return memo[item]
}
