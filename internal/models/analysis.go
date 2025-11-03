package models

type FileAnalysis struct {
	Filename string   `json:"filename"`
	Summary  string   `json:"summary"`
	Issues   []string `json:"issues"`
	Rating   int32    `json:"rating"`
}

type CommitAnalysis struct {
	Files []FileAnalysis `json:"files"`
}
