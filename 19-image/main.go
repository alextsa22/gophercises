package main

import (
	"fmt"
	"log"
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {
	f, err := os.OpenFile("demo.svg", os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := []struct {
		Month string
		Point int
	}{
		{Month: "Jan", Point: 171},
		{Month: "Feb", Point: 180},
		{Month: "Mar", Point: 100},
		{Month: "Apr", Point: 87},
		{Month: "May", Point: 66},
		{Month: "Jun", Point: 40},
		{Month: "Jul", Point: 32},
		{Month: "Aug", Point: 55},
		{Month: "Sep", Point: 0},
		{Month: "Oct", Point: 0},
		{Month: "Nov", Point: 0},
		{Month: "Dec", Point: 0},
	}

	canvas := svg.New(f)
	width, height := len(data)*60+10, 300
	threshold := 120

	max := 0
	for _, dp := range data {
		if max < dp.Point {
			max = dp.Point
		}
	}

	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:rgb(255, 255, 255)")
	for i, dp := range data {
		percent := dp.Point * (height - 50) / max
		canvas.Rect(i*60+10, (height-50)-percent, 50, percent,
			"fill:rgb(77, 200, 232)")
		canvas.Text(i*60+35, height-20, dp.Month,
			"font-size:14pt;fill:rgb(150, 150, 150);text-anchor:middle")
	}
	threshPercent := threshold * (height - 50) / max
	canvas.Line(0, height-threshPercent, width, height-threshPercent,
		"stroke:rgb(255, 100, 100);opacity:0.8;stroke-width:2")
	canvas.Rect(0, 0, width, height-threshPercent,
		"fill:rgb(255, 100, 100);opacity:0.2")
	canvas.Line(0, height-50, width, height-50,
		"stroke:rgb(150, 150, 150);stroke-width:2")
	canvas.End()

	fmt.Println("image created.")
}
