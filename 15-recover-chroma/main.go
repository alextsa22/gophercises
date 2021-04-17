package main

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", sourceCodeHandler)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(mux)))
}

func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := bytes.NewBuffer(nil)
	if _, err = io.Copy(b, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	quick.Highlight(w, b.String(), "go", "html", "monokai")
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, makeLinks(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func makeLinks(stack string) string {
	var filename, tmp string
	var bias int
	lines := strings.Split(stack, "\n")
	for li, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}

		tmp, bias = line, 0
		if strings.HasPrefix(line, "\tC:") {
			tmp, bias = line[3:], 3
		}

		for i, ch := range tmp {
			if ch == ':' {
				filename = line[:i+bias]
				break
			}
		}

		lines[li] = `<a href="/debug/?path=` + filename + `"/>` +
			filename + `</a>` + line[len(filename)+1:]
	}

	return strings.Join(lines, "\n")
}
