// randobot.go makes random moves for the given lockitdown game.
// You need to run randobot with a username, gameId, and server address.
//
// $> randobot -username=randobot -gameId=00000000-...-0000 -server=https://boardbots.dev
//
// It is recommend to run setupbots before using randobot.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/rwsargent/boardbots-go/client"
	"github.com/rwsargent/boardbots-go/internal"
	"github.com/rwsargent/boardbots-go/lockitdown"
)

var edges = []struct{ pos, dir internal.Position }{
	{pos: internal.Position{Q: 0, R: -5}, dir: internal.Position{Q: 0, R: 1}},
	{pos: internal.Position{Q: 5, R: -5}, dir: internal.Position{Q: -1, R: 1}},
	{pos: internal.Position{Q: 5, R: 0}, dir: internal.Position{Q: -1, R: 0}},
	{pos: internal.Position{Q: 0, R: 5}, dir: internal.Position{Q: 0, R: -1}},
	{pos: internal.Position{Q: -5, R: 5}, dir: internal.Position{Q: 1, R: -1}},
	{pos: internal.Position{Q: -5, R: 0}, dir: internal.Position{Q: 1, R: 0}},
}

func main() {

	server := flag.String("server", "http://localhost:8080", "Host of the boardbots server to play on.")
	username := flag.String("username", "", "Username")
	gameId := flag.String("gameId", "", "Game ID")

	flag.Parse()

	if *gameId == "" || *username == "" {
		fmt.Println("Require a game ID and username")
	}

	bbClient, err := client.NewBoardBotClient[lockitdown.TransportState](client.Credentials{
		Username: *username,
	}, *server)

	if err != nil {
		fmt.Printf("failed to start client, %s\n", err.Error())
		return
	}

	err = bbClient.Authenticate()
	if err != nil {
		fmt.Printf("could not authenticate %s", err.Error())
	}

	tGame, err := bbClient.Game(*gameId)
	if err != nil {
		panic(err)
	}
	playerPosition := getPlayerPosition(tGame, bbClient.Credentials.Username)

	game := lockitdown.StateFromTransport(&tGame.State)

	rand.Seed(63)

	for game.Winner < 0 {
		if playerPosition-1 != int(game.PlayerTurn) {
			fmt.Println("waiting for turn")
			time.Sleep(time.Second * 3)

			tGame, err := bbClient.Game(*gameId)
			if err != nil {
				panic(err)
			}

			game = lockitdown.StateFromTransport(&tGame.State)
			continue
		}

		playerMoves, err := movesForPlayer(bbClient, gameId, playerPosition)

		var moveCommand client.MoveCommand

		if len(playerMoves) == 0 {
			moveCommand = placeRobotMove(game, playerPosition)
		} else {
			moveCommand = client.MoveCommand{
				Json: playerMoves[rand.Intn(len(playerMoves))],
			}
		}

		fmt.Printf("%s making move: %+v\n", bbClient.Credentials.Username, moveCommand)
		state, err := bbClient.MakeMove(*gameId, moveCommand)
		if err != nil {
			panic(err)
		}
		game = lockitdown.StateFromTransport(&state)
	}
}

func movesForPlayer(bbClient *client.BoardBotClient[lockitdown.TransportState], gameId *string, playerPosition int) ([]client.MoveT, error) {
	moves, err := bbClient.GetPossibleMoves(*gameId)
	if err != nil {
		panic(err)
	}

	playerMoves := make([]client.MoveT, 0, len(moves))

	for _, move := range moves {
		if move.Player == playerPosition {
			playerMoves = append(playerMoves, move)
		}
	}
	return playerMoves, err
}

func placeRobotMove(game *lockitdown.GameState, playerPosition int) client.MoveCommand {
	placePosition := edges[rand.Intn(len(edges))]

	for bot := game.RobotAt(lockitdown.Pair{Q: placePosition.pos.Q, R: placePosition.pos.R}); bot != nil; {
		placePosition = edges[rand.Intn(len(edges))]
	}

	return client.MoveCommand{
		Json: client.MoveT{
			Player: playerPosition,
			Pos:    placePosition.pos,
			Action: lockitdown.PlaceRobotT{
				PlaceRobot: struct {
					Dir lockitdown.Pair `json:"dir"`
				}{Dir: lockitdown.Pair{Q: placePosition.dir.Q, R: placePosition.dir.R}},
			},
		},
	}
}

func getPlayerPosition[S any](game client.Game[S], username string) int {
	for idx, user := range game.Players {
		if user.Name == username {
			return idx + 1
		}
	}
	return -1
}
