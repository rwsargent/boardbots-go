// An implmentation of boardbots.dev LockItDown game. Users can use package
// to make and undo moves on an internal state of the board, and then apply
// those moves to a game hosted on boardbots.dev.
package lockitdown

import (
	"encoding/json"
	"fmt"
	"sync"
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
		currentMove      *GameMove
		saveStack        []SaveState
	}

	TurnDirection int

	TieBreak struct {
		Robots []*Robot
		State  string
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

	Cardinals = []Pair{E, SE, SW, W, NW, NE}

	moveBufferPool = sync.Pool{
		New: func() any {
			return make([]*GameMove, 0, 128)
		},
	}

	movePool = sync.Pool{
		New: func() any {
			return new(GameMove)
		},
	}
)

func (p *Pair) Plus(that Pair) {
	p.Q += that.Q
	p.R += that.R
}

func (p *Pair) Minus(that Pair) {
	p.Q -= that.Q
	p.R -= that.R
}

func (p Pair) String() string {
	return fmt.Sprintf("{%d, %d}", p.Q, p.R)
}

func (p Pair) S() int {
	return -p.Q - p.R
}

func (p Pair) Copy() Pair {
	return p
}

func (p Pair) Dist() int {
	return (intAbs(p.Q) + intAbs(p.R) + intAbs(p.S())) / 2
}

func (r *Robot) Disable() {
	r.IsBeamEnabled = false
	r.IsLockedDown = true
}

func (r *Robot) Enable() {
	r.IsBeamEnabled = true
	r.IsLockedDown = false
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
		saveStack:        make([]SaveState, 0),
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

func (game *GameState) Move(move *GameMove) error {
	game.currentMove = move
	if move.Player != PlayerPosition(game.PlayerTurn) {
		return fmt.Errorf("wrong player, expected %d, was %d", game.PlayerTurn, move.Player)
	}

	game.saveState()
	err := move.Move(game)

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

func (game *GameState) Undo(move *GameMove) error {
	game.currentMove = move

	err := move.Undo(game)
	game.revertState()

	if err != nil {
		return err
	}

	err = game.resolveMove()

	if err != nil {
		return err
	}

	return nil
}

func (game *GameState) resolveMove() error {
	for resolved := false; !resolved; {
		targeted := game.taretedRobots()

		if tiebreaks := game.checkForTieBreaks(targeted); len(tiebreaks) > 0 {
			game.RequiresTieBreak = true
			json, _ := game.ToJson()
			return TieBreak{
				Robots: tiebreaks,
				State:  json,
			}
		}

		resolved = game.updateLockedRobots(targeted)
	}
	return nil
}

func (game *GameState) updateLockedRobots(targeted map[Pair][]*Robot) bool {
	resolved := true
	// clear state on bots
	for _, robot := range game.Robots {
		robot.Enable()
	}

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
				locked.Disable()
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
	board := game.GameDef.Board.HexaBoard
	corridor := board.ArenaRadius + 1
	return pair.Dist() == corridor
}

func (game *GameState) shutdownRobot(hex Pair, attackers []*Robot) {
	game.saveRobot(game.Robots[hex])
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

// as an optimization, PossibleMoves takes a buffer to avoid allocating every call.
func (g *GameState) PossibleMoves(buf []*GameMove) []*GameMove {
	moves := buf
	if g.MovesThisTurn == 3 && g.playerBotsInCorridor() < 2 {
		edges := edges(g.GameDef.Board.HexaBoard.ArenaRadius + 1)
		for _, edge := range edges {
			if _, found := g.Robots[edge.position]; !found {
				m := movePool.Get().(*GameMove)
				m.Mover = &PlaceRobot{
					Robot:     edge.position,
					Direction: edge.direction,
				}
				m.Player = g.PlayerTurn
				moves = append(moves, m)
			}
		}
	}
	return moves
}

func (g *GameState) playerBotsInCorridor() int {
	corridorBots := 0
	for _, bot := range g.Robots {
		if bot.Player == g.PlayerTurn {
			if g.isCorridor(bot.Position) {
				corridorBots++
			}
		}
	}
	return corridorBots
}

func inBounds(size int, position Pair) bool {
	return position.Dist() <= size
}

func (state *GameState) saveState() {
	save := SaveState{
		players:        []Player{},
		shutdownRobots: []*Robot{},
		movesThisTurn:  0,
		player:         0,
	}
	for _, player := range state.Players {
		save.players = append(save.players, *player)
	}
	save.player = state.PlayerTurn
	save.movesThisTurn = state.MovesThisTurn
	state.saveStack = append(state.saveStack, save)
}

func (state *GameState) saveRobot(robot *Robot) {
	save := &state.saveStack[len(state.saveStack)-1]
	save.shutdownRobots = append(save.shutdownRobots, robot)
}

func (state *GameState) revertState() {
	save := state.saveStack[len(state.saveStack)-1]
	for _, bot := range save.shutdownRobots {
		state.Robots[bot.Position] = bot
	}

	for i, player := range save.players {
		state.Players[i].PlacedRobots = player.PlacedRobots
		state.Players[i].Points = player.Points
	}

	state.PlayerTurn = save.player
	state.MovesThisTurn = save.movesThisTurn

	state.saveStack = state.saveStack[:len(state.saveStack)-1]
}

func intAbs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func (tiereak TieBreak) Error() string {
	return fmt.Sprintf("Tiebreak: %s", tiereak.State)
}
