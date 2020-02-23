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

	fmt.Println("Querying the ledger..")

	if len(args) != 2 {
		var repo datastructures.Repo
		return repo, errors.New("Incorrect number of arguments. Expecting 2")
	}

	// getting the required information from first table.
	repoHash := GetRepoKey(args[0], args[1])
	repoData, err := stub.GetState(string(repoHash))
	if err != nil {
		var repo datastructures.Repo
		return repo, errors.New("Could not find requested Repo")
	}
	fmt.Println("Found this repo:", string(repoData))
	// unmarashaling the data
	structuredRepoData := map[string]string{}
	err = json.Unmarshal([]byte(repoData), &structuredRepoData)
	if err != nil {
		var repo datastructures.Repo
		return repo, errors.New("Could not unmarashal requested repo")
	}
	timestamp, _ := strconv.Atoi(structuredRepoData["timestamp"])
	repo, _ := datastructures.CreateNewRepo(structuredRepoData["repoName"], structuredRepoData["author"], timestamp, nil)

	// getting the repo branches
	branchQueryString := fmt.Sprintf("{\"selector\": {\"repoID\": \"%s\"},\"fields\": [\"repoID\", \"branchName\", \"author\", \"timeStamp\"], \"sort\": [{\"timeStamp\": \"asc\"}]}", repoHash)
	branchResultsIterator, err := stub.GetQueryResult(branchQueryString)
	if err != nil {
		var repo datastructures.Repo
		return repo, err
	}
	defer branchResultsIterator.Close()
	//iterating over branches
	for branchResultsIterator.HasNext() {
		branchString, err := branchResultsIterator.Next()
		if err != nil {
			var repo datastructures.Repo
			return repo, err
		}

		structuredBranchData := map[string]string{}
		err = json.Unmarshal([]byte(branchString.Value), &structuredBranchData)
		if err != nil {
			var repo datastructures.Repo
			return repo, errors.New("Could not unmarashal requested Branch")
		}
		branchTimestamp, _ := strconv.Atoi(structuredBranchData["timestamp"])
		branch, _ := datastructures.CreateNewRepoBranch(structuredBranchData["branchName"], structuredBranchData["author"], branchTimestamp, nil)

		//adding branch commits
		commitsQueryString := fmt.Sprintf("{\"selector\": {\"repoID\": \"%s\", \"branchName\": \"%s\"},\"fields\": [\"repoID\", \"branchName\", \"message\", \"author\", \"committer\", \"committerTimestamp\", \"CommitParenthashes\", \"signature\"], \"sort\": [{\"timeStamp\": \"asc\"}]}", repoHash, branch.Name)
		commitsResultsIterator, err := stub.GetQueryResult(commitsQueryString)
		if err != nil {
			var repo datastructures.Repo
			return repo, err
		}
		defer commitsResultsIterator.Close()
		//iterating over branches
		for commitsResultsIterator.HasNext() {
			commitString, err := commitsResultsIterator.Next()
			if err != nil {
				var repo datastructures.Repo
				return repo, err
			}

			structuredCommitData := map[string]string{}
			err = json.Unmarshal([]byte(commitString.Value), &structuredCommitData)
			if err != nil {
				var repo datastructures.Repo
				return repo, errors.New("Could not unmarashal requested commit")
			}
			committerTimestamp, _ := strconv.Atoi(structuredCommitData["committerTimestamp"])
			var ph []string
			_ = json.Unmarshal([]byte(structuredCommitData["parenthashes"]), &ph)
			var s []byte
			_ = json.Unmarshal([]byte(structuredCommitData["signature"]), &s)
			commit, _ := datastructures.CreateNewCommitLog(structuredCommitData["message"], structuredCommitData["author"], structuredCommitData["commiter"], committerTimestamp, structuredCommitData["hash"], ph, s)
			branch.AddCommitLog(commit)
		}

		repo.AddBranch(branch)
	}

	// buffer, err := constructQueryResponseFromIterator(resultsIterator)
	// if err != nil {
	// 	return nil, err
	// }

	// create repoInstance

	return repo, nil
}

func (contract *RepoContract) getRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName

	fmt.Println("Querying the ledger..")

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

// indexName := "index-Branch"
// branchIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, args[2]})
// // jsonKey, _ := json.Marshal(key)
// branchData, err := stub.GetState(string(branchIndexKey))
