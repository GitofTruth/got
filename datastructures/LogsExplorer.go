package datastructures

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type LogsExplorer struct {
	GitRepo *git.Repository
	Logs    object.CommitIter
}

func CreateNewLogsExplorer(path string) (LogsExplorer, error) {
	if path == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Using default directory:\t" + dir)
		path = dir
	}

	le := LogsExplorer{}
	_, err := le.OpenRepo(path)
	if err != nil {
		panic(err)
	}

	return le, nil
}

func (le *LogsExplorer) OpenRepo(path string) (bool, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	le.GitRepo = repo
	return true, nil
}

func (le *LogsExplorer) LoadLogs() bool {
	if le.GitRepo != nil {
		ops := git.LogOptions{}
		ops.All = true
		ops.Order = git.LogOrderCommitterTime
		logs, _ := le.GitRepo.Log(&ops)
		le.Logs = logs
	}
	return true
}

func (le *LogsExplorer) GetInternalRepo() (Repo, error) {
	//author first commit
	//branch >> commits

	le.LoadLogs()

	if le.Logs != nil {
		fmt.Println(le.Logs)
	}

	logs := make(map[string]CommitLog)
	auth := ""
	time := 0
	for {
		cmm, err := le.Logs.Next()
		if err != nil {
			break
		}
		parentHashes := make([]string, 0)
		for _, hashplumb := range cmm.ParentHashes {
			parentHashes = append(parentHashes, string([]byte(hashplumb[:])))
		}
		logs[string([]byte(cmm.Hash[:]))], _ = CreateNewCommitLog(cmm.Message, cmm.Author.Name, cmm.Committer.Name, cmm.Committer.When.Second(), string([]byte(cmm.Hash[:])), parentHashes, nil)
		auth = cmm.Author.Name
		time = cmm.Committer.When.Second()
	}

	//logs map[string]CommitLog
	branch, _ := CreateNewRepoBranch("master", auth, time, logs)
	branches := make(map[string]RepoBranch)
	branches[branch.Name] = branch

	return CreateNewRepo(auth,"", time, branches)
}

//get repo general info

//get current branch info

//get current branch logs

func (le *LogsExplorer) PrintAllLogs() bool {
	if le.Logs != nil {
		fmt.Println(le.Logs)
	}

	// le.Logs.ForEach(func(arg1 *object.Commit) error {
	// 	fmt.Println(arg1)
	// 	return nil
	// })

	for {
		fmt.Println("starting iterating")
		cmm, err := le.Logs.Next()
		if err != nil {
			break
		}
		fmt.Println(cmm)
	}

	return true
}