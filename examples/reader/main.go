// Package main demonstrates how to use the PDF Reader to read and inspect PDF files.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coregx/gxpdf/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <pdf-file>")
		fmt.Println("\nExample: go run main.go ../../testdata/pdfs/multipage.pdf")
		os.Exit(1)
	}

	pdfFile := os.Args[1]

	// Example 1: Quick info retrieval
	fmt.Println("=== Quick PDF Info ===")
	version, pageCount, err := parser.ReadPDFInfo(pdfFile)
	if err != nil {
		log.Fatalf("Failed to read PDF info: %v", err)
	}
	fmt.Printf("PDF Version: %s\n", version)
	fmt.Printf("Page Count: %d\n\n", pageCount)

	// Example 2: Full PDF reading
	fmt.Println("=== Full PDF Reading ===")
	reader, err := parser.OpenPDF(pdfFile)
	if err != nil {
		log.Fatalf("Failed to open PDF: %v", err)
	}

	// Get catalog
	catalog, err := reader.GetCatalog()
	if err != nil {
		_ = reader.Close()
		log.Fatalf("Failed to get catalog: %v", err)
	}
	fmt.Printf("Catalog: %v\n", catalog)

	// Get page tree
	pages, err := reader.GetPages()
	if err != nil {
		_ = reader.Close()
		log.Fatalf("Failed to get pages: %v", err)
	}
	fmt.Printf("Page Tree Root: %v\n\n", pages)

	// Example 3: Iterate through all pages
	fmt.Println("=== Page Inspection ===")
	count, _ := reader.GetPageCount()
	for i := 0; i < count; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			log.Printf("Failed to get page %d: %v", i, err)
			continue
		}

		fmt.Printf("Page %d:\n", i)

		// Get page dimensions
		mediaBoxObj := page.Get("MediaBox")
		if mediaBoxObj != nil {
			fmt.Printf("  MediaBox: %v\n", mediaBoxObj)
		}

		// Get page contents
		contentsObj := page.Get("Contents")
		if contentsObj != nil {
			fmt.Printf("  Contents: %v\n", contentsObj)
		}

		// Get page resources
		resourcesObj := page.Get("Resources")
		if resourcesObj != nil {
			fmt.Printf("  Resources: %v\n", resourcesObj)
		}

		fmt.Println()
	}

	// Example 4: Access specific objects
	fmt.Println("=== Object Access ===")
	trailer := reader.Trailer()
	if trailer != nil {
		fmt.Printf("Trailer Size: %d\n", trailer.GetInteger("Size"))

		if idArray := trailer.GetArray("ID"); idArray != nil {
			fmt.Printf("Document ID: %v\n", idArray)
		}
	}

	// Example 5: XRef table inspection
	xrefTable := reader.XRefTable()
	if xrefTable != nil {
		fmt.Printf("\nXRef Table:\n")
		fmt.Printf("  Total entries: %d\n", xrefTable.Size())
		fmt.Printf("  In-use entries: %d\n", len(xrefTable.GetInUseEntries()))
		fmt.Printf("  Free entries: %d\n", len(xrefTable.GetFreeEntries()))
	}

	fmt.Println("\n=== Reader State ===")
	fmt.Println(reader.String())

	// Close the reader
	_ = reader.Close()
}
