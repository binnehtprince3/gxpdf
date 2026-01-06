package creator

import (
	"github.com/coregx/gxpdf/internal/document"
)

// TextAnnotation represents a sticky note annotation in the Creator API.
//
// Text annotations appear as icons (sticky notes) on PDF pages.
// When clicked, they display pop-up text.
//
// Example:
//
//	note := creator.NewTextAnnotation(100, 700, "This is a comment")
//	note.SetAuthor("John Doe")
//	note.SetColor(creator.Yellow)
//	note.SetOpen(true)
//	page.AddTextAnnotation(note)
type TextAnnotation struct {
	x        float64 // X coordinate (from left)
	y        float64 // Y coordinate (from bottom)
	contents string  // Pop-up text
	author   string  // Author name
	color    Color   // Annotation color
	open     bool    // Open by default?
}

// NewTextAnnotation creates a new text annotation (sticky note).
//
// The annotation appears as a small icon (typically 20x20 points) at (x, y).
// When clicked, it displays the contents text in a pop-up.
//
// Parameters:
//   - x: Horizontal position in points (from left edge)
//   - y: Vertical position in points (from bottom edge)
//   - contents: Text to display in the pop-up
//
// Example:
//
//	note := creator.NewTextAnnotation(100, 700, "Review this section")
//	note.SetAuthor("Alice")
//	note.SetColor(creator.Yellow)
func NewTextAnnotation(x, y float64, contents string) *TextAnnotation {
	return &TextAnnotation{
		x:        x,
		y:        y,
		contents: contents,
		author:   "",
		color:    Yellow, // Default to yellow (like sticky notes)
		open:     false,
	}
}

// SetAuthor sets the author name for the annotation.
//
// This appears in the annotation properties and can be used to track
// who added the comment.
//
// Example:
//
//	note.SetAuthor("John Doe")
func (a *TextAnnotation) SetAuthor(author string) *TextAnnotation {
	a.author = author
	return a
}

// SetColor sets the annotation color.
//
// Common colors for sticky notes:
//   - Yellow (default)
//   - Red (important)
//   - Green (ok/approved)
//   - Blue (informational)
//
// Example:
//
//	note.SetColor(creator.Red) // Mark as important
func (a *TextAnnotation) SetColor(color Color) *TextAnnotation {
	a.color = color
	return a
}

// SetOpen sets whether the pop-up should be open by default.
//
// If true, the pop-up is visible when the PDF is opened.
// If false, the user must click the icon to see the text.
//
// Example:
//
//	note.SetOpen(true) // Show immediately
func (a *TextAnnotation) SetOpen(open bool) *TextAnnotation {
	a.open = open
	return a
}

// toDomain converts the Creator API annotation to a domain annotation.
func (a *TextAnnotation) toDomain() *document.TextAnnotation {
	// Icon size is typically 20x20 points.
	const iconSize = 20.0

	rect := [4]float64{
		a.x,            // x1 (left)
		a.y,            // y1 (bottom)
		a.x + iconSize, // x2 (right)
		a.y + iconSize, // y2 (top)
	}

	domainAnnot := document.NewTextAnnotation(rect, a.contents, a.author)
	domainAnnot.SetColor([3]float64{a.color.R, a.color.G, a.color.B})
	domainAnnot.SetOpen(a.open)

	return domainAnnot
}
