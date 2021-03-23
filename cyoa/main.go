package main

import (
	"flag"
	"fmt"
	"github.com/alextsa22/gophercises/cyoa/story"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	intro = "intro"
)

var (
	path = flag.String("path", "stories.json", "path to file containing stories")
)

func main() {
	flag.Parse()

	stories, err := story.NewStories(*path)
	if err != nil {
		log.Fatal(err)
	}

	temp, err := template.ParseFiles("template/index.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("template")),
	))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var path string
		if r.URL.Path == "/" {
			path = intro
		} else {
			path = strings.TrimLeft(r.URL.Path, "/")
		}
		story := stories[path]
		temp.Execute(w, story)
	})

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Starting the server on :8080")
	server.ListenAndServe()
}
