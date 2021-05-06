package cobra

import (
	"fmt"
	"github.com/spf13/cobra"

	secret "github.com/alextsa22/gophercises/17-secret"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.NewVault(encodingKey, secretsPath())
		if err := v.Set(args[0], args[1]); err != nil {
			fmt.Print(err)
			return
		}
		fmt.Print("secret set successfully")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
