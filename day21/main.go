package main

import (
	"bytes"
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/mpja69/aoc2024/day21/keypad"
)

func main() {
	data, _ := os.ReadFile("data.txt")
	data = bytes.TrimSpace(data)
	lines := strings.Split(string(data), "\n")

	p1(lines)
	p2(lines)
}

type Pair struct {
	a, b byte
}

var dirSeqs map[Pair][]keypad.Sequence
var dirLengths map[Pair]int

func p2(lines []string) {
	res := 0

	// Create the memoization cache map, (for dir-layout), and fill level 1 with values
	dp := keypad.NewKeypad(keypad.DirectionLayout)
	cache = map[seqItem]int{}
	dirSeqs = map[Pair][]keypad.Sequence{}
	dirLengths = map[Pair]int{}
	for k, v := range dp.Move2seq {
		a := byte(k.From)
		b := byte(k.To)
		dirSeqs[Pair{a, b}] = v
		dirLengths[Pair{a, b}] = len(v[0])
		// cache[item{a, b, 1}] = len(v[0])
	}

	// Go over the numeric lines...
	np := keypad.NewKeypad(keypad.NumberLayout)
	for _, line := range lines {
		dirSequences := np.GetPossibleSequences(line)
		length := math.MaxInt
		// ...and call the recursive function for each seuence
		for _, seq := range dirSequences {
			// ...and sum up the shortest one
			length = min(length, solve(seq, 25))
		}
		res += length * value(line)
	}
	fmt.Println("P2 (264518225304496):", res)
}

func pairsFromSeq(seq string) func(func(Pair) bool) {
	a := byte('A')
	return func(yield func(Pair) bool) {
		for _, b := range []byte(seq) {
			if !yield(Pair{a, b}) {
				return
			}
			a = b
		}
	}
}

type seqItem struct {
	seq string
	l   int
}

var cache map[seqItem]int

// Depth First approach - for a sequence (string) at a time
func solve(seq string, level int) int {
	if level == 1 {
		// Base case: On level 1, we use the pre-computed map, summed over the sequence
		length := 0
		for p := range pairsFromSeq(seq) {
			length += dirLengths[p]
		}
		return length
	}

	// Memoization
	if v, ok := cache[seqItem{seq, level}]; ok {
		return v
	}

	length := 0
	for p := range pairsFromSeq(seq) {
		minLength := math.MaxInt
		for _, subSeq := range dirSeqs[p] {
			minLength = min(minLength, solve(string(subSeq), level-1))
		}
		length += minLength
	}

	cache[seqItem{seq, level}] = length
	return length
}

// func runByte(seqs []string) int {
// 	kp := keypad.NewKeypad(keypad.DirectionLayout)
// 	optimalLength := 1_000_000_000_000_000_000
// 	for _, seq := range seqs {
// 		length := 0
// 		for _, b := range seq {
// 			length += solve(p, 25, kp)
// 			a = b
// 		}
// 		optimalLength = min(optimalLength, length)
// 	}
// 	return optimalLength
// }

// -------------------- Top-Down solution for pairwise
// type item struct {
// 	a, b byte
// 	l    int
// }
//
//	Depth First approach - working on moves between pairs
// func solve(p P, level int, kp *keypad.Keypad) int {
// 	// ==== Not needed since level is cached from the beginning ====
// 	// if level == 1 {
// 	// 	return len(kp.Move(a, b)[0])
// 	// }
//
// 	v, ok := cache[item{a, b, level}]
// 	if ok {
// 		return v
// 	}
// 	optimalLength := 1_000_000_000_000_000_000
// 	for _, seq := range kp.Move(a, b) {
// 		length := 0
//		for p := range pairsFromSeq(seq) {
// 			length += solve(p, level-1, kp)
// 		}
// 		optimalLength = min(optimalLength, length)
// 	}
// 	cache[item{a, b, level}] = optimalLength
// 	return optimalLength
// }

func p1(lines []string) {
	sum := 0
	for _, line := range lines {
		kp := keypad.NewKeypad(keypad.NumberLayout)
		lines := kp.GetPossibleSequences(line)
		sum += runLine(lines) * value(line)
	}
	fmt.Println("P1 (219254):", sum)
}

// Breadth First approach - doeasn't work if the stack depth is too deep
func runLine(lines []string) int {
	kp := keypad.NewKeypad(keypad.DirectionLayout)
	dirLines := []string{}
	for _, line := range lines {
		dirLines = append(dirLines, kp.GetPossibleSequences(line)...)
	}
	// Only pick the shortes ones
	lenCmp := func(a, b string) int { return cmp.Compare(len(a), len(b)) }
	slices.SortFunc(dirLines, lenCmp)

	// Finally another dir pad
	dirLines2 := []string{}
	for line := range onlyShortest(dirLines) {
		dirLines2 = append(dirLines2, kp.GetPossibleSequences(line)...)
	}

	// Find the length of the shortes ones
	return len(slices.MinFunc(dirLines2, lenCmp))
}

func onlyShortest(s []string) func(func(string) bool) {
	return func(yield func(string) bool) {
		for _, line := range s {
			if len(line) > len(s[0]) {
				break
			}
			yield(line)
		}
	}
}
func value(code string) int {
	c, _ := strconv.Atoi(string(code[:3]))
	return c
}
