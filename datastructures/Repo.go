package datastructures

import (
	"encoding/json"
	"fmt"
)

type Repo struct {
	Name         string                `json:"repoName"`
	Author       string                `json:"author"`
	DirectoryCID string 								`json:"directoryCID"`
	Timestamp    int                   `json:"timeStamp"`
	CommitHashes map[string]struct{}   `json:"hashes"`
	Branches     map[string]RepoBranch `json:"branches"`
}

func CreateNewRepo(name string, author string, directoryCID string , timestamp int, branches map[string]RepoBranch) (Repo, error) {
	var repo Repo

	repo.Name = name
	repo.Author = author
	repo.Timestamp = timestamp
	repo.DirectoryCID = directoryCID
	repo.CommitHashes = make(map[string]struct{})

	// the first commit
	var empty struct{}
	repo.CommitHashes["0000000000000000000000000000000000000000"] = empty
	//check what is built on this hash

	if branches == nil {
		fmt.Println("empty branch is created!")
		repo.Branches = make(map[string]RepoBranch)
	} else {
		repo.Branches = branches
		for _, branch := range branches {
			for _, log := range branch.Logs {
				repo.AddCommitHash(log)
			}
		}
	}

	fmt.Println("empty repo is created!")

	return repo, nil
}

func UnmarashalRepo(objectString string) (Repo, error) {
	var unmarashaledRepo Repo
	json.Unmarshal([]byte(objectString), &unmarashaledRepo)

	repo, _ := CreateNewRepo(unmarashaledRepo.Name, unmarashaledRepo.Author, unmarashaledRepo.DirectoryCID, unmarashaledRepo.Timestamp, nil)

	// var repo Repo

	// repo.Name = unmarashaledRepo.Name
	// repo.Author = unmarashaledRepo.Author
	// repo.Timestamp = unmarashaledRepo.Timestamp

	// repo.CommitHashes = make(map[string]struct{})
	// repo.Branches = make(map[string]RepoBranch)

	for _, branch := range unmarashaledRepo.Branches {
		newBranch, _ := CreateNewRepoBranch(branch.Name, branch.Author, branch.Timestamp, nil)
		repo.AddBranch(newBranch)
		for _, log := range branch.Logs {
			repo.AddCommitLog(log, branch.Name)
		}
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

func (repo *Repo) GetBranches() []string {

	keys := make([]string, 0, len(repo.Branches))
	for k := range repo.Branches {
		keys = append(keys, k)
	}
	return keys
}

//checks that all hash parents are hashes
func (repo *Repo) ValidCommitLog(commitLog CommitLog, branchName string) (bool, error) {

	return true, nil

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
	return true, nil

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
	fmt.Println("Trying to add new branch ")

	if valid, _ := repo.ValidBranch(branch); valid {
		fmt.Println("New branch is valid!")

		repo.Branches[branch.Name] = branch
		for _, log := range branch.Logs {
			repo.AddCommitHash(log)
		}

		return true, nil
	}

	return false, nil
}
