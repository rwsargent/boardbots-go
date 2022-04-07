package lockitdown

import "sort"

var cache map[int][]Placement = make(map[int][]Placement)

type Placement struct {
	position  Pair
	direction Pair
}

type ByCorner []Placement

func (a ByCorner) Len() int      { return len(a) }
func (a ByCorner) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCorner) Less(i, j int) bool {
	return min(a[i].position) < min(a[j].position)
}

func edges(ringSize int) []Placement {
	if cached, found := cache[ringSize]; found {
		return cached
	}

	edges := make([]Placement, 0, (3*6)+(6*4*(ringSize-1)))

	// Top left
	cursor := Pair{
		0,
		-ringSize,
	}
	// idx := 0
	for side := 0; side < 6; side++ {
		dir := Cardinals[side]
		for hex := 0; hex < ringSize; hex++ {
			cursor.Plus(dir)
			for _, placeDirection := range Cardinals {
				position := cursor.Copy()
				position.Plus(placeDirection)
				if inBounds(ringSize, position) {
					edges = append(edges, Placement{
						position:  cursor.Copy(),
						direction: placeDirection,
					})
					// idx++
				}
			}
		}
	}
	sort.Sort(ByCorner(edges))
	cache[ringSize] = edges
	return edges
}

func min(p Pair) int {
	abQ := intAbs(p.Q)
	abR := intAbs(p.R)
	abS := intAbs(p.S())
	if abQ < abR {
		if abQ < abS {
			return abQ
		} else {
			return abS
		}
	} else {
		if abR < abS {
			return abR
		} else {
			return abS
		}
	}
}
