package document

import "github.com/coregx/gxpdf/internal/models/types"

// PageSize represents standard PDF page sizes.
//
// Standard sizes are provided as constants for convenience.
// For custom sizes, use CustomPageSize().
type PageSize int

const (
	// ISO 216 A series (most common international sizes)

	// A4 is 210 × 297 mm (8.27 × 11.69 in) - Most common international paper size.
	A4 PageSize = iota

	// A3 is 297 × 420 mm (11.69 × 16.54 in) - Twice the area of A4.
	A3

	// A5 is 148 × 210 mm (5.83 × 8.27 in) - Half the area of A4.
	A5

	// ISO 216 B series

	// B4 is 250 × 353 mm (9.84 × 13.90 in) - Between A3 and A4.
	B4

	// B5 is 176 × 250 mm (6.93 × 9.84 in) - Between A4 and A5.
	B5

	// North American sizes

	// Letter is 8.5 × 11 in (215.9 × 279.4 mm) - Standard US/Canada paper size.
	Letter

	// Legal is 8.5 × 14 in (215.9 × 355.6 mm) - US legal documents.
	Legal

	// Tabloid is 11 × 17 in (279.4 × 431.8 mm) - Also known as Ledger.
	Tabloid

	// Custom indicates a custom page size (use CustomPageSize function).
	Custom
)

// ToRectangle converts PageSize to Rectangle (in points, 1 point = 1/72 inch).
//
// All standard page sizes are returned in portrait orientation.
// Use Page.SetRotation() for landscape orientation.
//
// Example:
//
//	rect := document.A4.ToRectangle()
//	// rect is 595×842 points (210×297mm)
func (ps PageSize) ToRectangle() types.Rectangle {
	switch ps {
	case A4:
		// 210mm × 297mm = 8.27in × 11.69in = 595.28pt × 841.89pt ≈ 595×842pt
		return types.MustRectangle(0, 0, 595, 842)

	case A3:
		// 297mm × 420mm = 11.69in × 16.54in = 841.89pt × 1190.55pt ≈ 842×1191pt
		return types.MustRectangle(0, 0, 842, 1191)

	case A5:
		// 148mm × 210mm = 5.83in × 8.27in = 419.53pt × 595.28pt ≈ 420×595pt
		return types.MustRectangle(0, 0, 420, 595)

	case B4:
		// 250mm × 353mm = 9.84in × 13.90in = 708.66pt × 1000.63pt ≈ 709×1001pt
		return types.MustRectangle(0, 0, 709, 1001)

	case B5:
		// 176mm × 250mm = 6.93in × 9.84in = 498.90pt × 708.66pt ≈ 499×709pt
		return types.MustRectangle(0, 0, 499, 709)

	case Letter:
		// 8.5in × 11in = 612pt × 792pt
		return types.MustRectangle(0, 0, 612, 792)

	case Legal:
		// 8.5in × 14in = 612pt × 1008pt
		return types.MustRectangle(0, 0, 612, 1008)

	case Tabloid:
		// 11in × 17in = 792pt × 1224pt
		return types.MustRectangle(0, 0, 792, 1224)

	default:
		// Default to A4 if unknown size
		return types.MustRectangle(0, 0, 595, 842)
	}
}

// String returns the name of the page size.
func (ps PageSize) String() string {
	switch ps {
	case A4:
		return "A4"
	case A3:
		return "A3"
	case A5:
		return "A5"
	case B4:
		return "B4"
	case B5:
		return "B5"
	case Letter:
		return "Letter"
	case Legal:
		return "Legal"
	case Tabloid:
		return "Tabloid"
	case Custom:
		return "Custom"
	default:
		return "Unknown"
	}
}

// CustomPageSize creates a custom page size in points.
//
// Points are 1/72 of an inch.
//
// Example:
//
//	// Create a custom 6×9 inch page
//	customSize := document.CustomPageSize(6*72, 9*72)
func CustomPageSize(widthPt, heightPt float64) types.Rectangle {
	return types.MustRectangle(0, 0, widthPt, heightPt)
}

// Common conversion constants for convenience

const (
	// PointsPerInch is the number of points in one inch.
	// 1 inch = 72 points (PostScript/PDF standard)
	PointsPerInch = 72.0

	// PointsPerMM is the number of points in one millimeter.
	// 1 mm = 72/25.4 ≈ 2.83465 points
	PointsPerMM = 72.0 / 25.4

	// PointsPerCM is the number of points in one centimeter.
	// 1 cm = 72/2.54 ≈ 28.3465 points
	PointsPerCM = 72.0 / 2.54
)

// InchesToPoints converts inches to points.
func InchesToPoints(inches float64) float64 {
	return inches * PointsPerInch
}

// MMToPoints converts millimeters to points.
func MMToPoints(mm float64) float64 {
	return mm * PointsPerMM
}

// CMToPoints converts centimeters to points.
func CMToPoints(cm float64) float64 {
	return cm * PointsPerCM
}

// PointsToInches converts points to inches.
func PointsToInches(points float64) float64 {
	return points / PointsPerInch
}

// PointsToMM converts points to millimeters.
func PointsToMM(points float64) float64 {
	return points / PointsPerMM
}

// PointsToCM converts points to centimeters.
func PointsToCM(points float64) float64 {
	return points / PointsPerCM
}
