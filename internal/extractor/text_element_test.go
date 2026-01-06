package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTextElement(t *testing.T) {
	elem := NewTextElement("Hello", 100, 200, 50, 12, "Helvetica", 12)

	assert.Equal(t, "Hello", elem.Text)
	assert.Equal(t, 100.0, elem.X)
	assert.Equal(t, 200.0, elem.Y)
	assert.Equal(t, 50.0, elem.Width)
	assert.Equal(t, 12.0, elem.Height)
	assert.Equal(t, "Helvetica", elem.FontName)
	assert.Equal(t, 12.0, elem.FontSize)
}

func TestTextElement_Boundaries(t *testing.T) {
	elem := NewTextElement("Test", 100, 200, 50, 12, "Arial", 12)

	assert.Equal(t, 100.0, elem.Left())
	assert.Equal(t, 150.0, elem.Right())
	assert.Equal(t, 200.0, elem.Bottom())
	assert.Equal(t, 212.0, elem.Top())
	assert.Equal(t, 125.0, elem.CenterX())
	assert.Equal(t, 206.0, elem.CenterY())
}

func TestTextElement_String(t *testing.T) {
	elem := NewTextElement("Test", 100, 200, 50, 12, "Arial", 12)
	str := elem.String()

	assert.Contains(t, str, "Test")
	assert.Contains(t, str, "100.00")
	assert.Contains(t, str, "200.00")
}

func TestNewRectangle(t *testing.T) {
	rect := NewRectangle(10, 20, 100, 50)

	assert.Equal(t, 10.0, rect.X)
	assert.Equal(t, 20.0, rect.Y)
	assert.Equal(t, 100.0, rect.Width)
	assert.Equal(t, 50.0, rect.Height)
}

func TestRectangle_Boundaries(t *testing.T) {
	rect := NewRectangle(10, 20, 100, 50)

	assert.Equal(t, 10.0, rect.Left())
	assert.Equal(t, 110.0, rect.Right())
	assert.Equal(t, 20.0, rect.Bottom())
	assert.Equal(t, 70.0, rect.Top())
}

func TestRectangle_Contains(t *testing.T) {
	rect := NewRectangle(10, 20, 100, 50)

	tests := []struct {
		name     string
		x, y     float64
		expected bool
	}{
		{"inside", 50, 40, true},
		{"on left edge", 10, 40, true},
		{"on right edge", 110, 40, true},
		{"on bottom edge", 50, 20, true},
		{"on top edge", 50, 70, true},
		{"outside left", 5, 40, false},
		{"outside right", 115, 40, false},
		{"outside bottom", 50, 15, false},
		{"outside top", 50, 75, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, rect.Contains(tt.x, tt.y))
		})
	}
}

func TestNewTextChunk(t *testing.T) {
	elements := []*TextElement{
		NewTextElement("Hello", 100, 200, 30, 12, "Arial", 12),
		NewTextElement(" ", 130, 200, 5, 12, "Arial", 12),
		NewTextElement("World", 135, 200, 30, 12, "Arial", 12),
	}

	chunk := NewTextChunk(elements)

	assert.Equal(t, 3, chunk.Len())
	assert.Equal(t, "Hello World", chunk.Text())
	assert.Equal(t, 100.0, chunk.Bounds.X)
	assert.Equal(t, 200.0, chunk.Bounds.Y)
	assert.Equal(t, 65.0, chunk.Bounds.Width) // 165 - 100
	assert.Equal(t, 12.0, chunk.Bounds.Height)
}

func TestTextChunk_Add(t *testing.T) {
	chunk := NewTextChunk([]*TextElement{})

	assert.Equal(t, 0, chunk.Len())

	elem1 := NewTextElement("Hello", 100, 200, 30, 12, "Arial", 12)
	chunk.Add(elem1)

	assert.Equal(t, 1, chunk.Len())
	assert.Equal(t, "Hello", chunk.Text())

	elem2 := NewTextElement("World", 140, 200, 30, 12, "Arial", 12)
	chunk.Add(elem2)

	assert.Equal(t, 2, chunk.Len())
	assert.Equal(t, "HelloWorld", chunk.Text())
	assert.Equal(t, 100.0, chunk.Bounds.X)
	assert.Equal(t, 70.0, chunk.Bounds.Width) // 170 - 100
}

func TestTextChunk_EmptyElements(t *testing.T) {
	chunk := NewTextChunk([]*TextElement{})

	assert.Equal(t, 0, chunk.Len())
	assert.Equal(t, "", chunk.Text())
	assert.Equal(t, 0.0, chunk.Bounds.Width)
	assert.Equal(t, 0.0, chunk.Bounds.Height)
}

func TestTextChunk_String(t *testing.T) {
	elements := []*TextElement{
		NewTextElement("Test", 100, 200, 30, 12, "Arial", 12),
	}
	chunk := NewTextChunk(elements)
	str := chunk.String()

	assert.Contains(t, str, "Test")
	assert.Contains(t, str, "elements=1")
}
