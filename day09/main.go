package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("d.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := strings.Trim(string(data), " \n")
	lines := strings.Split(input, "\n")

	fmt.Printf("Part 1: (1928, 6201130364722): %d\n", part1(lines[0]))
	fmt.Printf("Part 2: (2858, 6221662795602): %d\n", part2(lines[0]))
}

func part1(diskmap string) int {
	blocks, _, _ := inflateBlocks(diskmap)
	partitionSingleBlocks(blocks)
	return checkSum(blocks)
}

func part2(diskmap string) int {
	blocks, fileSizes, fileIdx := inflateBlocks(diskmap)
	partitionFiles(blocks, fileSizes, fileIdx)
	return checkSum(blocks)
}

func inflateBlocks(diskmap string) ([]int, []int, []int) {
	fileSizes := make([]int, len(diskmap)/2+1)
	fileIdx := make([]int, len(diskmap)/2+1)
	blocks := []int{}
	isFile := true
	id := 0
	for _, digit := range diskmap {
		nbrBlocks := int(digit) - '0'
		if isFile {
			fileIdx[id] = len(blocks)
			for block := 0; block < nbrBlocks; block++ {
				blocks = append(blocks, id)
			}
			fileSizes[id] = nbrBlocks
			id++
		} else {
			for block := 0; block < nbrBlocks; block++ {
				blocks = append(blocks, -1)
			}
		}
		isFile = !isFile
	}
	return blocks, fileSizes, fileIdx
}

func partitionSingleBlocks(b []int) {
	src := len(b) - 1
	dst := 0

	for {
		// Move forward to next free block
		for b[dst] >= 0 {
			dst++
		}
		//Move beckward to next block (to move)
		for b[src] < 0 {
			src--
		}

		// fmt.Printf("%v\n", b)
		// If now more blocks to move
		if dst >= src {
			return
		}
		// Swap blocks
		b[dst], b[src] = b[src], b[dst]
	}
}
func partitionFiles(blocks []int, sizes []int, indices []int) {
	// printBlocks(blocks)
	id := len(sizes) - 1
	for id >= 0 {
		// Always look for free space by starting from the beginning
		startFree := 0

		// Identify file to move
		nbrBlocks := sizes[id]
		idx := indices[id]

		// Search for a slot big enough
		for startFree < len(blocks) {
			// Find the start of free space
			for startFree < len(blocks) && blocks[startFree] >= 0 {
				startFree++
			}
			// Find end of free space
			endFree := startFree
			for endFree < len(blocks) && blocks[endFree] < 0 {
				endFree++
			}
			// If file ID can be moved
			if idx > startFree && endFree-startFree >= nbrBlocks {
				for i := 0; i < nbrBlocks; i++ {
					blocks[startFree+i] = id
					blocks[idx+i] = -1
				}
				break
			}
			// Continue search for a big enough free space
			startFree = endFree
		}
		// printBlocks(blocks)
		// Continue with next file
		id--
	}
}

func checkSum(blocks []int) int {
	sum := 0
	for i := range blocks {
		if blocks[i] < 0 {
			continue
		}
		id := blocks[i]
		prod := id * i
		sum += prod
	}
	return sum
}

// Utility functions
func printBlock(b int) string {
	if b < 0 {
		return "."
	}
	return strconv.Itoa(b)
}
func printBlocks(blocks []int) {
	for _, b := range blocks {
		fmt.Printf("%v", printBlock(b))
	}
	fmt.Println()
}
