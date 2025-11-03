package pranalysis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"Kaushik1766/PRReview/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/google/go-github/github"
)

type PRAnalysis struct {
	aiClient *bedrockruntime.Client
}

func NewPRAnalysis(aiClient *bedrockruntime.Client) (*PRAnalysis, error) {
	return &PRAnalysis{
		aiClient: aiClient,
	}, nil
}

func (service PRAnalysis) AnalyzePR(comp *github.CommitsComparison, ctx context.Context) (models.CommitAnalysis, error) {
	var diffs []string
	for _, file := range comp.Files {
		diffs = append(diffs, fmt.Sprintf("File: %s\nPatch:\n%s", file.GetFilename(), file.GetPatch()))
	}

	prompt := fmt.Sprintf(`
		You are a senior software engineer reviewing a single Git commit.

		Your job:
		1. Understand what the commit is doing and why.
		2. Review each file changed for:
			- correctness and logic errors
			- readability and maintainability
			- security or performance concerns
			- potential improvements or refactoring opportunities
		3. Rate each file on code quality and provide actionable feedback.

		Return your answer strictly in JSON format:
		[
			{
				"file": "<filename>",
				"summary": "<short summary of what changed>",
				"issues": ["<issue 1>", "<issue 2>", ...],
				"suggestions": ["<suggestion 1>", "<suggestion 2>", ...],
				"rating": <1-5 integer>
			}
		]

		Commit Diff:
		%s
	`, strings.Join(diffs, "\n\n"))

	payload := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        2048,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("failed to marshal payload: %v", err)
	}
	resp, err := service.
		aiClient.
		InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String("global.anthropic.claude-sonnet-4-5-20250929-v]byte(prompt),1:0"),
			ContentType: aws.String("application/json"),
			Body:        body,
		})
	if err != nil {
		return models.CommitAnalysis{}, fmt.Errorf("failed to generate content: %w", err)
	}

	fmt.Println(string(resp.Body))

	// fmt.Println(resp.Text())

	// var analysis []models.FileAnalysis
	// if err := json.Unmarshal([]byte(resp.Text()[7:len(resp.Text())-3]), &analysis); err != nil {
	// 	return models.CommitAnalysis{}, fmt.Errorf("failed to parse model output: %w", err)
	// }

	// for _, val := range analysis {
	// 	utils.PrettyPrint(val)
	// }
	// return models.CommitAnalysis{}, nil
	return models.CommitAnalysis{}, nil
}
