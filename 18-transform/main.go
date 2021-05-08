package main

import (
	"fmt"
	"github.com/alextsa22/gophercises/18-transform/primitive"
	"io"
	"log"
	"os"
)

var (
	projectPath = "./18-transform/"
	inPath      = fmt.Sprintf("%ssamurai.jpg", projectPath)
	outPath     = fmt.Sprintf("%sout.jpg", projectPath)
)

func main() {
	in, err := os.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}
	transformOut, err := primitive.Transform(in, 10)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.Remove(outPath); err != nil {
		log.Fatal()
	}
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if _, err = io.Copy(out, transformOut); err != nil {
		log.Fatal(err)
	}
}
