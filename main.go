package GoT

import (
	"fmt"

	"github.com/GitofTruth/GoT/client"
)

func main() {

	cli, _ := client.CreateNewClient(0)
	str, _ := cli.CreateBranchMessage(nil)
	fmt.Println(str)

}
