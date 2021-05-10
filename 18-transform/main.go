package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/18-transform/primitive"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	uploadForm = `<html>
			<body>
				<form action="/upload" method="POST" enctype="multipart/form-data">
					<input type="file" name="image">
					<button type="submit">upload image</button>
				</form>
			</body>
			</html>`

	imagesHTML = `<html><body>
			{{range .}}
				<img src="/{{.}}">
			{{end}}
			</body></html>`
	imagesTpl = template.Must(template.New("").Parse(imagesHTML))
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, uploadForm)
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

		circleImg, _ := generateImage(file, ext, 10, primitive.ModeCircle)
		file.Seek(0, 0)
		rectImg, _ := generateImage(file, ext, 10, primitive.ModeRect)
		file.Seek(0, 0)
		polygonImg, _ := generateImage(file, ext, 10, primitive.ModePolygon)

		images := []string{circleImg, rectImg, polygonImg}
		for i, img := range images {
			images[i] = strings.ReplaceAll(img, `\`, `/`)
		}
		imagesTpl.Execute(w, images)
	})
	mux.Handle(
		"/img/",
		http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))),
	)

	log.Println("server start...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func generateImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", fmt.Errorf("primitive.Transform(): %v", err)
	}

	outFile, err := tempFile("", ext)
	if err != nil {
		return "", fmt.Errorf("tempFile(); %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, out)
	return outFile.Name(), err
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
