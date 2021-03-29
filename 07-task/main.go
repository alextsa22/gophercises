package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/07-task/cmd"
	"github.com/alextsa22/gophercises/07-task/db"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

func main() {
	homeDir, _ := homedir.Dir()
	dbPath := filepath.Join(homeDir, "tasks.db")
	if err := db.Init(dbPath); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
