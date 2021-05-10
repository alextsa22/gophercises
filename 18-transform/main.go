package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/18-transform/primitive"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		out, err := primitive.Transform(file, ext, 10, primitive.WithMode(primitive.ModeCircle))
		if err != nil {
			log.Printf("Transform(): %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		outFile, err := tempFile("__out_", ext)
		if err != nil {
			log.Printf("tempFile(): %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer outFile.Close()

		if _, err = io.Copy(outFile, out); err != nil {
			log.Printf("io.Copy(): %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		redirectUrl := fmt.Sprintf("/%s", outFile.Name())
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	})
	mux.Handle(
		"/img/",
		http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))),
	)

	log.Println("server start...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func tempFile(prefix, ext string) (*os.File, error) {
	tmp, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	return os.Create(fmt.Sprintf("%s.%s", tmp.Name(), ext))
}
