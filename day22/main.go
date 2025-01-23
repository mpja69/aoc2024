package main

import (
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"time"
)

func main() {
	timer(p1)()
	timer(p2)()
}
func p1() {
	sum := 0
	for secret := range readNumberIterator("data.txt") {
		for range 2000 {
			secret = evolve(secret)
		}
		sum += secret
	}

	fmt.Printf("P1: %d\n", sum)
}

// Just for fun: A Timing-decorater
func timer(fn func()) func() {
	return func() {
		t := time.Now()
		fn()
		fmt.Printf("%v\n", time.Since(t))
	}
}

// type Seq [4]int

func p2() {
	all := map[[4]int]int{}

	// Loop over the secrets in the file
	for secret := range readNumberIterator("data.txt") {
		seen := map[[4]int]bool{}
		buf := NewRingBuffer[int]()
		lastPrice := 0

		// Evolve 2000 times, and store all ok sequences/prices
		for range 2000 {
			price := secret % 10
			delta := price - lastPrice
			lastPrice = price
			buf.Write(delta)
			secret = evolve(secret)

			sequence, ok := buf.Sequence()
			if !ok {
				continue // Skip the first 3 rounds, until the ring buffer is filled up
			}
			if seen[sequence] {
				continue // Only store the first occurance
			}
			seen[sequence] = true
			if _, ok := all[sequence]; !ok {
				all[sequence] = 0
			}
			all[sequence] += price

		}
	}

	bestPrice := slices.Max(slices.Collect(maps.Values(all)))
	fmt.Printf("P2: %d\n", bestPrice)
}

type RingBuffer[T ~int] struct {
	buf     [4]T
	idx     int
	touches int
}

func NewRingBuffer[T ~int]() *RingBuffer[T] {
	rb := RingBuffer[T]{}
	return &rb
}
func (rb *RingBuffer[T]) Write(val T) {
	rb.buf[rb.idx] = val
	rb.idx = (rb.idx + 1) % len(rb.buf)
	rb.touches++
}
func (rb *RingBuffer[T]) Sequence() ([4]T, bool) {
	if rb.touches < 4 {
		return [4]T{}, false
	}
	res := [4]T{}
	// Return the values in the correct order, starting at idx
	for i := range 4 {
		res[i] = rb.buf[rb.idx]
		rb.idx = (rb.idx + 1) % len(rb.buf)
	}
	return res, true
}

// Iterator that reads a files of lines as numbers
func readNumberIterator(s string) func(func(int) bool) {
	f, err := os.Open(s)
	if err != nil {
		log.Fatal("readNumber(): ", err)
	}
	return func(yield func(val int) bool) {
		defer f.Close()
		val := 0
		for _, err := fmt.Fscanln(f, &val); err != io.EOF; _, err = fmt.Fscanln(f, &val) {
			if err != nil {
				log.Fatal("readNumber(): ", err)
				return
			}
			if !yield(val) {
				return
			}
		}
	}

}

// Almost twice (2X) as fast!!
func evolve(secret int) int {
	secret = ((secret << 6) ^ secret) & 16777215
	secret = ((secret >> 5) ^ secret) & 16777215
	secret = ((secret << 11) ^ secret) & 16777215
	return secret
}

// The naive first version
func evolve1(secret int) int {
	secret = ((secret * 64) ^ secret) % 16777216
	secret = ((secret / 32) ^ secret) % 16777216
	secret = ((secret * 2048) ^ secret) % 16777216
	return secret
}
