package lockitdown

import (
	"errors"
	"fmt"
)

type (
	Mover interface {
		Move(*GameState, PlayerPosition) error
		Undo(*GameState, PlayerPosition) error
		ToTransport() BoardbotsMove
	}

	GameMove struct {
		Player PlayerPosition
		Mover
	}

	BoardbotsMove struct {
		Position Pair `json:"pos"`
		Action   any  `json:"action"`
	}

	TurnRobot struct {
		Robot     Pair
		Direction TurnDirection
	}

	InnerTurnRobotT struct {
		Side string `json:"side"`
	}
	TurnRobotT struct {
		Turn InnerTurnRobotT `json:"Turn"`
	}

	PlaceRobot struct {
		Hex, Direction Pair
	}

	InnerPlaceRobotT struct {
		Dir Pair `json:"dir"`
	}

	PlaceRobotT struct {
		PlaceRobot InnerPlaceRobotT
	}

	AdvanceRobot struct {
		Robot Pair
	}
)

func NewMove(m Mover, p PlayerPosition) *GameMove {
	return &GameMove{
		Player: p,
		Mover:  m,
	}
}
func (m *GameMove) Move(state *GameState) error {
	return m.Mover.Move(state, m.Player)
}
func (m *GameMove) Undo(state *GameState) error {
	err := m.Mover.Undo(state, m.Player)
	return err
}

func (m *AdvanceRobot) Move(game *GameState, player PlayerPosition) error {
	robot, found := game.Robots[m.Robot]
	if !found {
		return fmt.Errorf("no robot at location %v", m.Robot)
	}
	if robot.IsLockedDown {
		return errors.New("cannot advance, robot is locked down")
	}
	if robot.Player != player {
		return fmt.Errorf("cannot move %s, it belongs to Player %d", m.Robot.String(), robot.Player)
	}
	robot.Position.Plus(robot.Direction)
	if _, ok := game.Robots[robot.Position]; ok {
		// Undo move
		robot.Position.Minus((robot.Direction))
		return errors.New("cannot advance, another bot in the way")
	}
	delete(game.Robots, m.Robot)
	game.Robots[robot.Position] = robot

	m.Robot = robot.Position
	game.MovesThisTurn -= 1
	return nil
}

func (m *AdvanceRobot) Undo(game *GameState, player PlayerPosition) error {
	robot, found := game.Robots[m.Robot]
	if !found {
		json, _ := game.ToJson()
		panic(fmt.Sprintf("Undoing %v, P%d, with gamestate %s", m, player, json))
	}
	delete(game.Robots, m.Robot)

	robot.Position.Minus(robot.Direction)
	game.Robots[robot.Position] = robot
	m.Robot = robot.Position

	return nil
}

func (m AdvanceRobot) ToTransport() BoardbotsMove {
	return BoardbotsMove{
		Position: m.Robot,
		Action:   "Advance",
	}

}

func (m *PlaceRobot) Move(game *GameState, player PlayerPosition) error {
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

func (m *PlaceRobot) Undo(game *GameState, player PlayerPosition) error {
	delete(game.Robots, m.Hex)
	return nil
}

func (m PlaceRobot) ToTransport() BoardbotsMove {
	return BoardbotsMove{
		Position: m.Hex,
		Action: PlaceRobotT{
			PlaceRobot: InnerPlaceRobotT{
				Dir: m.Direction,
			},
		},
	}
}

func (m *TurnRobot) Move(game *GameState, player PlayerPosition) error {
	var robot *Robot
	var found bool
	if robot, found = game.Robots[m.Robot]; !found {
		return fmt.Errorf("cannot find robot %v", m.Robot)
	}
	if robot.Player != player {
		return fmt.Errorf("cannot move %s, it belongs to Player %d", m.Robot.String(), robot.Player)
	}
	robot.Direction.Rotate(m.Direction)

	game.MovesThisTurn -= 1
	return nil
}

func (m *TurnRobot) Undo(game *GameState, player PlayerPosition) error {
	// Left and Right are zero and one, so 1 - <direction> will
	// give the other direction. 1 - <0> = 1; 1 - <1> = 0
	game.Robots[m.Robot].Direction.Rotate(1 - m.Direction)
	return nil
}

func (m TurnRobot) ToTransport() BoardbotsMove {
	var turn string
	if m.Direction == Left {
		turn = "Left"
	} else {
		turn = "Right"
	}

	return BoardbotsMove{
		Position: m.Robot,
		Action: TurnRobotT{
			Turn: InnerTurnRobotT{
				Side: turn,
			},
		},
	}
}
