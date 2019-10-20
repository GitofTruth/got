package datastructures

type Repo struct {
	Author       string                `json:"author"`
	Timestamp    int                   `json:"timestamp"`
	CommitHashes map[string]struct{}   `json:"hashes"`
	Branches     map[string]RepoBranch `json:"branches"`
}

func CreateNewRepo(author string, timestamp int, branches map[string]RepoBranch) (Repo, error) {
	var repo Repo

	repo.Author = author
	repo.Timestamp = timestamp
	repo.CommitHashes = make(map[string]struct{})

	for _, branch := range branches {
		for _, log := range branch.Logs {
			repo.AddCommitHash(log)
		}
	}

	if branches == nil {
		repo.Branches = make(map[string]RepoBranch)
	} else {
		repo.Branches = branches
	}

	return repo, nil
}

func (repo *Repo) IsCommitHash(hashName string) bool {
	_, exist := repo.CommitHashes[hashName]
	return exist
}

func (repo *Repo) IsBranch(branchName string) bool {
	_, exist := repo.Branches[branchName]
	return exist
}

//checks that all hash parents are hashes
func (repo *Repo) ValidCommitLog(commitLog CommitLog, branchName string) (bool, error) {

	if repo.IsBranch(branchName) {
		branch := repo.Branches[branchName]
		if valid, _ := branch.ValidLog(commitLog); valid {

			allParentsAreHashes := true
			for _, parentHash := range commitLog.Parenthashes {
				allParentsAreHashes = allParentsAreHashes && repo.IsCommitHash(parentHash)
				if !allParentsAreHashes {
					break
				}
			}

			if allParentsAreHashes {
				return true, nil
			}
		}
	}

	return false, nil
}

//Check that all parents hashes are current hashes or previous ones in the new hashes list
func (repo *Repo) ValidCommitLogs(commitLogs []CommitLog, branchName string) (bool, error) {
	var previousHashes map[string]struct{}
	previousHashes = make(map[string]struct{})

	if repo.IsBranch(branchName) {
		branch := repo.Branches[branchName]

		for _, commitLog := range commitLogs {
			if valid, _ := branch.ValidLog(commitLog); valid {

				allParentsAreHashes := true
				for _, parentHash := range commitLog.Parenthashes {
					_, exist := previousHashes[parentHash]
					allParentsAreHashes = allParentsAreHashes && (repo.IsCommitHash(parentHash) || exist)
					if !allParentsAreHashes {
						break
					}
				}
				if !allParentsAreHashes {
					return false, nil
				}

				var empty struct{}
				previousHashes[commitLog.Hash] = empty
			} else {
				return false, nil
			}
		}
	}

	return true, nil
}

//need to check commits in the branch
func (repo *Repo) ValidBranch(Branch RepoBranch) (bool, error) {
	return !repo.IsBranch(Branch.Name), nil
}

func (repo *Repo) AddCommitHash(commitLog CommitLog) bool {
	var empty struct{}
	repo.CommitHashes[commitLog.Hash] = empty
	return true
}

func (repo *Repo) AddCommitHashes(commitLogs []CommitLog) bool {
	var empty struct{}
	for _, commitLog := range commitLogs {
		repo.CommitHashes[commitLog.Hash] = empty
	}
	return true
}

func (repo *Repo) AddCommitLog(commitLog CommitLog, branchName string) (bool, error) {

	if valid, _ := repo.ValidCommitLog(commitLog, branchName); valid {

		branch := repo.Branches[branchName]
		if done, _ := branch.AddCommitLog(commitLog); done {
			repo.Branches[branchName] = branch
			repo.AddCommitHash(commitLog)
			return true, nil
		}

	}

	return false, nil
}

//What if not all the commits were added? rolling back?
func (repo *Repo) AddCommitLogs(commitLogs []CommitLog, branchName string) (bool, error) {

	if valid, _ := repo.ValidCommitLogs(commitLogs, branchName); valid {

		branch := repo.Branches[branchName]
		fullDone := true
		for _, commitLog := range commitLogs {
			if done, _ := branch.AddCommitLog(commitLog); !done {
				fullDone = false
				break
			}
		}

		if fullDone {
			repo.Branches[branchName] = branch
			repo.AddCommitHashes(commitLogs)
		}
	}

	return false, nil
}

func (repo *Repo) AddBranch(branch RepoBranch) (bool, error) {

	if valid, _ := repo.ValidBranch(branch); valid {

		repo.Branches[branch.Name] = branch
		for _, log := range branch.Logs {
			repo.AddCommitHash(log)
		}

		return true, nil
	}

	return false, nil
}