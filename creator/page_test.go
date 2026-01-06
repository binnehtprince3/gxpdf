package creator

import (
	"testing"

	"github.com/coregx/gxpdf/internal/document"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPage_SetRotation(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Default rotation is 0
	assert.Equal(t, 0, page.Rotation())

	// Set valid rotations
	err = page.SetRotation(90)
	require.NoError(t, err)
	assert.Equal(t, 90, page.Rotation())

	err = page.SetRotation(180)
	require.NoError(t, err)
	assert.Equal(t, 180, page.Rotation())

	err = page.SetRotation(270)
	require.NoError(t, err)
	assert.Equal(t, 270, page.Rotation())

	err = page.SetRotation(0)
	require.NoError(t, err)
	assert.Equal(t, 0, page.Rotation())

	// Invalid rotation
	err = page.SetRotation(45)
	assert.Error(t, err)
}

func TestPage_Rotate(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Default rotation is 0
	assert.Equal(t, 0, page.Rotation())

	// Rotate to landscape (90 degrees)
	err = page.Rotate(90)
	require.NoError(t, err)
	assert.Equal(t, 90, page.Rotation())

	// Rotate to upside down (180 degrees)
	err = page.Rotate(180)
	require.NoError(t, err)
	assert.Equal(t, 180, page.Rotation())

	// Rotate to landscape reverse (270 degrees)
	err = page.Rotate(270)
	require.NoError(t, err)
	assert.Equal(t, 270, page.Rotation())

	// Back to portrait (0 degrees)
	err = page.Rotate(0)
	require.NoError(t, err)
	assert.Equal(t, 0, page.Rotation())

	// Invalid rotation values should fail
	testCases := []struct {
		name     string
		rotation int
	}{
		{"45 degrees", 45},
		{"135 degrees", 135},
		{"360 degrees", 360},
		{"negative", -90},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := page.Rotate(tc.rotation)
			assert.Error(t, err)
		})
	}
}

func TestPage_Dimensions(t *testing.T) {
	c := New()
	page, err := c.NewPage() // A4
	require.NoError(t, err)

	// Portrait (default)
	assert.Equal(t, 595.0, page.Width())
	assert.Equal(t, 842.0, page.Height())

	// Landscape (90 degrees)
	err = page.SetRotation(90)
	require.NoError(t, err)
	assert.Equal(t, 842.0, page.Width())  // Swapped
	assert.Equal(t, 595.0, page.Height()) // Swapped

	// Upside down (180 degrees)
	err = page.SetRotation(180)
	require.NoError(t, err)
	assert.Equal(t, 595.0, page.Width())
	assert.Equal(t, 842.0, page.Height())

	// Landscape (270 degrees)
	err = page.SetRotation(270)
	require.NoError(t, err)
	assert.Equal(t, 842.0, page.Width())  // Swapped
	assert.Equal(t, 595.0, page.Height()) // Swapped
}

func TestPage_SetMargins(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	// Default margins (from creator)
	margins := page.Margins()
	assert.Equal(t, 72.0, margins.Top)
	assert.Equal(t, 72.0, margins.Right)
	assert.Equal(t, 72.0, margins.Bottom)
	assert.Equal(t, 72.0, margins.Left)

	// Set page-specific margins
	err = page.SetMargins(36, 36, 36, 36)
	require.NoError(t, err)

	margins = page.Margins()
	assert.Equal(t, 36.0, margins.Top)
	assert.Equal(t, 36.0, margins.Right)
	assert.Equal(t, 36.0, margins.Bottom)
	assert.Equal(t, 36.0, margins.Left)
}

func TestPage_SetMargins_Negative(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	require.NoError(t, err)

	err = page.SetMargins(-10, 0, 0, 0)
	assert.ErrorIs(t, err, ErrInvalidMargins)
}

func TestPage_ContentDimensions(t *testing.T) {
	c := New()

	// Create page with custom margins
	err := c.SetMargins(50, 40, 30, 60)
	require.NoError(t, err)

	page, err := c.NewPage() // A4
	require.NoError(t, err)

	// Content width = page width - left - right = 595 - 60 - 40 = 495
	assert.Equal(t, 495.0, page.ContentWidth())

	// Content height = page height - top - bottom = 842 - 50 - 30 = 762
	assert.Equal(t, 762.0, page.ContentHeight())
}

func TestPageSize_String(t *testing.T) {
	tests := []struct {
		size PageSize
		want string
	}{
		{A4, "A4"},
		{Letter, "Letter"},
		{Legal, "Legal"},
		{Tabloid, "Tabloid"},
		{A3, "A3"},
		{A5, "A5"},
		{B4, "B4"},
		{B5, "B5"},
		{PageSize(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.size.String())
		})
	}
}

func TestPageSize_ToDomainSize(t *testing.T) {
	tests := []struct {
		size       PageSize
		domainSize document.PageSize
	}{
		{A4, document.A4},
		{Letter, document.Letter},
		{Legal, document.Legal},
		{Tabloid, document.Tabloid},
		{A3, document.A3},
		{A5, document.A5},
		{B4, document.B4},
		{B5, document.B5},
	}

	for _, tt := range tests {
		t.Run(tt.size.String(), func(t *testing.T) {
			assert.Equal(t, tt.domainSize, tt.size.toDomainSize())
		})
	}
}
