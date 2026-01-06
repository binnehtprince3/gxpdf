package fonts

import (
	"testing"
)

// TestNewFontSubset tests creating a new font subset.
func TestNewFontSubset(t *testing.T) {
	font := &TTFFont{
		UnitsPerEm:  1000,
		GlyphWidths: make(map[uint16]uint16),
		CharToGlyph: make(map[rune]uint16),
	}

	subset := NewFontSubset(font)

	if subset.BaseFont != font {
		t.Error("BaseFont not set correctly")
	}
	if subset.UsedChars == nil {
		t.Error("UsedChars not initialized")
	}
	if subset.GlyphMapping == nil {
		t.Error("GlyphMapping not initialized")
	}
}

// TestUseChar tests marking characters as used.
func TestUseChar(t *testing.T) {
	font := &TTFFont{
		UnitsPerEm:  1000,
		GlyphWidths: make(map[uint16]uint16),
		CharToGlyph: make(map[rune]uint16),
	}
	subset := NewFontSubset(font)

	// Mark some characters as used.
	subset.UseChar('A')
	subset.UseChar('B')
	subset.UseChar('C')

	// Verify characters are marked.
	if !subset.UsedChars['A'] {
		t.Error("character 'A' not marked as used")
	}
	if !subset.UsedChars['B'] {
		t.Error("character 'B' not marked as used")
	}
	if !subset.UsedChars['C'] {
		t.Error("character 'C' not marked as used")
	}
	if subset.UsedChars['D'] {
		t.Error("character 'D' incorrectly marked as used")
	}
}

// TestUseString tests marking string characters as used.
func TestUseString(t *testing.T) {
	font := &TTFFont{
		UnitsPerEm:  1000,
		GlyphWidths: make(map[uint16]uint16),
		CharToGlyph: make(map[rune]uint16),
	}
	subset := NewFontSubset(font)

	// Mark string characters as used.
	subset.UseString("Hello")

	// Verify characters are marked.
	for _, ch := range "Hello" {
		if !subset.UsedChars[ch] {
			t.Errorf("character %q not marked as used", ch)
		}
	}

	// Verify unused character.
	if subset.UsedChars['X'] {
		t.Error("character 'X' incorrectly marked as used")
	}
}

// TestGetCharWidth tests getting character widths.
func TestGetCharWidth(t *testing.T) {
	font := &TTFFont{
		UnitsPerEm: 1000,
		GlyphWidths: map[uint16]uint16{
			1: 500,  // Glyph 1 = 500 units.
			2: 750,  // Glyph 2 = 750 units.
			3: 1000, // Glyph 3 = 1000 units.
		},
		CharToGlyph: map[rune]uint16{
			'A': 1,
			'B': 2,
			'C': 3,
		},
	}
	subset := NewFontSubset(font)

	tests := []struct {
		char     rune
		expected uint16
	}{
		{'A', 500},
		{'B', 750},
		{'C', 1000},
		{'X', 0}, // Unknown character.
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			width := subset.GetCharWidth(tt.char)
			if width != tt.expected {
				t.Errorf("expected width %d, got %d", tt.expected, width)
			}
		})
	}
}

// TestMeasureString tests measuring string width.
func TestMeasureString(t *testing.T) {
	font := &TTFFont{
		UnitsPerEm: 1000,
		GlyphWidths: map[uint16]uint16{
			1: 500, // A.
			2: 600, // B.
			3: 700, // C.
		},
		CharToGlyph: map[rune]uint16{
			'A': 1,
			'B': 2,
			'C': 3,
		},
	}
	subset := NewFontSubset(font)

	tests := []struct {
		text     string
		size     float64
		expected float64
	}{
		{"A", 12.0, 6.0},    // 500 * 12 / 1000 = 6.
		{"AB", 12.0, 13.2},  // (500+600) * 12 / 1000 = 13.2.
		{"ABC", 10.0, 18.0}, // (500+600+700) * 10 / 1000 = 18.
		{"", 12.0, 0.0},     // Empty string.
		{"X", 12.0, 0.0},    // Unknown character.
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			width := subset.MeasureString(tt.text, tt.size)
			if width != tt.expected {
				t.Errorf("expected width %.2f, got %.2f", tt.expected, width)
			}
		})
	}
}

// TestGetWidths tests getting widths array for PDF.
func TestGetWidths(t *testing.T) {
	font := &TTFFont{
		UnitsPerEm: 1000,
		GlyphWidths: map[uint16]uint16{
			1: 500,
			2: 600,
			3: 700,
		},
		CharToGlyph: map[rune]uint16{
			'A': 1, // Code 65.
			'C': 3, // Code 67.
		},
	}
	subset := NewFontSubset(font)

	// Mark characters as used.
	subset.UseChar('A')
	subset.UseChar('C')

	// Get widths array.
	first, last, widths := subset.GetWidths()

	// Verify range.
	if first != 65 {
		t.Errorf("expected firstChar 65, got %d", first)
	}
	if last != 67 {
		t.Errorf("expected lastChar 67, got %d", last)
	}

	// Verify widths array (A=500, B=0, C=700).
	expected := []int{500, 0, 700}
	if len(widths) != len(expected) {
		t.Fatalf("expected %d widths, got %d", len(expected), len(widths))
	}
	for i, w := range widths {
		if w != expected[i] {
			t.Errorf("widths[%d]: expected %d, got %d", i, expected[i], w)
		}
	}
}

// TestIdentifyUsedGlyphs tests identifying used glyphs.
func TestIdentifyUsedGlyphs(t *testing.T) {
	font := &TTFFont{
		UnitsPerEm: 1000,
		GlyphWidths: map[uint16]uint16{
			0: 0,
			1: 500,
			2: 600,
		},
		CharToGlyph: map[rune]uint16{
			'A': 1,
			'B': 2,
		},
	}
	subset := NewFontSubset(font)

	// Mark character 'A' as used.
	subset.UseChar('A')

	// Identify used glyphs.
	glyphs := subset.identifyUsedGlyphs()

	// Should include glyph 0 (.notdef) and glyph 1 ('A').
	if len(glyphs) != 2 {
		t.Fatalf("expected 2 glyphs, got %d", len(glyphs))
	}

	// Verify glyphs (should be sorted).
	if glyphs[0] != 0 {
		t.Errorf("expected glyph 0 at index 0, got %d", glyphs[0])
	}
	if glyphs[1] != 1 {
		t.Errorf("expected glyph 1 at index 1, got %d", glyphs[1])
	}
}
