package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCell(t *testing.T) {
	cell := NewCell("Hello", 0, 1)

	assert.Equal(t, "Hello", cell.Text)
	assert.Equal(t, 0, cell.Row)
	assert.Equal(t, 1, cell.Column)
	assert.Equal(t, 1, cell.RowSpan)
	assert.Equal(t, 1, cell.ColSpan)
	assert.Equal(t, AlignLeft, cell.TextAlign)
}

func TestNewCellWithBounds(t *testing.T) {
	bounds := NewRectangle(10, 20, 100, 50)
	cell := NewCellWithBounds("Test", 1, 2, bounds)

	assert.Equal(t, "Test", cell.Text)
	assert.Equal(t, 1, cell.Row)
	assert.Equal(t, 2, cell.Column)
	assert.Equal(t, bounds, cell.Bounds)
}

func TestCell_IsMerged(t *testing.T) {
	tests := []struct {
		name     string
		rowSpan  int
		colSpan  int
		expected bool
	}{
		{"not merged", 1, 1, false},
		{"row merged", 2, 1, true},
		{"col merged", 1, 2, true},
		{"both merged", 2, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := &Cell{RowSpan: tt.rowSpan, ColSpan: tt.colSpan}
			assert.Equal(t, tt.expected, cell.IsMerged())
		})
	}
}

func TestCell_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"empty", "", true},
		{"non-empty", "text", false},
		{"whitespace", "   ", false}, // Whitespace is not empty
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := &Cell{Text: tt.text}
			assert.Equal(t, tt.expected, cell.IsEmpty())
		})
	}
}

func TestCell_WithRowSpan(t *testing.T) {
	original := NewCell("Test", 0, 0)

	// Valid span
	merged := original.WithRowSpan(3)
	assert.Equal(t, 3, merged.RowSpan)
	assert.Equal(t, 1, merged.ColSpan) // ColSpan unchanged
	assert.Equal(t, "Test", merged.Text)
	assert.Equal(t, 1, original.RowSpan) // Original unchanged (immutable)

	// Invalid span (< 1) should default to 1
	invalid := original.WithRowSpan(0)
	assert.Equal(t, 1, invalid.RowSpan)

	negative := original.WithRowSpan(-5)
	assert.Equal(t, 1, negative.RowSpan)
}

func TestCell_WithColSpan(t *testing.T) {
	original := NewCell("Test", 0, 0)

	// Valid span
	merged := original.WithColSpan(2)
	assert.Equal(t, 2, merged.ColSpan)
	assert.Equal(t, 1, merged.RowSpan) // RowSpan unchanged
	assert.Equal(t, "Test", merged.Text)
	assert.Equal(t, 1, original.ColSpan) // Original unchanged

	// Invalid span
	invalid := original.WithColSpan(0)
	assert.Equal(t, 1, invalid.ColSpan)
}

func TestCell_WithAlignment(t *testing.T) {
	original := NewCell("Test", 0, 0)

	centered := original.WithAlignment(AlignCenter)
	assert.Equal(t, AlignCenter, centered.TextAlign)
	assert.Equal(t, AlignLeft, original.TextAlign) // Original unchanged

	right := original.WithAlignment(AlignRight)
	assert.Equal(t, AlignRight, right.TextAlign)
}

func TestTextAlign_String(t *testing.T) {
	tests := []struct {
		align    TextAlign
		expected string
	}{
		{AlignLeft, "left"},
		{AlignCenter, "center"},
		{AlignRight, "right"},
		{TextAlign(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.align.String())
		})
	}
}

func TestCell_String(t *testing.T) {
	// Regular cell
	cell := NewCell("Hello", 1, 2)
	str := cell.String()
	assert.Contains(t, str, "Hello")
	assert.Contains(t, str, "row=1")
	assert.Contains(t, str, "col=2")

	// Merged cell
	merged := cell.WithRowSpan(2).WithColSpan(3)
	str = merged.String()
	assert.Contains(t, str, "rowSpan=2")
	assert.Contains(t, str, "colSpan=3")
}
