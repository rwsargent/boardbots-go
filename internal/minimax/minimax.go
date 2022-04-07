package minimax

import (
	"math"
	"sync"
)

type (
	Node interface {
		Evaluate()
		Children([]Node) []Node
		Release()
		Move()
		Undo()
		ShouldMaximize() bool
		Score() int
		SetScore(int)
	}
)

var nodePool = sync.Pool{
	New: func() any {
		buf := make([]Node, 0, 128)
		return &buf
	},
}

func Minimax(node Node, depth int) Node {
	nodeBuffer := *(nodePool.Get().(*[]Node))
	children := node.Children(nodeBuffer[:0])
	defer nodePool.Put(&children)
	if depth == 0 || len(children) == 0 {
		node.Evaluate()
		return node
	}

	var best Node
	if node.ShouldMaximize() {
		best = children[0]
		best.SetScore(math.MinInt)
		for _, child := range children {
			child.Move()
			childsBest := Minimax(child, depth-1)
			if childsBest.Score() > best.Score() {
				best.Release()
				best = child
				best.SetScore(childsBest.Score())
			}
			child.Undo()
			if child != best {
				child.Release()
			}
		}
	} else {
		best = children[0]
		best.SetScore(math.MaxInt)
		for _, child := range children {
			child.Move()
			childsBest := Minimax(child, depth-1)
			if childsBest.Score() < best.Score() {
				best.Release()
				best = child
				best.SetScore(childsBest.Score())
			}
			child.Undo()
			if child != best {
				child.Release()
			}
		}
	}

	defer best.Release()
	return best
}
