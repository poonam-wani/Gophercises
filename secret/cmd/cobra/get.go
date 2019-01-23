package cobra

import (
	"fmt"

	"github.com/poonam-wani/gophercises/secret"
	"github.com/spf13/cobra"
)

//GetCmd cobra command used to get the value for provided key
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {

		v := secret.File(encodingKey, secretsPath())

		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("no value set")
			return
		}
		fmt.Printf("%s=%s\n", key, value)
	},
}

func init() {
	RootCmd.AddCommand(GetCmd)
}
