package datastructures

import "encoding/json"

// This struct is used to model the required data to add a push
// to a got project
type PushLog struct {
	BranchName   string      `json:"branchName"`
	DirectoryCID string      `json:"directoryCID"`
	Logs         []CommitLog `json:"logs"`
}

// Helper function that creates a new object instance of a PushLog
func CreateNewPushLog(branchname string, directoryCID string, logs []CommitLog) (PushLog, error) {
	var pushLog PushLog
	pushLog.BranchName = branchname
	pushLog.DirectoryCID = directoryCID
	pushLog.Logs = logs

	return pushLog, nil
}

// This function takes a json string that represents the marshalling of PushLog
// and returns a PushLog.
// The returned data is valid and consistent
func UnmarashalPushLog(objectString string) (PushLog, error) {
	var pushLog PushLog

	json.Unmarshal([]byte(objectString), &pushLog)

	// validate a pushLog
	// ??

	return pushLog, nil
}
