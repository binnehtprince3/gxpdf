package creator

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

const (
	// Test constants.
	testFormatPNG = "png"
)

// TestLoadImage tests loading a JPEG image from a file.
func TestLoadImage(t *testing.T) {
	// Create a temporary JPEG file.
	tmpFile := createTempJPEG(t, 100, 80, color.RGBA{255, 0, 0, 255})
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	// Load the image.
	img, err := LoadImage(tmpFile)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	// Verify dimensions.
	if img.Width() != 100 {
		t.Errorf("expected width 100, got %d", img.Width())
	}
	if img.Height() != 80 {
		t.Errorf("expected height 80, got %d", img.Height())
	}

	// Verify color space.
	if img.ColorSpace() != ColorSpaceRGB {
		t.Errorf("expected RGB color space, got %s", img.ColorSpace())
	}

	// Verify data is not empty.
	if len(img.Data()) == 0 {
		t.Error("image data is empty")
	}
}

// TestLoadImageFromReader tests loading from an io.Reader.
func TestLoadImageFromReader(t *testing.T) {
	// Create JPEG data in memory.
	data := createJPEGData(t, 150, 100, color.RGBA{0, 255, 0, 255})
	reader := bytes.NewReader(data)

	// Load the image.
	img, err := LoadImageFromReader(reader)
	if err != nil {
		t.Fatalf("LoadImageFromReader failed: %v", err)
	}

	// Verify dimensions.
	if img.Width() != 150 {
		t.Errorf("expected width 150, got %d", img.Width())
	}
	if img.Height() != 100 {
		t.Errorf("expected height 100, got %d", img.Height())
	}
}

// TestLoadImageUnsupportedFormat tests loading non-JPEG images.
func TestLoadImageUnsupportedFormat(t *testing.T) {
	// Try to load a non-JPEG (just some random bytes).
	reader := bytes.NewReader([]byte("not a jpeg"))

	_, err := LoadImageFromReader(reader)
	if err == nil {
		t.Error("expected error for non-JPEG data, got nil")
	}
}

// TestDrawImage tests the DrawImage method.
func TestDrawImage(t *testing.T) {
	// Create test image.
	data := createJPEGData(t, 100, 80, color.RGBA{255, 0, 0, 255})
	img, err := LoadImageFromReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to load test image: %v", err)
	}

	// Create test page and draw image.
	page := createTestPage(t)
	err = page.DrawImage(img, 100, 500, 200, 150)
	if err != nil {
		t.Errorf("DrawImage failed: %v", err)
	}

	// Verify operation.
	verifyImageOperation(t, page, img)
}

// Helper: createTestPage creates a page for testing.
func createTestPage(t *testing.T) *Page {
	t.Helper()
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}
	return page
}

// Helper: verifyImageOperation verifies image drawing operation.
func verifyImageOperation(t *testing.T, page *Page, img *Image) {
	t.Helper()
	ops := page.GraphicsOperations()
	if len(ops) != 1 {
		t.Fatalf("expected 1 graphics operation, got %d", len(ops))
	}

	op := ops[0]
	if op.Type != GraphicsOpImage {
		t.Errorf("expected image operation, got type %v", op.Type)
	}
	if op.X != 100 || op.Y != 500 {
		t.Errorf("expected position (100, 500), got (%.0f, %.0f)", op.X, op.Y)
	}
	if op.Width != 200 || op.Height != 150 {
		t.Errorf("expected size (200, 150), got (%.0f, %.0f)", op.Width, op.Height)
	}
	if op.Image != img {
		t.Error("image reference not preserved in operation")
	}
}

// TestDrawImageInvalidDimensions tests validation of image dimensions.
func TestDrawImageInvalidDimensions(t *testing.T) {
	// Create test image.
	data := createJPEGData(t, 100, 80, color.RGBA{255, 0, 0, 255})
	img, err := LoadImageFromReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to load test image: %v", err)
	}

	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	tests := []struct {
		name   string
		width  float64
		height float64
	}{
		{"zero width", 0, 100},
		{"zero height", 100, 0},
		{"negative width", -100, 100},
		{"negative height", 100, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := page.DrawImage(img, 100, 100, tt.width, tt.height)
			if err == nil {
				t.Error("expected error for invalid dimensions, got nil")
			}
		})
	}
}

// TestDrawImageFit tests the DrawImageFit method.
func TestDrawImageFit(t *testing.T) {
	// Create test image (200x100 aspect ratio 2:1).
	data := createJPEGData(t, 200, 100, color.RGBA{0, 0, 255, 255})
	img, err := LoadImageFromReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to load test image: %v", err)
	}

	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Fit image into 100x100 box.
	// Should scale to 100x50 (maintaining 2:1 aspect ratio).
	err = page.DrawImageFit(img, 100, 500, 100, 100)
	if err != nil {
		t.Errorf("DrawImageFit failed: %v", err)
	}

	// Verify operation.
	ops := page.GraphicsOperations()
	if len(ops) != 1 {
		t.Fatalf("expected 1 graphics operation, got %d", len(ops))
	}

	op := ops[0]
	if op.Type != GraphicsOpImage {
		t.Errorf("expected image operation, got type %v", op.Type)
	}

	// Image should be scaled to 100x50 and centered.
	expectedWidth := 100.0
	expectedHeight := 50.0
	if op.Width != expectedWidth || op.Height != expectedHeight {
		t.Errorf("expected size (%.0f, %.0f), got (%.0f, %.0f)",
			expectedWidth, expectedHeight, op.Width, op.Height)
	}

	// Image should be centered vertically (25 points offset).
	expectedY := 500.0 + (100.0-50.0)/2.0 // y + (maxHeight - scaledHeight) / 2
	if op.Y != expectedY {
		t.Errorf("expected Y position %.0f (centered), got %.0f", expectedY, op.Y)
	}
}

// TestCalculateFitDimensions tests aspect ratio calculations.
func TestCalculateFitDimensions(t *testing.T) {
	tests := []struct {
		name       string
		imgW, imgH float64
		maxW, maxH float64
		expectW    float64
		expectH    float64
	}{
		{
			name: "fit width constrained",
			imgW: 200, imgH: 100,
			maxW: 100, maxH: 100,
			expectW: 100, expectH: 50,
		},
		{
			name: "fit height constrained",
			imgW: 100, imgH: 200,
			maxW: 100, maxH: 100,
			expectW: 50, expectH: 100,
		},
		{
			name: "fit exact match",
			imgW: 100, imgH: 100,
			maxW: 100, maxH: 100,
			expectW: 100, expectH: 100,
		},
		{
			name: "fit smaller than max",
			imgW: 50, imgH: 50,
			maxW: 100, maxH: 100,
			expectW: 100, expectH: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := calculateFitDimensions(tt.imgW, tt.imgH, tt.maxW, tt.maxH)
			if w != tt.expectW || h != tt.expectH {
				t.Errorf("expected (%.0f, %.0f), got (%.0f, %.0f)",
					tt.expectW, tt.expectH, w, h)
			}
		})
	}
}

// Helper: createJPEGData creates a JPEG image in memory.
func createJPEGData(t *testing.T, width, height int, c color.Color) []byte {
	t.Helper()

	// Create image.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	// Encode to JPEG.
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	if err != nil {
		t.Fatalf("failed to encode JPEG: %v", err)
	}

	return buf.Bytes()
}

// Helper: createTempJPEG creates a temporary JPEG file.
func createTempJPEG(t *testing.T, width, height int, c color.Color) string {
	t.Helper()

	// Create temp file.
	tmpFile, err := os.CreateTemp("", "test-*.jpg")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		_ = tmpFile.Close()
	}()

	// Create and encode image.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	err = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: 90})
	if err != nil {
		t.Fatalf("failed to encode JPEG: %v", err)
	}

	return tmpFile.Name()
}

// TestLoadPNGImage tests loading a PNG image.
func TestLoadPNGImage(t *testing.T) {
	// Create a temporary PNG file (RGB).
	tmpFile := createTempPNG(t, 100, 80, color.RGBA{255, 0, 0, 255})
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	// Load the image.
	img, err := LoadImage(tmpFile)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	// Verify format.
	if img.Format() != testFormatPNG {
		t.Errorf("expected format png, got %s", img.Format())
	}

	// Verify dimensions.
	if img.Width() != 100 {
		t.Errorf("expected width 100, got %d", img.Width())
	}
	if img.Height() != 80 {
		t.Errorf("expected height 80, got %d", img.Height())
	}

	// Verify color space.
	if img.ColorSpace() != ColorSpaceRGB {
		t.Errorf("expected RGB color space, got %s", img.ColorSpace())
	}

	// Verify data is not empty.
	if len(img.Data()) == 0 {
		t.Error("image data is empty")
	}

	// Verify components.
	if img.Components() != 3 {
		t.Errorf("expected 3 components, got %d", img.Components())
	}
}

// TestLoadPNGWithAlpha tests loading an RGBA PNG with transparency.
func TestLoadPNGWithAlpha(t *testing.T) {
	// Create PNG with semi-transparent pixels.
	data := createPNGData(t, 50, 50, color.RGBA{0, 0, 255, 128})
	reader := bytes.NewReader(data)

	// Load the image.
	img, err := LoadImageFromReader(reader)
	if err != nil {
		t.Fatalf("LoadImageFromReader failed: %v", err)
	}

	// Verify format.
	if img.Format() != testFormatPNG {
		t.Errorf("expected format png, got %s", img.Format())
	}

	// Verify alpha mask exists.
	if !img.HasAlpha() {
		t.Error("expected image to have alpha mask")
	}

	if img.AlphaMask() == nil {
		t.Error("alpha mask is nil")
	}

	// Verify color space is still RGB.
	if img.ColorSpace() != ColorSpaceRGB {
		t.Errorf("expected RGB color space, got %s", img.ColorSpace())
	}
}

// TestLoadPNGGrayscale tests loading a grayscale PNG.
func TestLoadPNGGrayscale(t *testing.T) {
	// Create grayscale PNG.
	data := createGrayscalePNGData(t, 60, 40, 128)
	reader := bytes.NewReader(data)

	// Load the image.
	img, err := LoadImageFromReader(reader)
	if err != nil {
		t.Fatalf("LoadImageFromReader failed: %v", err)
	}

	// Verify format.
	if img.Format() != testFormatPNG {
		t.Errorf("expected format png, got %s", img.Format())
	}

	// Verify dimensions.
	if img.Width() != 60 || img.Height() != 40 {
		t.Errorf("expected dimensions (60, 40), got (%d, %d)", img.Width(), img.Height())
	}

	// Verify color space.
	if img.ColorSpace() != ColorSpaceGray {
		t.Errorf("expected grayscale color space, got %s", img.ColorSpace())
	}

	// Verify components.
	if img.Components() != 1 {
		t.Errorf("expected 1 component, got %d", img.Components())
	}

	// Verify no alpha.
	if img.HasAlpha() {
		t.Error("grayscale image should not have alpha mask")
	}
}

// TestLoadPNGPaletted tests loading a paletted PNG.
func TestLoadPNGPaletted(t *testing.T) {
	// Create paletted PNG.
	data := createPalettedPNGData(t, 70, 50)
	reader := bytes.NewReader(data)

	// Load the image.
	img, err := LoadImageFromReader(reader)
	if err != nil {
		t.Fatalf("LoadImageFromReader failed: %v", err)
	}

	// Verify format.
	if img.Format() != testFormatPNG {
		t.Errorf("expected format png, got %s", img.Format())
	}

	// Paletted PNG should be converted to RGB.
	if img.ColorSpace() != ColorSpaceRGB {
		t.Errorf("expected RGB color space (converted), got %s", img.ColorSpace())
	}

	// Verify components.
	if img.Components() != 3 {
		t.Errorf("expected 3 components (converted to RGB), got %d", img.Components())
	}
}

// TestDetectImageFormat tests format detection.
func TestDetectImageFormat(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		expect string
	}{
		{
			name:   "JPEG",
			data:   []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46},
			expect: "jpeg",
		},
		{
			name:   "PNG",
			data:   []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A},
			expect: "png",
		},
		{
			name:   "unknown",
			data:   []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
			expect: "",
		},
		{
			name:   "too short",
			data:   []byte{0xFF, 0xD8},
			expect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := detectImageFormat(tt.data)
			if format != tt.expect {
				t.Errorf("expected format %q, got %q", tt.expect, format)
			}
		})
	}
}

// TestPNGWithFullyOpaqueAlpha tests that fully opaque RGBA PNG has no mask.
func TestPNGWithFullyOpaqueAlpha(t *testing.T) {
	// Create PNG with alpha=255 (fully opaque).
	data := createPNGData(t, 30, 30, color.RGBA{255, 0, 0, 255})
	reader := bytes.NewReader(data)

	// Load the image.
	img, err := LoadImageFromReader(reader)
	if err != nil {
		t.Fatalf("LoadImageFromReader failed: %v", err)
	}

	// Verify no alpha mask (optimization).
	if img.HasAlpha() {
		t.Error("fully opaque image should not have alpha mask")
	}
}

// Helper: createPNGData creates a PNG image in memory.
func createPNGData(t *testing.T, width, height int, c color.Color) []byte {
	t.Helper()

	// Create image.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	// Encode to PNG.
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("failed to encode PNG: %v", err)
	}

	return buf.Bytes()
}

// Helper: createTempPNG creates a temporary PNG file.
func createTempPNG(t *testing.T, width, height int, c color.Color) string {
	t.Helper()

	// Create temp file.
	tmpFile, err := os.CreateTemp("", "test-*.png")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		_ = tmpFile.Close()
	}()

	// Create and encode image.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	err = png.Encode(tmpFile, img)
	if err != nil {
		t.Fatalf("failed to encode PNG: %v", err)
	}

	return tmpFile.Name()
}

// Helper: createGrayscalePNGData creates a grayscale PNG.
func createGrayscalePNGData(t *testing.T, width, height int, gray uint8) []byte {
	t.Helper()

	// Create grayscale image.
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{Y: gray})
		}
	}

	// Encode to PNG.
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("failed to encode grayscale PNG: %v", err)
	}

	return buf.Bytes()
}

// Helper: createPalettedPNGData creates a paletted PNG.
func createPalettedPNGData(t *testing.T, width, height int) []byte {
	t.Helper()

	// Create palette with 3 colors.
	palette := color.Palette{
		color.RGBA{255, 0, 0, 255}, // Red
		color.RGBA{0, 255, 0, 255}, // Green
		color.RGBA{0, 0, 255, 255}, // Blue
	}

	// Create paletted image.
	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			//nolint:gosec // G115: Safe conversion, (x+y)%3 is always in range [0, 2].
			img.SetColorIndex(x, y, uint8((x+y)%3))
		}
	}

	// Encode to PNG.
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("failed to encode paletted PNG: %v", err)
	}

	return buf.Bytes()
}
