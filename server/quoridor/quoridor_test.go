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

func NewFourPlayerGame() *Game {
	game, _ := NewGame(TestIds[0], "Four Player Game")
	for i := 0; i < 4; i++ {
		_, _ = game.AddPlayer(TestIds[i+1], fmt.Sprint("Player ", i))
	}
	return game
}

func Test_AllPlayersMoveForward(t *testing.T) {
	game := NewFourPlayerGame()
	var testCases = []struct {
		position Position
		player   PlayerPosition
	}{
		{Position{X: 8, Y: 14}, PlayerOne},
		{Position{X: 8, Y: 2}, PlayerTwo},
		{Position{X: 14, Y: 8}, PlayerThree},
		{Position{X: 2, Y: 8}, PlayerFour},

		{Position{X: 8, Y: 12}, PlayerOne},
		{Position{X: 8, Y: 4}, PlayerTwo},
		{Position{X: 12, Y: 8}, PlayerThree},
		{Position{X: 4, Y: 8}, PlayerFour},

		{Position{X: 8, Y: 10}, PlayerOne},
		{Position{X: 8, Y: 6}, PlayerTwo},
		{Position{X: 10, Y: 8}, PlayerThree},
		{Position{X: 6, Y: 8}, PlayerFour},

		{Position{X: 8, Y: 8}, PlayerOne},
		{Position{X: 8, Y: 6}, PlayerTwo},
		{Position{X: 8, Y: 8}, PlayerThree},
		{Position{X: 6, Y: 8}, PlayerFour},
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
}

//func Test_FirstMoveIsValid(t *testing.T) {
//	game := NewGame(TestUUID, "")
//	board, err := game.MovePawn(Position{2, 8}, PlayerOne)
//	assert.Nil(t, err, "Valid Move")
//	assert.Len(t, board, 2, "Game has two pawns")
//
//	assert.Equal(t, game.Players[PlayerOne].Pawn.Position, Position{2, 8})
//}
//
//func Test_FirstMoveIsAnInvalidMove(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	_, err := game.MovePawn(Position{4, 4}, PlayerOne)
//	assert.EqualError(t, err, "the pawn cannot reach that square")
//}
//
//func Test_PlaceBarrier(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	placePosition := Position{1, 6}
//	board, err := game.PlaceBarrier(placePosition, PlayerOne)
//	assert.NotNil(t, board[placePosition])
//	assert.NotNil(t, board[Position{1, 7}])
//	assert.NotNil(t, board[Position{1, 8}])
//	_, present := board[Position{1, 9}]
//	assert.False(t, present)
//	assert.Nil(t, err, "Valid placement")
//}
//
//func Test_PlaceBarrierWithNoMoreBarriers(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	for i := 0; i < 18; i++ {
//		Y := 14 - (4 * (i / 8)) // offset after a Y as filled up
//		X := 1 + (2 * (i % 8))  // wrap around to the same X, since Y changes.
//		position := Position{Y: Y, X: X}
//		board, err := game.PlaceBarrier(position, PlayerPosition(i%2))
//		assert.Nil(t, err, "No error expected")
//		assert.NotNil(t, board[position], "board should be placed")
//		assert.NotNil(t, game.Board[position], "board should be placed")
//	}
//	game.PlaceBarrier(Position{Y: 0, X: 3}, PlayerOne)
//	board, err := game.PlaceBarrier(Position{Y: 0, X: 1}, PlayerTwo)
//	assert.Nil(t, err, "No error expected")
//	assert.NotNil(t, board[Position{Y: 0, X: 1}], "board should be placed")
//
//	illegalPosition := Position{Y: 0, X: 7}
//	board, err = game.PlaceBarrier(illegalPosition, PlayerOne)
//	assert.NotNil(t, err, "expect 11th barrier to thY error")
//	_, ok := board[illegalPosition]
//	assert.False(t, ok, "barrier should not be on board")
//	_, ok = game.Board[illegalPosition]
//	assert.False(t, ok)
//}
//
//func Test_PawnCannotMoveOverBarrier(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	game.PlaceBarrier(Position{1, 8}, PlayerOne)
//	game.PlaceBarrier(Position{5, 6}, PlayerTwo)
//	_, err := game.MovePawn(Position{2, 8}, PlayerOne)
//	assert.NotNil(t, err, "Expected error")
//	assert.Equal(t, "the pawn cannot reach that square", err.Error(), "Wrong message")
//}
//
//func Test_OneValidDirection(t *testing.T) {
//	board := make(Board)
//	placePawn(board)
//	placeBarrier(board)
//
//	moves := board.GetValidPawnMoves(Position{0, 0})
//	assert.Len(t, moves, 1)
//	assert.Equal(t, moves[0], Position{2, 0})
//}
//
//func Test_JumpPawn(t *testing.T) {
//	board := make(Board)
//	setupPawnBlockingBoard(board)
//	moves := board.GetValidPawnMoves(Position{0, 8})
//	assert.Len(t, moves, 1)
//	assert.Equal(t, moves[0], Position{4, 8})
//}
//
//func Test_JumpDiagonal(t *testing.T) {
//	board := make(Board)
//	board[Position{2, 8}] = Piece{}
//	board[Position{4, 8}] = Piece{}
//
//	//Barriers
//	board[Position{0, 9}] = Piece{}
//	board[Position{1, 9}] = Piece{}
//	board[Position{2, 9}] = Piece{}
//
//	board[Position{0, 7}] = Piece{}
//	board[Position{1, 7}] = Piece{}
//	board[Position{2, 7}] = Piece{}
//
//	board[Position{5, 8}] = Piece{}
//	board[Position{5, 9}] = Piece{}
//	board[Position{5, 10}] = Piece{}
//
//	moves := board.GetValidPawnMoves(Position{2, 8})
//	assert.Len(t, moves, 3)
//
//	expecteds := map[Position]bool{
//		{4, 10}: true,
//		{4, 6}:  true,
//		{0, 8}:  true,
//	}
//	for _, move := range moves {
//		if _, ok := expecteds[move]; !ok {
//			t.Errorf("Missing expected move %v", move)
//		}
//	}
//}
//
//func Test_FourDirectionsForFree(t *testing.T) {
//	board := make(Board)
//	board[Position{2, 8}] = Piece{}
//
//	moves := board.GetValidPawnMoves(Position{2, 8})
//
//	assert.Len(t, moves, 4)
//}
//
//func Test_GameNotOver(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	assert.Equal(t, PlayerPosition(-1), game.MaybeReturnWinnerPlayerPosition())
//}
//
//func Test_GameOver(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	p1Pawn := &game.Players[PlayerOne].Pawn
//	delete(game.Board, p1Pawn.Position)
//	p1Pawn.Position = Position{16, 8}
//	game.Board[p1Pawn.Position] = *p1Pawn
//
//	assert.Equal(t, PlayerOne, game.MaybeReturnWinnerPlayerPosition())
//}
//
//func Test_BarrierPlacementBlocksWin(t *testing.T) {
//	var board = `.......|0|.......
//.......|.|.......
//.......|.|.......
//.................
//.................
//.................
//.................
//.................
//.................
//.................
//.................
//.................
//.................
//.................
//.................
//.................
//................`
//
//	game, err := BuildQuoidorBoardFromString(board)
//	if err != nil {
//		t.Error(err.Error())
//	}
//
//	newBoard, err := game.PlaceBarrier(Position{3, 8}, PlayerOne)
//	assert.NotNil(t, err)
//	assertNoPiece(t, newBoard, Position{Y: 3, X: 8})
//	assertNoPiece(t, newBoard, Position{Y: 3, X: 9})
//	assertNoPiece(t, newBoard, Position{Y: 3, X: 10})
//}
//
//func assertNoPiece(t *testing.T, board Board, position Position) {
//	if _, exists := board[position]; exists {
//		t.Error("found unexpected piece at Position: ", position)
//	}
//}
//
//func TestAddPlayer(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	playerId := uuid.New()
//	player, err := game.AddPlayer(playerId, "Test Name")
//	assert.NoError(t, err)
//	assert.Equal(t, PlayerOne, player)
//
//	player, err = game.AddPlayer(uuid.New(), "Test Name")
//	assert.NoError(t, err)
//	assert.Equal(t, PlayerTwo, player)
//}
//
//func Test_NewGame_PlayerOnesTurn(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	assert.Equal(t, game.CurrentTurn, PlayerOne)
//}
//
//func Test_MovePawn_FailsIfWrongPlayer(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	_, err := game.MovePawn(Position{14, 8}, PlayerTwo)
//	assert.Error(t, err)
//}
//
//func Test_MovePawn_TurnCheck(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	game.MovePawn(Position{2, 8}, PlayerOne)
//
//	assert.Equal(t, game.CurrentTurn, PlayerTwo)
//}
//
//func Test_MovePawnTwice_TurnCheck(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//	game.MovePawn(Position{2, 8}, PlayerOne)
//	game.MovePawn(Position{14, 8}, PlayerTwo)
//
//	assert.Equal(t, game.CurrentTurn, PlayerOne)
//}
//
//func Test_MovePawnFourTimes_FourPersonGame(t *testing.T) {
//	var err error
//	game := NewFourPersonGame(TestUUID)
//	game.MovePawn(Position{2, 8}, PlayerOne)
//	game.MovePawn(Position{14, 8}, PlayerTwo)
//	_, err = game.MovePawn(Position{8, 2}, PlayerThree)
//	assert.NoError(t, err)
//	_, err = game.MovePawn(Position{8, 14}, PlayerFour)
//	assert.NoError(t, err)
//
//	assert.Equal(t, game.CurrentTurn, PlayerOne)
//}
//
//func Test_PlaceBarrier_WrongTurn(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//
//	_, err := game.PlaceBarrier(Position{1, 2}, PlayerTwo)
//
//	assert.Error(t, err)
//	assert.EqualError(t, err, "wrong turn, current turn is for Player: 0")
//}
//
//func Test_PlaceBarrierChangesTurnTwoPlayer(t *testing.T) {
//	game := NewTwoPersonGame(TestUUID)
//
//	game.PlaceBarrier(Position{1, 2}, PlayerOne)
//
//	assert.Equal(t, game.CurrentTurn, PlayerTwo)
//}
//
//func Test_MultiplePlaceBarrier(t *testing.T) {
//	game := NewFourPersonGame(TestUUID)
//
//	game.PlaceBarrier(Position{1, 2}, PlayerOne)
//	game.PlaceBarrier(Position{3, 2}, PlayerTwo)
//	game.PlaceBarrier(Position{5, 2}, PlayerThree)
//	game.PlaceBarrier(Position{7, 2}, PlayerFour)
//
//	assert.Equal(t, game.CurrentTurn, PlayerOne)
//
//}

func setupPawnBlockingBoard(board Board) {
	//Pawns
	board[Position{0, 8}] = Piece{}
	board[Position{2, 8}] = Piece{}

	//Barriers
	board[Position{0, 9}] = Piece{}
	board[Position{1, 9}] = Piece{}
	board[Position{2, 9}] = Piece{}
	board[Position{0, 7}] = Piece{}
	board[Position{1, 7}] = Piece{}
	board[Position{2, 7}] = Piece{}
}

func placeBarrier(board Board) {
	board[Position{0, 1}] = Piece{}
	board[Position{1, 1}] = Piece{}
	board[Position{1, 1}] = Piece{}
}

func placePawn(board Board) {
	board[Position{0, 0}] = Piece{}
}

func assertNoExtraPlayersCreated(t *testing.T, game *Game) {
	assert.Nil(t, game.Players[PlayerThree])
	assert.Nil(t, game.Players[PlayerFour])
}

func assertCorrectPlayerInit(t *testing.T, game *Game) {
	assert.NotNil(t, game.Players[PlayerOne], "Initialized Player one!")
	assert.NotNil(t, game.Players[PlayerTwo], "Initialized Player two!")
	assert.Equal(t, Position{0, 8}, game.Players[PlayerOne].Pawn.Position)
	assert.Equal(t, Position{16, 8}, game.Players[PlayerTwo].Pawn.Position)
	assert.Equal(t, 10, game.Players[PlayerOne].Barriers)
	assert.Equal(t, 10, game.Players[PlayerTwo].Barriers)
}

func assertBaseGame(t *testing.T, game *Game) {
	assert.NotNil(t, game.Board, "Board should be initialized")
	assert.NotNil(t, game.Board, "Pieces of board should be initialized")
	assert.Equal(t, 2, len(game.Players), "Should be 4 players")
}
