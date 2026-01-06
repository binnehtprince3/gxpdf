package fonts

import (
	"math"
	"testing"
)

// TestGetCharWidthHelvetica tests character width retrieval for Helvetica.
func TestGetCharWidthHelvetica(t *testing.T) {
	m := Helvetica.GetMetrics()
	if m == nil {
		t.Fatal("Helvetica metrics should not be nil")
	}

	tests := []struct {
		char  rune
		width int
	}{
		{'A', 667},
		{'a', 556},
		{' ', 278},
		{'M', 833},
		{'i', 222},
		{'W', 944},
	}

	for _, tt := range tests {
		got := m.GetCharWidth(tt.char)
		if got != tt.width {
			t.Errorf("Helvetica GetCharWidth(%q) = %d, want %d", tt.char, got, tt.width)
		}
	}
}

// TestGetCharWidthTimesRoman tests character width retrieval for Times-Roman.
func TestGetCharWidthTimesRoman(t *testing.T) {
	m := TimesRoman.GetMetrics()
	if m == nil {
		t.Fatal("Times-Roman metrics should not be nil")
	}

	tests := []struct {
		char  rune
		width int
	}{
		{'A', 722},
		{'a', 444},
		{' ', 250},
		{'M', 889},
		{'i', 278},
	}

	for _, tt := range tests {
		got := m.GetCharWidth(tt.char)
		if got != tt.width {
			t.Errorf("Times-Roman GetCharWidth(%q) = %d, want %d", tt.char, got, tt.width)
		}
	}
}

// TestGetCharWidthCourier tests that Courier is monospaced (all chars = 600).
func TestGetCharWidthCourier(t *testing.T) {
	m := Courier.GetMetrics()
	if m == nil {
		t.Fatal("Courier metrics should not be nil")
	}

	// All Courier characters should be 600 units wide
	chars := []rune{'A', 'a', ' ', 'M', 'i', 'W', '!', '@', '#'}
	for _, ch := range chars {
		got := m.GetCharWidth(ch)
		if got != 600 {
			t.Errorf("Courier GetCharWidth(%q) = %d, want 600 (monospace)", ch, got)
		}
	}
}

// TestMissingCharacterHandling tests that unknown characters use default width.
func TestMissingCharacterHandling(t *testing.T) {
	m := Helvetica.GetMetrics()
	if m == nil {
		t.Fatal("Helvetica metrics should not be nil")
	}

	// Chinese character - not in WinAnsiEncoding
	chineseChar := '中'
	got := m.GetCharWidth(chineseChar)
	if got != m.DefaultWidth {
		t.Errorf("Unknown char width = %d, want DefaultWidth %d", got, m.DefaultWidth)
	}

	// Verify default width is space width (278 for Helvetica)
	if m.DefaultWidth != 278 {
		t.Errorf("Helvetica DefaultWidth = %d, want 278", m.DefaultWidth)
	}
}

// TestMeasureStringHelvetica tests string measurement accuracy.
func TestMeasureStringHelvetica(t *testing.T) {
	m := Helvetica.GetMetrics()
	if m == nil {
		t.Fatal("Helvetica metrics should not be nil")
	}

	// "Hello" at 12pt
	// H=722, e=556, l=222, l=222, o=556 = 2278 font units
	// At 12pt: 2278 * 12 / 1000 = 27.336 points
	width := m.MeasureString("Hello", 12.0)
	expected := 27.336
	if !floatEquals(width, expected, 0.001) {
		t.Errorf("MeasureString('Hello', 12) = %f, want %f", width, expected)
	}

	// Empty string should return 0
	width = m.MeasureString("", 12.0)
	if width != 0 {
		t.Errorf("MeasureString('', 12) = %f, want 0", width)
	}

	// Size 0 should return 0
	width = m.MeasureString("Hello", 0)
	if width != 0 {
		t.Errorf("MeasureString('Hello', 0) = %f, want 0", width)
	}
}

// TestMeasureStringConvenienceFunction tests the package-level MeasureString.
func TestMeasureStringConvenienceFunction(t *testing.T) {
	// Test with valid font
	width := MeasureString("Helvetica", "Test", 10.0)
	if width <= 0 {
		t.Errorf("MeasureString('Helvetica', 'Test', 10) = %f, want > 0", width)
	}

	// Test with unknown font
	width = MeasureString("UnknownFont", "Test", 10.0)
	if width != 0 {
		t.Errorf("MeasureString('UnknownFont', 'Test', 10) = %f, want 0", width)
	}
}

// TestAllFontsHaveMetrics verifies all 14 standard fonts have metrics.
func TestAllFontsHaveMetrics(t *testing.T) {
	fonts := []*Standard14Font{
		Helvetica, HelveticaBold, HelveticaOblique, HelveticaBoldOblique,
		TimesRoman, TimesBold, TimesItalic, TimesBoldItalic,
		Courier, CourierBold, CourierOblique, CourierBoldOblique,
		Symbol, ZapfDingbats,
	}

	for _, f := range fonts {
		m := f.GetMetrics()
		if m == nil {
			t.Errorf("Font %s has nil metrics", f.Name)
			continue
		}

		// Verify basic metrics are set
		if m.Ascender == 0 {
			t.Errorf("Font %s has zero Ascender", f.Name)
		}
		if m.Descender == 0 {
			t.Errorf("Font %s has zero Descender", f.Name)
		}
		if m.CapHeight == 0 {
			t.Errorf("Font %s has zero CapHeight", f.Name)
		}
		if m.XHeight == 0 {
			t.Errorf("Font %s has zero XHeight", f.Name)
		}
		if m.DefaultWidth == 0 {
			t.Errorf("Font %s has zero DefaultWidth", f.Name)
		}
		if len(m.CharWidths) == 0 {
			t.Errorf("Font %s has empty CharWidths", f.Name)
		}
	}
}

// TestFontMetricsGetters tests the getter methods on FontMetrics.
func TestFontMetricsGetters(t *testing.T) {
	m := Helvetica.GetMetrics()
	if m == nil {
		t.Fatal("Helvetica metrics should not be nil")
	}

	// Test Helvetica standard metrics
	if m.GetAscender() != 718 {
		t.Errorf("Helvetica GetAscender() = %d, want 718", m.GetAscender())
	}
	if m.GetDescender() != -207 {
		t.Errorf("Helvetica GetDescender() = %d, want -207", m.GetDescender())
	}
	if m.GetCapHeight() != 718 {
		t.Errorf("Helvetica GetCapHeight() = %d, want 718", m.GetCapHeight())
	}
	if m.GetXHeight() != 523 {
		t.Errorf("Helvetica GetXHeight() = %d, want 523", m.GetXHeight())
	}
}

// TestGetMetricsByName tests looking up metrics by font name string.
func TestGetMetricsByName(t *testing.T) {
	tests := []struct {
		name   string
		expect bool
	}{
		{"Helvetica", true},
		{"Helvetica-Bold", true},
		{"Helvetica-Oblique", true},
		{"Helvetica-BoldOblique", true},
		{"Times-Roman", true},
		{"Times-Bold", true},
		{"Times-Italic", true},
		{"Times-BoldItalic", true},
		{"Courier", true},
		{"Courier-Bold", true},
		{"Courier-Oblique", true},
		{"Courier-BoldOblique", true},
		{"Symbol", true},
		{"ZapfDingbats", true},
		{"UnknownFont", false},
		{"Arial", false}, // Not a standard 14 font
	}

	for _, tt := range tests {
		m := GetMetrics(tt.name)
		if tt.expect && m == nil {
			t.Errorf("GetMetrics(%q) = nil, want non-nil", tt.name)
		}
		if !tt.expect && m != nil {
			t.Errorf("GetMetrics(%q) = non-nil, want nil", tt.name)
		}
	}
}

// TestCourierFamilyMonospace verifies all Courier variants are monospace.
func TestCourierFamilyMonospace(t *testing.T) {
	courierFonts := []*Standard14Font{
		Courier, CourierBold, CourierOblique, CourierBoldOblique,
	}

	for _, f := range courierFonts {
		m := f.GetMetrics()
		if m == nil {
			t.Errorf("Font %s has nil metrics", f.Name)
			continue
		}

		// Check that common characters all have width 600
		for ch := 'A'; ch <= 'Z'; ch++ {
			if w := m.GetCharWidth(ch); w != 600 {
				t.Errorf("%s: GetCharWidth(%q) = %d, want 600", f.Name, ch, w)
			}
		}
	}
}

// TestHelveticaFamilyConsistency tests related fonts share expected widths.
func TestHelveticaFamilyConsistency(t *testing.T) {
	// Helvetica and Helvetica-Oblique should have identical widths
	regular := Helvetica.GetMetrics()
	oblique := HelveticaOblique.GetMetrics()

	if regular == nil || oblique == nil {
		t.Fatal("Metrics should not be nil")
	}

	// Check a few characters
	chars := []rune{'A', 'a', 'M', 'm', '0', '9'}
	for _, ch := range chars {
		r := regular.GetCharWidth(ch)
		o := oblique.GetCharWidth(ch)
		if r != o {
			t.Errorf("Helvetica vs Oblique: GetCharWidth(%q) = %d vs %d", ch, r, o)
		}
	}
}

// TestSymbolicFontMetrics tests that Symbol and ZapfDingbats have metrics.
func TestSymbolicFontMetrics(t *testing.T) {
	// Symbol should have Greek letters
	sm := Symbol.GetMetrics()
	if sm == nil {
		t.Fatal("Symbol metrics should not be nil")
	}

	// Greek alpha
	if w := sm.GetCharWidth('α'); w == sm.DefaultWidth {
		t.Error("Symbol should have explicit width for Greek alpha")
	}

	// ZapfDingbats should have decorative symbols
	zm := ZapfDingbats.GetMetrics()
	if zm == nil {
		t.Fatal("ZapfDingbats metrics should not be nil")
	}

	// Common dingbat - checkmark
	if w := zm.GetCharWidth('✓'); w == 0 {
		t.Error("ZapfDingbats should have width for checkmark")
	}
}

// TestLineHeight tests calculating line height from metrics.
func TestLineHeight(t *testing.T) {
	m := Helvetica.GetMetrics()
	if m == nil {
		t.Fatal("Helvetica metrics should not be nil")
	}

	// Line height = (Ascender - Descender) * size / 1000
	// For Helvetica: (718 - (-207)) * 12 / 1000 = 925 * 12 / 1000 = 11.1
	size := 12.0
	lineHeight := float64(m.Ascender-m.Descender) * size / 1000.0
	expected := 11.1
	if !floatEquals(lineHeight, expected, 0.001) {
		t.Errorf("Line height at 12pt = %f, want %f", lineHeight, expected)
	}
}

// TestExtendedLatinCharacters tests that extended Latin chars are supported.
func TestExtendedLatinCharacters(t *testing.T) {
	m := Helvetica.GetMetrics()
	if m == nil {
		t.Fatal("Helvetica metrics should not be nil")
	}

	// Extended Latin characters commonly used in WinAnsiEncoding
	extendedChars := []rune{'é', 'ü', 'ñ', 'ø', 'ß', '€'}
	for _, ch := range extendedChars {
		w := m.GetCharWidth(ch)
		if w == 0 {
			t.Errorf("Extended char %q should have non-zero width", ch)
		}
	}
}

// floatEquals compares two floats with tolerance.
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
