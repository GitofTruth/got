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

	for branchName, branch := range branches{
		for commitHash, log := range branch.Logs{
			var empty struct{}
			repo.CommitHashes[commitHash] = empty
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

		allParentsAreHashes := true
		for index, parentHash := range commitLog.Parenthashes{
				allParentsAreHashes = allParentsAreHashes && repo.IsCommitHash(parentHash)
				if (!allParentsAreHashes){
					break
				}
		}

		if(allParentsAreHashes){
				return true, nil
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

		for _, commitLog := range commitLogs{
			if valid,_ := branch.ValidLog(commitLog); valid{

				allParentsAreHashes := true
				for _, parentHash := range commitLog.Parenthashes{

						allParentsAreHashes = allParentsAreHashes && repo.IsCommitHash(parentHash)
						if (!allParentsAreHashes){
							break
							}

			}

		}


		}

		if(allParentsAreHashes){
				return true, nil
		}

	}

	return false, nil
}

func (repo *Repo) ValidBranch(repoBranch RepoBranch) (bool, error) {
	return true, nil
}

func (repo *Repo) AddCommitLogInBranch(commitLog CommitLog, branchName string) (bool, error) {

	if repo.IsBranch(branchName) {
		branch := repo.Branches[branchName]

		allParentsAreHashes := true
		for index, parentHash := range commitLog.Parenthashes{
				allParentsAreHashes = allParentsAreHashes && repo.IsCommitHash(parentHash)
				if (!allParentsAreHashes){
					break
				}
		}

		if(allParentsAreHashes){
			if done, _ := branch.AddLog(commitLog); done {
				repo.Branches[branchName] = branch
				return true, nil
			}
		}

	}

	return false, nil
}


func (repo *Repo) AddBranch(branch RepoBranch) (bool, error) {

	if !repo.IsBranch(branch.Name) {
		repo.Branches[branch.Name] = branch
			for commitHash, log := range branch.Logs{
				repo.CommitHashes[commitHash] = struct{}
			}
		return true, nil
	}

	return false, nil
}
