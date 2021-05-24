package main

import (
	"fmt"
	"log"

	"github.com/jung-kurt/gofpdf"
)

const (
	bannerHeight = 94.0
	xIndent      = 45.0
	taxRate      = 0.09
)

type LineItem struct {
	UnitName       string
	PricePerUnit   int
	UnitsPurchased int
}

func main() {
	lineItems := []LineItem{
		{
			UnitName:       "2x6 Lumber - 8'",
			PricePerUnit:   375,
			UnitsPurchased: 220,
		},
		{
			UnitName:       "Drywall Sheet",
			PricePerUnit:   822,
			UnitsPurchased: 50,
		},
		{
			UnitName:       "Paint",
			PricePerUnit:   1455,
			UnitsPurchased: 3,
		},
		{
			UnitName: "This is a line item with a very long description to test" +
				"that our word wrapping is implemented and working as intended",
			PricePerUnit:   3211,
			UnitsPurchased: 3,
		},
		{
			UnitName:       "Paint",
			PricePerUnit:   5,
			UnitsPurchased: 3300,
		},
		{
			UnitName:       "Paint",
			PricePerUnit:   332,
			UnitsPurchased: 44,
		},
	}

	subTotal := 0
	for _, li := range lineItems {
		subTotal += li.PricePerUnit * li.UnitsPurchased
	}
	tax := int(float64(subTotal) * taxRate)
	total := subTotal + tax
	totalUSD := toUSD(total)

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

	// Line Items - headers
	pdf.SetFont("times", "", 14)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	x, y = xIndent-2.0, y+lineHeight
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		width/2.65+1.5,
		lineHeight,
		"Description",
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignLeft,
		false,
		0,
		"",
	)
	x = x + width/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		100.0,
		lineHeight,
		"Price Per Unit",
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		80.0,
		lineHeight,
		"Quantity",
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	x = width - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		119.5,
		lineHeight,
		"Amount",
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)

	// Line Items - real data
	y = y + lineHeight
	for _, li := range lineItems {
		x, y = lineItem(pdf, x, y, li)
	}

	// Subtotal etc
	x, y = width/1.75, y+lineHeight*2.25
	x, y = trailerLine(pdf, x, y, "Subtotal", subTotal)
	x, y = trailerLine(pdf, x, y, "Tax", tax)
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(x+20.0, y, x+220.0, y)
	y = y + lineHeight*0.5
	x, y = trailerLine(pdf, x, y, "Total", total)

	// Grid
	// drawGrid(pdf)

	if err := pdf.OutputFileAndClose("p4.pdf"); err != nil {
		log.Fatal(err)
	}
}

func trailerLine(pdf *gofpdf.Fpdf, x, y float64, label string, amount int) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		80.0,
		lineHt,
		label,
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	x = w - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(
		119.5,
		lineHt,
		toUSD(amount),
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	y = y + lineHt*1.5
	return origX, y
}

func toUSD(cents int) string {
	centsStr := fmt.Sprintf("%d", cents%100)
	if len(centsStr) < 2 {
		centsStr = "0" + centsStr
	}
	return fmt.Sprintf("$%d.%s", cents/100, centsStr)
}

func lineItem(pdf *gofpdf.Fpdf, x, y float64, lineItem LineItem) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(50, 50, 50)
	pdf.MoveTo(x, y)
	x, y = xIndent-2, y+lineHt*.75
	pdf.MoveTo(x, y)
	pdf.MultiCell(
		w/2.65+1.5,
		lineHt,
		lineItem.UnitName,
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false,
	)
	tmp := pdf.SplitLines([]byte(lineItem.UnitName), w/2.65+1.5)
	maxY := y + float64(len(tmp)-1)*lineHt
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		100,
		lineHt,
		toUSD(lineItem.PricePerUnit),
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	x = x + 100
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		80,
		lineHt,
		fmt.Sprintf("%d", lineItem.UnitsPurchased),
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	x = w - xIndent - 2 - 119.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(
		119.5,
		lineHt,
		toUSD(lineItem.PricePerUnit*lineItem.UnitsPurchased),
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignRight,
		false,
		0,
		"",
	)
	if maxY > y {
		y = maxY
	}
	y = y + lineHt*1.75
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(xIndent-10, y, w-xIndent+10, y)
	return origX, y
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
