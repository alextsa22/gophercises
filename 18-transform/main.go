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
	"strconv"
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
				<a href="/modify/{{.Name}}?mode={{.Mode}}">
					<img style="width: 22%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	imagesTpl = template.Must(template.New("").Parse(imagesHTML))

	imagesWithNumShapesHTML = `<html><body>
			{{range .}}
				<a href="/modify/{{.Name}}?mode={{.Mode}}&n={{.NumShapes}}">
					<img style="width: 22%;" src="/img/{{.Name}}">
				</a>
			{{end}}
			</body></html>`
	imagesWithNumShapesTpl = template.Must(template.New("").Parse(imagesWithNumShapesHTML))
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, uploadForm); err != nil {
			log.Printf("fmt.Fptint: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		file, header, err := r.FormFile("image")
		if err != nil {
			log.Printf("FormFile: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		ext := filepath.Ext(header.Filename)[1:]

		onDisk, err := tempFile("", ext)
		if err != nil {
			log.Printf("tempFile: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer onDisk.Close()
		if _, err = io.Copy(onDisk, file); err != nil {
			log.Printf("io.Copy: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/modify/"+filepath.Base(onDisk.Name()), http.StatusFound)
	})
	mux.HandleFunc("/modify/", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("./img/" + filepath.Base(r.URL.Path))
		if err != nil {
			log.Printf("os.Open: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		ext := filepath.Ext(f.Name())[1:]

		modeStr := r.FormValue("mode")
		if modeStr == "" {
			renderModeChoices(w, f, ext)
			return
		}
		mode, err := strconv.Atoi(modeStr)
		if err != nil {
			log.Printf("strconv.Atoi: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		numShapesStr := r.FormValue("numShapes")
		if numShapesStr == "" {
			renderNumShapeChoices(w, f, ext, primitive.Mode(mode))
			return
		}
		numShapes, err := strconv.Atoi(numShapesStr)
		if err != nil {
			log.Printf("strconv.Atoi: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = numShapes

		if _, err = io.Copy(w, f); err != nil {
			log.Printf("io.Copy: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.Handle(
		"/img/",
		http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))),
	)

	log.Println("server start...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func renderModeChoices(w http.ResponseWriter, rs io.ReadSeeker, ext string) {
	ns := 10
	opts := []generateOpts{
		{NumShapes: ns, Mode: primitive.ModeCircle},
		{NumShapes: ns, Mode: primitive.ModeRect},
		{NumShapes: ns, Mode: primitive.ModePolygon},
		{NumShapes: ns, Mode: primitive.ModeCombo},
	}
	images, err := generateImages(rs, ext, opts...)
	if err != nil {
		log.Printf("generateImages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type dataStruct struct {
		Name string
		Mode primitive.Mode
	}

	var data []dataStruct
	for i, img := range images {
		data = append(data, dataStruct{
			Name: filepath.Base(img),
			Mode: opts[i].Mode,
		})
	}

	for _, img := range data {
		img.Name = strings.ReplaceAll(img.Name, `\`, `/`)
	}

	if err = imagesTpl.Execute(w, data); err != nil {
		log.Printf("imagesTpl.Execute: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderNumShapeChoices(
	w http.ResponseWriter,
	rs io.ReadSeeker,
	ext string,
	mode primitive.Mode,
) {
	opts := []generateOpts{
		{NumShapes: 10, Mode: mode},
		{NumShapes: 20, Mode: mode},
		{NumShapes: 30, Mode: mode},
		{NumShapes: 40, Mode: mode},
	}

	images, err := generateImages(rs, ext, opts...)
	if err != nil {
		log.Printf("generateImages: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type dataStruct struct {
		Name      string
		NumShapes int
		Mode      primitive.Mode
	}

	var data []dataStruct
	for i, img := range images {
		data = append(data, dataStruct{
			Name:      filepath.Base(img),
			NumShapes: opts[i].NumShapes,
			Mode:      opts[i].Mode,
		})
	}

	for _, img := range data {
		img.Name = strings.ReplaceAll(img.Name, `\`, `/`)
	}

	if err = imagesWithNumShapesTpl.Execute(w, data); err != nil {
		log.Printf("imagesWithNumShapesTpl.Execute: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type generateOpts struct {
	NumShapes int
	Mode      primitive.Mode
}

func generateImages(rs io.ReadSeeker, ext string, opts ...generateOpts) ([]string, error) {
	var images []string
	for _, opt := range opts {
		if _, err := rs.Seek(0, 0); err != nil {
			return nil, fmt.Errorf("rs.Seek: %v", err)
		}
		img, err := generateImage(rs, ext, opt.NumShapes, opt.Mode)
		if err != nil {
			return nil, fmt.Errorf("generateImage: %v", err)
		}
		images = append(images, img)
	}
	return images, nil
}

func generateImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", fmt.Errorf("primitive.Transform: %v", err)
	}

	outFile, err := tempFile("", ext)
	if err != nil {
		return "", fmt.Errorf("tempFile; %v", err)
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
