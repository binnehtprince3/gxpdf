// Package main demonstrates page rotation in PDF documents.
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	fmt.Println("=== Page Rotation Example ===")

	// Example 1: Create new PDF with rotated pages.
	createRotatedPDF()

	// Example 2: Rotate existing PDF pages.
	rotateExistingPDF()

	fmt.Println("Examples completed successfully!")
}

// createRotatedPDF demonstrates creating pages with different rotations.
func createRotatedPDF() {
	fmt.Println("\n1. Creating PDF with rotated pages...")

	c := creator.New()
	c.SetTitle("Page Rotation Example")
	c.SetAuthor("gxpdf")

	// Create pages with different rotations.
	createPortraitPage(c)
	createLandscapePage(c)
	createUpsideDownPage(c)
	createLandscapeReversePage(c)

	// Save the PDF.
	if err := c.WriteToFile("rotation_new.pdf"); err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Println("✓ Created rotation_new.pdf with 4 pages (0°, 90°, 180°, 270°)")
}

func createPortraitPage(c *creator.Creator) {
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}
	if err := page.AddText("Portrait Page (0°)", 100, 700, creator.HelveticaBold, 24); err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}
}

func createLandscapePage(c *creator.Creator) {
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}
	if err := page.Rotate(90); err != nil {
		log.Fatalf("Failed to rotate page: %v", err)
	}
	if err := page.AddText("Landscape Page (90°)", 100, 400, creator.HelveticaBold, 24); err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}
}

func createUpsideDownPage(c *creator.Creator) {
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}
	if err := page.Rotate(180); err != nil {
		log.Fatalf("Failed to rotate page: %v", err)
	}
	if err := page.AddText("Upside Down Page (180°)", 100, 700, creator.HelveticaBold, 24); err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}
}

func createLandscapeReversePage(c *creator.Creator) {
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}
	if err := page.Rotate(270); err != nil {
		log.Fatalf("Failed to rotate page: %v", err)
	}
	if err := page.AddText("Landscape Reverse (270°)", 100, 400, creator.HelveticaBold, 24); err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}
}

// rotateExistingPDF demonstrates rotating pages in an existing PDF.
func rotateExistingPDF() {
	fmt.Println("\n2. Rotating pages in existing PDF...")

	// Create sample PDF.
	createSamplePDF()

	// Rotate pages in the PDF.
	rotateAllPages()

	fmt.Println("✓ Created rotation_original.pdf (3 portrait pages)")
	fmt.Println("✓ Created rotation_modified.pdf (all pages rotated to landscape)")
}

func createSamplePDF() {
	c := creator.New()
	c.SetTitle("Original Document")

	for i := 1; i <= 3; i++ {
		page, err := c.NewPage()
		if err != nil {
			log.Fatalf("Failed to create page %d: %v", i, err)
		}
		text := fmt.Sprintf("Page %d (Portrait)", i)
		if err := page.AddText(text, 100, 700, creator.Helvetica, 18); err != nil {
			log.Fatalf("Failed to add text: %v", err)
		}
	}

	if err := c.WriteToFile("rotation_original.pdf"); err != nil {
		log.Fatalf("Failed to write original PDF: %v", err)
	}
}

func rotateAllPages() {
	app, err := creator.NewAppender("rotation_original.pdf")
	if err != nil {
		log.Fatalf("Failed to open PDF: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Rotate all pages to landscape.
	for i := 0; i < app.PageCount(); i++ {
		rotatePage(app, i)
	}

	// Save the modified PDF.
	//nolint:gocritic // Example uses Fatal for simplicity
	if err := app.WriteToFile("rotation_modified.pdf"); err != nil {
		_ = app.Close() // Close before fatal
		log.Fatalf("Failed to write modified PDF: %v", err)
	}
}

func rotatePage(app *creator.Appender, pageIndex int) {
	page, err := app.GetPage(pageIndex)
	if err != nil {
		log.Printf("Failed to get page %d: %v", pageIndex, err)
		return
	}

	if err := page.Rotate(90); err != nil {
		log.Printf("Failed to rotate page %d: %v", pageIndex, err)
		return
	}

	// Add watermark to rotated page.
	text := fmt.Sprintf("ROTATED 90° - Page %d", pageIndex+1)
	if err := page.AddText(text, 300, 200, creator.HelveticaBold, 16); err != nil {
		log.Printf("Failed to add watermark: %v", err)
	}
}
