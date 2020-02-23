package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func (contract *RepoContract) getRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("Querying the ledger..")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	repoHash := GetRepoKey(args[0], args[1])

	// key := map[string]interface{}{"repoID": repoHash}
	// jsonKey, _ := json.Marshal(key)
	repoData, err := stub.GetState(string(repoHash))

	if err != nil {
		return shim.Error("Repo does not exist")
	}

	indexName := "index-Branch"
	branchIndexKey, _ := stub.CreateCompositeKey(indexName, []string{repoHash, args[2]})


	// jsonKey, _ := json.Marshal(key)

	branchData, err := stub.GetState(string(branchIndexKey))



	fmt.Println("Found this repo:", string(repoData))

	return shim.Success(branchData)
}
