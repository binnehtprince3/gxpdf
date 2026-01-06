// Package writer provides infrastructure for writing PDF files.
//
// This package implements the low-level PDF writing functionality,
// including object management, cross-reference tables, and file structure.
package writer

import (
	"fmt"
	"io"
)

// IndirectObject represents a PDF indirect object.
//
// In PDF format, indirect objects are uniquely identified by:
// - Object number (positive integer)
// - Generation number (usually 0 for new objects)
//
// Format in PDF file:
//
//	N G obj
//	... object data ...
//	endobj
//
// Example:
//
//	1 0 obj
//	<< /Type /Catalog /Pages 2 0 R >>
//	endobj
type IndirectObject struct {
	// Number is the object number (must be positive).
	Number int

	// Generation is the generation number (usually 0).
	Generation int

	// Data contains the serialized object data (dictionary, array, etc.).
	Data []byte
}

// NewIndirectObject creates a new indirect object.
//
// Parameters:
//   - number: Object number (must be positive)
//   - generation: Generation number (usually 0)
//   - data: Serialized object data
//
// Example:
//
//	catalogData := []byte("<< /Type /Catalog /Pages 2 0 R >>")
//	obj := NewIndirectObject(1, 0, catalogData)
func NewIndirectObject(number, generation int, data []byte) *IndirectObject {
	return &IndirectObject{
		Number:     number,
		Generation: generation,
		Data:       data,
	}
}

// WriteTo writes the indirect object to the writer.
//
// Format:
//
//	N G obj
//	<data>
//	endobj
//
// Returns the number of bytes written and any error.
func (o *IndirectObject) WriteTo(w io.Writer) (int64, error) {
	var totalBytes int64

	// Write object header: "N G obj\n"
	header := fmt.Sprintf("%d %d obj\n", o.Number, o.Generation)
	n, err := w.Write([]byte(header))
	if err != nil {
		return int64(n), fmt.Errorf("failed to write object header: %w", err)
	}
	totalBytes += int64(n)

	// Write object data
	n, err = w.Write(o.Data)
	if err != nil {
		return totalBytes + int64(n), fmt.Errorf("failed to write object data: %w", err)
	}
	totalBytes += int64(n)

	// Write newline after data if not present
	if len(o.Data) > 0 && o.Data[len(o.Data)-1] != '\n' {
		n, err = w.Write([]byte("\n"))
		if err != nil {
			return totalBytes + int64(n), fmt.Errorf("failed to write newline: %w", err)
		}
		totalBytes += int64(n)
	}

	// Write object footer: "endobj\n"
	n, err = w.Write([]byte("endobj\n"))
	if err != nil {
		return totalBytes + int64(n), fmt.Errorf("failed to write endobj: %w", err)
	}
	totalBytes += int64(n)

	return totalBytes, nil
}

// String returns a string representation of the object (for debugging).
func (o *IndirectObject) String() string {
	return fmt.Sprintf("Object %d %d: %d bytes", o.Number, o.Generation, len(o.Data))
}
