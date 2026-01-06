// Package main demonstrates watermark functionality in GxPDF.
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create a new PDF document.
	c := creator.New()
	c.SetTitle("Watermark Example")
	c.SetAuthor("GxPDF Library")

	// Example 1: Centered diagonal watermark (typical for CONFIDENTIAL).
	if err := createDiagonalWatermarkPage(c); err != nil {
		log.Fatalf("Failed to create page 1: %v", err)
	}

	// Example 2: Top-left corner watermark.
	if err := createCornerWatermarkPage(c); err != nil {
		log.Fatalf("Failed to create page 2: %v", err)
	}

	// Example 3: Multiple watermarks on same page.
	if err := createMultiWatermarkPage(c); err != nil {
		log.Fatalf("Failed to create page 3: %v", err)
	}

	// Save the PDF.
	outputPath := "watermark_example.pdf"
	if err := c.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Printf("PDF created successfully: %s\n", outputPath)
	fmt.Println("Open the PDF to see the watermarks!")
}

// createDiagonalWatermarkPage creates a page with a centered diagonal CONFIDENTIAL watermark.
func createDiagonalWatermarkPage(c *creator.Creator) error {
	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	// Add content.
	if err := page.AddText("Example 1: Centered Diagonal Watermark", 72, 700, creator.HelveticaBold, 18); err != nil {
		return fmt.Errorf("add title: %w", err)
	}
	if err := page.AddText("This page has a diagonal CONFIDENTIAL watermark in the center.", 72, 660, creator.Helvetica, 12); err != nil {
		return fmt.Errorf("add description: %w", err)
	}

	// Create and configure watermark.
	wm := creator.NewTextWatermark("CONFIDENTIAL")
	if err := configureWatermark(wm, creator.HelveticaBold, 72, creator.Gray, 0.3, 45, creator.WatermarkCenter); err != nil {
		return fmt.Errorf("configure watermark: %w", err)
	}

	// Apply watermark.
	if err := page.DrawWatermark(wm); err != nil {
		return fmt.Errorf("draw watermark: %w", err)
	}

	return nil
}

// createCornerWatermarkPage creates a page with a DRAFT watermark in top-left corner.
func createCornerWatermarkPage(c *creator.Creator) error {
	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	// Add content.
	if err := page.AddText("Example 2: Top-Left Corner Watermark", 72, 700, creator.HelveticaBold, 18); err != nil {
		return fmt.Errorf("add title: %w", err)
	}
	if err := page.AddText("This page has a DRAFT watermark in the top-left corner.", 72, 660, creator.Helvetica, 12); err != nil {
		return fmt.Errorf("add description: %w", err)
	}

	// Create and configure watermark.
	wm := creator.NewTextWatermark("DRAFT")
	if err := configureWatermark(wm, creator.TimesRoman, 48, creator.Red, 0.4, 0, creator.WatermarkTopLeft); err != nil {
		return fmt.Errorf("configure watermark: %w", err)
	}

	// Apply watermark.
	if err := page.DrawWatermark(wm); err != nil {
		return fmt.Errorf("draw watermark: %w", err)
	}

	return nil
}

// createMultiWatermarkPage creates a page with multiple watermarks in different corners.
func createMultiWatermarkPage(c *creator.Creator) error {
	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	// Add content.
	if err := page.AddText("Example 3: Multiple Watermarks", 72, 700, creator.HelveticaBold, 18); err != nil {
		return fmt.Errorf("add title: %w", err)
	}
	if err := page.AddText("This page demonstrates multiple watermarks in different corners.", 72, 660, creator.Helvetica, 12); err != nil {
		return fmt.Errorf("add description: %w", err)
	}

	// Top-left watermark.
	wmTL := creator.NewTextWatermark("SAMPLE")
	if err := configureWatermark(wmTL, creator.Courier, 36, creator.Blue, 0.25, 0, creator.WatermarkTopLeft); err != nil {
		return fmt.Errorf("configure top-left watermark: %w", err)
	}
	if err := page.DrawWatermark(wmTL); err != nil {
		return fmt.Errorf("draw top-left watermark: %w", err)
	}

	// Bottom-right watermark.
	wmBR := creator.NewTextWatermark("SAMPLE")
	if err := configureWatermark(wmBR, creator.CourierBold, 36, creator.Blue, 0.25, 0, creator.WatermarkBottomRight); err != nil {
		return fmt.Errorf("configure bottom-right watermark: %w", err)
	}
	if err := page.DrawWatermark(wmBR); err != nil {
		return fmt.Errorf("draw bottom-right watermark: %w", err)
	}

	return nil
}

// configureWatermark applies common watermark settings.
func configureWatermark(wm *creator.TextWatermark, font creator.FontName, size float64, color creator.Color, opacity, rotation float64, position creator.WatermarkPosition) error {
	if err := wm.SetFont(font, size); err != nil {
		return err
	}
	if err := wm.SetColor(color); err != nil {
		return err
	}
	if err := wm.SetOpacity(opacity); err != nil {
		return err
	}
	if err := wm.SetRotation(rotation); err != nil {
		return err
	}
	if err := wm.SetPosition(position); err != nil {
		return err
	}
	return nil
}
