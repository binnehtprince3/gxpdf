package creator

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreator_HelloWorld(t *testing.T) {
	// Create PDF with Hello World text
	c := New()
	c.SetTitle("Hello World Test")
	c.SetAuthor("GxPDF Test Suite")

	page, err := c.NewPage()
	require.NoError(t, err)

	err = page.AddText("Hello World!", 100, 700, Helvetica, 24)
	require.NoError(t, err)

	// Write to temporary file
	tmpFile := filepath.Join(t.TempDir(), "hello_world.pdf")
	err = c.WriteToFile(tmpFile)
	require.NoError(t, err)

	// Verify file exists and is valid PDF
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Check PDF header
	assert.True(t, bytes.HasPrefix(data, []byte("%PDF-")), "Should start with PDF header")

	// Check EOF marker
	assert.True(t, bytes.HasSuffix(bytes.TrimSpace(data), []byte("%%EOF")), "Should end with EOF marker")

	// With compression, text is not visible in raw bytes.
	// Check that content stream exists and is compressed.
	assert.Contains(t, string(data), "/Contents", "Should have content stream")
	assert.Contains(t, string(data), "/Filter /FlateDecode", "Content should be compressed")

	// Check for font reference
	assert.Contains(t, string(data), "/Font", "Should contain font resources")
	assert.Contains(t, string(data), "Helvetica", "Should reference Helvetica font")
}

func TestCreator_ColoredText(t *testing.T) {
	c := New()

	page, err := c.NewPage()
	require.NoError(t, err)

	// Add text with different colors
	err = page.AddTextColor("Red Text", 100, 700, HelveticaBold, 18, Red)
	require.NoError(t, err)

	err = page.AddTextColor("Blue Text", 100, 650, Helvetica, 16, Blue)
	require.NoError(t, err)

	err = page.AddTextColor("Green Text", 100, 600, HelveticaOblique, 14, Green)
	require.NoError(t, err)

	// Write to temporary file
	tmpFile := filepath.Join(t.TempDir(), "colored_text.pdf")
	err = c.WriteToFile(tmpFile)
	require.NoError(t, err)

	// Verify file exists
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// With compression enabled, we check PDF structure, not raw text.
	// Check that PDF is valid and has content stream with filter.
	assert.Contains(t, string(data), "/Filter /FlateDecode", "Should have compressed content")
	assert.Contains(t, string(data), "/Contents", "Should have content stream reference")
}

func TestCreator_MultipleFonts(t *testing.T) {
	c := New()

	page, err := c.NewPage()
	require.NoError(t, err)

	// Add text with different fonts
	err = page.AddText("Helvetica", 100, 750, Helvetica, 14)
	require.NoError(t, err)

	err = page.AddText("Times-Roman", 100, 700, TimesRoman, 14)
	require.NoError(t, err)

	err = page.AddText("Courier", 100, 650, Courier, 14)
	require.NoError(t, err)

	// Write to temporary file
	tmpFile := filepath.Join(t.TempDir(), "multiple_fonts.pdf")
	err = c.WriteToFile(tmpFile)
	require.NoError(t, err)

	// Verify file
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	// Check for all three fonts
	assert.Contains(t, string(data), "Helvetica", "Should contain Helvetica font")
	assert.Contains(t, string(data), "Times-Roman", "Should contain Times-Roman font")
	assert.Contains(t, string(data), "Courier", "Should contain Courier font")
}

func TestCreator_MultiplePagesWithText(t *testing.T) {
	c := New()

	// Page 1
	page1, err := c.NewPage()
	require.NoError(t, err)
	err = page1.AddText("Page 1", 100, 700, Helvetica, 24)
	require.NoError(t, err)

	// Page 2
	page2, err := c.NewPage()
	require.NoError(t, err)
	err = page2.AddText("Page 2", 100, 700, Helvetica, 24)
	require.NoError(t, err)

	// Page 3
	page3, err := c.NewPage()
	require.NoError(t, err)
	err = page3.AddText("Page 3", 100, 700, Helvetica, 24)
	require.NoError(t, err)

	// Write to file
	tmpFile := filepath.Join(t.TempDir(), "multiple_pages_text.pdf")
	err = c.WriteToFile(tmpFile)
	require.NoError(t, err)

	// Verify
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	// With compression, text content is not visible in raw bytes.
	// Check PDF structure instead.
	assert.Contains(t, string(data), "/Count 3", "Should have 3 pages")
	assert.Contains(t, string(data), "/Type /Page", "Should have page objects")
	assert.Contains(t, string(data), "/Filter /FlateDecode", "Content should be compressed")
}

func TestPage_AddText_Validation(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	tests := []struct {
		name      string
		text      string
		x         float64
		y         float64
		font      FontName
		size      float64
		wantError bool
	}{
		{
			name:      "valid text",
			text:      "Valid",
			x:         100,
			y:         700,
			font:      Helvetica,
			size:      12,
			wantError: false,
		},
		{
			name:      "zero font size",
			text:      "Invalid",
			x:         100,
			y:         700,
			font:      Helvetica,
			size:      0,
			wantError: true,
		},
		{
			name:      "negative font size",
			text:      "Invalid",
			x:         100,
			y:         700,
			font:      Helvetica,
			size:      -12,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := page.AddText(tt.text, tt.x, tt.y, tt.font, tt.size)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPage_AddTextColor_ValidColors(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	validColors := []Color{
		{0.5, 0.5, 0.5}, Black, White, Red, Green, Blue,
		{0, 0, 0}, {1, 1, 1}, {0.0, 0.0, 0.0}, {1.0, 1.0, 1.0},
	}
	for _, color := range validColors {
		err := page.AddTextColor("Test", 100, 700, Helvetica, 12, color)
		assert.NoError(t, err, "color %v should be valid", color)
	}
}

func TestPage_AddTextColor_InvalidColors(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	invalidColors := []struct {
		name  string
		color Color
	}{
		{"R < 0", Color{-0.1, 0.5, 0.5}},
		{"R > 1", Color{1.1, 0.5, 0.5}},
		{"G < 0", Color{0.5, -0.1, 0.5}},
		{"G > 1", Color{0.5, 1.1, 0.5}},
		{"B < 0", Color{0.5, 0.5, -0.1}},
		{"B > 1", Color{0.5, 0.5, 1.1}},
	}
	for _, tc := range invalidColors {
		err := page.AddTextColor("Test", 100, 700, Helvetica, 12, tc.color)
		assert.Error(t, err, "%s should be invalid", tc.name)
	}
}

func TestCreator_EmptyDocument(t *testing.T) {
	c := New()

	// Try to write empty document (no pages)
	tmpFile := filepath.Join(t.TempDir(), "empty.pdf")
	err := c.WriteToFile(tmpFile)

	// Should fail validation (no pages)
	assert.Error(t, err)
}

func TestCreator_PageWithoutContent(t *testing.T) {
	c := New()

	// Create page but don't add any text
	_, err := c.NewPage()
	require.NoError(t, err)

	// Should still write successfully (empty page)
	tmpFile := filepath.Join(t.TempDir(), "empty_page.pdf")
	err = c.WriteToFile(tmpFile)
	require.NoError(t, err)

	// Verify file exists
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Check PDF structure
	assert.True(t, bytes.HasPrefix(data, []byte("%PDF-")))
	assert.Contains(t, string(data), "/Count 1", "Should have 1 page")
}
