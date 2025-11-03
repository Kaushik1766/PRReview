package pr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"google.golang.org/genai"
)

type FileAnalysis struct {
	Filename string   `json:"filename"`
	Summary  string   `json:"summary"`
	Issues   []string `json:"issues"`
	Rating   int32    `json:"rating"`
}

type CommitAnalysis struct {
	Files []FileAnalysis `json:"files"`
}

type PRHandler struct{}

func NewPRHandler() *PRHandler {
	return &PRHandler{}
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
		handlePushEvent(event)
	default:
		fmt.Println("event not valid")
		return
	}
}

func handlePushEvent(pushEvent *github.PushEvent) {
	client := github.NewClient(nil)

	comp, _, err := client.Repositories.CompareCommits(context.Background(), "Kaushik1766", "TestRepo", pushEvent.GetBefore(), pushEvent.GetAfter())
	if err != nil {
		fmt.Println(err)
	}

	// for _, file := range comp.Files {
	// 	fmt.Printf("file: %s\nchanges:\n%s\n\n", file.GetFilename(), file.GetPatch())
	// }
	//
	err = getAIAnalysis(comp)
	if err != nil {
		fmt.Println(err)
	}
}

func getAIAnalysis(comp *github.CommitsComparison) error {
	apiKey := os.Getenv("GEMINI_API_KEY")
	fmt.Println(apiKey)

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}

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

	resp, err := client.Models.GenerateContent(context.Background(), "gemini-2.0-flash", []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}, nil)
	if err != nil {
		return err
	}

	// fmt.Println(resp.Text())

	var analysis []FileAnalysis
	if err := json.Unmarshal([]byte(resp.Text()[7:len(resp.Text())-3]), &analysis); err != nil {
		return fmt.Errorf("failed to parse model output: %w", err)
	}

	for _, val := range analysis {
		PrettyPrint(val)
	}
	return nil
}

func PrettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}
