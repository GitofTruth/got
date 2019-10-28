package common

import (
	"fmt"

	"github.com/GitofTruth/GoT/datastructures"
)

func GetPushesDifference(local datastructures.Repo, remote datastructures.Repo) {
	fmt.Println("hello")
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
