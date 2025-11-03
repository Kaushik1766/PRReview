package app

import (
	"net/http"
)

// basePath = "/api/v1"

var routes map[string]func(w http.ResponseWriter, r *http.Request)

func (app *App) registerRoutes() {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		"/": app.PRHandler.GetPREvent,
	}

	for route, handler := range routes {
		// pathArr := strings.Split(route, " ")
		// method := pathArr[0]
		// path := basePath + pathArr[1]
		app.apiMux.HandleFunc(route, handler)
	}
}
