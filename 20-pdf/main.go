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

	// summary - billed to, invoice #, date of issue
	_, sy := summaryBlock(
		pdf,
		xIndent,
		bannerHeight+lineHeight*2.0,
		"Billed To",
		"Client Name", "123 Client Address", "City, State, Country", "Postal Code",
	)
	summaryBlock(
		pdf,
		xIndent*2.0+lineHeight*12.5,
		bannerHeight+lineHeight*2.0,
		"Invoice Number",
		"0000000123",
	)
	summaryBlock(pdf,
		xIndent*2.0+lineHeight*12.5,
		bannerHeight+lineHeight*6.25,
		"Date of Issue",
		"05/29/2018",
	)

	// summary - invoice total
	x, y := width-xIndent-124.0, bannerHeight+lineHeight*2.25
	pdf.MoveTo(x, y)
	pdf.SetFont("times", "", 14)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.CellFormat(
		124.0,
		lineHeight,
		"Invoice Total",
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	x, y = x+2.0, y+lineHeight*1.5
	pdf.MoveTo(x, y)
	pdf.SetFont("times", "", 48)
	_, lineHeight = pdf.GetFontSize()
	alpha := 58
	pdf.SetTextColor(72+alpha, 42+alpha, 55+alpha)
	totalUSD := "$1234.56"
	pdf.CellFormat(
		124.0,
		lineHeight,
		totalUSD,
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	x, y = x-2.0, y+lineHeight*1.25

	if sy > y {
		y = sy
	}
	x, y = xIndent-20.0, y+30.0
	pdf.Rect(x, y, width-(xIndent*2.0)+40.0, 3.0, "F")

	// Grid
	// drawGrid(pdf)

	if err := pdf.OutputFileAndClose("p3.pdf"); err != nil {
		log.Fatal(err)
	}
}

func summaryBlock(pdf *gofpdf.Fpdf, x, y float64, title string, data ...string) (float64, float64) {
	pdf.SetFont("times", "", 14)
	pdf.SetTextColor(180, 180, 180)
	_, lineHt := pdf.GetFontSize()
	y = y + lineHt
	pdf.Text(x, y, title)
	y = y + lineHt*.25
	pdf.SetTextColor(50, 50, 50)
	for _, str := range data {
		y = y + lineHt*1.25
		pdf.Text(x, y, str)
	}
	return x, y
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
