package writer

import (
	"bytes"
	"crypto/rand"
	"strings"
	"testing"
)

// TestCompressStream tests basic compression functionality.
//
//nolint:funlen // Table-driven test with many cases
func TestCompressStream(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		level CompressionLevel
	}{
		{
			name:  "empty data",
			data:  []byte{},
			level: DefaultCompression,
		},
		{
			name:  "small text",
			data:  []byte("Hello World"),
			level: DefaultCompression,
		},
		{
			name:  "pdf content stream",
			data:  []byte("BT /F1 12 Tf 100 700 Td (Hello World) Tj ET\n"),
			level: DefaultCompression,
		},
		{
			name:  "repeated pattern",
			data:  []byte(strings.Repeat("A", 1000)),
			level: DefaultCompression,
		},
		{
			name:  "large text",
			data:  []byte(strings.Repeat("Lorem ipsum dolor sit amet ", 100)),
			level: DefaultCompression,
		},
		{
			name:  "no compression",
			data:  []byte("Test data"),
			level: NoCompression,
		},
		{
			name:  "best speed",
			data:  []byte(strings.Repeat("Test ", 200)),
			level: BestSpeed,
		},
		{
			name:  "best compression",
			data:  []byte(strings.Repeat("Test ", 200)),
			level: BestCompression,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compress
			compressed, err := CompressStream(tt.data, tt.level)
			if err != nil {
				t.Fatalf("CompressStream() error = %v", err)
			}

			// Verify empty input -> empty output
			if len(tt.data) == 0 {
				if len(compressed) != 0 {
					t.Errorf("empty data should compress to empty, got %d bytes", len(compressed))
				}
				return
			}

			// Verify compression produced output
			if len(compressed) == 0 {
				t.Error("compressed data is empty for non-empty input")
			}

			// Decompress
			decompressed, err := DecompressStream(compressed)
			if err != nil {
				t.Fatalf("DecompressStream() error = %v", err)
			}

			// Verify round-trip
			if !bytes.Equal(decompressed, tt.data) {
				t.Errorf("round-trip failed:\noriginal:  %q\ndecompressed: %q", tt.data, decompressed)
			}
		})
	}
}

// TestCompressStreamInvalidLevel tests invalid compression levels.
func TestCompressStreamInvalidLevel(t *testing.T) {
	data := []byte("test")

	invalidLevels := []CompressionLevel{-2, 10, 100}

	for _, level := range invalidLevels {
		t.Run("", func(t *testing.T) {
			_, err := CompressStream(data, level)
			if err == nil {
				t.Errorf("CompressStream with level %d should fail", level)
			}
		})
	}
}

// TestDecompressStream tests decompression.
func TestDecompressStream(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		decompressed, err := DecompressStream([]byte{})
		if err != nil {
			t.Fatalf("DecompressStream() error = %v", err)
		}
		if len(decompressed) != 0 {
			t.Errorf("empty compressed data should decompress to empty, got %d bytes", len(decompressed))
		}
	})

	t.Run("invalid zlib data", func(t *testing.T) {
		invalidData := []byte{0xFF, 0xFF, 0xFF, 0xFF}
		_, err := DecompressStream(invalidData)
		if err == nil {
			t.Error("DecompressStream should fail for invalid data")
		}
	})

	t.Run("truncated data", func(t *testing.T) {
		// Create valid compressed data
		original := []byte("test data")
		compressed, err := CompressStream(original, DefaultCompression)
		if err != nil {
			t.Fatalf("CompressStream() error = %v", err)
		}

		// Truncate it
		truncated := compressed[:len(compressed)/2]

		// Should fail
		_, err = DecompressStream(truncated)
		if err == nil {
			t.Error("DecompressStream should fail for truncated data")
		}
	})
}

// TestCompressionLevels verifies compression ratios for different levels.
func TestCompressionLevels(t *testing.T) {
	// Create highly compressible data
	data := []byte(strings.Repeat("AAAA", 500)) // 2000 bytes of 'A'

	tests := []struct {
		level         CompressionLevel
		maxSize       int // Maximum expected compressed size
		expectSmaller bool
	}{
		{NoCompression, 2100, false},    // Stored, slight overhead
		{BestSpeed, 100, true},          // Fast compression, should be small
		{DefaultCompression, 100, true}, // Default, should be small
		{BestCompression, 100, true},    // Best compression, should be smallest
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			compressed, err := CompressStream(data, tt.level)
			if err != nil {
				t.Fatalf("CompressStream(level=%d) error = %v", tt.level, err)
			}

			if len(compressed) > tt.maxSize {
				t.Errorf("compressed size %d exceeds max %d for level %d",
					len(compressed), tt.maxSize, tt.level)
			}

			if tt.expectSmaller && len(compressed) >= len(data) {
				t.Errorf("compression level %d should produce smaller output, got %d bytes (original: %d)",
					tt.level, len(compressed), len(data))
			}

			// Verify round-trip
			decompressed, err := DecompressStream(compressed)
			if err != nil {
				t.Fatalf("DecompressStream() error = %v", err)
			}

			if !bytes.Equal(decompressed, data) {
				t.Error("round-trip failed for compression level", tt.level)
			}
		})
	}
}

// TestCompressionRatios verifies compression effectiveness.
func TestCompressionRatios(t *testing.T) {
	tests := []struct {
		name          string
		data          []byte
		expectedRatio float64 // Maximum expected ratio (smaller = better compression)
	}{
		{
			name:          "highly repetitive",
			data:          []byte(strings.Repeat("A", 1000)),
			expectedRatio: 0.05, // Should compress very well
		},
		{
			name:          "text content",
			data:          []byte(strings.Repeat("Lorem ipsum dolor sit amet ", 50)),
			expectedRatio: 0.5, // Should compress reasonably
		},
		{
			name:          "pdf operators",
			data:          []byte(strings.Repeat("BT /F1 12 Tf 100 700 Td (Hello) Tj ET\n", 20)),
			expectedRatio: 0.4, // Should compress well (repeated operators)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := CompressStream(tt.data, DefaultCompression)
			if err != nil {
				t.Fatalf("CompressStream() error = %v", err)
			}

			ratio := float64(len(compressed)) / float64(len(tt.data))

			if ratio > tt.expectedRatio {
				t.Errorf("compression ratio %.2f exceeds expected %.2f (compressed: %d, original: %d)",
					ratio, tt.expectedRatio, len(compressed), len(tt.data))
			}
		})
	}
}

// TestBinaryData tests compression of binary data.
func TestBinaryData(t *testing.T) {
	// Create random binary data
	data := make([]byte, 1024)
	if _, err := rand.Read(data); err != nil {
		t.Fatalf("failed to generate random data: %v", err)
	}

	// Compress
	compressed, err := CompressStream(data, DefaultCompression)
	if err != nil {
		t.Fatalf("CompressStream() error = %v", err)
	}

	// Random data typically doesn't compress well, but should still work
	if len(compressed) == 0 {
		t.Error("compressed random data should not be empty")
	}

	// Decompress
	decompressed, err := DecompressStream(compressed)
	if err != nil {
		t.Fatalf("DecompressStream() error = %v", err)
	}

	// Verify round-trip
	if !bytes.Equal(decompressed, data) {
		t.Error("round-trip failed for binary data")
	}
}

// TestEstimateCompressionRatio tests compression ratio estimation.
func TestEstimateCompressionRatio(t *testing.T) {
	tests := []struct {
		name          string
		data          []byte
		expectedRatio float64
		tolerance     float64
	}{
		{
			name:          "empty data",
			data:          []byte{},
			expectedRatio: 1.0,
			tolerance:     0.0,
		},
		{
			name:          "highly compressible",
			data:          []byte(strings.Repeat("A", 1000)),
			expectedRatio: 0.05,
			tolerance:     0.1,
		},
		{
			name:          "moderately compressible",
			data:          []byte(strings.Repeat("Lorem ipsum ", 100)),
			expectedRatio: 0.1, // Repeated text compresses very well
			tolerance:     0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := EstimateCompressionRatio(tt.data)

			if ratio < tt.expectedRatio-tt.tolerance || ratio > tt.expectedRatio+tt.tolerance {
				t.Errorf("EstimateCompressionRatio() = %.2f, expected %.2f ± %.2f",
					ratio, tt.expectedRatio, tt.tolerance)
			}
		})
	}
}

// TestShouldCompress tests compression recommendation heuristic.
func TestShouldCompress(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "empty data",
			data:     []byte{},
			expected: false,
		},
		{
			name:     "very small data",
			data:     []byte("Hi"),
			expected: false,
		},
		{
			name:     "small data below threshold",
			data:     bytes.Repeat([]byte("A"), 49),
			expected: false,
		},
		{
			name:     "data at threshold",
			data:     bytes.Repeat([]byte("A"), 50),
			expected: true,
		},
		{
			name:     "large data",
			data:     bytes.Repeat([]byte("Lorem ipsum "), 100),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldCompress(tt.data)
			if result != tt.expected {
				t.Errorf("ShouldCompress() = %v, expected %v (size: %d)", result, tt.expected, len(tt.data))
			}
		})
	}
}

// TestLargeData tests compression of large data.
func TestLargeData(t *testing.T) {
	// Create 1MB of data
	data := bytes.Repeat([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. "), 20000)

	compressed, err := CompressStream(data, DefaultCompression)
	if err != nil {
		t.Fatalf("CompressStream() error = %v", err)
	}

	// Should compress significantly
	ratio := float64(len(compressed)) / float64(len(data))
	if ratio > 0.3 {
		t.Errorf("large text should compress well, got ratio %.2f", ratio)
	}

	// Verify round-trip
	decompressed, err := DecompressStream(compressed)
	if err != nil {
		t.Fatalf("DecompressStream() error = %v", err)
	}

	if !bytes.Equal(decompressed, data) {
		t.Error("round-trip failed for large data")
	}
}

// BenchmarkCompressSmall benchmarks compression of small data.
func BenchmarkCompressSmall(b *testing.B) {
	data := []byte("BT /F1 12 Tf 100 700 Td (Hello World) Tj ET\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompressStream(data, DefaultCompression)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCompressMedium benchmarks compression of medium data.
func BenchmarkCompressMedium(b *testing.B) {
	data := bytes.Repeat([]byte("BT /F1 12 Tf 100 700 Td (Lorem ipsum dolor sit amet) Tj ET\n"), 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompressStream(data, DefaultCompression)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCompressLarge benchmarks compression of large data.
func BenchmarkCompressLarge(b *testing.B) {
	data := bytes.Repeat([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. "), 10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompressStream(data, DefaultCompression)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecompress benchmarks decompression.
func BenchmarkDecompress(b *testing.B) {
	data := bytes.Repeat([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. "), 1000)
	compressed, err := CompressStream(data, DefaultCompression)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecompressStream(compressed)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCompressionLevels benchmarks different compression levels.
func BenchmarkCompressionLevels(b *testing.B) {
	data := bytes.Repeat([]byte("Lorem ipsum dolor sit amet "), 500)

	levels := []struct {
		name  string
		level CompressionLevel
	}{
		{"NoCompression", NoCompression},
		{"BestSpeed", BestSpeed},
		{"DefaultCompression", DefaultCompression},
		{"BestCompression", BestCompression},
	}

	for _, l := range levels {
		b.Run(l.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := CompressStream(data, l.level)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
