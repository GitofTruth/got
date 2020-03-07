package main

import (
	"github.com/GitofTruth/GoT/datastructures"

	"strconv"

	b64 "encoding/base64"
	"fmt"

	"crypto/sha256"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type LedgerPair struct {
	key   string
	value []byte
}

func GenerateRepoDBPair(stub shim.ChaincodeStubInterface, repo datastructures.Repo) ([]LedgerPair, error) {

	repoHash := GetRepoKey(repo.Author, repo.Name)

	list := make([]LedgerPair, 0)
	var pair LedgerPair

	// key := map[string]interface{}{"repoID": repoHash}
	// jsonKey, _ := json.Marshal(key)

	pair.key = string(repoHash)

	value := map[string]interface{}{"docName": "repo", "repoID": repoHash, "repoName": repo.Name, "author": repo.Author, "timeStamp": strconv.Itoa(repo.Timestamp)}
	pair.value, _ = json.Marshal(value)

	list = append(list, pair)

	return list, nil
}

func GenerateRepoBranchDBPair(stub shim.ChaincodeStubInterface, author string, repoName string, branch datastructures.RepoBranch) (LedgerPair, error) {

	repoHash := GetRepoKey(author, repoName)

	var pair LedgerPair

	// key := map[string]interface{}{"repoID": repoHash, "branchName": branch.Name}
	indexName := "index-Branch"
	branchIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, branch.Name})

	// jsonKey, _ := json.Marshal(key)

	fmt.Println("branchIndexKey : " + branchIndexKey)
	pair.key = string(branchIndexKey)

	value := map[string]interface{}{"docName": "branch", "repoID": repoHash, "branchName": branch.Name, "author": branch.Author, "timeStamp": strconv.Itoa(branch.Timestamp)}
	pair.value, _ = json.Marshal(value)

	return pair, nil
}

func GenerateRepoBranchesDBPair(stub shim.ChaincodeStubInterface, repo datastructures.Repo) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for _, branch := range repo.Branches {
		pair, _ := GenerateRepoBranchDBPair(stub, repo.Author, repo.Name, branch)
		list = append(list, pair)
	}

	return list, nil
}

func GenerateRepoBranchCommitDBPair(stub shim.ChaincodeStubInterface, author string, repoName string, branchName string, commitLog datastructures.CommitLog) (LedgerPair, error) {

	repoHash := GetRepoKey(author, repoName)

	var pair LedgerPair

	indexName := "index-BranchCommits"
	branchCommitIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, branchName, commitLog.Hash})

	// key := map[string]interface{}{"repoID": repoHash, "branchName": branchName, "hash": commitLog.Hash}
	// jsonKey, _ := json.Marshal(key)
	pair.key = string(branchCommitIndexKey)

	value := map[string]interface{}{"docName": "commit", "repoID": repoHash, "branchName": branchName, "hash": commitLog.Hash, "message": commitLog.Message, "author": commitLog.Author, "committer": commitLog.Committer, "committerTimestamp": strconv.Itoa(commitLog.CommitterTimestamp), "parentHashes": commitLog.Parenthashes, "signature": commitLog.Signature}
	pair.value, _ = json.Marshal(value)

	return pair, nil
}

func GenerateRepoBranchesCommitsDBPair(stub shim.ChaincodeStubInterface, repo datastructures.Repo) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)
	for _, branch := range repo.Branches {
		for _, log := range branch.Logs {

			pair, _ := GenerateRepoBranchCommitDBPair(stub, repo.Author, repo.Name, branch.Name, log)
			list = append(list, pair)
		}
	}

	return list, nil
}

func GenerateRepoBranchesCommitsDBPairUsingBranch(stub shim.ChaincodeStubInterface, author string, repoName string, repoBranch datastructures.RepoBranch) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for _, log := range repoBranch.Logs {
		pair, _ := GenerateRepoBranchCommitDBPair(stub, author, repoName, repoBranch.Name, log)
		list = append(list, pair)
	}

	return list, nil
}

func GenerateRepoBranchesCommitsDBPairUsingPushLog(stub shim.ChaincodeStubInterface, author string, repoName string, pushLog datastructures.PushLog) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for _, log := range pushLog.Logs {
		pair, _ := GenerateRepoBranchCommitDBPair(stub, author, repoName, pushLog.BranchName, log)
		list = append(list, pair)
	}

	return list, nil
}

func GetRepoKey(author string, repoName string) string {

	data := map[string]interface{}{"repoName": repoName, "author": author}
	js, _ := json.Marshal(data)

	repoHash := sha256.New()
	repoHash.Write(js)

	fmt.Println("Repo Hash: ", repoHash.Sum(nil))
	sEnc := b64.StdEncoding.EncodeToString([]byte(repoHash.Sum(nil)))
	fmt.Println("Repo Hash: ", sEnc)

	// fmt.Println("Repo String Hash: " + string(repoHash.Sum(nil)))
	//
	// keyBytes, _ := json.Marshal(string(repoHash.Sum(nil)))
	// fmt.Println("Repo String Hash keyBytes: ", keyBytes)
	// fmt.Println("Repo String Hash keyBytes String: " + string(keyBytes))

	return sEnc
}
