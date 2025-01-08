package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO:	Recursive func? Iterative
//
//	Loop until all outputs are correct (1-16):
//		Loop over 8 values -> to get the new-output match next-input:
//			Set reg A
//			Run machine 1 whole sequence, (7 instructions, skip jnz)
func p2(m *model, program []int) int {

	fmt.Printf("Program:  %v\n", program)

	a, ok, seq := findA(program, 0, []int{})
	if ok {
		fmt.Printf("Sequence: %v, %o\n", seq, a)
		m.reg[A] = convertBase8SequnceToBase10(seq)
		tea.NewProgram(m).Run()
		return a
	} else {
		return -1
	}

	return -1
}

// a = ans << 3 | t		// a = 8 + [ ]
// b = a % 8		// ersätts med b = t
// b = b ^ 2
// c = a >> b
// b = b ^ c
// b = b ^ 3
// out(b % 8)
// a = a >> 3
// if a != 0: jp 0

// bst A        ; B = A % 8, (last digit to B)
// bxl 3        ; B = B ^ 3, (this one have just another exp...doesn't matter)
// cdv B        ; C = A >> B
// bxc 0        ; B = B ^ C
// bxl 3        ; B = B ^ 3,  (these 2 switced place...doesn't matter)
// adv 3        ; A = A >> 3, (these 2 switced place...doesn't matter)
// out B        ; OUT = B % 8, (last digit to out)
// jnz 0        ; jp 0, if A!=0

// program =  [2 4 1 3 7 5 4 0 1 3 0 3 5 5 3 0]
func findA(program []int, ans int, seq []int) (int, bool, []int) {
	if len(program) == 0 {
		return ans, true, seq // NOTE: Needed to shift up one last time, (probably due to the the 0 in the start)
	}
	k := len(program) - 1
	for t := range 8 {
		// Init A register
		a := (ans << 3) + t // old answer (shifted to correct pos) + new
		if a == 0 {
			continue // Becasue 0 == 0<<3
		}
		//---Run rogram---
		b := a % 8
		b = b ^ 3
		c := a >> b
		b = b ^ c
		b = b ^ 3
		//----------------

		//a = a >> 3 // skip this ... probably becasue we want to accumulate a

		// out b % 8, becomes this comparison
		if b%8 == program[k] {
			// If we found a solution for program[k], call recursively, with rest of program[:k-1] and a as the answer
			sub, ok, candidate := findA(program[:k], a, append(seq, t))
			if ok {
				// If we found a solution, return it
				return sub, true, candidate
			}
		}
	} // jp 0, has the considion as the functions base condition
	return 0, false, nil // The 0 (zero) will be disregared from the caller
}

func rotateRight(sequence []int) []int {
	j := len(sequence) - 1
	l, r := sequence[:j], sequence[j:]
	return append(r, l...)
}
func rotateLeft(sequence []int) []int {
	l, r := sequence[:1], sequence[1:]
	return append(r, l...)
}

// func convertRevBase8ToBase10(sequence []int) int {
// 	res := 0
// 	for i, s := range sequence {
// 		ok := s << (3 * i)
// 		res += ok
// 	}
// 	return res
// }

func convertBase8SequnceToBase10(sequence []int) int {
	res := 0
	for _, s := range sequence {
		res = res*8 + s
	}
	return res
}
func (m *model) restartWith(seq []int) {
	m.reset()
	m.reg[A] = convertBase8SequnceToBase10(seq)
	m.testSequence = seq
}
