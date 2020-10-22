// Retired sub-package used only in testing now.
package common

import ()

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
