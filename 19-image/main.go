package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	data := []int{100, 33, 73, 64}

	width, height := len(data)*60+10, 100
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	bg := image.NewUniform(color.RGBA{R: 240, G: 240, B: 240, A: 255})

	draw.Draw(img, rect, bg, image.Point{}, draw.Src)

	mask := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var alpha uint8
			switch {
			case y < 30:
				alpha = 255
			case y < 50:
				alpha = 100
			}
			mask.Set(x, y, color.RGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: alpha,
			})
		}
	}

	for i, dp := range data {
		x0, y0 := i*60+10, 100-dp
		x1, y1 := (i+1)*60-1, 100
		bar := image.Rect(x0, y0, x1, y1)
		grey := image.NewUniform(color.RGBA{180, 180, 180, 255})
		draw.Draw(img, bar, grey, image.Point{}, draw.Src)

		red := image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		draw.DrawMask(img, bar, red, image.Point{}, mask, image.Point{X: x0, Y: y0}, draw.Over)
	}

	f, err := os.Create("demo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err = png.Encode(f, img); err != nil {
		log.Fatal(err)
	}

	fmt.Println("image created.")
}
