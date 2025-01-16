package keypad

import (
	"bytes"
)

type (
	Symbol   byte
	Sequence []byte
)

type Move struct {
	from Symbol
	to   Symbol
}

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
)

type Numpad struct {
	data     []byte            // To function as a reader
	current  Symbol            // Keep track of the current key
	layout   [][]byte          // The key layout -> To generate key2pos map
	keys     []Symbol          // The keys in this keypad -> To generate "moves" and move2seq map
	move2seq map[Move]Sequence // Given a move -> Get a sequence of symbols
	key2pos  map[Symbol]P
	// pos2key  map[P]Symbol
}

func NewNumpad(data []byte) *Numpad {
	np := &Numpad{data: data}
	np.current = 'A'
	np.keys = []Symbol{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'A'}
	np.layout = [][]byte{
		{0, 0, 0, 0, 0},
		{0, '7', '8', '9', 0},
		{0, '4', '5', '6', 0},
		{0, '1', '2', '3', 0},
		{0, 0, '0', 'A', 0},
		{0, 0, 0, 0, 0},
	}
	np.initKeyPosMaps()
	np.initMoves()

	return np
}
func NewDirpad(data Sequence) *Numpad {
	dp := &Numpad{data: data}
	dp.current = 'A'
	dp.keys = []Symbol{'^', 'A', '<', 'v', '>'}
	dp.layout = [][]byte{
		{0, 0, 0, 0, 0},
		{0, 0, '^', 'A', 0},
		{0, '<', 'v', '>', 0},
		{0, 0, 0, 0, 0},
	}
	dp.initKeyPosMaps()
	dp.initMoves()

	return dp
}
func (np *Numpad) Read(p []byte) (int, error) {
	i := 0
	for _, b := range np.data {
		if b == 10 {
			p[i] = b
			i++
		} else {
			seq := np.Press(Symbol(b))
			i += copy(p[i:], []byte(seq))
		}
	}
	return i, nil
}
func (np *Numpad) Press(key Symbol) Sequence {
	seq := np.MoveTo(key)
	seq = append(seq, 'A')
	return seq
}

func (np *Numpad) MoveTo(key Symbol) Sequence {
	seq := np.PeekTo(key)
	np.current = key
	return seq
}

func (np *Numpad) PeekTo(key Symbol) Sequence {
	seq := np.move2seq[Move{np.current, key}]
	return seq
}

func (np *Numpad) initMoves() {
	np.move2seq = map[Move]Sequence{}

	for _, from := range np.keys {
		for _, to := range np.keys {
			if from == to {
				continue
			}
			move := Move{from, to}
			np.move2seq[move] = np.getSeqFromPositions(np.key2pos[from], np.key2pos[to])
		}
	}

}

// BFS to find nearest path between 2 position
func (np *Numpad) getSeqFromPositions(start, end P) Sequence {
	type qItem struct {
		pos P
		seq Sequence
	}

	q := []qItem{{start, Sequence{}}}
	seen := map[P]bool{}

	for len(q) > 0 {
		// Queue
		curr := q[0]
		q = q[1:]
		p := curr.pos

		// Seen
		if seen[p] {
			continue
		}
		seen[p] = true

		// Exit condition
		if p == end {
			return curr.seq
		}

		// Neighbours
		for _, sym := range ">v<^" { // Using a string sequence instead of the map, to get correct order
			dir := Sym2dir[Symbol(sym)] // This var (map) is "global" within types...should it be somewhere else?
			n := p.Add(dir)
			if !np.validPos(n) {
				continue
			}
			if seen[n] {
				continue
			}
			seq := bytes.Clone(curr.seq)
			seq = append(seq, byte(sym))
			q = append(q, qItem{n, seq})
		}

	}
	return Sequence{}
}

func (np *Numpad) validPos(p P) bool {
	return np.layout[p.R][p.C] != 0
}

func (np *Numpad) initKeyPosMaps() {
	np.key2pos = map[Symbol]P{}
	// pos2key = map[P]Symbol{}

	for r := range np.layout {
		for c := range np.layout[0] {
			key := np.layout[r][c]
			pos := P{R: r, C: c}
			if key != 0 {
				np.key2pos[Symbol(key)] = pos
				// pos2key[pos] = Symbol(key)
			}
		}
	}
}

type P struct {
	R, C int
}

func (p P) Add(p2 P) P {
	return P{p.R + p2.R, p.C + p2.C}
}
