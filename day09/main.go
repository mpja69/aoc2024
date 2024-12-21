package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("s.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	input := strings.Trim(string(data), " \n")
	lines := strings.Split(input, "\n")

	fmt.Printf("Part 1: (1928, 6201130364722): %d\n", part1(lines[0]))
	fmt.Printf("Part 2: (2858, 6221662795602): %d\n", part2(lines[0]))
}

func part1(diskmap string) int {
	blocks, _ := inflateBlocks(diskmap)
	partitionSingleBlocks(blocks)
	return checkSum(blocks)
}

func part2(diskmap string) int {
	blocks, info := inflateBlocks(diskmap)
	partitionFiles(blocks, info)
	return checkSum(blocks)
}

type fileInfo struct {
	size int
	idx  int
}

func inflateBlocks(diskmap string) ([]int, []fileInfo) {
	info := make([]fileInfo, len(diskmap)/2+1)
	blocks := []int{}
	isFile := true
	id := 0
	for _, digit := range diskmap {
		nbrBlocks := int(digit) - '0'
		if isFile {
			info[id].idx = len(blocks)
			for block := 0; block < nbrBlocks; block++ {
				blocks = append(blocks, id)
			}
			info[id].size = nbrBlocks
			id++
		} else {
			for block := 0; block < nbrBlocks; block++ {
				blocks = append(blocks, -1)
			}
		}
		isFile = !isFile
	}
	return blocks, info
}

func partitionSingleBlocks(blocks []int) {
	src := len(blocks) - 1
	dst := 0

	for {
		// Move forward to next free block
		for blocks[dst] >= 0 {
			dst++
		}
		//Move beckward to next block (to move)
		for blocks[src] < 0 {
			src--
		}

		// printBlocks(b)
		// If now more blocks to move
		if dst >= src {
			return
		}

		// Move block and mark it as "free"
		blocks[dst] = blocks[src]
		blocks[src] = -1
	}
}

func partitionFiles(blocks []int, info []fileInfo) {
	id := len(info) - 1
	for id >= 0 {
		// Always look for free space by starting from the beginning
		dstIdx := 0

		// Identify file to move
		srcSize := info[id].size
		srcIdx := info[id].idx

		// Search for a slot big enough
		for dstIdx < len(blocks) {
			// Find the start of free space
			for dstIdx < len(blocks) && blocks[dstIdx] >= 0 {
				dstIdx++
			}
			// Find end of free space
			dstSize := 0
			for dstEnd := dstIdx; dstEnd < len(blocks) && blocks[dstEnd] < 0; dstEnd++ {
				dstSize++
			}

			if srcIdx > dstIdx && dstSize >= srcSize {
				// Found free space! Move the file (all the blocks)
				for i := 0; i < srcSize; i++ {
					blocks[dstIdx+i] = id
					blocks[srcIdx+i] = -1
				}
				// Do the next file
				break
			} else {
				// Free space to small! Jump forward and continue the search
				dstIdx += dstSize
			}
		}
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
