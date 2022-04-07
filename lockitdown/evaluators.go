package lockitdown

var attackAxes = map[Pair]func(Pair, Pair) bool{
	W: func(attacker Pair, bot Pair) bool {
		return attacker.R == bot.R && attacker.Q > bot.Q
	},
	NW: func(attacker Pair, bot Pair) bool {
		return attacker.Q == bot.Q && attacker.R > bot.R
	},
	NE: func(attacker Pair, bot Pair) bool {
		return attacker.S() == bot.S() && attacker.R > bot.R
	},
	E: func(attacker Pair, bot Pair) bool {
		return attacker.R == bot.R && attacker.Q < bot.Q
	},
	SE: func(attacker Pair, bot Pair) bool {
		return attacker.Q == bot.Q && attacker.R < bot.R
	},
	SW: func(attacker Pair, bot Pair) bool {
		return attacker.S() == bot.S() && attacker.Q > bot.Q
	},
}

func ScoreGameState(game *GameState, player PlayerPosition) int {
	score := 0

	for _, robot := range game.Robots {
		botScore := scoreRobot(robot, game)
		if robot.Player == player {
			score += botScore
		} else {
			score -= botScore
		}
	}

	score += playerScore(game.Players[game.PlayerTurn])

	return score
}

func playerScore(player *Player) int {
	return player.Points * 30
}

func scoreRobot(robot *Robot, game *GameState) int {
	botScore := 0

	if robot.IsLockedDown {
		botScore -= 100
	}

	botScore += scoreBotPosition(robot, game)

	return botScore
}

func scoreBotPosition(robot *Robot, game *GameState) int {
	score := 0

	// Boost score for hexes the are on the axis. This will help
	// prioritize corners.
	if robot.Position.Q == 0 || robot.Position.R == 0 || robot.Position.S() == 0 {
		score += 10
	}

	// Count number of attackable paths to bot position.
	// This will help prioiritize edges.
	if !game.isCorridor(robot.Position) {
		cursor := robot.Position.Copy()
		cursor.Plus(NW)

		attackableHexes := 0
		for _, dir := range Cardinals {
			if !game.isCorridor(cursor) {
				attackableHexes++
			}
			cursor.Plus(dir)
		}
		score -= attackableHexes
	}

	// Encourage enemy bots in range
	for _, bot := range game.Robots {
		if robot.Player != bot.Player {
			if attackAxes[robot.Direction](robot.Position, bot.Position) {
				score += 20
			}
		}
	}

	return score
}
