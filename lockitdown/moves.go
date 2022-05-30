package lockitdown

import (
	"errors"
	"fmt"
)

type (
	Mover interface {
		Move(*GameState, PlayerPosition) error
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
		Robot, Direction Pair
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

func (m *AdvanceRobot) Move(game *GameState, player PlayerPosition) error {
	robot := game.RobotAt(m.Robot)
	if robot == nil {
		return fmt.Errorf("no robot at location %v", m.Robot)
	}
	if robot.IsLockedDown {
		return errors.New("cannot advance, robot is locked down")
	}
	if robot.Player != player {
		return fmt.Errorf("cannot move %s, it belongs to Player %d", m.Robot.String(), robot.Player)
	}
	advanceSpot := robot.Position.Copy()
	advanceSpot.Plus(robot.Direction)
	if block := game.RobotAt(advanceSpot); block != nil {
		// Undo move
		return errors.New("cannot advance, another bot in the way")
	}
	robot.Position.Plus(robot.Direction)

	game.MovesThisTurn -= 1

	// // Evaluate state before turning on beam
	game.resolveMove()

	robot.IsBeamEnabled = !game.isCorridor(robot.Position) && !robot.IsLockedDown
	return nil
}

func (m AdvanceRobot) ToTransport() BoardbotsMove {
	return BoardbotsMove{
		Position: m.Robot,
		Action:   "Advance",
	}

}

func (m AdvanceRobot) String() string {
	return fmt.Sprintf("Move %s", m.Robot.String())
}

func (m *PlaceRobot) Move(game *GameState, player PlayerPosition) error {
	if game.MovesThisTurn != 3 {
		return errors.New("can only place a robot on your first action of the turn")
	}
	if !game.isCorridor(m.Robot) {
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

	game.Robots = append(game.Robots, Robot{
		Position:      m.Robot,
		Direction:     m.Direction,
		IsBeamEnabled: true,
		IsLockedDown:  false,
		Player:        player,
	})

	game.MovesThisTurn = 0
	game.Players[player].PlacedRobots += 1

	return nil
}

func (m PlaceRobot) String() string {
	return fmt.Sprintf("Place %s: dir: %s", m.Robot.String(), m.Direction.String())
}

func (m PlaceRobot) ToTransport() BoardbotsMove {
	return BoardbotsMove{
		Position: m.Robot,
		Action: PlaceRobotT{
			PlaceRobot: InnerPlaceRobotT{
				Dir: m.Direction,
			},
		},
	}
}

func (m *TurnRobot) Move(game *GameState, player PlayerPosition) error {
	var robot *Robot
	if robot = game.RobotAt(m.Robot); robot == nil {
		return fmt.Errorf("cannot find robot %v", m.Robot)
	}
	if robot.Player != player {
		return fmt.Errorf("cannot move %s, it belongs to Player %d", m.Robot.String(), robot.Player)
	}
	robot.IsBeamEnabled = false
	game.activeBot = robot
	robot.Direction.Rotate(m.Direction)
	game.resolveMove()

	robot.IsBeamEnabled = !robot.IsLockedDown
	game.activeBot = nil

	game.MovesThisTurn -= 1
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

func (m TurnRobot) String() string {
	var turn string
	if m.Direction == Left {
		turn = "Left"
	} else {
		turn = "Right"
	}
	return fmt.Sprintf("Turn %s %s", m.Robot.String(), turn)
}
