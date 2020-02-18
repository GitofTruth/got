package main

import (
	"github.com/GitofTruth/GoT/datastructures"

	"crypto/sha256"
	"encoding/json"
)

func (repo *datastructures.Repo) GenerateRepoDBPair() (string, []byte) {

	repoHash := GetRepoKey(repo.author, repo.Name)

	repoData := map[string]interface{}{"repoName": repo.Name, "author": repo.Author, "timestamp": repo.Timestamp}

	jsonData, _ := json.Marshal(repoData)

	return repoHash, repoData
}

func GetRepoKey(author string, repoName string) string {

	data := map[string]interface{}{"repoName": repoName, "author": author}
	js, _ := json.Marshal(data)

	repoHash := sha256.New()
	repoHash.Write(js)

	return repoHash.Sum(nil)
}
