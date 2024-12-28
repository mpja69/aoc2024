package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Pair struct {
	x, y int
}

type Config struct {
	a, b, p Pair
}

type Matrix [][]float64
type Vector []float64

func newMatrix(a, b Pair) Matrix {
	return Matrix{
		{float64(a.x), float64(b.x)},
		{float64(a.y), float64(b.y)},
	}
}
func newVector(a Pair) Vector {
	return Vector{float64(a.x), float64(a.y)}
}

func (m Matrix) det() float64 {
	return 1 / (m[0][0]*m[1][1] - m[0][1]*m[1][0])
}

func (m Matrix) scalar(n float64) Matrix {
	return Matrix{
		{m[0][0] * n, m[0][1] * n},
		{m[1][0] * n, m[1][1] * n},
	}
}
func (m Matrix) inv() Matrix {
	i := Matrix{
		{m[1][1], -m[0][1]},
		{-m[1][0], m[0][0]},
	}
	return i.scalar(m.det())
}
func (m Matrix) mul(v Vector) Vector {
	return []float64{
		m[0][0]*v[0] + m[0][1]*v[1],
		m[1][0]*v[0] + m[1][1]*v[1],
	}
}

// Hard coded "Gauss-Jordan elimination" for a 2x2 matrix
func gaussEliminationSolver(A Matrix, y Vector) (x Vector) {

	// First find the "echelon form", (zeros below the diagonal)
	s := A[1][0] / A[0][0]
	// R2 - s*R1 -> R2
	A[1][0] = 0
	A[1][1] -= s * A[0][1]
	y[1] -= s * y[0]

	// Continue with reducing the rows above the diagonal, ("reduced row echelon form")
	s = A[0][1] / A[1][1]
	// R1 - s*R2 -> R1
	A[0][0] -= s * A[1][0]
	A[0][1] = 0
	y[0] -= s * y[1]

	// Continue with "Back substitution" to find the answer
	x = make(Vector, 2)
	x[1] = y[1] / A[1][1]
	x[0] = y[0] / A[0][0]
	return x
}

func part1(data []Config) int {
	cost := 0
	for _, c := range data {
		A := newMatrix(c.a, c.b)
		x := A.inv().mul(newVector(c.p))

		// Round them...since they should be integers
		a := math.Round(x[0])
		b := math.Round(x[1])

		// Calulate the prize to verify the parameters
		y := A.mul(Vector{a, b})
		if y[0] == float64(c.p.x) && y[1] == float64(c.p.y) {
			cost += 3*int(a) + int(b)
		}
	}
	return cost
}
func part2(data []Config) int {
	cost := 0
	for _, c := range data {
		A := newMatrix(c.a, c.b)
		// for part2 add an extreme high number to the prize
		c.p.x += 10000000000000
		c.p.y += 10000000000000
		x := gaussEliminationSolver(A, newVector(c.p))

		// Round them...since the parameters should be integers
		a := math.Round(x[0])
		b := math.Round(x[1])

		// Need create the matrix again, since the solver changes it...
		A = newMatrix(c.a, c.b)
		//...calculate the prize, to verify correct parameters
		y := A.mul(Vector{a, b})
		if y[0] == float64(c.p.x) && y[1] == float64(c.p.y) {
			cost += 3*int(a) + int(b)
		}
	}
	return cost
}

func main() {
	data, err := os.ReadFile("s.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}
	list := []Config{}

	input := strings.Trim(string(data), " \n")
	groups := strings.Split(input, "\n\n")
	for _, group := range groups {
		lines := strings.Split(group, "\n")
		a := newButton(lines[0])
		b := newButton(lines[1])
		p := newPrize(lines[2])
		c := Config{a: a, b: b, p: p}
		list = append(list, c)
	}

	cost := part1(list)
	fmt.Printf("Part 1: (480, 35997) %d\n", cost)
	cost = part2(list)
	fmt.Printf("Part 2: (875318608908, 82510994362072) %d\n", cost)
}

func newButton(str string) Pair {
	parts := strings.Split(str, "+")
	x, _ := strings.CutSuffix(parts[1], ", Y")
	y := parts[2]
	pair := Pair{}
	pair.x, _ = strconv.Atoi(x)
	pair.y, _ = strconv.Atoi(y)
	return pair
}
func newPrize(str string) Pair {
	parts := strings.Split(str, "=")
	x, _ := strings.CutSuffix(parts[1], ", Y")
	y := parts[2]
	pair := Pair{}
	pair.x, _ = strconv.Atoi(x)
	pair.y, _ = strconv.Atoi(y)

	return pair
}
