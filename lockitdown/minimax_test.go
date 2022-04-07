package lockitdown

import (
	"fmt"
	"testing"

	"github.com/rwsargent/boardbots-go/internal/minimax"
	"github.com/stretchr/testify/assert"
)

func TestDepthOfOne(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)

	root := MinimaxNode{
		GameState: game,
		GameMove:  nil,
		Searcher:  0,
		Evaluator: ScoreGameState,
	}

	score := minimax.Minimax(&root, 1)
	move, _ := score.(*MinimaxNode)
	fmt.Printf("%T: %+v\n", move.GameMove, move.GameMove)
}

func TestDepthOfTwo(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)

	root := MinimaxNode{
		GameState: game,
		GameMove:  nil,
		Searcher:  0,
		Evaluator: ScoreGameState,
	}

	score := minimax.Minimax(&root, 2)
	move, _ := score.(*MinimaxNode)
	fmt.Printf("%T: %+v\n", move.GameMove, move.GameMove)
}

func TestDepthOf3(t *testing.T) {
	game := NewGame(TwoPlayerGameDef)
	originalJson, _ := game.ToJson()
	root := MinimaxNode{
		GameState: game,
		GameMove:  nil,
		Searcher:  0,
		Evaluator: ScoreGameState,
	}

	score := minimax.Minimax(&root, 3)
	move, _ := score.(*MinimaxNode)
	fmt.Printf("%T: %+v\n", move.GameMove, move.GameMove)
	searchedJson, _ := game.ToJson()
	assert.Equal(t, originalJson, searchedJson)
}

func BenchmarkMinimax3(b *testing.B) {

	game := NewGame(TwoPlayerGameDef)
	root := MinimaxNode{
		GameState: game,
		GameMove:  nil,
		Searcher:  0,
		Evaluator: ScoreGameState,
	}

	for n := 0; n < b.N; n++ {
		_ = minimax.Minimax(&root, 3)
	}

}

func BenchmarkMinimaxWithIterator(b *testing.B) {
	game := NewGame(TwoPlayerGameDef)
	root := MinimaxNode{
		GameState: game,
		GameMove:  nil,
		Searcher:  0,
		Evaluator: ScoreGameState,
	}

	for n := 0; n < b.N; n++ {
		_ = MinimaxWithIterator(&root, 3)
	}
}
