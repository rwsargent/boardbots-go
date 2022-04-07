package lockitdown

import (
	"fmt"
	"math"
	"sync"

	"github.com/rwsargent/boardbots-go/internal/minimax"
)

type (
	Evaluator func(*GameState, PlayerPosition) int

	MinimaxNode struct {
		GameState    *GameState
		GameMove     *GameMove
		Searcher     PlayerPosition
		Evaluator    Evaluator
		MinimaxValue int
	}
)

var nodePool = sync.Pool{
	New: func() any {
		return &MinimaxNode{}
	},
}

func (n *MinimaxNode) Evaluate() {
	n.MinimaxValue = n.Evaluator(n.GameState, n.Searcher)
}

func (n *MinimaxNode) Children(nodeBuffer []minimax.Node) []minimax.Node {
	moveBuffer := moveBufferPool.Get().(*[]*GameMove)
	defer moveBufferPool.Put(moveBuffer)

	nextMoves := n.GameState.PossibleMoves((*moveBuffer)[:0])

	for _, nextMove := range nextMoves {
		node := nodePool.Get().(*MinimaxNode)
		node.GameState = n.GameState
		node.GameMove = nextMove
		node.Searcher = n.Searcher
		node.Evaluator = n.Evaluator
		nodeBuffer = append(nodeBuffer, node)
	}
	return nodeBuffer
}

func (n *MinimaxNode) ShouldMaximize() bool {
	return n.Searcher == n.GameState.PlayerTurn
}

func (n *MinimaxNode) Move() {
	err := n.GameState.Move(n.GameMove)
	if err != nil {
		panic(err)
	}
}

func (n *MinimaxNode) Undo() {
	n.GameState.Undo(n.GameMove)
}

func (n *MinimaxNode) Score() int {
	return n.MinimaxValue
}

func (n *MinimaxNode) SetScore(score int) {
	n.MinimaxValue = score
}

func (n *MinimaxNode) Release() {
	movePool.Put(n.GameMove)
	nodePool.Put(n)
}

func MinimaxWithIterator(node *MinimaxNode, depth int) MinimaxNode {
	it := NewMoveIterator(node.GameState)
	if depth == 0 || !it.Next() {
		node.Evaluate()
		return *node
	}

	var best, child = MinimaxNode{}, MinimaxNode{
		GameState: node.GameState,
		Evaluator: node.Evaluator,
		Searcher:  node.Searcher,
	}
	var bestMove GameMove

	var comparator func(int, int) bool
	if node.ShouldMaximize() {
		best.SetScore(math.MinInt)
		comparator = gt
	} else {
		best.SetScore(math.MaxInt)
		comparator = lt
	}

	for it.Next() {
		child.GameMove = it.Get()

		if child.GameMove.Mover == nil {
			panic(fmt.Sprintf("depth: %d, parent: %+v", depth, node))
		}

		child.Move()
		childsBest := MinimaxWithIterator(&child, depth-1)
		if comparator(childsBest.Score(), best.Score()) {
			best = child
			best.SetScore(childsBest.Score())

			// Copy the move, so the iterator doesn't blow it out.
			bestMove = *child.GameMove
			best.GameMove = &bestMove
		}
		child.Undo()

		if child.GameMove.Mover != best.GameMove.Mover {
			ReleaseMover(child.GameMove.Mover)
		}
	}
	return best
}

func gt(a int, b int) bool {
	return a > b
}

func lt(a int, b int) bool {
	return a < b
}
