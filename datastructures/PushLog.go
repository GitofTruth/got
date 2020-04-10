package datastructures

import "encoding/json"

type PushLog struct {
	BranchName string      `json:"branchName"`
	DirectoryCID string			`json:"directoryCID"`
	Logs       []CommitLog `json:"logs"`
}

func CreateNewPushLog(branchname string, directoryCID string, logs []CommitLog) (PushLog, error) {
	var pushLog PushLog
	pushLog.BranchName = branchname
	pushLog.DirectoryCID = directoryCID
	pushLog.Logs = logs

	return pushLog, nil
}

func UnmarashalPushLog(objectString string) (PushLog, error) {
	var pushLog PushLog

	json.Unmarshal([]byte(objectString), &pushLog)

	// validate a pushLog
	// ??

	return pushLog, nil
}
