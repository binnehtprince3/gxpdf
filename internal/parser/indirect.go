package parser

import (
	"fmt"
	"io"
)

// IndirectObject represents an indirect PDF object.
// Format: objNum genNum obj ... endobj
//
// Example: 1 0 obj (Hello) endobj
//
// Reference: PDF 1.7 specification, Section 7.3.10 (Indirect Objects).
type IndirectObject struct {
	Number     int       // Object number
	Generation int       // Generation number
	Object     PdfObject // The actual object
}

// NewIndirectObject creates a new indirect object.
func NewIndirectObject(number, generation int, obj PdfObject) *IndirectObject {
	return &IndirectObject{
		Number:     number,
		Generation: generation,
		Object:     obj,
	}
}

// String returns a string representation of the indirect object.
func (o *IndirectObject) String() string {
	return fmt.Sprintf("%d %d obj %v endobj", o.Number, o.Generation, o.Object)
}

// WriteTo writes the PDF representation to w.
func (o *IndirectObject) WriteTo(w io.Writer) (int64, error) {
	var total int64

	// Write object header: "N G obj"
	header := fmt.Sprintf("%d %d obj\n", o.Number, o.Generation)
	n, err := w.Write([]byte(header))
	total += int64(n)
	if err != nil {
		return total, err
	}

	// Write the object
	written, err := o.Object.WriteTo(w)
	total += written
	if err != nil {
		return total, err
	}

	// Write newline and endobj
	n, err = w.Write([]byte("\nendobj\n"))
	total += int64(n)
	return total, err
}

// IndirectReference represents a reference to an indirect object.
// Format: objNum genNum R
//
// Example: 1 0 R (refers to object 1, generation 0)
//
// Reference: PDF 1.7 specification, Section 7.3.10 (Indirect Objects).
type IndirectReference struct {
	Number     int // Object number being referenced
	Generation int // Generation number being referenced
}

// NewIndirectReference creates a new indirect reference.
func NewIndirectReference(number, generation int) *IndirectReference {
	return &IndirectReference{
		Number:     number,
		Generation: generation,
	}
}

// String returns the PDF representation of the reference.
func (r *IndirectReference) String() string {
	return fmt.Sprintf("%d %d R", r.Number, r.Generation)
}

// WriteTo writes the PDF representation to w.
func (r *IndirectReference) WriteTo(w io.Writer) (int64, error) {
	str := r.String()
	n, err := w.Write([]byte(str))
	return int64(n), err
}

// Equals checks if two references point to the same object.
func (r *IndirectReference) Equals(other *IndirectReference) bool {
	if other == nil {
		return false
	}
	return r.Number == other.Number && r.Generation == other.Generation
}

// Clone creates a copy of the indirect reference.
func (r *IndirectReference) Clone() *IndirectReference {
	return &IndirectReference{
		Number:     r.Number,
		Generation: r.Generation,
	}
}
