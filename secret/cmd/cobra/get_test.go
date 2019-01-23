package cobra

import "testing"

func TestGetCmd(t *testing.T) {

	a := []string{"key1"} //Positive
	GetCmd.Run(GetCmd, a)

	a = []string{"aafreerg"} //Negative
	GetCmd.Run(GetCmd, a)
}
