package main

import(
  "fmt"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  "github.com/hyperledger/fabric/protos/peer"
  "github.com/hyperledger/fabric/protos/ledger/queryresult"


  "strconv"
)


func(contract *RepoContract) getRepo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
  fmt.Println("Querying the ledger..")
  
  if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

  repoHash := GetRepoKey(args[0], args[1])
  repoData, err := stub.GetState(repoHash)
  if err != nil{
    return shim.Error("Repo does not exist")
  }

  fmt.Println("Found this repo:", repoData)

  return shim.Success(repoData)
}