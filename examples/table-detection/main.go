// Package main demonstrates table detection from PDF documents.
//
// This example shows how to use the table detector to find table regions
// in PDF pages using both lattice mode (with ruling lines) and stream mode
// (without ruling lines).
package main

import (
	"fmt"
	"os"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/coregx/gxpdf/internal/parser"
	"github.com/coregx/gxpdf/internal/tabledetect"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: table-detection <pdf-file>")
		fmt.Println("Example: table-detection document.pdf")
		os.Exit(1)
	}

	pdfPath := os.Args[1]

	// Open PDF
	fmt.Printf("Opening PDF: %s\n", pdfPath)
	reader, err := parser.OpenPDF(pdfPath)
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

	fmt.Printf("PDF has %d pages\n\n", pageCount)

	// Create extractors and detector
	textExtractor := extractor.NewTextExtractor(reader)
	graphicsParser := extractor.NewGraphicsParser(reader)
	tableDetector := tabledetect.NewTableDetector()

	// Process each page
	for pageNum := 0; pageNum < pageCount; pageNum++ {
		fmt.Printf("=== Page %d ===\n", pageNum+1)

		// Extract text elements
		fmt.Println("Extracting text elements...")
		textElements, err := textExtractor.ExtractFromPage(pageNum)
		if err != nil {
			fmt.Printf("Error extracting text: %v\n", err)
			continue
		}
		fmt.Printf("Found %d text elements\n", len(textElements))

		// Extract graphics elements
		fmt.Println("Extracting graphics elements...")
		graphics, err := graphicsParser.ParseFromPage(pageNum)
		if err != nil {
			fmt.Printf("Error parsing graphics: %v\n", err)
			continue
		}
		fmt.Printf("Found %d graphics elements\n", len(graphics))

		// Detect tables
		fmt.Println("Detecting tables...")
		tables, err := tableDetector.DetectTables(textElements, graphics)
		if err != nil {
			fmt.Printf("Error detecting tables: %v\n", err)
			continue
		}

		// Display results
		if len(tables) == 0 {
			fmt.Println("No tables detected on this page")
		} else {
			fmt.Printf("Found %d table(s):\n", len(tables))

			for i, table := range tables {
				fmt.Printf("\nTable %d:\n", i+1)
				fmt.Printf("  Method: %s\n", table.Method.String())
				fmt.Printf("  Bounds: %s\n", table.Bounds.String())
				fmt.Printf("  Rows: %d\n", table.RowCount())
				fmt.Printf("  Columns: %d\n", table.ColumnCount())

				if table.HasRulingLines && table.Grid != nil {
					fmt.Println("  Grid structure:")
					fmt.Printf("    Row coordinates: %v\n", table.Grid.Rows)
					fmt.Printf("    Column coordinates: %v\n", table.Grid.Columns)
				} else {
					fmt.Println("  Stream mode structure:")
					fmt.Printf("    Row coordinates: %v\n", table.Rows)
					fmt.Printf("    Column coordinates: %v\n", table.Columns)
				}
			}
		}

		// Auto-detection mode
		fmt.Println("\nAuto-detected mode:")
		mode := tableDetector.DetectMode(textElements, graphics)
		fmt.Printf("  Recommended mode: %s\n", mode.String())

		fmt.Println()
	}

	_ = reader.Close()
	fmt.Println("Table detection complete!")
}
