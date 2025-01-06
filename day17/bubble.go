package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
)

type model struct {
	reg            Register        // Registers: A, B C
	dirty          map[rune]bool   // Dirty flags for the registers
	pc             int             // Program Counter
	output         string          // Comma separated string of operations
	inputCode      string          // Comma separated string of operations
	length         int             // Length of all the slices
	opString       string          // The string representation och the next operation
	stringFns      []func()        // The string generations of the opreations
	codeFns        []func()        // The operations as functions (closures with each operand)
	src            []string        // Disassembled source code
	NextCmd        tea.Cmd         // Bubbletea Command for running Next operation
	AllCmd         tea.Cmd         // Bubbletea Command for running All operations
	StringCmd      tea.Cmd         // Bubbletea Command for rendering the next (next) operation
	codeEnumerator list.Enumerator // Bubbletea
	codeStyle      list.StyleFunc  // Bubbletea
	RegStyle       table.StyleFunc // Bubbletea
	srcList        *list.List      // Bubbletea
}
type pcMsg int
type stringMsg string

func (m *model) Init() tea.Cmd {
	m.NextCmd = func() tea.Msg {
		if m.pc < m.length {
			m.dirty = map[rune]bool{}
			m.codeFns[m.pc]()
		}
		return pcMsg(m.pc)
	}

	m.AllCmd = func() tea.Msg {
		m.dirty = map[rune]bool{}
		for m.pc < m.length {
			m.codeFns[m.pc]()
		}
		return pcMsg(m.pc)
	}

	m.StringCmd = func() tea.Msg {
		if m.pc < m.length {
			m.stringFns[m.pc]() // pc must be updated, since this looks at the next operation
		}
		return stringMsg(m.opString)
	}

	m.codeEnumerator = func(_ list.Items, i int) string {
		if i == m.pc {
			return fmt.Sprintf("->%d: ", i)
		}
		return fmt.Sprintf("  %d: ", i)
	}

	m.codeStyle = func(_ list.Items, i int) lipgloss.Style {
		if i == m.pc {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		}
		return lipgloss.NewStyle()
	}
	m.RegStyle = func(row, _ int) lipgloss.Style {
		if m.dirty[rune('A'+row)] {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		}
		return lipgloss.NewStyle()
	}

	m.srcList = list.New(m.src).Enumerator(m.codeEnumerator).ItemStyleFunc(m.codeStyle)
	return m.StringCmd
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":
			return m, m.NextCmd
		case "a":
			return m, m.AllCmd
		}
	case pcMsg:
		if int(msg) >= m.length {
			return m, tea.Quit
		}
		return m, m.StringCmd
	}

	return m, nil
}

func (m *model) View() string {
	tReg := table.New().Headers("Reg", "        ").Rows(m.reg.String()...).StyleFunc(m.RegStyle)

	tIO := table.New().Headers("I/0", "")
	tIO.Row("in", m.inputCode)
	tIO.Row("out", m.output)
	tIO.Row("", "")

	tExec := table.New().Headers("            Execution            ")
	if m.pc < m.length {
		// currOp := m.operations[m.pc]
		// currVal := m.operands[m.pc]
		// format := codeFormats[currOp]
		// numbers := getNumbers(currOp, currVal, m.reg)
		// tExec.Row(fmt.Sprintf(format, numbers...))
		tExec.Row(m.opString)
	} else {
		tExec.Row("Finished!")
	}

	s := lipgloss.JoinHorizontal(lipgloss.Top, tReg.Render(), tIO.Render())
	s = lipgloss.JoinVertical(lipgloss.Center, s, tExec.Render())

	s += "\n\n" + m.srcList.String()

	s += "\n\n\nControls: [N]ext  [A]ll  [Q]uit\n\n"
	return s
}

// ------------ Helper functions to visualize the current operation and values ------------------
// var codeFormats = []string{
// 	"%d div 2^%d -> A",
// 	"%d xor %d -> B",
// 	"%d mod 8 -> B",
// 	"jnz %d",
// 	"%d xor %d -> B",
// 	"%d mod 8 -> out",
// 	"%d div 2^%d -> B",
// 	"%d div 2^%d -> C",
// }
//
// func getNumbers(operation int, operand int, reg Register) []any {
// 	switch operation {
// 	case 0:
// 		return []any{reg.A, combo(operand, reg)}
// 	case 1:
// 		return []any{reg.B, operand}
// 	case 2:
// 		return []any{combo(operand, reg)}
// 	case 3:
// 		return []any{operand}
// 	case 4:
// 		return []any{reg.B, reg.C}
// 	case 5:
// 		return []any{combo(operand, reg)}
// 	case 6:
// 		return []any{reg.A, combo(operand, reg)}
// 	case 7:
// 		return []any{reg.A, combo(operand, reg)}
// 	}
// 	return []any{}
// }

// func (m *model) View() string {
// 	s := "\nRegisters:\n"
// 	s += fmt.Sprintf("     A:  %d\n", M.register.A)
// 	s += fmt.Sprintf("     B:  %d\n", M.register.B)
// 	s += fmt.Sprintf("     C:  %d\n", M.register.C)
//
// 	s += "\nProgram counter:\n"
// 	msg := ""
// 	if M.PC >= len(M.program) {
// 		msg = "Outside scope!"
// 	}
// 	s += fmt.Sprintf("     PC: %d\t%s\n", M.PC, msg)
//
// 	s += "\nProgram Listing:\n"
// 	for i, c := range M.program {
// 		cursor := " "
// 		if M.PC == i {
// 			cursor = ">"
// 		}
// 		s += fmt.Sprintf("%s%d:  %s %d\n", cursor, i, opcodes[c.opcode], c.operand)
// 	}
//
//
// 	s += "\nOutput: \n"
// 	s += fmt.Sprintf("     %s", M.output)
//
// 	s += "\n\nControls: \n"
// 	s += "N: [N]ext\t"
// 	s += "A: [A]ll\t"
// 	s += "Q: [Q]uit\n"
// 	return s
// }
