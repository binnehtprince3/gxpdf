package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create a new PDF creator.
	c := creator.New()

	// Set document metadata.
	c.SetTitle("PDF Annotations Example")
	c.SetAuthor("GxPDF Library")
	c.SetSubject("Demonstrating all annotation types")

	// Create a new page.
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	// Build page content.
	if err := addTitle(page); err != nil {
		log.Fatalf("Failed to add title: %v", err)
	}

	if err := addTextAnnotations(page); err != nil {
		log.Fatalf("Failed to add text annotations: %v", err)
	}

	if err := addMarkupAnnotations(page); err != nil {
		log.Fatalf("Failed to add markup annotations: %v", err)
	}

	if err := addStampAnnotations(page); err != nil {
		log.Fatalf("Failed to add stamp annotations: %v", err)
	}

	if err := addSummary(page); err != nil {
		log.Fatalf("Failed to add summary: %v", err)
	}

	// Write to file.
	outputPath := "annotations_demo.pdf"
	if err := c.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Printf("PDF with annotations created successfully: %s\n", outputPath)
	fmt.Println("Open the file in Adobe Acrobat or another PDF viewer that supports annotations.")
}

// addTitle adds the main title to the page.
func addTitle(page *creator.Page) error {
	return page.AddText("PDF Annotations Demo", 50, 800, creator.HelveticaBold, 18)
}

// addTextAnnotations adds sticky note examples to the page.
func addTextAnnotations(page *creator.Page) error {
	y := 750.0

	if err := page.AddText("1. Text Annotations (Sticky Notes):", 50, y, creator.HelveticaBold, 14); err != nil {
		return err
	}

	y -= 30
	if err := page.AddText("Click the icons to see comments:", 70, y, creator.Helvetica, 12); err != nil {
		return err
	}

	// Yellow sticky note.
	note1 := creator.NewTextAnnotation(250, y-5, "This is a yellow sticky note")
	note1.SetAuthor("Alice").SetColor(creator.Yellow)
	if err := page.AddTextAnnotation(note1); err != nil {
		return err
	}

	// Red sticky note (important).
	note2 := creator.NewTextAnnotation(280, y-5, "IMPORTANT: This requires attention!")
	note2.SetAuthor("Manager").SetColor(creator.Red).SetOpen(true)
	return page.AddTextAnnotation(note2)
}

// addMarkupAnnotations adds highlight, underline, and strikeout examples.
func addMarkupAnnotations(page *creator.Page) error {
	y := 660.0

	if err := page.AddText("2. Markup Annotations:", 50, y, creator.HelveticaBold, 14); err != nil {
		return err
	}

	// Highlight.
	y -= 30
	if err := page.AddText("This text is highlighted in yellow", 70, y, creator.Helvetica, 12); err != nil {
		return err
	}
	highlight := creator.NewHighlightAnnotation(70, y-3, 240, y+12)
	highlight.SetColor(creator.Yellow).SetAuthor("Bob").SetNote("Key point")
	if err := page.AddHighlightAnnotation(highlight); err != nil {
		return err
	}

	// Underline.
	y -= 30
	if err := page.AddText("This text is underlined in blue", 70, y, creator.Helvetica, 12); err != nil {
		return err
	}
	underline := creator.NewUnderlineAnnotation(70, y-3, 220, y+12)
	underline.SetColor(creator.Blue).SetAuthor("Bob")
	if err := page.AddUnderlineAnnotation(underline); err != nil {
		return err
	}

	// Strikeout.
	y -= 30
	if err := page.AddText("This text is struck out in red", 70, y, creator.Helvetica, 12); err != nil {
		return err
	}
	strikeout := creator.NewStrikeOutAnnotation(70, y-3, 200, y+12)
	strikeout.SetColor(creator.Red).SetNote("Obsolete information")
	return page.AddStrikeOutAnnotation(strikeout)
}

// addStampAnnotations adds stamp examples to the page.
func addStampAnnotations(page *creator.Page) error {
	y := 510.0

	if err := page.AddText("3. Stamp Annotations:", 50, y, creator.HelveticaBold, 14); err != nil {
		return err
	}

	y -= 40

	// Approved stamp.
	stamp1 := creator.NewStampAnnotation(70, y, 100, 40, creator.StampApproved)
	stamp1.SetColor(creator.Green).SetAuthor("Manager").SetNote("Approved on 2025-01-06")
	if err := page.AddStampAnnotation(stamp1); err != nil {
		return err
	}

	// Draft stamp.
	stamp2 := creator.NewStampAnnotation(200, y, 100, 40, creator.StampDraft)
	stamp2.SetColor(creator.Yellow)
	if err := page.AddStampAnnotation(stamp2); err != nil {
		return err
	}

	// Confidential stamp.
	stamp3 := creator.NewStampAnnotation(330, y, 100, 40, creator.StampConfidential)
	stamp3.SetColor(creator.Red)
	return page.AddStampAnnotation(stamp3)
}

// addSummary adds the summary section to the page.
func addSummary(page *creator.Page) error {
	y := 390.0

	if err := page.AddText("Summary:", 50, y, creator.HelveticaBold, 14); err != nil {
		return err
	}

	y -= 25
	summary := []string{
		"- Text annotations appear as clickable icons",
		"- Highlight/underline/strikeout mark text",
		"- Stamps show predefined status messages",
		"- All annotations support colors and authors",
		"- Open the PDF in a viewer that supports annotations",
	}

	for _, line := range summary {
		if err := page.AddText(line, 70, y, creator.Helvetica, 11); err != nil {
			return err
		}
		y -= 20
	}

	return nil
}
