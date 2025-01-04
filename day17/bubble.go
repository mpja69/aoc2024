package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
)

type model struct {
	register Register
	// program  []Combo
	PC     int
	output string
	ops    []func()
	src    []string
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "n":
			return m, Next
		case "a":
			return m, All
		}
	}

	return m, nil
}

func Next() tea.Msg {
	if m.PC < len(m.ops) {
		m.ops[m.PC]()
	}
	return m.PC
}

func All() tea.Msg {
	for m.PC < len(m.ops) {
		m.ops[m.PC]()
	}
	return m.PC
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
	if m.PC >= len(m.ops) {
		msg = "Finished"
	}
	tpc.Row(msg)
	spc := tpc.String() + "\n"

	to := table.New().Headers("Output")
	so := to.Row(m.output).String()

	s := lipgloss.JoinHorizontal(lipgloss.Top, sr, spc, so)

	s += "\nProgram Listing:\n"
	l := list.New(m.src).Enumerator(ProgramEnumerator).ItemStyleFunc(ProgramStyle)
	s += l.String() + "\n"

	s += "\n\nControls:\n"
	s += "N: [N]ext\t"
	s += "A: [A]ll\t"
	s += "Q: [Q]uit\n"
	return s
}

func ProgramEnumerator(_ list.Items, i int) string {
	if i == m.PC {
		return fmt.Sprintf("-> %d:", i)
	}
	return fmt.Sprintf("   %d:", i)
}

func ProgramStyle(_ list.Items, i int) lipgloss.Style {
	if i == m.PC {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	}
	return lipgloss.NewStyle()
}

// func (m *model) View() string {
// 	s := "\nRegisters:\n"
// 	s += fmt.Sprintf("     A:  %d\n", m.register.A)
// 	s += fmt.Sprintf("     B:  %d\n", m.register.B)
// 	s += fmt.Sprintf("     C:  %d\n", m.register.C)
//
// 	s += "\nProgram counter:\n"
// 	msg := ""
// 	if m.PC >= len(m.program) {
// 		msg = "Outside scope!"
// 	}
// 	s += fmt.Sprintf("     PC: %d\t%s\n", m.PC, msg)
//
// 	s += "\nProgram Listing:\n"
// 	for i, c := range m.program {
// 		cursor := " "
// 		if m.PC == i {
// 			cursor = ">"
// 		}
// 		s += fmt.Sprintf("%s%d:  %s %d\n", cursor, i, opcodes[c.opcode], c.operand)
// 	}
//
//
// 	s += "\nOutput: \n"
// 	s += fmt.Sprintf("     %s", m.output)
//
// 	s += "\n\nControls: \n"
// 	s += "N: [N]ext\t"
// 	s += "A: [A]ll\t"
// 	s += "Q: [Q]uit\n"
// 	return s
// }
