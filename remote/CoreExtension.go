package remote

import (
	client "github.com/GitofTruth/got/client"
	"github.com/GitofTruth/got/datastructures"

	"strconv"
	"time"

	b64 "encoding/base64"
	"fmt"

	"crypto/sha256"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// This is the resulting key-value pair generated to store some the data
// in the hyperledger after the maping from complex data structurs is done.
type LedgerPair struct {
	key   string
	value []byte
}

func generateRepoDBPair(stub shim.ChaincodeStubInterface, repo datastructures.Repo) ([]LedgerPair, error) {

	repoHash := getRepoKey(repo.Author, repo.Name)

	list := make([]LedgerPair, 0)
	var pair LedgerPair

	// key := map[string]interface{}{"repoID": repoHash}
	// jsonKey, _ := json.Marshal(key)

	pair.key = string(repoHash)
	accessLogs, _ := json.Marshal(repo.AccessLogs)
	encryptionKey, _ := json.Marshal(repo.EncryptionKey)
	keyAnnouncements, _ := json.Marshal(repo.KeyAnnouncements)

	// fmt.Println("\n\nbos el awal hena")
	// fmt.Println(repo.KeyAnnouncements)
	// fmt.Println(keyAnnouncements)
	// keyAnnouncementsss, _ := datastructures.UnmarashalKeyAnnouncements(string(keyAnnouncements))
	// fmt.Println(keyAnnouncementsss)
	// fmt.Println("Kefaya bas y ostaz\n\n")

	value := map[string]interface{}{"docName": "repo", "repoID": repoHash, "repoName": repo.Name,
		"author": repo.Author, "directoryCID": repo.DirectoryCID,
		"timeStamp": strconv.Itoa(repo.Timestamp), "encryptionKey": string(encryptionKey),
		"accessLogs": string(accessLogs), "keyAnnouncements": string(keyAnnouncements)}

	pair.value, _ = json.Marshal(value)

	list = append(list, pair)

	return list, nil
}

func generateRepoBranchDBPair(stub shim.ChaincodeStubInterface, author string, repoName string, branch datastructures.RepoBranch) (LedgerPair, error) {

	repoHash := getRepoKey(author, repoName)

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

func generateRepoBranchesDBPair(stub shim.ChaincodeStubInterface, repo datastructures.Repo) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for _, branch := range repo.Branches {
		pair, _ := generateRepoBranchDBPair(stub, repo.Author, repo.Name, branch)
		list = append(list, pair)
	}

	return list, nil
}

func generateRepoBranchCommitDBPair(stub shim.ChaincodeStubInterface, author string, repoName string, branchName string, commitLog datastructures.CommitLog) (LedgerPair, error) {

	repoHash := getRepoKey(author, repoName)

	var pair LedgerPair

	indexName := "index-BranchCommits"
	branchCommitIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, branchName, commitLog.Hash})

	// key := map[string]interface{}{"repoID": repoHash, "branchName": branchName, "hash": commitLog.Hash}
	// jsonKey, _ := json.Marshal(key)
	pair.key = string(branchCommitIndexKey)

	storageHashes, _ := json.Marshal(commitLog.StorageHashes)
	value := map[string]interface{}{"docName": "commit", "repoID": repoHash, "branchName": branchName, "hash": commitLog.Hash, "message": commitLog.Message, "author": commitLog.Author, "committer": commitLog.Committer, "committerTimestamp": strconv.Itoa(commitLog.CommitterTimestamp), "parentHashes": commitLog.Parenthashes, "signature": commitLog.Signature, "encryptionKey": commitLog.EncryptionKey, "storageHashes": string(storageHashes)}
	pair.value, _ = json.Marshal(value)

	return pair, nil
}

func generateRepoBranchesCommitsDBPair(stub shim.ChaincodeStubInterface, repo datastructures.Repo) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)
	for _, branch := range repo.Branches {
		for _, log := range branch.Logs {
			pair, _ := generateRepoBranchCommitDBPair(stub, repo.Author, repo.Name, branch.Name, log)
			list = append(list, pair)
		}
	}

	return list, nil
}

func generateRepoBranchesCommitsDBPairUsingBranch(stub shim.ChaincodeStubInterface, author string, repoName string, repoBranch datastructures.RepoBranch) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for _, log := range repoBranch.Logs {
		pair, _ := generateRepoBranchCommitDBPair(stub, author, repoName, repoBranch.Name, log)
		list = append(list, pair)
	}

	return list, nil
}

func generateRepoBranchesCommitsDBPairUsingPushLog(stub shim.ChaincodeStubInterface, author string, repoName string, pushLog datastructures.PushLog) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for _, log := range pushLog.Logs {
		pair, _ := generateRepoBranchCommitDBPair(stub, author, repoName, pushLog.BranchName, log)
		list = append(list, pair)
	}

	return list, nil
}

func getRepoKey(author string, repoName string) string {

	data := map[string]interface{}{"repoName": repoName, "author": author}
	js, _ := json.Marshal(data)

	repoHash := sha256.New()
	repoHash.Write(js)

	// fmt.Println("Repo Hash: ", repoHash.Sum(nil))
	sEnc := b64.StdEncoding.EncodeToString([]byte(repoHash.Sum(nil)))
	fmt.Println("Repo Hash: ", sEnc)

	// fmt.Println("Repo String Hash: " + string(repoHash.Sum(nil)))
	//
	// keyBytes, _ := json.Marshal(string(repoHash.Sum(nil)))
	// fmt.Println("Repo String Hash keyBytes: ", keyBytes)
	// fmt.Println("Repo String Hash keyBytes String: " + string(keyBytes))

	return sEnc
}

func generateUserUpdateDBPairs(stub shim.ChaincodeStubInterface, userUpdate client.UserUpdate) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)
	fmt.Println(userUpdate.UserUpdateType)
	if userUpdate.UserUpdateType == client.CreateNewUser || userUpdate.UserUpdateType == client.ChangeUserPublicKey {
		var pair LedgerPair
		indexName := "index-user"

		userIndexKey, _ := stub.CreateCompositeKey(indexName, []string{userUpdate.UserName})
		pair.key = string(userIndexKey)
		fmt.Println("userIndexKey : \n" + pair.key + "\n")

		pubKey := userUpdate.PublicKey.(string)
		value := map[string]interface{}{"docName": "user", "userName": userUpdate.UserName, "publicKey": pubKey}
		pair.value, _ = json.Marshal(value)

		list = append(list, pair)
	} else if userUpdate.UserUpdateType == client.DeleteUser {
		var pair LedgerPair
		indexName := "index-user"

		userIndexKey, _ := stub.CreateCompositeKey(indexName, []string{userUpdate.UserName})
		pair.key = string(userIndexKey)
		fmt.Println("userIndexKey : \n" + pair.key + "\n")

		value := map[string]interface{}{"docName": "user", "userName": "", "publicKey": ""}
		pair.value, _ = json.Marshal(value)

		list = append(list, pair)
	} else if userUpdate.UserUpdateType == client.ChangeUserUserName {
		var pair LedgerPair
		indexName := "index-user"

		//removing old user
		userIndexKey, _ := stub.CreateCompositeKey(indexName, []string{userUpdate.OldUserName})
		pair.key = string(userIndexKey)

		value := map[string]interface{}{"docName": "user", "userName": "", "publicKey": ""}
		pair.value, _ = json.Marshal(value)

		list = append(list, pair)

		//adding new user
		userIndexKey, _ = stub.CreateCompositeKey(indexName, []string{userUpdate.UserName})
		pair.key = string(userIndexKey)
		fmt.Println("userIndexKey : \t" + pair.key)

		pubKey := userUpdate.PublicKey.(string)
		value = map[string]interface{}{"docName": "user", "userName": userUpdate.UserName, "publicKey": pubKey}
		pair.value, _ = json.Marshal(value)

		// TODO: update all the user has access to :D

		list = append(list, pair)
	}

	return list, nil
}

func generateRepoUserAccessDBPair(stub shim.ChaincodeStubInterface, author string, repoName string, authorized string, userAccess string, authorizer string) (LedgerPair, error) {

	repoHash := getRepoKey(author, repoName)

	var pair LedgerPair

	// key := map[string]interface{}{"repoID": repoHash, "branchName": branch.Name}
	indexName := "repoUserAccess-Branch"
	repoUserAccessIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, authorized})

	fmt.Println(indexName + " : \n" + repoUserAccessIndexKey)
	pair.key = string(repoUserAccessIndexKey)

	value := map[string]interface{}{"docName": "userAccess", "repoHash": repoHash, "authorized": authorized, "userAccess": userAccess, "authorizer": authorizer}
	pair.value, _ = json.Marshal(value)

	return pair, nil
}

func generateRepoUserAccessesDBPair(stub shim.ChaincodeStubInterface, repo datastructures.Repo) ([]LedgerPair, error) {

	list := make([]LedgerPair, 0)

	for authorized, userAccess := range repo.Users {
		pair, _ := generateRepoUserAccessDBPair(stub, repo.Author, repo.Name, authorized, strconv.Itoa(int(userAccess)), repo.Author)
		list = append(list, pair)
	}

	return list, nil
}

func checkCurrentTimeStamp(timeStamp int64) bool {
	now := time.Now()
	sec := now.Unix()

	return ((sec - timeStamp) < 259200) && ((sec - timeStamp) > -86400)
}
