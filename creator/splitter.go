// Package creator provides a high-level API for creating and modifying PDF documents.
package creator

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/coregx/gxpdf/internal/document"
	"github.com/coregx/gxpdf/internal/reader"
)

// Splitter provides functionality to split PDF files into smaller parts.
//
// Use Splitter when you need to:
// - Split a PDF into individual pages (one file per page)
// - Split a PDF into multiple files based on page ranges
// - Extract specific pages into a new document
//
// Example - Split into individual pages:
//
//	splitter, _ := creator.NewSplitter("large.pdf")
//	defer splitter.Close()
//	splitter.Split("output/") // Creates page_001.pdf, page_002.pdf, etc.
//
// Example - Split by ranges:
//
//	splitter, _ := creator.NewSplitter("large.pdf")
//	defer splitter.Close()
//	splitter.SplitByRanges(
//	    creator.PageRange{Start: 1, End: 5, Output: "part1.pdf"},
//	    creator.PageRange{Start: 6, End: 10, Output: "part2.pdf"},
//	)
//
// Example - Extract specific pages:
//
//	splitter, _ := creator.NewSplitter("large.pdf")
//	doc, _ := splitter.ExtractPages(1, 3, 5, 7)
type Splitter struct {
	// Source document path.
	sourcePath string

	// Source document.
	sourceDoc *document.Document

	// PDF reader for cleanup.
	reader *reader.PdfReader

	// Output filename pattern for Split().
	filenamePattern string
}

// PageRange defines a range of pages to extract.
//
// Page numbers are 1-based (1 = first page).
// Start and End are inclusive.
type PageRange struct {
	Start  int    // First page (1-based, inclusive)
	End    int    // Last page (1-based, inclusive)
	Output string // Output file path
}

// NewSplitter creates a new Splitter for the specified PDF file.
//
// Parameters:
//   - path: Path to the PDF file to split
//
// Returns an error if:
//   - File cannot be opened
//   - File is not a valid PDF
//
// Example:
//
//	splitter, err := creator.NewSplitter("large.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer splitter.Close()
func NewSplitter(path string) (*Splitter, error) {
	// Open and reconstruct document.
	doc, r, err := openAndReconstruct(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}

	return &Splitter{
		sourcePath:      path,
		sourceDoc:       doc,
		reader:          r,
		filenamePattern: "page_%03d.pdf",
	}, nil
}

// SetFilenamePattern sets the output filename pattern.
//
// The pattern must contain a single %d or %03d formatter for page number.
//
// Default pattern: "page_%03d.pdf"
//
// Example:
//
//	splitter.SetFilenamePattern("output_%04d.pdf")
//	// Produces: output_0001.pdf, output_0002.pdf, etc.
func (s *Splitter) SetFilenamePattern(pattern string) {
	s.filenamePattern = pattern
}

// Split splits the PDF into individual page files.
//
// Each page is written to a separate PDF file in the specified directory.
// Filenames are generated using the filename pattern (default: page_001.pdf).
//
// Parameters:
//   - outputDir: Directory where individual page files will be written
//
// Returns an error if:
//   - Output directory cannot be created
//   - Any page cannot be written
//
// Example:
//
//	splitter.Split("output/") // Creates page_001.pdf, page_002.pdf, etc.
func (s *Splitter) Split(outputDir string) error {
	ctx := context.Background()
	return s.SplitContext(ctx, outputDir)
}

// SplitContext splits with context support.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	splitter.SplitContext(ctx, "output/")
func (s *Splitter) SplitContext(ctx context.Context, outputDir string) error {
	// Check context.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate.
	pageCount := s.sourceDoc.PageCount()
	if pageCount == 0 {
		return fmt.Errorf("source document has no pages")
	}

	// Split each page.
	for i := 1; i <= pageCount; i++ {
		// Check context before each page.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Generate output filename.
		filename := fmt.Sprintf(s.filenamePattern, i)
		outputPath := filepath.Join(outputDir, filename)

		// Extract and write page.
		if err := s.extractAndWrite(outputPath, i); err != nil {
			return fmt.Errorf("failed to split page %d: %w", i, err)
		}
	}

	return nil
}

// SplitByRanges splits the PDF by page ranges.
//
// Each range is written to the specified output file.
// Page numbers are 1-based and inclusive.
//
// Parameters:
//   - ranges: One or more PageRange specifications
//
// Returns an error if:
//   - No ranges specified
//   - Any range is invalid
//   - Any output file cannot be written
//
// Example:
//
//	splitter.SplitByRanges(
//	    creator.PageRange{Start: 1, End: 5, Output: "part1.pdf"},
//	    creator.PageRange{Start: 6, End: 10, Output: "part2.pdf"},
//	)
func (s *Splitter) SplitByRanges(ranges ...PageRange) error {
	ctx := context.Background()
	return s.SplitByRangesContext(ctx, ranges...)
}

// SplitByRangesContext splits by ranges with context support.
func (s *Splitter) SplitByRangesContext(ctx context.Context, ranges ...PageRange) error {
	// Check context.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate.
	if len(ranges) == 0 {
		return fmt.Errorf("no ranges specified")
	}

	// Validate all ranges first.
	if err := s.validateRanges(ranges); err != nil {
		return err
	}

	// Process each range.
	for i, r := range ranges {
		// Check context before each range.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Extract pages in range.
		pages := make([]int, 0, r.End-r.Start+1)
		for p := r.Start; p <= r.End; p++ {
			pages = append(pages, p)
		}

		// Extract and write range.
		if err := s.extractPages(r.Output, pages); err != nil {
			return fmt.Errorf("failed to split range %d: %w", i+1, err)
		}
	}

	return nil
}

// ExtractPages extracts specific pages into a new document.
//
// This creates a new in-memory document with only the specified pages.
// The returned document can be modified or written to a file.
//
// Page numbers are 1-based.
//
// Parameters:
//   - pages: Page numbers to extract (1-based)
//
// Returns an error if:
//   - No pages specified
//   - Any page number is invalid
//
// Example:
//
//	doc, err := splitter.ExtractPages(1, 3, 5, 7)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use document...
func (s *Splitter) ExtractPages(pages ...int) (*document.Document, error) {
	// Validate.
	if len(pages) == 0 {
		return nil, fmt.Errorf("no pages specified")
	}

	// Validate page numbers.
	if err := s.validatePageNumbers(pages); err != nil {
		return nil, err
	}

	// Create output document.
	return s.createDocumentWithPages(pages)
}

// Close closes the splitter and releases resources.
//
// This should be called when done with the splitter (use defer).
func (s *Splitter) Close() error {
	if s.reader != nil {
		return s.reader.Close()
	}
	return nil
}

// extractAndWrite extracts a single page and writes it.
func (s *Splitter) extractAndWrite(outputPath string, pageNum int) error {
	return s.extractPages(outputPath, []int{pageNum})
}

// extractPages extracts pages and writes to file.
func (s *Splitter) extractPages(outputPath string, pages []int) error {
	// Create document with pages.
	doc, err := s.createDocumentWithPages(pages)
	if err != nil {
		return err
	}

	// Write document using merger's write logic.
	merger := &Merger{outputDoc: doc}
	return merger.writeOutput(outputPath)
}

// createDocumentWithPages creates a document with specified pages.
func (s *Splitter) createDocumentWithPages(pages []int) (*document.Document, error) {
	// Create output document.
	outputDoc := document.NewDocument()

	// Copy each page.
	for _, pageNum := range pages {
		if err := s.copyPage(outputDoc, pageNum); err != nil {
			return nil, fmt.Errorf("copy page %d: %w", pageNum, err)
		}
	}

	return outputDoc, nil
}

// copyPage copies a page to the output document.
func (s *Splitter) copyPage(outputDoc *document.Document, pageNum int) error {
	// Get source page (0-based).
	pages := s.sourceDoc.Pages()
	srcPage := pages[pageNum-1]

	// Get page size.
	mediaBox := srcPage.MediaBox()
	size := sizeFromMediaBox(mediaBox)

	// Add page to output.
	dstPage, err := outputDoc.AddPage(size)
	if err != nil {
		return fmt.Errorf("add page: %w", err)
	}

	// Copy rotation.
	if err := dstPage.SetRotation(srcPage.Rotation()); err != nil {
		return fmt.Errorf("set rotation: %w", err)
	}

	return nil
}

// validateRanges validates page ranges.
func (s *Splitter) validateRanges(ranges []PageRange) error {
	pageCount := s.sourceDoc.PageCount()

	for i, r := range ranges {
		// Check range bounds.
		if r.Start < 1 {
			return fmt.Errorf("range %d: start must be >= 1", i+1)
		}
		if r.End < r.Start {
			return fmt.Errorf("range %d: end must be >= start", i+1)
		}
		if r.End > pageCount {
			return fmt.Errorf("range %d: end %d exceeds page count %d",
				i+1, r.End, pageCount)
		}

		// Check output path.
		if r.Output == "" {
			return fmt.Errorf("range %d: output path is empty", i+1)
		}
	}

	return nil
}

// validatePageNumbers validates page numbers.
func (s *Splitter) validatePageNumbers(pages []int) error {
	pageCount := s.sourceDoc.PageCount()

	for _, pageNum := range pages {
		if pageNum < 1 || pageNum > pageCount {
			return fmt.Errorf("invalid page %d (file has %d pages)",
				pageNum, pageCount)
		}
	}

	return nil
}
