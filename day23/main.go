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

func main() {
	// p1("data.txt")
	p2("data.txt")
}

func p1(file string) {
	graph := buildGraph(file)

	res := allFullyConnectedPathsWithLength(graph, 3)
	lans := []string{}
	for _, path := range res {
		if hasNodeStartsWithT(path) {
			lan := strings.Join(path, ",")
			lans = append(lans, lan)
		}
	}
	slices.Sort(lans)
	i := 0
	for _, lan := range lans {
		fmt.Println(lan)
		i++
	}
	println("P1 (1108): ", i)

}

func p2(file string) {
	adjList := buildGraph(file)
	matrix := buildMatrix(file)
	path := longestFullyConnectedPath(adjList, matrix)
	slices.Sort(path)
	lan := strings.Join(path, ",")
	println("P2 (ab,cp,ep,fj,fl,ij,in,ng,pl,qr,rx,va,vf): ", lan)
}

func buildGraph(file string) map[string][]string {
	adjacencyList := map[string][]string{}

	for a, b := range readEdges(file) {
		if _, ok := adjacencyList[a]; !ok {
			adjacencyList[a] = []string{}
		}
		if _, ok := adjacencyList[b]; !ok {
			adjacencyList[b] = []string{}
		}
		adjacencyList[a] = append(adjacencyList[a], b)
		adjacencyList[b] = append(adjacencyList[b], a)
	}
	return adjacencyList
}

func buildMatrix(file string) map[string]map[string]bool {
	adjacencyMap := map[string]map[string]bool{}

	for a, b := range readEdges(file) {
		if _, ok := adjacencyMap[a]; !ok {
			adjacencyMap[a] = map[string]bool{}
		}
		if _, ok := adjacencyMap[b]; !ok {
			adjacencyMap[b] = map[string]bool{}
		}
		adjacencyMap[a][b] = true
		adjacencyMap[b][a] = true
	}
	return adjacencyMap
}

// ======================= Outer loop, ---------------->
func longestFullyConnectedPath(graph map[string][]string, matrix map[string]map[string]bool) []string {
	longestPathSoFar := []string{} // A list if sub graphs

	print("Length: ", len(graph))
	i := 0
	// Loop over all nodes
	for node := range graph {
		print("  ", i)
		i++
		res := exploreLongestSubgraph(graph, node, matrix)
		if len(res) > len(longestPathSoFar) {
			longestPathSoFar = res
		}
	}
	return longestPathSoFar
}

// -----------------------> ...Inner BFS
// Return the longest path, "starting" at start, where every node is interconnected
func exploreLongestSubgraph(graph map[string][]string, start string, matrix map[string]map[string]bool) []string {
	q := []Q{{start, []string{start}}}
	longestPathsForStart := []string{}
	seen := map[string]bool{} //start: true}

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]

		if seen[curr.id] {
			continue
		}

		// Goal
		if len(curr.path) > len(longestPathsForStart) {
			fmt.Printf("len Q: %d, len curr: %d, len longest: %d...", len(q), len(curr.path), len(longestPathsForStart))
			longestPathsForStart = append([]string{}, curr.path...)
			fmt.Println(longestPathsForStart)
			seen[curr.id] = true
		}

		// Neighbours
		for _, n := range graph[curr.id] {
			if slices.Contains(curr.path, n) { // Could use a map for performance?!
				// if curr.seen[n] { // Could use a map for performance?!
				continue
			}
			if !isConnectedWithAll(graph, n, curr.path, matrix) {
				continue
			}
			path := append([]string{n}, curr.path...)
			// seen := maps.Clone(curr.seen)
			// seen[n] = true
			q = append(q, Q{n, path})

		}

	}

	return longestPathsForStart
}

// ======================= Outer loop, ---------------->
func allFullyConnectedPathsWithLength(graph map[string][]string, k int) [][]string {
	seenPath := map[string]bool{} // To avoid duplicate sub graphs, needs to live outside BFS
	subgrapghs := [][]string{}    // A list if sub graphs

	// Loop over all nodes
	for node := range graph {
		subgraphs := exploreSubgraph(graph, k, node, seenPath)
		subgrapghs = append(subgrapghs, subgraphs...)
	}
	return subgrapghs
}

type Q struct {
	id   string
	path []string
}

// -----------------------> ...Inner BFS
func exploreSubgraph(graph map[string][]string, k int, start string, seen map[string]bool) [][]string {
	q := []Q{{start, []string{start}}}
	paths := [][]string{}
	visitedNode := map[string]bool{}

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]

		visitedNode[curr.id] = true

		// Goal: Length k...
		if len(curr.path) == k {

			// Avoid storing duplicate paths
			slices.Sort(curr.path)
			key := strings.Join(curr.path, ",")
			if seen[key] {
				continue
			}
			seen[key] = true

			// Avoid paths that are NOT inter-connected, (Could have been a simpler check for only 3 nodes)
			if !allInterConnected(graph, curr.path) {
				continue
			}

			// Add the path to the result
			paths = append(paths, curr.path)
			continue
		}

		// Neighbours
		for _, n := range graph[curr.id] {
			if visitedNode[n] {
				continue
			}
			path := append([]string{n}, curr.path...)
			q = append(q, Q{n, path})
		}
	}
	return paths
}

// ----------- Util------------
func hasNodeStartsWithT(path []string) bool {
	for _, node := range path {
		if node[0] == 't' {
			return true
		}
	}
	return false
}
func isConnectedWithAll(graph map[string][]string, node string, path []string, matrix map[string]map[string]bool) bool {
	// If any node is NOT connected to ALL other -> early return
	for _, pathNode := range path { // Loop over all nodes in the path so far
		// Check all other edges
		if !matrix[node][pathNode] {
			return false
		}
	}
	return true
}
func allInterConnected(graph map[string][]string, path []string) bool {
	// If any node is NOT connected to ALL other -> early return
	for _, from := range path { // Loop over all FROM nodes:  [A,B,C]
		for _, to := range path { // Loop over all TO nodes: [A,B,C]
			// Skip if same
			if from == to {
				continue
			}
			// Check all other edges
			if !slices.Contains(graph[from], to) {
				return false
			}
		}
	}
	return true
}

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
