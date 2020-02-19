package main

import (
  "github.com/GitofTruth/GoT/datastructures"
  
  "fmt"

	"crypto/sha256"
	"encoding/json"
)

func GenerateRepoDBPair(repo datastructures.Repo) (string, []byte) {

	repoHash := GetRepoKey(repo.Author, repo.Name)

	repoData := map[string]interface{}{"repoName": repo.Name, "author": repo.Author, "timestamp": repo.Timestamp}

	jsonData, _ := json.Marshal(repoData)

	return repoHash, jsonData
}

func GetRepoKey(author string, repoName string) string {

	data := map[string]interface{}{"repoName": repoName, "author": author}
	js, _ := json.Marshal(data)

	repoHash := sha256.New()
  repoHash.Write(js)
  
  fmt.Println("Repo Hash: " + string(repoHash.Sum(nil)))

	return string(repoHash.Sum(nil))
}
