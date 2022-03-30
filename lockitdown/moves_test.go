package lockitdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTurn(t *testing.T) {

	testcases := []struct {
		direction TurnDirection
		expected  Pair
	}{
		{Left,
			Pair{0, -1}},
		{Left,
			Pair{-1, 0}},
		{Left,
			Pair{-1, 1}},
		{Left,
			Pair{0, 1}},

		// Roll it back!

		{Right,
			Pair{-1, 1}},
		{Right,
			Pair{-1, 0}},
		{Right,
			Pair{0, -1}},
	}

	direction := Pair{1, -1}
	for _, tc := range testcases {
		direction.Rotate(tc.direction)
		assert.EqualValues(t, tc.expected, direction, "Wrong turn!")
	}
}
