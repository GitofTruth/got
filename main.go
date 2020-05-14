package got

import (
	"fmt"
	"os"

	"github.com/GitofTruth/got/client"
)

func main() {

	cli, _ := client.CreateNewClient(0)
	// str, _ := cli.CreatePushMessage(1)
	newRepoMsg := cli.CreateAddNewRepoMessage()
	os.Setenv("CC_INVOKE_ARGS", newRepoMsg)
	// os.Setenv("CC_QUERY_ARGS", "{\"Args\":[\"getBetween\", \"0\", \"1\"]}")
	os.Setenv("CC_QUERY_ARGS", "{\"Args\":[\"getRepo\", \"mickey\", \"GoT\"]}")
	fmt.Println("CC_INVOKE_ARGS:", os.Getenv("CC_INVOKE_ARGS"))
	fmt.Println("CC_QUERY_ARGS:", os.Getenv("CC_QUERY_ARGS"))

}
