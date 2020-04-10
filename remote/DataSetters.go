package main

import (
	"fmt"

	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func applyPair(stub shim.ChaincodeStubInterface, pair LedgerPair) bool {
	stub.PutState(pair.key, pair.value)
	return true
}

func applyPairs(stub shim.ChaincodeStubInterface, pairs []LedgerPair) bool {
	for _, pair := range pairs {

		fmt.Println("Key: " + string(pair.key))
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

func (contract *RepoContract) addNewBranch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, branchBinary

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fmt.Println("trying to process:\t", args)

	// TODO: needs read branches & commits
	// generate Repo & check validation

	repoBranch, err := datastructures.UnmarashalRepoBranch(args[2])
	if err != nil {
		return shim.Error("RepoBranch is invalid!")
	}

	branchPair, _ := GenerateRepoBranchDBPair(stub, args[0], args[1], repoBranch)
	applyPair(stub, branchPair)

	commitsPairs, _ := GenerateRepoBranchesCommitsDBPairUsingBranch(stub, args[0], args[1], repoBranch)
	applyPairs(stub, commitsPairs)

	return shim.Success(nil)
}

func (contract *RepoContract) addCommits(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoName, repoAuthor, PushLogBinary

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fmt.Println("trying to process:\t", args)

	// TODO: needs read branches & commits
	// generate Repo & check validation

	pushLog, err := datastructures.UnmarashalPushLog(args[2])
	if err != nil {
		return shim.Error("PushLog is invalid!")
	}

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	repo.DirectoryCID = pushLog.DirectoryCID

	repoPairs, _ := GenerateRepoDBPair(stub, repo)
	applyPairs(stub, repoPairs)


	commitsPairs, _ := GenerateRepoBranchesCommitsDBPairUsingPushLog(stub, args[0], args[1], pushLog)
	applyPairs(stub, commitsPairs)

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
