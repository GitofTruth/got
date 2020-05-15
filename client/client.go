package client

// should get info about current user, connect to send transactions to a contract and send transactions accordingly
// Go Client is currently obsolete, only nodeJS is currently supported.
type Client struct {
	// LastPush   int
	// LocalRepo  datastructures.Repo
	// RemoteRepo datastructures.Repo
	// Explorer   LogsExplorer
}

/*
func CreateNewClient(lastPush int) (Client, error) {
	var cli Client

	exp, err := CreateNewLogsExplorer("")
	if err != nil {
		panic(err)
	}
	cli.Explorer = exp
	// cli.LocalRepo, _ = cli.Explorer.GetInternalRepo()

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

	push, _ := datastructures.CreateNewPushLog("repo", branch.Name, branchcommits)
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

func (cli *Client) CreateCommitLogMessage(branch *datastructures.RepoBranch) (string, error) {
	// commit, _ := datastructures.CreateNewCommitLog("Testing the contract", "mickey", "mickey", 0, "COMMITHASH", nil, nil, nil, nil)
	pushes := make([]datastructures.CommitLog, 1)
	// pushes[0] = commit
	push, _ := datastructures.CreateNewPushLog("repo", "master", pushes)
	pushasbytes, _ := json.Marshal(push)
	argsStr, _ := json.Marshal(common.CreateNewArgsList("push", string(pushasbytes)))

	return string(argsStr), nil
}

func (cli *Client) CreateAddNewRepoMessage() string {

	// x, _ := datastructures.CreateNewRepo("repo", "GoT", "hassan", 0, nil, nil, nil)
	// branch, _ := datastructures.CreateNewRepoBranch("master", "masterCreator", 1, nil)
	// x.AddBranch(branch)
	// commit, _ := datastructures.CreateNewCommitLog("message", "mickey", "mickeyAsCommiter", 3, "*************", nil, nil, nil, nil)
	// x.AddCommitLog(commit, "master")
	// str, _ := json.Marshal(x)

	argsStr, _ := json.Marshal(common.CreateNewArgsList("addNewRepo", string("str")))

	return string(argsStr)
}
*/
