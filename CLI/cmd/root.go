package cmd

import "github.com/spf13/cobra"

// RootCmd root command
var RootCmd = &cobra.Command{
	Use:   "add",
	Short: "Add the todo list tasks",
}
