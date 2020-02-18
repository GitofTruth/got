package main

import(
  "fmt"

  "github.com/GitofTruth/GoT/datastructures"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  "github.com/hyperledger/fabric/protos/peer"
  "github.com/hyperledger/fabric/protos/ledger/queryresult"


  "strconv"
)


func(contract *RepoContract) addNewRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
  // Input assumption
  // Example:
  // Required Tables:
  // Index(es) Used and primary key

  if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  
  fmt.Println("trying to process:\t", args[0])


  repo, err:= UnmarashalRepo(args[0])
  if err != nil{
    return shim.Error("Repo is invalid!")
  }

  k,v:= repo.GenerateRepoDBPair()
  stub.PutState(k,v)


  // Set entries in tables

    //data encoding


    // Composite keys encoding


    // Set the data


  return shim.Success(nil)
}












func(contract *RepoContract) addNewBranch(stub shim.ChaincodeStubInterface, args []string) peer.Response {
  // Input assumption
  // Example:
  // Required Tables
  // Index(ex) Used and primary key


  // Set entries in tables

    //data encoding


    // Composite keys encoding


    // Set the data

}
