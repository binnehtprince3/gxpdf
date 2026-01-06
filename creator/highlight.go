package creator

import (
	"github.com/coregx/gxpdf/internal/document"
)

// HighlightAnnotation represents a highlight markup annotation.
//
// Highlight annotations mark text with a colored overlay (typically yellow).
//
// Example:
//
//	highlight := creator.NewHighlightAnnotation(100, 650, 300, 670)
//	highlight.SetColor(creator.Yellow)
//	highlight.SetAuthor("John Doe")
//	page.AddHighlightAnnotation(highlight)
type HighlightAnnotation struct {
	x1     float64 // Left X coordinate
	y1     float64 // Bottom Y coordinate
	x2     float64 // Right X coordinate
	y2     float64 // Top Y coordinate
	color  Color   // Highlight color
	author string  // Author name
	note   string  // Optional note text
}

// NewHighlightAnnotation creates a new highlight annotation.
//
// The highlight covers the rectangular area from (x1, y1) to (x2, y2).
//
// Parameters:
//   - x1: Left X coordinate (from left edge)
//   - y1: Bottom Y coordinate (from bottom edge)
//   - x2: Right X coordinate (from left edge)
//   - y2: Top Y coordinate (from bottom edge)
//
// Example:
//
//	// Highlight text from (100, 650) to (300, 670)
//	highlight := creator.NewHighlightAnnotation(100, 650, 300, 670)
//	highlight.SetColor(creator.Yellow)
func NewHighlightAnnotation(x1, y1, x2, y2 float64) *HighlightAnnotation {
	return &HighlightAnnotation{
		x1:     x1,
		y1:     y1,
		x2:     x2,
		y2:     y2,
		color:  Yellow, // Default to yellow
		author: "",
		note:   "",
	}
}

// SetColor sets the highlight color.
//
// Common highlight colors:
//   - Yellow (default)
//   - Cyan (informational)
//   - Magenta (important)
//   - Green (approved)
//
// Example:
//
//	highlight.SetColor(creator.Yellow)
func (a *HighlightAnnotation) SetColor(color Color) *HighlightAnnotation {
	a.color = color
	return a
}

// SetAuthor sets the author name.
//
// Example:
//
//	highlight.SetAuthor("John Doe")
func (a *HighlightAnnotation) SetAuthor(author string) *HighlightAnnotation {
	a.author = author
	return a
}

// SetNote sets an optional note text.
//
// This text appears when the user hovers over or clicks the highlight.
//
// Example:
//
//	highlight.SetNote("Important point")
func (a *HighlightAnnotation) SetNote(note string) *HighlightAnnotation {
	a.note = note
	return a
}

// toDomain converts the Creator API annotation to a domain annotation.
func (a *HighlightAnnotation) toDomain() *document.MarkupAnnotation {
	rect := [4]float64{a.x1, a.y1, a.x2, a.y2}

	// QuadPoints define the area to highlight.
	// Format: [x1, y1, x2, y2, x3, y3, x4, y4]
	// Points are: top-left, top-right, bottom-left, bottom-right.
	quadPoints := [][8]float64{
		{a.x1, a.y2, a.x2, a.y2, a.x1, a.y1, a.x2, a.y1},
	}

	domainAnnot := document.NewMarkupAnnotation(
		document.AnnotationTypeHighlight,
		rect,
		quadPoints,
	)
	domainAnnot.SetColor([3]float64{a.color.R, a.color.G, a.color.B})
	domainAnnot.SetAuthor(a.author)
	domainAnnot.SetContents(a.note)

	return domainAnnot
}

// UnderlineAnnotation represents an underline markup annotation.
//
// Underline annotations draw a line under text.
//
// Example:
//
//	underline := creator.NewUnderlineAnnotation(100, 650, 300, 670)
//	underline.SetColor(creator.Blue)
//	page.AddUnderlineAnnotation(underline)
type UnderlineAnnotation struct {
	x1     float64 // Left X coordinate
	y1     float64 // Bottom Y coordinate
	x2     float64 // Right X coordinate
	y2     float64 // Top Y coordinate
	color  Color   // Underline color
	author string  // Author name
	note   string  // Optional note text
}

// NewUnderlineAnnotation creates a new underline annotation.
//
// The underline is drawn under the rectangular area from (x1, y1) to (x2, y2).
//
// Parameters:
//   - x1: Left X coordinate (from left edge)
//   - y1: Bottom Y coordinate (from bottom edge)
//   - x2: Right X coordinate (from left edge)
//   - y2: Top Y coordinate (from bottom edge)
//
// Example:
//
//	underline := creator.NewUnderlineAnnotation(100, 650, 300, 670)
//	underline.SetColor(creator.Blue)
func NewUnderlineAnnotation(x1, y1, x2, y2 float64) *UnderlineAnnotation {
	return &UnderlineAnnotation{
		x1:     x1,
		y1:     y1,
		x2:     x2,
		y2:     y2,
		color:  Blue, // Default to blue
		author: "",
		note:   "",
	}
}

// SetColor sets the underline color.
//
// Example:
//
//	underline.SetColor(creator.Blue)
func (a *UnderlineAnnotation) SetColor(color Color) *UnderlineAnnotation {
	a.color = color
	return a
}

// SetAuthor sets the author name.
//
// Example:
//
//	underline.SetAuthor("John Doe")
func (a *UnderlineAnnotation) SetAuthor(author string) *UnderlineAnnotation {
	a.author = author
	return a
}

// SetNote sets an optional note text.
//
// Example:
//
//	underline.SetNote("Check this")
func (a *UnderlineAnnotation) SetNote(note string) *UnderlineAnnotation {
	a.note = note
	return a
}

// toDomain converts the Creator API annotation to a domain annotation.
func (a *UnderlineAnnotation) toDomain() *document.MarkupAnnotation {
	rect := [4]float64{a.x1, a.y1, a.x2, a.y2}

	quadPoints := [][8]float64{
		{a.x1, a.y2, a.x2, a.y2, a.x1, a.y1, a.x2, a.y1},
	}

	domainAnnot := document.NewMarkupAnnotation(
		document.AnnotationTypeUnderline,
		rect,
		quadPoints,
	)
	domainAnnot.SetColor([3]float64{a.color.R, a.color.G, a.color.B})
	domainAnnot.SetAuthor(a.author)
	domainAnnot.SetContents(a.note)

	return domainAnnot
}

// StrikeOutAnnotation represents a strikeout markup annotation.
//
// StrikeOut annotations draw a line through text.
//
// Example:
//
//	strikeout := creator.NewStrikeOutAnnotation(100, 650, 300, 670)
//	strikeout.SetColor(creator.Red)
//	page.AddStrikeOutAnnotation(strikeout)
type StrikeOutAnnotation struct {
	x1     float64 // Left X coordinate
	y1     float64 // Bottom Y coordinate
	x2     float64 // Right X coordinate
	y2     float64 // Top Y coordinate
	color  Color   // StrikeOut color
	author string  // Author name
	note   string  // Optional note text
}

// NewStrikeOutAnnotation creates a new strikeout annotation.
//
// The strikeout line is drawn through the rectangular area from (x1, y1) to (x2, y2).
//
// Parameters:
//   - x1: Left X coordinate (from left edge)
//   - y1: Bottom Y coordinate (from bottom edge)
//   - x2: Right X coordinate (from left edge)
//   - y2: Top Y coordinate (from bottom edge)
//
// Example:
//
//	strikeout := creator.NewStrikeOutAnnotation(100, 650, 300, 670)
//	strikeout.SetColor(creator.Red)
func NewStrikeOutAnnotation(x1, y1, x2, y2 float64) *StrikeOutAnnotation {
	return &StrikeOutAnnotation{
		x1:     x1,
		y1:     y1,
		x2:     x2,
		y2:     y2,
		color:  Red, // Default to red
		author: "",
		note:   "",
	}
}

// SetColor sets the strikeout color.
//
// Example:
//
//	strikeout.SetColor(creator.Red)
func (a *StrikeOutAnnotation) SetColor(color Color) *StrikeOutAnnotation {
	a.color = color
	return a
}

// SetAuthor sets the author name.
//
// Example:
//
//	strikeout.SetAuthor("John Doe")
func (a *StrikeOutAnnotation) SetAuthor(author string) *StrikeOutAnnotation {
	a.author = author
	return a
}

// SetNote sets an optional note text.
//
// Example:
//
//	strikeout.SetNote("Obsolete")
func (a *StrikeOutAnnotation) SetNote(note string) *StrikeOutAnnotation {
	a.note = note
	return a
}

// toDomain converts the Creator API annotation to a domain annotation.
func (a *StrikeOutAnnotation) toDomain() *document.MarkupAnnotation {
	rect := [4]float64{a.x1, a.y1, a.x2, a.y2}

	quadPoints := [][8]float64{
		{a.x1, a.y2, a.x2, a.y2, a.x1, a.y1, a.x2, a.y1},
	}

	domainAnnot := document.NewMarkupAnnotation(
		document.AnnotationTypeStrikeOut,
		rect,
		quadPoints,
	)
	domainAnnot.SetColor([3]float64{a.color.R, a.color.G, a.color.B})
	domainAnnot.SetAuthor(a.author)
	domainAnnot.SetContents(a.note)

	return domainAnnot
}
