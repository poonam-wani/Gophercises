package cobra

import (
	"testing"
)

func TestSetCmd(t *testing.T) {

	a := []string{"key1", "123value"} // positive
	SetCmd.Run(SetCmd, a)

	b := []string{"", "%s"} //negative
	SetCmd.Run(SetCmd, b)

}
