package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type RepoContract struct {
}

func (contract *RepoContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("initializing ledger")
	return shim.Success(nil)
}

func (contract *RepoContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	fmt.Println("****************************************\nStarting invokation .. \nfunctionName:\t"+function+"\nargs:\t\n", args)
	defer fmt.Println("Invokation end\n\n")

	if function == "addNewRepo" {
		return contract.addNewRepo(stub, args)
	} else if function == "queryRepo" {
		return contract.queryRepo(stub, args)
	} else if function == "clone" {
		return contract.clone(stub, args)
	} else if function == "addNewBranch" {
		return contract.addNewBranch(stub, args)
	} else if function == "queryBranches" {
		return contract.queryBranches(stub, args)
	} else if function == "queryBranch" {
		return contract.queryBranch(stub, args)
	} else if function == "addCommits" {
		return contract.addCommits(stub, args)
	} else if function == "queryBranchCommits" {
		return contract.queryBranchCommits(stub, args)
	} else if function == "queryLastBranchCommit" {
		return contract.queryLastBranchCommit(stub, args)
	} else if function == "addUserUpdate" {
		return contract.addUserUpdate(stub, args)
	} else if function == "queryUser" {
		return contract.queryUser(stub, args)
	} else if function == "queryUsers" {
		return contract.queryUsers(stub, args)
	} else if function == "updateRepoUserAccess" {
		return contract.updateRepoUserAccess(stub, args)
	} else if function == "queryRepoUserAccess" {
		return contract.queryRepoUserAccess(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(RepoContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
