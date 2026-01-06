package encoding

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

// createTestJPEG creates a test JPEG image with solid color.
func createTestJPEG(width, height int, c color.Color, quality int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	return buf.Bytes()
}

// createTestGrayJPEG creates a test grayscale JPEG image.
func createTestGrayJPEG(width, height int, value uint8, quality int) []byte {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{Y: value})
		}
	}

	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	return buf.Bytes()
}

func TestDCTDecoder_Decode_RGB(t *testing.T) {
	decoder := NewDCTDecoder()

	// Create test JPEG (red).
	jpegData := createTestJPEG(100, 100, color.RGBA{255, 0, 0, 255}, 90)

	// Decode.
	result, err := decoder.DecodeWithMetadata(jpegData)
	if err != nil {
		t.Fatalf("Failed to decode JPEG: %v", err)
	}

	// Verify metadata.
	if result.Width != 100 {
		t.Errorf("Expected width 100, got %d", result.Width)
	}
	if result.Height != 100 {
		t.Errorf("Expected height 100, got %d", result.Height)
	}
	if result.Components != 3 {
		t.Errorf("Expected 3 components, got %d", result.Components)
	}
	if result.BitsPerComponent != 8 {
		t.Errorf("Expected 8 bits per component, got %d", result.BitsPerComponent)
	}

	// Verify data length.
	expectedLen := 100 * 100 * 3
	if len(result.Data) != expectedLen {
		t.Errorf("Expected data length %d, got %d", expectedLen, len(result.Data))
	}

	// Check center pixel is approximately red (JPEG is lossy).
	centerIdx := (50*100 + 50) * 3
	r, g, b := result.Data[centerIdx], result.Data[centerIdx+1], result.Data[centerIdx+2]
	if r < 200 || g > 50 || b > 50 {
		t.Errorf("Expected red-ish pixel, got RGB(%d, %d, %d)", r, g, b)
	}
}

func TestDCTDecoder_Decode_Grayscale(t *testing.T) {
	decoder := NewDCTDecoder()

	// Create test grayscale JPEG.
	jpegData := createTestGrayJPEG(50, 50, 128, 90)

	// Decode.
	result, err := decoder.DecodeWithMetadata(jpegData)
	if err != nil {
		t.Fatalf("Failed to decode grayscale JPEG: %v", err)
	}

	// Verify metadata.
	if result.Width != 50 {
		t.Errorf("Expected width 50, got %d", result.Width)
	}
	if result.Height != 50 {
		t.Errorf("Expected height 50, got %d", result.Height)
	}
	if result.Components != 1 {
		t.Errorf("Expected 1 component for grayscale, got %d", result.Components)
	}

	// Verify data length.
	expectedLen := 50 * 50
	if len(result.Data) != expectedLen {
		t.Errorf("Expected data length %d, got %d", expectedLen, len(result.Data))
	}

	// Check center pixel is approximately mid-gray.
	centerIdx := 25*50 + 25
	gray := result.Data[centerIdx]
	if gray < 100 || gray > 160 {
		t.Errorf("Expected gray ~128, got %d", gray)
	}
}

func TestDCTDecoder_Decode_InvalidData(t *testing.T) {
	decoder := NewDCTDecoder()

	// Try to decode invalid data.
	_, err := decoder.Decode([]byte("not a jpeg"))
	if err == nil {
		t.Error("Expected error for invalid JPEG data")
	}
}

func TestDCTDecoder_DecodeToImage(t *testing.T) {
	decoder := NewDCTDecoder()

	// Create test JPEG.
	jpegData := createTestJPEG(30, 30, color.RGBA{0, 255, 0, 255}, 90)

	// Decode to image.
	img, err := decoder.DecodeToImage(jpegData)
	if err != nil {
		t.Fatalf("Failed to decode to image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 30 || bounds.Dy() != 30 {
		t.Errorf("Expected 30x30 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestDCTDecoder_Encode_RGB(t *testing.T) {
	decoder := NewDCTDecoder()

	// Create raw RGB data (blue).
	width, height := 10, 10
	data := make([]byte, width*height*3)
	for i := 0; i < len(data); i += 3 {
		data[i] = 0     // R
		data[i+1] = 0   // G
		data[i+2] = 255 // B
	}

	// Encode.
	jpegData, err := decoder.Encode(data, width, height, 90)
	if err != nil {
		t.Fatalf("Failed to encode JPEG: %v", err)
	}

	// Verify we got valid JPEG data.
	if len(jpegData) < 100 {
		t.Error("JPEG data seems too small")
	}

	// Decode and verify.
	result, err := decoder.DecodeWithMetadata(jpegData)
	if err != nil {
		t.Fatalf("Failed to decode encoded JPEG: %v", err)
	}

	if result.Width != 10 || result.Height != 10 {
		t.Errorf("Expected 10x10, got %dx%d", result.Width, result.Height)
	}

	// Check center pixel is approximately blue.
	centerIdx := (5*10 + 5) * 3
	r, g, b := result.Data[centerIdx], result.Data[centerIdx+1], result.Data[centerIdx+2]
	if r > 50 || g > 50 || b < 200 {
		t.Errorf("Expected blue-ish pixel, got RGB(%d, %d, %d)", r, g, b)
	}
}

func TestDCTDecoder_EncodeGray(t *testing.T) {
	decoder := NewDCTDecoder()

	// Create raw grayscale data.
	width, height := 20, 20
	data := make([]byte, width*height)
	for i := range data {
		data[i] = 200 // Light gray
	}

	// Encode.
	jpegData, err := decoder.EncodeGray(data, width, height, 90)
	if err != nil {
		t.Fatalf("Failed to encode grayscale JPEG: %v", err)
	}

	// Decode and verify.
	result, err := decoder.DecodeWithMetadata(jpegData)
	if err != nil {
		t.Fatalf("Failed to decode encoded grayscale JPEG: %v", err)
	}

	if result.Width != 20 || result.Height != 20 {
		t.Errorf("Expected 20x20, got %dx%d", result.Width, result.Height)
	}
}

func TestDCTDecoder_Encode_InvalidData(t *testing.T) {
	decoder := NewDCTDecoder()

	// Wrong data length.
	data := make([]byte, 100) // Not divisible by 3
	_, err := decoder.Encode(data, 10, 10, 90)
	if err == nil {
		t.Error("Expected error for invalid data length")
	}
}

func TestDCTDecoder_DefaultQuality(t *testing.T) {
	decoder := NewDCTDecoder()

	// Create RGB data.
	width, height := 5, 5
	data := make([]byte, width*height*3)
	for i := range data {
		data[i] = 128
	}

	// Encode with quality 0 (should use default).
	jpegData, err := decoder.Encode(data, width, height, 0)
	if err != nil {
		t.Fatalf("Failed to encode with default quality: %v", err)
	}

	if len(jpegData) == 0 {
		t.Error("Expected non-empty JPEG data")
	}
}

func TestNewDCTDecoderWithParams(t *testing.T) {
	decoder := NewDCTDecoderWithParams(0)
	if decoder.ColorTransform != 0 {
		t.Errorf("Expected ColorTransform 0, got %d", decoder.ColorTransform)
	}

	decoder = NewDCTDecoderWithParams(1)
	if decoder.ColorTransform != 1 {
		t.Errorf("Expected ColorTransform 1, got %d", decoder.ColorTransform)
	}
}

func BenchmarkDCTDecoder_Decode(b *testing.B) {
	decoder := NewDCTDecoder()
	jpegData := createTestJPEG(200, 200, color.RGBA{100, 150, 200, 255}, 85)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := decoder.Decode(jpegData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDCTDecoder_Encode(b *testing.B) {
	decoder := NewDCTDecoder()
	width, height := 200, 200
	data := make([]byte, width*height*3)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := decoder.Encode(data, width, height, 85)
		if err != nil {
			b.Fatal(err)
		}
	}
}
