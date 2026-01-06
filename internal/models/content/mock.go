package content

import (
	"io"

	"github.com/coregx/gxpdf/internal/models/types"
)

// MockContent is a mock implementation of Content for testing.
//
// This is exported so other packages can use it in their tests.
// It provides a simple way to test code that depends on the Content interface.
//
// Example usage:
//
//	mock := &content.MockContent{
//	    BoundsValue: types.MustRectangle(0, 0, 100, 100),
//	    TypeValue:   content.ContentTypeText,
//	}
//	page.AddContent(mock)
type MockContent struct {
	// RenderFunc is called when Render() is invoked.
	// If nil, a default implementation writes "mock content".
	RenderFunc func(w io.Writer) error

	// BoundsValue is returned by Bounds().
	// If zero, a default 100x100 rectangle is returned.
	BoundsValue types.Rectangle

	// ValidateFunc is called when Validate() is invoked.
	// If nil, returns no error.
	ValidateFunc func() error

	// TypeValue is returned by Type().
	// If empty, returns ContentTypeText.
	TypeValue ContentType
}

// Render implements Content.Render().
func (m *MockContent) Render(w io.Writer) error {
	if m.RenderFunc != nil {
		return m.RenderFunc(w)
	}
	_, err := w.Write([]byte("mock content"))
	return err
}

// Bounds implements Content.Bounds().
func (m *MockContent) Bounds() types.Rectangle {
	if m.BoundsValue.Width() == 0 {
		// Return default bounds if not set
		return types.MustRectangle(0, 0, 100, 100)
	}
	return m.BoundsValue
}

// Validate implements Content.Validate().
func (m *MockContent) Validate() error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc()
	}
	return nil
}

// Type implements Content.Type().
func (m *MockContent) Type() ContentType {
	if m.TypeValue == "" {
		return ContentTypeText
	}
	return m.TypeValue
}

// Ensure MockContent implements Content interface at compile time.
var _ Content = (*MockContent)(nil)
