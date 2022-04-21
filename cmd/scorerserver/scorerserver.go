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

	http.HandleFunc("/api/score", score)

	fmt.Printf("Now listening on port %s\n", *port)
	http.ListenAndServe(*port, nil)

	fmt.Println("done")
}

func score(w http.ResponseWriter, req *http.Request) {
	var scoreReqest ScoreRequest
	err := json.NewDecoder(req.Body).Decode(&scoreReqest)
	defer req.Body.Close()

	fmt.Printf("Request:\n%+v\n", scoreReqest)

	if err != nil {
		fmt.Printf("error reading body, %v\n", err)
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
	fmt.Printf("Response:\n%+v\n", resp)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		fmt.Printf("error writing response, %v", err)
	}
}
