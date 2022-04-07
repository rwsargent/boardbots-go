package lockitdown

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TwoPlayerGameDef = GameDef{
	Players: 2,
	Board: Board{
		HexaBoard: BoardType{4},
	},
	RobotsPerPlayer: 6,
	WinCondition:    "Elimination",
	MovesPerTurn:    3,
}

func TestNewGame(t *testing.T) {
	game := NewGame(GameDef{
		Players: 2,
	})
	if game.PlayerTurn != 0 {
		t.Errorf("Wrong player turn")
	}
	if len(game.Players) != 2 {
		t.Errorf(("Wrong number of players"))
	}
	if len(game.Robots) != 0 {
		t.Error("Improperly initialized robots")
	}
}
func TestMoves(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)

	tests := []struct {
		move   Mover
		player PlayerPosition
		err    error
	}{
		{&PlaceRobot{
			Hex:       Pair{0, 5},
			Direction: Pair{0, -1},
		}, 0, nil},
		{&PlaceRobot{
			Hex:       Pair{5, 0},
			Direction: Pair{-1, 0},
		}, 1, nil},
		{&AdvanceRobot{
			Robot: Pair{0, 5},
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Left,
		}, 0, nil},
		{&PlaceRobot{
			Hex:       Pair{-5, 4},
			Direction: Pair{1, 0},
		}, 1, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&AdvanceRobot{
			Robot: Pair{-5, 4},
		}, 1, nil},
		{&AdvanceRobot{
			Robot: Pair{5, 0},
		}, 1, nil},
		{&TurnRobot{
			Robot:     Pair{4, 0},
			Direction: Left,
		}, 1, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&PlaceRobot{
			Hex:       Pair{0, -5},
			Direction: Pair{0, 1},
		}, 1, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{&TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{
			move: &AdvanceRobot{
				Robot: Pair{0, -5},
			},
			player: 1,
			err:    nil,
		},
	}

	for _, tt := range tests {
		m := NewMove(tt.move, tt.player)
		err := game.Move(m)
		if tt.err == nil && err != nil {
			t.Errorf("expected no error, got %s", err.Error())
		}
		if tt.err != nil && err == nil {
			t.Errorf("expected error %s, none recieved", tt.err.Error())
		}
		if tt.err != nil && err != nil && tt.err.Error() != err.Error() {
			t.Errorf("expected error %s, recieved %s", tt.err.Error(), err.Error())
		}
	}

	assert.Equal(t, 1, game.Players[0].PlacedRobots, "wrong number of player 1 robots")
	assert.Equal(t, 3, game.Players[1].PlacedRobots, "wrong number of player 2 robots")
	assert.Equal(t, 3, game.Players[1].Points, "wrong number of player 2 points")
}

func TestGameOver(t *testing.T) {
	gameState := GameState{
		GameDef: TwoPlayerGameDef,
		Players: []*Player{
			{
				Points:       0,
				PlacedRobots: 6,
			},
			{
				Points:       9,
				PlacedRobots: 3,
			},
		},
		Robots: map[Pair]*Robot{
			{-4, 4}: {
				Position:      Pair{-4, 4},
				Direction:     NE,
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
			{0, -4}: {
				Position:      Pair{0, -4},
				Direction:     SE,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        1,
			},
			{4, -4}: {
				Position:      Pair{4, -4},
				Direction:     SW,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        0,
			},
			{5, -5}: {
				Position:      Pair{5, -5},
				Direction:     SW,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        0,
			},
			{0, 5}: {
				Position:      Pair{0, 5},
				Direction:     NW,
				IsBeamEnabled: true,
				IsLockedDown:  false,
				Player:        0,
			},
		},
		PlayerTurn:       0,
		MovesThisTurn:    3,
		RequiresTieBreak: false,
		Winner:           -1,
	}

	testcases := []struct {
		move   Mover
		player PlayerPosition
		result error
	}{
		{
			&AdvanceRobot{
				Robot: Pair{4, -4},
			},
			0,
			nil,
		},
		{
			&AdvanceRobot{
				Robot: Pair{5, -5},
			}, 0,
			nil,
		},
		{
			&AdvanceRobot{
				Robot: Pair{0, 5},
			},
			0,
			nil,
		},
	}

	for _, tc := range testcases {
		m := NewMove(tc.move, tc.player)
		err := gameState.Move(m)
		assert.Nilf(t, err, "unexpected err for move %v", tc.move)
		assert.Equal(t, -1, gameState.Winner)
		if err != nil {
			fmt.Println(gameState.ToJson())
		}
	}

	err := gameState.Move(NewMove(&TurnRobot{
		Robot:     Pair{-4, 4},
		Direction: Right,
	}, 1))

	fmt.Println(gameState.ToJson())

	assert.EqualError(t, err, "winner is 2")
	assert.Equal(t, 1, gameState.Winner)
}

func TestPairCopy(t *testing.T) {
	p1 := Pair{
		Q: 32,
		R: 12,
	}

	p2 := p1.Copy()

	p1.Q = 2
	p1.R = 4

	p2.Q = 5
	p2.R = 7

	assert.Equal(t, p1, Pair{2, 4})
	assert.Equal(t, p2, Pair{5, 7})
}

func TestInBounds(t *testing.T) {
	testCases := []struct {
		p        Pair
		s        int
		inBounds bool
	}{
		{
			p:        Pair{-1, 1},
			s:        1,
			inBounds: true,
		},
	}

	for _, tc := range testCases {
		assert.Truef(t, inBounds(tc.s, tc.p), "%+v is not InBounds", tc.p)
	}
}

func TestPossibleMoves(t *testing.T) {

	game := NewGame(TwoPlayerGameDef)
	initMoves := []*GameMove{
		NewMove(&PlaceRobot{
			Hex:       Pair{-5, 0},
			Direction: E,
		}, 0),
		NewMove(&PlaceRobot{
			Hex:       Pair{5, 0},
			Direction: W,
		}, 1),
		NewMove(&PlaceRobot{
			Hex:       Pair{0, -5},
			Direction: SE,
		}, 0),
		NewMove(&PlaceRobot{
			Hex:       Pair{0, 5},
			Direction: NW,
		}, 1),
		NewMove(&AdvanceRobot{
			Robot: Pair{-5, 0},
		}, 0),
		NewMove(&TurnRobot{
			Robot:     Pair{-4, 0},
			Direction: Left,
		}, 0),
		NewMove(&TurnRobot{
			Robot:     Pair{-4, 0},
			Direction: Left,
		}, 0),
	}

	for _, initmove := range initMoves {
		err := game.Move(initmove)
		assert.Nil(t, err)
	}

	possibleMoves := game.PossibleMoves([]*GameMove{})

	for _, possibleMove := range possibleMoves {
		player := game.PlayerTurn
		assert.Equal(t, PlayerPosition(1), player)
		err := game.Move(possibleMove)
		assert.Nil(t, err)
		game.Undo(possibleMove)
	}

	game.Undo(initMoves[6])
	game.Undo(initMoves[5])
	game.Undo(initMoves[4])

	assert.Equal(t, PlayerPosition(0), game.PlayerTurn)

	game.Undo(initMoves[3])

	assert.Equal(t, PlayerPosition(1), game.PlayerTurn)

	nextMoves := game.PossibleMoves([]*GameMove{})
	next := nextMoves[rand.Intn(len(nextMoves))]

	err := game.Move(next)
	assert.Nil(t, err)
	game.Undo(next)
}

func TestFakeMinimaxStressTest(t *testing.T) {

	game := NewGame(TwoPlayerGameDef)

	var recur func(*GameState, int)
	recur = func(game *GameState, depth int) {
		if depth == 0 {
			return
		}
		moves := game.PossibleMoves([]*GameMove{})
		for _, move := range moves {
			tState := ConvertToTransport(game)
			err := game.Move(move)
			assert.Nil(t, err)

			recur(game, depth-1)

			err = game.Undo(move)
			undoneTState := ConvertToTransport(game)
			assert.Nil(t, err)

			assert.Equal(t, tState.GameDef, undoneTState.GameDef)
			assert.ElementsMatch(t, tState.Players, undoneTState.Players)
			assert.Equal(t, len(tState.Robots), len(undoneTState.Robots))
			assert.Equal(t, tState.PlayerTurn, undoneTState.PlayerTurn)
		}
	}

	recur(game, 3)
}

func TestPossibleMovesFromState(t *testing.T) {
	jsonState := `{
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
			"placedRobots": 1
		  }
		],
		"robots": [
		  [
			{
			  "q": -1,
			  "r": -4
			},
			{
			  "player": 1,
			  "dir": {
				"q": 0,
				"r": 1
			  },
			  "isLocked": false,
			  "isBeamEnabled": false
			}
		  ],
		  [
			{
			  "q": -4,
			  "r": -1
			},
			{
			  "player": 2,
			  "dir": {
				"q": 0,
				"r": 1
			  },
			  "isLocked": false,
			  "isBeamEnabled": false
			}
		  ],
		  [
			{
			  "q": -1,
			  "r": 5
			},
			{
			  "player": 1,
			  "dir": {
				"q": 1,
				"r": 0
			  },
			  "isLocked": false,
			  "isBeamEnabled": false
			}
		  ]
		],
		"playerTurn": 2,
		"status": "OnGoing",
		"movesThisTurn": 0,
		"requiresTieBreak": false
	  }`

	var tState TransportState
	json.Unmarshal([]byte(jsonState), &tState)
	state := StateFromTransport(&tState)

	moves := state.PossibleMoves([]*GameMove{})
	fmt.Printf("%v\n", moves)
}

func BenchmarkPossibleMoves(b *testing.B) {
	b.StopTimer()
	game := NewGame(TwoPlayerGameDef)

	game.Robots = map[Pair]*Robot{
		{0, -4}: {
			Position:      Pair{0, -4},
			Direction:     SE,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        0,
		},
		{4, -4}: {
			Position:      Pair{4, -4},
			Direction:     SW,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        0,
		},
		{4, 0}: {
			Position:      Pair{4, 0},
			Direction:     W,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        0,
		},
		{0, 4}: {
			Position:      Pair{0, 4},
			Direction:     NW,
			IsBeamEnabled: true,
			IsLockedDown:  false,
			Player:        1,
		},
		{-4, 4}: {
			Position:      Pair{-4, 4},
			Direction:     NE,
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
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		game.PossibleMoves([]*GameMove{})
	}
}

func TestNilGameMoveFromState(t *testing.T) {
	jsonState := `{
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
			"placedRobots": 1
		  },
		  {
			"points": 0,
			"placedRobots": 1
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
				"r": 0
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
		  ]
		],
		"playerTurn": 1,
		"status": "OnGoing",
		"movesThisTurn": 0,
		"requiresTieBreak": false
	  }`

	var tState TransportState
	json.Unmarshal([]byte(jsonState), &tState)
	state := StateFromTransport(&tState)

	root := MinimaxNode{
		GameState:    state,
		GameMove:     nil,
		Searcher:     1,
		Evaluator:    ScoreGameState,
		MinimaxValue: 0,
	}

	move := MinimaxWithIterator(&root, 3)
	assert.NotNil(t, move)
	assert.NotNil(t, move.GameMove)
	assert.NotNil(t, move.GameMove.Mover)
}

func TestFromState(t *testing.T) {

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

	root := MinimaxNode{
		GameState:    state,
		GameMove:     nil,
		Searcher:     1,
		Evaluator:    ScoreGameState,
		MinimaxValue: 0,
	}

	move := MinimaxWithIterator(&root, 3)
	assert.NotNil(t, move)
	assert.NotNil(t, move.GameMove)
	assert.NotNil(t, move.GameMove.Mover)

}
