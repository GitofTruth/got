package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type RepoContract struct {
}

//how to pass variables for initialization?
func (contract *RepoContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("initializing ledger")
	pushNumber := 0
	pushNumberBytes, _ := json.Marshal(pushNumber)
	APIstub.PutState("PushNumber", pushNumberBytes)
	fmt.Println("Ledger initalized push number is ")
	fmt.Println(pushNumberBytes)
	return shim.Success(nil)
}

func (contract *RepoContract) getCurrentRepoState(APIstub shim.ChaincodeStubInterface) (datastructures.Repo, int) {
	repo, _ := datastructures.CreateNewRepo("", 0, nil)
	master, _ := datastructures.CreateNewRepoBranch("master", "client", 0, nil)
	repo.AddBranch(master)
	pushes := contract.getAllPushes(APIstub)

	for _, push := range pushes {
		repo.AddCommitLogs(push.Logs, push.BranchName)
	}

	pushNumberBytes, _ := APIstub.GetState("PushNumber")
	var pushNumber int
	json.Unmarshal(pushNumberBytes, &pushNumber)

	return repo, pushNumber
}

func (contract *RepoContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "push" {
		return contract.addPush(APIstub, args)
	} else if function == "getBetween" {
		return contract.getPushes(APIstub, args)
	} else if function == "addBranch" {
		return contract.addBranch(APIstub, args)
	} else if function == "getBranches" {
		return contract.getBranches(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

//pushAsBytes, _ := json.Marshal(args[1])
//APIstub.PutState(args[0], pushAsBytes)

// carAsBytes, _ := APIstub.GetState(args[0])
// return shim.Success(carAsBytes)

//need to generate hash as key instead of just the same object
func (contract *RepoContract) addPush(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	repo, pushNumber := contract.getCurrentRepoState(APIstub)

	//it should be Marshaled on submission
	//pushAsBytes, _ := json.Marshal(args[1]) //changes from json object to bytes (string)

	var pushLog datastructures.PushLog
	fmt.Println("trying to process:\t", args[0])
	json.Unmarshal([]byte(args[0]), &pushLog)

	if done, _ := repo.AddCommitLogs(pushLog.Logs, pushLog.BranchName); done {
		startKeyBytes, _ := json.Marshal(pushNumber)
		APIstub.PutState(string(startKeyBytes), []byte(args[0]))
		pushNumber = pushNumber + 1
		pushNumberBytes, _ := json.Marshal(pushNumber)
		APIstub.PutState("PushNumber", pushNumberBytes)
		return shim.Success(nil)
	}

	return shim.Error("Invalid push Log!")
}

func (contract *RepoContract) addBranch(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	repo, _ := contract.getCurrentRepoState(APIstub)

	var branch datastructures.RepoBranch
	json.Unmarshal([]byte(args[0]), &branch)
	fmt.Println("unmarshaling done!")
	fmt.Println(branch)

	if done, _ := repo.AddBranch(branch); done {
		fmt.Println("New branch added:\t" + branch.Name)
		APIstub.PutState(branch.Name, []byte(args[0]))
		return shim.Success(nil)
	}

	return shim.Error("Invalid Branch")
}

//
// func (contract *RepoContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	return shim.Success(nil)
// }

func (contract *RepoContract) getAllPushes(APIstub shim.ChaincodeStubInterface) []datastructures.PushLog {

	startKeyBytes, _ := json.Marshal(0)
	endKeyBytes, _ := APIstub.GetState("PushNumber")

	resultsIterator, err := APIstub.GetStateByRange(string(startKeyBytes), string(endKeyBytes))
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

func (contract *RepoContract) getPushes(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("Querying the ledger..")

	startKey, _ := strconv.Atoi(args[0])
	endKey, _ := strconv.Atoi(args[1])

	startKeyBytes, _ := json.Marshal(startKey)
	endKeyBytes, _ := json.Marshal(endKey)

	resultsIterator, err := APIstub.GetStateByRange(string(startKeyBytes), string(endKeyBytes))
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

func (contract *RepoContract) getBranches(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("Querying the ledger..")
	repo, _ := contract.getCurrentRepoState(APIstub)

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
