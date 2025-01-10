package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	A = iota
	B
	C
)

type Register struct {
	A, B, C int
}

func (m model) RegStrings() [][]string {
	f := map[bool]string{true: "%d", false: "%o"}
	return [][]string{
		{fmt.Sprintf("A: "+f[m.dec], m.reg[A])},
		{fmt.Sprintf("B: "+f[m.dec], m.reg[B])},
		{fmt.Sprintf("C: "+f[m.dec], m.reg[C])},
	}
}

func (m model) StackStrings() [][]string {
	f := map[bool]string{true: "%d", false: "%o"}
	str := [][]string{}
	for _, i := range m.stack {
		str = append(str, []string{fmt.Sprintf(f[m.dec], i.answer), fmt.Sprintf("%v%v", m.seq[:i.idx+1], m.seq[i.idx+1:])})
		// tIO.Row("test", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(m.testSequence)), ","), "[]"))
	}
	return str
}

type model struct {
	next           func() (int, bool)
	hasNext        bool
	stack          Stack
	seq            []int           // The slice of program being used to find settings of the A register
	reg            [3]int          // Registers
	dirty          [3]bool         // Dirty flags for the registers
	dec            bool            // true if octal, otherwise decimal
	pc             int             // Program Counter
	output         []string        // A slice of the output
	newOutput      int             // Latest output
	inputProg      string          // Comma separated string of operations
	currOpString   string          // The string representation och the next operation
	length         int             // Length of the following slices
	stringFns      []func()        // The functions that generate the the currOpString
	codeFns        []func()        // The operations as functions (closures with each operand)
	src            []string        // Disassembled source code
	StepCmd        tea.Cmd         // Bubbletea Command for step one operation
	OnceCmd        tea.Cmd         // Bubbletea Command for running all operation ONCE
	AllCmd         tea.Cmd         // Bubbletea Command for running All operations
	StringCmd      tea.Cmd         // Bubbletea Command for visualizing the next operation
	NextCmd        tea.Cmd         // Bubbletea Command for finding next input
	codeEnumerator list.Enumerator // Bubbletea
	codeStyle      list.StyleFunc  // Bubbletea
	RegStyle       table.StyleFunc // Bubbletea
	srcList        *list.List      // Bubbletea
}
type pcUpdatedMsg int
type stringMsg string
type findMsg bool

func (m *model) Init() tea.Cmd {
	m.NextCmd = func() tea.Msg {
		tmp := 0
		if m.hasNext {
			m.dirty = [3]bool{}
			tmp, m.hasNext = m.next()
			if tmp > 0 {
				m.newOutput = tmp
			}
			if m.hasNext == false {
				m.reg[A] = m.newOutput
				m.stack.push(Item{m.newOutput, -1})
				m.dirty[A] = true
				m.reg[B] = 0
				m.reg[C] = 0
				m.newOutput = -1
			}
		}
		return findMsg(m.hasNext)
	}

	m.StepCmd = func() tea.Msg {
		if m.pc < m.length {
			m.dirty = [3]bool{}
			m.newOutput = -1
			m.codeFns[m.pc]()
		}
		return pcUpdatedMsg(m.pc)
	}

	m.OnceCmd = func() tea.Msg {
		m.dirty = [3]bool{}
		m.newOutput = -1
		for m.pc < m.length-1 {
			m.codeFns[m.pc]()
		}
		m.codeFns[m.pc]()
		return pcUpdatedMsg(m.pc)
	}
	m.AllCmd = func() tea.Msg {
		m.dirty = [3]bool{}
		m.newOutput = -1
		for m.pc < m.length {
			m.codeFns[m.pc]()
		}
		return pcUpdatedMsg(m.pc)
	}

	m.StringCmd = func() tea.Msg {
		if m.pc < m.length {
			m.stringFns[m.pc]() // pc must be updated, since this looks at the next operation
		}
		return stringMsg(m.currOpString)
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
		if row >= 0 && m.dirty[row] {
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
		if msg.String() == "q" {
			return m, tea.Quit
		}
		if msg.String() == "d" {
			m.dec = !m.dec
		}
		if m.hasNext {
			if msg.String() == "n" {
				return m, m.NextCmd
			}
		} else {
			switch msg.String() {
			case "s":
				return m, m.StepCmd
			case "a":
				return m, m.AllCmd
			case "o":
				return m, m.OnceCmd
			}
		}
	case pcUpdatedMsg:
		if int(msg) >= m.length {
			return m, tea.Quit
		}
		return m, m.StringCmd
	}

	return m, nil
}

func (m *model) View() string {
	tStack := table.New().Rows(m.StackStrings()...).Width(80)
	tStack.Headers("Stack")
	tReg := table.New().Rows(m.RegStrings()...).StyleFunc(m.RegStyle).Width(40)
	tReg.Headers("Registers")

	tIO := table.New().Headers("I/0", "").Width(40)
	if m.newOutput < 0 {
		tIO.Row("new", "")
	} else {
		tIO.Row("new", fmt.Sprintf("%d", m.newOutput))
	}
	tIO.Row("in", m.inputProg)
	tIO.Row("out", strings.Join(m.output, ","))

	tExec := table.New().Headers("Next opration to be executed").Width(80)
	if m.pc < m.length {
		tExec.Row(m.currOpString)
	} else {
		tExec.Row("Finished!")
	}

	s := lipgloss.JoinHorizontal(lipgloss.Top, tReg.Render(), tIO.Render())
	s = lipgloss.JoinVertical(lipgloss.Center, tStack.Render(), s, tExec.Render())

	s += "\n\n" + m.srcList.String()

	if m.hasNext {
		s += "\n\n\nControls: [N]ext [D]ec/Oct [Q]uit\n\n"
	} else {
		s += "\n\n\nControls: [S]tep [O]nce [A]ll [D]ec/Oct [Q]uit\n\n"
	}
	return s
}

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
