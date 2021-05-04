package main

import (
	"fmt"
	"log"

	secret "github.com/alextsa22/gophercises/17-secret"
)

func main() {
	v := secret.NewVault("my-fake-key", ".secrets")
	
	err := v.Set("demo_key1", "123 some crazy value")
	if err != nil {
		log.Fatal(err)
	}

	err = v.Set("demo_key2", "456 some crazy value")
	if err != nil {
		log.Fatal(err)
	}

	err = v.Set("demo_key3", "789 some crazy value")
	if err != nil {
		log.Fatal(err)
	}

	plain, err := v.Get("demo_key1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("plain:", plain)

	plain, err = v.Get("demo_key2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("plain:", plain)

	plain, err = v.Get("demo_key3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("plain:", plain)
}
