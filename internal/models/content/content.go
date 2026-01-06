// Package content defines the domain model for PDF page content.
//
// This package contains the Content interface and related types that represent
// elements that can be placed on PDF pages (text, images, shapes, etc.).
package content

import (
	"io"

	"github.com/coregx/gxpdf/internal/models/types"
)

// Content represents any element that can be placed on a PDF page.
//
// This is the core interface in the content domain. All renderable elements
// (text, images, shapes, tables, forms, etc.) must implement this interface.
//
// Design principles:
// - Each content element knows how to render itself (Rich Domain Model)
// - Content is position-aware (has bounds)
// - Content is validatable
// - Content provides metadata for debugging/logging
//
// Example implementations:
// - TextContent: Rendered text with font, size, position
// - ImageContent: Embedded images with scaling
// - PathContent: Vector graphics (lines, rectangles, etc.)
// - TableContent: Complex table layouts
type Content interface {
	// Render writes the PDF content stream operators for this element.
	//
	// The implementation should write valid PDF operators to the writer.
	// For example, text might write: BT, Tf, Tm, Tj, ET
	//
	// Returns an error if rendering fails.
	Render(w io.Writer) error

	// Bounds returns the bounding rectangle of the content in page coordinates.
	//
	// The rectangle defines the area occupied by this content element.
	// This is useful for:
	// - Layout calculations
	// - Collision detection
	// - Content positioning
	Bounds() types.Rectangle

	// Validate checks if the content is valid and ready to be rendered.
	//
	// Returns an error if:
	// - Required fields are missing (e.g., text without font)
	// - Values are out of valid range
	// - Content is malformed
	Validate() error

	// Type returns the content type identifier.
	//
	// This is primarily for debugging, logging, and type discrimination.
	// Examples: "text", "image", "path", "table"
	Type() ContentType
}

// ContentType identifies the type of content.
//
// This is used for:
// - Debugging and logging
// - Type assertions when needed
// - Content statistics
type ContentType string

const (
	// ContentTypeText represents text content.
	ContentTypeText ContentType = "text"

	// ContentTypeImage represents image content.
	ContentTypeImage ContentType = "image"

	// ContentTypePath represents vector graphics (paths, lines, shapes).
	ContentTypePath ContentType = "path"

	// ContentTypeTable represents table content.
	ContentTypeTable ContentType = "table"

	// ContentTypeForm represents form field content.
	ContentTypeForm ContentType = "form"

	// ContentTypeAnnotation represents annotations and markup.
	ContentTypeAnnotation ContentType = "annotation"
)

// String returns the string representation of the content type.
func (ct ContentType) String() string {
	return string(ct)
}

// IsValid checks if the content type is valid.
func (ct ContentType) IsValid() bool {
	switch ct {
	case ContentTypeText, ContentTypeImage, ContentTypePath, ContentTypeTable, ContentTypeForm, ContentTypeAnnotation:
		return true
	default:
		return false
	}
}
