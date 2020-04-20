package quoridor

import (
	"container/heap"
	"math"
)

/**
 * A struct for the priority queue, holds the Position and priority of the Node in the Board graph
 */
type PQNode struct {
	position Position
	prev     *PQNode

	distance, priority int
}

func absInt(val int) int {
	y := val >> 31
	return (val ^ y) - y
}

// Calculates the priority of the node as sum of distance to goal + path so far
func (node *PQNode) setPriority(goal Position) {
	if goal.Y < 0 {
		node.priority = absInt(goal.X-node.position.X) + node.distance
	} else if goal.X < 0 {
		node.priority = absInt(goal.Y-node.position.Y) + node.distance
	}
}

/**
 * Priority Queue methods for the heap
 */
type PriorityQueue []*PQNode

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(left, right int) bool {
	return pq[left].priority < pq[right].priority
}
func (pq PriorityQueue) Swap(left, right int) {
	pq[left], pq[right] = pq[right], pq[left]
}

func (pq *PriorityQueue) Push(item interface{}) {
	*pq = append(*pq, item.(*PQNode))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

/**
 * Simple A* Algorithm, implementing a best first search by calculating distance to goal plus
 * the length of path from the pawn.
 */
func (game *Game) FindPath(start, goal Position) []Position {
	pq := &PriorityQueue{}
	heap.Init(pq)
	//breath-first / best-first
	pq.Push(newPQNode(start, goal))
	visited := make(map[Position]bool)
	for pq.Len() != 0 {
		node := heap.Pop(pq).(*PQNode)
		if reachedGoal(node.position, goal) {
			return buildPath(node)
		}
		if _, seen := visited[node.position]; seen {
			continue
		}
		visited[node.position] = true
		neighbors := getReachableNeighbors(node, goal, game)
		for _, neighbor := range neighbors {
			heap.Push(pq, neighbor)
		}
	}
	return nil
}

func reachedGoal(current Position, goal Position) bool {
	return current.Y == goal.Y || current.X == goal.X
}

func getReachableNeighbors(node *PQNode, goal Position, game *Game) []*PQNode {
	neighbors := make([]*PQNode, 0, 4)
	for _, dir := range directions {
		neighborPositions := game.Board.getValidMoveByDirection(node.position, dir)
		if neighborPositions != nil {
			for _, neighborPos := range neighborPositions {
				node := &PQNode{position: neighborPos, prev: node, priority: math.MaxInt32, distance: node.distance + 1}
				node.setPriority(goal)
				neighbors = append(neighbors, node)
			}
		}
	}
	return neighbors
}

func buildPath(node *PQNode) []Position {
	path := make([]Position, 0)
	cursor := node
	for cursor.prev != nil {
		path = append(path, cursor.position)
		cursor = cursor.prev
	}
	//reverse
	for idx := len(path)/2 - 1; idx >= 0; idx-- {
		opp := len(path) - 1 - idx
		path[idx], path[opp] = path[opp], path[idx]
	}
	return path
}

func newPQNode(position Position, goal Position) *PQNode {
	node := &PQNode{
		position: position,
	}
	node.setPriority(goal)
	return node
}
