package document

import (
	"errors"
	"fmt"

	"github.com/coregx/gxpdf/internal/models/content"
	"github.com/coregx/gxpdf/internal/models/types"
)

// Page represents a single page in a PDF document.
//
// Pages contain content elements and have properties like size, rotation, etc.
// This is an entity within the Document aggregate.
//
// Example:
//
//	page := document.NewPage(0, document.A4)
//	page.SetRotation(90)  // Landscape
//	width := page.Width()
type Page struct {
	// Identity
	number int // Page number (0-based)

	// Properties
	mediaBox types.Rectangle  // Page dimensions
	cropBox  *types.Rectangle // Visible area (optional)
	rotation int              // Rotation angle (0, 90, 180, 270)

	// Content
	contents []content.Content // Content elements on the page

	// Annotations (different types)
	linkAnnotations   []*LinkAnnotation   // Link annotations
	textAnnotations   []*TextAnnotation   // Text (sticky note) annotations
	markupAnnotations []*MarkupAnnotation // Markup annotations (highlight, underline, strikeout)
	stampAnnotations  []*StampAnnotation  // Stamp annotations

	// Form fields (interactive form widgets)
	formFields []*FormField // Form field annotations
}

// NewPage creates a new page with the specified size.
//
// The page number is 0-based and should correspond to its position
// in the document.
//
// Example:
//
//	page := document.NewPage(0, document.A4)
func NewPage(number int, size PageSize) *Page {
	return &Page{
		number:            number,
		mediaBox:          size.ToRectangle(),
		rotation:          0,
		contents:          make([]content.Content, 0),
		linkAnnotations:   make([]*LinkAnnotation, 0),
		textAnnotations:   make([]*TextAnnotation, 0),
		markupAnnotations: make([]*MarkupAnnotation, 0),
		stampAnnotations:  make([]*StampAnnotation, 0),
		formFields:        make([]*FormField, 0),
	}
}

// Number returns the page number (0-based).
func (p *Page) Number() int {
	return p.number
}

// MediaBox returns the page's media box (page dimensions).
func (p *Page) MediaBox() types.Rectangle {
	return p.mediaBox
}

// CropBox returns the page's crop box (visible area).
//
// Returns nil if no crop box is set (media box is used).
func (p *Page) CropBox() *types.Rectangle {
	return p.cropBox
}

// SetCropBox sets the crop box (visible area of the page).
//
// The crop box must be within or equal to the media box.
func (p *Page) SetCropBox(box types.Rectangle) error {
	// Get coordinates from rectangles
	boxLLX, boxLLY := box.LowerLeft()
	boxURX, boxURY := box.UpperRight()
	mediaLLX, mediaLLY := p.mediaBox.LowerLeft()
	mediaURX, mediaURY := p.mediaBox.UpperRight()

	// Validate crop box is within media box
	if boxLLX < mediaLLX || boxLLY < mediaLLY || boxURX > mediaURX || boxURY > mediaURY {
		return ErrCropBoxOutOfBounds
	}

	p.cropBox = &box
	return nil
}

// SetRotation sets the page rotation (0, 90, 180, 270 degrees).
//
// Rotation is applied clockwise.
//
// Returns an error if the rotation is not one of the valid values.
func (p *Page) SetRotation(degrees int) error {
	if degrees != 0 && degrees != 90 && degrees != 180 && degrees != 270 {
		return fmt.Errorf("%w: got %d, want 0, 90, 180, or 270", ErrInvalidRotation, degrees)
	}
	p.rotation = degrees
	return nil
}

// Rotation returns the current page rotation in degrees.
func (p *Page) Rotation() int {
	return p.rotation
}

// Width returns the page width in points.
//
// If the page is rotated 90 or 270 degrees, width and height are swapped.
func (p *Page) Width() float64 {
	if p.rotation == 90 || p.rotation == 270 {
		return p.mediaBox.Height()
	}
	return p.mediaBox.Width()
}

// Height returns the page height in points.
//
// If the page is rotated 90 or 270 degrees, width and height are swapped.
func (p *Page) Height() float64 {
	if p.rotation == 90 || p.rotation == 270 {
		return p.mediaBox.Width()
	}
	return p.mediaBox.Height()
}

// AddContent adds a content element to the page.
//
// Returns an error if:
// - Content is nil
// - Content validation fails
//
// Example:
//
//	text := &TextContent{...}
//	err := page.AddContent(text)
func (p *Page) AddContent(c content.Content) error {
	if c == nil {
		return ErrNilContent
	}

	// Validate content before adding
	if err := c.Validate(); err != nil {
		return fmt.Errorf("content validation failed: %w", err)
	}

	p.contents = append(p.contents, c)
	return nil
}

// Contents returns all content elements on the page.
//
// The returned slice is a copy to prevent external modifications.
func (p *Page) Contents() []content.Content {
	result := make([]content.Content, len(p.contents))
	copy(result, p.contents)
	return result
}

// ContentCount returns the number of content elements on the page.
func (p *Page) ContentCount() int {
	return len(p.contents)
}

// ClearContent removes all content from the page.
func (p *Page) ClearContent() {
	p.contents = make([]content.Content, 0)
}

// AddAnnotation adds a link annotation to the page.
//
// Deprecated: Use AddLinkAnnotation instead.
// This method is kept for backward compatibility.
//
// Returns an error if:
// - Annotation is nil
// - Annotation validation fails
//
// Example:
//
//	link := NewLinkAnnotation([4]float64{100, 690, 200, 710}, "https://example.com")
//	err := page.AddAnnotation(link)
func (p *Page) AddAnnotation(a *LinkAnnotation) error {
	return p.AddLinkAnnotation(a)
}

// AddLinkAnnotation adds a link annotation to the page.
//
// Returns an error if:
// - Annotation is nil
// - Annotation validation fails
//
// Example:
//
//	link := NewLinkAnnotation([4]float64{100, 690, 200, 710}, "https://example.com")
//	err := page.AddLinkAnnotation(link)
func (p *Page) AddLinkAnnotation(a *LinkAnnotation) error {
	if a == nil {
		return ErrNilAnnotation
	}

	// Validate annotation before adding.
	if err := a.Validate(); err != nil {
		return fmt.Errorf("link annotation validation failed: %w", err)
	}

	p.linkAnnotations = append(p.linkAnnotations, a)
	return nil
}

// AddTextAnnotation adds a text (sticky note) annotation to the page.
//
// Returns an error if:
// - Annotation is nil
// - Annotation validation fails
//
// Example:
//
//	note := NewTextAnnotation([4]float64{100, 700, 120, 720}, "Important!", "Alice")
//	err := page.AddTextAnnotation(note)
func (p *Page) AddTextAnnotation(a *TextAnnotation) error {
	if a == nil {
		return ErrNilAnnotation
	}

	if err := a.Validate(); err != nil {
		return fmt.Errorf("text annotation validation failed: %w", err)
	}

	p.textAnnotations = append(p.textAnnotations, a)
	return nil
}

// AddMarkupAnnotation adds a markup annotation (highlight, underline, strikeout) to the page.
//
// Returns an error if:
// - Annotation is nil
// - Annotation validation fails
//
// Example:
//
//	highlight := NewMarkupAnnotation(
//	    AnnotationTypeHighlight,
//	    [4]float64{100, 650, 300, 670},
//	    [][4]float64{{100, 670, 300, 670, 100, 650, 300, 650}},
//	)
//	err := page.AddMarkupAnnotation(highlight)
func (p *Page) AddMarkupAnnotation(a *MarkupAnnotation) error {
	if a == nil {
		return ErrNilAnnotation
	}

	if err := a.Validate(); err != nil {
		return fmt.Errorf("markup annotation validation failed: %w", err)
	}

	p.markupAnnotations = append(p.markupAnnotations, a)
	return nil
}

// AddStampAnnotation adds a stamp annotation to the page.
//
// Returns an error if:
// - Annotation is nil
// - Annotation validation fails
//
// Example:
//
//	stamp := NewStampAnnotation([4]float64{300, 700, 400, 750}, StampApproved)
//	err := page.AddStampAnnotation(stamp)
func (p *Page) AddStampAnnotation(a *StampAnnotation) error {
	if a == nil {
		return ErrNilAnnotation
	}

	if err := a.Validate(); err != nil {
		return fmt.Errorf("stamp annotation validation failed: %w", err)
	}

	p.stampAnnotations = append(p.stampAnnotations, a)
	return nil
}

// AddFormField adds a form field annotation to the page.
//
// Returns an error if:
// - Form field is nil
// - Form field validation fails
//
// Example:
//
//	field := NewFormField("Tx", "username", [4]float64{100, 700, 300, 720})
//	err := page.AddFormField(field)
func (p *Page) AddFormField(f *FormField) error {
	if f == nil {
		return ErrNilFormField
	}

	if err := f.Validate(); err != nil {
		return fmt.Errorf("form field validation failed: %w", err)
	}

	p.formFields = append(p.formFields, f)
	return nil
}

// Annotations returns all link annotations on the page.
//
// Deprecated: Use LinkAnnotations instead.
// This method is kept for backward compatibility.
//
// The returned slice is a copy to prevent external modifications.
func (p *Page) Annotations() []*LinkAnnotation {
	return p.LinkAnnotations()
}

// LinkAnnotations returns all link annotations on the page.
//
// The returned slice is a copy to prevent external modifications.
func (p *Page) LinkAnnotations() []*LinkAnnotation {
	result := make([]*LinkAnnotation, len(p.linkAnnotations))
	copy(result, p.linkAnnotations)
	return result
}

// TextAnnotations returns all text annotations on the page.
//
// The returned slice is a copy to prevent external modifications.
func (p *Page) TextAnnotations() []*TextAnnotation {
	result := make([]*TextAnnotation, len(p.textAnnotations))
	copy(result, p.textAnnotations)
	return result
}

// MarkupAnnotations returns all markup annotations on the page.
//
// The returned slice is a copy to prevent external modifications.
func (p *Page) MarkupAnnotations() []*MarkupAnnotation {
	result := make([]*MarkupAnnotation, len(p.markupAnnotations))
	copy(result, p.markupAnnotations)
	return result
}

// StampAnnotations returns all stamp annotations on the page.
//
// The returned slice is a copy to prevent external modifications.
func (p *Page) StampAnnotations() []*StampAnnotation {
	result := make([]*StampAnnotation, len(p.stampAnnotations))
	copy(result, p.stampAnnotations)
	return result
}

// FormFields returns all form field annotations on the page.
//
// The returned slice is a copy to prevent external modifications.
func (p *Page) FormFields() []*FormField {
	result := make([]*FormField, len(p.formFields))
	copy(result, p.formFields)
	return result
}

// AnnotationCount returns the total number of annotations on the page.
func (p *Page) AnnotationCount() int {
	return len(p.linkAnnotations) + len(p.textAnnotations) +
		len(p.markupAnnotations) + len(p.stampAnnotations) + len(p.formFields)
}

// ClearAnnotations removes all annotations from the page.
func (p *Page) ClearAnnotations() {
	p.linkAnnotations = make([]*LinkAnnotation, 0)
	p.textAnnotations = make([]*TextAnnotation, 0)
	p.markupAnnotations = make([]*MarkupAnnotation, 0)
	p.stampAnnotations = make([]*StampAnnotation, 0)
	p.formFields = make([]*FormField, 0)
}

// Validate checks page consistency.
//
// Returns an error if:
// - Crop box is out of bounds
// - Rotation is invalid
//
// Note: Page dimensions are always valid because Rectangle value objects
// enforce validity at construction time.
func (p *Page) Validate() error {
	// Note: No need to check media box dimensions - Rectangle enforces validity

	// Check crop box if set
	if p.cropBox != nil {
		cropLLX, cropLLY := p.cropBox.LowerLeft()
		cropURX, cropURY := p.cropBox.UpperRight()
		mediaLLX, mediaLLY := p.mediaBox.LowerLeft()
		mediaURX, mediaURY := p.mediaBox.UpperRight()

		if cropLLX < mediaLLX || cropLLY < mediaLLY || cropURX > mediaURX || cropURY > mediaURY {
			return ErrCropBoxOutOfBounds
		}
	}

	// Check rotation
	if p.rotation != 0 && p.rotation != 90 && p.rotation != 180 && p.rotation != 270 {
		return fmt.Errorf("%w: %d", ErrInvalidRotation, p.rotation)
	}

	// Validate all content elements
	for i, c := range p.contents {
		if c == nil {
			return fmt.Errorf("content at index %d is nil", i)
		}
		if err := c.Validate(); err != nil {
			return fmt.Errorf("content at index %d validation failed: %w", i, err)
		}
	}

	// Validate all annotations.
	for i, a := range p.linkAnnotations {
		if a == nil {
			return fmt.Errorf("link annotation at index %d is nil", i)
		}
		if err := a.Validate(); err != nil {
			return fmt.Errorf("link annotation at index %d validation failed: %w", i, err)
		}
	}

	for i, a := range p.textAnnotations {
		if a == nil {
			return fmt.Errorf("text annotation at index %d is nil", i)
		}
		if err := a.Validate(); err != nil {
			return fmt.Errorf("text annotation at index %d validation failed: %w", i, err)
		}
	}

	for i, a := range p.markupAnnotations {
		if a == nil {
			return fmt.Errorf("markup annotation at index %d is nil", i)
		}
		if err := a.Validate(); err != nil {
			return fmt.Errorf("markup annotation at index %d validation failed: %w", i, err)
		}
	}

	for i, a := range p.stampAnnotations {
		if a == nil {
			return fmt.Errorf("stamp annotation at index %d is nil", i)
		}
		if err := a.Validate(); err != nil {
			return fmt.Errorf("stamp annotation at index %d validation failed: %w", i, err)
		}
	}

	for i, f := range p.formFields {
		if f == nil {
			return fmt.Errorf("form field at index %d is nil", i)
		}
		if err := f.Validate(); err != nil {
			return fmt.Errorf("form field at index %d validation failed: %w", i, err)
		}
	}

	return nil
}

// Domain errors
var (
	// ErrInvalidRotation is returned when rotation is not 0, 90, 180, or 270.
	ErrInvalidRotation = errors.New("rotation must be 0, 90, 180, or 270")

	// ErrInvalidPageSize is returned when page dimensions are invalid.
	ErrInvalidPageSize = errors.New("invalid page size: width and height must be positive")

	// ErrCropBoxOutOfBounds is returned when crop box exceeds media box.
	ErrCropBoxOutOfBounds = errors.New("crop box must be within media box bounds")

	// ErrNilContent is returned when trying to add nil content to a page.
	ErrNilContent = errors.New("content cannot be nil")

	// ErrNilAnnotation is returned when trying to add nil annotation to a page.
	ErrNilAnnotation = errors.New("annotation cannot be nil")

	// ErrNilFormField is returned when trying to add nil form field to a page.
	ErrNilFormField = errors.New("form field cannot be nil")
)
