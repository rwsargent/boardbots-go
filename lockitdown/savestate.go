package lockitdown

type SaveState struct {
	players        []Player
	shutdownRobots []*Robot
	movesThisTurn  int
	player         PlayerPosition
}
