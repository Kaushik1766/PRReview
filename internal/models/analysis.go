package models

type FileAnalysis struct {
	Filename string   `json:"filename"`
	Summary  string   `json:"summary"`
	Issues   []string `json:"issues"`
	Severity int32    `json:"severity"`
}

type CommitAnalysis struct {
	Files []FileAnalysis `json:"files"`
}

type Message struct {
	Model   string    `json:"model"`
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
