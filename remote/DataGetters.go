package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	client "github.com/GitofTruth/GoT/client"
	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func (contract *RepoContract) getRepoInstance(stub shim.ChaincodeStubInterface, args []string) (datastructures.Repo, error) {
	// repoAuthor, repoName

	fmt.Println("\n---------------------------\nQuerying the ledger .. getRepoInstance", args)
	defer fmt.Println("---------------------------")

	if len(args) < 2 {
		var repo datastructures.Repo
		fmt.Println("Incorrect number of arguments. Expecting 2")
		return repo, errors.New("Incorrect number of arguments. Expecting 2")
	}

	// getting the required information from first table.
	repoHash := GetRepoKey(args[0], args[1])
	repoData, err := stub.GetState(string(repoHash))
	if err != nil {
		var repo datastructures.Repo
		fmt.Println("Could not find requested Repo: ", err)
		return repo, errors.New("Could not find requested Repo")
	}
	fmt.Println("Found this repo:", string(repoData))
	// unmarashaling the data
	structuredRepoData := map[string]string{}
	err = json.Unmarshal([]byte(repoData), &structuredRepoData)
	if err != nil {
		var repo datastructures.Repo
		fmt.Println("Could not unmarashal requested repo: ", err)
		return repo, errors.New("Could not unmarashal requested repo")
	}
	timestamp, _ := strconv.Atoi(structuredRepoData["timeStamp"])
	users, _ := contract.getRepoUsers(stub, repoHash)
	encryptionKey, _ := datastructures.UnmarashalKeyAnnouncement(structuredRepoData["encryptionKey"])
	keyAnnouncements, _ := datastructures.UnmarashalKeyAnnouncements(structuredRepoData["keyAnnouncements"])
	repo, _ := datastructures.CreateNewRepo(structuredRepoData["repoName"], structuredRepoData["author"], structuredRepoData["directoryCID"], timestamp, nil, encryptionKey, users)
	repo.KeyAnnouncements = keyAnnouncements

	// getting the repo branches
	branchQueryString := fmt.Sprintf("{\"selector\": {\"docName\": \"branch\", \"repoID\": \"%s\"},\"fields\": [\"repoID\", \"branchName\", \"author\", \"timeStamp\"]}", repoHash)
	branchResultsIterator, err := stub.GetQueryResult(branchQueryString)
	if err != nil {
		fmt.Println("Couldnot find Requested Branch: ", err)
		var repo datastructures.Repo
		return repo, err
	}
	defer branchResultsIterator.Close()
	//iterating over branches
	for branchResultsIterator.HasNext() {
		branchString, err := branchResultsIterator.Next()
		if err != nil {
			fmt.Println("Couldnot proceed to next branch: ", err)
			var repo datastructures.Repo
			return repo, err
		}

		structuredBranchData := map[string]string{}
		err = json.Unmarshal([]byte(branchString.Value), &structuredBranchData)
		fmt.Println("Found This Branch: ", structuredBranchData)
		if err != nil {
			var repo datastructures.Repo
			fmt.Println("Could not unmarashal requested Branch: ", err)
			return repo, errors.New("Could not unmarashal requested Branch")
		}
		branchTimestamp, _ := strconv.Atoi(structuredBranchData["timeStamp"])
		branch, _ := datastructures.CreateNewRepoBranch(structuredBranchData["branchName"], structuredBranchData["author"], branchTimestamp, nil)
		repo.AddBranch(branch)

		//adding branch commits
		commitsQueryString := fmt.Sprintf("{\"selector\": {\"docName\": \"commit\", \"repoID\": \"%s\", \"branchName\": \"%s\"},\"fields\": [\"repoID\", \"branchName\",\"hash\", \"message\", \"author\", \"committer\", \"committerTimestamp\", \"CommitParenthashes\", \"signature\", \"storageHashes\"]}", repoHash, branch.Name)
		commitsResultsIterator, err := stub.GetQueryResult(commitsQueryString)
		if err != nil {
			var repo datastructures.Repo
			fmt.Println("Could not find requested commit: ", err)
			return repo, err
		}
		defer commitsResultsIterator.Close()
		//iterating over branches
		for commitsResultsIterator.HasNext() {
			commitString, err := commitsResultsIterator.Next()
			if err != nil {
				fmt.Println("Could not proceed to requested commit: ", err)
				var repo datastructures.Repo
				return repo, err
			}

			structuredCommitData := map[string]string{}
			err = json.Unmarshal([]byte(commitString.Value), &structuredCommitData)
			fmt.Println("Found This Commit: ", structuredCommitData)
			if err != nil {
				var repo datastructures.Repo
				fmt.Println("Could not unmarashal requested commit: ", err)
				return repo, errors.New("Could not unmarashal requested commit")
			}
			committerTimestamp, _ := strconv.Atoi(structuredCommitData["committerTimestamp"])
			var ph []string
			_ = json.Unmarshal([]byte(structuredCommitData["parentHashes"]), &ph)
			var s []byte
			_ = json.Unmarshal([]byte(structuredCommitData["signature"]), &s)
			var enc string
			_ = json.Unmarshal([]byte(structuredCommitData["encryptionKey"]), &enc)
			var sh map[string]string
			_ = json.Unmarshal([]byte(structuredCommitData["storageHashes"]), &sh)
			commit, _ := datastructures.CreateNewCommitLog(structuredCommitData["message"], structuredCommitData["author"], structuredCommitData["commiter"], committerTimestamp, structuredCommitData["hash"], ph, s, enc, sh)
			repo.AddCommitLog(commit, branch.Name, true)
		}
	}

	fmt.Println("This is the final fetched Repo: ", repo)
	return repo, nil
}

func (contract *RepoContract) queryRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName

	fmt.Println("Querying the ledger .. getRepo", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	repoHash := GetRepoKey(args[0], args[1])

	// key := map[string]interface{}{"repoID": repoHash}
	// jsonKey, _ := json.Marshal(key)
	repoData, err := stub.GetState(string(repoHash))

	if err != nil {
		return shim.Error("Repo does not exist")
	}

	fmt.Println("Found this repo:", string(repoData))

	return shim.Success(repoData)
}

func (contract *RepoContract) clone(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, currentUserName

	fmt.Println("Querying the ledger .. clone", args)

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3.")
	}

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	repo.UpdateCommitsForUser(args[2])

	fmt.Println("Found this repo:", repo)

	j, _ := json.Marshal(repo)
	return shim.Success([]byte(string(j)))
}

func (contract *RepoContract) queryBranches(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, currentUserName

	fmt.Println("Querying the ledger .. queryBranches", args)

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3.")
	}

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}
	repo.UpdateCommitsForUser(args[2])

	fmt.Println("Found these branches:", repo.GetBranches())

	j, _ := json.Marshal(repo.GetBranches())
	return shim.Success([]byte(string(j)))
}

func (contract *RepoContract) queryBranch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, branchName, currentUserName

	fmt.Println("Querying the ledger .. queryBranch", args)

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4.")
	}

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	if !repo.IsBranch(args[2]) {
		fmt.Println("Requested Branch Not found")
		return shim.Error("Requested Branch Not found")
	}
	repo.UpdateCommitsForUser(args[3])

	branch := repo.Branches[args[2]]
	fmt.Println("Found these branches:", branch)

	seralized, _ := json.Marshal(branch)
	return shim.Success([]byte(string(seralized)))
}

func (contract *RepoContract) queryBranchCommits(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, branchName, lastcommit, currentUserName

	fmt.Println("Querying the ledger .. queryBranchCommits", args)

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5.")
	}

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	if !repo.IsBranch(args[2]) {
		fmt.Println("Requested Branch Not found")
		return shim.Error("Requested Branch Not found")
	}
	repo.UpdateCommitsForUser(args[4])

	branch := repo.Branches[args[2]]
	fmt.Println("Found this branch:", branch)

	commits := make([]datastructures.CommitLog, 0)

	//last commit exist?
	t := 0
	if !repo.IsCommitHash(args[3]) && args[3] != "" {
		fmt.Println("Requested BranchCommit Not found")
		return shim.Error("Requested BranchCommit Not found")
	}

	if args[3] != "" {
		t = branch.Logs[args[3]].CommitterTimestamp
	}

	// get all commits after this time
	for _, log := range branch.Logs {
		if log.CommitterTimestamp > t {
			commits = append(commits, log)
		}
	}

	seralized, _ := json.Marshal(commits)
	return shim.Success([]byte(string(seralized)))
}

func (contract *RepoContract) queryLastBranchCommit(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, branchName

	fmt.Println("Querying the ledger .. queryBranchCommits", args)

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3.")
	}

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	if !repo.IsBranch(args[2]) {
		fmt.Println("Requested Branch Not found")
		return shim.Error("Requested Branch Not found")
	}

	branch := repo.Branches[args[2]]
	fmt.Println("Found this branch:", branch)

	t := 0
	hash := ""
	// get all commits after this time
	for _, log := range branch.Logs {
		if log.CommitterTimestamp >= t {
			t = log.CommitterTimestamp
			hash = log.Hash
		}
	}

	return shim.Success([]byte(hash))
}

func (contract *RepoContract) getUserInfo(stub shim.ChaincodeStubInterface, userName string) (client.UserInfo, peer.Response) {
	var userInfo client.UserInfo

	userQueryString := fmt.Sprintf("{\"selector\": {\"docName\": \"user\", \"userName\": \"%s\"},\"fields\": [\"userName\", \"publicKey\"]}", userName)
	userResultsIterator, err := stub.GetQueryResult(userQueryString)
	if err != nil {
		fmt.Println("Could not find Requested User: ", err)
		return userInfo, shim.Error("User does not exist")
	}
	defer userResultsIterator.Close()

	for userResultsIterator.HasNext() {
		userString, err := userResultsIterator.Next()
		if err != nil {
			fmt.Println("Could not proceed to next user: ", err)
			return userInfo, shim.Error("Could not proceed to next user")
		}

		structuredUserData := map[string]string{}
		err = json.Unmarshal([]byte(userString.Value), &structuredUserData)
		fmt.Println("Found This User: \t", structuredUserData)
		if err != nil {
			fmt.Println("Could not unmarashal requested Branch: ", err)
			return userInfo, shim.Error("Could not unmarashal requested Branch")
		}
		userInfo.UserName = structuredUserData["userName"]
		userInfo.PublicKey = structuredUserData["publicKey"]
	}

	return userInfo, shim.Success([]byte(""))
}

func (contract *RepoContract) queryUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// userName

	fmt.Println("Querying the ledger .. queryUser", args)

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}

	userInfo, failMessage := contract.getUserInfo(stub, args[0])
	if failMessage.Message != "" {
		return failMessage
	}

	seralized, _ := json.Marshal(userInfo)
	return shim.Success([]byte(string(seralized)))
}

func (contract *RepoContract) queryUsers(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// [userName1 userName2]

	fmt.Println("Querying the ledger .. queryUser", args)

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}

	userInfos := make([]client.UserInfo, 0)
	var userNames []string
	_ = json.Unmarshal([]byte(args[0]), &userNames)

	for _, username := range userNames {
		userInfo, failMessage := contract.getUserInfo(stub, username)
		if failMessage.Message != "" {
			continue
		}
		userInfos = append(userInfos, userInfo)
	}

	seralized, _ := json.Marshal(userInfos)
	return shim.Success([]byte(string(seralized)))
}

func (contract *RepoContract) getRepoUsers(stub shim.ChaincodeStubInterface, repoHash string) (map[string]datastructures.UserAccess, peer.Response) {

	users := make(map[string]datastructures.UserAccess, 0)

	accessQueryString := fmt.Sprintf("{\"selector\": {\"docName\": \"userAccess\", \"repoHash\": \"%s\"},\"fields\": [\"authorized\", \"userAccess\", \"authorizer\"]}", repoHash)
	accessResultsIterator, err := stub.GetQueryResult(accessQueryString)
	if err != nil {
		fmt.Println("Could not find Repo Access: ", err)
		return users, shim.Error("Repo Access does not exist")
	}
	defer accessResultsIterator.Close()

	for accessResultsIterator.HasNext() {
		accessString, err := accessResultsIterator.Next()
		if err != nil {
			fmt.Println("Could not proceed to user access: ", err)
			return users, shim.Error("Could not proceed to next user access")
		}

		structuredAccessData := map[string]string{}
		err = json.Unmarshal([]byte(accessString.Value), &structuredAccessData)
		fmt.Println("Found This User: \t", structuredAccessData)
		if err != nil {
			fmt.Println("Could not unmarashal requested Branch: ", err)
			return users, shim.Error("Could not unmarashal requested Branch")
		}
		access, err := strconv.Atoi(structuredAccessData["userAccess"])
		if err != nil {
			fmt.Println("Could not parse UserAcess: ", access, err)
			return users, shim.Error("Could not parse UserAcess")
		} else {
			users[structuredAccessData["authorized"]] = datastructures.UserAccess(access)
		}
	}

	return users, shim.Success([]byte(""))
}

func (contract *RepoContract) queryRepoUserAccess(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName

	fmt.Println("Querying the ledger .. queryRepoUserAccess", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	repoHash := GetRepoKey(args[0], args[1])

	users, failMessage := contract.getRepoUsers(stub, repoHash)
	if failMessage.Message != "" {
		return failMessage
	}

	seralized, _ := json.Marshal(users)
	return shim.Success([]byte(string(seralized)))
}

// indexName := "index-Branch"
// branchIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, args[2]})
// // jsonKey, _ := json.Marshal(key)
// branchData, err := stub.GetState(string(branchIndexKey))
