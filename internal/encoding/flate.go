// Package encoding implements PDF stream encoding and decoding filters.
package encoding

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

// FlateDecoder implements FlateDecode (zlib/deflate) stream decompression.
//
// FlateDecode is the most common compression filter in PDF files,
// using the zlib/deflate algorithm (RFC 1950/1951).
//
// Reference: PDF 1.7 specification, Section 7.4.4 (FlateDecode Filter).
type FlateDecoder struct{}

// NewFlateDecoder creates a new Flate decoder.
func NewFlateDecoder() *FlateDecoder {
	return &FlateDecoder{}
}

// Decode decompresses Flate-encoded data.
//
// This is a straightforward zlib decompression without predictor support.
// Predictors (like PNG filters) are typically not used for xref streams.
//
// Parameters:
//   - data: Compressed data bytes
//
// Returns: Decompressed data bytes, or error if decompression fails.
func (d *FlateDecoder) Decode(data []byte) (result []byte, err error) {
	// Create zlib reader
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer func() {
		if closeErr := reader.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close zlib reader: %w", closeErr)
		}
	}()

	// Read all decompressed data
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	return buf.Bytes(), nil
}

// Encode compresses data using Flate encoding.
//
// This is for future PDF writing support.
func (d *FlateDecoder) Encode(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		_ = writer.Close() // Best effort close on write error
		return nil, fmt.Errorf("failed to compress data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zlib writer: %w", err)
	}

	return buf.Bytes(), nil
}
