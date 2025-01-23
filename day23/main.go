package main

import (
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"slices"
	"strings"
)

// Implementation of Graph from the book
type Node struct {
	idx   int
	edges map[int][]Edge
	label string
}

type Edge struct {
	from  int
	to    int
	weigt float64
}

type Graph struct {
	nodes []Node
}

func main() {

	graph := buildGraph("sample.txt")
	// fmt.Println(len(graph))
	// fmt.Println(graph)

	res := allFullyConnectedPathsWithLength(graph, 3)
	for _, sl := range res {
		fmt.Println(sl)
	}

}

func buildGraph(file string) map[string][]string {
	adjacencyList := map[string][]string{} // NOTE: Har man inte en adj-list i grafen ovan??

	for a, b := range readEdges(file) {
		if _, ok := adjacencyList[a]; !ok {
			adjacencyList[a] = []string{}
		}
		adjacencyList[a] = append(adjacencyList[a], b)
		if _, ok := adjacencyList[b]; !ok {
			adjacencyList[b] = []string{}
		}
		adjacencyList[b] = append(adjacencyList[b], a)
	}
	return adjacencyList
}

func getNodes(file string) []string {
	nodeList := []string{}
	seen := map[string]bool{}
	for a, b := range readEdges(file) {
		if !seen[a] {
			nodeList = append(nodeList, a)
			seen[a] = true
		}
		if !seen[b] {
			nodeList = append(nodeList, b)
			seen[b] = true
		}
	}
	return nodeList
}

type QI struct {
	id  string
	lan []string
	// length int
}

// ======================= Vanligt mönster....med en yttre loop, ---------------->
func allFullyConnectedPathsWithLength(graph map[string][]string, k int) [][]string {
	seen := map[string]bool{} // Unique string representation of each sub graph -> To avoid duplicate sub graphs
	res := [][]string{}       // A list if sub graphs

	// Loop over all nodes
	for node := range graph {
		subgraphs := exploreSubgraph(graph, k, node, seen)
		res = append(res, subgraphs...)
	}
	return res
}

type StringSet map[string]bool
type SetSet map[*StringSet]bool

// -----------------------> ...som kallar den inre
// 1. Starta på en nod
// 2. Hitta alla _unika_ paths med längd K...
func exploreSubgraph(graph map[string][]string, k int, start string, seen map[string]bool) [][]string {
	q := []QI{{start, []string{start}}}
	res := [][]string{}
	visited := map[string]bool{}

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]

		// Avoid going back
		// if visited[curr.id] {
		// 	continue
		// }
		visited[curr.id] = true

		// Goal: Length k...
		if len(curr.lan) == k {
			path := curr.lan

			// Avoid storing duplicate sub graphs
			slices.Sort(path)
			key := strings.Join(path, ",")
			if seen[key] {
				continue
			}
			seen[key] = true

			// Avoid paths that are NOT inter-connected
			if !allInterConnected(graph, path) {
				continue
			}

			// Add the path to the reult list
			res = append(res, path)
			continue
		}

		// Neighbours
		for _, n := range graph[curr.id] {
			if visited[n] {
				continue
			}
			path := append(slices.Clone(curr.lan), n)
			q = append(q, QI{n, path})
		}
	}
	return res
}

func allInterConnected(graph map[string][]string, path []string) bool {
	// If any node is NOT connected to ALL other -> early return
	for _, from := range path { // Loop over all FROM nodes:  [A,B,C]
		for _, to := range path { // Loop over all TO nodes: [A,B,C]
			// Skip if same
			if from == to {
				continue
			}
			// Check all other edges TODO: Om path var ett set gåt kollen fortare
			if !slices.Contains(graph[from], to) {
				return false
			}
		}
	}
	return true
}

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

// ----------- Util------------
// Iterator för att läsa in kanterna
func readEdges(s string) iter.Seq2[string, string] {
	f, err := os.Open(s)
	if err != nil {
		log.Fatal("readNumber(): ", err)
	}
	return func(yield func(string, string) bool) {
		defer f.Close()
		line := "      "
		for _, err := fmt.Fscanln(f, &line); err != io.EOF; _, err = fmt.Fscan(f, &line) {
			if err != nil {
				log.Fatal("readIterator(): ", err)
				return
			}
			if !yield(line[:2], line[3:]) {
				return
			}
		}
	}

}

// func foo() {
// 	// List of all LANS (internally ordered)...5 computers in each
// 	lans := [][]string{}
// 	for _, node := range nodeList {
// 		set := append(adjacencyList[node], node)
// 		slices.Sort(set)
// 		lans = append(lans, set)
// 	}
// 	fmt.Println("lans:", len(lans))
//
// 	// Set of LANS
// 	lanSets := map[string][]string{}
// 	for _, lan := range lans {
// 		lanStr := strings.Join(lan, "")
// 		lanSets[lanStr] = lan
// 	}
// 	fmt.Println("Set of lans:", len(lanSets))
//
// 	// lanSets := map[string]bool
// 	for _, lan := range lanSets {
// 		ans := kCombinations(lan, 3)
// 		fmt.Println(ans)
// 	}
// }
