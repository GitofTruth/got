package client

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

var {
  cc="RepoContract"
}
func invokeCC(client *channel.Client, newValue string) {
	fmt.Println("Invoke cc with new value:", newValue)
	invokeArgs := [][]byte{[]byte("test-push"), []byte(newValue)}

	_, err := client.Execute(channel.Request{
		ChaincodeID: cc,
		Fcn:         "push",
		Args:        invokeArgs,
	})

	if err != nil {
		fmt.Printf("Failed to invoke: %+v\n", err)
	}
}
