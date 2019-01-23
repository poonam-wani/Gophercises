package cmd

import (
	"fmt"

	"github.com/poonam-wani/gophercises/CLI/db"
	"github.com/spf13/cobra"
)

// FunctionList Created the function as variable for test coverage
var FunctionList = db.GetAllLists

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list is used to list down all the tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := FunctionList()
		if err != nil {
			fmt.Println("Something went wrong", err)
			//os.Exit(1)
		}
		// Checking the condition when list is empty
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete")
			return
		}

		fmt.Println("You have following list of tasks:")
		for i, task := range tasks {
			fmt.Printf("%d.%s\n", i+1, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
