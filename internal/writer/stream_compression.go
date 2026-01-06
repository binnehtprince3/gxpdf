// Package writer implements PDF writing infrastructure.
package writer

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

// CompressionLevel defines compression level for streams.
//
// PDF uses FlateDecode filter, which is based on zlib (deflate with header and checksum).
//
// Levels:
//   - NoCompression: Store data without compression (fastest, largest)
//   - BestSpeed: Fast compression with lower compression ratio
//   - DefaultCompression: Balanced speed and compression (recommended)
//   - BestCompression: Maximum compression (slowest, smallest)
//
// Reference: PDF 1.7 Specification, Table 3.4 (Standard Filters).
type CompressionLevel int

const (
	// NoCompression disables compression (level 0).
	NoCompression CompressionLevel = 0

	// BestSpeed uses fastest compression (level 1).
	BestSpeed CompressionLevel = 1

	// DefaultCompression uses default zlib compression (level -1).
	// This is the recommended setting for most use cases.
	DefaultCompression CompressionLevel = -1

	// BestCompression uses maximum compression (level 9).
	BestCompression CompressionLevel = 9
)

// CompressStream compresses data using zlib (FlateDecode in PDF terminology).
//
// PDF's FlateDecode filter uses zlib compression, which includes:
//   - 2-byte zlib header
//   - Deflate-compressed data
//   - 4-byte Adler-32 checksum
//
// This is different from raw deflate compression.
//
// Parameters:
//   - data: Uncompressed data to compress
//   - level: Compression level (NoCompression, BestSpeed, DefaultCompression, BestCompression)
//
// Returns:
//   - compressed: Compressed data (zlib format)
//   - error: Any error that occurred
//
// Example:
//
//	content := []byte("BT /F1 12 Tf 100 700 Td (Hello) Tj ET")
//	compressed, err := CompressStream(content, DefaultCompression)
//
// Reference: PDF 1.7 Specification, Section 7.4.4 (FlateDecode Filter).
func CompressStream(data []byte, level CompressionLevel) ([]byte, error) {
	if len(data) == 0 {
		// Empty data compresses to empty (avoid creating zlib header for nothing)
		return []byte{}, nil
	}

	// Validate compression level
	if !isValidCompressionLevel(level) {
		return nil, fmt.Errorf("invalid compression level: %d (must be -1, 0-9)", level)
	}

	var buf bytes.Buffer

	// Create zlib writer with specified compression level
	w, err := zlib.NewWriterLevel(&buf, int(level))
	if err != nil {
		return nil, fmt.Errorf("create zlib writer: %w", err)
	}

	// Write uncompressed data
	if _, err := w.Write(data); err != nil {
		// Close writer to release resources (ignore error, we already have an error)
		_ = w.Close()
		return nil, fmt.Errorf("write data to zlib: %w", err)
	}

	// Close writer to flush and write checksum
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("close zlib writer: %w", err)
	}

	return buf.Bytes(), nil
}

// DecompressStream decompresses zlib data (FlateDecode filter).
//
// This is the inverse of CompressStream, used for reading compressed streams.
//
// Parameters:
//   - data: Compressed data (zlib format)
//
// Returns:
//   - decompressed: Decompressed data
//   - error: Any error that occurred (invalid format, checksum mismatch, etc.)
//
// Example:
//
//	decompressed, err := DecompressStream(compressed)
//
// Reference: PDF 1.7 Specification, Section 7.4.4 (FlateDecode Filter).
func DecompressStream(data []byte) ([]byte, error) {
	if len(data) == 0 {
		// Empty compressed data -> empty result
		return []byte{}, nil
	}

	// Create zlib reader
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create zlib reader: %w", err)
	}
	defer func() {
		// Best effort close (error already reported if read failed)
		_ = r.Close()
	}()

	// Read all decompressed data with size limit to prevent decompression bombs
	var buf bytes.Buffer
	// Limit decompressed size to 100MB (reasonable for PDF content streams)
	const maxDecompressedSize = 100 * 1024 * 1024
	limitedReader := io.LimitReader(r, maxDecompressedSize)

	n, err := io.Copy(&buf, limitedReader)
	if err != nil {
		return nil, fmt.Errorf("decompress zlib data: %w", err)
	}

	// Check if we hit the limit (potential decompression bomb)
	if n >= maxDecompressedSize {
		return nil, fmt.Errorf("decompressed data exceeds maximum size (%d bytes)", maxDecompressedSize)
	}

	// Close reader to verify checksum
	if err := r.Close(); err != nil {
		return nil, fmt.Errorf("verify zlib checksum: %w", err)
	}

	return buf.Bytes(), nil
}

// isValidCompressionLevel checks if compression level is valid.
//
// Valid levels: -1 (default), 0 (no compression), 1-9 (compression levels).
func isValidCompressionLevel(level CompressionLevel) bool {
	return level == DefaultCompression || (level >= NoCompression && level <= BestCompression)
}

// EstimateCompressionRatio estimates the compression ratio for given data.
//
// This is useful for deciding whether to compress a stream:
//   - Ratio < 0.9: Good compression, use FlateDecode
//   - Ratio >= 0.9: Poor compression, store uncompressed
//
// Parameters:
//   - data: Uncompressed data
//
// Returns:
//   - ratio: Estimated compression ratio (compressed size / original size)
//
// Note: This actually compresses the data to get an accurate estimate.
// For a fast heuristic, check if data contains repeated patterns.
func EstimateCompressionRatio(data []byte) float64 {
	if len(data) == 0 {
		return 1.0
	}

	compressed, err := CompressStream(data, DefaultCompression)
	if err != nil {
		return 1.0 // Assume no compression on error
	}

	return float64(len(compressed)) / float64(len(data))
}

// ShouldCompress determines if data should be compressed based on size and content.
//
// Heuristic:
//   - Data < 50 bytes: Don't compress (overhead not worth it)
//   - Data >= 50 bytes: Compress (likely worth it for text content)
//
// For more precise control, use EstimateCompressionRatio.
//
// Parameters:
//   - data: Uncompressed data
//
// Returns:
//   - bool: true if compression is recommended
func ShouldCompress(data []byte) bool {
	// Small streams have too much overhead relative to compression gains
	const minSizeForCompression = 50
	return len(data) >= minSizeForCompression
}
