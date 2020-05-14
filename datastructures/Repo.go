package datastructures

import (
	"encoding/json"
	"fmt"
	"sort"
)

// This enum represents the type of access a user has for a repo
type UserAccess int

const (
	ReadWriteAccess    UserAccess = 1
	CollaboratorAccess UserAccess = 2
	OwnerAccess        UserAccess = 3
	ReovkedAccess      UserAccess = 4
	NeverSetAccess     UserAccess = 5
)

// A struct that contains the required data to keep track about who is responsible
// of another user's access in the repository.
type AccessLog struct {
	Authorizer string `json:"authorizer"`
	Authorized string `json:"authorized"`
	UserAccess `json:"userAccess"`
}

// KeyAnnouncement is struct that is used to keep information about
// the avaiable symmetric encryption keys.
// it maps the has of the symmetric encryption key to the encrypted
// version of the key that a user can then decrypt using their private key.
type KeyAnnouncement struct {
	KeyHash       string            `json:"keyHash"`
	EncryptedKeys map[string]string `json:"encryptedKeys"`
}

// This function takes a json string that represents the marshalling
// of map that has an encryption key's hash as its key and a KeyAnnouncement
// as its value. It returns this specified map.
// The returned data is valid and consistent
func UnmarashalKeyAnnouncements(objectString string) (map[string]KeyAnnouncement, error) {
	var keyAnnouncements map[string]KeyAnnouncement
	json.Unmarshal([]byte(objectString), &keyAnnouncements)
	return keyAnnouncements, nil
}

// This function takes a json string that represents
// the marshalling of KeyAnnouncement and returns a KeyAnnouncement.
// The returned data is valid and consistent
func UnmarashalKeyAnnouncement(objectString string) (KeyAnnouncement, error) {
	var keyAnnouncement KeyAnnouncement
	json.Unmarshal([]byte(objectString), &keyAnnouncement)
	return keyAnnouncement, nil
}

// TODO: should we check for mutex or the hyperledger handles this???
// TODO: now handle the 	EncryptionKey interface{}, Users map[string]UserAccessAccessLogs, []AccessLog	into Coachdb

// This structis the entry point to store all the needed data for
// repo to ensure that it's working smoothly whether it's metadata
// or data required for access control or storage access.
type Repo struct {
	Name         string `json:"repoName"`
	Author       string `json:"author"`
	DirectoryCID string `json:"directoryCID"`
	Timestamp    int    `json:"timeStamp"`

	CommitHashes map[string]struct{}   `json:"hashes"`
	Branches     map[string]RepoBranch `json:"branches"`

	EncryptionKey    KeyAnnouncement            `json:"encryptionKey"`
	Users            map[string]UserAccess      `json:"users"`
	AccessLogs       []AccessLog                `json:"accessLogs"`
	KeyAnnouncements map[string]KeyAnnouncement `json:"keyAnnouncements"`
}

// returns the current access type of the specified user's userName
func (repo *Repo) GetUserAccess(userName string) UserAccess {
	if val, exist := repo.Users[userName]; exist {
		return val
	}
	return NeverSetAccess
}

// checks if the mentioned userName  belongs to
// user that is authorized to do read/write for the repository.
// anyone included in the repo can read / write
func (repo *Repo) CanEdit(userName string) bool {
	if val, exist := repo.Users[userName]; exist {
		return val != ReovkedAccess
	}
	return false
}

// checks if the mentioned userName  belongs to
// user that is authorized to authorize or revoke ReadWriteAccess
// for the repository.
// Only collaborators and owners CanAuthorize
func (repo *Repo) CanAuthorize(userName string) bool {
	if val, exist := repo.Users[userName]; exist {
		return val != ReovkedAccess && val != ReadWriteAccess
	}
	return false
}

// checks if the mentioned userName belongs to
// user that is authorized to authorize or revoke
// CollaboratorAccess for the repository.
// Only owners CanAuthorizeCollaborator
func (repo *Repo) CanAuthorizeCollaborator(userName string) bool {
	if val, exist := repo.Users[userName]; exist {
		return val == OwnerAccess
	}
	return false
}

// check if the authorized has enough permissions to update the
// authorized's current access type to the mentioned value.
func (repo *Repo) ValidUpdateAccess(authorized string, userAccess UserAccess, authorizer string) bool {

	if userAccess == ReadWriteAccess || userAccess == ReovkedAccess {
		if (repo.CanAuthorize(authorizer) && !repo.CanAuthorize(authorized)) || repo.CanAuthorizeCollaborator(authorizer) {
			//return repo.GetUserAccess(authorized) != NeverSetAccess
			return true
		}
	} else if repo.CanAuthorizeCollaborator(authorizer) {
		// TODO: check if you want to have many owners. maybe owners only revoke themselves?
		// return repo.GetUserAccess(authorized) != NeverSetAccess
		return true
	}

	return false
}

// It does the required data writing work to update a user's
// access type in case, the user access update is valid.
func (repo *Repo) UpdateAccess(authorized string, userAccess UserAccess, authorizer string, keyAnnouncement KeyAnnouncement) bool {

	if repo.ValidUpdateAccess(authorized, userAccess, authorizer) {
		if val, exist := repo.Users[authorized]; exist {
			if val == userAccess {
				return false
			}
		}

		// you are actually doing something here
		var accessLog AccessLog
		accessLog.Authorizer = authorizer
		accessLog.Authorized = authorized
		accessLog.UserAccess = userAccess

		repo.AccessLogs = append(repo.AccessLogs, accessLog)
		repo.Users[authorized] = userAccess

		if keyAnnouncement.KeyHash != "" {
			repo.EncryptionKey = keyAnnouncement
			repo.KeyAnnouncements[keyAnnouncement.KeyHash] = keyAnnouncement
		}

		return true
	}

	return false
}

// helper function that is needed to create a new Repo instance
func CreateNewRepo(name string, author string, directoryCID string, timestamp int, branches map[string]RepoBranch, encryptionKey KeyAnnouncement, users map[string]UserAccess) (Repo, error) {
	var repo Repo

	repo.Name = name
	repo.Author = author
	repo.Timestamp = timestamp
	repo.DirectoryCID = directoryCID
	repo.CommitHashes = make(map[string]struct{})
	repo.EncryptionKey = encryptionKey

	if users != nil {
		repo.Users = users
	} else {
		repo.Users = make(map[string]UserAccess, 0)
		repo.Users[repo.Author] = OwnerAccess
	}

	// the first commit
	// var empty struct{}
	// repo.CommitHashes["0000000000000000000000000000000000000000"] = empty
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

// This function takes a json string that represents the marshalling of Repo
// and returns a Repo.
// The returned data is valid and consistent
func UnmarashalRepo(objectString string) (Repo, error) {
	var unmarashaledRepo Repo
	json.Unmarshal([]byte(objectString), &unmarashaledRepo)

	// TODO: check number of owenrs for first creation in the contract

	repo, _ := CreateNewRepo(unmarashaledRepo.Name, unmarashaledRepo.Author, unmarashaledRepo.DirectoryCID, unmarashaledRepo.Timestamp, nil, unmarashaledRepo.EncryptionKey, unmarashaledRepo.Users)
	repo.KeyAnnouncements = unmarashaledRepo.KeyAnnouncements

	// var repo Repo

	// repo.Name = unmarashaledRepo.Name
	// repo.Author = unmarashaledRepo.Author
	// repo.Timestamp = unmarashaledRepo.Timestamp

	// repo.CommitHashes = make(map[string]struct{})
	// repo.Branches = make(map[string]RepoBranch)

	// TODO: now you need to initialize a repoBranch with its master first
	for _, branch := range unmarashaledRepo.Branches {
		newBranch, _ := CreateNewRepoBranch(branch.Name, branch.Author, branch.Timestamp, nil)
		repo.AddBranch(newBranch)

		logsList := make([]CommitLog, 0, len(branch.Logs))

		for _, v := range branch.Logs {
			logsList = append(logsList, v)
		}

		// sort the logs with time ascending
		sort.Slice(logsList, func(i, j int) bool {
			return logsList[i].CommitterTimestamp > logsList[j].CommitterTimestamp
		})

		for i := len(logsList) - 1; i >= 0; i-- {
			repo.AddCommitLog(logsList[i], branch.Name, false)
		}

		// for _, log := range branch.Logs {
		// 	repo.AddCommitLog(log, branch.Name)
		// }
	}

	return repo, nil
}

// checks if the provided hash has belonged to one of the repo's
// branches
func (repo *Repo) IsCommitHash(hashName string) bool {
	_, exist := repo.CommitHashes[hashName]
	return exist
}

// checks if the mentioned branch Name belongs to this repo.
func (repo *Repo) IsBranch(branchName string) bool {
	_, exist := repo.Branches[branchName]
	return exist
}

// returns a list that contains the names of branches in a project
func (repo *Repo) GetBranches() []string {

	keys := make([]string, 0, len(repo.Branches))
	for k := range repo.Branches {
		keys = append(keys, k)
	}
	return keys
}

//checks that all hash parents are hashes in the repository
func (repo *Repo) ValidCommitLog(commitLog CommitLog, branchName string) (bool, error) {

	if repo.IsBranch(branchName) {
		branch := repo.Branches[branchName]
		if len(repo.CommitHashes) == 0 {
			return true, nil
		}
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

// Adds a commitLog to a branch if it creates a new valid state
func (repo *Repo) AddCommitLog(commitLog CommitLog, branchName string, passValidation bool) (bool, error) {

	if valid, _ := repo.ValidCommitLog(commitLog, branchName); valid || passValidation {

		branch := repo.Branches[branchName]
		if done, _ := branch.AddCommitLog(commitLog, passValidation); done {
			repo.Branches[branchName] = branch
			repo.AddCommitHash(commitLog)
			return true, nil
		}

	}

	return false, nil
}

// Adds a list of commitLogs to a branch if it creates a new valid state
func (repo *Repo) AddCommitLogs(commitLogs []CommitLog, branchName string, passValidation bool) (bool, error) {

	if valid, _ := repo.ValidCommitLogs(commitLogs, branchName); valid || passValidation {

		branch := repo.Branches[branchName]
		fullDone := true
		for _, commitLog := range commitLogs {
			if done, _ := branch.AddCommitLog(commitLog, passValidation); !done {
				fullDone = false
				break
			}
		}

		if fullDone {
			repo.Branches[branchName] = branch
			repo.AddCommitHashes(commitLogs)
			return true, nil
		}
	}

	return false, nil
}

// Adds a new beanch to the repo if it creates a new valid state
// a new branch must have a new unique name and it must be consistent
// and builds on the current repo state
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

// Originally the encryptedKey inside each CommitLog contains the
// hash of the encrytion key
// this functions replaces each hash with a user's encryted version
// of the key.
func (repo *Repo) UpdateCommitsForUser(userName string) bool {
	fmt.Println("\n\nBefore Update:", repo)
	for branchName := range repo.Branches {
		for commitHash := range repo.Branches[branchName].Logs {
			keyHash := repo.Branches[branchName].Logs[commitHash].EncryptionKey
			log := repo.Branches[branchName].Logs[commitHash]
			log.EncryptionKey = repo.KeyAnnouncements[keyHash].EncryptedKeys[userName]
			repo.Branches[branchName].Logs[commitHash] = log
		}
	}
	fmt.Println("\n\nAfter Update:", repo)
	return true
}
