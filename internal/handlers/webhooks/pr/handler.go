package pr

import (
	"context"
	"fmt"
	"net/http"

	pranalysis "Kaushik1766/PRReview/internal/services/pr-analysis"
	"Kaushik1766/PRReview/utils"

	"github.com/google/go-github/github"
)

type PRHandler struct {
	prAnalyzer   pranalysis.PRAnalyzer
	githubClient *github.Client
}

func NewPRHandler(analyzer pranalysis.PRAnalyzer, githubClient *github.Client) *PRHandler {
	return &PRHandler{
		prAnalyzer:   analyzer,
		githubClient: githubClient,
	}
}

func (handler *PRHandler) GetPREvent(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte("asdf"))
	if err != nil {
		fmt.Println(err)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch event := event.(type) {
	case *github.PushEvent:
		comp, _, err := handler.githubClient.Repositories.CompareCommits(context.Background(), event.GetRepo().GetOwner().GetName(), event.GetRepo().GetName(), event.GetBefore(), event.GetAfter())
		if err != nil {
			fmt.Println(err)
			return
		}

		analysis, err := handler.prAnalyzer.AnalyzePR(comp, r.Context())
		if err != nil {
			fmt.Println(err)
			return
		}
		utils.PrettyPrint(analysis)

	default:
		fmt.Println("event not valid")
		return
	}
}
