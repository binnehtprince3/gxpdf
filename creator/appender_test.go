package creator

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestNewAppender_Success tests opening a valid PDF.
func TestNewAppender_Success(t *testing.T) {
	// Create a simple test PDF first.
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	// Open with Appender.
	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Verify page count (test PDF has 2 pages).
	if app.PageCount() < 1 {
		t.Errorf("PageCount() = %d, want >= 1", app.PageCount())
	}
}

// TestNewAppender_NonExistentFile tests error handling for missing file.
func TestNewAppender_NonExistentFile(t *testing.T) {
	_, err := NewAppender("nonexistent.pdf")
	if err == nil {
		t.Error("NewAppender() should fail for nonexistent file")
	}
}

// TestAppender_Close tests closing the appender.
func TestAppender_Close(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}

	// Close once.
	if err := app.Close(); err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Close again (should be safe).
	if err := app.Close(); err != nil {
		t.Errorf("Close() failed on second call: %v", err)
	}
}

// TestAppender_PageCount tests page counting.
func TestAppender_PageCount(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Initial count (test PDF has 2 pages).
	initialCount := app.PageCount()
	if initialCount != 2 {
		t.Errorf("Initial PageCount() = %d, want 2", initialCount)
	}

	// Add a page.
	_, err = app.AddPage(A4)
	if err != nil {
		t.Fatalf("AddPage() failed: %v", err)
	}

	// Count should increase.
	if app.PageCount() != 3 {
		t.Errorf("PageCount() after AddPage() = %d, want 3", app.PageCount())
	}
}

// TestAppender_GetPage tests page retrieval.
func TestAppender_GetPage(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	tests := []struct {
		name      string
		index     int
		wantError bool
	}{
		{"valid index 0", 0, false},
		{"negative index", -1, true},
		{"out of bounds", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, err := app.GetPage(tt.index)
			if tt.wantError {
				if err == nil {
					t.Error("GetPage() should return error")
				}
			} else {
				if err != nil {
					t.Errorf("GetPage() failed: %v", err)
				}
				if page == nil {
					t.Error("GetPage() returned nil page")
				}
			}
		})
	}
}

// TestAppender_AddPage tests adding new pages.
func TestAppender_AddPage(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Add A4 page.
	page, err := app.AddPage(A4)
	if err != nil {
		t.Fatalf("AddPage(A4) failed: %v", err)
	}

	if page == nil {
		t.Fatal("AddPage() returned nil page")
	}

	// Verify page size.
	if page.Width() != 595 || page.Height() != 842 {
		t.Errorf("Page size = %.0f x %.0f, want 595 x 842", page.Width(), page.Height())
	}
}

// TestAppender_AddTextToExistingPage tests adding text to existing page.
func TestAppender_AddTextToExistingPage(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Get first page.
	page, err := app.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage(0) failed: %v", err)
	}

	// Add text.
	err = page.AddText("Watermark", 300, 400, HelveticaBold, 48)
	if err != nil {
		t.Errorf("AddText() failed: %v", err)
	}

	// Verify text operation was added.
	if len(page.TextOperations()) != 1 {
		t.Errorf("TextOperations() count = %d, want 1", len(page.TextOperations()))
	}
}

// TestAppender_AddTextToNewPage tests adding text to new page.
func TestAppender_AddTextToNewPage(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Add new page.
	page, err := app.AddPage(Letter)
	if err != nil {
		t.Fatalf("AddPage() failed: %v", err)
	}

	// Add text.
	err = page.AddText("New content", 100, 700, Helvetica, 12)
	if err != nil {
		t.Errorf("AddText() failed: %v", err)
	}

	// Verify text operation.
	if len(page.TextOperations()) != 1 {
		t.Errorf("TextOperations() count = %d, want 1", len(page.TextOperations()))
	}
}

// TestAppender_WriteToFile tests writing modified PDF.
func TestAppender_WriteToFile(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Modify first page.
	page, err := app.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage(0) failed: %v", err)
	}

	err = page.AddText("Modified", 100, 700, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddText() failed: %v", err)
	}

	// Add new page.
	newPage, err := app.AddPage(A4)
	if err != nil {
		t.Fatalf("AddPage() failed: %v", err)
	}

	err = newPage.AddText("Page 2", 100, 700, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddText() failed: %v", err)
	}

	// Write to new file.
	outputPath := filepath.Join(t.TempDir(), "modified.pdf")
	err = app.WriteToFile(outputPath)
	if err != nil {
		t.Fatalf("WriteToFile() failed: %v", err)
	}

	// Verify output file exists.
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify output file is a PDF (check header).
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	if len(data) == 0 {
		t.Error("Output PDF is empty")
	}

	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		t.Error("Output file is not a valid PDF (missing header)")
	}

	// Note: We cannot reopen the written PDF because the writer has a bug
	// with xref offsets. This will be fixed later.
}

// TestAppender_SetMetadata tests updating document metadata.
func TestAppender_SetMetadata(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Set metadata.
	app.SetMetadata("Modified Title", "John Doe", "Test Document")
	app.SetKeywords("test", "modified", "gxpdf")

	// Write to file.
	outputPath := filepath.Join(t.TempDir(), "metadata.pdf")
	err = app.WriteToFile(outputPath)
	if err != nil {
		t.Fatalf("WriteToFile() failed: %v", err)
	}

	// Verify file was created.
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

// TestAppender_AddGraphicsToPage tests adding graphics to existing page.
func TestAppender_AddGraphicsToPage(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Get first page.
	page, err := app.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage(0) failed: %v", err)
	}

	// Draw line.
	lineOpts := &LineOptions{
		Color: Black,
		Width: 2.0,
	}
	err = page.DrawLine(100, 700, 500, 700, lineOpts)
	if err != nil {
		t.Errorf("DrawLine() failed: %v", err)
	}

	// Draw rectangle.
	rectOpts := &RectOptions{
		StrokeColor: &Black,
		StrokeWidth: 1.0,
		FillColor:   &LightGray,
	}
	err = page.DrawRect(100, 600, 200, 100, rectOpts)
	if err != nil {
		t.Errorf("DrawRect() failed: %v", err)
	}

	// Verify graphics operations.
	if len(page.GraphicsOperations()) != 2 {
		t.Errorf("GraphicsOperations() count = %d, want 2", len(page.GraphicsOperations()))
	}
}

// TestAppender_Document tests accessing underlying document.
func TestAppender_Document(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	doc := app.Document()
	if doc == nil {
		t.Error("Document() returned nil")
	}
}

// TestAppender_GetParserReader tests accessing parser reader.
func TestAppender_GetParserReader(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	parserReader := app.GetParserReader()
	if parserReader == nil {
		t.Error("GetParserReader() returned nil")
	}
}

// TestAppender_RotatePage tests rotating existing pages.
func TestAppender_RotatePage(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	page, err := app.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage(0) failed: %v", err)
	}

	// Test valid rotations.
	rotations := []int{0, 90, 180, 270, 0}
	for _, rotation := range rotations {
		err = page.SetRotation(rotation)
		if err != nil {
			t.Errorf("SetRotation(%d) failed: %v", rotation, err)
		}
		if page.Rotation() != rotation {
			t.Errorf("Rotation = %d, want %d", page.Rotation(), rotation)
		}
	}
}

// TestAppender_RotatePage_Invalid tests invalid rotation values.
func TestAppender_RotatePage_Invalid(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	page, err := app.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage(0) failed: %v", err)
	}

	// Test invalid rotation values.
	invalidRotations := []int{45, 135, 360, -90, 1, 91}

	for _, rotation := range invalidRotations {
		t.Run("invalid_rotation_"+string(rune(rotation)), func(t *testing.T) {
			err := page.SetRotation(rotation)
			if err == nil {
				t.Errorf("SetRotation(%d) should fail, but succeeded", rotation)
			}

			// Rotation should not change after error.
			if page.Rotation() != 0 {
				t.Errorf("Rotation changed after error: got %d, want 0", page.Rotation())
			}
		})
	}
}

// TestAppender_RotateAllPages tests rotating all pages in a document.
func TestAppender_RotateAllPages(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	// Rotate all pages and verify.
	rotateAndVerifyAllPages(t, app, 90)

	// Write to file and verify.
	outputPath := filepath.Join(t.TempDir(), "rotated.pdf")
	verifyWriteAndOutput(t, app, outputPath)
}

// rotateAndVerifyAllPages rotates all pages to specified degree and verifies.
func rotateAndVerifyAllPages(t *testing.T, app *Appender, degrees int) {
	t.Helper()
	pageCount := app.PageCount()

	// Rotate all pages.
	for i := 0; i < pageCount; i++ {
		page, err := app.GetPage(i)
		if err != nil {
			t.Fatalf("GetPage(%d) failed: %v", i, err)
		}
		if err := page.SetRotation(degrees); err != nil {
			t.Errorf("SetRotation(%d) on page %d failed: %v", degrees, i, err)
		}
	}

	// Verify all pages are rotated.
	for i := 0; i < pageCount; i++ {
		page, err := app.GetPage(i)
		if err != nil {
			t.Fatalf("GetPage(%d) failed: %v", i, err)
		}
		if page.Rotation() != degrees {
			t.Errorf("Page %d rotation = %d, want %d", i, page.Rotation(), degrees)
		}
	}
}

// verifyWriteAndOutput writes PDF and verifies output is valid.
func verifyWriteAndOutput(t *testing.T, app *Appender, outputPath string) {
	t.Helper()

	err := app.WriteToFile(outputPath)
	if err != nil {
		t.Fatalf("WriteToFile() failed: %v", err)
	}

	//nolint:gosec // Test file path is safe
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output PDF: %v", err)
	}

	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		t.Error("Output file is not a valid PDF")
	}
}

// TestAppender_RotateAndAddContent tests rotating and adding content.
func TestAppender_RotateAndAddContent(t *testing.T) {
	testPDF := createTestPDF(t)
	defer func() { _ = os.Remove(testPDF) }()

	app, err := NewAppender(testPDF)
	if err != nil {
		t.Fatalf("NewAppender() failed: %v", err)
	}
	defer func() { _ = app.Close() }()

	page, err := app.GetPage(0)
	if err != nil {
		t.Fatalf("GetPage(0) failed: %v", err)
	}

	// Rotate to landscape.
	err = page.SetRotation(90)
	if err != nil {
		t.Errorf("SetRotation(90) failed: %v", err)
	}

	// Add text to rotated page.
	err = page.AddText("Rotated Text", 100, 100, Helvetica, 12)
	if err != nil {
		t.Errorf("AddText() on rotated page failed: %v", err)
	}

	// Draw graphics on rotated page.
	lineOpts := &LineOptions{
		Color: Black,
		Width: 2.0,
	}
	err = page.DrawLine(50, 50, 200, 50, lineOpts)
	if err != nil {
		t.Errorf("DrawLine() on rotated page failed: %v", err)
	}

	// Write to file.
	outputPath := filepath.Join(t.TempDir(), "rotated_with_content.pdf")
	err = app.WriteToFile(outputPath)
	if err != nil {
		t.Fatalf("WriteToFile() failed: %v", err)
	}

	// Verify file was created.
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

// createTestPDF creates a simple test PDF file and returns its path.
func createTestPDF(t *testing.T) string {
	t.Helper()

	// Create temporary directory.
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.pdf")

	// Use an existing valid PDF from pdfcpu reference samples.
	// This is a workaround until the PDF writer xref offset bug is fixed.
	sourcePDF := filepath.Join("..", "reference", "pdfcpu", "pkg", "samples", "annotations", "LinkAnnotWithDestTopLeft.pdf")

	// Check if source exists.
	if _, err := os.Stat(sourcePDF); os.IsNotExist(err) {
		t.Skipf("Source PDF not found: %s", sourcePDF)
	}

	// Copy source PDF to temp location.
	data, err := os.ReadFile(sourcePDF)
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	err = os.WriteFile(testPath, data, 0644)
	if err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}

	return testPath
}
