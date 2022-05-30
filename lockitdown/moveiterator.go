package lockitdown

import "sync"

type (
	MoveIterator struct {
		game        *GameState
		currentMove *GameMove
		moveBuf     [3]GameMove
		moveIdx     int
		edgeIndex   int
		robotIndex  int
		botCache    map[Pair]*Robot
	}
)

var (
	advancePool = sync.Pool{
		New: func() any { return new(AdvanceRobot) },
	}

	turnPool = sync.Pool{
		New: func() any { return new(TurnRobot) },
	}

	placePool = sync.Pool{
		New: func() any { return new(PlaceRobot) },
	}
)

func NewMoveIterator(game *GameState) *MoveIterator {
	botCache := make(map[Pair]*Robot)
	for i := 0; i < len(game.Robots); i++ {
		robot := &game.Robots[i]
		botCache[robot.Position] = &game.Robots[i]
	}
	it := &MoveIterator{
		game:        game,
		currentMove: new(GameMove),
		moveBuf:     [3]GameMove{},
		edgeIndex:   0,
		moveIdx:     -1,
		robotIndex:  0,
		botCache:    botCache,
	}
	return it
}

func (it *MoveIterator) Get() *GameMove {
	return it.currentMove
}

func (it *MoveIterator) Next() bool {
	it.findNext()
	return it.currentMove != nil
}

func (it *MoveIterator) findNext() {
	// Check to see if we have any buffered moves
	// already calculated
	for it.moveIdx >= 0 && it.moveIdx < len(it.moveBuf) {
		next := &it.moveBuf[it.moveIdx]
		it.currentMove = next

		it.moveIdx++
		if it.currentMove.Mover != nil {
			return
		}
	}

	// Check all owned bots. robotIndex refers to the next bot
	// to buffer moves for.

	ringSize := it.game.GameDef.Board.HexaBoard.ArenaRadius + 1
	botIdx := -1
	for _, bot := range it.game.Robots {
		if bot.Player == it.game.PlayerTurn && !bot.IsLockedDown {
			botIdx++
			if botIdx == it.robotIndex {
				it.robotIndex++
				it.moveIdx = 0
				// Advance bot
				advancePosition := Pair{
					Q: bot.Position.Q,
					R: bot.Position.R,
				}
				advancePosition.Plus(bot.Direction)

				if _, blocked := it.botCache[advancePosition]; !blocked &&
					inBounds(it.game.GameDef.Board.HexaBoard.ArenaRadius+1, advancePosition) {
					advance := advancePool.Get().(*AdvanceRobot)
					advance.Robot = bot.Position
					it.moveBuf[0].Mover = advance
					it.moveBuf[0].Player = it.game.PlayerTurn
				} else {
					it.moveBuf[0].Mover = nil
				}

				addTurn(&it.moveBuf[1], Left, ringSize, bot, it.game)
				addTurn(&it.moveBuf[2], Right, ringSize, bot, it.game)

				// Find the first non-nil move
				for it.moveIdx < len(it.moveBuf) && it.moveBuf[it.moveIdx].Mover == nil {
					it.moveIdx++
				}
				if it.moveIdx >= len(it.moveBuf) {
					// no moves for this robot
					continue
				}
				it.currentMove = &it.moveBuf[it.moveIdx]
				it.moveIdx++ // advance cursor for prep next
				return
			}
		}
	}

	if it.game.MovesThisTurn == 3 && it.game.playerBotsInCorridor() < 2 {
		edges := edges(it.game.GameDef.Board.HexaBoard.ArenaRadius + 1)
		for it.edgeIndex < len(edges) {
			edge := edges[it.edgeIndex]
			it.edgeIndex++
			if _, blocked := it.botCache[edge.position]; !blocked {
				place := placePool.Get().(*PlaceRobot)
				place.Robot = edge.position
				place.Direction = edge.direction
				it.currentMove.Mover = place
				it.currentMove.Player = it.game.PlayerTurn
				return
			}
		}
	}
	movePool.Put(it.currentMove)
	it.currentMove = nil
	return
}

func addTurn(gameMove *GameMove, direction TurnDirection, arenaRadius int, bot Robot, g *GameState) {
	turn := bot.Direction.Copy()
	turn.Rotate(direction)

	pos := bot.Position.Copy()
	pos.Plus(turn)

	if inBounds(arenaRadius, pos) {
		turn := turnPool.Get().(*TurnRobot)
		turn.Direction = direction
		turn.Robot = bot.Position
		gameMove.Mover = turn
		gameMove.Player = g.PlayerTurn
	} else {
		gameMove.Mover = nil
	}
}

func ReleaseMover(mover Mover) {
	switch v := mover.(type) {
	case *AdvanceRobot:
		advancePool.Put(v)
	case *TurnRobot:
		turnPool.Put(v)
	case *PlaceRobot:
		placePool.Put(v)
	}
}
