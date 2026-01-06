package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageSize_ToRectangle(t *testing.T) {
	tests := []struct {
		name       string
		pageSize   PageSize
		wantWidth  float64
		wantHeight float64
	}{
		{
			name:       "A4",
			pageSize:   A4,
			wantWidth:  595.0,
			wantHeight: 842.0,
		},
		{
			name:       "A3",
			pageSize:   A3,
			wantWidth:  842.0,
			wantHeight: 1191.0,
		},
		{
			name:       "A5",
			pageSize:   A5,
			wantWidth:  420.0,
			wantHeight: 595.0,
		},
		{
			name:       "Letter",
			pageSize:   Letter,
			wantWidth:  612.0,
			wantHeight: 792.0,
		},
		{
			name:       "Legal",
			pageSize:   Legal,
			wantWidth:  612.0,
			wantHeight: 1008.0,
		},
		{
			name:       "Tabloid",
			pageSize:   Tabloid,
			wantWidth:  792.0,
			wantHeight: 1224.0,
		},
		{
			name:       "Unknown (defaults to A4)",
			pageSize:   PageSize(999),
			wantWidth:  595.0,
			wantHeight: 842.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rect := tt.pageSize.ToRectangle()

			llx, lly := rect.LowerLeft()
			assert.Equal(t, 0.0, llx, "lower-left X should be 0")
			assert.Equal(t, 0.0, lly, "lower-left Y should be 0")
			assert.Equal(t, tt.wantWidth, rect.Width(), "width mismatch")
			assert.Equal(t, tt.wantHeight, rect.Height(), "height mismatch")
		})
	}
}

func TestPageSize_String(t *testing.T) {
	tests := []struct {
		pageSize PageSize
		want     string
	}{
		{A4, "A4"},
		{A3, "A3"},
		{A5, "A5"},
		{Letter, "Letter"},
		{Legal, "Legal"},
		{Tabloid, "Tabloid"},
		{Custom, "Custom"},
		{PageSize(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pageSize.String())
		})
	}
}

func TestCustomPageSize(t *testing.T) {
	tests := []struct {
		name       string
		width      float64
		height     float64
		wantWidth  float64
		wantHeight float64
	}{
		{
			name:       "6x9 inches in points",
			width:      432.0, // 6 * 72
			height:     648.0, // 9 * 72
			wantWidth:  432.0,
			wantHeight: 648.0,
		},
		{
			name:       "custom square",
			width:      500.0,
			height:     500.0,
			wantWidth:  500.0,
			wantHeight: 500.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rect := CustomPageSize(tt.width, tt.height)

			llx, lly := rect.LowerLeft()
			assert.Equal(t, 0.0, llx)
			assert.Equal(t, 0.0, lly)
			assert.Equal(t, tt.wantWidth, rect.Width())
			assert.Equal(t, tt.wantHeight, rect.Height())
		})
	}
}

func TestConversionFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
		convert  func(float64) float64
	}{
		{
			name:     "1 inch to points",
			input:    1.0,
			expected: 72.0,
			convert:  InchesToPoints,
		},
		{
			name:     "8.5 inches to points",
			input:    8.5,
			expected: 612.0,
			convert:  InchesToPoints,
		},
		{
			name:     "72 points to inches",
			input:    72.0,
			expected: 1.0,
			convert:  PointsToInches,
		},
		{
			name:     "25.4 mm to points (approximately 72)",
			input:    25.4,
			expected: 72.0,
			convert:  MMToPoints,
		},
		{
			name:     "210 mm to points (A4 width)",
			input:    210.0,
			expected: 595.27559055118, // 210 * 72/25.4
			convert:  MMToPoints,
		},
		{
			name:     "595 points to mm (approximately 210)",
			input:    595.0,
			expected: 209.90277777778, // 595 * 25.4/72
			convert:  PointsToMM,
		},
		{
			name:     "1 cm to points",
			input:    1.0,
			expected: 28.346456692913, // 72/2.54
			convert:  CMToPoints,
		},
		{
			name:     "72 points to cm",
			input:    72.0,
			expected: 2.54,
			convert:  PointsToCM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.convert(tt.input)
			assert.InDelta(t, tt.expected, result, 0.0001, "conversion mismatch")
		})
	}
}

func TestConversionConstants(t *testing.T) {
	assert.Equal(t, 72.0, PointsPerInch, "1 inch = 72 points")
	assert.InDelta(t, 2.83465, PointsPerMM, 0.00001, "1 mm ≈ 2.83465 points")
	assert.InDelta(t, 28.3465, PointsPerCM, 0.0001, "1 cm ≈ 28.3465 points")
}

func TestRealWorldSizes(t *testing.T) {
	t.Run("A4 dimensions match standard", func(t *testing.T) {
		// A4 is 210mm × 297mm
		widthPt := MMToPoints(210.0)
		heightPt := MMToPoints(297.0)

		// Should be approximately 595×842 points
		assert.InDelta(t, 595.0, widthPt, 1.0, "A4 width")
		assert.InDelta(t, 842.0, heightPt, 1.0, "A4 height")
	})

	t.Run("Letter dimensions match standard", func(t *testing.T) {
		// Letter is 8.5in × 11in
		widthPt := InchesToPoints(8.5)
		heightPt := InchesToPoints(11.0)

		// Should be 612×792 points
		assert.Equal(t, 612.0, widthPt, "Letter width")
		assert.Equal(t, 792.0, heightPt, "Letter height")
	})
}

// Benchmark tests
func BenchmarkPageSize_ToRectangle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = A4.ToRectangle()
	}
}

func BenchmarkCustomPageSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CustomPageSize(595, 842)
	}
}

func BenchmarkInchesToPoints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = InchesToPoints(8.5)
	}
}

func BenchmarkMMToPoints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MMToPoints(210.0)
	}
}
