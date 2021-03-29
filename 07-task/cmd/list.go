package cmd

import (
	"fmt"
	"github.com/alextsa22/gophercises/07-task/db"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all of your tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Printf("something went wrong: %s\n", err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("you have no tasks to complete! why not take a vacation?")
			return
		}

		fmt.Println("you have the following tasks:")
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}