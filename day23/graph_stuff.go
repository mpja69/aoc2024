package main

import "slices"

// =========== De vanliga mönstren på : , CANconstruct, COUNTconstruct, HOWconstruct...
// ----------- På grafer: HASpath, ...
func hasPath(graph map[string][]string, start string, end string, seen map[string]bool) bool {
	if seen[start] {
		return false
	}
	if start == end {
		return true
	}
	seen[start] = true
	for _, n := range graph[start] {
		if hasPath(graph, n, end, seen) {
			return true
		}
	}
	return false
}

// === Exempel på hur man kan kombinera K element ur en lista, med "back tracking"
func kCombinations[T any](list []T, k int) [][]T {
	// TODO: Testa att använda en channel till att mata ut svaret.
	var ans [][]T
	var sol []T
	var backtrack func(int)

	if k < 0 || k > len(list) {
		return [][]T{{}}
	}

	backtrack = func(n int) {
		if len(sol) == k {
			ans = append(ans, slices.Clone(sol))
			return
		}
		left := len(list) - n
		stillNeed := k - len(sol)
		if left > stillNeed {
			backtrack(n + 1)
		}
		sol = append(sol, list[n])
		backtrack(n + 1)
		sol = sol[:len(sol)-1]
	}
	backtrack(0)
	return ans
}

// === Exempel på en rekursiv som finner alla kombinationer ur en lista, i olika längder
func allCombinations(list []int) [][]int {
	if len(list) == 0 {
		return [][]int{{}} // NOTE: To create a list with an empty list [ [] ], initialize it "all" the way
	}

	first := list[0]
	rest := list[1:]
	combsRest := allCombinations(rest)
	combsAll := [][]int{}

	for _, comb := range combsRest {
		oneCombWithFirst := append(slices.Clone(comb), first) // Skapa en temp-lista med den första och de andra
		combsAll = append(combsAll, oneCombWithFirst)         // Lägg till listan i listan
	}
	return slices.Concat(combsRest, combsAll)
}
