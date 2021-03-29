package cmd

import (
	"fmt"
	"github.com/alextsa22/gophercises/07-task/db"
	"github.com/spf13/cobra"
	"strconv"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("failed to parse the argument: %s\n", arg)
			} else {
				ids = append(ids, id)
			}
		}

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Printf("something went wrong: %s\n", err)
			return
		}

		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Printf("invalid task number: %d\n", id)
				continue
			}

			task := tasks[id-1]
			if err := db.DeleteTask(task.Key); err != nil {
				fmt.Printf("failed to mark \"%d\" as completed, error: %s\n", id, err)
			} else {
				fmt.Printf("marked \"%d\" as completed\n", id)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
