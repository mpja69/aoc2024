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
	p1()
	// p2()
}

func p1() {
	graph := buildGraph("data.txt")

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
	println("P1: ", i)

}

func buildGraph(file string) map[string][]string {
	adjacencyList := map[string][]string{} // NOTE: Har man inte en adj-list i grafen ovan??

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
