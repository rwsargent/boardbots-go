package lockitdown

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorThirdPly(t *testing.T) {

	gameState := NewGame(TwoPlayerGameDef)

	gameState.Robots = map[Pair]*Robot{
		{-5, 0}: {
			Position:      Pair{-5, 0},
			Direction:     E,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        0,
		},
		{5, 0}: {
			Position:      Pair{5, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        0,
		},
		{-5, 5}: {
			Position:      Pair{-5, 5},
			Direction:     NE,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        1,
		},
		{5, -5}: {
			Position:      Pair{5, -5},
			Direction:     SW,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        1,
		},
	}
	gameState.PlayerTurn = PlayerPosition(0)

	it := NewMoveIterator(gameState)

	for i := 0; i < 6; i++ {
		assert.True(t, it.Next())
		m := it.Get()
		assert.NotNil(t, m)
		assert.NotNilf(t, m.Mover, "failed on %d iteration", i)
	}

	assert.False(t, it.Next())
}

func TestNewGameIterator(t *testing.T) {

	game := NewGame(TwoPlayerGameDef)

	it := NewMoveIterator(game)

	for i := 0; i < (6*3)+(6*4*4); i++ {
		assert.True(t, it.Next())
		m := it.Get()
		assert.NotNil(t, m)
		assert.NotNilf(t, m.Mover, "failed on %d iteration", i)
	}
	assert.False(t, it.Next())
}

func TestFullMoveIterator(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)
	it := NewMoveIterator(game)

	game.Robots = map[Pair]*Robot{
		{-4, 0}: {
			Position:      Pair{-4, 0},
			Direction:     E,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        0,
		},
		{4, 0}: {
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: false,
			IsLockedDown:  false,
			Player:        1,
		},
	}

	for i := 0; i < (6*3)+(6*4*4)+3; i++ {
		assert.Truef(t, it.Next(), "failed on %d iteration", i)
		m := it.Get()
		assert.NotNil(t, m)
		assert.NotNilf(t, m.Mover, "failed on %d iteration", i)
	}

	assert.False(t, it.Next())
}

func TestIteratorFromState(t *testing.T) {

	gamjson := `{
		"gameDef": {
		  "board": {
			"HexaBoard": {
			  "arenaRadius": 4
			}
		  },
		  "maxRobotsInStaging": 2,
		  "winCondition": "Elimination",
		  "movesPerTurn": 3,
		  "robotsPerPlayer": 6,
		  "numPlayers": 2
		},
		"players": [
		  {
			"points": 0,
			"placedRobots": 2
		  },
		  {
			"points": 0,
			"placedRobots": 2
		  }
		],
		"robots": [
		  [
			{
			  "q": -3,
			  "r": -2
			},
			{
			  "player": 1,
			  "dir": {
				"q": 1,
				"r": -1
			  },
			  "isLocked": false,
			  "isBeamEnabled": false
			}
		  ],
		  [
			{
			  "q": 2,
			  "r": -5
			},
			{
			  "player": 2,
			  "dir": {
				"q": 1,
				"r": 0
			  },
			  "isLocked": false,
			  "isBeamEnabled": false
			}
		  ],
		  [
			{
			  "q": 5,
			  "r": -5
			},
			{
			  "player": 2,
			  "dir": {
				"q": -1,
				"r": 0
			  },
			  "isLocked": false,
			  "isBeamEnabled": false
			}
		  ],
		  [
			{
			  "q": 0,
			  "r": -4
			},
			{
			  "player": 1,
			  "dir": {
				"q": 1,
				"r": 0
			  },
			  "isLocked": false,
			  "isBeamEnabled": true
			}
		  ]
		],
		"playerTurn": 2,
		"status": "OnGoing",
		"movesThisTurn": 2,
		"requiresTieBreak": false
	  }`

	var tGame TransportState

	err := json.Unmarshal([]byte(gamjson), &tGame)
	assert.Nil(t, err)

	state := StateFromTransport(&tGame)

	it := NewMoveIterator(state)

	i := 0
	for it.Next() {
		m := it.Get()
		assert.NotNil(t, m)
		assert.NotNilf(t, m.Mover, "iter: %d\n", i)
		i++
	}
}
