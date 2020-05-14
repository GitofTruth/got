package datastructures

import "encoding/json"

type RepoBranch struct {
	Name      string               `json:"branchName"`
	Author    string               `json:"author"`
	Timestamp int                  `json:"timeStamp"`
	Logs      map[string]CommitLog `json:"logs"`
}

func CreateNewRepoBranch(name string, author string, timestamp int, logs map[string]CommitLog) (RepoBranch, error) {
	var repoBranch RepoBranch
	repoBranch.Name = name
	repoBranch.Author = author
	repoBranch.Timestamp = timestamp

	if logs == nil {
		repoBranch.Logs = make(map[string]CommitLog)
	} else {
		repoBranch.Logs = logs
	}

	return repoBranch, nil
}

// This function takes a json string that represents the marshalling of RepoBranch
// and returns a RepoBranch.
// The returned data is valid and consistent
func UnmarashalRepoBranch(objectString string) (RepoBranch, error) {
	var repoBranch RepoBranch

	json.Unmarshal([]byte(objectString), &repoBranch)

	// validate a branch
	// ??

	return repoBranch, nil
}

//Check if the hash has been added to the branch before.
func (branch *RepoBranch) IsCommitHash(hashName string) bool {
	_, exist := branch.Logs[hashName]
	return exist
}

//Checks if a log has at least one parent in the branch.
func (branch *RepoBranch) ValidLog(commitLog CommitLog) (bool, error) {

	if !branch.IsCommitHash(commitLog.Hash) {
		if len(branch.Logs) == 0 {
			return true, nil
		}

		for _, hash := range commitLog.Parenthashes {
			if branch.IsCommitHash(hash) {
				timeProgressing := branch.Logs[hash].CommitterTimestamp < commitLog.CommitterTimestamp
				return timeProgressing, nil
			}
		}
	}

	return false, nil
}

//Adds a CommitLog to the branch if it can be added according to the info avaiable to the branch.
func (branch *RepoBranch) AddCommitLog(commitLog CommitLog, passValidation bool) (bool, error) {

	if valid, _ := branch.ValidLog(commitLog); valid || passValidation {
		branch.Logs[commitLog.Hash] = commitLog
		return true, nil
	}

	return false, nil
}

//Removes a Log from the branch it exists, useful for rolling back cases.
func (branch *RepoBranch) RemoveLog(commitLogName string) (bool, error) {
	if branch.IsCommitHash(commitLogName) {
		delete(branch.Logs, commitLogName)
		return true, nil
	}

	return false, nil
}
