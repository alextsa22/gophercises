package main

import (
	"flag"
	"fmt"
	"github.com/alextsa22/gophercises/04-link/link"
	"log"
	"os"
)

var (
	path = flag.String("path", "examples/ex1.html", "path to html file")
)

func main() {
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	links, err := link.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	for _, l := range links {
		fmt.Println("href: ", l.Href)
		fmt.Println("text: ", l.Text)
	}
}
