package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
)

type model struct {
	register       Register
	PC             int
	output         string
	code           []func()
	src            []string
	Next           tea.Cmd
	All            tea.Cmd
	codeEnumerator list.Enumerator
	codeStyle      list.StyleFunc
}
type pcMsg int

func (m *model) Init() tea.Cmd {
	m.Next = func() tea.Msg {
		if m.PC < len(m.code) {
			m.code[m.PC]()
		}
		return pcMsg(m.PC)
	}

	m.All = func() tea.Msg {
		for m.PC < len(m.code) {
			m.code[m.PC]()
		}
		return pcMsg(m.PC)
	}

	m.codeEnumerator = func(_ list.Items, i int) string {
		if i == m.PC {
			return fmt.Sprintf("-> %d:", i)
		}
		return fmt.Sprintf("   %d:", i)
	}

	m.codeStyle = func(_ list.Items, i int) lipgloss.Style {
		if i == m.PC {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		}
		return lipgloss.NewStyle()
	}

	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":
			return m, m.Next
		case "a":
			return m, m.All
		}
	case pcMsg:
		if int(msg) >= len(m.code) {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) View() string {
	tr := table.New()
	tr.Headers("Registers", "Value")
	tr.Row("A", fmt.Sprint(m.register.A))
	tr.Row("B", fmt.Sprint(m.register.B))
	tr.Row("C", fmt.Sprint(m.register.C))
	sr := tr.String() + "\n"

	tpc := table.New().Headers("Program Counter")
	msg := fmt.Sprintf("%d", m.PC)
	if m.PC >= len(m.code) {
		msg = "Finished"
	}
	tpc.Row(msg)
	spc := tpc.String() + "\n"

	to := table.New().Headers("Output")
	so := to.Row(m.output).String()

	s := lipgloss.JoinHorizontal(lipgloss.Top, sr, spc, so)

	s += "\nProgram Listing:\n"
	l := list.New(m.src).Enumerator(m.codeEnumerator).ItemStyleFunc(m.codeStyle)
	s += l.String() + "\n\n"

	s += lipgloss.JoinHorizontal(lipgloss.Top, "Controls: ", "[N]ext ", "[A]ll ", "[Q]uit")
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
