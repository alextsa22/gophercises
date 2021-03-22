package main

import (
	"flag"
	"fmt"
	"github.com/alextsa22/gophercises/urlshort/shorter"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var (
	path = flag.String("path", "redirect.yml", "path to file containing shortened paths to URL's")
)

func main() {
	flag.Parse()

	mux := defaultMux()

	ext := filepath.Ext(*path)

	var (
		handler http.Handler
		err     error
	)

	if ext == ".yml" || ext == ".yaml" {
		handler, err = shorter.YAMLHandler(*path, mux)
		if err != nil {
			log.Fatal(err)
		}
	} else if ext == ".json" {
		handler, err = shorter.JSONHandler(*path, mux)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("unsupported %s file type", ext)
	}

	server := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Starting the server on :8080")
	server.ListenAndServe()
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	return mux
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "default page")
}
