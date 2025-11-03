package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"Kaushik1766/PRReview/internal/handlers/webhooks/pr"
	pranalysis "Kaushik1766/PRReview/internal/services/pr-analysis"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/google/go-github/github"
)

type App struct {
	apiMux *http.ServeMux

	PRHandler *pr.PRHandler

	PRAnalyzer pranalysis.PRAnalyzer
}

func NewApp() (*App, error) {
	app := App{}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	githubClient := github.NewClient(nil)
	prAnalyzer, err := pranalysis.NewPRAnalysis(client)
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
