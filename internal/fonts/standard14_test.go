package fonts

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// TestAllStandard14FontsDefined verifies that all 14 standard fonts are defined.
func TestAllStandard14FontsDefined(t *testing.T) {
	fonts := []*Standard14Font{
		// Serif (Times).
		TimesRoman,
		TimesBold,
		TimesItalic,
		TimesBoldItalic,
		// Sans-Serif (Helvetica).
		Helvetica,
		HelveticaBold,
		HelveticaOblique,
		HelveticaBoldOblique,
		// Monospace (Courier).
		Courier,
		CourierBold,
		CourierOblique,
		CourierBoldOblique,
		// Symbolic.
		Symbol,
		ZapfDingbats,
	}

	if len(fonts) != 14 {
		t.Fatalf("expected 14 fonts, got %d", len(fonts))
	}

	// Verify all fonts are non-nil.
	for i, f := range fonts {
		if f == nil {
			t.Errorf("font at index %d is nil", i)
		}
	}
}

// TestStandard14FontProperties verifies the properties of each standard font.
func TestStandard14FontProperties(t *testing.T) {
	tests := []struct {
		font       *Standard14Font
		name       string
		family     string
		weight     string
		style      string
		isSymbolic bool
	}{
		// Serif (Times).
		{TimesRoman, "Times-Roman", "Times", "Regular", "Normal", false},
		{TimesBold, "Times-Bold", "Times", "Bold", "Normal", false},
		{TimesItalic, "Times-Italic", "Times", "Regular", "Italic", false},
		{TimesBoldItalic, "Times-BoldItalic", "Times", "Bold", "Italic", false},
		// Sans-Serif (Helvetica).
		{Helvetica, "Helvetica", "Helvetica", "Regular", "Normal", false},
		{HelveticaBold, "Helvetica-Bold", "Helvetica", "Bold", "Normal", false},
		{HelveticaOblique, "Helvetica-Oblique", "Helvetica", "Regular", "Oblique", false},
		{HelveticaBoldOblique, "Helvetica-BoldOblique", "Helvetica", "Bold", "Oblique", false},
		// Monospace (Courier).
		{Courier, "Courier", "Courier", "Regular", "Normal", false},
		{CourierBold, "Courier-Bold", "Courier", "Bold", "Normal", false},
		{CourierOblique, "Courier-Oblique", "Courier", "Regular", "Oblique", false},
		{CourierBoldOblique, "Courier-BoldOblique", "Courier", "Bold", "Oblique", false},
		// Symbolic.
		{Symbol, "Symbol", "Symbol", "Regular", "Normal", true},
		{ZapfDingbats, "ZapfDingbats", "ZapfDingbats", "Regular", "Normal", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.font.Name != tt.name {
				t.Errorf("Name: got %q, want %q", tt.font.Name, tt.name)
			}
			if tt.font.Family != tt.family {
				t.Errorf("Family: got %q, want %q", tt.font.Family, tt.family)
			}
			if tt.font.Weight != tt.weight {
				t.Errorf("Weight: got %q, want %q", tt.font.Weight, tt.weight)
			}
			if tt.font.Style != tt.style {
				t.Errorf("Style: got %q, want %q", tt.font.Style, tt.style)
			}
			if tt.font.IsSymbolic != tt.isSymbolic {
				t.Errorf("IsSymbolic: got %v, want %v", tt.font.IsSymbolic, tt.isSymbolic)
			}
		})
	}
}

// TestPDFName verifies that PDFName returns the correct PostScript name.
func TestPDFName(t *testing.T) {
	tests := []struct {
		font *Standard14Font
		want string
	}{
		{Helvetica, "Helvetica"},
		{HelveticaBold, "Helvetica-Bold"},
		{TimesRoman, "Times-Roman"},
		{Courier, "Courier"},
		{Symbol, "Symbol"},
		{ZapfDingbats, "ZapfDingbats"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.font.PDFName()
			if got != tt.want {
				t.Errorf("PDFName() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestWriteFontObject_NonSymbolic tests font object generation for non-symbolic fonts.
func TestWriteFontObject_NonSymbolic(t *testing.T) {
	var buf bytes.Buffer
	err := Helvetica.WriteFontObject(5, &buf)
	if err != nil {
		t.Fatalf("WriteFontObject() error = %v", err)
	}

	got := buf.String()

	// Check required components.
	requiredParts := []string{
		"5 0 obj",
		"/Type /Font",
		"/Subtype /Type1",
		"/BaseFont /Helvetica",
		"/Encoding /WinAnsiEncoding",
		">>",
		"endobj",
	}

	for _, part := range requiredParts {
		if !strings.Contains(got, part) {
			t.Errorf("WriteFontObject() missing %q\nGot:\n%s", part, got)
		}
	}
}

// TestWriteFontObject_Symbolic tests font object generation for symbolic fonts.
func TestWriteFontObject_Symbolic(t *testing.T) {
	tests := []struct {
		font *Standard14Font
		name string
	}{
		{Symbol, "Symbol"},
		{ZapfDingbats, "ZapfDingbats"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.font.WriteFontObject(10, &buf)
			if err != nil {
				t.Fatalf("WriteFontObject() error = %v", err)
			}

			got := buf.String()

			// Check required components.
			requiredParts := []string{
				"10 0 obj",
				"/Type /Font",
				"/Subtype /Type1",
				"/BaseFont /" + tt.name,
				">>",
				"endobj",
			}

			for _, part := range requiredParts {
				if !strings.Contains(got, part) {
					t.Errorf("WriteFontObject() missing %q\nGot:\n%s", part, got)
				}
			}

			// Symbolic fonts should NOT have /Encoding.
			if strings.Contains(got, "/Encoding") {
				t.Errorf("WriteFontObject() should not contain /Encoding for symbolic font\nGot:\n%s", got)
			}
		})
	}
}

// TestWriteFontObject_AllFonts tests that all 14 fonts can generate valid objects.
//
//nolint:cyclop // Test function needs to verify all 14 fonts
func TestWriteFontObject_AllFonts(t *testing.T) {
	fonts := []*Standard14Font{
		TimesRoman, TimesBold, TimesItalic, TimesBoldItalic,
		Helvetica, HelveticaBold, HelveticaOblique, HelveticaBoldOblique,
		Courier, CourierBold, CourierOblique, CourierBoldOblique,
		Symbol, ZapfDingbats,
	}

	for i, font := range fonts {
		t.Run(font.Name, func(t *testing.T) {
			var buf bytes.Buffer
			objNum := i + 1

			err := font.WriteFontObject(objNum, &buf)
			if err != nil {
				t.Fatalf("WriteFontObject() error = %v", err)
			}

			got := buf.String()

			// Verify object structure.
			if !strings.HasPrefix(got, string(rune(objNum+'0'))) && objNum < 10 {
				t.Errorf("WriteFontObject() should start with object number %d", objNum)
			}

			if !strings.Contains(got, "endobj") {
				t.Errorf("WriteFontObject() missing endobj")
			}

			// Verify font-specific parts.
			if !strings.Contains(got, "/BaseFont /"+font.Name) {
				t.Errorf("WriteFontObject() missing /BaseFont /%s", font.Name)
			}

			// Check encoding presence.
			hasEncoding := strings.Contains(got, "/Encoding")
			if font.IsSymbolic && hasEncoding {
				t.Errorf("WriteFontObject() symbolic font should not have /Encoding")
			}
			if !font.IsSymbolic && !hasEncoding {
				t.Errorf("WriteFontObject() non-symbolic font should have /Encoding")
			}
		})
	}
}

// TestWriteFontObject_DifferentObjectNumbers tests writing with various object numbers.
func TestWriteFontObject_DifferentObjectNumbers(t *testing.T) {
	tests := []struct {
		name   string
		objNum int
	}{
		{"obj-1", 1},
		{"obj-5", 5},
		{"obj-10", 10},
		{"obj-100", 100},
		{"obj-999", 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Helvetica.WriteFontObject(tt.objNum, &buf)
			if err != nil {
				t.Fatalf("WriteFontObject() error = %v", err)
			}

			got := buf.String()

			// Check object header contains the correct number.
			expectedStart := fmt.Sprintf("%d 0 obj", tt.objNum)
			if !strings.HasPrefix(got, expectedStart) {
				t.Errorf("WriteFontObject() should start with %q, got:\n%s", expectedStart, got)
			}
		})
	}
}

// TestWriteFontObject_OutputFormat verifies the exact format of generated objects.
func TestWriteFontObject_OutputFormat(t *testing.T) {
	tests := []struct {
		font      *Standard14Font
		objNum    int
		wantParts []string
	}{
		{
			font:   Helvetica,
			objNum: 5,
			wantParts: []string{
				"5 0 obj\n",
				"<< /Type /Font\n",
				"/Subtype /Type1\n",
				"/BaseFont /Helvetica\n",
				"/Encoding /WinAnsiEncoding\n",
				">>\n",
				"endobj\n",
			},
		},
		{
			font:   Symbol,
			objNum: 10,
			wantParts: []string{
				"10 0 obj\n",
				"<< /Type /Font\n",
				"/Subtype /Type1\n",
				"/BaseFont /Symbol\n",
				">>\n",
				"endobj\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.font.Name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.font.WriteFontObject(tt.objNum, &buf)
			if err != nil {
				t.Fatalf("WriteFontObject() error = %v", err)
			}

			got := buf.String()

			for _, part := range tt.wantParts {
				if !strings.Contains(got, part) {
					t.Errorf("WriteFontObject() missing part %q\nGot:\n%s", part, got)
				}
			}
		})
	}
}

// TestFontFamilies verifies that fonts are grouped correctly by family.
func TestFontFamilies(t *testing.T) {
	timesFonts := []*Standard14Font{TimesRoman, TimesBold, TimesItalic, TimesBoldItalic}
	helveticaFonts := []*Standard14Font{Helvetica, HelveticaBold, HelveticaOblique, HelveticaBoldOblique}
	courierFonts := []*Standard14Font{Courier, CourierBold, CourierOblique, CourierBoldOblique}

	// Verify Times family.
	for _, f := range timesFonts {
		if f.Family != "Times" {
			t.Errorf("Font %s: expected family Times, got %s", f.Name, f.Family)
		}
	}

	// Verify Helvetica family.
	for _, f := range helveticaFonts {
		if f.Family != "Helvetica" {
			t.Errorf("Font %s: expected family Helvetica, got %s", f.Name, f.Family)
		}
	}

	// Verify Courier family.
	for _, f := range courierFonts {
		if f.Family != "Courier" {
			t.Errorf("Font %s: expected family Courier, got %s", f.Name, f.Family)
		}
	}
}

// TestSymbolicFonts verifies that only Symbol and ZapfDingbats are marked as symbolic.
func TestSymbolicFonts(t *testing.T) {
	allFonts := []*Standard14Font{
		TimesRoman, TimesBold, TimesItalic, TimesBoldItalic,
		Helvetica, HelveticaBold, HelveticaOblique, HelveticaBoldOblique,
		Courier, CourierBold, CourierOblique, CourierBoldOblique,
		Symbol, ZapfDingbats,
	}

	symbolicCount := 0
	for _, f := range allFonts {
		if f.IsSymbolic {
			symbolicCount++
			if f.Name != "Symbol" && f.Name != "ZapfDingbats" {
				t.Errorf("Font %s is marked as symbolic but should not be", f.Name)
			}
		}
	}

	if symbolicCount != 2 {
		t.Errorf("Expected exactly 2 symbolic fonts, got %d", symbolicCount)
	}

	// Verify Symbol and ZapfDingbats are symbolic.
	if !Symbol.IsSymbolic {
		t.Error("Symbol font should be marked as symbolic")
	}
	if !ZapfDingbats.IsSymbolic {
		t.Error("ZapfDingbats font should be marked as symbolic")
	}
}
