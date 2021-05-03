package main

import (
	"fmt"
	secret "github.com/alextsa22/gophercises/17-secret"
	"log"
)

func main() {
	demoKey := "demo_key"

	v := secret.NewVault("fake key")
	if err := v.Set(demoKey, "demo_value"); err != nil {
		log.Fatal(err)
	}

	value, err := v.Get(demoKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("value:", value)
}
