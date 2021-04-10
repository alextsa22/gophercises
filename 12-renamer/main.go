package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	dry           = flag.Bool("dry", true, "whether of not this should be a real or dry run")
	reg           = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
	replaceString = "$2 - $1 - $3 of $4.$5"
)

func main() {
	flag.Parse()

	walkDir := "sample"
	var toRename []string
	filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}

		return nil
	})

	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		filename := filepath.Base(oldPath)
		newFilename, _ := match(filename)
		newPath := filepath.Join(dir, newFilename)

		fmt.Printf("mv %s => %s\n", oldPath, newPath)
		if !*dry {
			if err := os.Rename(oldPath, newPath); err != nil {
				fmt.Printf("error renaming: %s => %s %s\n", oldPath, newPath, err)
			}
		}
	}
}

// match returns the new filename, or an error if the filename
// didn't match our pattern.
func match(filename string) (string, error) {
	if !reg.MatchString(filename) {
		return "", fmt.Errorf("%s didnt' match our pattern", filename)
	}

	return reg.ReplaceAllString(filename, replaceString), nil
}
