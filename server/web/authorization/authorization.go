package authorization

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rwsargent/boardbots/server/users"
	"net/http"
	"strings"
)

type (
	Authenticator interface {
		ValidCredentials(username, password string) bool
		ValidateToken(token string) bool
	}

	credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Response struct {
		Token string `json:"token"`
	}

	handler struct {
		authenticator Authenticator
		users users.UserFinder
	}
)

func RegisterRoutes(e *echo.Echo, authenticator Authenticator, users users.UserFinder) {
	handler := handler{
		authenticator: authenticator,
		users : users,
	}
	// Global authorization handlers
	e.Use(handler.AuthorizeClient)
	e.Use(handler.AuthorizeUser)

	g := e.Group("/auth")
	g.POST("/login", handler.Login)
	g.POST("/validate", handler.ValidateSession)
}

func (handler *handler) Login(context echo.Context) error {
	credentials := new(credentials)
	if err := context.Bind(credentials); err != nil {
		return err
	}
	context.Logger().Info(fmt.Sprintf("creds: %v", credentials))
	valid := handler.authenticator.ValidCredentials(credentials.Username, credentials.Password)
	if !valid {
		return context.String(http.StatusTeapot, "unauthorized")
	}
	return context.JSON(http.StatusOK, Response{
		Token: handler.users.FindByName(credentials.Username).Token,
	})
}

func (handler *handler) AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		if strings.HasPrefix(context.Request().URL.Path, "/auth/login") {
			// don't authorize the login - how else will people log in!
			return next(context)
		}
		authorization := context.Request().Header.Get("Authorization")
		parts := strings.Split(authorization, " ")
		if parts[0] != "Bearer" {
			return context.String(http.StatusBadRequest, "unexpected authorization type")
		}
		if !handler.authenticator.ValidateToken(parts[1]) {
			return context.String(http.StatusUnauthorized, "invalid token")
		}
		context.Set("user", handler.users.FindByToken(parts[1]))
		return next(context)
	}
}

func (handler *handler) ValidateSession (ctx echo.Context) error {
	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		return ctx.String(http.StatusUnauthorized, "no session found")
	}
	if !handler.authenticator.ValidateToken(token) {
		return ctx.String(http.StatusUnauthorized, "invalid session")
	}
	return ctx.String(http.StatusOK, "ok")
}

func (handler *handler) AuthorizeClient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		clientId := context.Request().Header.Get("X-Client-Id")
		if clientId != "IN-DEVELOPMENT" {
			context.Logger().Info(fmt.Sprintf("illegal client id: %s", clientId))
			return context.String(http.StatusUnauthorized, "invalid client")
		}
		return next(context)
	}
}
