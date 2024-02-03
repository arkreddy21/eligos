package main

import (
	"context"
	"fmt"
	"github.com/arkreddy21/eligos/internal/http"
	"github.com/arkreddy21/eligos/internal/postgres"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	app := newApp()
	app.run()
	fmt.Println("app started")

	// Wait for interrupt signal.
	<-ctx.Done()
	fmt.Println("\napp stopped")
	// Clean up program.
	if err := app.close(); err != nil {
		log.Fatal("failed to shutdown gracefully: ", err)
	}
}

type App struct {
	DB         *postgres.DB
	HTTPServer *http.Server
}

func newApp() *App {
	return &App{
		HTTPServer: http.NewServer(),
	}
}

func (app *App) run() {
	app.DB = postgres.NewDB()
	app.HTTPServer.UserService = postgres.NewUserService(app.DB)
	app.HTTPServer.SpaceService = postgres.NewSpaceService(app.DB)
	app.HTTPServer.MessageService = postgres.NewMessageService(app.DB)
	app.HTTPServer.Open()
}

func (app *App) close() error {
	err := app.HTTPServer.Close()
	if err != nil {
		return err
	}
	if app.DB != nil {
		app.DB.Close()
	}
	return nil
}
