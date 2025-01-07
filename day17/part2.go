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
func p2(m *model, sequence []int) int {

	fmt.Printf("origianal base8 seq: %v\n", sequence)

	sequence = rotateRight(sequence)
	fmt.Printf("rotated base8 seq: %v\n", sequence)

	// sequence = []int{0, 0, 3}
	// Calulate base 10 of the reverse sequence
	base10 := convertBase8SequnceToBase10(sequence)
	fmt.Printf("Reversed seq base8: %o\n", base10)

	// m.run(sequence)
	m.reg[A] = base10
	tea.NewProgram(m).Run()
	return base10
}
func rotateRight(sequence []int) []int {
	j := len(sequence) - 1
	l, r := sequence[:j], sequence[j:]
	return append(r, l...)
}

func convertBase8SequnceToBase10(sequence []int) int {
	res := 0
	for i, s := range sequence {
		ok := s << (3 * i)
		res += ok
	}
	return res
}

// func convertBase8ToBase10(base8 int) int {
// 	base10 := 0
// 	for base8 > 0 {
// 		o := base8 & 7
// 		base10 = base10 * 10 + o
// 		base8 = base8 >> 3
// 	}
// 	return base10
// }

func (m *model) run(seq []int) int {
	octalSeq := []int{}
	for i := range len(seq) - 1 {
		for {
			fmt.Printf("Octals so far: %v...", octalSeq)
			fmt.Printf("Searching for output code: %d, at pos %d...", seq[i], i)

			octal, ok := m.loopOverOctalValues2(octalSeq, seq[i])
			if ok {
				fmt.Printf("Found octal: %o for #%d, (%d)", octal, i, seq[i])
				octalSeq = append(octalSeq, octal)
				break
			} else {
				fmt.Printf("Not matching last...add more")
				octalSeq = append(octalSeq, 0) //FIX: hur kan jag använda 2 platser för att söka efter nästa nya
			}
		}
	}
	return 0
}

// // Ska denna vara rekursiv??
// func (m *model) loopOverOctalValues(octalSeq []int, nextIput int) (int, bool) {
// 	cpy := append(octalSeq, 0)
//
// 	for octal := range 8 {
// 		cpy[len(cpy)-1] = octal
// 		fmt.Printf("%v\t", cpy)
// 		oRes, ok := m.loopOverOctalValues2(cpy, nextIput)
// 		if ok {
// 			return oRes, true
// 		}
// 	}
// 	return 0, false
// }

// Ska denna vara rekursiv??
func (m *model) loopOverOctalValues2(octalSeq []int, nextIput int) (int, bool) {
	cpy := append(octalSeq, 0)
	println()
	for octal := range 8 {
		cpy[len(cpy)-1] = octal
		fmt.Printf("%v\n", cpy)
		a := convertBase8SequnceToBase10(cpy)
		m.resetReg(a, 0, 0)
		// fmt.Printf("\tA: m.reg.A: %d, %o\n", m.reg.A, m.reg.A)
		if m.runOnce(nextIput) {
			fmt.Printf("octal: %d, nextInput: %d", octal, nextIput)
			return octal, true
		}
	}
	return 0, false
}

func (m *model) runOnce(nextInput int) bool {
	m.pc = 0                // Start from the top
	for m.pc < m.length-1 { // Run program 1 time, (skip the jnz 0)
		m.codeFns[m.pc]()
	}
	// Now check new-output against the next input (idx)
	return m.newOutput == nextInput
}

func (m *model) resetReg(a, b, c int) { m.reg = []int{a, b, c} }
func (m *model) setRegA(a int)        { m.reg[A] = a }
