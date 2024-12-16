package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Equation struct {
	answer   int
	operands []int
}

func main() {
	data, err := os.ReadFile("s.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := strings.Trim(string(data), " \n")
	lines := strings.Split(input, "\n")
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
	fmt.Println("Part 1 (3749, 4364915411363):", p1(equations))
}

func p1(equations []Equation) int {
	res := 0
	for _, equation := range equations {
		if isValidEquation(equation) {
			res += equation.answer
		}
	}
	return res
}

func isValidEquation(eq Equation) bool {
	// answers := linearCalcAllCombos(eq.operands)
	answers := recursiveCalcAllCombos(eq.operands)
	isValid := slices.Contains(answers, eq.answer)
	// fmt.Printf("%d:\t%v\t%v\n", eq.answer, isValid, answers)
	return isValid
}

// ------------------------ Linear, Naive approach -------------------------------
func linearCalcAllCombos(operands []int) []int {
	answers := []int{}
	nbrBits := len(operands) - 1
	nbrCombinations := int(math.Pow(2, float64(nbrBits)))
	for combo := 0; combo < nbrCombinations; combo++ {
		answer := 0
		// Calculate the first 2
		answer = calculate(operands[0], operands[1], combo&1)
		// Calculate the rest
		for bit := 1; bit < len(operands)-1; bit++ {
			mask := (int(math.Pow(2, float64(bit))))
			answer = calculate(operands[bit+1], answer, (combo&mask)>>bit)
		}
		answers = append(answers, answer)
	}
	return answers
}

func calculate(a, b int, op int) int {
	if op == 0 {
		return a + b
	} else {
		return a * b
	}
}

// ------------------------ Recursive approach---------------------------------
func recursiveCalcAllCombos(operands []int) []int {
	answers := []int{}
	nbrBits := len(operands) - 1
	nbrCombinations := int(math.Pow(2, float64(nbrBits)))
	for combo := 0; combo < nbrCombinations; combo++ {
		binaryStringCombo := fmt.Sprintf("%0*b", nbrBits, combo)
		answer := recursiveCalculateForCombination(operands, binaryStringCombo)
		answers = append(answers, answer)
	}
	return answers
}

func recursiveCalculateForCombination(operands []int, operations string) int {
	last := len(operands) - 1
	if last == 0 {
		return operands[0]
	}
	if operations[0] == '0' {
		return recursiveCalculateForCombination(operands[:last], operations[1:]) + operands[last]
	} else {
		return recursiveCalculateForCombination(operands[:last], operations[1:]) * operands[last]
	}
}
