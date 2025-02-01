package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	data, _ := os.ReadFile("data.txt")
	data = bytes.TrimSpace(data)
	units := strings.Split(string(data), "\n\n")

	locks := [][]int{}
	keys := [][]int{}
	for _, unit := range units {

		lines := strings.Split(string(unit), "\n")
		// R := len(lines)
		// C := len(lines[0])
		keyLock := []string{}
		for _, line := range lines {
			keyLock = append(keyLock, line)
		}
		if keyLock[0][0] == '#' {
			// Lock
			locks = append(locks, decodeLock(keyLock))
		} else if keyLock[0][0] == '.' {
			keys = append(keys, decodeKey(keyLock))
			// Lock
		} else {
			log.Fatal("Error parsing lock and keys: ", keyLock)
		}
	}

	p1(locks, keys)
}

func decodeLock(lockString []string) []int {
	pins := len(lockString[0])
	lock := make([]int, pins)
	for _, line := range lockString[1:] {
		for pin := range pins {
			if line[pin] == '#' {
				lock[pin]++
			}
		}
	}
	return lock
}

func decodeKey(keyString []string) []int {
	pins := len(keyString[0])
	key := make([]int, pins)
	for i := len(keyString) - 2; i >= 0; i-- {
		for pin := range pins {
			if keyString[i][pin] == '#' {
				key[pin]++
			}
		}
	}
	return key
}

func p1(locks, keys [][]int) {
	nbr := 0
	for _, lock := range locks {
		for _, key := range keys {
			if fits(key, lock) {
				nbr++
			}
		}
	}
	fmt.Println("P1: (2854) ", nbr)
}

func fits(key, lock []int) bool {
	for i := range key {
		if key[i]+lock[i] > 5 {
			return false
		}
	}
	return true
}
