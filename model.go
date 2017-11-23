package main

// Committer is used to retrieve committer info.
type Committer struct {
	Committer string `json:"login"`
	Commits   int    `json:"contributions"`
}

// GitObject receives all information and is marshalled to JSON
type GitObject struct {
	Owner     string   `json:"owner"`
	Project   string   `json:"project"`
	Committer string   `json:"committer"`
	Commits   int      `json:"commits"`
	Languages []string `json:"language"`
}
