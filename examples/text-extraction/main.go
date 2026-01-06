// Package main demonstrates text extraction with positional information.
//
// This example shows how to extract text from PDF pages with X, Y coordinates.
// Positional information is critical for table extraction and layout analysis.
//
// Usage:
//
//	go run main.go <pdf-file>
//
// Example:
//
//	go run main.go sample.pdf
package main

import (
	"fmt"
	"os"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/coregx/gxpdf/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <pdf-file>")
		fmt.Println()
		fmt.Println("Example: go run main.go sample.pdf")
		os.Exit(1)
	}

	pdfFile := os.Args[1]

	// Open PDF file
	fmt.Printf("Opening PDF: %s\n", pdfFile)
	reader, err := parser.OpenPDF(pdfFile)
	if err != nil {
		fmt.Printf("Error opening PDF: %v\n", err)
		os.Exit(1)
	}

	// Get page count
	pageCount, err := reader.GetPageCount()
	if err != nil {
		_ = reader.Close()
		fmt.Printf("Error getting page count: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("PDF has %d pages\n", pageCount)
	fmt.Println()

	// Create text extractor
	textExtractor := extractor.NewTextExtractor(reader)

	// Extract text from first page
	pageNum := 0
	fmt.Printf("Extracting text from page %d...\n", pageNum+1)
	fmt.Println()

	elements, err := textExtractor.ExtractFromPage(pageNum)
	if err != nil {
		_ = reader.Close()
		fmt.Printf("Error extracting text: %v\n", err)
		os.Exit(1)
	}

	// Display extracted text elements
	if len(elements) == 0 {
		_ = reader.Close()
		fmt.Println("No text found on page.")
		return
	}

	fmt.Printf("Found %d text elements:\n", len(elements))
	fmt.Println()

	// Print first 20 elements with details
	maxElements := 20
	if len(elements) < maxElements {
		maxElements = len(elements)
	}

	fmt.Println("Position details for first", maxElements, "elements:")
	fmt.Println("-------------------------------------------------------")
	for i, elem := range elements[:maxElements] {
		fmt.Printf("[%d] Text: %q\n", i+1, elem.Text)
		fmt.Printf("    Position: (%.2f, %.2f)\n", elem.X, elem.Y)
		fmt.Printf("    Size: %.2f x %.2f\n", elem.Width, elem.Height)
		fmt.Printf("    Font: %s, Size: %.1fpt\n", elem.FontName, elem.FontSize)
		fmt.Println()
	}

	// Print all extracted text (concatenated)
	fmt.Println("-------------------------------------------------------")
	fmt.Println("All extracted text:")
	fmt.Println("-------------------------------------------------------")
	for _, elem := range elements {
		fmt.Print(elem.Text)
	}
	fmt.Println()
	fmt.Println("-------------------------------------------------------")

	// Statistics
	fmt.Println()
	fmt.Printf("Statistics:\n")
	fmt.Printf("  Total text elements: %d\n", len(elements))

	// Count unique fonts
	fonts := make(map[string]int)
	for _, elem := range elements {
		fonts[elem.FontName]++
	}
	fmt.Printf("  Unique fonts: %d\n", len(fonts))
	for font, count := range fonts {
		fmt.Printf("    %s: %d elements\n", font, count)
	}

	// Calculate bounding box
	if len(elements) > 0 {
		minX, minY := elements[0].X, elements[0].Y
		maxX, maxY := elements[0].Right(), elements[0].Top()

		for _, elem := range elements[1:] {
			if elem.X < minX {
				minX = elem.X
			}
			if elem.Y < minY {
				minY = elem.Y
			}
			if elem.Right() > maxX {
				maxX = elem.Right()
			}
			if elem.Top() > maxY {
				maxY = elem.Top()
			}
		}

		fmt.Printf("  Text bounding box:\n")
		fmt.Printf("    Bottom-left: (%.2f, %.2f)\n", minX, minY)
		fmt.Printf("    Top-right: (%.2f, %.2f)\n", maxX, maxY)
		fmt.Printf("    Dimensions: %.2f x %.2f points\n", maxX-minX, maxY-minY)
	}

	_ = reader.Close()
}
