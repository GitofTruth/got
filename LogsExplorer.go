package GoT

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type LogsExplorer struct {
	Repo *git.Repository
	Logs object.CommitIter
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
	le.Repo = repo
	return true, nil
}

func (le *LogsExplorer) LoadLogs() bool {
	if le.Repo != nil {
		ops := git.LogOptions{}
		ops.All = true
		ops.Order = git.LogOrderCommitterTime
		logs, _ := le.Repo.Log(&ops)
		le.Logs = logs
		fmt.Println(logs)
	}
	return true
}


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
		if err!=nil{
			break
		}
		fmt.Println(cmm)
	}

	return true
}
