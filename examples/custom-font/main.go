package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	fmt.Println("Custom Font Example")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("This example demonstrates how to use custom TTF/OTF fonts in GxPDF.")
	fmt.Println()

	// Note: This example requires a TTF font file.
	// For testing purposes, you can use any TTF font from your system:
	// - Windows: C:\\Windows\\Fonts\\arial.ttf
	// - Linux: /usr/share/fonts/truetype/dejavu/DejaVuSans.ttf
	// - macOS: /Library/Fonts/Arial.ttf

	fontPath := "path/to/your/font.ttf" // Update this path.

	// Load custom font.
	fmt.Printf("Loading font: %s\n", fontPath)
	customFont, err := creator.LoadFont(fontPath)
	if err != nil {
		log.Printf("Error loading font: %v\n", err)
		log.Println("Please update the fontPath variable with a valid TTF file path.")
		return
	}

	fmt.Println("Font loaded successfully!")
	fmt.Printf("PostScript name: %s\n", customFont.PostScriptName())
	fmt.Printf("Units per em: %d\n", customFont.UnitsPerEm())
	fmt.Println()

	// Mark some characters as used.
	testText := "Hello, World! Текст на русском языке. 你好世界"
	customFont.UseString(testText)

	// Measure text width.
	width := customFont.MeasureString(testText, 12)
	fmt.Printf("Text: %s\n", testText)
	fmt.Printf("Width at 12pt: %.2f points\n", width)
	fmt.Println()

	// Build font subset.
	fmt.Println("Building font subset...")
	if err := customFont.Build(); err != nil {
		log.Fatalf("Error building font subset: %v", err)
	}

	fmt.Println("Font subset built successfully!")
	fmt.Println()

	// Demonstrate measurement.
	fmt.Println("Measuring different strings:")
	testStrings := []string{
		"A",
		"Hello",
		"Привет",
		"The quick brown fox jumps over the lazy dog",
	}

	for _, str := range testStrings {
		w := customFont.MeasureString(str, 12)
		fmt.Printf("  %-50s: %.2f pt\n", str, w)
	}
	fmt.Println()

	fmt.Println("Custom font example completed successfully!")
	fmt.Println()
	fmt.Println("Note: Full PDF generation with custom fonts requires PDF writing support.")
	fmt.Println("This example demonstrates the font loading and measurement functionality.")
}
