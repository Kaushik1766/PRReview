package pranalysis

import (
	"context"

	"Kaushik1766/PRReview/internal/models"

	"github.com/google/go-github/github"
)

type PRAnalyzer interface {
	AnalyzePR(commit *github.CommitsComparison, ctx context.Context) (models.CommitAnalysis, error)
}
