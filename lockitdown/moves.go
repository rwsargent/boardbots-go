package lockitdown

import (
	"errors"
	"fmt"
)

type (
	Move interface {
		Move(*GameState, PlayerPosition) error
		Undo(*GameState, PlayerPosition)
	}

	TurnRobot struct {
		Robot     Pair
		Direction TurnDirection
	}

	PlaceRobot struct {
		Hex, Direction Pair
	}

	PlaceRobotT struct {
		PlaceRobot struct {
			Dir Pair `json:"dir"`
		}
	}

	AdvanceRobot struct {
		Robot Pair
	}
)

func (m AdvanceRobot) Move(game *GameState, player PlayerPosition) error {
	robot := game.Robots[m.Robot]
	robot.Position.Plus(robot.Direction)
	if _, ok := game.Robots[robot.Position]; ok {
		// Undo move
		robot.Position.Minus((robot.Direction))
		return errors.New("cannot advance, another bot the way")
	}
	delete(game.Robots, m.Robot)
	game.Robots[robot.Position] = robot

	game.MovesThisTurn -= 1
	return nil
}

func (m AdvanceRobot) Undo(game *GameState, player PlayerPosition) {
	robot := game.Robots[m.Robot]
	delete(game.Robots, m.Robot)

	robot.Position.Minus(robot.Direction)
	game.Robots[robot.Position] = robot

	game.MovesThisTurn += 1
}

func (m PlaceRobot) Move(game *GameState, player PlayerPosition) error {
	if game.MovesThisTurn != 3 {
		return errors.New("can only place a robot on your first action of the turn")
	}
	if !game.isCorridor(m.Hex) {
		return errors.New("must place robot in corridor")
	}

	robotsInCorridor := 0
	for _, robot := range game.Robots {
		if robot.Player == player {
			if game.isCorridor(robot.Position) {
				robotsInCorridor++
			}
		}
	}
	if robotsInCorridor > 1 {
		return errors.New("can only have two robots in the corridor at a time")
	}

	game.Robots[m.Hex] = &Robot{
		Position:      m.Hex,
		Direction:     m.Direction,
		IsBeamEnabled: true,
		IsLockedDown:  false,
		Player:        player,
	}

	game.MovesThisTurn = 0
	game.Players[player].PlacedRobots += 1

	return nil
}

func (m PlaceRobot) Undo(game *GameState, player PlayerPosition) {
	delete(game.Robots, m.Hex)

	game.Players[player].PlacedRobots -= 1

	game.MovesThisTurn = 3
}

func (m TurnRobot) Move(game *GameState, player PlayerPosition) error {
	var robot *Robot
	var found bool
	if robot, found = game.Robots[m.Robot]; !found {
		return fmt.Errorf("cannot find robot %v", *robot)
	}
	robot.Direction.Rotate(m.Direction)

	game.MovesThisTurn -= 1
	return nil
}

func (m TurnRobot) Undo(game *GameState, player PlayerPosition) {
	// Left and Right are zero and one, so 1 - <direction> will
	// give the other direction. 1 - <0> = 1; 1 - <1> = 0
	game.Robots[m.Robot].Direction.Rotate(1 - m.Direction)

	game.MovesThisTurn += 1
}
