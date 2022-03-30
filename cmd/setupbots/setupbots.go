// Setupbots authenticates two users and starts a lockitdown game.
// Requires two usernames, one as host and the other as the opponent,
// and a server address.
//
// $> setupbos -host=host-bot -oppo=oppo-bot -server=https://boardbots.dev
package main

import (
	"flag"
	"fmt"

	"github.com/rwsargent/boardbots-go/client"
	"github.com/rwsargent/boardbots-go/lockitdown"
)

func main() {
	hostName := flag.String("host", "host-bot", "Username of the host of the game.")
	oppoName := flag.String("oppo", "oppo-bot", "Username of the opponent")
	server := flag.String("server", "http://localhost:8080", "Url of boardbot server")

	flag.Parse()

	hostClient, err := client.NewBoardBotClient[lockitdown.TransportState](client.Credentials{Username: *hostName}, *server)

	if err != nil {
		fmt.Println(err)
		return
	}

	oppoClient, err := client.NewBoardBotClient[lockitdown.TransportState](client.Credentials{Username: *oppoName}, *server)

	if err != nil {
		fmt.Println(err)
		return
	}

	Must(hostClient.Authenticate())

	lobby, err := hostClient.CreateLobby()

	if err != nil {
		fmt.Println(err)
		return
	}

	Must(oppoClient.Authenticate())

	Must(oppoClient.JoinLobby(lobby.Id.String()))

	game, err := hostClient.StartGame(lobby.Id.String())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(game.Id.String())
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
