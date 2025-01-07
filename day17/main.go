package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
	data, err := os.ReadFile("s2.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}

	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllString(string(data), -1)
	regs, prog := matches[:3], matches[3:]

	//register := Register{}
	register := make([]int, 3)

	register[A], _ = strconv.Atoi(regs[0])
	register[B], _ = strconv.Atoi(regs[1])
	register[C], _ = strconv.Atoi(regs[2])

	codeFns := []func(){}
	src := []string{}
	stringFns := []func(){}
	m := &model{}

	length := 0
	sequence := []int{}
	for i := 0; i < len(prog); i += 2 {
		op, _ := strconv.Atoi(prog[i])
		val, _ := strconv.Atoi(prog[i+1])

		// Store the functions to execute and show operations, and show the source code listing
		codeFns = append(codeFns, getCodeFunc(op, val, m))
		stringFns = append(stringFns, getStringFunc(op, val, m))
		src = append(src, getSource(op, val))

		sequence = append(sequence, op, val)
		length++
	}
	m.reg = register
	m.codeFns = codeFns
	m.src = src
	m.stringFns = stringFns
	m.length = length
	m.inputProg = strings.Join(prog, ",")
	m.dirty = [3]bool{}

	// fmt.Printf("P1: (1,5,7,4,1,6,0,3,0) %s\n", p1(m))
	fmt.Printf("P2: Base10 Reg.A: %d\n", p2(m, sequence))
}

func p1(m *model) string {
	// Run Bubble Tea
	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}

	return strings.Join(m.output, ",")
}

// ----------- Wrapper for the operations --------------

// This func that returns a func that calls a func
// Usually a wrapper like this...Is a func that returns a func that is declared within
func getCodeFunc(operation, operand int, m *model) func() {
	functions := []func(int, *model){adv, bxl, bst, jnz, bxc, out, bdv, cdv}
	return func() { functions[operation](operand, m) }
}

func getStringFunc(operation, operand int, m *model) func() {
	functions := []func(int, *model){advString, bxlString, bstString, jnzString, bxcString, outString, bdvString, cdvString}
	return func() { functions[operation](operand, m) }
}

func getSource(operation, operand int) string {
	mnemonics := []string{"adv", "bxl", "bst", "jnz", "bxc", "out", "bdv", "cdv"}
	comments := []string{
		"; A >> %s -> A",
		"; B xor %s -> B",
		"; %s & 7 -> B, (last digit to B)",
		"; jp %s, if A!=0",
		"; B xor C -> B, (not using the operand: %s)",
		"; %s & 7 -> OUT, (last digit to out)",
		"; A >> %s -> B",
		"; A >> %s -> C"}
	litOrReg := []string{"0", "1", "2", "3", "A", "B", "C", "N/A"}
	format := comments[operation]
	comment := fmt.Sprintf(format, litOrReg[operand])
	return fmt.Sprintf("%d %d\t\t%s %s\t\t%s", operation, operand, mnemonics[operation], litOrReg[operand], comment)
}

// ------------ The functions for each operation --------------------
// 0
func adv(operand int, m *model) {
	val := combo(operand, m.reg)
	m.reg[A] = m.reg[A] >> val
	m.dirty[A] = true
	m.pc++
}
func advString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.currOpString = fmt.Sprintf("%o >> %o  -> A", m.reg[A], val)
}

// 1
func bxl(operand int, m *model) {
	m.reg[B] ^= operand
	m.dirty[B] = true
	m.pc++
}
func bxlString(operand int, m *model) {
	m.currOpString = fmt.Sprintf("%o xor %o  -> B", m.reg[B], operand)
}

// 2
func bst(operand int, m *model) {
	val := combo(operand, m.reg)
	// m.reg[B] = val % 8
	m.reg[B] = val & 7
	m.dirty[B] = true
	m.pc++
}
func bstString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.currOpString = fmt.Sprintf("%o & 7  -> B", val)
}

// 3
func jnz(operand int, m *model) {
	if m.reg[A] != 0 {
		m.pc = operand
	} else {
		m.pc++
	}
}
func jnzString(operand int, m *model) {
	if m.reg[A] != 0 {
		m.currOpString = fmt.Sprintf("%o  -> PC", operand)
	} else {
		m.currOpString = fmt.Sprintf("%o  -> PC", m.pc+1)
	}
}

// 4
func bxc(operand int, m *model) {
	m.reg[B] ^= m.reg[C]
	m.dirty[B] = true
	m.pc++
}

func bxcString(operand int, m *model) {
	m.currOpString = fmt.Sprintf("%o xor %o  -> B", m.reg[B], m.reg[C])
}

// 5
func out(operand int, m *model) {
	val := combo(operand, m.reg)
	val &= 7
	m.output = append(m.output, string(val+'0'))
	m.newOutput = val
	m.pc++
}
func outString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.currOpString = fmt.Sprintf("%o & 7  -> out", val)
}

// 6
func bdv(operand int, m *model) {
	val := combo(operand, m.reg)
	m.reg[B] = m.reg[A] >> val
	m.dirty[B] = true
	m.pc++
}
func bdvString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.currOpString = fmt.Sprintf("%o >> %o  -> B", m.reg[A], val)
}

func cdv(operand int, m *model) {
	val := combo(operand, m.reg)
	m.reg[C] = m.reg[A] >> val
	m.dirty[C] = true
	m.pc++
}
func cdvString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.currOpString = fmt.Sprintf("%o >> %o  -> C", m.reg[A], val)
}

// ----------- Helper functions ------------
func combo(operand int, reg []int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return reg[A]
	case 5:
		return reg[B]
	case 6:
		return reg[C]
	}
	log.Fatalf("Err: Invalid operand! op=%d", operand)
	return 0
}

// NOTE: Not used. Using strings.Join instead
func writeAnyToString(dst *string, src int) {
	if len(*dst) > 0 {
		*dst += ","
	}
	*dst += fmt.Sprintf("%v", src)

}
