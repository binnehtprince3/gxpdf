// Package main demonstrates table extraction from PDFs using gxpdf.
//
// This example shows how to:
//   - Open a PDF and extract tables using gxpdf.Open
//   - Access table data programmatically
//   - Export tables to CSV, JSON, and Excel formats
//
// Usage:
//
//	go run main.go
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/coregx/gxpdf"
	"github.com/coregx/gxpdf/export"
	"github.com/coregx/gxpdf/internal/models/table"
)

func main() {
	fmt.Println("=== GxPDF Table Extraction Example ===")
	fmt.Println()

	// Example 1: Create a sample table programmatically
	fmt.Println("Example 1: Creating and exporting a sample table")
	exampleCreateAndExportTable()
	fmt.Println()

	// Example 2: Using the new gxpdf API (if you have a PDF file)
	fmt.Println("Example 2: gxpdf API usage")
	exampleGxPDFAPI()
	fmt.Println()

	// Example 3: Export formats comparison
	fmt.Println("Example 3: Export formats comparison")
	exampleExportFormats()
	fmt.Println()

	fmt.Println("All examples completed successfully!")
}

// exampleCreateAndExportTable demonstrates creating a table and exporting to different formats.
func exampleCreateAndExportTable() {
	// Create a sample table
	tbl, err := createSampleTable()
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Printf("Created table: %dx%d (%s mode)\n", tbl.RowCount, tbl.ColCount, tbl.Method)

	// Export using the export package directly
	csvExporter := export.NewCSVExporter()

	// Export to file
	file, err := os.Create("example_table.csv")
	if err != nil {
		log.Printf("Warning: Failed to create CSV file: %v", err)
	} else {
		if err := csvExporter.Export(tbl, file); err != nil {
			log.Printf("Warning: Failed to export CSV: %v", err)
		} else {
			fmt.Println("Exported to example_table.csv")
		}
		file.Close()
	}

	// Export to JSON
	jsonExporter := export.NewJSONExporter().WithPrettyPrint(true).WithMetadata(true)
	file, err = os.Create("example_table.json")
	if err != nil {
		log.Printf("Warning: Failed to create JSON file: %v", err)
	} else {
		if err := jsonExporter.Export(tbl, file); err != nil {
			log.Printf("Warning: Failed to export JSON: %v", err)
		} else {
			fmt.Println("Exported to example_table.json")
		}
		file.Close()
	}

	// Display table content
	fmt.Println("\nTable content:")
	grid := tbl.ToStringGrid()
	for i, row := range grid {
		fmt.Printf("  Row %d: %v\n", i, row)
	}
}

// exampleGxPDFAPI demonstrates the new gxpdf API.
func exampleGxPDFAPI() {
	// This example shows the API, but won't run without a real PDF
	fmt.Println("gxpdf API example (code demonstration):")
	fmt.Println()
	fmt.Println("  // Open a PDF document")
	fmt.Println("  doc, err := gxpdf.Open(\"document.pdf\")")
	fmt.Println("  if err != nil {")
	fmt.Println("      log.Fatal(err)")
	fmt.Println("  }")
	fmt.Println("  defer doc.Close()")
	fmt.Println()
	fmt.Println("  // Get document info")
	fmt.Printf("  // fmt.Printf(\"Pages: %%d\\n\", doc.PageCount())\n")
	fmt.Println()
	fmt.Println("  // Extract all tables")
	fmt.Println("  tables := doc.ExtractTables()")
	fmt.Println("  for _, t := range tables {")
	_, _ = os.Stdout.WriteString("      fmt.Printf(\"Table on page %d: %d rows x %d cols\\n\",\n")
	fmt.Println("          t.PageNumber(), t.RowCount(), t.ColumnCount())")
	fmt.Println()
	fmt.Println("      // Access rows")
	fmt.Println("      for _, row := range t.Rows() {")
	fmt.Println("          fmt.Println(row)")
	fmt.Println("      }")
	fmt.Println()
	fmt.Println("      // Export to CSV")
	fmt.Println("      csv, _ := t.ToCSV()")
	fmt.Println("      fmt.Println(csv)")
	fmt.Println("  }")

	// Show version
	fmt.Printf("\ngxpdf version: %s\n", gxpdf.Version)
}

// exampleExportFormats demonstrates different export options.
func exampleExportFormats() {
	tbl, _ := createSampleTable()

	// CSV with default delimiter
	csvExporter := export.NewCSVExporter()
	csvContent, _ := csvExporter.ExportToString(tbl)
	fmt.Println("CSV format:")
	fmt.Println(csvContent)

	// CSV with semicolon delimiter
	csvSemicolon := export.NewCSVExporter().WithDelimiter(";")
	csvSemicolonContent, _ := csvSemicolon.ExportToString(tbl)
	fmt.Println("CSV with semicolon:")
	fmt.Println(csvSemicolonContent)

	// JSON with pretty print
	jsonExporter := export.NewJSONExporter().WithPrettyPrint(true).WithMetadata(true)
	jsonContent, _ := jsonExporter.ExportToString(tbl)
	fmt.Println("JSON format:")
	fmt.Println(jsonContent)

	// Excel format info
	excelExporter := export.NewExcelExporter()
	fmt.Printf("Excel format: %s\n", excelExporter.ContentType())
	fmt.Printf("File extension: %s\n", excelExporter.FileExtension())

	// Export to buffer (for serving over HTTP)
	var buf bytes.Buffer
	if err := csvExporter.Export(tbl, &buf); err != nil {
		log.Printf("Failed to export to buffer: %v", err)
	} else {
		fmt.Printf("\nExported %d bytes to buffer\n", buf.Len())
	}
}

// createSampleTable creates a sample table for demonstration.
func createSampleTable() (*table.Table, error) {
	tbl, err := table.NewTable(4, 3)
	if err != nil {
		return nil, err
	}

	// Set metadata
	tbl.Method = "Lattice"
	tbl.PageNum = 0
	tbl.Bounds = table.NewRectangle(50, 100, 500, 400)

	// Header row
	tbl.SetCell(0, 0, table.NewCell("Product", 0, 0))
	tbl.SetCell(0, 1, table.NewCell("Quantity", 0, 1))
	tbl.SetCell(0, 2, table.NewCell("Price", 0, 2))

	// Data rows
	tbl.SetCell(1, 0, table.NewCell("Widget A", 1, 0))
	tbl.SetCell(1, 1, table.NewCell("10", 1, 1).WithAlignment(table.AlignRight))
	tbl.SetCell(1, 2, table.NewCell("$99.99", 1, 2).WithAlignment(table.AlignRight))

	tbl.SetCell(2, 0, table.NewCell("Widget B", 2, 0))
	tbl.SetCell(2, 1, table.NewCell("5", 2, 1).WithAlignment(table.AlignRight))
	tbl.SetCell(2, 2, table.NewCell("$149.99", 2, 2).WithAlignment(table.AlignRight))

	tbl.SetCell(3, 0, table.NewCell("Total", 3, 0))
	tbl.SetCell(3, 1, table.NewCell("15", 3, 1).WithAlignment(table.AlignRight))
	tbl.SetCell(3, 2, table.NewCell("$1,249.85", 3, 2).WithAlignment(table.AlignRight))

	return tbl, nil
}
