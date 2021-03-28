package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	var input string
	if _, err := fmt.Scan(&input); err != nil {
		log.Fatal(err)
	}

	countWords := camelcase(input)
	fmt.Printf("word count: %d", countWords)
}

// camelcase returns the number of words
func camelcase(input string) int {
	countWords := 1
	for _, ch := range input {
		str := string(ch)
		if strings.ToUpper(str) == str {
			countWords++
		}
	}

	return countWords
}
