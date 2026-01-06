// Package main demonstrates CMYK color support in the Creator API.
//
// CMYK (Cyan, Magenta, Yellow, blacK) is a subtractive color model used in
// professional printing. This example shows how to use CMYK colors for text
// and graphics.
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create a new PDF creator
	c := creator.New()
	c.SetTitle("CMYK Colors Example")
	c.SetAuthor("GxPDF Creator API")

	// Create a page
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	// --- Text with CMYK colors ---
	y := 750.0

	// Pure CMYK colors
	page.AddTextColorCMYK("Pure Cyan (C=100%)", 50, y, creator.HelveticaBold, 14, creator.CMYKCyan)
	y -= 30

	page.AddTextColorCMYK("Pure Magenta (M=100%)", 50, y, creator.HelveticaBold, 14, creator.CMYKMagenta)
	y -= 30

	page.AddTextColorCMYK("Pure Yellow (Y=100%)", 50, y, creator.HelveticaBold, 14, creator.CMYKYellow)
	y -= 30

	page.AddTextColorCMYK("Pure Black (K=100%)", 50, y, creator.HelveticaBold, 14, creator.CMYKBlack)
	y -= 40

	// Mixed CMYK colors
	page.AddTextColorCMYK("Red (M+Y)", 50, y, creator.Helvetica, 12, creator.CMYKRed)
	y -= 25

	page.AddTextColorCMYK("Green (C+Y)", 50, y, creator.Helvetica, 12, creator.CMYKGreen)
	y -= 25

	page.AddTextColorCMYK("Blue (C+M)", 50, y, creator.Helvetica, 12, creator.CMYKBlue)
	y -= 40

	// Custom CMYK color
	customCMYK := creator.NewColorCMYK(0.5, 0.3, 0.7, 0.1)
	page.AddTextColorCMYK("Custom CMYK (C=50% M=30% Y=70% K=10%)", 50, y, creator.Helvetica, 12, customCMYK)
	y -= 40

	// --- Graphics with CMYK colors ---

	// Rectangle with CMYK fill
	rectCMYK := creator.CMYKCyan
	rectOpts := &creator.RectOptions{
		FillColorCMYK:   &rectCMYK,
		StrokeColorCMYK: &creator.CMYKBlack,
		StrokeWidth:     2,
	}
	page.DrawRect(50, y-50, 150, 40, rectOpts)
	page.AddText("CMYK Fill", 80, y-30, creator.Helvetica, 10)
	y -= 100

	// Circle with CMYK colors
	magenta := creator.CMYKMagenta
	black := creator.CMYKBlack
	circleOpts := &creator.CircleOptions{
		FillColorCMYK:   &magenta,
		StrokeColorCMYK: &black,
		StrokeWidth:     2,
	}
	page.DrawCircle(125, y-30, 30, circleOpts)
	y -= 80

	// Polygon with CMYK colors
	yellow := creator.CMYKYellow
	polygonOpts := &creator.PolygonOptions{
		FillColorCMYK:   &yellow,
		StrokeColorCMYK: &black,
		StrokeWidth:     2,
	}
	vertices := []creator.Point{
		{X: 50, Y: y},
		{X: 100, Y: y - 40},
		{X: 150, Y: y},
		{X: 125, Y: y + 20},
		{X: 75, Y: y + 20},
	}
	page.DrawPolygon(vertices, polygonOpts)
	y -= 80

	// --- Color conversion example ---
	y -= 20
	page.AddText("Color Conversion Examples:", 50, y, creator.HelveticaBold, 12)
	y -= 25

	// RGB to CMYK
	rgbColor := creator.Color{R: 1.0, G: 0.5, B: 0.0} // Orange
	cmykFromRGB := rgbColor.ToCMYK()
	page.AddTextColorCMYK(
		fmt.Sprintf("RGB(%.1f,%.1f,%.1f) → CMYK(%.2f,%.2f,%.2f,%.2f)",
			rgbColor.R, rgbColor.G, rgbColor.B,
			cmykFromRGB.C, cmykFromRGB.M, cmykFromRGB.Y, cmykFromRGB.K),
		50, y, creator.Courier, 9, cmykFromRGB)
	y -= 25

	// CMYK to RGB
	cmykColor := creator.NewColorCMYK(0.0, 0.5, 1.0, 0.0) // Yellowish-Orange
	rgbFromCMYK := cmykColor.ToRGB()
	page.AddTextColor(
		fmt.Sprintf("CMYK(%.2f,%.2f,%.2f,%.2f) → RGB(%.2f,%.2f,%.2f)",
			cmykColor.C, cmykColor.M, cmykColor.Y, cmykColor.K,
			rgbFromCMYK.R, rgbFromCMYK.G, rgbFromCMYK.B),
		50, y, creator.Courier, 9, rgbFromCMYK)

	// Write to file
	err = c.WriteToFile("cmyk_colors.pdf")
	if err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Println("PDF created successfully: cmyk_colors.pdf")
	fmt.Println("\nThis PDF demonstrates:")
	fmt.Println("- CMYK text colors")
	fmt.Println("- CMYK fill and stroke colors for shapes")
	fmt.Println("- RGB ↔ CMYK color conversion")
	fmt.Println("\nCMYK colors are ideal for professional printing workflows")
}
