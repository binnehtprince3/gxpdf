package content

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/coregx/gxpdf/internal/models/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockContent(t *testing.T) {
	t.Run("default behavior", func(t *testing.T) {
		mock := &MockContent{}

		// Test Render
		var buf bytes.Buffer
		err := mock.Render(&buf)
		require.NoError(t, err)
		assert.Equal(t, "mock content", buf.String())

		// Test Bounds
		bounds := mock.Bounds()
		assert.Equal(t, 100.0, bounds.Width())
		assert.Equal(t, 100.0, bounds.Height())

		// Test Validate
		err = mock.Validate()
		assert.NoError(t, err)

		// Test Type
		assert.Equal(t, ContentTypeText, mock.Type())
	})

	t.Run("custom behavior", func(t *testing.T) {
		customBounds := types.MustRectangle(10, 20, 110, 120)
		renderError := errors.New("render failed")
		validateError := errors.New("validation failed")

		mock := &MockContent{
			RenderFunc: func(w io.Writer) error {
				return renderError
			},
			BoundsValue: customBounds,
			ValidateFunc: func() error {
				return validateError
			},
			TypeValue: ContentTypeImage,
		}

		// Test custom Render
		var buf bytes.Buffer
		err := mock.Render(&buf)
		assert.ErrorIs(t, err, renderError)

		// Test custom Bounds
		bounds := mock.Bounds()
		assert.Equal(t, customBounds, bounds)

		// Test custom Validate
		err = mock.Validate()
		assert.ErrorIs(t, err, validateError)

		// Test custom Type
		assert.Equal(t, ContentTypeImage, mock.Type())
	})
}

func TestContentType_String(t *testing.T) {
	tests := []struct {
		name string
		ct   ContentType
		want string
	}{
		{
			name: "text",
			ct:   ContentTypeText,
			want: "text",
		},
		{
			name: "image",
			ct:   ContentTypeImage,
			want: "image",
		},
		{
			name: "path",
			ct:   ContentTypePath,
			want: "path",
		},
		{
			name: "table",
			ct:   ContentTypeTable,
			want: "table",
		},
		{
			name: "form",
			ct:   ContentTypeForm,
			want: "form",
		},
		{
			name: "annotation",
			ct:   ContentTypeAnnotation,
			want: "annotation",
		},
		{
			name: "custom",
			ct:   ContentType("custom"),
			want: "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ct.String())
		})
	}
}

func TestContentType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		ct    ContentType
		valid bool
	}{
		{
			name:  "text is valid",
			ct:    ContentTypeText,
			valid: true,
		},
		{
			name:  "image is valid",
			ct:    ContentTypeImage,
			valid: true,
		},
		{
			name:  "path is valid",
			ct:    ContentTypePath,
			valid: true,
		},
		{
			name:  "table is valid",
			ct:    ContentTypeTable,
			valid: true,
		},
		{
			name:  "form is valid",
			ct:    ContentTypeForm,
			valid: true,
		},
		{
			name:  "annotation is valid",
			ct:    ContentTypeAnnotation,
			valid: true,
		},
		{
			name:  "empty is invalid",
			ct:    ContentType(""),
			valid: false,
		},
		{
			name:  "unknown is invalid",
			ct:    ContentType("unknown"),
			valid: false,
		},
		{
			name:  "custom is invalid",
			ct:    ContentType("custom"),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.ct.IsValid())
		})
	}
}

func TestContentInterface(t *testing.T) {
	// This test verifies that our mock correctly implements the Content interface
	var _ Content = (*MockContent)(nil)

	t.Run("interface compliance", func(t *testing.T) {
		mock := &MockContent{
			BoundsValue: types.MustRectangle(50, 50, 150, 150),
			TypeValue:   ContentTypePath,
		}

		// Verify all interface methods work
		var buf bytes.Buffer
		err := mock.Render(&buf)
		assert.NoError(t, err)

		bounds := mock.Bounds()
		assert.NotNil(t, bounds)
		assert.Equal(t, 100.0, bounds.Width())

		err = mock.Validate()
		assert.NoError(t, err)

		contentType := mock.Type()
		assert.Equal(t, ContentTypePath, contentType)
		assert.True(t, contentType.IsValid())
	})
}

// Example test showing how to use MockContent in other packages
func ExampleMockContent() {
	// Create a mock content element
	mock := &MockContent{
		BoundsValue: types.MustRectangle(100, 200, 300, 400),
		TypeValue:   ContentTypeText,
	}

	// Use it as Content interface
	var content Content = mock

	// Test rendering
	var buf bytes.Buffer
	_ = content.Render(&buf)

	// Get bounds
	bounds := content.Bounds()
	_ = bounds.Width()  // 200
	_ = bounds.Height() // 200

	// Validate
	_ = content.Validate()

	// Get type
	_ = content.Type() // ContentTypeText
}

// Benchmark tests
func BenchmarkContentType_String(b *testing.B) {
	ct := ContentTypeText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ct.String()
	}
}

func BenchmarkContentType_IsValid(b *testing.B) {
	ct := ContentTypeText
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ct.IsValid()
	}
}

func BenchmarkMockContent_Render(b *testing.B) {
	mock := &MockContent{}
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = mock.Render(&buf)
	}
}
