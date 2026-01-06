// Package table provides domain entities for PDF table extraction.
//
// This is the Domain layer in DDD/Clean Architecture.
// It contains the core business logic for representing extracted tables.
package table

import "fmt"

// Cell represents a single cell in an extracted table.
//
// A cell contains:
//   - Text content (potentially multi-line)
//   - Position within the table (row, column)
//   - Span information for merged cells
//   - Bounding rectangle
//   - Text alignment
//
// Cells are value objects in DDD - they are compared by value, not identity.
type Cell struct {
	Text      string    // Text content (may contain newlines)
	Row       int       // Row index (0-based)
	Column    int       // Column index (0-based)
	RowSpan   int       // Number of rows this cell spans (1 = no merge)
	ColSpan   int       // Number of columns this cell spans (1 = no merge)
	Bounds    Rectangle // Bounding rectangle
	TextAlign TextAlign // Text alignment within cell
}

// TextAlign represents text alignment within a cell.
type TextAlign int

const (
	// AlignLeft indicates text is left-aligned.
	AlignLeft TextAlign = iota
	// AlignCenter indicates text is center-aligned.
	AlignCenter
	// AlignRight indicates text is right-aligned.
	AlignRight
)

// String returns a string representation of the text alignment.
func (ta TextAlign) String() string {
	switch ta {
	case AlignLeft:
		return "left"
	case AlignCenter:
		return "center"
	case AlignRight:
		return "right"
	default:
		return "unknown"
	}
}

// NewCell creates a new Cell with the given text and position.
//
// By default, cells have RowSpan=1, ColSpan=1 (not merged),
// and TextAlign=AlignLeft.
func NewCell(text string, row, col int) *Cell {
	return &Cell{
		Text:      text,
		Row:       row,
		Column:    col,
		RowSpan:   1,
		ColSpan:   1,
		TextAlign: AlignLeft,
	}
}

// NewCellWithBounds creates a new Cell with text, position, and bounds.
func NewCellWithBounds(text string, row, col int, bounds Rectangle) *Cell {
	return &Cell{
		Text:      text,
		Row:       row,
		Column:    col,
		RowSpan:   1,
		ColSpan:   1,
		Bounds:    bounds,
		TextAlign: AlignLeft,
	}
}

// IsMerged returns true if this cell is merged (spans multiple rows or columns).
func (c *Cell) IsMerged() bool {
	return c.RowSpan > 1 || c.ColSpan > 1
}

// IsEmpty returns true if the cell has no text content.
func (c *Cell) IsEmpty() bool {
	return len(c.Text) == 0
}

// WithRowSpan returns a new Cell with the specified row span.
func (c *Cell) WithRowSpan(rowSpan int) *Cell {
	if rowSpan < 1 {
		rowSpan = 1
	}
	return &Cell{
		Text:      c.Text,
		Row:       c.Row,
		Column:    c.Column,
		RowSpan:   rowSpan,
		ColSpan:   c.ColSpan,
		Bounds:    c.Bounds,
		TextAlign: c.TextAlign,
	}
}

// WithColSpan returns a new Cell with the specified column span.
func (c *Cell) WithColSpan(colSpan int) *Cell {
	if colSpan < 1 {
		colSpan = 1
	}
	return &Cell{
		Text:      c.Text,
		Row:       c.Row,
		Column:    c.Column,
		RowSpan:   c.RowSpan,
		ColSpan:   colSpan,
		Bounds:    c.Bounds,
		TextAlign: c.TextAlign,
	}
}

// WithAlignment returns a new Cell with the specified text alignment.
func (c *Cell) WithAlignment(align TextAlign) *Cell {
	return &Cell{
		Text:      c.Text,
		Row:       c.Row,
		Column:    c.Column,
		RowSpan:   c.RowSpan,
		ColSpan:   c.ColSpan,
		Bounds:    c.Bounds,
		TextAlign: align,
	}
}

// String returns a string representation of the cell.
func (c *Cell) String() string {
	if c.IsMerged() {
		return fmt.Sprintf("Cell{text=%q, row=%d, col=%d, rowSpan=%d, colSpan=%d}",
			c.Text, c.Row, c.Column, c.RowSpan, c.ColSpan)
	}
	return fmt.Sprintf("Cell{text=%q, row=%d, col=%d}", c.Text, c.Row, c.Column)
}
