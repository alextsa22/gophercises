package cobra

import (
	"fmt"
	"github.com/spf13/cobra"

	secret "github.com/alextsa22/gophercises/17-secret"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.NewVault(encodingKey, secretsPath())
		value, err := v.Get(args[0])
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Print(value)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
