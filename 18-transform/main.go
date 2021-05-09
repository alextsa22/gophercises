package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/18-transform/primitive"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

var html = `<html>
<body>
	<form action="/upload" method="POST" enctype="multipart/form-data">
		<input type="file" name="image">
		<button type="submit">upload image</button>
	</form>
</body>
</html>`

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		file, header, err := r.FormFile("image")
		if err != nil {
			log.Printf("FormFile(): %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)[1:]
		switch ext {
		case "jpg":
			fallthrough
		case "jpeg":
			ext = "jpeg"
		case "png":
		default:
			log.Printf("invalid image type: %s", ext)
			http.Error(w, "invalid image type", http.StatusBadRequest)
			return
		}

		out, err := primitive.Transform(file, ext, 10)
		if err != nil {
			log.Printf("Transform(): %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", fmt.Sprintf("image/%s", ext))
		io.Copy(w, out)
	})

	log.Println("server start...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
