package datastructures

import (
	"fmt"
	"testing"
)

func TestLogsExplorer(t *testing.T) {
	fmt.Println("Starting Test \t\t LogsExplorer")

	le, err := CreateNewLogsExplorer("")
	if err != nil {
		panic(err)
	}

	le.LoadLogs()
	le.PrintAllLogs()

}
