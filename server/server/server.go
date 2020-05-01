package server

import (
	"github.com/rwsargent/boardbots/server/manager"
	"github.com/rwsargent/boardbots/server/users"
	"github.com/rwsargent/boardbots/server/web/authorization"
)

type (
	Server struct {
		Manager manager.GameManager
		Authenticator authorization.Authenticator
		UserFinder users.UserFinder
	}
)

func NewServer() *Server {
	users := users.NewDevUsers("users.json")
	return &Server{
		Manager: manager.NewInMemoryGameManager(),
		Authenticator: users,
		UserFinder: users,
	}
}
