package creator

import (
	"math"
	"testing"
)

// TestColorCMYK_ToRGB tests CMYK to RGB conversion.
func TestColorCMYK_ToRGB(t *testing.T) {
	tests := []struct {
		name     string
		cmyk     ColorCMYK
		expected Color
	}{
		{
			name:     "Pure black (100% K)",
			cmyk:     ColorCMYK{0, 0, 0, 1},
			expected: Color{0, 0, 0},
		},
		{
			name:     "White (no ink)",
			cmyk:     ColorCMYK{0, 0, 0, 0},
			expected: Color{1, 1, 1},
		},
		{
			name:     "Pure cyan",
			cmyk:     ColorCMYK{1, 0, 0, 0},
			expected: Color{0, 1, 1},
		},
		{
			name:     "Pure magenta",
			cmyk:     ColorCMYK{0, 1, 0, 0},
			expected: Color{1, 0, 1},
		},
		{
			name:     "Pure yellow",
			cmyk:     ColorCMYK{0, 0, 1, 0},
			expected: Color{1, 1, 0},
		},
		{
			name:     "Red (M + Y)",
			cmyk:     ColorCMYK{0, 1, 1, 0},
			expected: Color{1, 0, 0},
		},
		{
			name:     "Green (C + Y)",
			cmyk:     ColorCMYK{1, 0, 1, 0},
			expected: Color{0, 1, 0},
		},
		{
			name:     "Blue (C + M)",
			cmyk:     ColorCMYK{1, 1, 0, 0},
			expected: Color{0, 0, 1},
		},
		{
			name:     "50% gray (50% K)",
			cmyk:     ColorCMYK{0, 0, 0, 0.5},
			expected: Color{0.5, 0.5, 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgb := tt.cmyk.ToRGB()
			if !colorsEqual(rgb, tt.expected) {
				t.Errorf("ToRGB() = %+v, expected %+v", rgb, tt.expected)
			}
		})
	}
}

// TestColor_ToCMYK tests RGB to CMYK conversion.
func TestColor_ToCMYK(t *testing.T) {
	tests := []struct {
		name     string
		rgb      Color
		expected ColorCMYK
	}{
		{
			name:     "Pure black",
			rgb:      Color{0, 0, 0},
			expected: ColorCMYK{0, 0, 0, 1},
		},
		{
			name:     "White",
			rgb:      Color{1, 1, 1},
			expected: ColorCMYK{0, 0, 0, 0},
		},
		{
			name:     "Red",
			rgb:      Color{1, 0, 0},
			expected: ColorCMYK{0, 1, 1, 0},
		},
		{
			name:     "Green",
			rgb:      Color{0, 1, 0},
			expected: ColorCMYK{1, 0, 1, 0},
		},
		{
			name:     "Blue",
			rgb:      Color{0, 0, 1},
			expected: ColorCMYK{1, 1, 0, 0},
		},
		{
			name:     "Cyan",
			rgb:      Color{0, 1, 1},
			expected: ColorCMYK{1, 0, 0, 0},
		},
		{
			name:     "Magenta",
			rgb:      Color{1, 0, 1},
			expected: ColorCMYK{0, 1, 0, 0},
		},
		{
			name:     "Yellow",
			rgb:      Color{1, 1, 0},
			expected: ColorCMYK{0, 0, 1, 0},
		},
		{
			name:     "50% gray",
			rgb:      Color{0.5, 0.5, 0.5},
			expected: ColorCMYK{0, 0, 0, 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmyk := tt.rgb.ToCMYK()
			if !cmykColorsEqual(cmyk, tt.expected) {
				t.Errorf("ToCMYK() = %+v, expected %+v", cmyk, tt.expected)
			}
		})
	}
}

// TestCMYK_RGB_Roundtrip tests that converting CMYK -> RGB -> CMYK preserves the color
// for special colors (primaries, black, white, grays).
//
// Note: RGB ↔ CMYK conversion is not perfectly reversible for all colors due to
// different color gamuts. This test only checks colors that should roundtrip correctly.
func TestCMYK_RGB_Roundtrip(t *testing.T) {
	tests := []ColorCMYK{
		{0, 0, 0, 1},   // Black
		{0, 0, 0, 0},   // White
		{1, 0, 0, 0},   // Cyan
		{0, 1, 0, 0},   // Magenta
		{0, 0, 1, 0},   // Yellow
		{0, 1, 1, 0},   // Red
		{1, 0, 1, 0},   // Green
		{1, 1, 0, 0},   // Blue
		{0, 0, 0, 0.5}, // 50% gray
	}

	for _, original := range tests {
		t.Run("CMYK->RGB->CMYK", func(t *testing.T) {
			rgb := original.ToRGB()
			cmyk := rgb.ToCMYK()
			if !cmykColorsEqual(cmyk, original) {
				t.Errorf("Roundtrip failed: original=%+v, rgb=%+v, final=%+v", original, rgb, cmyk)
			}
		})
	}
}

// TestRGB_CMYK_Roundtrip tests that RGB -> CMYK -> RGB preserves the color.
func TestRGB_CMYK_Roundtrip(t *testing.T) {
	tests := []Color{
		{0, 0, 0},         // Black
		{1, 1, 1},         // White
		{1, 0, 0},         // Red
		{0, 1, 0},         // Green
		{0, 0, 1},         // Blue
		{0, 1, 1},         // Cyan
		{1, 0, 1},         // Magenta
		{1, 1, 0},         // Yellow
		{0.5, 0.5, 0.5},   // 50% gray
		{0.25, 0.75, 0.5}, // Mixed color
	}

	for _, original := range tests {
		t.Run("RGB->CMYK->RGB", func(t *testing.T) {
			cmyk := original.ToCMYK()
			rgb := cmyk.ToRGB()
			if !colorsEqual(rgb, original) {
				t.Errorf("Roundtrip failed: original=%+v, cmyk=%+v, final=%+v", original, cmyk, rgb)
			}
		})
	}
}

// TestNewColorCMYK tests the ColorCMYK constructor.
func TestNewColorCMYK(t *testing.T) {
	cmyk := NewColorCMYK(0.5, 0.3, 0.2, 0.1)
	expected := ColorCMYK{C: 0.5, M: 0.3, Y: 0.2, K: 0.1}
	if cmyk != expected {
		t.Errorf("NewColorCMYK() = %+v, expected %+v", cmyk, expected)
	}
}

// TestPredefinedCMYKColors tests the predefined CMYK color constants.
func TestPredefinedCMYKColors(t *testing.T) {
	tests := []struct {
		name     string
		color    ColorCMYK
		expected ColorCMYK
	}{
		{"CMYKBlack", CMYKBlack, ColorCMYK{0, 0, 0, 1}},
		{"CMYKWhite", CMYKWhite, ColorCMYK{0, 0, 0, 0}},
		{"CMYKCyan", CMYKCyan, ColorCMYK{1, 0, 0, 0}},
		{"CMYKMagenta", CMYKMagenta, ColorCMYK{0, 1, 0, 0}},
		{"CMYKYellow", CMYKYellow, ColorCMYK{0, 0, 1, 0}},
		{"CMYKRed", CMYKRed, ColorCMYK{0, 1, 1, 0}},
		{"CMYKGreen", CMYKGreen, ColorCMYK{1, 0, 1, 0}},
		{"CMYKBlue", CMYKBlue, ColorCMYK{1, 1, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color != tt.expected {
				t.Errorf("%s = %+v, expected %+v", tt.name, tt.color, tt.expected)
			}
		})
	}
}

// TestPage_AddTextColorCMYK tests adding CMYK-colored text.
func TestPage_AddTextColorCMYK(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("NewPage() failed: %v", err)
	}

	// Test valid CMYK text
	cyan := NewColorCMYK(1.0, 0.0, 0.0, 0.0)
	err = page.AddTextColorCMYK("CMYK Text", 100, 700, Helvetica, 12, cyan)
	if err != nil {
		t.Errorf("AddTextColorCMYK() failed: %v", err)
	}

	// Verify text operation was added
	ops := page.TextOperations()
	if len(ops) != 1 {
		t.Fatalf("Expected 1 text operation, got %d", len(ops))
	}

	// Verify CMYK color was set
	if ops[0].ColorCMYK == nil {
		t.Errorf("ColorCMYK was not set")
	} else if *ops[0].ColorCMYK != cyan {
		t.Errorf("ColorCMYK = %+v, expected %+v", *ops[0].ColorCMYK, cyan)
	}
}

// TestPage_AddTextColorCMYK_Validation tests validation of CMYK text colors.
func TestPage_AddTextColorCMYK_Validation(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("NewPage() failed: %v", err)
	}

	tests := []struct {
		name      string
		color     ColorCMYK
		expectErr bool
	}{
		{"Valid CMYK", ColorCMYK{0.5, 0.5, 0.5, 0.5}, false},
		{"Invalid C (negative)", ColorCMYK{-0.1, 0, 0, 0}, true},
		{"Invalid C (> 1)", ColorCMYK{1.1, 0, 0, 0}, true},
		{"Invalid M (negative)", ColorCMYK{0, -0.1, 0, 0}, true},
		{"Invalid M (> 1)", ColorCMYK{0, 1.1, 0, 0}, true},
		{"Invalid Y (negative)", ColorCMYK{0, 0, -0.1, 0}, true},
		{"Invalid Y (> 1)", ColorCMYK{0, 0, 1.1, 0}, true},
		{"Invalid K (negative)", ColorCMYK{0, 0, 0, -0.1}, true},
		{"Invalid K (> 1)", ColorCMYK{0, 0, 0, 1.1}, true},
		{"All zeros (white)", ColorCMYK{0, 0, 0, 0}, false},
		{"All ones", ColorCMYK{1, 1, 1, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := page.AddTextColorCMYK("Test", 100, 700, Helvetica, 12, tt.color)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error for invalid color %+v, got nil", tt.color)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error for valid color %+v, got: %v", tt.color, err)
			}
		})
	}
}

// colorsEqual compares two RGB colors with a small tolerance for floating point errors.
func colorsEqual(c1, c2 Color) bool {
	const tolerance = 0.0001
	return math.Abs(c1.R-c2.R) < tolerance &&
		math.Abs(c1.G-c2.G) < tolerance &&
		math.Abs(c1.B-c2.B) < tolerance
}

// cmykColorsEqual compares two CMYK colors with a small tolerance for floating point errors.
func cmykColorsEqual(c1, c2 ColorCMYK) bool {
	const tolerance = 0.0001
	return math.Abs(c1.C-c2.C) < tolerance &&
		math.Abs(c1.M-c2.M) < tolerance &&
		math.Abs(c1.Y-c2.Y) < tolerance &&
		math.Abs(c1.K-c2.K) < tolerance
}
