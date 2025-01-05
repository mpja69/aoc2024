package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Register struct {
	A, B, C int
}

func main() {
	// Setup logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	log.Printf("-------------------------------------------")

	// Read the input
	data, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Parse the the input
	input := strings.Trim(string(data), " \n")
	parts := strings.Split(input, "\n\n")
	regStrings := strings.Split(parts[0], "\n")
	register := Register{}

	A := strings.Trim(strings.Split(regStrings[0], ":")[1], " ")
	register.A, _ = strconv.Atoi(A)
	B := strings.Trim(strings.Split(regStrings[1], ":")[1], " ")
	register.B, _ = strconv.Atoi(B)
	C := strings.Trim(strings.Split(regStrings[2], ":")[1], " ")
	register.C, _ = strconv.Atoi(C)

	cleanStrings := strings.Trim(strings.Split(parts[1], ":")[1], " ")
	comboStrings := strings.Split(cleanStrings, ",")

	log.Printf("%v", comboStrings)
	code := []func(){}
	src := []string{}

	m := &model{register: register, code: code, src: src}
	for i := 0; i < len(comboStrings); i += 2 {
		op, _ := strconv.Atoi(comboStrings[i])
		val, _ := strconv.Atoi(comboStrings[i+1])

		m.code = append(m.code, getCode(op, val, m))
		m.src = append(m.src, getSource(op, val))
	}

	fmt.Printf("P1: (1,5,7,4,1,6,0,3,0) %s\n", p1(m))
}

func p1(m *model) string {
	// Run Bubble Tea
	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}

	return m.output
}

// ----------- Wrapper for the operations --------------
func getCode(operation, operand int, m *model) func() {
	functions := []func(int, *model){adv, bxl, bst, jnz, bxc, out, bdv, cdv}
	return func() { functions[operation](operand, m) }
}

func getSource(operation, operand int) string {
	mnemonics := []string{"0.adv", "1.bxl", "2.bst", "3.jnz", "4.bxc", "5.out", "6.bdv", "7.cdv"}
	return fmt.Sprintf("%s %d", mnemonics[operation], operand)
}

// ------------ The functions for each operation --------------------
// 0
func adv(operand int, m *model) {
	val := getVal(operand, m.register)
	m.register.A = dividePow(m.register.A, val)
	m.PC++
}

// 1
func bxl(operand int, m *model) {
	log.Printf("bxl: operand %d\n", operand)
	m.register.B ^= operand
	m.PC++
}

// 2
func bst(operand int, m *model) {
	log.Printf("bst: operand %d\n", operand)
	val := getVal(operand, m.register)
	m.register.B = val % 8
	m.PC++
}

// 3
func jnz(operand int, m *model) {
	log.Printf("jnz: operand %d\n", operand)
	if m.register.A != 0 {
		m.PC = operand
	} else {
		m.PC++
	}
}

// 4
func bxc(operand int, m *model) {
	log.Printf("bxc: operand %d\n", operand)
	m.register.B ^= m.register.C
	m.PC++
}

// 5
func out(operand int, m *model) {
	log.Printf("out: operand %d\n", operand)
	val := getVal(operand, m.register) % 8
	if len(m.output) > 0 {
		m.output += ","
	}
	m.output += fmt.Sprintf("%d", val)
	m.PC++
}

// 6
func bdv(operand int, m *model) {
	log.Printf("bdv: operand %d\n", operand)
	val := getVal(operand, m.register)
	m.register.B = dividePow(m.register.A, val)
	m.PC++
}

// 7
func cdv(operand int, m *model) {
	log.Printf("cdv: operand %d\n", operand)
	val := getVal(operand, m.register)
	m.register.C = dividePow(m.register.A, val)
	m.PC++
}

// ----------- Helper functions ------------
func getVal(operand int, reg Register) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return reg.A
	case 5:
		return reg.B
	case 6:
		return reg.C
	}
	log.Fatalf("Err: Invalid operand! op=%d", operand)
	return 0
}

func dividePow(numerator, logDenominator int) int {
	res := float64(numerator)
	for range logDenominator {
		res /= 2.0
	}
	return int(math.Trunc(res))
}
