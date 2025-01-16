package types

type P struct {
	R, C int
}

func (p P) Add(p2 P) P {
	return P{p.R + p2.R, p.C + p2.C}
}

// These types exist to clarify the interaction with the pads
type (
	Sequence string // Sequence of Symbols. För att kunna skilja på symbol och sekvens...och ha metoder?
	Symbol   string // Symbol of Directions. Bara 1 tecken i en sequence
	// Ska man ha ett interfce istället?
)

// These variables (maps) are "global" within types...should it be somewhere else? Used by numpad, dirpad,...
var (
	Dir2sym = map[P]Symbol{
		{0, 1}:  ">",
		{1, 0}:  "v",
		{0, -1}: "<",
		{-1, 0}: "^",
	}
	Sym2dir = map[Symbol]P{
		">": {0, 1},
		"v": {1, 0},
		"<": {0, -1},
		"^": {-1, 0},
	}
)
