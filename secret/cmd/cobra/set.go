package cobra

import (
	"fmt"

	"github.com/poonam-wani/gophercises/secret"
	"github.com/spf13/cobra"
)

// SetCmd function set the value for the  provided key
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			return
		}
		fmt.Println("value set successfully")
	},
}

func init() {
	RootCmd.AddCommand(SetCmd)
}
