package creator

import "errors"

// Block represents a rectangular area for drawing content.
//
// Blocks are used as containers for header/footer content and provide
// a bounded region where Drawable elements can be placed. Content is
// positioned relative to the block's origin (top-left corner).
//
// Example:
//
//	block := NewBlock(500, 50)
//	block.SetMargins(Margins{Left: 10, Right: 10})
//	p := NewParagraph("Header Text")
//	block.Draw(p)
type Block struct {
	width     float64
	height    float64
	drawables []DrawablePosition
	margins   Margins
	currentX  float64
	currentY  float64
}

// DrawablePosition stores a drawable with its position within the block.
type DrawablePosition struct {
	Drawable Drawable
	X        float64
	Y        float64
}

// NewBlock creates a new block with the specified dimensions.
//
// Parameters:
//   - width: Block width in points
//   - height: Block height in points
//
// Example:
//
//	block := NewBlock(500, 50)  // 500pt wide, 50pt tall
func NewBlock(width, height float64) *Block {
	return &Block{
		width:     width,
		height:    height,
		drawables: make([]DrawablePosition, 0, 4),
		margins:   Margins{},
		currentX:  0,
		currentY:  0,
	}
}

// Width returns the block width in points.
func (b *Block) Width() float64 {
	return b.width
}

// Height returns the block height in points.
func (b *Block) Height() float64 {
	return b.height
}

// SetWidth sets the block width in points.
func (b *Block) SetWidth(width float64) {
	b.width = width
}

// SetHeight sets the block height in points.
func (b *Block) SetHeight(height float64) {
	b.height = height
}

// Margins returns the block margins.
func (b *Block) Margins() Margins {
	return b.margins
}

// SetMargins sets the block margins.
//
// Margins define the padding inside the block, reducing the available
// drawing area.
func (b *Block) SetMargins(m Margins) {
	b.margins = m
}

// ContentWidth returns the usable width inside the block (width minus margins).
func (b *Block) ContentWidth() float64 {
	return b.width - b.margins.Left - b.margins.Right
}

// ContentHeight returns the usable height inside the block (height minus margins).
func (b *Block) ContentHeight() float64 {
	return b.height - b.margins.Top - b.margins.Bottom
}

// Draw adds a drawable element to the block at the current position.
//
// The drawable is positioned at the block's current cursor position,
// which starts at (0, 0) relative to the block's content area.
//
// Returns an error if the drawable is nil.
func (b *Block) Draw(d Drawable) error {
	if d == nil {
		return errors.New("drawable cannot be nil")
	}
	b.drawables = append(b.drawables, DrawablePosition{
		Drawable: d,
		X:        b.currentX + b.margins.Left,
		Y:        b.currentY + b.margins.Top,
	})
	return nil
}

// DrawAt adds a drawable element at a specific position within the block.
//
// Coordinates are relative to the block's top-left corner (before margins).
//
// Parameters:
//   - d: The drawable element to add
//   - x: Horizontal position from left edge of block
//   - y: Vertical position from top edge of block
//
// Returns an error if the drawable is nil.
func (b *Block) DrawAt(d Drawable, x, y float64) error {
	if d == nil {
		return errors.New("drawable cannot be nil")
	}
	b.drawables = append(b.drawables, DrawablePosition{
		Drawable: d,
		X:        x,
		Y:        y,
	})
	return nil
}

// GetDrawables returns all drawable elements with their positions.
//
// This is used internally to render the block contents to a page.
func (b *Block) GetDrawables() []DrawablePosition {
	return b.drawables
}

// Clear removes all drawables from the block.
func (b *Block) Clear() {
	b.drawables = b.drawables[:0]
	b.currentX = 0
	b.currentY = 0
}

// SetCursor sets the current drawing position within the block.
//
// This affects subsequent Draw() calls.
func (b *Block) SetCursor(x, y float64) {
	b.currentX = x
	b.currentY = y
}

// MoveCursor moves the current drawing position by the specified delta.
func (b *Block) MoveCursor(dx, dy float64) {
	b.currentX += dx
	b.currentY += dy
}

// GetLayoutContext creates a LayoutContext for this block.
//
// This allows Drawable elements to query the available space and position
// themselves correctly within the block.
func (b *Block) GetLayoutContext() *LayoutContext {
	return &LayoutContext{
		PageWidth:  b.width,
		PageHeight: b.height,
		Margins:    b.margins,
		CursorX:    b.margins.Left,
		CursorY:    0,
	}
}
