package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	dry = flag.Bool("dry", true, "whether of not this should be a real or dry run")
)

func main() {
	flag.Parse()

	walkDir := "sample"
	toRename := make(map[string][]string)
	filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		curDir := filepath.Dir(path)
		if m, err := match(info.Name()); err == nil {
			key := filepath.Join(curDir, fmt.Sprintf("%s.%s", m.base, m.ext))
			toRename[key] = append(toRename[key], info.Name())
		}

		return nil
	})

	for key, files := range toRename {
		dir := filepath.Dir(key)
		n := len(files)
		sort.Strings(files)
		for i, filename := range files {
			result, _ := match(filename)
			newFilename := fmt.Sprintf("%s - %d of %d.%s", result.base, i+1, n, result.ext)

			oldPath := filepath.Join(dir, filename)
			newPath := filepath.Join(dir, newFilename)
			fmt.Printf("mv %s => %s\n", oldPath, newPath)

			if !*dry {
				if err := os.Rename(oldPath, newPath); err != nil {
					fmt.Printf("error renaming: %s => %s %s\n", oldPath, newPath, err)
				}
			}
		}
	}
}

type matchResult struct {
	base  string
	index int
	ext   string
}

// match returns the new filename, or an error if the filename
// didn't match our pattern.
func match(filename string) (*matchResult, error) {
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return nil, fmt.Errorf("%s didnt' match our pattern", filename)
	}

	return &matchResult{strings.Title(name), number, ext}, nil
}
