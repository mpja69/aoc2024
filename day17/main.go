package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type Register struct {
	A, B, C int
}

func (r Register) String() [][]string {
	return [][]string{
		{"A", fmt.Sprint(r.A)},
		{"B", fmt.Sprint(r.B)},
		{"C", fmt.Sprint(r.C)},
	}
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

	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllString(string(data), -1)

	register := Register{}
	register.A, _ = strconv.Atoi(matches[0])
	register.B, _ = strconv.Atoi(matches[1])
	register.C, _ = strconv.Atoi(matches[2])

	codeFns := []func(){}
	src := []string{}
	stringFns := []func(){}
	m := &model{reg: register}
	length := 0
	for i := 3; i < len(matches); i += 2 {
		op, _ := strconv.Atoi(matches[i])
		val, _ := strconv.Atoi(matches[i+1])

		// Store the functions to execute and show operations, and show the source code listing
		codeFns = append(codeFns, getCodeFunc(op, val, m))
		stringFns = append(stringFns, getStringFunc(op, val, m))
		src = append(src, getSource(op, val))

		// Store the input code sequence
		writeAnyToString(&m.inputCode, op)

		length++
	}
	m.reg = register
	m.codeFns = codeFns
	m.src = src
	m.stringFns = stringFns
	m.length = length

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
		"; %s mod 8  -> B",
		"; jp %s, if A!=0",
		"; B xor C -> B, (skip %s)",
		"; %s mod 8 -> OUT",
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
	m.reg.A = m.reg.A >> val
	m.dirty['A'] = true
	m.pc++
}
func advString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.opString = fmt.Sprintf("%d >> %d  -> A", m.reg.A, val)
}

// 1
func bxl(operand int, m *model) {
	m.reg.B ^= operand
	m.dirty['B'] = true
	m.pc++
}
func bxlString(operand int, m *model) {
	m.opString = fmt.Sprintf("%d xor %d  -> B", m.reg.B, operand)
}

// 2
func bst(operand int, m *model) {
	val := combo(operand, m.reg)
	m.reg.B = val % 8
	m.dirty['B'] = true
	m.pc++
}
func bstString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.opString = fmt.Sprintf("%d mod 8  -> B", val)
}

// 3
func jnz(operand int, m *model) {
	if m.reg.A != 0 {
		m.pc = operand
	} else {
		m.pc++
	}
}
func jnzString(operand int, m *model) {
	if m.reg.A != 0 {
		m.opString = fmt.Sprintf("%d  -> PC", operand)
	} else {
		m.opString = fmt.Sprintf("%d  -> PC", m.pc+1)
	}
}

// 4
func bxc(operand int, m *model) {
	m.reg.B ^= m.reg.C
	m.dirty['B'] = true
	m.pc++
}

func bxcString(operand int, m *model) {
	m.opString = fmt.Sprintf("%d xor %d -> B", m.reg.B, m.reg.C)
}

// 5
func out(operand int, m *model) {
	val := combo(operand, m.reg)
	val %= 8
	writeAnyToString(&m.output, val)
	m.pc++
}
func outString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.opString = fmt.Sprintf("%d mod 8 -> out", val)
}

// 6
func bdv(operand int, m *model) {
	val := combo(operand, m.reg)
	m.reg.B = m.reg.A >> val
	m.dirty['B'] = true
	m.pc++
}
func bdvString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.opString = fmt.Sprintf("%d >> %d -> B", m.reg.A, val)
}

func cdv(operand int, m *model) {
	val := combo(operand, m.reg)
	m.reg.C = m.reg.A >> val
	m.dirty['C'] = true
	m.pc++
}
func cdvString(operand int, m *model) {
	val := combo(operand, m.reg)
	m.opString = fmt.Sprintf("%d >> %d -> C", m.reg.A, val)
}

// ----------- Helper functions ------------
func combo(operand int, reg Register) int {
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

func writeAnyToString(dst *string, src int) {
	if len(*dst) > 0 {
		*dst += ","
	}
	*dst += fmt.Sprintf("%v", src)

}
