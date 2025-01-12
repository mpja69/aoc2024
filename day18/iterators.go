package main

import "fmt"

func iterator_test() {

	println("Fib + DP: Using the iterator together with range")
	// DP (table), and with push iterator w/o the import
	for f := range FibDP(10) {
		fmt.Print(f, "  ")
	}

	println()
	println("Fib + DP: Using the iterator as a 'normal' closure...over-complicated!")
	FibDP2(10)(func(curr, end int) bool {
		fmt.Print(curr, "  ")
		// if curr == end {
		// 	return false
		// }
		return true
	})

	println()
	println("Fib + DP: Using an iterator like a map closure")
	AllFib(10)(func(curr int) {
		fmt.Print(curr, "  ")
	})

	println()
	println("Fib+ BU: Same kind of iterator as FibDP, but a naive Bottom-Up Fib function")
	for f := range FibBU(10) {
		fmt.Print(f, "  ")
	}

	println()
	println("Even: A homebrewed pull iterator, with 2 callbacks")
	next, value := getEvenFunc(10)
	for next() {
		fmt.Print(value(), "  ")
	}

}

// NOTE: Can use the normal iterator, FibDB2, but just omitting the bool and the check,...in the
func AllFib(n int) func(func(int)) {
	return func(yield func(int)) {
		t := make([]int, n+2)
		t[0] = 0
		t[1] = 1

		for i := 1; i < n; i++ {
			t[i+1] += t[i]
			t[i+2] += t[i]
			yield(t[i+1])
		}
	}
}

func FibDP2(n int) func(func(int, int) bool) {
	return func(yield func(int, int) bool) {
		t := make([]int, n+2)
		t[0] = 0
		t[1] = 1

		for i := 1; i < n; i++ {
			t[i+1] += t[i]
			t[i+2] += t[i]
			if !yield(t[i+1], n) {
				return
			}
		}
	}
}

// A getFunc returuning 2 funcs: Next() and Valur()
// Compare this to an object with 2 methods.
func getEvenFunc(end int) (func() bool, func() int) {
	i := 0
	next := func() bool {
		i += 2
		if i < end {
			return true
		}
		return false
	}
	value := func() int {
		return i
	}

	return next, value
}

func FibBU(n int) func(func(int) bool) {
	return func(yield func(int) bool) {
		a := 0
		b := 1

		for i := 1; i < n; i++ {
			a, b = b, a+b
			if !yield(b) {
				return
			}
		}
	}
}

func FibDP(n int) func(func(int) bool) {
	return func(yield func(int) bool) {
		t := make([]int, n+2)
		t[0] = 0
		t[1] = 1

		for i := 1; i < n; i++ {
			t[i+1] += t[i]
			t[i+2] += t[i]
			if !yield(t[i+1]) {
				return
			}
		}
	}
}

func sieve(n int) []int {
	t := make([]bool, n)
	p := []int{}
	t[2] = true // not prime

	for i := 2; i < n; i++ {
		for j := i + i; j < n; j += i {
			t[j] = true
		}
	}
	for i, b := range t {
		if !b {
			p = append(p, i)
		}
	}
	return p
}
