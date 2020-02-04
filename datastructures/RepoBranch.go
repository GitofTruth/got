package datastructures

type RepoBranch struct {
	Name      string               `json:"branchName"`
	Author    string               `json:"author"`
	Timestamp int                  `json:"timestamp"`
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

//Check if the hash has been added to the branch before.
func (branch *RepoBranch) IsCommitHash(hashName string) bool {
	_, exist := branch.Logs[hashName]
	return exist
}

//Checks if a log has at least one parent in the branch.
func (branch *RepoBranch) ValidLog(commitLog CommitLog) (bool, error) {

	if !branch.IsCommitHash(commitLog.Hash) {
		for _, hash := range commitLog.Parenthashes {
			if branch.IsCommitHash(hash) {
				return true, nil
			}
		}
	}

	return false, nil
}

//Adds a CommitLog to the branch if it can be added according to the info avaiable to the branch.
func (branch *RepoBranch) AddCommitLog(commitLog CommitLog) (bool, error) {

	if valid, _ := branch.ValidLog(commitLog); valid {
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
