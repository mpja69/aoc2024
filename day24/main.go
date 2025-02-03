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
//
//	An XOR of the input x and y of this bit
//	An AND of the input x and y of this bit
//	an XOR of (1) and the incoming carry (this should be the z for this bit)
//	an AND of (1) and the incoming carry
//	an OR of both ANDs (this is the output carry; it will be the input carry for the next bit).
func p2(exprStrings []string) {

	// Make a map of the expressions
	expr = map[string]string{}
	for _, line := range exprStrings {
		l := strings.Split(line, " -> ")
		expr[l[1]] = l[0]
	}

	getWiresWithErrors2(exprStrings)
}

// From the python example
func getWiresWithErrors2(expressions []string) { //map[string]string) {

	wires := []string{}
	for _, expr := range expressions {
		a, op, b, out := parseExpression(expr)

		// x != "XOR" and c[0] == 'z' and c != "z45"]
		if op != "XOR" && out[0] == 'z' && out != "z45" {
			wires = append(wires, out)
		}

		// x == "XOR" and all(d[0] not in 'xyz' for d in (a, b, c)) or
		if op == "XOR" && a[0] != 'x' && a[0] != 'y' && out[0] != 'z' {
			wires = append(wires, out)
		}

		// x == "AND" and not "x00" in (a, b) and r(c, 'XOR') or
		if op == "AND" && a != "x00" && b != "x00" && connectedTo(out, "XOR") {
			wires = append(wires, out)
		}

		// x == "XOR" and not "x00" in (a, b) and r(c, 'OR') or
		if op == "XOR" && a != "x00" && b != "x00" && connectedTo(out, "OR") {
			wires = append(wires, out)
		}
	}

	slices.Sort(wires)
	wires = slices.Compact(wires)
	fmt.Println(strings.Join(wires, ","))
}

// Helper function - used in some if-stmts below
// True if any expr exist, that has any input 'ab' AND the operation 'op'
func connectedTo(ab, op string) bool {
	for _, v := range expr {
		a, x, b := parse(v)
		if op == x && (ab == a || ab == b) {
			return true
		}
	}
	return false
}

// From the bash-example
func getWiresWithErrors(expr map[string]string) {
	// get the z-outputs, except the last one, which expr has not "XOR"
	// $(cat input.txt | grep '> z' | grep -v 'XOR' | grep -v 'z45$' | awk '{print $5}')
	// INFO: expr.filter(k[0]=='z').filter(k!="z45").filter("XOR" not in v).map(k)
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

	// get the the"non-z-outputs" of all  expr with "XOR" except expr with "x" or "y"
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

	// get a list of "OR"-ops, but get JUST the 2 inputs
	// INPUT_OF_OR=$(cat input.txt | grep ' OR ' | awk '{ print $1; print $3 }' | sort -u)
	INPUT_OF_OR := []string{}
	for _, v := range expr {
		a, op, b := parse(v)
		if op != "OR" {
			continue
		}
		INPUT_OF_OR = append(INPUT_OF_OR, a, b)
	}

	// Find the outputs of expr that has "AND" but not the first inputs
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

	// Take the items that exclusively exists in each list (of the above inputs and outputs)
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

	// Finally take all candidates, sort them, and keep unique ones
	// echo $CANDIDATE_1 $CANDIDATE_2 $CANDIDATE_3 | tr ' ' '\n' | sort -u | tr '\n' ',' | sed -e 's/,$//'
	// INFO: allUniqueAndSorted(...)
	wires := slices.Concat(CANDIDATE_1, CANDIDATE_2, CANDIDATE_3)
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

func parseExpression(e string) (a, op, b, out string) {
	p := strings.Split(e, " -> ")
	s := strings.Split(p[0], " ")
	return s[0], s[1], s[2], p[1]

}

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
