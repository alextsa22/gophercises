package main

import (
	"fmt"
	"log"
)

func main() {
	var input string
	if _, err := fmt.Scan(&input); err != nil {
		log.Fatal(err)
	}

	var k int
	if _, err := fmt.Scan(&k); err != nil {
		log.Fatal(err)
	}

	encrypted := caesarCipher(input, k)
	fmt.Printf("encrypted: %s", encrypted)
}

// caesarCipher encrypts the string using the Caesar cipher
func caesarCipher(s string, k int) string {
	var rs []rune
	for _, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			rs = append(rs, rotate(ch, 'A', k))
		}
		if ch >= 'a' && ch <= 'z' {
			rs = append(rs, rotate(ch, 'a', k))
		}
	}

	return string(rs)
}

func rotate(r rune, base, k int) rune {
	tmp := int(r) - base
	tmp = (tmp + k) % 26
	return rune(tmp + base)
}
