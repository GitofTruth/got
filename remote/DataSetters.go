package main

import(
  "encoding/json"
  "fmt"

  "github.com/GitofTruth/GoT/datastructures"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  "github.com/hyperledger/fabric/protos/peer"


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
  
  // var repo datastructures.Repo

  repo, err:= datastructures.UnmarashalRepo(args[0])
  if err != nil{
    return shim.Error("Repo is invalid!")
  }

  
  k,v:= GenerateRepoDBPair(repo)
  keyBytes, _ := json.Marshal(k)
  stub.PutState(string(keyBytes),[]byte(v))


  // Set entries in tables

    //data encoding


    // Composite keys encoding


    // Set the data


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
