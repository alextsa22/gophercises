package shorter

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

type pathToURL struct {
	Path string
	URL  string
}

func buildMap(pathsToURLs []pathToURL) (builtMap map[string]string) {
	builtMap = make(map[string]string)
	for _, ptu := range pathsToURLs {
		builtMap[ptu.Path] = ptu.URL
	}
	return
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(filename string, fallback http.Handler) (http.HandlerFunc, error) {
	ymlData, err := getFileBytes(filename)
	if err != nil {
		return nil, err
	}

	parsedYaml, err := parseYAML(ymlData)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yamlData []byte) (pathsToURLs []pathToURL, err error) {
	err = yaml.Unmarshal(yamlData, &pathsToURLs)
	return
}

func JSONHandler(filename string, fallback http.Handler) (http.HandlerFunc, error) {
	jsonData, err := getFileBytes(filename)
	if err != nil {
		return nil, err
	}

	parsedJSON, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsonData []byte) (pathsToURLs []pathToURL, err error) {
	err = json.Unmarshal(jsonData, &pathsToURLs)
	return
}

func getFileBytes(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
