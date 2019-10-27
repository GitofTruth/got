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
	contract.Repo, _ = datastructures.CreateNewRepo("", 0, nil)
	contract.PushNumber = 0
	return shim.Success(nil)
}

func (contract *RepoContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "push" {
		return contract.addPush(APIstub, args)
	} else if function == "getpushes" {
		return contract.getpushes(APIstub, args)
	} else if function == "addbranch" {

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
	//pushAsBytes, _ := json.Marshal(args[0])
	APIstub.PutState(string(contract.PushNumber), []byte(args[0]))
	contract.PushNumber = contract.PushNumber + 1

	var pushLog datastructures.PushLog
	json.Unmarshal([]byte(args[0]), &pushLog)

	done, _ := contract.AddCommitLogs(pushLog.Logs, pushLog.BranchName)

	if done {
		return shim.Success(nil)
	}

	return shim.Error("Invalid push Log")

}

//
// func (contract *RepoContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
// 	return shim.Success(nil)
// }

func (contract *RepoContract) getpushes(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

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
	return shim.Success(pushlogsjson)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(RepoContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
