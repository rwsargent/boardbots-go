package lockitdown

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	assert.Len(t, edges(1), 18)
	assert.Len(t, edges(2), 42)
}

func TestCardinalDirections(t *testing.T) {
	edges := edges(1)
	directions := make(map[Pair][]Pair)
	for _, edge := range edges {
		if dir, found := directions[edge.position]; found {
			directions[edge.position] = append(dir, edge.direction)
		} else {
			dirs := make([]Pair, 0, 3)
			directions[edge.position] = append(dirs, edge.direction)
		}
	}
	for pos, dirs := range directions {
		assert.Lenf(t, dirs, 3, "%+v has %d directions", pos, len(dirs))
	}
}

func TestSort(t *testing.T) {
	edges := edges(2)

	display := make([]string, len(edges))
	for i, edge := range edges {
		display[i] = fmt.Sprintf("{%d,%d,%d}->%s", edge.position.Q, edge.position.R, edge.position.S(), edge.direction.String())
	}
	fmt.Printf("%v\n", display)
}

func TestDirection(t *testing.T) {
	edges := edges(3)

	for _, edge := range edges {
		next := edge.position.Copy()
		next.Plus(edge.direction)

		assert.LessOrEqual(t, next.Dist(), 3, "%s with direciton %s", edge.position.String(), edge.direction.String())
	}
}
