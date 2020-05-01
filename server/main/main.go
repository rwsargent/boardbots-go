package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rwsargent/boardbots/server/server"
	"github.com/rwsargent/boardbots/server/web/authorization"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler : func(e echo.Context, reqBody []byte, resBody []byte) {
			fmt.Printf(fmt.Sprintf("body: %s\n", string(reqBody)))
		},
	}))

	server := server.NewServer()

	authorization.RegisterRoutes(e, server.Authenticator, server.UserFinder)

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	gracefulShutdown(e)
}

func gracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
