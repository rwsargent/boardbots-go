package lockitdown

type SaveState struct {
	players       []Player
	bots          []Robot
	movesThisTurn int
	player        PlayerPosition
}
