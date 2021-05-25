package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFOption func(*gofpdf.Fpdf)

func FillColor(c color.RGBA) PDFOption {
	return func(pdf *gofpdf.Fpdf) {
		r, g, b := rgb(c)
		pdf.SetFillColor(r, g, b)
	}
}

func rgb(c color.RGBA) (int, int, int) {
	alpha := float64(c.A) / 255.0
	alphaWhite := int(255 * (1.0 - alpha))
	r := int(float64(c.R)*alpha) + alphaWhite
	g := int(float64(c.G)*alpha) + alphaWhite
	b := int(float64(c.B)*alpha) + alphaWhite
	return r, g, b
}

type PDF struct {
	fpdf *gofpdf.Fpdf
	x, y float64
}

func (p *PDF) Move(xDelta, yDelta float64) {
	p.x, p.y = p.x+xDelta, p.y+yDelta
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) MoveAbs(x, y float64) {
	p.x, p.y = x, y
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) Text(text string) {
	p.fpdf.Text(p.x, p.y, text)
}

func (p *PDF) Polygon(pts []gofpdf.PointType, opts ...PDFOption) {
	for _, opt := range opts {
		opt(p.fpdf)
	}
	p.fpdf.Polygon(pts, "F")
}

// go run main.go -name="Alexander Tsapkov"
func main() {
	name := flag.String("name", "", "the name of the person who completed the course")
	flag.Parse()

	fpdf := gofpdf.New(gofpdf.OrientationLandscape, gofpdf.UnitPoint, gofpdf.PageSizeLetter, "")
	width, height := fpdf.GetPageSize()
	fpdf.AddPage()
	pdf := PDF{
		fpdf: fpdf,
	}

	primary := color.RGBA{103, 60, 79, 255}
	secondary := color.RGBA{103, 60, 79, 220}

	// top and bottom graphics
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{0, height / 9.0},
		{width - (width / 6.0), 0},
	}, FillColor(secondary))
	pdf.Polygon([]gofpdf.PointType{
		{width / 6.0, 0},
		{width, 0},
		{width, height / 9.0},
	}, FillColor(primary))
	pdf.Polygon([]gofpdf.PointType{
		{width, height},
		{width, height - height/8.0},
		{width / 6, height},
	}, FillColor(secondary))
	pdf.Polygon([]gofpdf.PointType{
		{0, height},
		{0, height - height/8.0},
		{width - (width / 6), height},
	}, FillColor(primary))

	fpdf.SetFont("times", "B", 50)
	fpdf.SetTextColor(50, 50, 50)
	pdf.MoveAbs(0, 100)
	_, lineHeight := fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHeight, "Certificate of Completion", gofpdf.AlignCenter)
	pdf.Move(0, lineHeight*2.0)

	fpdf.SetFont("arial", "", 28)
	_, lineHeight = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHeight, "This certificate is awarded to", gofpdf.AlignCenter)
	pdf.Move(0, lineHeight*2.0)

	fpdf.SetFont("times", "B", 42)
	_, lineHeight = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHeight, *name, gofpdf.AlignCenter)
	pdf.Move(0, lineHeight*1.75)

	fpdf.SetFont("arial", "", 22)
	_, lineHeight = fpdf.GetFontSize()
	fpdf.WriteAligned(
		0,
		lineHeight*1.5,
		"For successfully completing all twenty programming exercises"+
			" in the Gophercises programming course for budding Gophers (Go developers)",
		gofpdf.AlignCenter,
	)
	pdf.Move(0, lineHeight*4.5)

	fpdf.ImageOptions("images/jump.png", width/2.0-50.0, pdf.y, 100.0, 0, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")

	pdf.Move(0, 65.0)
	fpdf.SetFillColor(100, 100, 100)
	fpdf.Rect(60.0, pdf.y, 240.0, 1.0, "F")
	fpdf.Rect(490.0, pdf.y, 240.0, 1.0, "F")

	fpdf.SetFont("arial", "", 12)
	pdf.Move(0, lineHeight/1.5)
	fpdf.SetTextColor(100, 100, 100)
	pdf.MoveAbs(60.0+105.0, pdf.y)
	pdf.Text("Date")
	pdf.MoveAbs(490.0+60.0, pdf.y)
	pdf.Text("Instructor - Jon Calhoun")
	pdf.MoveAbs(60.0, pdf.y-lineHeight/1.5)
	fpdf.SetFont("times", "", 22)
	_, lineHeight = fpdf.GetFontSize()
	pdf.Move(0, -lineHeight)
	fpdf.SetTextColor(50, 50, 50)
	year, month, day := time.Now().Date()
	dateStr := fmt.Sprintf("%d/%d/%d", month, day, year)
	fpdf.CellFormat(
		240.0,
		lineHeight,
		dateStr,
		gofpdf.BorderNone,
		gofpdf.LineBreakNone,
		gofpdf.AlignCenter,
		false,
		0,
		"",
	)
	pdf.MoveAbs(490.0, pdf.y)
	sig, err := gofpdf.SVGBasicFileParse("images/signature.svg")
	if err != nil {
		log.Fatal(err)
	}
	pdf.Move(0, -(sig.Ht*.45 - lineHeight))
	fpdf.SVGBasicWrite(&sig, 0.5)

	// Grid
	// drawGrid(fpdf)

	err = fpdf.OutputFileAndClose("certificate.pdf")
	if err != nil {
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
		pdf.SetTextColor(80, 80, 80)
		pdf.Line(0, y, width, y)
		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
	}
}
