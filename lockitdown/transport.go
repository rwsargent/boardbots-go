package lockitdown

import (
	"fmt"
	"strconv"
)

type (
	TransportRobot struct {
		Player        int  `json:"player"`
		Dir           Pair `json:"dir"`
		IsLocked      bool `json:"isLocked"`
		IsBeamEnabled bool `json:"isBeamEnabled"`
	}
	TransportRobots []interface{}
	TransportState  struct {
		GameDef          GameDef           `json:"gameDef"`
		Players          []Player          `json:"players"`
		Robots           []TransportRobots `json:"robots"`
		PlayerTurn       int               `json:"playerTurn"`
		Status           interface{}       `json:"status"`
		MovesThisTurn    int               `json:"movesThisTurn"`
		RequiresTieBreak bool              `json:"requiresTieBreak"`
	}
)

func ConvertToTransport(game *GameState) *TransportState {
	players := make([]Player, len(game.Players))
	for i, player := range game.Players {
		players[i] = *player
	}

	robots := make([]TransportRobots, len(game.Robots))
	idx := 0
	for _, robot := range game.Robots {
		robots[idx] = []interface{}{
			robot.Position,
			TransportRobot{
				Player:        int(robot.Player) + 1,
				Dir:           robot.Direction,
				IsLocked:      robot.IsLockedDown,
				IsBeamEnabled: robot.IsBeamEnabled,
			},
		}
		idx++
	}

	var status string
	if game.Winner < 0 {
		status = "OnGoing"
	} else {
		status = fmt.Sprintf("%d", game.Winner)
	}

	return &TransportState{
		GameDef:          game.GameDef,
		Players:          players,
		Robots:           robots,
		MovesThisTurn:    game.GameDef.MovesPerTurn - game.MovesThisTurn,
		Status:           status,
		RequiresTieBreak: game.RequiresTieBreak,
		PlayerTurn:       int(game.PlayerTurn) + 1,
	}
}

func StateFromTransport(tState *TransportState) *GameState {
	players := make([]*Player, 0, len(tState.Players))
	for _, player := range tState.Players {
		players = append(players, &Player{
			Points:       player.Points,
			PlacedRobots: player.PlacedRobots,
		})
	}

	robots := make([]Robot, len(tState.Robots))
	for i, robot := range tState.Robots {
		tRobot := TransportRobotFromMap(robot[1].(map[string]interface{}))
		position := PairFromMap(robot[0].(map[string]interface{}))
		robots[i] = Robot{
			Position:      position,
			Direction:     tRobot.Dir,
			IsBeamEnabled: tRobot.IsBeamEnabled,
			IsLockedDown:  tRobot.IsLocked,
			Player:        PlayerPosition(tRobot.Player - 1),
		}
	}

	winner := -1
	if tState.Status != "OnGoing" {
		winner, _ = strconv.Atoi(tState.Status.(string))
	}

	return &GameState{
		GameDef:          tState.GameDef,
		Players:          players,
		Robots:           robots,
		PlayerTurn:       PlayerPosition(tState.PlayerTurn - 1),
		MovesThisTurn:    tState.GameDef.MovesPerTurn - tState.MovesThisTurn,
		RequiresTieBreak: tState.RequiresTieBreak,
		Winner:           winner,
	}
}

func PairFromMap(json map[string]interface{}) Pair {
	return Pair{
		Q: int(json["q"].(float64)),
		R: int(json["r"].(float64)),
	}
}

func TransportRobotFromMap(json map[string]interface{}) TransportRobot {
	dir := json["dir"].(map[string]interface{})
	return TransportRobot{
		Player: int(json["player"].(float64)),
		Dir: Pair{
			Q: int(dir["q"].(float64)),
			R: int(dir["r"].(float64)),
		},
		IsLocked:      json["isLocked"].(bool),
		IsBeamEnabled: json["isBeamEnabled"].(bool),
	}
}
