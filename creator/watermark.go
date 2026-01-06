package creator

import (
	"errors"
	"math"
)

// WatermarkPosition defines the position of a watermark on a page.
type WatermarkPosition int

const (
	// WatermarkCenter positions the watermark at the page center.
	WatermarkCenter WatermarkPosition = iota

	// WatermarkTopLeft positions the watermark at the top-left corner.
	WatermarkTopLeft

	// WatermarkTopRight positions the watermark at the top-right corner.
	WatermarkTopRight

	// WatermarkBottomLeft positions the watermark at the bottom-left corner.
	WatermarkBottomLeft

	// WatermarkBottomRight positions the watermark at the bottom-right corner.
	WatermarkBottomRight
)

// TextWatermark represents a text watermark to be applied to PDF pages.
//
// A watermark is semi-transparent text (typically "CONFIDENTIAL", "DRAFT", etc.)
// rendered on a page, often rotated diagonally.
//
// Example:
//
//	wm := creator.NewTextWatermark("CONFIDENTIAL")
//	wm.SetFont(creator.HelveticaBold, 72)
//	wm.SetColor(creator.Gray)
//	wm.SetOpacity(0.3)
//	wm.SetRotation(45)
//	page.DrawWatermark(wm)
type TextWatermark struct {
	text     string
	font     FontName
	fontSize float64
	color    Color
	opacity  float64
	rotation float64
	position WatermarkPosition
}

// NewTextWatermark creates a new text watermark with default settings.
//
// Default settings:
//   - Font: HelveticaBold
//   - Font size: 48 points
//   - Color: Gray (0.5, 0.5, 0.5)
//   - Opacity: 0.5 (50% transparent)
//   - Rotation: 45 degrees (diagonal)
//   - Position: Center
//
// Example:
//
//	wm := creator.NewTextWatermark("DRAFT")
func NewTextWatermark(text string) *TextWatermark {
	return &TextWatermark{
		text:     text,
		font:     HelveticaBold,
		fontSize: 48,
		color:    Gray,
		opacity:  0.5,
		rotation: 45,
		position: WatermarkCenter,
	}
}

// SetFont sets the watermark font and size.
//
// Parameters:
//   - font: Font name (one of the Standard 14 fonts)
//   - size: Font size in points (must be > 0)
//
// Example:
//
//	wm.SetFont(creator.HelveticaBold, 72)
func (w *TextWatermark) SetFont(font FontName, size float64) error {
	if size <= 0 {
		return errors.New("font size must be positive")
	}
	w.font = font
	w.fontSize = size
	return nil
}

// SetColor sets the watermark text color.
//
// Color components must be in the range [0.0, 1.0].
//
// Example:
//
//	wm.SetColor(creator.Red)
//	wm.SetColor(creator.Color{R: 0.8, G: 0.0, B: 0.0})
func (w *TextWatermark) SetColor(color Color) error {
	if err := validateColor(color); err != nil {
		return err
	}
	w.color = color
	return nil
}

// SetOpacity sets the watermark transparency.
//
// Opacity must be in the range [0.0, 1.0]:
//   - 0.0 = fully transparent (invisible)
//   - 1.0 = fully opaque (no transparency)
//   - 0.3-0.5 = typical watermark range
//
// Example:
//
//	wm.SetOpacity(0.3) // 30% opaque, 70% transparent
func (w *TextWatermark) SetOpacity(opacity float64) error {
	if opacity < 0 || opacity > 1 {
		return errors.New("opacity must be in range [0.0, 1.0]")
	}
	w.opacity = opacity
	return nil
}

// SetRotation sets the watermark rotation angle in degrees (clockwise).
//
// Common values:
//   - 0 = horizontal
//   - 45 = diagonal (typical for "CONFIDENTIAL")
//   - 90 = vertical
//
// Example:
//
//	wm.SetRotation(45) // Diagonal
func (w *TextWatermark) SetRotation(degrees float64) error {
	w.rotation = degrees
	return nil
}

// SetPosition sets the watermark position on the page.
//
// Example:
//
//	wm.SetPosition(creator.WatermarkCenter)
//	wm.SetPosition(creator.WatermarkTopRight)
func (w *TextWatermark) SetPosition(position WatermarkPosition) error {
	w.position = position
	return nil
}

// Text returns the watermark text.
func (w *TextWatermark) Text() string {
	return w.text
}

// Font returns the watermark font.
func (w *TextWatermark) Font() FontName {
	return w.font
}

// FontSize returns the watermark font size.
func (w *TextWatermark) FontSize() float64 {
	return w.fontSize
}

// Color returns the watermark color.
func (w *TextWatermark) Color() Color {
	return w.color
}

// Opacity returns the watermark opacity.
func (w *TextWatermark) Opacity() float64 {
	return w.opacity
}

// Rotation returns the watermark rotation in degrees.
func (w *TextWatermark) Rotation() float64 {
	return w.rotation
}

// Position returns the watermark position.
func (w *TextWatermark) Position() WatermarkPosition {
	return w.position
}

// DrawWatermark applies a text watermark to the page.
//
// The watermark is rendered as semi-transparent text positioned
// according to the watermark's settings.
//
// Example:
//
//	wm := creator.NewTextWatermark("CONFIDENTIAL")
//	wm.SetOpacity(0.3)
//	page.DrawWatermark(wm)
func (p *Page) DrawWatermark(wm *TextWatermark) error {
	if wm == nil {
		return errors.New("watermark cannot be nil")
	}

	// Calculate position.
	x, y := calculateWatermarkPosition(p, wm)

	// Store watermark as a graphics operation.
	// We use a special operation type for watermarks to handle
	// the opacity and rotation transformation in the content stream writer.
	p.graphicsOps = append(p.graphicsOps, GraphicsOperation{
		Type:        GraphicsOpWatermark,
		X:           x,
		Y:           y,
		WatermarkOp: wm,
	})

	return nil
}

// calculateWatermarkPosition calculates the watermark position based on
// the page dimensions and watermark position setting.
func calculateWatermarkPosition(p *Page, wm *TextWatermark) (float64, float64) {
	pageWidth := p.Width()
	pageHeight := p.Height()

	// Measure text width for positioning.
	textWidth := measureTextWidth(string(wm.font), wm.text, wm.fontSize)

	// Calculate position based on setting.
	var x, y float64

	switch wm.position {
	case WatermarkCenter:
		// Center of page.
		x = pageWidth / 2
		y = pageHeight / 2

	case WatermarkTopLeft:
		// Top-left corner with padding.
		padding := wm.fontSize * 0.5
		x = padding
		y = pageHeight - padding

	case WatermarkTopRight:
		// Top-right corner with padding.
		padding := wm.fontSize * 0.5
		x = pageWidth - padding - textWidth
		y = pageHeight - padding

	case WatermarkBottomLeft:
		// Bottom-left corner with padding.
		padding := wm.fontSize * 0.5
		x = padding
		y = padding + wm.fontSize

	case WatermarkBottomRight:
		// Bottom-right corner with padding.
		padding := wm.fontSize * 0.5
		x = pageWidth - padding - textWidth
		y = padding + wm.fontSize

	default:
		// Default to center.
		x = pageWidth / 2
		y = pageHeight / 2
	}

	return x, y
}

// rotationMatrix calculates the transformation matrix for rotation.
//
// This rotates the coordinate system around the origin point (x, y)
// by the specified angle in degrees (clockwise in PDF coordinate system).
//
// Returns: [a, b, c, d, e, f] for PDF transformation matrix:
//
//	| a  b  0 |
//	| c  d  0 |
//	| e  f  1 |
func rotationMatrix(x, y, degrees float64) [6]float64 {
	// Convert degrees to radians.
	radians := degrees * math.Pi / 180.0

	cos := math.Cos(radians)
	sin := math.Sin(radians)

	// Transformation matrix for rotation around point (x, y):
	// 1. Translate to origin: (-x, -y)
	// 2. Rotate: (cos, sin, -sin, cos, 0, 0)
	// 3. Translate back: (x, y)
	//
	// Combined matrix:
	// a = cos, b = sin, c = -sin, d = cos
	// e = x - x*cos + y*sin
	// f = y - x*sin - y*cos

	return [6]float64{
		cos,               // a
		sin,               // b
		-sin,              // c
		cos,               // d
		x - x*cos + y*sin, // e
		y - x*sin - y*cos, // f
	}
}
