package datastructures

type PushLog struct {
	BranchName string
	Logs       []CommitLog
}

func CreateNewPushLog(branchname string, logs []CommitLog) (PushLog, error) {
	var pushLog PushLog
	pushLog.BranchName = branchname
	pushLog.Logs = logs

	return pushLog, nil
}
