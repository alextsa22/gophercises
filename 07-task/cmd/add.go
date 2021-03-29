package cmd

import (
	"fmt"
	"github.com/alextsa22/gophercises/07-task/db"
	"github.com/spf13/cobra"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		if _, err := db.CreateTask(task); err != nil {
			fmt.Printf("something went wrong: %s\n", err)
			return
		}

		fmt.Printf("added \"%s\" to your task list\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
