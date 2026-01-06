package creator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextAnnotation(t *testing.T) {
	tests := []struct {
		name        string
		x, y        float64
		contents    string
		setupAnnot  func(*TextAnnotation)
		expectError bool
		validatePDF func(*testing.T, string)
	}{
		{
			name:     "basic sticky note",
			x:        100,
			y:        700,
			contents: "This is a comment",
			setupAnnot: func(a *TextAnnotation) {
				a.SetAuthor("John Doe")
			},
			expectError: false,
		},
		{
			name:     "yellow sticky note",
			x:        150,
			y:        650,
			contents: "Important note",
			setupAnnot: func(a *TextAnnotation) {
				a.SetColor(Yellow).SetAuthor("Alice")
			},
			expectError: false,
		},
		{
			name:     "red sticky note (important)",
			x:        200,
			y:        600,
			contents: "URGENT!",
			setupAnnot: func(a *TextAnnotation) {
				a.SetColor(Red).SetAuthor("Manager").SetOpen(true)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, err := c.NewPage()
			require.NoError(t, err)

			// Add some text to the page.
			err = page.AddText("Sample text with annotation", 100, 750, Helvetica, 12)
			require.NoError(t, err)

			// Create and configure annotation.
			annot := NewTextAnnotation(tt.x, tt.y, tt.contents)
			if tt.setupAnnot != nil {
				tt.setupAnnot(annot)
			}

			// Add annotation to page.
			err = page.AddTextAnnotation(annot)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Write to file for manual verification (optional).
			// Uncomment to generate test PDFs.
			//tmpfile := filepath.Join(os.TempDir(), "test_text_annotation_"+tt.name+".pdf")
			//err = c.WriteToFile(tmpfile)
			//require.NoError(t, err)
			//t.Logf("PDF written to: %s", tmpfile)
		})
	}
}

func TestHighlightAnnotation(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Add text.
	err = page.AddText("This text should be highlighted", 100, 700, Helvetica, 12)
	require.NoError(t, err)

	// Highlight annotation.
	highlight := NewHighlightAnnotation(100, 695, 300, 710)
	highlight.SetColor(Yellow).SetAuthor("Bob").SetNote("Important text")

	err = page.AddHighlightAnnotation(highlight)
	require.NoError(t, err)

	// Verify page has annotation.
	assert.Equal(t, 1, page.page.AnnotationCount())
}

func TestUnderlineAnnotation(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Add text.
	err = page.AddText("This text should be underlined", 100, 650, Helvetica, 12)
	require.NoError(t, err)

	// Underline annotation.
	underline := NewUnderlineAnnotation(100, 645, 300, 660)
	underline.SetColor(Blue)

	err = page.AddUnderlineAnnotation(underline)
	require.NoError(t, err)

	// Verify page has annotation.
	assert.Equal(t, 1, page.page.AnnotationCount())
}

func TestStrikeOutAnnotation(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Add text.
	err = page.AddText("This text should be struck out", 100, 600, Helvetica, 12)
	require.NoError(t, err)

	// StrikeOut annotation.
	strikeout := NewStrikeOutAnnotation(100, 595, 300, 610)
	strikeout.SetColor(Red).SetNote("Obsolete")

	err = page.AddStrikeOutAnnotation(strikeout)
	require.NoError(t, err)

	// Verify page has annotation.
	assert.Equal(t, 1, page.page.AnnotationCount())
}

func TestStampAnnotation(t *testing.T) {
	tests := []struct {
		name       string
		stampName  string
		color      Color
		expectPass bool
	}{
		{
			name:       "Approved stamp",
			stampName:  string(StampApproved),
			color:      Green,
			expectPass: true,
		},
		{
			name:       "Draft stamp",
			stampName:  string(StampDraft),
			color:      Yellow,
			expectPass: true,
		},
		{
			name:       "Confidential stamp",
			stampName:  string(StampConfidential),
			color:      Red,
			expectPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, err := c.NewPage()
			require.NoError(t, err)

			// Create stamp.
			stamp := NewStampAnnotation(300, 700, 100, 50, StampApproved)
			stamp.SetColor(tt.color).SetAuthor("Manager")

			err = page.AddStampAnnotation(stamp)
			if tt.expectPass {
				require.NoError(t, err)
				assert.Equal(t, 1, page.page.AnnotationCount())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestMultipleAnnotationsOnPage(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Add text.
	err = page.AddText("Document with multiple annotations", 100, 750, HelveticaBold, 14)
	require.NoError(t, err)

	// Text annotation.
	note := NewTextAnnotation(100, 720, "Review this")
	note.SetAuthor("Reviewer").SetColor(Yellow)
	err = page.AddTextAnnotation(note)
	require.NoError(t, err)

	// Highlight.
	highlight := NewHighlightAnnotation(100, 690, 300, 705)
	highlight.SetColor(Yellow)
	err = page.AddHighlightAnnotation(highlight)
	require.NoError(t, err)

	// Stamp.
	stamp := NewStampAnnotation(400, 720, 80, 40, StampApproved)
	stamp.SetColor(Green)
	err = page.AddStampAnnotation(stamp)
	require.NoError(t, err)

	// Verify all annotations were added.
	assert.Equal(t, 3, page.page.AnnotationCount())

	// Write to file for manual verification.
	tmpfile := os.TempDir() + "/test_multiple_annotations.pdf"
	err = c.WriteToFile(tmpfile)
	require.NoError(t, err)
	t.Logf("PDF with multiple annotations written to: %s", tmpfile)

	// Clean up.
	defer func() {
		_ = os.Remove(tmpfile)
	}()
}

func TestAnnotationChaining(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Test method chaining.
	note := NewTextAnnotation(100, 700, "Chained annotation")
	note.SetAuthor("Alice").SetColor(Yellow).SetOpen(true)

	err = page.AddTextAnnotation(note)
	require.NoError(t, err)

	assert.Equal(t, 1, page.page.AnnotationCount())
}
