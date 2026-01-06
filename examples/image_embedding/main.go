package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create a Creator instance.
	c := creator.New()
	c.SetTitle("Image Embedding Example")
	c.SetAuthor("GxPDF")

	// Create a page.
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	// Add title.
	err = page.AddText("JPEG Image Embedding Demo", 100, 750, creator.HelveticaBold, 24)
	if err != nil {
		log.Fatalf("Failed to add title: %v", err)
	}

	// Create a test JPEG image (red square).
	img := createTestImage(200, 150, color.RGBA{255, 0, 0, 255})

	// Example 1: Draw image with specific dimensions.
	err = page.AddText("1. Fixed size (200x150 pts):", 100, 700, creator.Helvetica, 12)
	if err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}

	err = page.DrawImage(img, 100, 500, 200, 150)
	if err != nil {
		log.Fatalf("Failed to draw image: %v", err)
	}

	// Example 2: Draw image scaled to fit (maintains aspect ratio).
	err = page.AddText("2. Fit to box (150x150 pts):", 350, 700, creator.Helvetica, 12)
	if err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}

	err = page.DrawImageFit(img, 350, 500, 150, 150)
	if err != nil {
		log.Fatalf("Failed to draw image fit: %v", err)
	}

	// Example 3: Multiple images of different sizes.
	err = page.AddText("3. Multiple sizes:", 100, 450, creator.Helvetica, 12)
	if err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}

	// Small.
	err = page.DrawImage(img, 100, 350, 80, 60)
	if err != nil {
		log.Fatalf("Failed to draw small image: %v", err)
	}

	// Medium.
	err = page.DrawImage(img, 200, 350, 120, 90)
	if err != nil {
		log.Fatalf("Failed to draw medium image: %v", err)
	}

	// Large.
	err = page.DrawImage(img, 340, 350, 160, 120)
	if err != nil {
		log.Fatalf("Failed to draw large image: %v", err)
	}

	// Add footer.
	err = page.AddText("Images are embedded as JPEG (DCTDecode) for efficient storage.",
		100, 50, creator.Helvetica, 10)
	if err != nil {
		log.Fatalf("Failed to add footer: %v", err)
	}

	// Write to file.
	outputPath := "image_embedding_example.pdf"
	fmt.Printf("Writing PDF to: %s\n", outputPath)

	err = c.WriteToFile(outputPath)
	if err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Printf("✓ PDF created successfully: %s\n", outputPath)
	fmt.Printf("  - 1 page with 5 embedded JPEG images\n")
	fmt.Printf("  - Demonstrates DrawImage() and DrawImageFit()\n")
}

// createTestImage creates a JPEG image for testing.
func createTestImage(width, height int, c color.Color) *creator.Image {
	// Create image.
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgba.Set(x, y, c)
		}
	}

	// Encode to JPEG.
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, rgba, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatalf("Failed to encode JPEG: %v", err)
	}

	// Load as creator.Image.
	creatorImg, err := creator.LoadImageFromReader(&buf)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	return creatorImg
}
