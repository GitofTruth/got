package main

import (
	"errors"
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

func (contract *RepoContract) validateUserArgsMessage(stub shim.ChaincodeStubInterface, args []string, argsNumber int) (client.UserMessage, []string, error) {

	if len(args) != 1 {
		var userMessage client.UserMessage
		var innerArgs []string
		return userMessage, innerArgs, errors.New("Incorrect number of arguments. Expecting 1")
	}

	userMessage, err := client.UnmarashalUserMessage(args[0])
	if err != nil {
		var userMessage client.UserMessage
		var innerArgs []string
		return userMessage, innerArgs, errors.New("userMessage is invalid!")
	}

	innerArgs, err := client.UnmarashalInnerArgs(userMessage.Content)
	if err != nil {
		var userMessage client.UserMessage
		var innerArgs []string
		return userMessage, innerArgs, errors.New("innerArgs is invalid!")
	}

	if len(innerArgs) != argsNumber {
		var userMessage client.UserMessage
		var innerArgs []string
		return userMessage, innerArgs, errors.New("Incorrect number of inner arguments. Expecting " + strconv.Itoa(argsNumber))
	}

	// verifying signature
	// getting claimed user
	userInfo, failMessage := contract.getUserInfo(stub, userMessage.UserName)
	if failMessage.Message != "" {
		var userMessage client.UserMessage
		var innerArgs []string
		return userMessage, innerArgs, errors.New(failMessage.Message)
	}

	if !userMessage.VerifySignature(userInfo.PublicKey) {
		var userMessage client.UserMessage
		var innerArgs []string
		return userMessage, innerArgs, errors.New("User Signature is not valid")
	}

	return userMessage, innerArgs, nil
}

func (contract *RepoContract) addNewRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	userMessage, innerArgs, err := contract.validateUserArgsMessage(stub, args, 1)
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}

	repo, err := datastructures.UnmarashalRepo(innerArgs[0])
	if err != nil {
		return shim.Error("Repo is invalid!")
	}

	// checking that the creator is whom they claim to be
	if userMessage.UserName != repo.Author {
		return shim.Error("Repo creator is not the signing user")
	}

	// check if repo already exists
	repoArgsList := make([]string, 2)
	repoArgsList[0] = repo.Author
	repoArgsList[1] = repo.Name
	_, err = contract.getRepoInstance(stub, repoArgsList)
	if err == nil {
		return shim.Error("Repo already exists")
	}

	repoPairs, _ := GenerateRepoDBPair(stub, repo)
	applyPairs(stub, repoPairs)

	accessPairs, _ := GenerateRepoUserAccessesDBPair(stub, repo)
	applyPairs(stub, accessPairs)

	branchPairs, _ := GenerateRepoBranchesDBPair(stub, repo)
	applyPairs(stub, branchPairs)

	branchCommitPairs, _ := GenerateRepoBranchesCommitsDBPair(stub, repo)
	applyPairs(stub, branchCommitPairs)

	return shim.Success([]byte("The repo has been added successfully to the blockchain."))
}

func (contract *RepoContract) addNewBranch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, branchBinary

	userMessage, innerArgs, err := contract.validateUserArgsMessage(stub, args, 3)
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}

	repoBranch, err := datastructures.UnmarashalRepoBranch(innerArgs[2])
	if err != nil {
		return shim.Error("RepoBranch is invalid!")
	}

	// generate Repo & check validation
	repo, err := contract.getRepoInstance(stub, innerArgs)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	// check authorization
	isAuthorized := repo.CanEdit(userMessage.UserName)
	if !isAuthorized {
		return shim.Error("User is not authorized to edit this repo")
	}

	valid, err := repo.ValidBranch(repoBranch)
	if err != nil || !valid {
		return shim.Error("RepoBranch could not be added!")
	}

	branchPair, _ := GenerateRepoBranchDBPair(stub, innerArgs[0], innerArgs[1], repoBranch)
	applyPair(stub, branchPair)

	commitsPairs, _ := GenerateRepoBranchesCommitsDBPairUsingBranch(stub, innerArgs[0], innerArgs[1], repoBranch)
	applyPairs(stub, commitsPairs)

	return shim.Success([]byte("The branch has been added successfully to its corresponding repo!"))
}

func (contract *RepoContract) addCommits(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoName, repoAuthor, PushLogBinary

	userMessage, innerArgs, err := contract.validateUserArgsMessage(stub, args, 3)
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}

	pushLog, err := datastructures.UnmarashalPushLog(innerArgs[2])
	if err != nil {
		return shim.Error("PushLog is invalid!")
	}

	// generate Repo & check validation
	repo, err := contract.getRepoInstance(stub, innerArgs)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	// check authorization
	isAuthorized := repo.CanEdit(userMessage.UserName)
	if !isAuthorized {
		return shim.Error("User is not authorized to edit this repo")
	}

	if len(pushLog.Logs) < 1 {
		return shim.Error("Could not find any commit log")
	}

	branchDidNotExist := !repo.IsBranch(pushLog.BranchName)
	if branchDidNotExist {
		newbranch, _ := datastructures.CreateNewRepoBranch(pushLog.BranchName, userMessage.UserName, pushLog.Logs[0].CommitterTimestamp, nil)
		branchDidNotExist, _ = repo.AddBranch(newbranch)
	}

	valid, err := repo.AddCommitLogs(pushLog.Logs, pushLog.BranchName, false)
	if err != nil || !valid {
		return shim.Error("Logs could not be added!")
	}

	repoPairs, _ := GenerateRepoDBPair(stub, repo)
	applyPairs(stub, repoPairs)

	if branchDidNotExist {
		branchPair, _ := GenerateRepoBranchDBPair(stub, innerArgs[0], innerArgs[1], repo.Branches[pushLog.BranchName])
		applyPair(stub, branchPair)
	}

	commitsPairs, _ := GenerateRepoBranchesCommitsDBPairUsingPushLog(stub, innerArgs[0], innerArgs[1], pushLog)
	applyPairs(stub, commitsPairs)

	return shim.Success([]byte("The commits have been added successfully to the blockchain"))
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

	// check signature
	userUpdate, err := client.UnmarashalUserUpdate(userMessage.Content)
	if err != nil {
		return shim.Error("userUpdate is invalid!")
	}

	userInfo, _ := contract.getUserInfo(stub, userUpdate.UserName)
	if userUpdate.UserUpdateType == client.CreateNewUser && userInfo.UserName != "" {
		return shim.Error("User already exist")
	}

	pubKey := userUpdate.PublicKey
	if userInfo.UserName != "" {
		pubKey = userInfo.PublicKey
	}

	userNameMatchingNoChange := (userMessage.UserName == userUpdate.UserName) || (userUpdate.UserUpdateType != client.ChangeUserUserName)
	userNameMatchingChange := (userMessage.UserName == userUpdate.OldUserName) && (userUpdate.UserUpdateType == client.ChangeUserUserName)
	userNameMatching := userNameMatchingNoChange || userNameMatchingChange
	if userMessage.VerifySignature(pubKey) && userNameMatching {
		pairs, _ := GenerateUserUpdateDBPairs(stub, userUpdate)
		applyPairs(stub, pairs)
	}

	return shim.Success([]byte("The requested user update has been processed successfully!"))
}

func (contract *RepoContract) updateRepoUserAccess(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// repoAuthor, repoName, authorized, userAccess, authorizer, encryptionKey/nil

	userMessage, innerArgs, err := contract.validateUserArgsMessage(stub, args, 6)
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}

	if userMessage.UserName != innerArgs[4] {
		return shim.Error("authorizer is not the signing user")
	}

	repo, err := contract.getRepoInstance(stub, innerArgs)
	if err != nil {
		return shim.Error("Repo does not exist")
	}

	access, err := strconv.Atoi(innerArgs[3])
	if err != nil {
		return shim.Error("could not parse access")
	}

	retrievedEncKey, _ := datastructures.UnmarashalKeyAnnouncement(innerArgs[5])
	if repo.UpdateAccess(innerArgs[2], datastructures.UserAccess(access), innerArgs[4], retrievedEncKey) {
		repoPairs, _ := GenerateRepoDBPair(stub, repo)
		applyPairs(stub, repoPairs)
		pair, _ := GenerateRepoUserAccessDBPair(stub, innerArgs[0], innerArgs[1], innerArgs[2], innerArgs[3], innerArgs[4])
		applyPair(stub, pair)

		return shim.Success([]byte("Access to the repo has been updated successfully!"))
	}

	return shim.Error("UserAccess was not set! Your access type does not permit you to do the required task")
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
