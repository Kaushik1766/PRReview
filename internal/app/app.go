package app

import (
	"fmt"
	"net/http"

	"Kaushik1766/PRReview/internal/handlers/webhooks/pr"
)

type App struct {
	apiMux *http.ServeMux

	PRHandler *pr.PRHandler
}

func NewApp() (*App, error) {
	app := App{
		apiMux:    http.NewServeMux(),
		PRHandler: pr.NewPRHandler(),
	}

	app.registerRoutes()
	return &app, nil
}

func (app *App) Run() {
	fmt.Println("Server started at localhost:3000")
	http.ListenAndServe("localhost:3000", app.apiMux)
}
