package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:funlen // Table-driven test with many cases
func TestNewRectangle(t *testing.T) {
	tests := []struct {
		name      string
		llx, lly  float64
		urx, ury  float64
		wantErr   bool
		errTarget error
	}{
		{
			name:    "valid A4 rectangle",
			llx:     0,
			lly:     0,
			urx:     595,
			ury:     842,
			wantErr: false,
		},
		{
			name:    "valid with offset",
			llx:     100,
			lly:     100,
			urx:     200,
			ury:     200,
			wantErr: false,
		},
		{
			name:      "invalid: urx <= llx",
			llx:       100,
			lly:       100,
			urx:       100,
			ury:       200,
			wantErr:   true,
			errTarget: ErrInvalidRectangle,
		},
		{
			name:      "invalid: ury <= lly",
			llx:       100,
			lly:       100,
			urx:       200,
			ury:       100,
			wantErr:   true,
			errTarget: ErrInvalidRectangle,
		},
		{
			name:      "invalid: negative dimensions",
			llx:       200,
			lly:       200,
			urx:       100,
			ury:       100,
			wantErr:   true,
			errTarget: ErrInvalidRectangle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rect, err := NewRectangle(tt.llx, tt.lly, tt.urx, tt.ury)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errTarget)
			} else {
				require.NoError(t, err)
				llx, lly := rect.LowerLeft()
				urx, ury := rect.UpperRight()
				assert.Equal(t, tt.llx, llx)
				assert.Equal(t, tt.lly, lly)
				assert.Equal(t, tt.urx, urx)
				assert.Equal(t, tt.ury, ury)
			}
		})
	}
}

func TestRectangle_Dimensions(t *testing.T) {
	tests := []struct {
		name  string
		rect  Rectangle
		wantW float64
		wantH float64
	}{
		{
			name:  "A4 dimensions",
			rect:  A4,
			wantW: 595.276,
			wantH: 841.890,
		},
		{
			name:  "Letter dimensions",
			rect:  Letter,
			wantW: 612,
			wantH: 792,
		},
		{
			name:  "custom rectangle",
			rect:  MustRectangle(0, 0, 100, 200),
			wantW: 100,
			wantH: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.wantW, tt.rect.Width(), 0.001)
			assert.InDelta(t, tt.wantH, tt.rect.Height(), 0.001)
		})
	}
}

func TestRectangle_Contains(t *testing.T) {
	rect := MustRectangle(0, 0, 100, 100)

	tests := []struct {
		name string
		x, y float64
		want bool
	}{
		{"inside center", 50, 50, true},
		{"on lower-left corner", 0, 0, true},
		{"on upper-right corner", 100, 100, true},
		{"on left edge", 0, 50, true},
		{"outside left", -1, 50, false},
		{"outside right", 101, 50, false},
		{"outside bottom", 50, -1, false},
		{"outside top", 50, 101, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rect.Contains(tt.x, tt.y)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRectangle_WithOffset(t *testing.T) {
	original := MustRectangle(0, 0, 100, 100)
	offset := original.WithOffset(50, 50)

	// Original should be unchanged (immutability)
	assert.Equal(t, 0.0, original.llx)
	assert.Equal(t, 0.0, original.lly)
	assert.Equal(t, 100.0, original.urx)
	assert.Equal(t, 100.0, original.ury)

	// New rectangle should have offset
	assert.Equal(t, 50.0, offset.llx)
	assert.Equal(t, 50.0, offset.lly)
	assert.Equal(t, 150.0, offset.urx)
	assert.Equal(t, 150.0, offset.ury)

	// Dimensions should remain the same
	assert.Equal(t, original.Width(), offset.Width())
	assert.Equal(t, original.Height(), offset.Height())
}

func TestRectangle_Equals(t *testing.T) {
	r1 := MustRectangle(0, 0, 100, 100)
	r2 := MustRectangle(0, 0, 100, 100)
	r3 := MustRectangle(0, 0, 100, 101)

	assert.True(t, r1.Equals(r2), "identical rectangles should be equal")
	assert.False(t, r1.Equals(r3), "different rectangles should not be equal")
}

func TestRectangle_String(t *testing.T) {
	rect := MustRectangle(10, 20, 110, 220)
	str := rect.String()

	assert.Contains(t, str, "10")
	assert.Contains(t, str, "20")
	assert.Contains(t, str, "110")
	assert.Contains(t, str, "220")
}

func TestMustRectangle_Panic(t *testing.T) {
	assert.Panics(t, func() {
		MustRectangle(100, 100, 100, 200) // Invalid: urx <= llx
	})
}

// Example test (doubles as documentation).
func ExampleRectangle_Width() {
	rect := A4
	width := rect.Width()
	println(width) // Approximately 595.276
}

// Benchmark tests.
func BenchmarkNewRectangle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewRectangle(0, 0, 595, 842)
	}
}

func BenchmarkRectangle_Contains(b *testing.B) {
	rect := MustRectangle(0, 0, 100, 100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = rect.Contains(50, 50)
	}
}
