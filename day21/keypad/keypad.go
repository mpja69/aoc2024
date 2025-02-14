package keypad

import (
	"bytes"
	"io"
)

// INFO: --------- Testing stuff to learn Go --------------
//			Testing to use a separate package/module
//			Testing Reader and Writer interfaces
//			Testing type definition and type aliases
//			Testing a function from StackOverflow: cartesianProduct

// NOTE: Should maybe have used type aliases here instead?
// type (
// 	Symbol   = byte
// 	Sequence = []byte
// )

type (
	Symbol   byte
	Sequence []byte
)

type Move struct {
	From Symbol
	To   Symbol
}

type Layout int

const (
	NumberLayout Layout = iota
	DirectionLayout
)

var (
	Dir2sym = map[P]Symbol{
		{0, 1}:  '>',
		{1, 0}:  'v',
		{0, -1}: '<',
		{-1, 0}: '^',
	}
	Sym2dir = map[Symbol]P{
		'>': {0, 1},
		'v': {1, 0},
		'<': {0, -1},
		'^': {-1, 0},
	}
	Dir2dir = map[string]string{
		"A": "A",
		"^": "<A>A",
		">": "vA^A",
		"v": "v<A>^A",
		"<": "v<<A>>^A",
	}

	numLayout = [][]byte{
		{0, 0, 0, 0, 0},
		{0, '7', '8', '9', 0},
		{0, '4', '5', '6', 0},
		{0, '1', '2', '3', 0},
		{0, 0, '0', 'A', 0},
		{0, 0, 0, 0, 0},
	}
	dirLayout = [][]byte{
		{0, 0, 0, 0, 0},
		{0, 0, '^', 'A', 0},
		{0, '<', 'v', '>', 0},
		{0, 0, 0, 0, 0},
	}
)

type Keypad struct {
	current  Symbol              // Keep track of the current key
	layout   [][]byte            // The key layout -> To generate key2pos map
	keys     []Symbol            // The keys in this keypad -> To generate "moves" and move2seq map
	Move2seq map[Move][]Sequence // Given a move -> Get a sequence of symbols
	key2pos  map[Symbol]P
}

// ==================== Keypad ======================
func NewKeypad(layout Layout) *Keypad {
	kp := &Keypad{}
	if layout == NumberLayout {
		kp.layout = numLayout
	} else {
		kp.layout = dirLayout
	}
	kp.current = 'A'
	kp.initKeyPosMaps()
	kp.initMoves()
	return kp
}

// NOTE: Instead of taking string, (or a []byte), a io.Reader could been used! (And...Keypad doesn't have to be a Reader!)
func (kp *Keypad) GetPossibleSequences(input string) []string {
	combinations := [][]Sequence{}
	for _, b := range input {
		combo := kp.MoveTo(Symbol(b))
		combinations = append(combinations, combo)
	}

	matrix := cartesianProduct(combinations)
	inputs := []string{}
	for _, r := range matrix {
		input := ""
		for _, c := range r {
			input = input + string(c)
		}
		inputs = append(inputs, input)
	}
	return inputs
}

// INFO: If I had used type aliases for Symbol, it might have been used here a bit easier?
func (kp *Keypad) Move(a, b byte) []Sequence {
	return kp.Move2seq[Move{Symbol(a), Symbol(b)}]
}
func (kp *Keypad) MoveTo(key Symbol) []Sequence {
	seq := kp.Move2seq[Move{kp.current, key}]
	kp.current = key
	return seq
}

func (kp *Keypad) initKeyPosMaps() {
	kp.key2pos = map[Symbol]P{}
	kp.keys = []Symbol{}

	for r := range kp.layout {
		for c := range kp.layout[0] {
			key := kp.layout[r][c]
			pos := P{R: r, C: c}
			if key != 0 {
				kp.key2pos[Symbol(key)] = pos
				kp.keys = append(kp.keys, Symbol(key))
			}
		}
	}
}
func (kp *Keypad) initMoves() {
	kp.Move2seq = map[Move][]Sequence{}

	for _, from := range kp.keys {
		for _, to := range kp.keys {
			move := Move{from, to}
			if from == to {
				kp.Move2seq[move] = []Sequence{[]byte{'A'}}
			} else {
				kp.Move2seq[move] = kp.getSeqFromPositions(kp.key2pos[from], kp.key2pos[to])
			}
		}
	}
}

// BFS to find nearest path between 2 position, (works even if it's the same position)
func (kp *Keypad) getSeqFromPositions(start, end P) []Sequence {
	type qItem struct {
		pos P
		seq Sequence
	}

	q := []qItem{{start, Sequence{}}}
	res := []Sequence{}
	optLen := 100

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]
		p := curr.pos

		// Target condition
		if p == end {
			if len(curr.seq) > optLen {
				return res
			}
			optLen = len(curr.seq)
			curr.seq = append(curr.seq, 'A')
			res = append(res, curr.seq)
		}

		// Neighbours
		for _, sym := range ">v<^" { // Using a string sequence instead of the map, to get correct order
			dir := Sym2dir[Symbol(sym)] // This var (map) is "global" within types...should it be somewhere else?
			n := p.Add(dir)
			if kp.layout[n.R][n.C] == 0 {
				continue
			}
			seq := bytes.Clone(curr.seq)
			seq = append(seq, byte(sym))
			q = append(q, qItem{n, seq})
		}

	}
	return []Sequence{}
}

// ----------------- Util ---------------
type P struct {
	R, C int
}

func (p P) Add(p2 P) P {
	return P{p.R + p2.R, p.C + p2.C}
}

// TODO: Rewrite this to my own. (From https://stackoverflow.com/questions/29002724/implement-ruby-style-cartesian-product-in-go)
// cartesianProduct returns the cartesian product
// of a given matrix
func cartesianProduct[T any](matrix [][]T) [][]T {
	// nextIndex sets ix to the lexicographically next value,
	// such that for each i>0, 0 <= ix[i] < lens(i).
	nextIndex := func(ix []int, lens func(i int) int) {
		for j := len(ix) - 1; j >= 0; j-- {
			ix[j]++

			if j == 0 || ix[j] < lens(j) {
				return
			}

			ix[j] = 0
		}
	}

	lens := func(i int) int { return len(matrix[i]) }

	results := make([][]T, 0, len(matrix))
	for indexes := make([]int, len(matrix)); indexes[0] < lens(0); nextIndex(indexes, lens) {
		var temp []T

		for j, k := range indexes {
			temp = append(temp, matrix[j][k])
		}

		results = append(results, temp)
	}

	return results
}

// From stackoverflow
//
//	func Zip[T, U any](t []T, u []U) iter.Seq2[T, U] {
//	    return func(yield func(T, U) bool) {
//	        // Starting with Go 1.22,
//	        // we can also range over integer:
//	        for i := range min(len(t), len(u)) {
//	            if !yield(t[i], u[i]) {
//	                return
//	            }
//	        }
//	    }
//	}
//

// Example of using a reader
func tt() {
	t(bytes.NewBufferString("Hello"))
}
func t(r io.Reader) {
	var a []byte
	r.Read(a)
}
