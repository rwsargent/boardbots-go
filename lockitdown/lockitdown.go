// An implmentation of boardbots.dev LockItDown game. Users can use package
// to make and undo moves on an internal state of the board, and then apply
// those moves to a game hosted on boardbots.dev.
package lockitdown

import (
	"encoding/json"
	"fmt"
)

type (
	Pair struct {
		Q int `json:"q"`
		R int `json:"r"`
	}
	Player struct {
		Points       int `json:"points"`
		PlacedRobots int `json:"placedRobots"`
	}

	PlayerPosition int

	BoardType struct {
		ArenaRadius int `json:"arenaRadius"`
	}

	Board struct {
		HexaBoard BoardType
	}

	GameDef struct {
		Board           Board  `json:"board"`
		Players         int    `json:"numOfPlayers"`
		MovesPerTurn    int    `json:"movesPerTurn"`
		RobotsPerPlayer int    `json:"robotsPerPlayer"`
		WinCondition    string `json:"winCondition"`
	}

	Robot struct {
		Position                    Pair
		Direction                   Pair
		IsBeamEnabled, IsLockedDown bool
		Player                      PlayerPosition
	}

	GameState struct {
		GameDef          GameDef
		Players          []*Player
		Robots           map[Pair]*Robot
		PlayerTurn       PlayerPosition
		MovesThisTurn    int
		RequiresTieBreak bool
		Winner           int
	}

	TurnDirection int

	TieBreak struct {
		Robots []*Robot
	}
)

const (
	Left TurnDirection = iota
	Right
)

var (
	NW Pair = Pair{0, -1}
	NE Pair = Pair{1, -1}
	E  Pair = Pair{1, 0}
	SE Pair = Pair{0, 1}
	SW Pair = Pair{-1, 1}
	W  Pair = Pair{-1, 0}
)

func (p *Pair) Plus(that Pair) {
	p.Q += that.Q
	p.R += that.R
}

func (p *Pair) Minus(that Pair) {
	p.Q -= that.Q
	p.R -= that.R
}

func (p Pair) S() int {
	return -p.Q - p.R
}

func NewGame(gameDef GameDef) *GameState {
	players := make([]*Player, gameDef.Players)
	for i := 0; i < len(players); i++ {
		players[i] = &Player{}
	}
	return &GameState{
		GameDef:          gameDef,
		Players:          players,
		Robots:           make(map[Pair]*Robot),
		PlayerTurn:       0,
		MovesThisTurn:    gameDef.MovesPerTurn,
		RequiresTieBreak: false,
		Winner:           -1,
	}
}

// Only intended for Unit pairs
func (p *Pair) Rotate(direction TurnDirection) {
	s := p.S()
	if direction == Right {
		p.Q = -p.R
		p.R = -s
	} else {
		p.R = -p.Q
		p.Q = -s
	}
}

func (game *GameState) Move(move Move, player PlayerPosition) error {
	if player != PlayerPosition(game.PlayerTurn) {
		return fmt.Errorf("wrong player, expected %d, was %d", game.PlayerTurn, player)
	}
	err := move.Move(game, player)
	if err != nil {
		return err
	}

	// Resolve move
	if err = game.resolveMove(); err != nil {
		return err
	}

	if game.MovesThisTurn == 0 {
		game.PlayerTurn = PlayerPosition((int(game.PlayerTurn) + 1) % len(game.Players))
		game.MovesThisTurn = 3
	}

	if over, winner := game.checkGameOver(); over {
		game.Winner = winner
		return fmt.Errorf("winner is %d", winner+1)
	}
	return nil
}

func (game *GameState) resolveMove() error {
	for resolved := false; !resolved; {
		targeted := game.taretedRobots()

		if tiebreaks := game.checkForTieBreaks(targeted); len(tiebreaks) > 0 {
			game.RequiresTieBreak = true
			return TieBreak{
				Robots: tiebreaks,
			}
		}

		resolved = game.updateLockedRobots(targeted)
	}
	return nil
}

func (game *GameState) updateLockedRobots(targeted map[Pair][]*Robot) bool {
	resolved := true
	for hex, attackers := range targeted {
		if len(attackers) == 1 {
			continue
		}
		if len(attackers) == 3 {
			game.shutdownRobot(hex, attackers)
			resolved = false
		}
		if len(attackers) == 2 {
			if locked, ok := game.Robots[hex]; ok {
				locked.IsLockedDown = true
			}
		}
	}
	return resolved
}

// If any "doomed" robots (locked or shutdown) are also part of a lock or shut down,
// we need to break a tie.
func (game *GameState) checkForTieBreaks(targeted map[Pair][]*Robot) []*Robot {
	tiebreaks := make([]*Robot, 0, 2)
	for doomed, attackers := range targeted {
		if len(attackers) > 1 {
			for _, attacker := range attackers {
				for doomedAttacker, doomedAttackerAttackers := range targeted {
					if doomedAttacker == attacker.Position && len(doomedAttackerAttackers) > 1 {
						tiebreaks = append(tiebreaks, game.Robots[doomed])
					}
				}
			}
		}
	}
	return tiebreaks
}

// Returns a map of locations of robots and which robots are pointing
// at them.
func (game *GameState) taretedRobots() map[Pair][]*Robot {
	targeted := make(map[Pair][]*Robot)
	for _, attacker := range game.Robots {
		if !attacker.IsBeamEnabled || attacker.IsLockedDown || game.isCorridor(attacker.Position) {
			continue
		}

		// add hexes to contended
		cursor := Pair{
			Q: attacker.Position.Q,
			R: attacker.Position.R,
		}
		cursor.Plus(attacker.Direction)

		for ; !game.isCorridor(cursor); cursor.Plus(attacker.Direction) {
			if targetedBot, hit := game.Robots[cursor]; hit {
				if targetedBot.Player == attacker.Player {
					break
				}
				// Add to attackers list
				if attackers, found := targeted[cursor]; found {
					targeted[cursor] = append(attackers, attacker)
				} else {
					l := make([]*Robot, 0, 10)
					targeted[cursor] = append(l, attacker)
				}
				break
			}
		}
	}

	return targeted
}

func (game *GameState) isCorridor(pair Pair) bool {
	s := pair.S()
	board := game.GameDef.Board.HexaBoard
	corridor := board.ArenaRadius + 1
	return intAbs(s) == corridor || intAbs(pair.Q) == corridor || intAbs(pair.R) == corridor
}

func (game *GameState) shutdownRobot(hex Pair, attackers []*Robot) {
	for _, attacker := range attackers {
		game.Players[attacker.Player].Points += 1
	}

	delete(game.Robots, hex)
}

func (game *GameState) checkGameOver() (bool, int) {
	if game.GameDef.WinCondition == "Elimination" {
		winner := 0
		bots := make(map[int]int)
		// Count all robots on the board
		for _, robot := range game.Robots {
			bots[int(robot.Player)] += 1
		}
		eliminated := 0
		for position, player := range game.Players {
			// XOR player, we'll un-XOR later to get survivor
			winner ^= position
			// If only two robots remain, the player is eliminated
			if game.GameDef.RobotsPerPlayer-player.PlacedRobots+bots[position] <= 2 {
				eliminated++
				// remove it from winner aggregator
				winner ^= position
			}
		}
		// If all but one player is eliminated, the remaining player is the winner!
		if eliminated+1 == len(game.Players) {
			return true, winner
		}
	}
	return false, 0
}

func (g *GameState) ToJson() (string, error) {
	transportState := ConvertToTransport(g)
	b, err := json.Marshal(transportState)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func intAbs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func (tiereak TieBreak) Error() string {
	return "We have a tiebreak! TODO ME"
}
