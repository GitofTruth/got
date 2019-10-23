package datastructures

type PushLog struct {
	BranchName string
	Logs       map[string]CommitLog
}

func CreateNewPushLog(branchname string, logs map[string]CommitLog) (PushLog, error) {
	var pushLog PushLog
	pushLog.BranchName = branchname
	pushLog.Logs = logs

	return pushLog, nil
}
