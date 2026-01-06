// Package valueobjects provides immutable value objects for the PDF domain.
package types

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidRectangle is returned when rectangle dimensions are invalid.
	ErrInvalidRectangle = errors.New("invalid rectangle: upper-right must be greater than lower-left")
)

// Rectangle represents an immutable PDF rectangle.
// In PDF, a rectangle is defined by two points: lower-left (llx, lly) and upper-right (urx, ury).
// This is a Value Object in DDD terms - immutable and compared by value.
type Rectangle struct {
	llx, lly float64 // Lower-left corner
	urx, ury float64 // Upper-right corner
}

// NewRectangle creates a new Rectangle with validation.
// Returns error if dimensions are invalid (urx <= llx or ury <= lly).
func NewRectangle(llx, lly, urx, ury float64) (Rectangle, error) {
	if urx <= llx {
		return Rectangle{}, fmt.Errorf("%w: urx (%f) <= llx (%f)", ErrInvalidRectangle, urx, llx)
	}
	if ury <= lly {
		return Rectangle{}, fmt.Errorf("%w: ury (%f) <= lly (%f)", ErrInvalidRectangle, ury, lly)
	}
	return Rectangle{llx, lly, urx, ury}, nil
}

// MustRectangle creates a Rectangle and panics on error.
// Use only when you're certain the dimensions are valid (e.g., constants).
func MustRectangle(llx, lly, urx, ury float64) Rectangle {
	r, err := NewRectangle(llx, lly, urx, ury)
	if err != nil {
		panic(err)
	}
	return r
}

// Width returns the width of the rectangle.
func (r Rectangle) Width() float64 {
	return r.urx - r.llx
}

// Height returns the height of the rectangle.
func (r Rectangle) Height() float64 {
	return r.ury - r.lly
}

// LowerLeft returns the lower-left corner coordinates.
func (r Rectangle) LowerLeft() (x, y float64) {
	return r.llx, r.lly
}

// UpperRight returns the upper-right corner coordinates.
func (r Rectangle) UpperRight() (x, y float64) {
	return r.urx, r.ury
}

// Contains checks if a point (x, y) is within the rectangle.
func (r Rectangle) Contains(x, y float64) bool {
	return x >= r.llx && x <= r.urx && y >= r.lly && y <= r.ury
}

// WithOffset returns a new Rectangle offset by (dx, dy).
// This maintains immutability - returns new instance instead of modifying.
func (r Rectangle) WithOffset(dx, dy float64) Rectangle {
	return Rectangle{
		llx: r.llx + dx,
		lly: r.lly + dy,
		urx: r.urx + dx,
		ury: r.ury + dy,
	}
}

// String returns a string representation of the rectangle.
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle[(%f, %f), (%f, %f)]", r.llx, r.lly, r.urx, r.ury)
}

// Equals checks if two rectangles are equal (value comparison).
func (r Rectangle) Equals(other Rectangle) bool {
	return r.llx == other.llx &&
		r.lly == other.lly &&
		r.urx == other.urx &&
		r.ury == other.ury
}

// Common page sizes as constants.
var (
	// A4 represents A4 page size (210mm x 297mm = 595pt x 842pt).
	A4 = MustRectangle(0, 0, 595.276, 841.890)

	// Letter represents US Letter size (8.5in x 11in = 612pt x 792pt).
	Letter = MustRectangle(0, 0, 612, 792)

	// Legal represents US Legal size (8.5in x 14in = 612pt x 1008pt).
	Legal = MustRectangle(0, 0, 612, 1008)

	// A3 represents A3 page size (297mm x 420mm = 842pt x 1191pt).
	A3 = MustRectangle(0, 0, 841.890, 1190.551)

	// A5 represents A5 page size (148mm x 210mm = 420pt x 595pt).
	A5 = MustRectangle(0, 0, 419.528, 595.276)
)
