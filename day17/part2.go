package main

import (
	"fmt"
	"log"
	"math"

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

	m.next = DFSModelFunc(m, program)
	m.hasNext = true
	// a := findDFS(program)
	// a,_ seq := findRec(program, 0)
	m.reg[A] = 0
	tea.NewProgram(m).Run()
	return m.stack.pop().answer

}

type Item struct {
	answer int
	idx    int
}

type Stack []Item

func (s *Stack) push(i Item) {
	*s = append(*s, i)
}
func (s *Stack) pop() Item {
	old := *s
	i := old[len(old)-1]
	*s = old[:len(old)-1]
	return i
}

func DFSModelFunc(m *model, program []int) func() (int, bool) {
	m.stack = Stack{}
	m.stack.push(Item{answer: 0, idx: len(program) - 1})
	m.seq = program
	answer := math.MaxInt

	return func() (int, bool) {
		log.Printf("Enter DFS")

		if len(m.stack) == 0 {
			log.Printf("Should only happen when called after finished?!")
			return -1, false
		}
		curr := m.stack.pop()

		// Check "done". (correct node, etc)
		if curr.idx < 0 {
			// HACK: This check shouldn't be needed, since we search from low to high
			if curr.answer < answer {
				answer = curr.answer
			}
			m.reg[A] = answer
			return answer, true
		}

		// Work on current node.
		for t := range 8 {
			// Init A register
			m.reg[A] = (curr.answer << 3) + t // old answer (shifted to correct pos) + new
			if m.reg[A] == 0 {
				continue // NOTE: Needed to skip if 0 in the start, becasue 0 << 3 | 0 -> 0
			}
			m.dirty[A] = true

			m.codeFns[m.pc]() // b := a % 8
			m.codeFns[m.pc]() // b = b ^ 3
			m.codeFns[m.pc]() // c := a >> b
			m.codeFns[m.pc]() // b = b ^ c
			m.codeFns[m.pc]() // b = b ^ 3

			if m.reg[B]%8 == program[curr.idx] {
				m.stack.push(Item{answer: m.reg[A], idx: curr.idx - 1})
			}
			m.pc = 0
		}

		return -1, true
	}
}

func DFSFunc(program []int) func() (int, bool) {
	stack := Stack{}
	stack.push(Item{answer: 0, idx: len(program) - 1})
	answer := math.MaxInt

	return func() (int, bool) {
		if len(stack) == 0 {
			return answer, false
		}
		// for len(stack) > 0 {
		curr := stack.pop()

		// Check "done". (correct node, etc)
		if curr.idx < 0 {
			println("New answer:", curr.answer)
			// HACK: This check shouldn't be needed, since we search from low to high
			if curr.answer < answer {
				answer = curr.answer
			}
			return answer, false
		}

		// Work on current node.
		for t := range 8 {
			// Init A register
			a := (curr.answer << 3) + t // old answer (shifted to correct pos) + new
			if a == 0 {
				continue // NOTE: Needed to skip if 0 in the start, becasue 0 << 3 | 0 -> 0
			}

			//---Run rogram---
			b := a % 8
			b = b ^ 3
			c := a >> b
			b = b ^ c
			b = b ^ 3
			//----------------

			// Find edges to new nodes (eg. neighbours)
			if b%8 == program[curr.idx] {
				stack.push(Item{answer: a, idx: curr.idx - 1})
			}
		}
		return answer, true
	}
	// return answer
}

func findDFS(program []int) int {
	stack := Stack{}
	stack.push(Item{answer: 0, idx: len(program) - 1})
	answer := math.MaxInt

	for len(stack) > 0 {
		curr := stack.pop()

		// Check "done". (correct node, etc)
		if curr.idx < 0 {
			println("New answer:", curr.answer)
			// HACK: This check shouldn't be needed, since we search from low to high
			if curr.answer < answer {
				answer = curr.answer
			}
			break
		}

		// Work on current node.
		for t := range 8 {
			// Init A register
			a := (curr.answer << 3) + t // old answer (shifted to correct pos) + new
			if a == 0 {
				continue // NOTE: Needed to skip if 0 in the start, becasue 0 << 3 | 0 -> 0
			}

			//---Run rogram---
			b := a % 8
			b = b ^ 3
			c := a >> b
			b = b ^ c
			b = b ^ 3
			//----------------

			// Find edges to new nodes (eg. neighbours)
			if b%8 == program[curr.idx] {
				stack.push(Item{answer: a, idx: curr.idx - 1})
			}
		}
	}
	return answer
}

func findRec(program []int, ans int) (int, bool) {
	if len(program) == 0 {
		return ans, true
	}
	k := len(program) - 1
	for t := range 8 {
		// Init A register
		a := (ans << 3) + t // old answer (shifted to correct pos) + new
		if a == 0 {
			continue // NOTE: Needed to skip if 0 in the start, becasue 0 << 3 | 0 -> 0
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
			sub, ok := findRec(program[:k], a)
			if ok {
				// If we found a solution, return it
				return sub, true
			}
		}
	} // jp 0, has the considion as the functions base condition
	return 0, false
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
