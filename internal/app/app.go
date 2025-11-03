package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"Kaushik1766/PRReview/internal/handlers/webhooks/pr"
	pranalysis "Kaushik1766/PRReview/internal/services/pr-analysis"

	"github.com/google/go-github/github"
	"google.golang.org/genai"
)

type App struct {
	apiMux *http.ServeMux

	PRHandler *pr.PRHandler

	PRAnalyzer pranalysis.PRAnalyzer
}

func NewApp() (*App, error) {

	app := App{}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	aiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}

	githubClient := github.NewClient(nil)
	prAnalyzer, err := pranalysis.NewPRAnalysis(aiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create PR analyzer: %w", err)
	}

	app.PRAnalyzer = *prAnalyzer

	app.PRHandler = pr.NewPRHandler(prAnalyzer, githubClient)

	app.apiMux = http.NewServeMux()

	app.registerRoutes()
	return &app, nil
}

func (app *App) Run() {
	fmt.Println("Server started at localhost:3000")
	http.ListenAndServe("localhost:3000", app.apiMux)
}
