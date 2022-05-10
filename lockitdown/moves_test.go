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

func TestAdvance(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)
	state.Robots = map[Pair]*Robot{
		{2, 3}: {
			Position:      Pair{2, 3},
			Direction:     NW,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        0,
		},
	}
	move := NewMove(&AdvanceRobot{
		Robot: Pair{2, 3},
	}, 0)
	err := state.Move(move)
	assert.Nil(t, err)
	assert.Equal(t, 2, state.MovesThisTurn)

	err = state.Undo(move)

	assert.Nil(t, err)
	assert.Equal(t, 3, state.MovesThisTurn)
	bot, found := state.Robots[Pair{2, 3}]
	assert.True(t, found)
	assert.Equal(t, Pair{2, 3}, bot.Position)
}

func TestAdvanceBlocksLockdown(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)
	state.Robots = map[Pair]*Robot{
		{4, 0}: {
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0,
		},
		{4, -4}: {
			Position:      Pair{4, -4},
			Direction:     SE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{0, 4}: {
			Position:      Pair{0, 4},
			Direction:     NE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{2, 3}: {
			Position:      Pair{2, 3},
			Direction:     W,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        0,
		},
	}

	move := NewMove(&AdvanceRobot{
		Robot: Pair{2, 3},
	}, 0)
	err := state.Move(move)

	assert.Nil(t, err)
	assert.False(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.True(t, state.Robots[Pair{4, 0}].IsBeamEnabled)

	err = state.Undo(move)
	assert.Nil(t, err)
	assert.True(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.False(t, state.Robots[Pair{4, 0}].IsBeamEnabled)
}

func TestAdvanceRemovesBot(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)
	state.Robots = map[Pair]*Robot{
		{4, 0}: {
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0,
		},
		{4, -4}: {
			Position:      Pair{4, -4},
			Direction:     SE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{-4, 0}: {
			Position:      Pair{-4, 0},
			Direction:     E,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{-1, 5}: {
			Position:      Pair{-1, 5},
			Direction:     NE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
	}

	state.PlayerTurn = 1
	move := NewMove(&AdvanceRobot{
		Robot: Pair{-1, 5},
	}, 1)

	err := state.Move(move)
	assert.Nil(t, err)

	assert.Nil(t, state.Robots[Pair{4, 0}])
	assert.Equal(t, 3, state.Players[1].Points)
	assert.Len(t, state.Robots, 3)

	state.Undo(move)

	assert.NotNil(t, state.Robots[Pair{4, 0}])
	assert.Equal(t, 0, state.Players[1].Points)
	assert.Len(t, state.Robots, 4)
}

func TestTurnLockUnlock(t *testing.T) {
	state := NewGame(TwoPlayerGameDef)

	state.Robots = map[Pair]*Robot{
		{4, 0}: {
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0,
		},
		{4, -4}: {
			Position:      Pair{4, -4},
			Direction:     SE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{-4, 0}: {
			Position:      Pair{-4, 0},
			Direction:     E,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{0, -4}: {
			Position:      Pair{0, -4},
			Direction:     E,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        0,
		},
	}
	state.PlayerTurn = 1
	state.MovesThisTurn = 3

	move1 := NewMove(&TurnRobot{
		Robot:     Pair{-4, 0},
		Direction: Left,
	}, 1)

	err := state.Move(move1)
	assert.Nil(t, err)
	assert.False(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.False(t, state.Robots[Pair{0, -4}].IsLockedDown)

	move2 := NewMove(&TurnRobot{
		Robot:     Pair{4, -4},
		Direction: Right,
	}, 1)
	err = state.Move(move2)
	assert.Nil(t, err)
	assert.False(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.False(t, state.Robots[Pair{0, -4}].IsLockedDown)

	move3 := NewMove(&TurnRobot{
		Robot:     Pair{4, -4},
		Direction: Right,
	}, 1)
	err = state.Move(move3)
	assert.Nil(t, err)
	assert.False(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.True(t, state.Robots[Pair{0, -4}].IsLockedDown)

	// REVERSE!

	err = state.Undo(move3)
	assert.Nil(t, err)
	assert.False(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.False(t, state.Robots[Pair{0, -4}].IsLockedDown)

	err = state.Undo(move2)
	assert.Nil(t, err)
	assert.False(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.False(t, state.Robots[Pair{0, -4}].IsLockedDown)

	err = state.Undo(move1)
	assert.Nil(t, err)
	assert.True(t, state.Robots[Pair{4, 0}].IsLockedDown)
	assert.False(t, state.Robots[Pair{0, -4}].IsLockedDown)
}

func TestRemovedToLock(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)
	game.Robots = map[Pair]*Robot{
		{0, 4}: {
			Position:      Pair{0, 4},
			Direction:     NW,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0,
		},
		{-4, 4}: {
			Position:      Pair{-4, 4},
			Direction:     E,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{4, 0}: {
			Position:      Pair{4, 0},
			Direction:     SW,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{0, 0}: {
			Position:      Pair{0, 0},
			Direction:     E,
			IsBeamEnabled: false,
			IsLockedDown:  true,
			Player:        0,
		},
	}
}

func TestPlacedRobots(t *testing.T) {

	game := NewGame(TwoPlayerGameDef)

	m1 := NewMove(&PlaceRobot{
		Robot:     Pair{0, 5},
		Direction: NW,
	}, 0)

	err := game.Move(m1)
	assert.Nil(t, err)
	assert.Equal(t, 1, game.Players[0].PlacedRobots)

	m2 := NewMove(&PlaceRobot{
		Robot:     Pair{0, -5},
		Direction: SE,
	}, 1)

	err = game.Move(m2)
	assert.Nil(t, err)
	assert.Equal(t, 1, game.Players[1].PlacedRobots)

	m3 := NewMove(&PlaceRobot{
		Robot:     Pair{-5, 5},
		Direction: SW,
	}, 0)

	err = game.Move(m3)
	assert.Nil(t, err)
	assert.Equal(t, 2, game.Players[0].PlacedRobots)

	m4 := NewMove(&PlaceRobot{
		Robot:     Pair{5, -5},
		Direction: NE,
	}, 1)

	err = game.Move(m4)
	assert.Nil(t, err)
	assert.Equal(t, 2, game.Players[1].PlacedRobots)

	game.Undo(m4)
	game.Undo(m3)
	game.Undo(m2)
	game.Undo(m1)

	assert.Equal(t, 0, game.Players[0].PlacedRobots)
	assert.Equal(t, 0, game.Players[1].PlacedRobots)
}
