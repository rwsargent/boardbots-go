package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/rwsargent/boardbots-go/client"
	"github.com/rwsargent/boardbots-go/internal"
	"github.com/rwsargent/boardbots-go/lockitdown"
)

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

		root := &lockitdown.MinimaxNode{
			GameState: game,
			GameMove:  lockitdown.GameMove{},
			Searcher:  lockitdown.PlayerPosition(playerPosition - 1),
			Evaluator: lockitdown.ScoreGameState,
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc()
		best := lockitdown.AlphaBeta(ctx, root, 10)

		fmt.Printf("%s making move: %+v\n", bbClient.Credentials.Username, best.GameMove)
		movet := best.GameMove.ToTransport()
		moveCommand := client.MoveCommand{
			Json: client.MoveT{
				Player: playerPosition,
				Pos: internal.Position{
					Q: movet.Position.Q,
					R: movet.Position.R,
				},
				Action: movet.Action,
			},
		}
		fmt.Printf("Move Command: %v\n", moveCommand)
		state, err := bbClient.MakeMove(*gameId, moveCommand)
		if err != nil {
			panic(err)
		}
		game = lockitdown.StateFromTransport(&state)
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
