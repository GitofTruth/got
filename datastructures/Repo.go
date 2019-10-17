package datastructures

type Repo struct {
	Author    string                `json:"author"`
	Timestamp int                   `json:"timestamp"`
	Branches  map[string]RepoBranch `json:"branches"`
}
