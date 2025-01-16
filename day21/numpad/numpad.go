package numpad

type Numpad struct {
	current  string
	layout   [][]byte
	keys     []string
	move2seq map[Move]string
	key2pos  map[string]P
	// pos2key  map[P]string
}

func NewNumpad() *Numpad {
	np := &Numpad{}
	np.current = "A"
	np.keys = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "A"}
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

func (np *Numpad) MoveTo(nbr string) string {
	seq := np.move2seq[Move{np.current, nbr}]
	np.current = nbr
	return seq
}

type P struct {
	r, c int
}
type Move struct {
	from string
	to   string
}

func (np *Numpad) initMoves() {
	np.move2seq = map[Move]string{}

	for _, from := range np.keys {
		for _, to := range np.keys {
			if from == to {
				continue
			}
			move := Move{from, to}
			np.move2seq[move] = np.getSeqFromPosToPos(np.key2pos[from], np.key2pos[to])
		}
	}

}

func (np *Numpad) getSeqFromPosToPos(start, end P) string {
	type Item struct {
		pos P
		seq string
	}

	q := []Item{{start, ""}}
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
		for d, s := range dir2sym {
			r, c := p.r+d.r, p.c+d.c
			if np.layout[r][c] == 0 {
				continue
			}
			n := P{r, c}
			if seen[n] {
				continue
			}
			seq := curr.seq + s
			q = append(q, Item{n, seq})
		}

	}
	return ""
}

func (np *Numpad) initKeyPosMaps() {
	np.key2pos = map[string]P{}
	// pos2key = map[P]string{}

	for r := range np.layout {
		for c := range np.layout[0] {
			key := np.layout[r][c]
			pos := P{r, c}
			if key != 0 {
				np.key2pos[string(key)] = pos
				// pos2key[pos] = string(key)
			}
		}
	}
}
