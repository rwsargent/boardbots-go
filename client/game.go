package client

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rwsargent/boardbots-go/internal"
)

// game.go is responsible for the game specific APIs in boardbots.dev

type (
	Game[S any] struct {
		Id        uuid.UUID   `json:"id"`
		LobbyId   uuid.UUID   `json:"lobbyId"`
		Players   []User      `json:"players"`
		GameType  string      `json:"gameType"`
		State     S           `json:"state"`
		Status    string      `json:"status"`
		NumMoves  int         `json:"numMoves"`
		StartedAt json.Number `json:"startedAt"`
	}

	Player struct {
		Player   int         `json:"player"`
		Username string      `json:"username"`
		UserId   json.Number `json:"userId"`
	}

	MoveCommand struct {
		Json any `json:"json"`
	}

	MoveResp struct {
		Index  int    `json:"index"`
		Player Player `json:"player"`
	}

	MoveT struct {
		Player int               `json:"player"`
		Pos    internal.Position `json:"pos"`
		Action interface{}       `json:"action"`
	}
)

func (c *BoardBotClient[S]) Game(gameId string) (Game[S], error) {
	return Get[Game[S]](c, fmt.Sprintf("/api/game/%s", gameId))
}

func (c *BoardBotClient[S]) MakeMove(gameId string, move MoveCommand) (S, error) {
	return Post[MoveCommand, S](c, fmt.Sprintf("/api/game/%s/move", gameId), move)
}

func (c *BoardBotClient[S]) GetPossibleMoves(gameId string) ([]MoveT, error) {
	return Get[[]MoveT](c, fmt.Sprintf("/api/game/%s/potential-moves", gameId))
}
