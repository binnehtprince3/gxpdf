// Package main demonstrates PDF bookmark/outline support.
//
// This example creates a PDF with a hierarchical bookmark structure:
// - Chapter 1
//   - Section 1.1
//   - Section 1.2
//
// - Chapter 2
//   - Section 2.1
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create new document.
	c := creator.New()
	c.SetTitle("Document with Bookmarks")
	c.SetAuthor("GxPDF Example")

	// Create pages for chapters.
	page1, err := c.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	page2, err := c.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	page3, err := c.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	page4, err := c.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	// Add content to pages.
	addChapter1Content(page1)
	addSection11Content(page2)
	addSection12Content(page3)
	addChapter2Content(page4)

	// Add bookmarks (hierarchical structure).
	if err := c.AddBookmark("Chapter 1", 0, 0); err != nil {
		log.Fatal(err)
	}
	if err := c.AddBookmark("Section 1.1", 1, 1); err != nil {
		log.Fatal(err)
	}
	if err := c.AddBookmark("Section 1.2", 2, 1); err != nil {
		log.Fatal(err)
	}
	if err := c.AddBookmark("Chapter 2", 3, 0); err != nil {
		log.Fatal(err)
	}

	// Verify bookmarks were added.
	bookmarks := c.Bookmarks()
	fmt.Printf("Added %d bookmarks:\n", len(bookmarks))
	for i, bm := range bookmarks {
		indent := ""
		for j := 0; j < bm.Level; j++ {
			indent += "  "
		}
		fmt.Printf("%d. %s%s (page %d, level %d)\n",
			i+1, indent, bm.Title, bm.PageIndex+1, bm.Level)
	}

	// Write PDF.
	outputPath := "bookmarks_example.pdf"
	if err := c.WriteToFile(outputPath); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nPDF created successfully: %s\n", outputPath)
	fmt.Println("Note: Bookmark rendering in PDF will be implemented in Phase 3 (Writer)")
}

func addChapter1Content(page *creator.Page) {
	if err := page.AddText("Chapter 1", 100, 700, creator.HelveticaBold, 24); err != nil {
		log.Fatal(err)
	}
	if err := page.AddText("Introduction to GxPDF", 100, 660, creator.Helvetica, 14); err != nil {
		log.Fatal(err)
	}
}

func addSection11Content(page *creator.Page) {
	if err := page.AddText("Section 1.1", 100, 700, creator.HelveticaBold, 18); err != nil {
		log.Fatal(err)
	}
	if err := page.AddText("Getting Started", 100, 670, creator.Helvetica, 14); err != nil {
		log.Fatal(err)
	}
}

func addSection12Content(page *creator.Page) {
	if err := page.AddText("Section 1.2", 100, 700, creator.HelveticaBold, 18); err != nil {
		log.Fatal(err)
	}
	if err := page.AddText("Basic Usage", 100, 670, creator.Helvetica, 14); err != nil {
		log.Fatal(err)
	}
}

func addChapter2Content(page *creator.Page) {
	if err := page.AddText("Chapter 2", 100, 700, creator.HelveticaBold, 24); err != nil {
		log.Fatal(err)
	}
	if err := page.AddText("Advanced Features", 100, 660, creator.Helvetica, 14); err != nil {
		log.Fatal(err)
	}
}
