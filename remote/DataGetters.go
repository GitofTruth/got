package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func (contract *RepoContract) getRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName

	fmt.Println("Querying the ledger..")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	repoHash := GetRepoKey(args[0], args[1])

	key := map[string]interface{}{"repoID": repoHash}
	jsonKey, _ := json.Marshal(key)
	repoData, err := stub.GetState(string(jsonKey))

	if err != nil {
		return shim.Error("Repo does not exist")
	}

	fmt.Println("Found this repo:", repoData)

	return shim.Success(repoData)
}
