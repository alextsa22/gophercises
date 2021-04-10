package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dir := "sample"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	var toRename []string
	for _, file := range files {
		if !file.IsDir() {
			if _, err := match(file.Name(), 0); err == nil {
				count++
				toRename = append(toRename, file.Name())
			}
		}
	}

	for _, origFilename := range toRename {
		origPath := filepath.Join(dir, origFilename)
		newFilename, err := match(origFilename, count)
		if err != nil {
			log.Fatal(err)
		}

		newPath := filepath.Join(dir, newFilename)
		fmt.Printf("mv %s => %s\n", origPath, newPath)

		if err = os.Rename(origPath, newPath); err != nil {
			log.Fatal(err)
		}
	}
}

// match returns the new filename, or an error if the filename
// didn't match our pattern.
func match(filename string, total int) (string, error) {
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didnt' match our pattern", filename)
	}

	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
}
