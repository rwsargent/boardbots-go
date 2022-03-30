package client

import (
	"fmt"
	"testing"

	"github.com/rwsargent/boardbots-go/lockitdown"
)

func TestMoves(t *testing.T) {

	bbClient, err := NewBoardBotClient[lockitdown.TransportState](Credentials{
		Username: "tester",
	}, "http://localhost:8080")

	if err != nil {
		fmt.Printf("failed to start client, %s\n", err.Error())
		return
	}

	err = bbClient.Authenticate()
	if err != nil {
		fmt.Printf("could not authenticate %s", err.Error())
	}

	resp, err := bbClient.GetPossibleMoves("3155e207-88c7-437e-bc9c-df820961ed07")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%v+\n", resp)
}
