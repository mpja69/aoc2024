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
	numbers := strings.Split(input, " ")

	stones := make([]int, len(numbers))
	for i, nbr := range numbers {
		num, _ := strconv.Atoi(nbr)
		stones[i] = int(num)
	}

	fmt.Printf("Part 1: (55321, 200446): %d\n", part1(stones, 25))
	fmt.Printf("Part 2: (65601038650482, 238317474993392): %d\n", part2(stones, 75))
}

// ----------------------------- Part 1 - Linear solution using slices ----------------------------------------
func part1(stones []int, steps int) int {
	for range steps {
		stones = applyRules(stones)
	}
	return len(stones)
}

func applyRules(stones []int) []int {
	res := []int{}
	for _, val := range stones {
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
	stone int
	steps int
}

var memo map[Item]int

func part2(stones []int, steps int) int {
	memo = make(map[Item]int)
	res := 0
	for _, stone := range stones {
		res += applyRulesRec(stone, steps)
	}
	return res
}

func applyRulesRec(stone, steps int) (res int) {
	item := Item{stone, steps}
	if _, ok := memo[item]; ok {
		res = memo[item]
		return
	}

	if steps == 0 {
		res = 1
	} else if stone == 0 {
		res = applyRulesRec(1, steps-1)
	} else if length := lenInt(stone); length%2 == 0 {
		a, b := divideInt(stone)
		res = applyRulesRec(a, steps-1) + applyRulesRec(b, steps-1)
	} else {
		res = applyRulesRec(2024*stone, steps-1)
	}
	memo[item] = res
	return
}
