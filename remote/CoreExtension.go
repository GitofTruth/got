package main

import (
	"github.com/GitofTruth/GoT/datastructures"

	"fmt"

	"crypto/sha256"
	"encoding/json"
)

type LedgerPair struct {
	key   string
	value []byte
}

func GenerateRepoDBPair(repo datastructures.Repo) ([]LedgerPair, error) {

	repoHash := GetRepoKey(repo.Author, repo.Name)

	list := make([]LedgerPair, 0)
	var pair LedgerPair

	key := map[string]interface{}{"repoID": repoHash}
	jsonKey, _ := json.Marshal(key)
	pair.key = string(jsonKey)

	value := map[string]interface{}{"repoID": repoHash, "repoName": repo.Name, "author": repo.Author, "timestamp": repo.Timestamp, "hashes": repo.CommitHashes}
	pair.value, _ = json.Marshal(value)

	list = append(list, pair)

	return list, nil
}

func GenerateRepoBranchDBPair(author string, repoName string, branch datastructures.RepoBranch) (LedgerPair, error) {

	repoHash := GetRepoKey(author, repoName)

	var pair LedgerPair

	key := map[string]interface{}{"repoID": repoHash, "branchName": branch.Name}
	jsonKey, _ := json.Marshal(key)
	pair.key = string(jsonKey)

	value := map[string]interface{}{"repoID": repoHash, "branchName": branch.Name, "author": branch.Author, "timeStamp": branch.Timestamp}
	pair.value, _ = json.Marshal(value)

	return pair, nil
}

func GenerateRepoBranchesDBPair(repo datastructures.Repo) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for _, branch := range repo.Branches {
		pair, _ := GenerateRepoBranchDBPair(repo.Author, repo.Name, branch)
		list = append(list, pair)
	}

	return list, nil
}

func GenerateRepoBranchCommitDBPair(author string, repoName string, branchName string, commitLog datastructures.CommitLog) (LedgerPair, error) {

	repoHash := GetRepoKey(author, repoName)

	var pair LedgerPair

	key := map[string]interface{}{"repoID": repoHash, "branchName": branchName, "hash": commitLog.Hash}
	jsonKey, _ := json.Marshal(key)
	pair.key = string(jsonKey)

	value := map[string]interface{}{"repoID": repoHash, "branchName": branchName, "hash": commitLog.Hash, "message": commitLog.Message, "author": commitLog.Author, "committer": commitLog.Committer, "committerTimestamp": commitLog.CommitterTimestamp, "parenthashes": commitLog.Parenthashes, "signature": commitLog.Signature}
	pair.value, _ = json.Marshal(value)

	return pair, nil
}

func GenerateRepoBranchesCommitsDBPair(repo datastructures.Repo) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)
	for _, branch := range repo.Branches {
		for _, log := range branch.Logs {

			pair, _ := GenerateRepoBranchCommitDBPair(repo.Author, repo.Name, branch.Name, log)
			list = append(list, pair)
		}
	}

	return list, nil
}

func GetRepoKey(author string, repoName string) string {

	data := map[string]interface{}{"repoName": repoName, "author": author}
	js, _ := json.Marshal(data)

	repoHash := sha256.New()
	repoHash.Write(js)

	fmt.Println("Repo Hash: " + string(repoHash.Sum(nil)))

	keyBytes, _ := json.Marshal(string(repoHash.Sum(nil)))

	return string(keyBytes)
}
