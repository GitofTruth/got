package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func (contract *RepoContract) getRepoInstance(stub shim.ChaincodeStubInterface, args []string) (datastructures.Repo, error) {
	// repoAuthor, repoName

	fmt.Println("Querying the ledger .. getRepoInstance", args)

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
	timestamp, _ := strconv.Atoi(structuredRepoData["timestamp"])
	repo, _ := datastructures.CreateNewRepo(structuredRepoData["repoName"], structuredRepoData["author"], timestamp, nil)

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
		branchTimestamp, _ := strconv.Atoi(structuredBranchData["timestamp"])
		branch, _ := datastructures.CreateNewRepoBranch(structuredBranchData["branchName"], structuredBranchData["author"], branchTimestamp, nil)
		repo.AddBranch(branch)

		//adding branch commits
		commitsQueryString := fmt.Sprintf("{\"selector\": {\"docName\": \"commit\", \"repoID\": \"%s\", \"branchName\": \"%s\"},\"fields\": [\"repoID\", \"branchName\",\"hash\", \"message\", \"author\", \"committer\", \"committerTimestamp\", \"CommitParenthashes\", \"signature\"]}", repoHash, branch.Name)
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
			_ = json.Unmarshal([]byte(structuredCommitData["parenthashes"]), &ph)
			var s []byte
			_ = json.Unmarshal([]byte(structuredCommitData["signature"]), &s)
			commit, _ := datastructures.CreateNewCommitLog(structuredCommitData["message"], structuredCommitData["author"], structuredCommitData["commiter"], committerTimestamp, structuredCommitData["hash"], ph, s)
			repo.AddCommitLog(commit, branch.Name)
		}
	}

	fmt.Println("This is the final fetched Repo: ", repo)
	return repo, nil
}

func (contract *RepoContract) getRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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

func (contract *RepoContract) queryBranches(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName

	fmt.Println("Querying the ledger .. queryBranches", args)

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	fmt.Println("Found these branches:", repo.GetBranches())

	j, _ := json.Marshal(repo.GetBranches())
	return shim.Success([]byte(string(j)))
}

func (contract *RepoContract) queryBranch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, branchName

	fmt.Println("Querying the ledger .. queryBranch", args)

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
	fmt.Println("Found these branches:", branch)

	seralized, _ := json.Marshal(branch)
	return shim.Success([]byte(string(seralized)))
}

// indexName := "index-Branch"
// branchIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, args[2]})
// // jsonKey, _ := json.Marshal(key)
// branchData, err := stub.GetState(string(branchIndexKey))
