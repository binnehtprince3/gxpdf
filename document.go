package gxpdf

import (
	"context"
	"fmt"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/coregx/gxpdf/internal/parser"
	"github.com/coregx/gxpdf/internal/tabledetect"
)

// Document represents an opened PDF document.
//
// Document provides methods for reading document properties and extracting content.
// It must be closed after use to release resources.
//
// Example:
//
//	doc, err := gxpdf.Open("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer doc.Close()
//
//	fmt.Printf("Pages: %d\n", doc.PageCount())
//	tables := doc.ExtractTables()
type Document struct {
	reader *parser.Reader
	ctx    context.Context
	path   string
}

// Close closes the document and releases resources.
//
// It is safe to call Close multiple times.
func (d *Document) Close() error {
	if d.reader != nil {
		return d.reader.Close()
	}
	return nil
}

// Path returns the file path of the document.
func (d *Document) Path() string {
	return d.path
}

// PageCount returns the total number of pages in the document.
func (d *Document) PageCount() int {
	count, err := d.reader.GetPageCount()
	if err != nil {
		return 0
	}
	return count
}

// Page returns the page at the given index (0-based).
//
// Returns nil if the index is out of bounds.
func (d *Document) Page(index int) *Page {
	if index < 0 || index >= d.PageCount() {
		return nil
	}
	return &Page{
		doc:   d,
		index: index,
	}
}

// Pages returns an iterator over all pages.
//
// Example:
//
//	for _, page := range doc.Pages() {
//	    text := page.ExtractText()
//	    fmt.Println(text)
//	}
func (d *Document) Pages() []*Page {
	count := d.PageCount()
	pages := make([]*Page, count)
	for i := 0; i < count; i++ {
		pages[i] = &Page{doc: d, index: i}
	}
	return pages
}

// ExtractTables extracts all tables from all pages.
//
// This is the simplest way to extract tables - uses automatic detection
// with the 4-Pass Hybrid algorithm for best accuracy.
//
// Example:
//
//	tables := doc.ExtractTables()
//	for _, t := range tables {
//	    fmt.Printf("Table on page %d: %d rows x %d cols\n",
//	        t.PageNumber(), t.RowCount(), t.ColumnCount())
//	}
func (d *Document) ExtractTables() []*Table {
	tables, _ := d.ExtractTablesWithOptions(nil)
	return tables
}

// ExtractTablesWithOptions extracts tables with custom options.
//
// Example:
//
//	opts := &gxpdf.ExtractionOptions{
//	    Method: gxpdf.MethodLattice,
//	    Pages:  []int{0, 1, 2},
//	}
//	tables, err := doc.ExtractTablesWithOptions(opts)
func (d *Document) ExtractTablesWithOptions(opts *ExtractionOptions) ([]*Table, error) {
	if opts == nil {
		opts = DefaultExtractionOptions()
	}

	// Determine pages to process
	pages := opts.Pages
	if len(pages) == 0 {
		count := d.PageCount()
		pages = make([]int, count)
		for i := 0; i < count; i++ {
			pages[i] = i
		}
	}

	// Create text extractor
	textExtractor := extractor.NewTextExtractor(d.reader)

	var allTables []*Table

	for _, pageIndex := range pages {
		// Check context cancellation
		select {
		case <-d.ctx.Done():
			return allTables, d.ctx.Err()
		default:
		}

		// Extract text elements
		textElements, err := textExtractor.ExtractFromPage(pageIndex)
		if err != nil {
			return nil, fmt.Errorf("gxpdf: failed to extract text from page %d: %w", pageIndex, err)
		}

		// Detect tables
		tableDetector := tabledetect.NewDefaultTableDetector()

		var detectedTables []*tabledetect.TableRegion
		var graphicsElements []*extractor.GraphicsElement

		switch opts.Method {
		case MethodLattice:
			detectedTables, err = tableDetector.DetectTablesLattice(textElements, graphicsElements)
		case MethodStream:
			detectedTables, err = tableDetector.DetectTablesStream(textElements)
		default:
			detectedTables, err = tableDetector.DetectTables(textElements, graphicsElements)
		}

		if err != nil {
			return nil, fmt.Errorf("gxpdf: failed to detect tables on page %d: %w", pageIndex, err)
		}

		// Extract table data
		tableExtractor := tabledetect.NewTableExtractor(textElements)
		for _, region := range detectedTables {
			extracted, err := tableExtractor.ExtractTable(region)
			if err != nil {
				continue
			}
			extracted.PageNum = pageIndex

			allTables = append(allTables, &Table{internal: extracted})
		}
	}

	return allTables, nil
}

// GetImages extracts all images from all pages in the document.
//
// This is the simplest way to extract images - returns all images found
// across all pages.
//
// Example:
//
//	images := doc.GetImages()
//	for i, img := range images {
//	    fmt.Printf("Image %d: %dx%d, %s\n", i, img.Width(), img.Height(), img.ColorSpace())
//	    img.SaveToFile(fmt.Sprintf("image_%d.jpg", i))
//	}
func (d *Document) GetImages() []*Image {
	images, _ := d.GetImagesWithError()
	return images
}

// GetImagesWithError extracts all images from all pages, returning any errors.
//
// Use this when you need error handling for image extraction.
func (d *Document) GetImagesWithError() ([]*Image, error) {
	imageExtractor := extractor.NewImageExtractor(d.reader)
	internalImages, err := imageExtractor.ExtractFromDocument()
	if err != nil {
		return nil, fmt.Errorf("gxpdf: failed to extract images: %w", err)
	}

	// Wrap internal images in public API
	images := make([]*Image, len(internalImages))
	for i, internal := range internalImages {
		images[i] = &Image{internal: internal}
	}

	return images, nil
}

// Info returns document metadata.
func (d *Document) Info() *DocumentInfo {
	pinfo := d.reader.GetDocumentInfo()
	return &DocumentInfo{
		PageCount: d.PageCount(),
		Path:      d.path,
		Version:   pinfo.Version,
		Title:     pinfo.Title,
		Author:    pinfo.Author,
		Subject:   pinfo.Subject,
		Keywords:  pinfo.Keywords,
		Creator:   pinfo.Creator,
		Producer:  pinfo.Producer,
		Encrypted: pinfo.Encrypted,
	}
}

// Version returns the PDF version (e.g., "1.7").
func (d *Document) Version() string {
	return d.reader.GetDocumentInfo().Version
}

// Title returns the document title.
func (d *Document) Title() string {
	return d.reader.GetDocumentInfo().Title
}

// Author returns the document author.
func (d *Document) Author() string {
	return d.reader.GetDocumentInfo().Author
}

// Subject returns the document subject.
func (d *Document) Subject() string {
	return d.reader.GetDocumentInfo().Subject
}

// Keywords returns the document keywords.
func (d *Document) Keywords() string {
	return d.reader.GetDocumentInfo().Keywords
}

// Creator returns the application that created the document.
func (d *Document) Creator() string {
	return d.reader.GetDocumentInfo().Creator
}

// Producer returns the PDF producer.
func (d *Document) Producer() string {
	return d.reader.GetDocumentInfo().Producer
}

// IsEncrypted returns true if the document is encrypted.
func (d *Document) IsEncrypted() bool {
	return d.reader.GetDocumentInfo().Encrypted
}

// ExtractTextFromPage extracts text from a specific page (1-based).
func (d *Document) ExtractTextFromPage(pageNum int) (string, error) {
	if pageNum < 1 || pageNum > d.PageCount() {
		return "", fmt.Errorf("page %d out of range (1-%d)", pageNum, d.PageCount())
	}
	page := d.Page(pageNum - 1)
	if page == nil {
		return "", fmt.Errorf("page %d not found", pageNum)
	}
	return page.ExtractText(), nil
}

// ExtractTablesFromPage extracts tables from a specific page (1-based).
func (d *Document) ExtractTablesFromPage(pageNum int) []*Table {
	if pageNum < 1 || pageNum > d.PageCount() {
		return nil
	}
	page := d.Page(pageNum - 1)
	if page == nil {
		return nil
	}
	return page.ExtractTables()
}

// DocumentInfo contains metadata about a PDF document.
type DocumentInfo struct {
	PageCount int
	Path      string
	Version   string
	Title     string
	Author    string
	Subject   string
	Keywords  string
	Creator   string
	Producer  string
	Encrypted bool
}
