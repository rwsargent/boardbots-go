package manager

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	q "github.com/rwsargent/boardbots/server/quoridor"
	"math"
	"sync"
)

type GameManager interface {
   NewGame(gameId uuid.UUID, name string) (q.Game, error)
   AddPlayer(gameId, playerId uuid.UUID, playerName string) (q.Game, error)
   GetGame(gameId uuid.UUID) (q.Game, error)
   MovePawn(gameId, playerId uuid.UUID, position q.Position) (q.Game, error)
   PlaceBarrier(gameId, playerId uuid.UUID, position q.Position) (q.Game, error)
}

type InMemoryGameManager struct {
	games map[uuid.UUID] *q.Game
	locks map[uint8] *sync.RWMutex
}

func NewInMemoryGameManager() *InMemoryGameManager {
	return &InMemoryGameManager{
		games : make(map[uuid.UUID]*q.Game),
		locks : initLocks(),
	}
}

func initLocks() map[uint8] *sync.RWMutex {
	locks := make(map[uint8] *sync.RWMutex)
	for i := uint8(0); i < math.MaxUint8; i++ {
		locks[i] = &sync.RWMutex{}
	}
	return locks
}

func lockId(gameId uuid.UUID) uint8{
	return (uint8)(gameId.ID() & 0xFF)
}

func (manager *InMemoryGameManager) NewGame(gameId uuid.UUID, name string) (q.Game, error) {
	manager.locks[lockId(gameId)].Lock()
	defer manager.locks[lockId(gameId)].Unlock()

	if _, present := manager.games[gameId]; present {
		return q.Game{}, errors.New(fmt.Sprintf("game with id %s already exists", gameId.String()))
	}
	game, err := q.NewGame(gameId, name)
	if err != nil {
		return q.Game{}, err
	}
	manager.games[game.Id] = game
	return game.Copy(), nil
}

func (manager *InMemoryGameManager) AddPlayer(gameId, playerId uuid.UUID, playerName string) (q.Game, error) {
	manager.locks[lockId(gameId)].Lock()
	defer manager.locks[lockId(gameId)].Unlock()

	game, present := manager.games[gameId]
	if !present {
		return q.Game{}, errors.New(fmt.Sprintf("game with id %s already exists", gameId.String()))
	}
	_, err := game.AddPlayer(playerId, playerName)
	if err != nil {
		return q.Game{}, err
	}
	return game.Copy(), nil
}

func (manager *InMemoryGameManager) GetGame(gameId uuid.UUID) (q.Game, error) {
	manager.locks[lockId(gameId)].RLock()
	defer manager.locks[lockId(gameId)].RUnlock()

	game, present := manager.games[gameId]
	if !present {
		return q.Game{}, errors.New(fmt.Sprintf("game with id %s already exists", gameId.String()))
	}

	return game.Copy(), nil
}

func (manager *InMemoryGameManager) MovePawn(gameId, playerId uuid.UUID, position q.Position) (q.Game, error) {
	manager.locks[lockId(gameId)].Lock()
	defer manager.locks[lockId(gameId)].Unlock()
	game, present := manager.games[gameId]
	if !present {
		return q.Game{}, errors.New(fmt.Sprintf("game with id %s already exists", gameId.String()))
	}
	player := game.GetPlayerPosition(playerId)
	if player == q.InvalidPlayer {
		return q.Game{}, errors.New("player does not exist on this game")
	}

	err := game.MovePawn(position, player)
	if err != nil {
		return q.Game{}, err
	}
	return game.Copy(), nil
}

func (manager *InMemoryGameManager) PlaceBarrier(gameId, playerId uuid.UUID, position q.Position) (q.Game, error) {
	manager.locks[lockId(gameId)].Lock()
	defer manager.locks[lockId(gameId)].Unlock()
	game, present := manager.games[gameId]
	if !present {
		return q.Game{}, errors.New(fmt.Sprintf("game with id %s already exists", gameId.String()))
	}
	player := game.GetPlayerPosition(playerId)
	if player == q.InvalidPlayer {
		return q.Game{}, errors.New("player does not exist on this game")
	}

	err := game.PlaceBarrier(position, player)
	if err != nil {
		return q.Game{}, err
	}
	return game.Copy(), nil
}