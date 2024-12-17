package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Equation struct {
	answer   int
	operands []int
}

type operation func(int, int) int

func plus(a, b int) int { return a + b }
func mult(a, b int) int { return a * b }
func concat(a, b int) int {
	for t := b; t > 0; t /= 10 {
		a *= 10
	}
	return a + b
}

func main() {
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := strings.Trim(string(data), " \n")
	lines := strings.Split(input, "\n")
	equations := parseEquations(lines)
	fmt.Println("Part 1 (3749, 4364915411363):", countValidEquations(equations, []operation{plus, mult}))
	fmt.Println("Part 2 (11387, 38322057216320):", countValidEquations(equations, []operation{plus, mult, concat}))
}

func parseEquations(lines []string) []Equation {
	equations := make([]Equation, 0, len(lines))
	for _, line := range lines {
		equation := Equation{}
		parts := strings.Split(line, ":")
		equation.answer, _ = strconv.Atoi(parts[0])
		ops := strings.Trim(parts[1], " ")
		for _, op := range strings.Split(ops, " ") {
			operand, _ := strconv.Atoi(op)
			equation.operands = append(equation.operands, operand)
		}
		equations = append(equations, equation)
	}
	return equations
}
func countValidEquations(equations []Equation, oprations []operation) int {
	res := 0
	for _, equation := range equations {
		if isValidEquation(equation.answer, equation.operands, oprations) {
			res += equation.answer
		}
	}
	return res
}

// --------- Fourth attempt: Recursive solution with operation-loop --------------
func isValidEquation(answer int, operands []int, operations []operation) bool {
	// Base case: Just 1 value
	if len(operands) == 1 {
		return answer == operands[0]
	}

	// Make sure to make a deep-copy
	cpy := make([]int, len(operands))
	copy(cpy, operands)

	// Loop over the operations
	for _, operation := range operations {
		cpy[1] = operation(operands[0], operands[1])
		if isValidEquation(answer, cpy[1:], operations) {
			return true
		}
	}
	return false
}

// --------------Third attempt: Recursive solution ---------------------
// func isValidEquation(answer int, op []int) bool {
// 	// Base case: Just 1 value
// 	if len(op) == 1 {
// 		return answer == op[0]
// 	}
//
// 	// Make sure to make a deep-copy
// 	cpy := make([]int, len(op))
// 	copy(cpy, op)
//
// 	// Case 1: Add and store in "first" place
// 	cpy[1] = op[0] + op[1]
// 	if isValidEquation(answer, cpy[1:]) {
// 		return true
// 	}
//
// 	// Case 2: Multiply and store in "first" place
// 	cpy[1] = op[0] * op[1]
// 	return isValidEquation(answer, cpy[1:])
// }
//

// ---------------------- OLD STUFF BELOW ----------------------------

// func isValidEquation(eq Equation) bool {
// 	// answers := linearCalcAllCombos(eq.operands)
// 	answers := recursiveCalcAllCombos(eq.operands)
// 	isValid := slices.Contains(answers, eq.answer)
// 	return isValid
// }

// --------------- First attempt: Linear, Naive approach -------------------------
// func linearCalcAllCombos(operands []int) []int {
// 	answers := []int{}
// 	nbrBits := len(operands) - 1
// 	nbrCombinations := int(math.Pow(2, float64(nbrBits)))
// 	for combo := 0; combo < nbrCombinations; combo++ {
// 		answer := 0
// 		// Calculate the first 2
// 		answer = calculate(operands[0], operands[1], combo&1)
// 		// Calculate the rest
// 		for bit := 1; bit < len(operands)-1; bit++ {
// 			mask := (int(math.Pow(2, float64(bit))))
// 			answer = calculate(operands[bit+1], answer, (combo&mask)>>bit)
// 		}
// 		answers = append(answers, answer)
// 	}
// 	return answers
// }
// func calculate(a, b int, op int) int {
// 	if op == 0 {
// 		return a + b
// 	} else {
// 		return a * b
// 	}
// }

// ------------------ Second attempt: Kind of recursive approach -------------------------
// func recursiveCalcAllCombos(operands []int) []int {
// 	answers := []int{}
// 	nbrBits := len(operands) - 1
// 	nbrCombinations := int(math.Pow(2, float64(nbrBits)))
// 	for combo := 0; combo < nbrCombinations; combo++ {
// 		binaryStringCombo := fmt.Sprintf("%0*b", nbrBits, combo)
// 		answer := recursiveCalculateForCombination(operands, binaryStringCombo)
// 		answers = append(answers, answer)
// 	}
// 	return answers
// }
// func recursiveCalculateForCombination(operands []int, operations string) int {
// 	last := len(operands) - 1
// 	if last == 0 {
// 		return operands[0]
// 	}
// 	if operations[0] == '0' {
// 		return recursiveCalculateForCombination(operands[:last], operations[1:]) + operands[last]
// 	} else {
// 		return recursiveCalculateForCombination(operands[:last], operations[1:]) * operands[last]
// 	}
// }
