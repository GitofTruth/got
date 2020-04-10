package datastructures

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRepo(t *testing.T) {
	fmt.Println("Starting Test \t\t Repo")

	x, _ := CreateNewRepo("GoT", "mickey", "DIRECOTRYCID" ,0, nil)
	str, _ := json.Marshal(x)
	fmt.Println(string(str))
}
