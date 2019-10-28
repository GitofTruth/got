package client

import (
	"encoding/json"

	"github.com/GitofTruth/GoT/datastructures"
)

//should get info about current user, connect to send transactions to a contract and send transactions accordingly
type Client struct {
	LastPush   int
	LocalRepo  datastructures.Repo
	RemoteRepo datastructures.Repo
}

type ArgsList struct {
	Args []string
}

func CreateNewArgsList(funcName string, args string) ArgsList {
	var argsList ArgsList
	argsList.Args = make([]string, 2)
	argsList.Args[0] = funcName
	argsList.Args[1] = args

	return argsList
}

func CreateNewClient(lastPush int) (Client, error) {
	var cli Client

	cli.LocalRepo, _ = datastructures.CreateNewRepo("mickey", 0, nil)
	//master, _ := datastructures.CreateNewRepoBranch("master", "mickey", 1, nil)
	//user.LocalRepo.AddBranch(master)

	//commit, _ := datastructures.CreateNewCommitLog("Testing the contract", "mickey", "mickey", 0, "COMMITHASH", nil, nil)

	cli.LastPush = 0

	return cli, nil
}

func (cli *Client) CreateBranchMessage(branch *datastructures.RepoBranch) (string, error) {
	master, _ := datastructures.CreateNewRepoBranch("master", "mickey", 1, nil)
	masterasbytes, _ := json.Marshal(master)
	argsStr, _ := json.Marshal(CreateNewArgsList("addBranch", string(masterasbytes)))

	return string(argsStr), nil
}

func (cli *Client) CreateCommitLogMessage(branch *datastructures.RepoBranch) (string, error) {
	commit, _ := datastructures.CreateNewCommitLog("Testing the contract", "mickey", "mickey", 0, "COMMITHASH", nil, nil)
	pushes := make([]datastructures.CommitLog, 1)
	pushes[0] = commit
	push, _ := datastructures.CreateNewPushLog("master", pushes)
	pushasbytes, _ := json.Marshal(push)
	argsStr, _ := json.Marshal(CreateNewArgsList("push", string(pushasbytes)))

	return string(argsStr), nil
}

//to invoke
//peer chaincode invoke -o "$ORDERER_ADDRESS" -C "$CC_CHANNEL_ID" -n "$CC_NAME"  -c "$CC_INVOKE_ARGS"
//CC_INVOKE_ARGS?
//'{"Args":["push","",""]}'

//CC_QUERY_ARGS?
//to query
//QUERY_RESULT=$(peer chaincode query -C "$CC_CHANNEL_ID" -n "$CC_NAME"  -c "$CC_QUERY_ARGS")
