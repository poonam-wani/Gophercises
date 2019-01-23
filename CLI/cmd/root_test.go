package cmd

import (
	"fmt"
	"testing"
)

// TestRootCmd tests for root command
func TestRootCmd(t *testing.T) {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println("Error in root")
	}
}
