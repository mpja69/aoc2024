package numpad

import (
	"aoc2024/day21/types"
)

type (
	P        = types.P
	Symbol   = types.Symbol
	Sequence = types.Sequence
)

type AlphaNumSymbol string
type AlphaNumSequence string

type Numpad struct {
	current  AlphaNumSymbol
	layout   [][]byte
	keys     []AlphaNumSymbol
	move2seq map[Move]Sequence // A move -> A sequence of direction-symbols
	key2pos  map[AlphaNumSymbol]P
	// pos2key  map[P]AlphaNumSymbol
}

func New() *Numpad {
	np := &Numpad{}
	np.current = "A"
	np.keys = []AlphaNumSymbol{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "A"}
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

// HACK: What is the best way to ensure that the given string is just 1 alhpanumeric char? (Is it even necessary?!)
func (np *Numpad) MoveTo(nbr AlphaNumSymbol) Sequence {
	if len(nbr) != 1 {
		panic("AlphaNumSymbol longer than 1: " + nbr)
	}
	if !((nbr[0] >= '0' && nbr[0] <= '9') || nbr[0] == 'A') {
		panic("Not an AlphaNumSymbol: " + nbr)
	}
	seq := np.move2seq[Move{np.current, nbr}]
	np.current = nbr
	return seq
}

func (np *Numpad) PeekTo(nbr AlphaNumSymbol) Sequence {
	if len(nbr) != 1 {
		panic("AlphaNumSymbol longer than 1: " + nbr)
	}
	if !((nbr[0] >= '0' && nbr[0] <= '9') || nbr[0] == 'A') {
		panic("Not an AlphaNumSymbol: " + nbr)
	}
	seq := np.move2seq[Move{np.current, nbr}]
	return seq
}

type Move struct {
	from AlphaNumSymbol
	to   AlphaNumSymbol
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

	q := []qItem{{start, ""}}
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
			dir := types.Sym2dir[Symbol(sym)] // This var (map) is "global" within types...should it be somewhere else?
			n := p.Add(dir)
			if !np.validPos(n) {
				continue
			}
			if seen[n] {
				continue
			}
			seq := curr.seq + Sequence(sym)
			q = append(q, qItem{n, seq})
		}

	}
	return ""
}

func (np *Numpad) validPos(p P) bool {
	return np.layout[p.R][p.C] != 0
}

func (np *Numpad) initKeyPosMaps() {
	np.key2pos = map[AlphaNumSymbol]P{}
	// pos2key = map[P]AlphaNumSymbol{}

	for r := range np.layout {
		for c := range np.layout[0] {
			key := np.layout[r][c]
			pos := P{R: r, C: c}
			if key != 0 {
				np.key2pos[AlphaNumSymbol(key)] = pos
				// pos2key[pos] = AlphaNumSymbol(key)
			}
		}
	}
}
