package main

import (
	"cmp"
	"fmt"
	"io"
	"iter"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
)

func main() {
	// p1("data.txt")
	p1Alternative("data.txt")
	// p2("data.txt")
	p2Alternative("data.txt")
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

// Hyper Neutrino
func p1Alternative(file string) {
	m := buildMatrix(file)

	set := map[string]bool{}
	// Loop over all nodes
	for node1 := range m {

		// Loop over its neighbours
		for node2 := range m[node1] {
			// Avoid going back
			if node2 == node1 {
				continue
			}

			// ...one more time
			for node3 := range m[node2] {
				// Avoid going back
				if node3 == node1 {
					continue
				}

				// Exclude loops thar are NOT closed
				if !m[node1][node3] {
					continue
				}

				if node1[0] == 't' || node2[0] == 't' || node3[0] == 't' {
					lan := []string{node1, node2, node3}
					slices.Sort(lan)
					set[strings.Join(lan, ",")] = true
				}
			}
		}
	}
	println("P1 (1108): ", len(set))
}

func p2(file string) {
	adjList := buildGraph(file)
	matrix := buildMatrix(file)
	path := longestFullyConnectedPath(adjList, matrix)
	slices.Sort(path)
	lan := strings.Join(path, ",")
	println("P2 (ab,cp,ep,fj,fl,ij,in,ng,pl,qr,rx,va,vf): ", lan)
}

// Hyper Neutrino
func p2Alternative(file string) {
	m := buildMatrix(file)
	sets = map[string]bool{}
	for start := range m {
		recursiveSearch(m, start, map[string]bool{start: true})
	}

	length := func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	}

	max := slices.MaxFunc(slices.Collect(maps.Keys(sets)), length)

	fmt.Printf("P2 (ab,cp,ep,fj,fl,ij,in,ng,pl,qr,rx,va,vf): %s (%d)\n", max, len(sets))
}

var sets map[string]bool

func recursiveSearch(m map[string]map[string]bool, node string, req map[string]bool) {
	key := strings.Join(slices.Sorted(maps.Keys(req)), ",")
	if sets[key] {
		return
	}
	sets[key] = true
	for neighbour := range m[node] {
		if req[neighbour] {
			continue
		}
		all := true
		for n := range req {
			if !m[n][neighbour] {
				all = false
			}
		}
		if !all {
			continue
		}
		cpy := maps.Clone(req)
		cpy[neighbour] = true
		recursiveSearch(m, neighbour, cpy)
	}

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
	accepted := map[string]bool{}

	print("Total: ", len(graph))
	i := 0
	// Loop over all nodes
	for node := range graph {
		if accepted[node] {
			continue
		}
		print("  ", i)
		i++

		res := exploreLongestSubgraph(graph, node, matrix, accepted)

		for _, r := range res {
			accepted[r] = true
		}
		if len(res) > len(longestPathSoFar) {
			longestPathSoFar = res
		}
	}
	return longestPathSoFar
}

// -----------------------> ...Inner BFS
// Return the longest path, "starting" at start, where every node is interconnected
func exploreLongestSubgraph(graph map[string][]string, start string, matrix map[string]map[string]bool, accepted map[string]bool) []string {
	q := []Q{{start, []string{start}}}
	longestPathsForStart := []string{}
	seen := map[string]bool{}

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]

		if seen[curr.id] {
			continue
		}

		// Goal
		if len(curr.path) > len(longestPathsForStart) {
			longestPathsForStart = append([]string{}, curr.path...)
			seen[curr.id] = true // ?? Not really sure about this check...would have wanted it outside??
		}

		// Neighbours
		for _, n := range graph[curr.id] {
			if accepted[n] {
				continue // Skip if neighbour inside an already accepted solution
			}
			if slices.Contains(curr.path, n) {
				continue // Skip if neighbour is with current path
			}
			if !isConnectedWithAll(graph, n, curr.path, matrix) {
				continue
			}
			path := append([]string{n}, curr.path...)
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

// func allisConnected(node string, path map[string]bool, matrix map[string]map[string]bool) bool
//
//		// If any node is NOT connected to ALL other -> early return
//		for _, pathNode := range path { // Loop over all nodes in the path so far
//			// Check all other edges
//			if !matrix[node][pathNode] {
//				return false
//			}
//		}
//		return true
//	}
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
