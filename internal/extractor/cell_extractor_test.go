package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCellExtractor(t *testing.T) {
	elements := []*TextElement{
		NewTextElement("test", 10, 20, 30, 10, "Arial", 12),
	}

	extractor := NewCellExtractor(elements)
	assert.NotNil(t, extractor)
	assert.Len(t, extractor.textElements, 1)
}

func TestCellExtractor_ExtractCellContent_Empty(t *testing.T) {
	// Empty extractor
	extractor := NewCellExtractor([]*TextElement{})
	bounds := NewRectangle(0, 0, 100, 100)

	content := extractor.ExtractCellContent(bounds)
	assert.Equal(t, "", content)
}

func TestCellExtractor_ExtractCellContent_SingleElement(t *testing.T) {
	elements := []*TextElement{
		NewTextElement("Hello", 10, 10, 30, 10, "Arial", 12),
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 50, 50)

	content := extractor.ExtractCellContent(bounds)
	assert.Equal(t, "Hello", content)
}

func TestCellExtractor_ExtractCellContent_MultipleElementsOneLine(t *testing.T) {
	// Elements on same line (same Y), should be joined with space
	elements := []*TextElement{
		NewTextElement("Hello", 10, 20, 30, 10, "Arial", 12),
		NewTextElement("World", 50, 20, 30, 10, "Arial", 12),
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 100, 50)

	content := extractor.ExtractCellContent(bounds)
	assert.Equal(t, "Hello World", content)
}

func TestCellExtractor_ExtractCellContent_MultipleLines(t *testing.T) {
	// Elements on different lines (different Y), should be joined with newline
	// PDF Y increases upward, so line 1 is higher (Y=30), line 2 is lower (Y=10)
	elements := []*TextElement{
		NewTextElement("Line 1", 10, 30, 40, 10, "Arial", 12), // Top line
		NewTextElement("Line 2", 10, 10, 40, 10, "Arial", 12), // Bottom line
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 100, 50)

	content := extractor.ExtractCellContent(bounds)
	// Should be ordered top to bottom
	assert.Equal(t, "Line 1\nLine 2", content)
}

func TestCellExtractor_ExtractCellContent_OutsideBounds(t *testing.T) {
	elements := []*TextElement{
		NewTextElement("Inside", 10, 10, 30, 10, "Arial", 12),
		NewTextElement("Outside", 200, 200, 30, 10, "Arial", 12),
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 50, 50)

	content := extractor.ExtractCellContent(bounds)
	// Only "Inside" should be included
	assert.Equal(t, "Inside", content)
}

func TestCellExtractor_ExtractCellContent_AdjacentWords(t *testing.T) {
	// Words that are immediately adjacent (no gap)
	elements := []*TextElement{
		NewTextElement("Hello", 10, 20, 30, 10, "Arial", 12),
		NewTextElement("World", 40, 20, 30, 10, "Arial", 12), // Exactly adjacent (10+30=40)
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 100, 50)

	content := extractor.ExtractCellContent(bounds)
	// Should not add space between adjacent words
	assert.Equal(t, "HelloWorld", content)
}

func TestCellExtractor_ExtractCellContent_ComplexTable(t *testing.T) {
	// Simulate a complex cell with multiple lines and words
	elements := []*TextElement{
		// Line 1 (Y=50)
		NewTextElement("Product", 10, 50, 40, 10, "Arial", 12),
		NewTextElement("Name:", 55, 50, 30, 10, "Arial", 12),
		// Line 2 (Y=35)
		NewTextElement("Widget", 10, 35, 35, 10, "Arial", 12),
		NewTextElement("Pro", 50, 35, 20, 10, "Arial", 12),
		// Line 3 (Y=20)
		NewTextElement("v2.0", 10, 20, 25, 10, "Arial", 12),
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 100, 70)

	content := extractor.ExtractCellContent(bounds)
	expected := "Product Name:\nWidget Pro\nv2.0"
	assert.Equal(t, expected, content)
}

func TestCellExtractor_FindElementsInBounds(t *testing.T) {
	elements := []*TextElement{
		NewTextElement("A", 10, 10, 10, 10, "Arial", 12), // Center at (15, 15)
		NewTextElement("B", 50, 50, 10, 10, "Arial", 12), // Center at (55, 55)
		NewTextElement("C", 90, 90, 10, 10, "Arial", 12), // Center at (95, 95) - outside
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 80, 80)

	found := extractor.FindElementsInBounds(bounds)
	require.Len(t, found, 2)
	assert.Equal(t, "A", found[0].Text)
	assert.Equal(t, "B", found[1].Text)
}

func TestCellExtractor_groupByLine(t *testing.T) {
	elements := []*TextElement{
		NewTextElement("A", 10, 50, 10, 10, "Arial", 12), // Y=50
		NewTextElement("B", 30, 51, 10, 10, "Arial", 12), // Y=51 (same line)
		NewTextElement("C", 10, 30, 10, 10, "Arial", 12), // Y=30 (different line)
		NewTextElement("D", 30, 29, 10, 10, "Arial", 12), // Y=29 (same line as C)
	}

	extractor := NewCellExtractor(elements)
	lines := extractor.groupByLine(elements)

	require.Len(t, lines, 2)

	// Lines should contain correct elements
	// Note: Order of lines is not guaranteed at this point
	var line1, line2 *textLine
	for _, line := range lines {
		if len(line.elements) == 2 {
			if line.elements[0].Text == "A" || line.elements[0].Text == "B" {
				line1 = line
			} else {
				line2 = line
			}
		}
	}

	require.NotNil(t, line1)
	require.NotNil(t, line2)
	assert.Len(t, line1.elements, 2)
	assert.Len(t, line2.elements, 2)
}

func TestCellExtractor_sortLines(t *testing.T) {
	// Create lines in random order
	lines := []*textLine{
		{y: 10, elements: []*TextElement{NewTextElement("Bottom", 10, 10, 10, 10, "Arial", 12)}},
		{y: 50, elements: []*TextElement{NewTextElement("Top", 10, 50, 10, 10, "Arial", 12)}},
		{y: 30, elements: []*TextElement{NewTextElement("Middle", 10, 30, 10, 10, "Arial", 12)}},
	}

	extractor := NewCellExtractor(nil)
	extractor.sortLines(lines)

	// Should be sorted top to bottom (descending Y)
	require.Len(t, lines, 3)
	assert.Equal(t, "Top", lines[0].elements[0].Text)
	assert.Equal(t, "Middle", lines[1].elements[0].Text)
	assert.Equal(t, "Bottom", lines[2].elements[0].Text)
}

func TestCellExtractor_sortLines_WithinLine(t *testing.T) {
	// Elements within a line should be sorted left to right
	lines := []*textLine{
		{
			y: 50,
			elements: []*TextElement{
				NewTextElement("C", 50, 50, 10, 10, "Arial", 12), // X=50
				NewTextElement("A", 10, 50, 10, 10, "Arial", 12), // X=10
				NewTextElement("B", 30, 50, 10, 10, "Arial", 12), // X=30
			},
		},
	}

	extractor := NewCellExtractor(nil)
	extractor.sortLines(lines)

	require.Len(t, lines[0].elements, 3)
	assert.Equal(t, "A", lines[0].elements[0].Text)
	assert.Equal(t, "B", lines[0].elements[1].Text)
	assert.Equal(t, "C", lines[0].elements[2].Text)
}

func TestCellExtractor_buildTextFromLines(t *testing.T) {
	lines := []*textLine{
		{
			y: 50,
			elements: []*TextElement{
				NewTextElement("Hello", 10, 50, 30, 10, "Arial", 12),
				NewTextElement("World", 50, 50, 30, 10, "Arial", 12),
			},
		},
		{
			y: 30,
			elements: []*TextElement{
				NewTextElement("Second", 10, 30, 35, 10, "Arial", 12),
				NewTextElement("Line", 50, 30, 25, 10, "Arial", 12),
			},
		},
	}

	extractor := NewCellExtractor(nil)
	text := extractor.buildTextFromLines(lines)

	expected := "Hello World\nSecond Line"
	assert.Equal(t, expected, text)
}

func TestCellExtractor_calculateAverageFontSize(t *testing.T) {
	tests := []struct {
		name     string
		elements []*TextElement
		expected float64
	}{
		{
			name:     "empty",
			elements: []*TextElement{},
			expected: 12.0, // Default
		},
		{
			name: "single element",
			elements: []*TextElement{
				NewTextElement("test", 0, 0, 10, 10, "Arial", 14),
			},
			expected: 14.0,
		},
		{
			name: "multiple elements",
			elements: []*TextElement{
				NewTextElement("a", 0, 0, 10, 10, "Arial", 10),
				NewTextElement("b", 0, 0, 10, 10, "Arial", 12),
				NewTextElement("c", 0, 0, 10, 10, "Arial", 14),
			},
			expected: 12.0, // (10+12+14)/3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewCellExtractor(nil)
			result := extractor.calculateAverageFontSize(tt.elements)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAbs(t *testing.T) {
	assert.Equal(t, 5.0, abs(5.0))
	assert.Equal(t, 5.0, abs(-5.0))
	assert.Equal(t, 0.0, abs(0.0))
}

// TestCellExtractor_ExtractCellContent_VTB_MultiLineCells tests VTB-style multi-line cells.
//
// VTB bank statements have cells with multi-line content (2-3 lines per cell).
// Line spacing is ~12-15pt for 10pt font.
//
// Example VTB cell:
//
//	Line 1 (Y=220): "17.06.2025"
//	Line 2 (Y=208): "19:14:57"   (12pt below - should be merged)
//	Line 3 (Y=195): "Описание"   (13pt below - should be merged)
//
// With old threshold (0.5x = 5pt): Lines treated as separate rows ❌
// With new threshold (1.5x = 15pt): Lines correctly merged ✅
//
// See: BUG_REPORT_VTB_TABLE_STRUCTURE.md, ANALYSIS_VTB_TABLE_MULTI_LINE_CELLS.md
func TestCellExtractor_ExtractCellContent_VTB_MultiLineCells(t *testing.T) {
	// Simulate VTB cell with 3 lines (date, time, description)
	// Font size = 10pt, line spacing = 12-13pt
	elements := []*TextElement{
		// Line 1 (top): Date
		NewTextElement("17.06.2025", 20, 220, 50, 10, "Arial", 10), // Y=220
		// Line 2 (middle): Time - 12pt below
		NewTextElement("19:14:57", 20, 208, 45, 10, "Arial", 10), // Y=208 (220-12=208)
		// Line 3 (bottom): Description - 13pt below
		NewTextElement("Описание", 20, 195, 55, 10, "Arial", 10), // Y=195 (208-13=195)
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 100, 250)

	content := extractor.ExtractCellContent(bounds)

	// All 3 lines should be merged with newlines
	expected := "17.06.2025\n19:14:57\nОписание"
	assert.Equal(t, expected, content, "VTB multi-line cell should merge all lines")

	// Verify groupByLine creates 3 separate lines (12-13pt spacing > 5pt threshold)
	lines := extractor.groupByLine(elements)
	assert.Len(t, lines, 3, "Elements with 12-13pt spacing should be in 3 separate lines")
	// Each line should contain one element
	for i, line := range lines {
		assert.Len(t, line.elements, 1, "Line %d should contain 1 element", i)
	}
}

// TestCellExtractor_ExtractCellContent_AlfaBank_Regression tests Alfa-Bank single-line cells.
//
// Alfa-Bank statements have mostly single-line cells with minimal spacing (2-3pt).
// This test ensures we didn't break Alfa-Bank support with the threshold change.
//
// Previous behavior (0.5x = 5pt): ✅ Worked
// New behavior (1.5x = 15pt): ✅ Should still work
func TestCellExtractor_ExtractCellContent_AlfaBank_Regression(t *testing.T) {
	// Simulate Alfa-Bank cell with elements on same line (Y difference < 3pt)
	elements := []*TextElement{
		NewTextElement("01.02.24", 20, 220, 45, 10, "Arial", 10), // Y=220
		NewTextElement("Операция", 70, 221, 55, 10, "Arial", 10), // Y=221 (1pt diff)
		NewTextElement("1000.00", 130, 220, 50, 10, "Arial", 10), // Y=220 (same)
	}

	extractor := NewCellExtractor(elements)
	bounds := NewRectangle(0, 0, 200, 250)

	content := extractor.ExtractCellContent(bounds)

	// All elements on same line - should be joined with spaces (no newlines)
	expected := "01.02.24 Операция 1000.00"
	assert.Equal(t, expected, content, "Alfa-Bank single-line cell should work as before")

	// Verify groupByLine creates single line group
	lines := extractor.groupByLine(elements)
	assert.Len(t, lines, 1, "Elements with Y diff 0-1pt should be in 1 line")
	assert.Len(t, lines[0].elements, 3, "Line should contain all 3 elements")
}

// TestCellExtractor_groupByLine_ThresholdBehavior tests threshold edge cases.
//
// UPDATED (2025-01): Algorithm uses 0.3x font size threshold (3pt for 10pt font)
// This is optimized for Alfa-Bank tight line spacing.
// Verifies that:
//   - Elements within threshold (< 3pt for 10pt font) are grouped together
//   - Elements beyond threshold (>= 3pt) are treated as separate lines
func TestCellExtractor_groupByLine_ThresholdBehavior(t *testing.T) {
	tests := []struct {
		name          string
		elements      []*TextElement
		expectedLines int
		description   string
	}{
		{
			name: "spacing 2pt - same line",
			elements: []*TextElement{
				NewTextElement("A", 10, 100, 10, 10, "Arial", 10), // Y=100
				NewTextElement("B", 10, 98, 10, 10, "Arial", 10),  // Y=98 (2pt diff)
			},
			expectedLines: 1,
			description:   "2pt < 3pt threshold - should be 1 line",
		},
		{
			name: "spacing 3pt - edge case (separate)",
			elements: []*TextElement{
				NewTextElement("A", 10, 100, 10, 10, "Arial", 10), // Y=100
				NewTextElement("B", 10, 97, 10, 10, "Arial", 10),  // Y=97 (3pt diff)
			},
			expectedLines: 2,
			description:   "3pt >= 3pt threshold - should be 2 lines (edge case)",
		},
		{
			name: "spacing 12pt - different lines (VTB)",
			elements: []*TextElement{
				NewTextElement("A", 10, 100, 10, 10, "Arial", 10), // Y=100
				NewTextElement("B", 10, 88, 10, 10, "Arial", 10),  // Y=88 (12pt diff)
			},
			expectedLines: 2,
			description:   "12pt > 3pt threshold - should be 2 lines (VTB multi-line cell)",
		},
		{
			name: "spacing 20pt - different lines",
			elements: []*TextElement{
				NewTextElement("A", 10, 100, 10, 10, "Arial", 10), // Y=100
				NewTextElement("B", 10, 80, 10, 10, "Arial", 10),  // Y=80 (20pt diff)
			},
			expectedLines: 2,
			description:   "20pt > 3pt threshold - should be 2 lines (different rows)",
		},
		{
			name: "mixed spacing - all separate",
			elements: []*TextElement{
				NewTextElement("A", 10, 100, 10, 10, "Arial", 10), // Y=100
				NewTextElement("B", 10, 88, 10, 10, "Arial", 10),  // Y=88 (12pt - separate)
				NewTextElement("C", 10, 76, 10, 10, "Arial", 10),  // Y=76 (12pt - separate)
				NewTextElement("D", 10, 50, 10, 10, "Arial", 10),  // Y=50 (26pt - separate)
			},
			expectedLines: 4,
			description:   "Elements with 12pt+ spacing should each be separate line (> 3pt threshold)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewCellExtractor(tt.elements)
			lines := extractor.groupByLine(tt.elements)

			assert.Len(t, lines, tt.expectedLines, tt.description)
		})
	}
}
