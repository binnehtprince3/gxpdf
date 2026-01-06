package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphicsElement_String(t *testing.T) {
	elem := &GraphicsElement{
		Type:   GraphicsTypeLine,
		Points: []Point{{X: 0, Y: 0}, {X: 100, Y: 0}},
		Color:  NewColor(0, 0, 0),
		Width:  1.0,
	}

	str := elem.String()
	assert.Contains(t, str, "Line")
	assert.Contains(t, str, "0.00")
	assert.Contains(t, str, "100.00")
}

func TestGraphicsType_String(t *testing.T) {
	tests := []struct {
		typ      GraphicsType
		expected string
	}{
		{GraphicsTypeLine, "Line"},
		{GraphicsTypeRectangle, "Rectangle"},
		{GraphicsTypePath, "Path"},
		{GraphicsType(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.typ.String())
		})
	}
}

func TestPoint_NewPoint(t *testing.T) {
	p := NewPoint(10.5, 20.3)
	assert.Equal(t, 10.5, p.X)
	assert.Equal(t, 20.3, p.Y)
}

func TestPoint_String(t *testing.T) {
	p := NewPoint(10.5, 20.3)
	str := p.String()
	assert.Contains(t, str, "10.50")
	assert.Contains(t, str, "20.30")
}

func TestColor_NewColor(t *testing.T) {
	c := NewColor(0.5, 0.7, 0.9)
	assert.Equal(t, 0.5, c.R)
	assert.Equal(t, 0.7, c.G)
	assert.Equal(t, 0.9, c.B)
}

func TestColor_IsBlack(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expected bool
	}{
		{"pure black", NewColor(0, 0, 0), true},
		{"almost black", NewColor(0.05, 0.05, 0.05), true},
		{"dark gray", NewColor(0.15, 0.15, 0.15), false},
		{"white", NewColor(1, 1, 1), false},
		{"red", NewColor(1, 0, 0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.color.IsBlack())
		})
	}
}

func TestColor_String(t *testing.T) {
	c := NewColor(0.5, 0.7, 0.9)
	str := c.String()
	assert.Contains(t, str, "0.50")
	assert.Contains(t, str, "0.70")
	assert.Contains(t, str, "0.90")
}

func TestNewGraphicsState(t *testing.T) {
	state := NewGraphicsState()

	require.NotNil(t, state)
	assert.NotNil(t, state.CurrentPath)
	assert.Equal(t, 1.0, state.LineWidth)
	assert.True(t, state.StrokeColor.IsBlack())
	assert.True(t, state.FillColor.IsBlack())
}

func TestGraphicsParser_isRectangle(t *testing.T) {
	gp := &GraphicsParser{state: NewGraphicsState()}

	tests := []struct {
		name     string
		points   []Point
		expected bool
	}{
		{
			name: "valid rectangle horizontal first",
			points: []Point{
				{X: 0, Y: 0},
				{X: 100, Y: 0},
				{X: 100, Y: 50},
				{X: 0, Y: 50},
				{X: 0, Y: 0},
			},
			expected: true,
		},
		{
			name: "valid rectangle vertical first",
			points: []Point{
				{X: 0, Y: 0},
				{X: 0, Y: 50},
				{X: 100, Y: 50},
				{X: 100, Y: 0},
				{X: 0, Y: 0},
			},
			expected: true,
		},
		{
			name: "too few points",
			points: []Point{
				{X: 0, Y: 0},
				{X: 100, Y: 0},
			},
			expected: false,
		},
		{
			name: "not closed",
			points: []Point{
				{X: 0, Y: 0},
				{X: 100, Y: 0},
				{X: 100, Y: 50},
				{X: 0, Y: 50},
				{X: 10, Y: 10}, // Not back to start
			},
			expected: false,
		},
		{
			name: "oblique shape",
			points: []Point{
				{X: 0, Y: 0},
				{X: 100, Y: 10},
				{X: 100, Y: 50},
				{X: 0, Y: 40},
				{X: 0, Y: 0},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gp.isRectangle(tt.points)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGraphicsParser_clearPath(t *testing.T) {
	gp := &GraphicsParser{state: NewGraphicsState()}

	gp.state.CurrentPath = []Point{{X: 0, Y: 0}, {X: 100, Y: 100}}
	assert.Len(t, gp.state.CurrentPath, 2)

	gp.clearPath()
	assert.Len(t, gp.state.CurrentPath, 0)
}

func TestGraphicsParser_closePath(t *testing.T) {
	gp := &GraphicsParser{state: NewGraphicsState()}

	gp.state.CurrentPath = []Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 100, Y: 100}}
	gp.closePath()

	assert.Len(t, gp.state.CurrentPath, 4)
	assert.Equal(t, gp.state.CurrentPath[0], gp.state.CurrentPath[3])
}

func TestGraphicsParser_strokePath_line(t *testing.T) {
	gp := &GraphicsParser{
		state:    NewGraphicsState(),
		elements: []*GraphicsElement{},
	}

	// Create a simple line
	gp.state.CurrentPath = []Point{{X: 0, Y: 0}, {X: 100, Y: 0}}
	gp.strokePath()

	require.Len(t, gp.elements, 1)
	assert.Equal(t, GraphicsTypeLine, gp.elements[0].Type)
	assert.Len(t, gp.elements[0].Points, 2)
}

func TestGraphicsParser_strokePath_rectangle(t *testing.T) {
	gp := &GraphicsParser{
		state:    NewGraphicsState(),
		elements: []*GraphicsElement{},
	}

	// Create a rectangle
	gp.state.CurrentPath = []Point{
		{X: 0, Y: 0},
		{X: 100, Y: 0},
		{X: 100, Y: 50},
		{X: 0, Y: 50},
		{X: 0, Y: 0},
	}
	gp.strokePath()

	require.Len(t, gp.elements, 1)
	assert.Equal(t, GraphicsTypeRectangle, gp.elements[0].Type)
	assert.Len(t, gp.elements[0].Points, 5)
}

func TestGraphicsParser_strokePath_multiSegment(t *testing.T) {
	gp := &GraphicsParser{
		state:    NewGraphicsState(),
		elements: []*GraphicsElement{},
	}

	// Create a multi-segment path (not a rectangle)
	gp.state.CurrentPath = []Point{
		{X: 0, Y: 0},
		{X: 50, Y: 0},
		{X: 100, Y: 50},
	}
	gp.strokePath()

	// Should extract 2 line segments
	require.Len(t, gp.elements, 2)
	assert.Equal(t, GraphicsTypeLine, gp.elements[0].Type)
	assert.Equal(t, GraphicsTypeLine, gp.elements[1].Type)
}
