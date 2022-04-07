package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/rwsargent/boardbots-go/lockitdown"
)

type (
	ScoreRequest struct {
		GameType  string                    `json:"gameType"`
		GameState lockitdown.TransportState `json:"state"`
		Strategy  string                    `json:"strategy"`
		Player    int                       `json:"player"`
	}

	ScoreResponse struct {
		Score int `json:"score"`
	}
)

func main() {

	port := flag.String("port", ":8888", "server port")
	flag.Parse()

	http.HandleFunc("/score", score)

	http.ListenAndServe(*port, nil)
}

func score(w http.ResponseWriter, req *http.Request) {
	var scoreReqest ScoreRequest
	err := json.NewDecoder(req.Body).Decode(&scoreReqest)
	defer req.Body.Close()

	if err != nil {
		fmt.Printf("error reading body, %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	state := lockitdown.StateFromTransport(&scoreReqest.GameState)
	var score int
	switch scoreReqest.Strategy {
	default:
		score = lockitdown.ScoreGameState(state, lockitdown.PlayerPosition(scoreReqest.Player-1))
	}

	resp := ScoreResponse{
		score,
	}
	out, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error marshaling response, %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = w.Write(out)
	if err != nil {
		fmt.Printf("error writing response, %v", err)
	}
}
