package creator

// Alignment represents text alignment within a block.
type Alignment int

const (
	// AlignLeft aligns text to the left edge.
	AlignLeft Alignment = iota

	// AlignCenter centers text horizontally.
	AlignCenter

	// AlignRight aligns text to the right edge.
	AlignRight

	// AlignJustify stretches text to fill the full width.
	AlignJustify
)

// LayoutContext provides positioning information for layout operations.
//
// It tracks the current cursor position and available space within
// the page margins. The cursor Y is measured from the top of the
// content area (intuitive "document flow" model).
//
// Example:
//
//	ctx := page.GetLayoutContext()
//	paragraph.Draw(ctx, page)  // Draws at cursor position
//	// Cursor automatically advances after drawing
type LayoutContext struct {
	// Page dimensions (in points).
	PageWidth  float64
	PageHeight float64

	// Page margins (in points).
	Margins Margins

	// CursorX is the current horizontal position (from left edge).
	CursorX float64

	// CursorY is the current vertical position from top of content area.
	// 0 = top of content area (below top margin)
	// Increases downward.
	CursorY float64
}

// Drawable is an interface for elements that can be drawn on a page.
//
// Elements implementing this interface can be used with Page.Draw() for
// automatic layout and positioning.
type Drawable interface {
	// Draw renders the element on the page using the layout context.
	// The context's cursor position is updated after drawing.
	Draw(ctx *LayoutContext, page *Page) error

	// Height returns the pre-calculated height of the element.
	// This is used for page break detection and layout planning.
	Height(ctx *LayoutContext) float64
}

// AvailableWidth returns the width available for content (excluding margins).
func (ctx *LayoutContext) AvailableWidth() float64 {
	return ctx.PageWidth - ctx.Margins.Left - ctx.Margins.Right
}

// AvailableHeight returns the height remaining from cursor to bottom margin.
func (ctx *LayoutContext) AvailableHeight() float64 {
	contentHeight := ctx.PageHeight - ctx.Margins.Top - ctx.Margins.Bottom
	return contentHeight - ctx.CursorY
}

// ContentLeft returns the X coordinate of the left content edge.
func (ctx *LayoutContext) ContentLeft() float64 {
	return ctx.Margins.Left
}

// ContentRight returns the X coordinate of the right content edge.
func (ctx *LayoutContext) ContentRight() float64 {
	return ctx.PageWidth - ctx.Margins.Right
}

// ContentTop returns the Y coordinate (PDF coordinates) of the top content edge.
func (ctx *LayoutContext) ContentTop() float64 {
	return ctx.PageHeight - ctx.Margins.Top
}

// ContentBottom returns the Y coordinate (PDF coordinates) of the bottom content edge.
func (ctx *LayoutContext) ContentBottom() float64 {
	return ctx.Margins.Bottom
}

// CurrentPDFY converts the cursor Y (from top) to PDF Y coordinate (from bottom).
func (ctx *LayoutContext) CurrentPDFY() float64 {
	return ctx.ContentTop() - ctx.CursorY
}

// MoveCursor moves the cursor by the specified delta values.
//
// Positive dx moves right, positive dy moves down.
func (ctx *LayoutContext) MoveCursor(dx, dy float64) {
	ctx.CursorX += dx
	ctx.CursorY += dy
}

// SetCursor sets the cursor to specific coordinates.
//
// x is measured from the left edge of the page.
// y is measured from the top of the content area (below top margin).
func (ctx *LayoutContext) SetCursor(x, y float64) {
	ctx.CursorX = x
	ctx.CursorY = y
}

// NewLine moves the cursor to the start of the next line.
//
// The cursor X is reset to the left content edge.
// The cursor Y advances by the specified line height.
func (ctx *LayoutContext) NewLine(lineHeight float64) {
	ctx.CursorX = ctx.ContentLeft()
	ctx.CursorY += lineHeight
}

// ResetX moves the cursor X back to the left content edge.
func (ctx *LayoutContext) ResetX() {
	ctx.CursorX = ctx.ContentLeft()
}

// CanFit checks if an element of the given height fits in the remaining space.
func (ctx *LayoutContext) CanFit(height float64) bool {
	return ctx.AvailableHeight() >= height
}
