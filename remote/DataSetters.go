package main

import (
	"fmt"
	"strconv"

	client "github.com/GitofTruth/GoT/client"
	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func applyPair(stub shim.ChaincodeStubInterface, pair LedgerPair) bool {
	fmt.Println("Key:\t" + string(pair.key))
	fmt.Println("Value:\t" + string(pair.value))

	stub.PutState(string(pair.key), pair.value)
	return true
}

func applyPairs(stub shim.ChaincodeStubInterface, pairs []LedgerPair) bool {
	for ind, pair := range pairs {
		fmt.Println("Adding index:\t", ind)
		applyPair(stub, pair)
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

	accessPairs, _ := GenerateRepoUserAccessesDBPair(stub, repo)
	applyPairs(stub, accessPairs)

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

func (contract *RepoContract) addUserUpdate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args[0] >> UserMessage with content as UserUpdate

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userMessage, err := client.UnmarashalUserMessage(args[0])
	if err != nil {
		return shim.Error("userMessage is invalid!")
	}

	// TODO: check userDoesnot exist

	// check signature
	userUpdate, err := client.UnmarashalUserUpdate(userMessage.Content)
	if err != nil {
		return shim.Error("userUpdate is invalid!")
	}

	userInfo, _ := contract.getUserInfo(stub, userUpdate.UserName)
	if userUpdate.UserUpdateType == client.CreateNewUser && userInfo.UserName != "" {
		return shim.Error("User already exist")
	}

	// TODO: publicKey retrieval
	pubKey := userUpdate.PublicKey
	userNameMatchingNoChange := (userMessage.UserName == userUpdate.UserName) || (userUpdate.UserUpdateType != client.ChangeUserUserName)
	userNameMatchingChange := (userMessage.UserName == userUpdate.OldUserName) && (userUpdate.UserUpdateType == client.ChangeUserUserName)
	userNameMatching := userNameMatchingNoChange || userNameMatchingChange
	if userMessage.VerifySignature(pubKey) && userNameMatching {
		pairs, _ := GenerateUserUpdateDBPairs(stub, userUpdate)
		applyPairs(stub, pairs)
	}

	return shim.Success(nil)
}

func (contract *RepoContract) updateRepoUserAccess(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, authorized, userAccess, authorizer, encryptionKey/nil

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	fmt.Println("trying to process:\t", args)

	repo, err := contract.getRepoInstance(stub, args)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	// TODO: check userDoesnot exist
	access, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("could not parse access")
	}
	if repo.UpdateAccess(args[2], datastructures.UserAccess(access), args[4], args[5]) {
		repoPairs, _ := GenerateRepoDBPair(stub, repo)
		applyPairs(stub, repoPairs)
		pair, _ := GenerateRepoUserAccessDBPair(stub, args[0], args[1], args[2], args[3], args[4])
		applyPair(stub, pair)
		return shim.Success(nil)
	}

	return shim.Error("UserAccess was not set!")
}

// func(contract *RepoContract) NewSetterFunction(stub shim.ChaincodeStubInterface, args []string) peer.Response {
//   // Input assumption
//   // Example:
//   // Required Tables
//   // Index(ex) Used and primary key

//   // Set entries in tables

//     //data encoding

//     // Composite keys encoding

//     // Set the data

// }
