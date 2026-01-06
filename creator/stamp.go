package creator

import (
	"github.com/coregx/gxpdf/internal/document"
)

// StampAnnotation represents a rubber stamp annotation.
//
// Stamp annotations display predefined stamps like "Approved", "Draft", etc.
//
// Example:
//
//	stamp := creator.NewStampAnnotation(300, 700, 100, 50, creator.StampApproved)
//	stamp.SetColor(creator.Green)
//	stamp.SetAuthor("John Doe")
//	page.AddStampAnnotation(stamp)
type StampAnnotation struct {
	x      float64            // X coordinate (from left)
	y      float64            // Y coordinate (from bottom)
	width  float64            // Stamp width
	height float64            // Stamp height
	name   document.StampName // Stamp name (Approved, Draft, etc.)
	color  Color              // Stamp color
	author string             // Author name
	note   string             // Optional note text
}

// Predefined stamp names (exported for user convenience).
const (
	// StampApproved represents an "Approved" stamp.
	StampApproved = document.StampApproved
	// StampNotApproved represents a "Not Approved" stamp.
	StampNotApproved = document.StampNotApproved
	// StampDraft represents a "Draft" stamp.
	StampDraft = document.StampDraft
	// StampFinal represents a "Final" stamp.
	StampFinal = document.StampFinal
	// StampConfidential represents a "Confidential" stamp.
	StampConfidential = document.StampConfidential
	// StampForComment represents a "For Comment" stamp.
	StampForComment = document.StampForComment
	// StampForPublicRelease represents a "For Public Release" stamp.
	StampForPublicRelease = document.StampForPublicRelease
	// StampAsIs represents an "As Is" stamp.
	StampAsIs = document.StampAsIs
	// StampDepartmental represents a "Departmental" stamp.
	StampDepartmental = document.StampDepartmental
	// StampExperimental represents an "Experimental" stamp.
	StampExperimental = document.StampExperimental
	// StampExpired represents an "Expired" stamp.
	StampExpired = document.StampExpired
	// StampNotForPublicRelease represents a "Not For Public Release" stamp.
	StampNotForPublicRelease = document.StampNotForPublicRelease
)

// NewStampAnnotation creates a new stamp annotation.
//
// The stamp appears at (x, y) with the specified width and height.
//
// Parameters:
//   - x: Horizontal position in points (from left edge)
//   - y: Vertical position in points (from bottom edge)
//   - width: Stamp width in points
//   - height: Stamp height in points
//   - name: Stamp name (e.g., creator.StampApproved)
//
// Example:
//
//	// Create an "Approved" stamp at (300, 700), 100x50 points
//	stamp := creator.NewStampAnnotation(300, 700, 100, 50, creator.StampApproved)
//	stamp.SetColor(creator.Green)
func NewStampAnnotation(x, y, width, height float64, name document.StampName) *StampAnnotation {
	return &StampAnnotation{
		x:      x,
		y:      y,
		width:  width,
		height: height,
		name:   name,
		color:  Red, // Default to red
		author: "",
		note:   "",
	}
}

// SetColor sets the stamp color.
//
// Common stamp colors:
//   - Red (default, for important/rejected)
//   - Green (approved/ok)
//   - Blue (informational)
//   - Yellow (warning)
//
// Example:
//
//	stamp.SetColor(creator.Green) // Approved
func (a *StampAnnotation) SetColor(color Color) *StampAnnotation {
	a.color = color
	return a
}

// SetAuthor sets the author name.
//
// Example:
//
//	stamp.SetAuthor("John Doe")
func (a *StampAnnotation) SetAuthor(author string) *StampAnnotation {
	a.author = author
	return a
}

// SetNote sets an optional note text.
//
// This text appears when the user hovers over or clicks the stamp.
//
// Example:
//
//	stamp.SetNote("Approved on 2025-01-06")
func (a *StampAnnotation) SetNote(note string) *StampAnnotation {
	a.note = note
	return a
}

// toDomain converts the Creator API annotation to a domain annotation.
func (a *StampAnnotation) toDomain() *document.StampAnnotation {
	rect := [4]float64{
		a.x,            // x1 (left)
		a.y,            // y1 (bottom)
		a.x + a.width,  // x2 (right)
		a.y + a.height, // y2 (top)
	}

	domainAnnot := document.NewStampAnnotation(rect, a.name)
	domainAnnot.SetColor([3]float64{a.color.R, a.color.G, a.color.B})
	domainAnnot.SetAuthor(a.author)
	domainAnnot.SetContents(a.note)

	return domainAnnot
}
