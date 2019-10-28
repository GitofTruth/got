package main

import (
	"encoding/json"
	"fmt"

	"github.com/GitofTruth/GoT/datastructures"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type RepoContract struct {
	datastructures.Repo
	PushNumber int
}

//how to pass variables for initialization?
func (contract *RepoContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("initializing ledger")

	contract.Repo, _ = datastructures.CreateNewRepo("", 0, nil)
	contract.PushNumber = 0
	return shim.Success(nil)
}

func (contract *RepoContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "push" {
		return contract.addPush(APIstub, args)
	} else if function == "getbetween" {
		return contract.getpushes(APIstub, args)
	} else if function == "addbranch" {
		return contract.addBranch(APIstub, args)
	} else if function == "getbranches" {
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

	//it should be Marshaled on submission
	//pushAsBytes, _ := json.Marshal(args[1]) //changes from json object to bytes (string)
	var pushLog datastructures.PushLog
	json.Unmarshal([]byte(args[0]), &pushLog)

	if done, _ := contract.AddCommitLogs(pushLog.Logs, pushLog.BranchName); done {
		APIstub.PutState(string(contract.PushNumber), []byte(args[0]))
		contract.PushNumber = contract.PushNumber + 1

		return shim.Success(nil)
	}

	return shim.Error("Invalid push Log!")
}

func (contract *RepoContract) addBranch(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var branch datastructures.RepoBranch
	json.Unmarshal([]byte(args[0]), &branch)

	if done, _ := contract.AddBranch(branch); done {
		fmt.Println("New branch added:\t" + branch.Name)
		APIstub.PutState(branch.Name, []byte(branchAsBytes))
		return shim.Success(nil)
	}

	return shim.Error("Invalid Branch")
}

//
// func (contract *RepoContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	return shim.Success(nil)
// }

func (contract *RepoContract) getpushes(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("Querying the ledger..")
	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
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

		pushlogs = append(pushlogs, pushlog)
	}

	defer resultsIterator.Close()

	pushlogsjson, _ := json.Marshal(pushlogs)

	fmt.Println(pushlogs)
	return shim.Success(pushlogsjson)
}

func (contract *RepoContract) getbranches(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("Querying the ledger..")

	brancgesjson, _ := json.Marshal(contract.GetBranches())

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
