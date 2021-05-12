package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	data := []int{100, 33, 73, 64}

	width, height := len(data)*60+10, 100
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255}) // white
		}
	}

	for i, dp := range data {
		for x := i*60 + 10; x < (i+1)*60; x++ {
			for y := height; y >= (100 - dp); y-- {
				img.Set(x, y, color.RGBA{180, 180, 250, 255})
			}
		}
	}

	f, err := os.Create("demo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err = png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
