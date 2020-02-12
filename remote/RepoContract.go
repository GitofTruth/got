package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type RepoContract struct {
}

//how to pass variables for initialization?
func (contract *RepoContract) Init(stub shim.ChaincodeStubInterface) peer.Response {

	fmt.Println("initializing ledger")

	// Add none repo?


	// pushNumber := 0
	// pushNumberBytes, _ := json.Marshal(pushNumber)
	// stub.PutState("PushNumber", pushNumberBytes)
	// fmt.Println("Ledger initalized push number is ")
	// fmt.Println(pushNumberBytes)

	return shim.Success(nil)
}


func (contract *RepoContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "addNewRepo" {
		return contract.addNewRepo(stub, args)
	} else if function == "addNewBranch" {
		return contract.addNewBranch(stub, args)
	} else if function == "addCommits" {
		return contract.addCommits(stub, args)
	} else if function == "addCollabrator" {
		return contract.addCollabrator(stub, args)
	}

	else if function == "queryRepo" {
		return contract.queryRepo(stub, args)
	} else if function == "queryBranch" {
		return contract.queryBranch(stub, args)
	} else if function == "queryRepoHashes" {
		return contract.queryRepoHashes(stub, args)
	} else if function == "queryCommits" {
		return contract.queryCommits(stub, args)
	} else if function == "queryCollabrators" {
		return contract.queryCollabrators(stub, args)
	}


	return shim.Error("Invalid Smart Contract function name.")
}

func(contract *RepoContract) addNewRepo(stub shim.ChaincodeStubInterface) peer.Response {

	
}

func (contract *RepoContract) getCurrentRepoState(stub shim.ChaincodeStubInterface) (datastructures.Repo, int) {
	repo, _ := datastructures.CreateNewRepo("", 0, nil)
	master, _ := datastructures.CreateNewRepoBranch("master", "client", 0, nil)
	repo.AddBranch(master)
	pushes := contract.getAllPushes(stub)

	for _, push := range pushes {
		repo.AddCommitLogs(push.Logs, push.BranchName)
	}

	pushNumberBytes, _ := stub.GetState("PushNumber")
	var pushNumber int
	json.Unmarshal(pushNumberBytes, &pushNumber)

	return repo, pushNumber
}



//pushAsBytes, _ := json.Marshal(args[1])
//stub.PutState(args[0], pushAsBytes)

// carAsBytes, _ := stub.GetState(args[0])
// return shim.Success(carAsBytes)

//need to generate hash as key instead of just the same object
func (contract *RepoContract) addPush(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	repo, pushNumber := contract.getCurrentRepoState(stub)

	//it should be Marshaled on submission
	//pushAsBytes, _ := json.Marshal(args[1]) //changes from json object to bytes (string)

	var pushLog datastructures.PushLog
	fmt.Println("trying to process:\t", args[0])
	json.Unmarshal([]byte(args[0]), &pushLog)

	if done, _ := repo.AddCommitLogs(pushLog.Logs, pushLog.BranchName); done {
		startKeyBytes, _ := json.Marshal(pushNumber)
		stub.PutState(string(startKeyBytes), []byte(args[0]))

		pushNumber = pushNumber + 1
		pushNumberBytes, _ := json.Marshal(pushNumber)
		stub.PutState("PushNumber", pushNumberBytes)

		return shim.Success(nil)
	}

	return shim.Error("Invalid push Log!")
}

func (contract *RepoContract) addBranch(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	repo, _ := contract.getCurrentRepoState(stub)

	var branch datastructures.RepoBranch
	json.Unmarshal([]byte(args[0]), &branch)
	fmt.Println("unmarshaling done!")
	fmt.Println(branch)

	if done, _ := repo.AddBranch(branch); done {
		fmt.Println("New branch added:\t" + branch.Name)
		stub.PutState(branch.Name, []byte(args[0]))
		return shim.Success(nil)
	}

	return shim.Error("Invalid Branch")
}

//
// func (contract *RepoContract) initLedger(stub shim.ChaincodeStubInterface) peer.Response {
// 	return shim.Success(nil)
// }

func (contract *RepoContract) getAllPushes(stub shim.ChaincodeStubInterface) []datastructures.PushLog {

	startKeyBytes, _ := json.Marshal(0)
	endKeyBytes, _ := stub.GetState("PushNumber")

	resultsIterator, err := stub.GetStateByRange(string(startKeyBytes), string(endKeyBytes))
	if err != nil {
		return make([]datastructures.PushLog, 0)
	}

	var pushlogs []datastructures.PushLog
	pushlogs = make([]datastructures.PushLog, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return make([]datastructures.PushLog, 0)
		}

		var pushlog datastructures.PushLog
		json.Unmarshal(queryResponse.Value, &pushlog)

		pushlogs = append(pushlogs, pushlog)
	}

	defer resultsIterator.Close()
	return pushlogs
}

func (contract *RepoContract) getPushes(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("Querying the ledger..")

	startKey, _ := strconv.Atoi(args[0])
	endKey, _ := strconv.Atoi(args[1])

	startKeyBytes, _ := json.Marshal(startKey)
	endKeyBytes, _ := json.Marshal(endKey)

	resultsIterator, err := stub.GetStateByRange(string(startKeyBytes), string(endKeyBytes))
	if err != nil {
		return shim.Error(err.Error())
	}

	var pushlogs []datastructures.PushLog
	pushlogs = make([]datastructures.PushLog, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var pushlog datastructures.PushLog
		json.Unmarshal(queryResponse.Value, &pushlog)

		fmt.Println("trying to query:\t", pushlog)

		pushlogs = append(pushlogs, pushlog)
	}

	defer resultsIterator.Close()

	fmt.Println("pushlogs before json marshalling: ", pushlogs)

	pushlogsjson, _ := json.Marshal(pushlogs)

	return shim.Success(pushlogsjson)
}

func (contract *RepoContract) getBranches(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("Querying the ledger..")
	repo, _ := contract.getCurrentRepoState(stub)

	brancgesjson, _ := json.Marshal(repo.GetBranches())

	fmt.Println(brancgesjson)
	return shim.Success(brancgesjson)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(RepoContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
