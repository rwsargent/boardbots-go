package quoridor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var TestIds = []uuid.UUID{
	uuid.MustParse("98ae983e-3f04-42ab-928a-c399d6d18375"),
	uuid.MustParse("5341acab-6e28-4d28-8530-8716e0c3dd2e"),
	uuid.MustParse("790bcc3f-6e72-4a0e-a6ea-bc806aa8aa03"),
	uuid.MustParse("6c8420b5-e7f5-4328-ae29-4dbdf7537612"),
	uuid.MustParse("f3282245-f546-4c71-92ca-5bada1f9c037"),
	uuid.MustParse("9cae2aa0-d21a-48ab-a877-4b78942259e4"),
	uuid.MustParse("0ad943b2-6ea9-45ad-9098-f67714652fcd"),
	uuid.MustParse("93ded37f-57d3-4b43-8933-1164e086a881"),
	uuid.MustParse("5b399bd3-aa3e-4754-bb51-175b30b77400"),
	uuid.MustParse("f7ea9019-033b-41e7-a671-26231952cd8c"),
}

func Test_NewGame(t *testing.T) {
	var testCases = []struct {
		id             uuid.UUID
		gameName       string
		expectGame     bool
		expectedErrMsg string
	}{
		{
			TestIds[0], "game 1", true, "",
		},
		{
			uuid.Nil, "game name", false, "unable to create game, need valid id",
		},
		{
			TestIds[0], "", false, "unable to create game, need non-empty name",
		},
	}
	for _, tc := range testCases {
		game, err := NewGame(tc.id, tc.gameName)
		if tc.expectGame {
			if game == nil {
				t.Fail()
			}
			assert.NotNil(t, game.Board)
			assert.Empty(t, game.Board)
			assert.NotNil(t, game.Players)
			assert.Empty(t, game.Players)
			assert.Equal(t, tc.gameName, game.Name)
			assert.Equal(t, PlayerOne, game.CurrentTurn)
			assert.Equal(t, tc.id, game.Id)
			assert.Equal(t, PlayerPosition(-1), game.Winner)
		} else {
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedErrMsg, err.Error())
		}
	}
}

func Test_AddPlayer(t *testing.T) {
	var testCases = []struct {
		name                 string
		expectedBarrier      int
		expectedPawnLocation Position
	}{
		{"playerOne", 10, Position{X: 8, Y: 16}},
		{"playerTwo", 10, Position{X: 8, Y: 0}},
		{"playerThree", 5, Position{X: 0, Y: 8}},
		{"playerFour", 5, Position{X: 16, Y: 8}},
	}
	game, _ := NewGame(TestIds[0], "AddPlayerGame")
	for idx, tc := range testCases {
		playerPosition, err := game.AddPlayer(TestIds[idx], tc.name)
		assert.Nil(t, err)
		assert.Equal(t, PlayerPosition(idx), playerPosition)
		assert.Len(t, game.Players, idx+1)
		for _, p := range game.Players {
			assert.Equal(t, tc.expectedBarrier, p.Barriers)
		}

		addedPlayer := game.Players[PlayerPosition(idx)]
		assert.Equal(t, tc.name, addedPlayer.PlayerName)
		assert.Equal(t, TestIds[idx], addedPlayer.PlayerId)
		assert.Equal(t, Piece{
			Position: tc.expectedPawnLocation,
			Owner:    PlayerPosition(idx),
			Type:     Pawn,
		}, addedPlayer.Pawn)

		assert.Len(t, game.Board, idx+1)
	}

	// With a full board, we can test error cases
	// First, can't add a player with the same id.
	_, err := game.AddPlayer(TestIds[1], "player one again")
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("player with id %s alreayd in this game", TestIds[1].String()))

	// Can't add more than 4 players
	_, err = game.AddPlayer(TestIds[5], "player five")
	assert.Error(t, err)
	assert.EqualError(t, err, "no open player positions in this game")

	// Can't add after a game has started.
	_ = game.StartGame()
	_, err = game.AddPlayer(TestIds[6], "late comer")
	assert.Error(t, err, fmt.Sprintf("cannot add player %s, game has already started", TestIds[6].String()))
}

func NewGameWithPlayers(players int) *Game {
	game, _ := NewGame(TestIds[0], "Four Player Game")
	for i := 0; i < players; i++ {
		_, _ = game.AddPlayer(TestIds[i+1], fmt.Sprint("Player ", i))
	}
	return game
}
func Test_AllPlayersMoveForward(t *testing.T) {
	game := NewGameWithPlayers(4)
	var testCases = []struct {
		position Position
		player   PlayerPosition
	}{
		{Position{X: 8, Y: 14}, PlayerOne},
		{Position{X: 8, Y: 2}, PlayerTwo},
		{Position{X: 2, Y: 8}, PlayerThree},
		{Position{X: 14, Y: 8}, PlayerFour},

		{Position{X: 8, Y: 12}, PlayerOne},
		{Position{X: 8, Y: 4}, PlayerTwo},
		{Position{X: 4, Y: 8}, PlayerThree},
		{Position{X: 12, Y: 8}, PlayerFour},

		{Position{X: 8, Y: 10}, PlayerOne},
		{Position{X: 8, Y: 6}, PlayerTwo},
		{Position{X: 6, Y: 8}, PlayerThree},
		{Position{X: 10, Y: 8}, PlayerFour},

		{Position{X: 8, Y: 8}, PlayerOne},
	}

	for idx, tc := range testCases {
		move := fmt.Sprint("Move ", idx)
		err := game.MovePawn(tc.position, tc.player)
		assert.Nil(t, err, move)
		pawn, ok := game.Board[tc.position]
		if ok {
			assert.Equal(t, tc.player, pawn.Owner, move)
			assert.Equal(t, tc.position, pawn.Position, move)
			assert.Equal(t, Pawn, pawn.Type, move)
		} else {
			t.Fail()
		}
		assert.Equal(t, game.Players[tc.player].Pawn, pawn, move)
		assert.Len(t, game.Board, 4, move) // No new pieces are ever added
	}

	err := game.MovePawn(Position{X: 8, Y: 8}, PlayerTwo)
	assert.EqualError(t, err, "the Pawn cannot reach that square")
	assert.Equal(t, PlayerTwo, game.CurrentTurn)

	// Test jump
	err = game.MovePawn(Position{X: 8, Y: 10}, PlayerTwo)
	assert.NoError(t, err)
	assert.Equal(t, PlayerThree, game.CurrentTurn)
}

func Test_InvalidMoves(t *testing.T) {
	board :=
		`........1........
.................
.................
.................
.................
.................
.................
.................
.................
.................
.................
.................
.................
.................
.......|.........
.......|---......
.......|0........`
	game, err := BuildQuoridorBoardFromString(board)
	assert.NoError(t, err)

	err = game.MovePawn(Position{X: 8, Y: 14}, PlayerOne)
	assert.EqualError(t, err, "the Pawn cannot reach that square")

	err = game.MovePawn(Position{X: 6, Y: 16}, PlayerOne)
	assert.EqualError(t, err, "the Pawn cannot reach that square")

	err = game.MovePawn(Position{X: 9, Y: 16}, PlayerOne)
	assert.EqualError(t, err, "invalid Pawn location")
}

func Test_DiagonalJump(t *testing.T) {
	board :=
		`.................
.................
.................
.................
.................
.................
.................
........---......
.......|1........
.......|.........
.......|0........
.................
.................
.................
.................
.................
.................`

	game, err := BuildQuoridorBoardFromString(board)
	assert.NoError(t, err)

	// Try going to the wrong diagonal.
	err = game.MovePawn(Position{X: 6, Y: 8}, PlayerOne)
	assert.EqualError(t, err, "the Pawn cannot reach that square")

	// Go to the correct diagonal
	err = game.MovePawn(Position{X: 10, Y: 8}, PlayerOne)
	assert.NoError(t, err)
	assert.Equal(t, Position{X: 10, Y: 8}, game.Players[PlayerOne].Pawn.Position)

	// Player two tries to use diagonal without a valid back barrier.
	err = game.MovePawn(Position{X: 10, Y: 6}, PlayerTwo)
	assert.EqualError(t, err, "the Pawn cannot reach that square")
}

func Test_PlaceVerticalBarrier(t *testing.T) {
	testCases := []struct {
		position Position
		player   PlayerPosition
	}{
		{Position{X: 1, Y: 0}, PlayerOne},
		{Position{X: 3, Y: 0}, PlayerTwo},
		{Position{X: 5, Y: 0}, PlayerThree},
		{Position{X: 7, Y: 0}, PlayerFour},

		{Position{X: 9, Y: 0}, PlayerOne},
		{Position{X: 11, Y: 0}, PlayerTwo},
		{Position{X: 13, Y: 0}, PlayerThree},
		{Position{X: 15, Y: 0}, PlayerFour},

		{Position{X: 1, Y: 4}, PlayerOne},
		{Position{X: 3, Y: 4}, PlayerTwo},
		{Position{X: 5, Y: 4}, PlayerThree},
		{Position{X: 7, Y: 4}, PlayerFour},

		{Position{X: 9, Y: 4}, PlayerOne},
		{Position{X: 11, Y: 4}, PlayerTwo},
		{Position{X: 13, Y: 4}, PlayerThree},
		{Position{X: 15, Y: 4}, PlayerFour},

		{Position{X: 1, Y: 8}, PlayerOne},
		{Position{X: 3, Y: 8}, PlayerTwo},
		{Position{X: 5, Y: 8}, PlayerThree},
		{Position{X: 7, Y: 8}, PlayerFour},
	}

	game := NewGameWithPlayers(4)

	for idx, tc := range testCases {
		err := game.PlaceBarrier(tc.position, tc.player)
		assert.NoError(t, err)
		for offset := 0; offset < 3; offset++ {
			expectedPosition := tc.position
			expectedPosition.Y += offset
			placedPiece, ok := game.Board[expectedPosition]
			assert.True(t, ok, "expected piece at %v", expectedPosition)
			assert.Equal(t, tc.player, placedPiece.Owner)
			assert.Equal(t, Barrier, placedPiece.Type)
		}
		assert.Equal(t, 4-(idx/4), game.Players[tc.player].Barriers)
	}
}

func Test_PlaceHorizontalBarrier(t *testing.T) {
	testCases := []struct {
		position Position
		player   PlayerPosition
	}{
		{Position{X: 0, Y: 1}, PlayerOne},
		{Position{X: 0, Y: 3}, PlayerTwo},
		{Position{X: 0, Y: 5}, PlayerThree},
		{Position{X: 0, Y: 7}, PlayerFour},

		{Position{X: 0, Y: 9}, PlayerOne},
		{Position{X: 0, Y: 11}, PlayerTwo},
		{Position{X: 0, Y: 13}, PlayerThree},
		{Position{X: 0, Y: 15}, PlayerFour},

		{Position{X: 4, Y: 1}, PlayerOne},
		{Position{X: 4, Y: 3}, PlayerTwo},
		{Position{X: 4, Y: 5}, PlayerThree},
		{Position{X: 4, Y: 7}, PlayerFour},

		{Position{X: 4, Y: 9}, PlayerOne},
		{Position{X: 4, Y: 11}, PlayerTwo},
		{Position{X: 4, Y: 13}, PlayerThree},
		{Position{X: 4, Y: 15}, PlayerFour},

		{Position{X: 8, Y: 1}, PlayerOne},
		{Position{X: 8, Y: 3}, PlayerTwo},
		{Position{X: 8, Y: 5}, PlayerThree},
		{Position{X: 8, Y: 7}, PlayerFour},
	}

	game := NewGameWithPlayers(4)

	for idx, tc := range testCases {
		err := game.PlaceBarrier(tc.position, tc.player)
		assert.NoError(t, err)
		for offset := 0; offset < 3; offset++ {
			expectedPosition := tc.position
			expectedPosition.X += offset
			placedPiece, ok := game.Board[expectedPosition]
			assert.True(t, ok, "expected piece at %v", expectedPosition)
			assert.Equal(t, tc.player, placedPiece.Owner)
			assert.Equal(t, Barrier, placedPiece.Type)
		}
		assert.Equal(t, 4-(idx/4), game.Players[tc.player].Barriers)
	}
}

func Test_PlaceBarrierErrors(t *testing.T) {
	board :=
		`........1........
---.---.---.---..
.................
..............---
.................
.................
.................
.................
3...............4
.................
.................
.................
.................
.................
.................
.................
........0........`
	game, err := BuildQuoridorBoardFromString(board)
	assert.NoError(t, err)

	err = game.PlaceBarrier(Position{X: 13, Y: 2}, PlayerOne)
	assert.EqualError(t, err, "the barrier prevents a players victory")

	err = game.PlaceBarrier(Position{X: 9, Y: 9}, PlayerOne)
	assert.EqualError(t, err, "invalid location for a barrier")

	err = game.PlaceBarrier(Position{X: 1, Y: 0}, PlayerOne)
	assert.EqualError(t, err, "the new barrier intersects with another")
}

func Test_WinCondition(t *testing.T) {
	game := NewGameWithPlayers(2)
	turns := []struct {
		position Position
		player   PlayerPosition
	}{
		{Position{X: 8, Y: 14}, PlayerOne},
		{Position{X: 8, Y: 2}, PlayerTwo},

		{Position{X: 8, Y: 12}, PlayerOne},
		{Position{X: 8, Y: 4}, PlayerTwo},

		{Position{X: 8, Y: 10}, PlayerOne},
		{Position{X: 8, Y: 6}, PlayerTwo},

		{Position{X: 8, Y: 8}, PlayerOne},
		{Position{X: 8, Y: 10}, PlayerTwo},

		{Position{X: 8, Y: 6}, PlayerOne},
		{Position{X: 8, Y: 12}, PlayerTwo},

		{Position{X: 8, Y: 4}, PlayerOne},
		{Position{X: 8, Y: 14}, PlayerTwo},

		{Position{X: 8, Y: 2}, PlayerOne},
		{Position{X: 8, Y: 16}, PlayerTwo},
	}

	for _, turn := range turns {
		err := game.MovePawn(turn.position, turn.player)
		assert.NoError(t, err)
	}

	err := game.MovePawn(Position{X: 8, Y: 0}, PlayerOne)
	assert.EqualError(t, err, "invalid move, game is already over")

	err = game.PlaceBarrier(Position{X: 1, Y: 0}, PlayerOne)
	assert.EqualError(t, err, "invalid move, game is already over")

	assert.Equal(t, PlayerTwo, game.Winner)
	assert.False(t, game.EndDate.IsZero())
}
