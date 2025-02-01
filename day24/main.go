package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var cache map[string]int
var expr map[string]string

func main() {
	data, _ := os.ReadFile("data.txt")
	parts := strings.Split(strings.TrimSpace(string(data)), "\n\n")
	// inputStrings := strings.Split(parts[0], "\n")
	exprStrings := strings.Split(parts[1], "\n")

	// p1(inputStrings, exprStrings)
	p2(exprStrings)
}

func p1(inputStrings, exprStrings []string) {
	// Parse the the inputs, and store in the cache
	cache = map[string]int{}
	for _, line := range inputStrings {
		l := strings.Split(line, ": ")
		val, _ := strconv.Atoi(l[1])
		cache[l[0]] = val
	}

	// Make a map of the expressions
	expr = map[string]string{}
	for _, line := range exprStrings {
		l := strings.Split(line, " -> ")
		expr[l[1]] = l[0]
	}

	// Solve the expressions
	for k := range expr {
		solve(k)
	}
	fmt.Println("P1: ", getVal('z'))
}

// For a full adder we expect the following gates:
//		An XOR of the input x and y of this bit
//		An AND of the input x and y of this bit
//		an XOR of (1) and the incoming carry (this should be the z for this bit)
//		an AND of (1) and the incoming carry
//		an OR of both ANDs (this is the output carry; it will be the input carry for the next bit).

// Hyper Neutrino:	Skip the values and only analyze the grid...and how it relates to an adder, and find anomolies
// Another idea:	Fill the input with with: zeros, ones and ordered mix. -> analyze output and find anomolies
func p2(exprStrings []string) {

	// Make a map of the expressions
	expr = map[string]string{}
	for _, line := range exprStrings {
		l := strings.Split(line, " -> ")
		expr[l[1]] = l[0]
	}

	getErrWires2(expr)
}

// From the bash-example
func getErrWires(expr map[string]string) {
	// get a list where: output is  "> z", but no "XOR"-ops, and not the last bit...get JUST the output-names
	// $(cat input.txt | grep '> z' | grep -v 'XOR' | grep -v 'z45$' | awk '{print $5}')
	// INFO: expr.filter(k[0]=='z').filter(k!="z45").filter("XOR" not in v).map(k)
	// python x != "XOR" and c[0] == 'z' and c != "z45"]

	CANDIDATE_1 := []string{}
	for k, v := range expr {
		if k[0] != 'z' {
			continue
		}
		if strings.Contains(v, "XOR") {
			continue
		}
		if k == "z45" {
			continue
		}
		CANDIDATE_1 = append(CANDIDATE_1, k)
	}
	fmt.Println("1:", CANDIDATE_1)

	// get a list of all "XOR"-ops, but not the lines starting with "x" or "y", and not output "> z"...get JUST the output-names
	// CANDIDATE_2=$(cat input.txt | grep ' XOR ' | grep -v '^x' | grep -v '^y' | grep -v '> z' | awk '{print $5}')
	// INFO: expr.filter(v[0]!='x').filter(v[0]!='y').filter("XOR" not in v).map(k)
	// python x == "XOR" and all(d[0] not in 'xyz' for d in (a, b, c)) or
	CANDIDATE_2 := []string{}
	for k, v := range expr {
		if !strings.Contains(v, "XOR") {
			continue
		}
		if v[0] == 'x' || v[0] == 'y' || k[0] == 'z' {
			continue
		}
		CANDIDATE_2 = append(CANDIDATE_2, k)
	}
	fmt.Println("2:", CANDIDATE_2)

	// get a list of "OR"-ops, but get JUST the 2 inputs, and sort them and remove duplilcates
	// INPUT_OF_OR=$(cat input.txt | grep ' OR ' | awk '{ print $1; print $3 }' | sort -u)
	INPUT_OF_OR := []string{}
	for _, v := range expr {
		a, op, b := parse(v)
		if op != "OR" {
			continue
		}
		INPUT_OF_OR = append(INPUT_OF_OR, a, b)
	}

	// # Exclude the first "input"-"AND", but include all other "AND"-ops, ...get JUST the output-names (sort and remove duplicates)
	// OUTPUT_OF_AND=$(cat input.txt | grep -v 'x00 AND y00' | grep ' AND ' | awk '{ print $5 }' | sort -u)
	OUTPUT_OF_AND := []string{}
	for k, v := range expr {
		_, op, _ := parse(v)
		if v == "x00 AND y00" {
			continue
		}
		if op != "AND" {
			continue
		}
		OUTPUT_OF_AND = append(OUTPUT_OF_AND, k)
	}

	// # From all "OR"s and "AND"s, (replace WHITE with NEWLINE), get just the unique ones
	// CANDIDATE_3=$(comm -3 <(echo $INPUT_OF_OR | tr ' ' '\n') <(echo $OUTPUT_OF_AND | tr ' ' '\n'))
	// INFO: CANDIDATE_3 := XOR(INPUT_OF_OR, OUTPUT_OF_AND)
	CANDIDATE_3 := []string{}
	for _, item := range append(INPUT_OF_OR, OUTPUT_OF_AND...) {
		if slices.Contains(OUTPUT_OF_AND, item) && slices.Contains(INPUT_OF_OR, item) {
			continue
		}
		CANDIDATE_3 = append(CANDIDATE_3, item)
	}
	fmt.Println("3:", CANDIDATE_3)

	// # From all candidates: (replace WHITE with NL), sort and remove duplicates, then replace NL with ",", and remove last ","
	// echo $CANDIDATE_1 $CANDIDATE_2 $CANDIDATE_3 | tr ' ' '\n' | sort -u | tr '\n' ',' | sed -e 's/,$//'
	// INFO: allUniqueAndSorted(...)
	wires := slices.Concat(CANDIDATE_1, CANDIDATE_2, CANDIDATE_3)
	slices.Sort(wires)
	wires = slices.Compact(wires)

	fmt.Println(strings.Join(wires, ","))

}

// From the python example
func getErrWires2(expr map[string]string) {

	r := func(output, op string) bool {
		for _, v := range expr {
			a, x, b := parse(v)
			if op == x && (output == a || output == b) {
				return true
			}
		}
		return false
	}
	wires := []string{}
	for k, v := range expr {
		a, op, b := parse(v)

		// x != "XOR" and c[0] == 'z' and c != "z45"]
		if op != "XOR" && k[0] == 'z' && k != "z45" {
			wires = append(wires, k)
		}

		// x == "XOR" and all(d[0] not in 'xyz' for d in (a, b, c)) or
		if op == "XOR" && v[0] != 'x' && v[0] != 'y' && k[0] != 'z' {
			wires = append(wires, k)
		}

		// x == "AND" and not "x00" in (a, b) and r(c, 'XOR') or
		if op == "AND" && a != "x00" && b != "x00" && r(k, "XOR") {
			wires = append(wires, k)
		}

		// x == "XOR" and not "x00" in (a, b) and r(c, 'OR') or
		if op == "XOR" && a != "x00" && b != "x00" && r(k, "OR") {
			wires = append(wires, k)
		}
	}

	slices.Sort(wires)
	wires = slices.Compact(wires)
	fmt.Println(strings.Join(wires, ","))
}

func solve(ans string) int {
	if val, ok := cache[ans]; ok {
		return val
	}

	a, op, b := parse(expr[ans])
	if _, ok := cache[a]; !ok {
		cache[a] = solve(a)
	}

	if _, ok := cache[b]; !ok {
		cache[b] = solve(b)
	}

	switch op {
	case "AND":
		cache[ans] = cache[a] & cache[b]
	case "OR":
		cache[ans] = cache[a] | cache[b]
	case "XOR":
		cache[ans] = cache[a] ^ cache[b]
	}

	return cache[ans]
}

// func calulate(a, b string) int {}
func parse(e string) (a, op, b string) {
	s := strings.Split(e, " ")
	return s[0], s[1], s[2]

}

func getVal(keyStart byte) int {
	keys := []string{}

	// Find all starting with "key"
	for k := range cache {
		if k[0] != keyStart {
			continue
		}
		keys = append(keys, k)
	}

	// order them
	slices.Sort(keys)

	// Calculate value
	res := 0
	for _, key := range slices.Backward(keys) {
		res = res | cache[key]
		res = res << 1
	}
	return res >> 1
}
