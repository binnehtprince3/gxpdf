package creator

import (
	"testing"

	"github.com/coregx/gxpdf/internal/document"
)

// TestDefaultLinkStyle tests the default link style.
func TestDefaultLinkStyle(t *testing.T) {
	style := DefaultLinkStyle()

	// Check defaults.
	if style.Font != Helvetica {
		t.Errorf("expected font Helvetica, got %s", style.Font)
	}
	if style.Size != 12 {
		t.Errorf("expected size 12, got %.2f", style.Size)
	}
	if style.Color != Blue {
		t.Errorf("expected color Blue, got %+v", style.Color)
	}
	if !style.Underline {
		t.Error("expected underline true, got false")
	}
}

// TestPage_AddLink tests adding a simple link.
func TestPage_AddLink(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Add link.
	err = page.AddLink("Visit Google", "https://google.com", 100, 700, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddLink failed: %v", err)
	}

	// Check that text was added.
	textOps := page.TextOperations()
	if len(textOps) == 0 {
		t.Fatal("expected text operations, got none")
	}
	if textOps[0].Text != "Visit Google" {
		t.Errorf("expected text 'Visit Google', got '%s'", textOps[0].Text)
	}

	// Check that annotation was added to domain page.
	annotations := page.page.Annotations()
	if len(annotations) == 0 {
		t.Fatal("expected annotations, got none")
	}
	if annotations[0].URI != "https://google.com" {
		t.Errorf("expected URI 'https://google.com', got '%s'", annotations[0].URI)
	}
	if annotations[0].IsInternal {
		t.Error("expected external link, got internal")
	}
}

// TestPage_AddLinkStyled tests adding a styled link.
func TestPage_AddLinkStyled(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Custom style.
	style := LinkStyle{
		Font:      HelveticaBold,
		Size:      14,
		Color:     Red,
		Underline: false,
	}

	// Add styled link.
	err = page.AddLinkStyled("Click here", "https://example.com", 100, 650, style)
	if err != nil {
		t.Fatalf("AddLinkStyled failed: %v", err)
	}

	// Check text color.
	textOps := page.TextOperations()
	if len(textOps) == 0 {
		t.Fatal("expected text operations, got none")
	}
	if textOps[0].Color != Red {
		t.Errorf("expected color Red, got %+v", textOps[0].Color)
	}
	if textOps[0].Font != HelveticaBold {
		t.Errorf("expected font HelveticaBold, got %s", textOps[0].Font)
	}

	// Check graphics ops (should be empty since no underline).
	graphicsOps := page.GraphicsOperations()
	if len(graphicsOps) != 0 {
		t.Errorf("expected no graphics ops (no underline), got %d", len(graphicsOps))
	}
}

// TestPage_AddInternalLink tests adding an internal page link.
func TestPage_AddInternalLink(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Add internal link to page 3 (0-based = 2).
	err = page.AddInternalLink("See page 3", 2, 100, 600, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddInternalLink failed: %v", err)
	}

	// Check annotation.
	annotations := page.page.Annotations()
	if len(annotations) == 0 {
		t.Fatal("expected annotations, got none")
	}
	if !annotations[0].IsInternal {
		t.Error("expected internal link, got external")
	}
	if annotations[0].DestPage != 2 {
		t.Errorf("expected dest page 2, got %d", annotations[0].DestPage)
	}
}

// TestLinkAnnotation_Rect tests that annotation rect is calculated correctly.
func TestLinkAnnotation_Rect(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Add link.
	err = page.AddLink("Test", "https://test.com", 100, 700, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddLink failed: %v", err)
	}

	// Get annotation.
	annotations := page.page.Annotations()
	if len(annotations) == 0 {
		t.Fatal("expected annotations, got none")
	}

	rect := annotations[0].Rect
	// Check that rect has reasonable values.
	if rect[0] != 100 { // x1 = x
		t.Errorf("expected x1=100, got %.2f", rect[0])
	}
	if rect[2] <= 100 { // x2 > x1
		t.Errorf("expected x2 > 100, got %.2f", rect[2])
	}
	if rect[1] >= 700 { // y1 < y (below baseline)
		t.Errorf("expected y1 < 700, got %.2f", rect[1])
	}
	if rect[3] <= 700 { // y2 > y (above baseline)
		t.Errorf("expected y2 > 700, got %.2f", rect[3])
	}
}

// TestLinkAnnotation_URI tests that URI is stored correctly.
func TestLinkAnnotation_URI(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	testURL := "https://example.com/path?query=value"
	err = page.AddLink("Link", testURL, 100, 700, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddLink failed: %v", err)
	}

	annotations := page.page.Annotations()
	if annotations[0].URI != testURL {
		t.Errorf("expected URI '%s', got '%s'", testURL, annotations[0].URI)
	}
}

// TestMultipleLinksOnPage tests adding multiple links to one page.
func TestMultipleLinksOnPage(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Add three links.
	err = page.AddLink("Link 1", "https://link1.com", 100, 700, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddLink 1 failed: %v", err)
	}

	err = page.AddLink("Link 2", "https://link2.com", 100, 650, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddLink 2 failed: %v", err)
	}

	err = page.AddInternalLink("Link 3", 1, 100, 600, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddInternalLink failed: %v", err)
	}

	// Check that all three annotations exist.
	annotations := page.page.Annotations()
	if len(annotations) != 3 {
		t.Fatalf("expected 3 annotations, got %d", len(annotations))
	}

	// Check first annotation.
	if annotations[0].URI != "https://link1.com" {
		t.Errorf("expected first URI 'https://link1.com', got '%s'", annotations[0].URI)
	}

	// Check second annotation.
	if annotations[1].URI != "https://link2.com" {
		t.Errorf("expected second URI 'https://link2.com', got '%s'", annotations[1].URI)
	}

	// Check third annotation (internal).
	if !annotations[2].IsInternal {
		t.Error("expected third annotation to be internal")
	}
	if annotations[2].DestPage != 1 {
		t.Errorf("expected third dest page 1, got %d", annotations[2].DestPage)
	}
}

// TestLinkUnderline tests that underline is drawn when requested.
func TestLinkUnderline(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Default style has underline.
	err = page.AddLink("Underlined", "https://example.com", 100, 700, Helvetica, 12)
	if err != nil {
		t.Fatalf("AddLink failed: %v", err)
	}

	// Check that a line was drawn (graphics operation).
	graphicsOps := page.GraphicsOperations()
	if len(graphicsOps) == 0 {
		t.Fatal("expected graphics ops for underline, got none")
	}

	// Check that it's a line operation.
	if graphicsOps[0].Type != GraphicsOpLine {
		t.Errorf("expected line operation, got type %d", graphicsOps[0].Type)
	}
}

// TestLink_EmptyText tests error handling for empty text.
func TestLink_EmptyText(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Empty text should error.
	err = page.AddLink("", "https://example.com", 100, 700, Helvetica, 12)
	if err == nil {
		t.Fatal("expected error for empty text, got nil")
	}
}

// TestLink_EmptyURL tests error handling for empty URL.
func TestLink_EmptyURL(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Empty URL should error.
	err = page.AddLink("Text", "", 100, 700, Helvetica, 12)
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

// TestInternalLink_NegativePage tests error handling for negative page number.
func TestInternalLink_NegativePage(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Negative page should error.
	err = page.AddInternalLink("Text", -1, 100, 700, Helvetica, 12)
	if err == nil {
		t.Fatal("expected error for negative page, got nil")
	}
}

// TestLink_ZeroFontSize tests error handling for zero font size.
func TestLink_ZeroFontSize(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Zero font size should error.
	err = page.AddLink("Text", "https://example.com", 100, 700, Helvetica, 0)
	if err == nil {
		t.Fatal("expected error for zero font size, got nil")
	}
}

// TestLinkAnnotation_Validate tests the validation of LinkAnnotation.
func TestLinkAnnotation_Validate(t *testing.T) {
	tests := []struct {
		name        string
		annotation  *document.LinkAnnotation
		expectError bool
	}{
		{
			name: "valid external link",
			annotation: document.NewLinkAnnotation(
				[4]float64{100, 690, 200, 710},
				"https://example.com",
			),
			expectError: false,
		},
		{
			name: "valid internal link",
			annotation: document.NewInternalLinkAnnotation(
				[4]float64{100, 690, 200, 710},
				2,
			),
			expectError: false,
		},
		{
			name: "invalid rect (x1 >= x2)",
			annotation: &document.LinkAnnotation{
				Rect:       [4]float64{200, 690, 100, 710},
				URI:        "https://example.com",
				IsInternal: false,
			},
			expectError: true,
		},
		{
			name: "invalid rect (y1 >= y2)",
			annotation: &document.LinkAnnotation{
				Rect:       [4]float64{100, 710, 200, 690},
				URI:        "https://example.com",
				IsInternal: false,
			},
			expectError: true,
		},
		{
			name: "external link with empty URI",
			annotation: &document.LinkAnnotation{
				Rect:       [4]float64{100, 690, 200, 710},
				URI:        "",
				IsInternal: false,
			},
			expectError: true,
		},
		{
			name: "internal link with negative page",
			annotation: &document.LinkAnnotation{
				Rect:       [4]float64{100, 690, 200, 710},
				DestPage:   -1,
				IsInternal: true,
			},
			expectError: true,
		},
		{
			name: "negative border width",
			annotation: &document.LinkAnnotation{
				Rect:        [4]float64{100, 690, 200, 710},
				URI:         "https://example.com",
				IsInternal:  false,
				BorderWidth: -1,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.annotation.Validate()
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
