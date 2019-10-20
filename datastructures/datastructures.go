package datastructures

import "github.com/GitofTruth/GoT/datastructures"

type PushLog struct {
	BranchName string
	Logs       []datastructures.CommitLog
}
