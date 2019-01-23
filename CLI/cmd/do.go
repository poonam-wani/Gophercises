package cmd

import (
	"fmt"
	"strconv"

	"github.com/poonam-wani/gophercises/CLI/db"
	"github.com/spf13/cobra"
)

// FunctionDo Created the function as a variable for test coverage
var FunctionDo = db.DeleteTasks

// DoCmd represents the do command
var DoCmd = &cobra.Command{
	Use:   "do",
	Short: "do is used to mark the task as complete",
	Long:  `Complete the task by do command`,
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg) // Atoi converts the string (which we have passed with do command through cmd) to integer
			if err != nil {
				fmt.Println("Failed to parse the argument", arg)
			} else {
				ids = append(ids, id)
			}
		}
		tasks, err := FunctionList() // It will lists all tasks present in the DB
		if err != nil {
			fmt.Println("Something went wrong", err)
			//os.Exit(1)
		}

		for _, id := range ids {
			if id <= 0 || id > len(ids) { // this checks the id provided with do command is valid or not
				fmt.Println("Invalid task number", id)
				continue
			}
			task := tasks[id-1]
			err = FunctionDo(task.Key)
			if err != nil {
				fmt.Printf("Failed to mark \"%d\" task as complete, error= %s\n", id, err)
			} else {
				fmt.Printf("Marked \"%d\" task as complete\n", id)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(DoCmd)
}
