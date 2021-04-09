package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	filename := "birthday_001.txt"
	newName, err := match(filename, 4)
	if err != nil {
		log.Fatalf("no match, received: %s", err)
	}

	fmt.Println(newName)
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
