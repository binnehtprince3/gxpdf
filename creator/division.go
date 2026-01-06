package creator

// Border defines the border style for a division.
//
// Border controls the appearance of a division's outline, including
// width and color. Individual borders can be set per side.
//
// Example:
//
//	border := Border{Width: 1.0, Color: Black}
//	div.SetBorder(border)
type Border struct {
	// Width is the border line width in points.
	Width float64

	// Color is the border color (RGB, 0.0 to 1.0 range).
	Color Color
}

// Division is a container for grouping multiple Drawables.
//
// Division provides a box model with support for:
// - Background color
// - Borders (all sides or individual)
// - Padding (inner spacing)
// - Margins (outer spacing)
// - Width and minimum height control
//
// Example:
//
//	div := NewDivision()
//	div.SetBackground(White)
//	div.SetBorder(Border{Width: 1, Color: Gray})
//	div.SetPaddingAll(10)
//	div.Add(NewParagraph("Content"))
//	page.Draw(div)
type Division struct {
	// Content.
	drawables []Drawable

	// Styling.
	background *Color // nil = transparent

	// Borders (nil = no border).
	border       *Border // Default for all sides
	borderTop    *Border // Overrides border
	borderRight  *Border // Overrides border
	borderBottom *Border // Overrides border
	borderLeft   *Border // Overrides border

	// Spacing.
	padding Margins // Inner spacing
	margins Margins // Outer spacing

	// Dimensions.
	width     float64 // 0 = auto (use available width)
	minHeight float64 // Minimum height
}

// NewDivision creates a new empty division.
//
// The division starts with no content, transparent background,
// no borders, no padding, and no margins.
//
// Example:
//
//	div := NewDivision()
//	div.Add(NewParagraph("Hello"))
func NewDivision() *Division {
	return &Division{
		drawables: make([]Drawable, 0),
		padding:   Margins{},
		margins:   Margins{},
		width:     0, // Auto width
		minHeight: 0,
	}
}

// SetBackground sets the background color for the division.
//
// The background fills the entire division area (including padding,
// but excluding margins).
//
// Example:
//
//	div.SetBackground(White)
func (d *Division) SetBackground(c Color) *Division {
	d.background = &c
	return d
}

// SetBorder sets the border for all sides of the division.
//
// Individual borders (top, right, bottom, left) can override this
// if set separately.
//
// Example:
//
//	div.SetBorder(Border{Width: 2, Color: Black})
func (d *Division) SetBorder(b Border) *Division {
	d.border = &b
	return d
}

// SetBorderTop sets the top border, overriding the default border.
//
// Example:
//
//	div.SetBorderTop(Border{Width: 3, Color: Red})
func (d *Division) SetBorderTop(b Border) *Division {
	d.borderTop = &b
	return d
}

// SetBorderRight sets the right border, overriding the default border.
//
// Example:
//
//	div.SetBorderRight(Border{Width: 1, Color: Gray})
func (d *Division) SetBorderRight(b Border) *Division {
	d.borderRight = &b
	return d
}

// SetBorderBottom sets the bottom border, overriding the default border.
//
// Example:
//
//	div.SetBorderBottom(Border{Width: 3, Color: Red})
func (d *Division) SetBorderBottom(b Border) *Division {
	d.borderBottom = &b
	return d
}

// SetBorderLeft sets the left border, overriding the default border.
//
// Example:
//
//	div.SetBorderLeft(Border{Width: 4, Color: Blue})
func (d *Division) SetBorderLeft(b Border) *Division {
	d.borderLeft = &b
	return d
}

// SetPadding sets padding for all sides individually.
//
// Padding is the inner spacing between the border and content.
//
// Example:
//
//	div.SetPadding(10, 15, 10, 15) // top, right, bottom, left
func (d *Division) SetPadding(top, right, bottom, left float64) *Division {
	d.padding = Margins{
		Top:    top,
		Right:  right,
		Bottom: bottom,
		Left:   left,
	}
	return d
}

// SetPaddingAll sets the same padding for all sides.
//
// This is a convenience method for uniform padding.
//
// Example:
//
//	div.SetPaddingAll(15) // 15 points on all sides
func (d *Division) SetPaddingAll(p float64) *Division {
	d.padding = Margins{
		Top:    p,
		Right:  p,
		Bottom: p,
		Left:   p,
	}
	return d
}

// SetMargins sets margins for the division.
//
// Margins are the outer spacing around the division (outside the border).
//
// Example:
//
//	div.SetMargins(Margins{Top: 10, Right: 5, Bottom: 10, Left: 5})
func (d *Division) SetMargins(m Margins) *Division {
	d.margins = m
	return d
}

// SetWidth sets an explicit width for the division.
//
// If width is 0, the division uses the full available width.
//
// Example:
//
//	div.SetWidth(300) // 300 points wide
func (d *Division) SetWidth(w float64) *Division {
	d.width = w
	return d
}

// SetMinHeight sets the minimum height for the division.
//
// The division will be at least this tall, even if content is shorter.
//
// Example:
//
//	div.SetMinHeight(100) // At least 100 points tall
func (d *Division) SetMinHeight(h float64) *Division {
	d.minHeight = h
	return d
}

// Add adds a drawable element to the division.
//
// Elements are drawn in the order they are added.
//
// Example:
//
//	div.Add(NewParagraph("First"))
//	div.Add(NewParagraph("Second"))
func (d *Division) Add(drawable Drawable) *Division {
	d.drawables = append(d.drawables, drawable)
	return d
}

// Clear removes all drawable elements from the division.
//
// This does not affect styling (background, borders, padding, margins).
//
// Example:
//
//	div.Clear() // Remove all content
func (d *Division) Clear() *Division {
	d.drawables = d.drawables[:0]
	return d
}

// Drawables returns all drawable elements in the division.
//
// The returned slice is a direct reference, not a copy.
func (d *Division) Drawables() []Drawable {
	return d.drawables
}

// Background returns the background color, or nil if transparent.
func (d *Division) Background() *Color {
	return d.background
}

// Padding returns the padding margins.
func (d *Division) Padding() Margins {
	return d.padding
}

// ContentWidth calculates the available width for content.
//
// This accounts for padding and border widths.
func (d *Division) ContentWidth(ctx *LayoutContext) float64 {
	// Start with division width (0 = full available width).
	width := d.width
	if width == 0 {
		width = ctx.AvailableWidth()
	}

	// Subtract padding.
	width -= d.padding.Left + d.padding.Right

	// Subtract border widths.
	width -= d.getBorderLeftWidth() + d.getBorderRightWidth()

	return width
}

// Height calculates the total height of the division.
//
// This includes content height, padding, and borders.
func (d *Division) Height(ctx *LayoutContext) float64 {
	// Calculate content height.
	contentHeight := d.calculateContentHeight(ctx)

	// Add padding.
	totalHeight := contentHeight + d.padding.Top + d.padding.Bottom

	// Add border heights.
	totalHeight += d.getBorderTopWidth() + d.getBorderBottomWidth()

	// Apply minimum height.
	if totalHeight < d.minHeight {
		totalHeight = d.minHeight
	}

	return totalHeight
}

// Draw renders the division and its contents on the page.
//
// Drawing sequence:
// 1. Apply margins (adjust cursor)
// 2. Draw background
// 3. Draw borders
// 4. Draw content (with padding)
// 5. Update cursor position.
func (d *Division) Draw(ctx *LayoutContext, page *Page) error {
	// Apply top margin.
	ctx.CursorY += d.margins.Top

	// Calculate division dimensions.
	divWidth := d.calculateDivisionWidth(ctx)
	divHeight := d.Height(ctx)

	// Calculate position in PDF coordinates.
	x := ctx.ContentLeft() + d.margins.Left
	y := ctx.CurrentPDFY() - divHeight + d.margins.Top

	// Draw background if set.
	if err := d.drawBackground(page, x, y, divWidth, divHeight); err != nil {
		return err
	}

	// Draw borders if set.
	if err := d.drawBorders(page, x, y, divWidth, divHeight); err != nil {
		return err
	}

	// Draw content with padding.
	if err := d.drawContent(ctx, page, divWidth); err != nil {
		return err
	}

	// Update cursor position (move past division + bottom margin).
	ctx.CursorY += divHeight + d.margins.Bottom

	return nil
}

// calculateDivisionWidth calculates the total division width.
func (d *Division) calculateDivisionWidth(ctx *LayoutContext) float64 {
	if d.width > 0 {
		return d.width
	}
	return ctx.AvailableWidth() - d.margins.Left - d.margins.Right
}

// calculateContentHeight calculates the total height of all content.
func (d *Division) calculateContentHeight(ctx *LayoutContext) float64 {
	// Create inner context with padding applied.
	innerCtx := d.createInnerContext(ctx)

	totalHeight := 0.0
	for _, drawable := range d.drawables {
		totalHeight += drawable.Height(innerCtx)
	}

	return totalHeight
}

// drawBackground draws the background rectangle if set.
func (d *Division) drawBackground(page *Page, x, y, width, height float64) error {
	if d.background == nil {
		return nil
	}

	opts := &RectOptions{
		FillColor: d.background,
	}
	return page.DrawRect(x, y, width, height, opts)
}

// drawBorders draws all borders if set.
func (d *Division) drawBorders(page *Page, x, y, width, height float64) error {
	// Draw top border.
	if err := d.drawBorderTop(page, x, y, width, height); err != nil {
		return err
	}

	// Draw right border.
	if err := d.drawBorderRight(page, x, y, width, height); err != nil {
		return err
	}

	// Draw bottom border.
	if err := d.drawBorderBottom(page, x, y, width, height); err != nil {
		return err
	}

	// Draw left border.
	if err := d.drawBorderLeft(page, x, y, width, height); err != nil {
		return err
	}

	return nil
}

// drawBorderTop draws the top border if set.
func (d *Division) drawBorderTop(page *Page, x, y, width, height float64) error {
	border := d.getEffectiveBorderTop()
	if border == nil {
		return nil
	}

	lineOpts := &LineOptions{
		Width: border.Width,
		Color: border.Color,
	}
	return page.DrawLine(x, y+height, x+width, y+height, lineOpts)
}

// drawBorderRight draws the right border if set.
func (d *Division) drawBorderRight(page *Page, x, y, width, height float64) error {
	border := d.getEffectiveBorderRight()
	if border == nil {
		return nil
	}

	lineOpts := &LineOptions{
		Width: border.Width,
		Color: border.Color,
	}
	return page.DrawLine(x+width, y, x+width, y+height, lineOpts)
}

// drawBorderBottom draws the bottom border if set.
func (d *Division) drawBorderBottom(page *Page, x, y, width, _ float64) error {
	border := d.getEffectiveBorderBottom()
	if border == nil {
		return nil
	}

	lineOpts := &LineOptions{
		Width: border.Width,
		Color: border.Color,
	}
	return page.DrawLine(x, y, x+width, y, lineOpts)
}

// drawBorderLeft draws the left border if set.
func (d *Division) drawBorderLeft(page *Page, x, y, _ float64, height float64) error {
	border := d.getEffectiveBorderLeft()
	if border == nil {
		return nil
	}

	lineOpts := &LineOptions{
		Width: border.Width,
		Color: border.Color,
	}
	return page.DrawLine(x, y, x, y+height, lineOpts)
}

// drawContent draws all content elements with padding applied.
func (d *Division) drawContent(ctx *LayoutContext, page *Page, _ float64) error {
	// Create inner context with adjusted cursor for padding and borders.
	innerCtx := d.createInnerContext(ctx)

	// Adjust cursor for top padding and top border.
	innerCtx.CursorY += d.padding.Top + d.getBorderTopWidth()

	// Draw each drawable sequentially.
	for _, drawable := range d.drawables {
		if err := drawable.Draw(innerCtx, page); err != nil {
			return err
		}
	}

	return nil
}

// createInnerContext creates a layout context for content.
//
// This context has adjusted dimensions to account for padding and borders.
func (d *Division) createInnerContext(ctx *LayoutContext) *LayoutContext {
	innerCtx := &LayoutContext{
		PageWidth:  ctx.PageWidth,
		PageHeight: ctx.PageHeight,
		Margins: Margins{
			Top:    ctx.Margins.Top + d.margins.Top + d.padding.Top + d.getBorderTopWidth(),
			Right:  ctx.Margins.Right + d.margins.Right + d.padding.Right + d.getBorderRightWidth(),
			Bottom: ctx.Margins.Bottom + d.margins.Bottom + d.padding.Bottom + d.getBorderBottomWidth(),
			Left:   ctx.Margins.Left + d.margins.Left + d.padding.Left + d.getBorderLeftWidth(),
		},
		CursorX: ctx.ContentLeft() + d.margins.Left + d.padding.Left + d.getBorderLeftWidth(),
		CursorY: ctx.CursorY,
	}
	return innerCtx
}

// getEffectiveBorderTop returns the top border (individual or default).
func (d *Division) getEffectiveBorderTop() *Border {
	if d.borderTop != nil {
		return d.borderTop
	}
	return d.border
}

// getEffectiveBorderRight returns the right border (individual or default).
func (d *Division) getEffectiveBorderRight() *Border {
	if d.borderRight != nil {
		return d.borderRight
	}
	return d.border
}

// getEffectiveBorderBottom returns the bottom border (individual or default).
func (d *Division) getEffectiveBorderBottom() *Border {
	if d.borderBottom != nil {
		return d.borderBottom
	}
	return d.border
}

// getEffectiveBorderLeft returns the left border (individual or default).
func (d *Division) getEffectiveBorderLeft() *Border {
	if d.borderLeft != nil {
		return d.borderLeft
	}
	return d.border
}

// getBorderTopWidth returns the top border width or 0 if no border.
func (d *Division) getBorderTopWidth() float64 {
	if border := d.getEffectiveBorderTop(); border != nil {
		return border.Width
	}
	return 0
}

// getBorderRightWidth returns the right border width or 0 if no border.
func (d *Division) getBorderRightWidth() float64 {
	if border := d.getEffectiveBorderRight(); border != nil {
		return border.Width
	}
	return 0
}

// getBorderBottomWidth returns the bottom border width or 0 if no border.
func (d *Division) getBorderBottomWidth() float64 {
	if border := d.getEffectiveBorderBottom(); border != nil {
		return border.Width
	}
	return 0
}

// getBorderLeftWidth returns the left border width or 0 if no border.
func (d *Division) getBorderLeftWidth() float64 {
	if border := d.getEffectiveBorderLeft(); border != nil {
		return border.Width
	}
	return 0
}
