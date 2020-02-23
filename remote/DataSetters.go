package main

import (
	"fmt"

	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func applyPairs(stub shim.ChaincodeStubInterface, pairs []LedgerPair) bool {
	for _, pair := range pairs {

		fmt.Println("Key: "+ string(pair.key))
		fmt.Println("Value: " + string(pair.value))

		stub.PutState(string(pair.key), pair.value)
	}
	return true
}

func (contract *RepoContract) addNewRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("trying to process:\t", args[0])

	repo, err := datastructures.UnmarashalRepo(args[0])
	if err != nil {
		return shim.Error("Repo is invalid!")
	}

	repoPairs, _ := GenerateRepoDBPair(stub, repo)
	applyPairs(stub, repoPairs)

	branchPairs, _ := GenerateRepoBranchesDBPair(stub, repo)
	applyPairs(stub, branchPairs)

	branchCommitPairs, _ := GenerateRepoBranchesCommitsDBPair(stub, repo)
	applyPairs(stub, branchCommitPairs)

	return shim.Success(nil)
}

// func(contract *RepoContract) addNewBranch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
//   // Input assumption
//   // Example:
//   // Required Tables
//   // Index(ex) Used and primary key

//   // Set entries in tables

//     //data encoding

//     // Composite keys encoding

//     // Set the data

// }
