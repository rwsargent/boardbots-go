package lockitdown

import (
	"fmt"
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
		move   Move
		player PlayerPosition
		err    error
	}{
		{PlaceRobot{
			Hex:       Pair{0, 5},
			Direction: Pair{0, -1},
		}, 0, nil},
		{PlaceRobot{
			Hex:       Pair{5, 0},
			Direction: Pair{-1, 0},
		}, 1, nil},
		{AdvanceRobot{
			Robot: Pair{0, 5},
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Left,
		}, 0, nil},
		{PlaceRobot{
			Hex:       Pair{-5, 4},
			Direction: Pair{1, 0},
		}, 1, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{AdvanceRobot{
			Robot: Pair{-5, 4},
		}, 1, nil},
		{AdvanceRobot{
			Robot: Pair{5, 0},
		}, 1, nil},
		{TurnRobot{
			Robot:     Pair{4, 0},
			Direction: Left,
		}, 1, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{PlaceRobot{
			Hex:       Pair{0, -5},
			Direction: Pair{0, 1},
		}, 1, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{TurnRobot{
			Robot:     Pair{0, 4},
			Direction: Right,
		}, 0, nil},
		{
			move: AdvanceRobot{
				Robot: Pair{0, -5},
			},
			player: 1,
			err:    nil,
		},
	}

	for _, tt := range tests {
		err := game.Move(tt.move, tt.player)
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
		move   Move
		player PlayerPosition
		result error
	}{
		{
			AdvanceRobot{
				Robot: Pair{4, -4},
			},
			0,
			nil,
		},
		{
			AdvanceRobot{
				Robot: Pair{5, -5},
			}, 0,
			nil,
		},
		{
			AdvanceRobot{
				Robot: Pair{0, 5},
			},
			0,
			nil,
		},
	}

	for _, tc := range testcases {
		err := gameState.Move(tc.move, tc.player)
		assert.Nilf(t, err, "unexpected err for move %v", tc.move)
		assert.Equal(t, -1, gameState.Winner)
		if err != nil {
			fmt.Println(gameState.ToJson())
		}
	}

	err := gameState.Move(TurnRobot{
		Robot:     Pair{-4, 4},
		Direction: Right,
	}, 1)

	fmt.Println(gameState.ToJson())

	assert.EqualError(t, err, "winner is 2")
	assert.Equal(t, 1, gameState.Winner)
}
