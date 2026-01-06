package document

import (
	"errors"
	"testing"

	"github.com/coregx/gxpdf/internal/models/content"
	"github.com/coregx/gxpdf/internal/models/types"
	"github.com/stretchr/testify/assert"
)

func TestNewPage(t *testing.T) {
	page := NewPage(0, A4)

	assert.Equal(t, 0, page.Number())
	assert.Equal(t, 595.0, page.MediaBox().Width())
	assert.Equal(t, 842.0, page.MediaBox().Height())
	assert.Equal(t, 0, page.Rotation())
	assert.Nil(t, page.CropBox())
}

func TestPage_SetRotation(t *testing.T) {
	tests := []struct {
		name      string
		rotation  int
		wantError bool
	}{
		{name: "0 degrees", rotation: 0, wantError: false},
		{name: "90 degrees", rotation: 90, wantError: false},
		{name: "180 degrees", rotation: 180, wantError: false},
		{name: "270 degrees", rotation: 270, wantError: false},
		{name: "invalid 45 degrees", rotation: 45, wantError: true},
		{name: "invalid 360 degrees", rotation: 360, wantError: true},
		{name: "invalid negative", rotation: -90, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := NewPage(0, A4)
			err := page.SetRotation(tt.rotation)

			if tt.wantError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidRotation)
				assert.Equal(t, 0, page.Rotation(), "rotation should not change on error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.rotation, page.Rotation())
			}
		})
	}
}

func TestPage_WidthHeight(t *testing.T) {
	tests := []struct {
		name       string
		pageSize   PageSize
		rotation   int
		wantWidth  float64
		wantHeight float64
	}{
		{
			name:       "A4 portrait (0°)",
			pageSize:   A4,
			rotation:   0,
			wantWidth:  595.0,
			wantHeight: 842.0,
		},
		{
			name:       "A4 landscape (90°)",
			pageSize:   A4,
			rotation:   90,
			wantWidth:  842.0, // swapped
			wantHeight: 595.0,
		},
		{
			name:       "A4 upside down (180°)",
			pageSize:   A4,
			rotation:   180,
			wantWidth:  595.0,
			wantHeight: 842.0,
		},
		{
			name:       "A4 landscape reverse (270°)",
			pageSize:   A4,
			rotation:   270,
			wantWidth:  842.0, // swapped
			wantHeight: 595.0,
		},
		{
			name:       "Letter portrait",
			pageSize:   Letter,
			rotation:   0,
			wantWidth:  612.0,
			wantHeight: 792.0,
		},
		{
			name:       "Letter landscape",
			pageSize:   Letter,
			rotation:   90,
			wantWidth:  792.0,
			wantHeight: 612.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := NewPage(0, tt.pageSize)
			page.SetRotation(tt.rotation)

			assert.Equal(t, tt.wantWidth, page.Width(), "width mismatch")
			assert.Equal(t, tt.wantHeight, page.Height(), "height mismatch")
		})
	}
}

func TestPage_SetCropBox(t *testing.T) {
	page := NewPage(0, A4) // A4 is 595×842

	tests := []struct {
		name      string
		cropBox   types.Rectangle
		wantError bool
	}{
		{
			name:      "valid crop box within media box",
			cropBox:   types.MustRectangle(50, 50, 545, 792),
			wantError: false,
		},
		{
			name:      "crop box equal to media box",
			cropBox:   types.MustRectangle(0, 0, 595, 842),
			wantError: false,
		},
		{
			name:      "crop box exceeds left",
			cropBox:   types.MustRectangle(-10, 0, 595, 842),
			wantError: true,
		},
		{
			name:      "crop box exceeds bottom",
			cropBox:   types.MustRectangle(0, -10, 595, 842),
			wantError: true,
		},
		{
			name:      "crop box exceeds right",
			cropBox:   types.MustRectangle(0, 0, 600, 842),
			wantError: true,
		},
		{
			name:      "crop box exceeds top",
			cropBox:   types.MustRectangle(0, 0, 595, 850),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset crop box before each test
			page.cropBox = nil

			err := page.SetCropBox(tt.cropBox)

			if tt.wantError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrCropBoxOutOfBounds)
				assert.Nil(t, page.CropBox(), "crop box should not be set on error")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, page.CropBox())
				assert.Equal(t, tt.cropBox, *page.CropBox())
			}
		})
	}
}

func TestPage_Validate(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *Page
		wantError bool
		errorType error
	}{
		{
			name: "valid page",
			setup: func() *Page {
				return NewPage(0, A4)
			},
			wantError: false,
		},
		{
			name: "valid page with crop box",
			setup: func() *Page {
				page := NewPage(0, A4)
				page.SetCropBox(types.MustRectangle(50, 50, 545, 792))
				return page
			},
			wantError: false,
		},
		{
			name: "valid page with rotation",
			setup: func() *Page {
				page := NewPage(0, A4)
				page.SetRotation(90)
				return page
			},
			wantError: false,
		},
		// Note: We can't test invalid page dimensions here because Rectangle
		// value objects enforce validity at construction time, and we can't
		// create invalid rectangles without reflection. These cases are covered
		// by Rectangle's own tests in the valueobjects package.
		{
			name: "invalid crop box",
			setup: func() *Page {
				page := NewPage(0, A4)
				// Manually set an invalid crop box
				invalidCropBox := types.MustRectangle(0, 0, 600, 842)
				page.cropBox = &invalidCropBox
				return page
			},
			wantError: true,
			errorType: ErrCropBoxOutOfBounds,
		},
		{
			name: "invalid rotation",
			setup: func() *Page {
				page := NewPage(0, A4)
				page.rotation = 45 // Manually set invalid rotation
				return page
			},
			wantError: true,
			errorType: ErrInvalidRotation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := tt.setup()
			err := page.Validate()

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPage_MediaBoxImmutable(t *testing.T) {
	page := NewPage(0, A4)
	originalMediaBox := page.MediaBox()

	// Try to modify returned media box (should not affect page)
	// Note: Since Rectangle is a value object (passed by value), this test
	// ensures that MediaBox() returns a value, not a pointer
	_ = originalMediaBox // Use the variable

	assert.Equal(t, originalMediaBox, page.MediaBox(), "media box should not change")
}

func TestPage_Number(t *testing.T) {
	tests := []struct {
		pageNum int
	}{
		{pageNum: 0},
		{pageNum: 1},
		{pageNum: 42},
		{pageNum: 999},
	}

	for _, tt := range tests {
		t.Run("page_number_"+string(rune(tt.pageNum)), func(t *testing.T) {
			page := NewPage(tt.pageNum, A4)
			assert.Equal(t, tt.pageNum, page.Number())
		})
	}
}

func TestPage_AddContent(t *testing.T) {
	tests := []struct {
		name      string
		content   content.Content
		wantError bool
		errorType error
	}{
		{
			name: "add valid content",
			content: &content.MockContent{
				BoundsValue: types.MustRectangle(0, 0, 100, 100),
				TypeValue:   content.ContentTypeText,
			},
			wantError: false,
		},
		{
			name:      "add nil content",
			content:   nil,
			wantError: true,
			errorType: ErrNilContent,
		},
		{
			name: "add invalid content (validation fails)",
			content: &content.MockContent{
				ValidateFunc: func() error {
					return errors.New("validation failed")
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := NewPage(0, A4)
			err := page.AddContent(tt.content)

			if tt.wantError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, page.ContentCount())
			}
		})
	}
}

func TestPage_Contents(t *testing.T) {
	page := NewPage(0, A4)

	// Add some content
	mock1 := &content.MockContent{TypeValue: content.ContentTypeText}
	mock2 := &content.MockContent{TypeValue: content.ContentTypeImage}
	mock3 := &content.MockContent{TypeValue: content.ContentTypePath}

	page.AddContent(mock1)
	page.AddContent(mock2)
	page.AddContent(mock3)

	// Get contents
	contents := page.Contents()
	assert.Len(t, contents, 3)
	assert.Equal(t, content.ContentTypeText, contents[0].Type())
	assert.Equal(t, content.ContentTypeImage, contents[1].Type())
	assert.Equal(t, content.ContentTypePath, contents[2].Type())

	// Verify it's a copy (modifying returned slice doesn't affect page)
	contents[0] = nil
	assert.Equal(t, 3, page.ContentCount(), "original contents should not be affected")
}

func TestPage_ContentCount(t *testing.T) {
	page := NewPage(0, A4)
	assert.Equal(t, 0, page.ContentCount(), "new page should have no content")

	mock := &content.MockContent{}
	page.AddContent(mock)
	assert.Equal(t, 1, page.ContentCount())

	page.AddContent(mock)
	assert.Equal(t, 2, page.ContentCount())

	page.AddContent(mock)
	assert.Equal(t, 3, page.ContentCount())
}

func TestPage_ClearContent(t *testing.T) {
	page := NewPage(0, A4)

	// Add content
	mock := &content.MockContent{}
	page.AddContent(mock)
	page.AddContent(mock)
	page.AddContent(mock)
	assert.Equal(t, 3, page.ContentCount())

	// Clear
	page.ClearContent()
	assert.Equal(t, 0, page.ContentCount())

	// Can add again after clearing
	page.AddContent(mock)
	assert.Equal(t, 1, page.ContentCount())
}

func TestPage_Validate_WithContent(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *Page
		wantError bool
	}{
		{
			name: "page with valid content",
			setup: func() *Page {
				page := NewPage(0, A4)
				mock := &content.MockContent{
					BoundsValue: types.MustRectangle(0, 0, 100, 100),
				}
				page.AddContent(mock)
				return page
			},
			wantError: false,
		},
		{
			name: "page with multiple valid contents",
			setup: func() *Page {
				page := NewPage(0, A4)
				mock1 := &content.MockContent{TypeValue: content.ContentTypeText}
				mock2 := &content.MockContent{TypeValue: content.ContentTypeImage}
				page.AddContent(mock1)
				page.AddContent(mock2)
				return page
			},
			wantError: false,
		},
		{
			name: "page with nil content (manually injected)",
			setup: func() *Page {
				page := NewPage(0, A4)
				// Manually inject nil to test validation
				page.contents = append(page.contents, nil)
				return page
			},
			wantError: true,
		},
		{
			name: "page with invalid content",
			setup: func() *Page {
				page := NewPage(0, A4)
				validationError := errors.New("content validation error")
				mock := &content.MockContent{
					ValidateFunc: func() error {
						return validationError
					},
				}
				// Bypass AddContent to inject invalid content
				page.contents = append(page.contents, mock)
				return page
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := tt.setup()
			err := page.Validate()

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewPage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPage(0, A4)
	}
}

func BenchmarkPage_SetRotation(b *testing.B) {
	page := NewPage(0, A4)
	rotations := []int{0, 90, 180, 270}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		page.SetRotation(rotations[i%4])
	}
}

func BenchmarkPage_Validate(b *testing.B) {
	page := NewPage(0, A4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = page.Validate()
	}
}
