// Package main demonstrates how to extract images from PDF documents.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/coregx/gxpdf"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <pdf-file> [output-dir]")
		fmt.Println("\nExamples:")
		fmt.Println("  go run main.go document.pdf")
		fmt.Println("  go run main.go document.pdf ./extracted_images")
		fmt.Println("\nThis example extracts all images from a PDF and saves them to disk.")
		os.Exit(1)
	}

	pdfFile := os.Args[1]
	outputDir := "."
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Example 1: Extract all images from document
	fmt.Println("=== Example 1: Extract All Images from Document ===")
	extractAllImages(pdfFile, outputDir)

	// Example 2: Extract images from specific page
	fmt.Println("\n=== Example 2: Extract Images from Specific Page ===")
	extractImagesFromPage(pdfFile, 0, outputDir) // First page

	// Example 3: Process images without saving
	fmt.Println("\n=== Example 3: Process Images (Metadata Only) ===")
	processImageMetadata(pdfFile)
}

// extractAllImages extracts all images from all pages in the PDF.
func extractAllImages(pdfFile, outputDir string) {
	// Open PDF
	doc, err := gxpdf.Open(pdfFile)
	if err != nil {
		log.Fatalf("Failed to open PDF: %v", err)
	}
	defer doc.Close()

	// Get all images
	images := doc.GetImages()
	fmt.Printf("Found %d images in document\n", len(images))

	// Save each image
	for i, img := range images {
		// Determine file extension based on filter
		ext := ".png"
		if img.Filter() == "/DCTDecode" {
			ext = ".jpg"
		}

		filename := filepath.Join(outputDir, fmt.Sprintf("image_%d%s", i, ext))

		if err := img.SaveToFile(filename); err != nil {
			log.Printf("Failed to save image %d: %v", i, err)
			continue
		}

		fmt.Printf("  Saved: %s (%dx%d, %s, %s)\n",
			filename, img.Width(), img.Height(), img.ColorSpace(), img.Filter())
	}
}

// extractImagesFromPage extracts images from a specific page.
func extractImagesFromPage(pdfFile string, pageIndex int, outputDir string) {
	// Open PDF
	doc, err := gxpdf.Open(pdfFile)
	if err != nil {
		log.Fatalf("Failed to open PDF: %v", err)
	}

	// Get page
	page := doc.Page(pageIndex)
	if page == nil {
		_ = doc.Close()
		log.Fatalf("Page %d not found", pageIndex)
	}

	// Get images from page
	images := page.GetImages()
	fmt.Printf("Found %d images on page %d\n", len(images), page.Number())

	// Save each image
	for i, img := range images {
		// Determine file extension based on filter
		ext := ".png"
		if img.Filter() == "/DCTDecode" {
			ext = ".jpg"
		}

		filename := filepath.Join(outputDir, fmt.Sprintf("page%d_image_%d%s", page.Number(), i, ext))

		if err := img.SaveToFile(filename); err != nil {
			log.Printf("Failed to save image %d: %v", i, err)
			continue
		}

		fmt.Printf("  Saved: %s (%dx%d)\n", filename, img.Width(), img.Height())
	}

	_ = doc.Close()
}

// processImageMetadata processes images without saving them to disk.
func processImageMetadata(pdfFile string) {
	// Open PDF
	doc, err := gxpdf.Open(pdfFile)
	if err != nil {
		log.Fatalf("Failed to open PDF: %v", err)
	}
	defer doc.Close()

	// Process each page
	for _, page := range doc.Pages() {
		images := page.GetImages()

		if len(images) == 0 {
			continue
		}

		fmt.Printf("Page %d:\n", page.Number())

		for i, img := range images {
			fmt.Printf("  Image %d:\n", i)
			fmt.Printf("    Name: %s\n", img.Name())
			fmt.Printf("    Dimensions: %dx%d pixels\n", img.Width(), img.Height())
			fmt.Printf("    Color Space: %s\n", img.ColorSpace())
			fmt.Printf("    Bits per Component: %d\n", img.BitsPerComponent())
			fmt.Printf("    Filter: %s\n", img.Filter())

			// You can also convert to Go image for processing
			goImg, err := img.ToGoImage()
			if err != nil {
				log.Printf("    Failed to convert to Go image: %v", err)
				continue
			}

			bounds := goImg.Bounds()
			fmt.Printf("    Go Image Bounds: %v\n", bounds)
		}
	}
}
