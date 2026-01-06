package document

import "errors"

// AnnotationType represents the type of PDF annotation.
type AnnotationType int

const (
	// AnnotationTypeLink represents a clickable link annotation.
	AnnotationTypeLink AnnotationType = iota
	// AnnotationTypeText represents a sticky note annotation.
	AnnotationTypeText
	// AnnotationTypeHighlight represents a highlight markup annotation.
	AnnotationTypeHighlight
	// AnnotationTypeUnderline represents an underline markup annotation.
	AnnotationTypeUnderline
	// AnnotationTypeStrikeOut represents a strikeout markup annotation.
	AnnotationTypeStrikeOut
	// AnnotationTypeStamp represents a rubber stamp annotation.
	AnnotationTypeStamp
)

// LinkAnnotation represents a clickable link in a PDF.
//
// Link annotations create clickable areas (hot spots) on PDF pages.
// They can point to external URLs or internal page destinations.
//
// Example:
//
//	// External link
//	link := NewLinkAnnotation([4]float64{100, 690, 200, 710}, "https://example.com")
//
//	// Internal link to page 3
//	link := NewInternalLinkAnnotation([4]float64{100, 690, 200, 710}, 2)
type LinkAnnotation struct {
	// Rect defines the clickable area [x1, y1, x2, y2] in PDF coordinates.
	// (x1, y1) is the lower-left corner, (x2, y2) is the upper-right corner.
	Rect [4]float64

	// URI is the target URL (for external links).
	// Empty for internal links.
	URI string

	// DestPage is the target page number (for internal links, 0-based).
	// -1 for external links.
	DestPage int

	// IsInternal indicates if this is an internal page link.
	// true = internal page link (use DestPage)
	// false = external URL link (use URI)
	IsInternal bool

	// BorderWidth is the width of the border around the clickable area.
	// 0 = no visible border (default for most links).
	BorderWidth float64
}

// NewLinkAnnotation creates a new URL link annotation.
//
// The rect parameter defines the clickable area in PDF coordinates:
// [x1, y1, x2, y2] where (x1, y1) is lower-left, (x2, y2) is upper-right.
//
// Example:
//
//	link := NewLinkAnnotation([4]float64{100, 690, 200, 710}, "https://google.com")
func NewLinkAnnotation(rect [4]float64, uri string) *LinkAnnotation {
	return &LinkAnnotation{
		Rect:        rect,
		URI:         uri,
		DestPage:    -1,
		IsInternal:  false,
		BorderWidth: 0,
	}
}

// NewInternalLinkAnnotation creates a new internal page link.
//
// The destPage parameter is 0-based (0 = first page, 1 = second page, etc.).
//
// Example:
//
//	link := NewInternalLinkAnnotation([4]float64{100, 690, 200, 710}, 2) // Link to page 3
func NewInternalLinkAnnotation(rect [4]float64, destPage int) *LinkAnnotation {
	return &LinkAnnotation{
		Rect:        rect,
		URI:         "",
		DestPage:    destPage,
		IsInternal:  true,
		BorderWidth: 0,
	}
}

// Validate checks if the link annotation is valid.
//
// Returns an error if:
// - Rectangle is invalid (x1 >= x2 or y1 >= y2)
// - External link has empty URI
// - Internal link has invalid destination page (< 0)
// - Border width is negative
func (a *LinkAnnotation) Validate() error {
	// Validate rectangle dimensions.
	if a.Rect[0] >= a.Rect[2] || a.Rect[1] >= a.Rect[3] {
		return ErrInvalidAnnotationRect
	}

	// Validate border width.
	if a.BorderWidth < 0 {
		return ErrInvalidBorderWidth
	}

	// Validate link target based on type.
	if a.IsInternal {
		if a.DestPage < 0 {
			return ErrInvalidDestPage
		}
	} else {
		if a.URI == "" {
			return ErrEmptyURI
		}
	}

	return nil
}

// TextAnnotation represents a sticky note annotation (/Subtype /Text).
//
// Text annotations appear as icons (sticky notes) on the page.
// They display pop-up text when clicked.
//
// Example:
//
//	note := NewTextAnnotation([4]float64{100, 700, 120, 720}, "This is a comment", "John Doe")
//	note.SetColor([3]float64{1, 1, 0}) // Yellow
type TextAnnotation struct {
	// Rect defines the icon location [x1, y1, x2, y2] in PDF coordinates.
	// Typically a small square (e.g., 20x20 points).
	Rect [4]float64

	// Contents is the text displayed in the pop-up.
	Contents string

	// Title is the author name (T field in PDF).
	Title string

	// Color is the annotation color in RGB (0.0 to 1.0 range).
	// Default: [1, 1, 0] (yellow).
	Color [3]float64

	// Open indicates if the pop-up should be open by default.
	Open bool
}

// NewTextAnnotation creates a new text (sticky note) annotation.
//
// Example:
//
//	note := NewTextAnnotation([4]float64{100, 700, 120, 720}, "Important!", "Alice")
func NewTextAnnotation(rect [4]float64, contents, title string) *TextAnnotation {
	return &TextAnnotation{
		Rect:     rect,
		Contents: contents,
		Title:    title,
		Color:    [3]float64{1, 1, 0}, // Yellow
		Open:     false,
	}
}

// SetColor sets the annotation color.
func (a *TextAnnotation) SetColor(color [3]float64) {
	a.Color = color
}

// SetOpen sets whether the pop-up is open by default.
func (a *TextAnnotation) SetOpen(open bool) {
	a.Open = open
}

// Validate checks if the text annotation is valid.
func (a *TextAnnotation) Validate() error {
	if a.Rect[0] >= a.Rect[2] || a.Rect[1] >= a.Rect[3] {
		return ErrInvalidAnnotationRect
	}
	if !isValidColor(a.Color) {
		return ErrInvalidColor
	}
	return nil
}

// MarkupAnnotation represents a markup annotation (highlight, underline, strikeout).
//
// Markup annotations highlight, underline, or strike through text.
// They use QuadPoints to define the area (4 points per quadrilateral).
//
// Example:
//
//	// Highlight from (100, 650) to (300, 670)
//	highlight := NewMarkupAnnotation(
//	    AnnotationTypeHighlight,
//	    [4]float64{100, 650, 300, 670},
//	    [][8]float64{{100, 670, 300, 670, 100, 650, 300, 650}},
//	)
type MarkupAnnotation struct {
	// Type is the markup type (Highlight, Underline, StrikeOut).
	Type AnnotationType

	// Rect defines the bounding box [x1, y1, x2, y2] in PDF coordinates.
	Rect [4]float64

	// QuadPoints defines the area to mark up.
	// Each quadrilateral is [x1, y1, x2, y2, x3, y3, x4, y4].
	// Points go: top-left, top-right, bottom-left, bottom-right.
	QuadPoints [][8]float64

	// Color is the annotation color in RGB (0.0 to 1.0 range).
	Color [3]float64

	// Title is the author name (T field in PDF).
	Title string

	// Contents is the optional note text.
	Contents string
}

// NewMarkupAnnotation creates a new markup annotation.
//
// The quadPoints should be in PDF coordinates (top-left, top-right, bottom-left, bottom-right).
//
// Example:
//
//	highlight := NewMarkupAnnotation(
//	    AnnotationTypeHighlight,
//	    [4]float64{100, 650, 300, 670},
//	    [][8]float64{{100, 670, 300, 670, 100, 650, 300, 650}},
//	)
func NewMarkupAnnotation(annotType AnnotationType, rect [4]float64, quadPoints [][8]float64) *MarkupAnnotation {
	return &MarkupAnnotation{
		Type:       annotType,
		Rect:       rect,
		QuadPoints: quadPoints,
		Color:      [3]float64{1, 1, 0}, // Yellow default
		Title:      "",
		Contents:   "",
	}
}

// SetColor sets the markup color.
func (a *MarkupAnnotation) SetColor(color [3]float64) {
	a.Color = color
}

// SetAuthor sets the author name.
func (a *MarkupAnnotation) SetAuthor(author string) {
	a.Title = author
}

// SetContents sets the note text.
func (a *MarkupAnnotation) SetContents(contents string) {
	a.Contents = contents
}

// Validate checks if the markup annotation is valid.
func (a *MarkupAnnotation) Validate() error {
	if a.Rect[0] >= a.Rect[2] || a.Rect[1] >= a.Rect[3] {
		return ErrInvalidAnnotationRect
	}
	if !isValidColor(a.Color) {
		return ErrInvalidColor
	}
	if len(a.QuadPoints) == 0 {
		return ErrMissingQuadPoints
	}
	return nil
}

// StampAnnotation represents a rubber stamp annotation (/Subtype /Stamp).
//
// Stamp annotations display predefined stamps like "Approved", "Draft", etc.
//
// Example:
//
//	stamp := NewStampAnnotation([4]float64{300, 700, 400, 750}, "Approved")
//	stamp.SetColor([3]float64{0, 1, 0}) // Green
type StampAnnotation struct {
	// Rect defines the stamp location [x1, y1, x2, y2] in PDF coordinates.
	Rect [4]float64

	// Name is the stamp name (e.g., "Approved", "Draft", "Confidential").
	Name string

	// Color is the annotation color in RGB (0.0 to 1.0 range).
	Color [3]float64

	// Title is the author name (T field in PDF).
	Title string

	// Contents is the optional note text.
	Contents string
}

// StampName represents predefined stamp names.
type StampName string

const (
	// StampApproved represents an "Approved" stamp.
	StampApproved StampName = "Approved"
	// StampNotApproved represents a "Not Approved" stamp.
	StampNotApproved StampName = "NotApproved"
	// StampDraft represents a "Draft" stamp.
	StampDraft StampName = "Draft"
	// StampFinal represents a "Final" stamp.
	StampFinal StampName = "Final"
	// StampConfidential represents a "Confidential" stamp.
	StampConfidential StampName = "Confidential"
	// StampForComment represents a "For Comment" stamp.
	StampForComment StampName = "ForComment"
	// StampForPublicRelease represents a "For Public Release" stamp.
	StampForPublicRelease StampName = "ForPublicRelease"
	// StampAsIs represents an "As Is" stamp.
	StampAsIs StampName = "AsIs"
	// StampDepartmental represents a "Departmental" stamp.
	StampDepartmental StampName = "Departmental"
	// StampExperimental represents an "Experimental" stamp.
	StampExperimental StampName = "Experimental"
	// StampExpired represents an "Expired" stamp.
	StampExpired StampName = "Expired"
	// StampNotForPublicRelease represents a "Not For Public Release" stamp.
	StampNotForPublicRelease StampName = "NotForPublicRelease"
)

// NewStampAnnotation creates a new stamp annotation.
//
// Example:
//
//	stamp := NewStampAnnotation([4]float64{300, 700, 400, 750}, StampApproved)
func NewStampAnnotation(rect [4]float64, name StampName) *StampAnnotation {
	return &StampAnnotation{
		Rect:     rect,
		Name:     string(name),
		Color:    [3]float64{1, 0, 0}, // Red default
		Title:    "",
		Contents: "",
	}
}

// SetColor sets the stamp color.
func (a *StampAnnotation) SetColor(color [3]float64) {
	a.Color = color
}

// SetAuthor sets the author name.
func (a *StampAnnotation) SetAuthor(author string) {
	a.Title = author
}

// SetContents sets the note text.
func (a *StampAnnotation) SetContents(contents string) {
	a.Contents = contents
}

// Validate checks if the stamp annotation is valid.
func (a *StampAnnotation) Validate() error {
	if a.Rect[0] >= a.Rect[2] || a.Rect[1] >= a.Rect[3] {
		return ErrInvalidAnnotationRect
	}
	if !isValidColor(a.Color) {
		return ErrInvalidColor
	}
	if a.Name == "" {
		return ErrMissingStampName
	}
	return nil
}

// isValidColor checks if all color components are in range [0, 1].
func isValidColor(c [3]float64) bool {
	for i := 0; i < 3; i++ {
		if c[i] < 0 || c[i] > 1 {
			return false
		}
	}
	return true
}

// Annotation errors.
var (
	// ErrInvalidAnnotationRect is returned when annotation rect is invalid.
	ErrInvalidAnnotationRect = errors.New("annotation rectangle must have x1 < x2 and y1 < y2")

	// ErrInvalidBorderWidth is returned when border width is negative.
	ErrInvalidBorderWidth = errors.New("border width must be non-negative")

	// ErrInvalidDestPage is returned when internal link destination is invalid.
	ErrInvalidDestPage = errors.New("destination page must be >= 0")

	// ErrEmptyURI is returned when external link has no URI.
	ErrEmptyURI = errors.New("external link must have a URI")

	// ErrInvalidColor is returned when color components are out of range.
	ErrInvalidColor = errors.New("color components must be in range [0.0, 1.0]")

	// ErrMissingQuadPoints is returned when markup annotation has no QuadPoints.
	ErrMissingQuadPoints = errors.New("markup annotation must have QuadPoints")

	// ErrMissingStampName is returned when stamp annotation has no name.
	ErrMissingStampName = errors.New("stamp annotation must have a name")
)
