package main

import (
	"fmt"
	"log"
	"os"
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
	data, err := os.ReadFile("sample.txt")
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

	comboStrings := strings.Split(strings.Split(parts[1], ":")[1], ",")
	// combo := []Combo{}
	fns := []func(){} // HACK:--------------------- En lista med Closures (funktioner med "minne")
	src := []string{}

	for i := 0; i < len(comboStrings); i += 2 {
		op, _ := strconv.Atoi(comboStrings[i])
		val, _ := strconv.Atoi(comboStrings[i+1])
		// instruction := Combo{mnemonic: mnemonics[op], operand: val, fn: methods[op]}
		// combo = append(combo, instruction)
		fn := operationFunc(op, val) // HACK:--------------------- Som en "factory" som ger "rätt" closure och "sparar" input
		fns = append(fns, fn)        // HACK:---------------------
		src = append(src, srcCode(op, val))
	}

	fmt.Printf("P1: () %d\n", p1(register, fns, src))
}

var mnemonics = []string{
	"adv",
	"bxl",
	"bst",
	"jnz",
	"bxc",
	"out",
	"bdv",
	"cdv",
}

var methods = []func(Combo){
	Combo.adv,
	Combo.bxl,
	Combo.bst,
	Combo.jnz,
	Combo.bxc,
	Combo.out,
	Combo.bdv,
	Combo.cdv,
}

func operationFunc(operation, operand int) func() {

	log.Printf("ops %d, oper %d", operation, operand)
	switch operation {
	case 0:
		return func() { adv(operand) }
	case 1:
		return func() { bxl(operand) }
	case 2:
		return func() { bst(operand) }
	case 3:
		return func() { jnz(operand) }
	case 4:
		return func() { bxc(operand) }
	case 5:
		return func() { out(operand) }
	case 6:
		return func() { bdv(operand) }
	case 7:
		return func() { cdv(operand) }
	}
	return nil
}

func srcCode(operation, operand int) string {
	return fmt.Sprintf("%s %d", mnemonics[operation], operand)
}

// 0
func adv(operand int) {
	log.Printf("adv (%v): operand %d\n", adv, operand)
	val := getVal(operand)
	m.register.A /= (1 << val)
	m.PC++
}

// 1
func bxl(operand int) {
	log.Printf("bxl: operand %d\n", operand)
	m.register.B ^= operand
	m.PC++
}

// 2 FIX: ska vara out-5
func bst(operand int) {
	log.Printf("bst: operand %d\n", operand)
	val := getVal(operand)
	m.register.B = val % 8
	m.PC++
}

// 3
func jnz(operand int) {
	log.Printf("jnz: operand %d\n", operand)
	if m.register.A != 0 {
		m.PC = operand
	} else {
		m.PC++
	}
}

// 4 FIX: ska vara 3- jnz
func bxc(operand int) {
	log.Printf("bxc: operand %d\n", operand)
	m.register.B ^= m.register.C
	m.PC++
}

// 5
func out(operand int) {
	log.Printf("out: operand %d\n", operand)
	val := getVal(operand) % 8
	if len(m.output) > 0 {
		m.output += ","
	}
	m.output += fmt.Sprintf("%d", val)
	m.PC++
}

// 6
func bdv(operand int) {
	log.Printf("bdv: operand %d\n", operand)
	val := getVal(operand)
	m.register.B = m.register.A / (1 << val)
	m.PC++
}

// 7
func cdv(operand int) {
	log.Printf("cdv: operand %d\n", operand)
	val := getVal(operand)
	m.register.C = m.register.A / (1 << val)
	m.PC++
}

// -------------------- Cmds --------------------
// Men detta funkar ju inte...för jag måste ju ge operanden dynamiskt.
// 0
func advFunc(operand int) tea.Cmd {
	return func() tea.Msg {
		adv(operand)
		return m.PC
	}
}

func getVal(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return m.register.A
	case 5:
		return m.register.B
	case 6:
		return m.register.C
	}
	log.Fatalf("Err: Invalid operand! op=%d", operand)
	return 0
}

type Register struct {
	A, B, C int
}

type Combo struct {
	mnemonic string
	operand  int
	fn       func(c Combo)
}

func (c Combo) String() string {
	return fmt.Sprintf("%s %d", c.mnemonic, c.operand)
}

func (c Combo) adv() { //adv(c.operand) }
	val := getVal(c.operand)
	m.register.A /= (1 << val)
	m.PC++
}

func (c Combo) bxl() { //bxl(c.operand) }
	m.register.B ^= c.operand
	m.PC++
}

func (c Combo) bst() { //bst(c.operand) }
	val := getVal(c.operand)
	m.register.B = val % 8
	m.PC++
}

func (c Combo) jnz() { //jnz(c.operand) }
	if m.register.A != 0 {
		m.PC = c.operand
	} else {
		m.PC++
	}
}

func (c Combo) bxc() { //bxc(c.operand) }
	m.register.B ^= m.register.C
	m.PC++
}

func (c Combo) out() { //out(c.operand) }
	val := getVal(c.operand) % 8
	if len(m.output) > 0 {
		m.output += ","
	}
	m.output += fmt.Sprintf("%d", val)
	m.PC++
}

func (c Combo) bdv() { //bdv(c.operand) }
	val := getVal(c.operand)
	m.register.B = m.register.A / (1 << val)
	m.PC++
}

func (c Combo) cdv() { // cdv(c.operand) }
	val := getVal(c.operand)
	m.register.C = m.register.A / (1 << val)
	m.PC++
}

var m *model

func p1(regs Register, ops []func(), src []string) int {

	m = &model{register: regs, ops: ops, src: src}

	// Run Bubble Tea
	if _, err := tea.NewProgram(m).Run(); err != nil {
		os.Exit(1)
	}

	return 0
}
