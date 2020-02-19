package client

import (
	"encoding/json"

	"github.com/GitofTruth/GoT/common"
	"github.com/GitofTruth/GoT/datastructures"
)

//should get info about current user, connect to send transactions to a contract and send transactions accordingly
type Client struct {
	LastPush   int
	LocalRepo  datastructures.Repo
	RemoteRepo datastructures.Repo
	Explorer   datastructures.LogsExplorer
}

func CreateNewClient(lastPush int) (Client, error) {
	var cli Client

	exp, err := datastructures.CreateNewLogsExplorer("")
	if err != nil {
		panic(err)
	}
	cli.Explorer = exp
	cli.LocalRepo, _ = cli.Explorer.GetInternalRepo()

	cli.LastPush = 0
	return cli, nil
}

func (cli *Client) CreatePushMessage(m int) (string, error) {
	branch := cli.LocalRepo.Branches["master"]
	branchcommits := make([]datastructures.CommitLog, 0)

	n := 0
	for _, val := range branch.Logs {
		if n == m {
			break
		}
		branchcommits = append(branchcommits, val)
		n = n + 1
	}

	push, _ := datastructures.CreateNewPushLog(branch.Name, branchcommits)
	pushasbytes, _ := json.Marshal(push)
	argsStr, _ := json.Marshal(common.CreateNewArgsList("push", string(pushasbytes)))

	return string(argsStr), nil
}

func (cli *Client) CreateBranchMessage(branch *datastructures.RepoBranch) (string, error) {
	master, _ := datastructures.CreateNewRepoBranch("master", "mickey", 1, nil)
	masterasbytes, _ := json.Marshal(master)
	argsStr, _ := json.Marshal(common.CreateNewArgsList("addBranch", string(masterasbytes)))

	return string(argsStr), nil
}

func (cli * Client) CreateCommitLogMessage(branch *datastructures.RepoBranch) (string, error) {
	commit, _ := datastructures.CreateNewCommitLog("Testing the contract", "mickey", "mickey", 0, "COMMITHASH", nil, nil)
	pushes := make([]datastructures.CommitLog, 1)
	pushes[0] = commit
	push, _ := datastructures.CreateNewPushLog("master", pushes)
	pushasbytes, _ := json.Marshal(push)
	argsStr, _ := json.Marshal(common.CreateNewArgsList("push", string(pushasbytes)))

	return string(argsStr), nil
}

func (cli * Client) CreateAddNewRepoMessage() string {

	x, _ := datastructures.CreateNewRepo("GoT", "mickey", 0, nil)
	str, _ := json.Marshal(x)

	argsStr, _ := json.Marshal(common.CreateNewArgsList("addNewRepo", string(str)))


	return string(argsStr)
}


//to invoke
//peer chaincode invoke -o "$ORDERER_ADDRESS" -C "$CC_CHANNEL_ID" -n "$CC_NAME"  -c "$CC_INVOKE_ARGS"
//CC_INVOKE_ARGS?
//'{"Args":["push","",""]}'

//CC_QUERY_ARGS?
//to query
//QUERY_RESULT=$(peer chaincode query -C "$CC_CHANNEL_ID" -n "$CC_NAME"  -c "$CC_QUERY_ARGS")
