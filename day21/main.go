package main

import (
	"bytes"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/mpja69/aoc2024/day21/keypad"
)

// TODO:	Prova Writer istället
//   - WriteByte('0') och låt den gå genom kedjan, sedan WriteByte('2')
//   - WriteString("029A")
//     Exempel:
//     var buf = bytes.Buffer,		//	skapa en Buffer
//     p := keypad.NewNumad(buf)	//	skapa min "numpad writer"
//     p.WriteByte('0')				//	Skriv genom min writer -> buf
//     buf.WriteTo(os.Stdout)		//	skriv ut från buf till skärmen

func main() {
	data, _ := os.ReadFile("data.txt")
	data = bytes.TrimSpace(data)
	lines := strings.Split(string(data), "\n")

	sum := 0
	for _, line := range lines {
		sum += runLine(line)
	}
	fmt.Println("P1: ", sum)
}

func runLine(line string) int {
	fmt.Printf("%s: ", line)

	// First a numpad
	kp := keypad.NewKeypad(keypad.NumberLayout)
	numLines := kp.GetPossibleInputsWithOutput(line)

	// Then a dir pad
	kp = keypad.NewKeypad(keypad.DirectionLayout)
	dirLines := []string{}
	for _, line := range numLines {
		dirLines = append(dirLines, kp.GetPossibleInputsWithOutput(line)...)
	}
	// Only pick the shortes ones
	lenCmp := func(a, b string) int { return cmp.Compare(len(a), len(b)) }
	slices.SortFunc(dirLines, lenCmp)

	// Finally another dir pad
	dirLines2 := []string{}
	for line := range onlyShortest(dirLines) {
		dirLines2 = append(dirLines2, kp.GetPossibleInputsWithOutput(line)...)
	}

	// Find the length of the shortes ones
	length := len(slices.MinFunc(dirLines2, lenCmp))
	value := value(line)
	println(length, value, length*value)
	return length * value
}

func onlyShortest(s []string) func(func(string) bool) {
	return func(yield func(string) bool) {
		for _, line := range s {
			if len(line) > len(s[0]) {
				break
			}
			yield(line)
		}
	}
}
func value(code string) int {
	c, _ := strconv.Atoi(string(code[:3]))
	return c
}
