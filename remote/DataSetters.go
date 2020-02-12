package main

import(
  "fmt"

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


  // Set entries in tables

    //data encoding


    // Composite keys encoding


    // Set the data
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
