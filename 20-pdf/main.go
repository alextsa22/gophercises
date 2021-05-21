package main

import (
	"fmt"
	"log"

	"github.com/jung-kurt/gofpdf"
)

const (
	bannerHeight = 94.0
	xIndent      = 45.0
)

func main() {
	pdf := gofpdf.New(
		gofpdf.OrientationPortrait,
		gofpdf.UnitPoint,
		gofpdf.PageSizeLetter,
		"",
	)
	width, height := pdf.GetPageSize()
	fmt.Printf("width=%v, height=%v\n", width, height)
	pdf.AddPage()

	// banner
	pdf.SetFillColor(103, 60, 79)
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{width, 0},
		{width, bannerHeight},
		{0, bannerHeight * 0.9},
	}, "F")
	pdf.Polygon([]gofpdf.PointType{
		{0, height},
		{0, height - (bannerHeight * 0.2)},
		{width, height - (bannerHeight * 0.1)},
		{width, height},
	}, "F")

	// banner - invoice
	pdf.SetFont("arial", "B", 40)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight := pdf.GetFontSize()
	pdf.Text(xIndent, bannerHeight-(bannerHeight/2)+lineHeight/3.1, "INVOICE")

	// banner - phone, email, domain
	pdf.SetFont("arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight = pdf.GetFontSize()
	pdf.MoveTo(width-xIndent-2*124, (bannerHeight-(lineHeight*1.5*3.0))/2)
	pdf.MultiCell(
		124.0,
		lineHeight*1.5,
		"(123) 456-7890\njon@calhoun.io\nGophercises.com",
		gofpdf.BorderNone,
		gofpdf.AlignRight,
		false,
	)

	// banner - address
	pdf.SetFont("arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight = pdf.GetFontSize()
	pdf.MoveTo(width-xIndent-124, (bannerHeight-(lineHeight*1.5*3.0))/2)
	pdf.MultiCell(
		124.0,
		lineHeight*1.5,
		"123 Fake St\nSome Town, PA\n12345",
		gofpdf.BorderNone,
		gofpdf.AlignRight,
		false,
	)


	// Grid
	// drawGrid(pdf)

	if err := pdf.OutputFileAndClose("p2.pdf"); err != nil {
		log.Fatal(err)
	}
}

func drawGrid(pdf *gofpdf.Fpdf) {
	width, height := pdf.GetPageSize()
	pdf.SetFont("courier", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetDrawColor(200, 200, 200)

	for x := 0.0; x < width; x = x + (width / 20.0) {
		pdf.SetTextColor(200, 200, 200)
		pdf.Line(x, 0, x, height)
		_, lineHt := pdf.GetFontSize()
		pdf.Text(x, lineHt, fmt.Sprintf("%d", int(x)))
	}

	for y := 0.0; y < height; y = y + (width / 20.0) {
		if y < bannerHeight*0.9 {
			pdf.SetTextColor(200, 200, 200)
		} else {
			pdf.SetTextColor(80, 80, 80)
		}
		pdf.Line(0, y, width, y)
		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
	}
}
