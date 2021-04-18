package main

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
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
	lineParam := r.FormValue("line")
	line, err := strconv.Atoi(lineParam)
	if err != nil {
		line = -1
	}

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

	var lines [][2]int
	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}

	lexer := lexers.Get("go")
	iterator, err := lexer.Tokenise(nil, b.String())
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(
		html.TabWidth(2),
		html.WithLineNumbers(true),
		html.HighlightLines(lines),
	)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<style>pre { font-size: 1.2em; }</style>")
	formatter.Format(w, style, iterator)
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

		var lineStr strings.Builder
		for i := len(filename) + 1; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}
			lineStr.WriteByte(line[i])
		}

		lines[li] = `<a href="/debug/?path=` + filename +
			`&line=` + lineStr.String() + `"/>` +
			filename + ":" + lineStr.String() + `</a>` +
			line[len(filename)+1+len(lineStr.String()):]
	}

	return strings.Join(lines, "\n")
}
