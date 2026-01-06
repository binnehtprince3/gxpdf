package creator

import "github.com/coregx/gxpdf/internal/document"

// PageSize represents standard PDF page sizes.
//
// Common page sizes are provided as constants (A4, Letter, etc.).
// Custom sizes can be created using the CustomSize function.
type PageSize int

const (
	// A4 paper size (210 × 297 mm or 595 × 842 points).
	// This is the most common paper size worldwide.
	A4 PageSize = iota

	// Letter paper size (8.5 × 11 inches or 612 × 792 points).
	// This is the standard size in North America.
	Letter

	// Legal paper size (8.5 × 14 inches or 612 × 1008 points).
	Legal

	// Tabloid paper size (11 × 17 inches or 792 × 1224 points).
	// Also known as Ledger when in landscape orientation.
	Tabloid

	// A3 paper size (297 × 420 mm or 842 × 1191 points).
	// Twice the size of A4.
	A3

	// A5 paper size (148 × 210 mm or 420 × 595 points).
	// Half the size of A4.
	A5

	// B4 paper size (250 × 353 mm or 709 × 1001 points).
	B4

	// B5 paper size (176 × 250 mm or 499 × 709 points).
	B5
)

// toDomainSize converts creator PageSize to domain PageSize.
//
// This is an internal method used by the Creator to work with the domain layer.
func (ps PageSize) toDomainSize() document.PageSize {
	switch ps {
	case A4:
		return document.A4
	case Letter:
		return document.Letter
	case Legal:
		return document.Legal
	case Tabloid:
		return document.Tabloid
	case A3:
		return document.A3
	case A5:
		return document.A5
	case B4:
		return document.B4
	case B5:
		return document.B5
	default:
		return document.A4 // Default to A4
	}
}

// String returns the name of the page size.
func (ps PageSize) String() string {
	switch ps {
	case A4:
		return "A4"
	case Letter:
		return "Letter"
	case Legal:
		return "Legal"
	case Tabloid:
		return "Tabloid"
	case A3:
		return "A3"
	case A5:
		return "A5"
	case B4:
		return "B4"
	case B5:
		return "B5"
	default:
		return "Unknown"
	}
}
