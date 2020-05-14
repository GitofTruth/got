package got

import (
	"fmt"
	"os"
	"testing"

	"github.com/GitofTruth/got/client"
)

func TestBasicgroup(t *testing.T) {
	fmt.Println("Starting Test \t\t Main")

	cli, _ := client.CreateNewClient(0)
	str, _ := cli.CreatePushMessage(2)
	os.Setenv("CC_INVOKE_ARGS", str)
	os.Setenv("CC_QUERY_ARGS", "{\"Args\":[\"getBetween\", \"0\", \"2\"]")
	fmt.Println("CC_INVOKE_ARGS:", os.Getenv("CC_INVOKE_ARGS"))
	fmt.Println("CC_QUERY_ARGS:", os.Getenv("CC_QUERY_ARGS"))

	/*

		fmt.Println("Reading connection profile..")
		c := config.FromFile("/home/hkandil/fabric-samples/first-network/connection-org1.yaml")
		sdk, err := fabsdk.New(c)
		if err != nil {
			fmt.Printf("Failed to create new SDK: %s\n", err)
			os.Exit(1)
		}
		defer sdk.Close()

		setupLogLevel()
		enrollUser(sdk)

		clientChannelContext := sdk.ChannelContext(channelName, fabsdk.WithUser(user))
		ledgerClient, err := ledger.New(clientChannelContext)
		if err != nil {
			fmt.Printf("Failed to create channel [%s] client: %#v", channelName, err)
			os.Exit(1)
		}

		fmt.Printf("\n===== Channel: %s ===== \n", channelName)
		queryChannelInfo(ledgerClient)
		queryChannelConfig(ledgerClient)

		fmt.Println("\n====== Chaincode =========")

		client, err := channel.New(clientChannelContext)
		if err != nil {
			fmt.Printf("Failed to create channel [%s]:", channelName)
		}

		invokeCC(client, "100")
		old := queryCC(client, []byte("john"))

		oldInt, _ := strconv.Atoi(old)
		invokeCC(client, strconv.Itoa(oldInt+1))

		queryCC(client, []byte("john"))

		fmt.Println("===============")
		fmt.Println("Done.")

	*/

}
