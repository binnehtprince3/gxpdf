package creator

import (
	"testing"
)

func TestDrawPolygon(t *testing.T) {
	tests := []struct {
		name        string
		vertices    []Point
		opts        *PolygonOptions
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid triangle with stroke and fill",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
			},
			opts: &PolygonOptions{
				StrokeColor: &Black,
				StrokeWidth: 2.0,
				FillColor:   &Blue,
			},
			expectError: false,
		},
		{
			name: "valid pentagon with fill only",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
				{X: 175, Y: 150},
				{X: 125, Y: 150},
			},
			opts: &PolygonOptions{
				FillColor: &Green,
			},
			expectError: false,
		},
		{
			name: "valid polygon with stroke only",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 200, Y: 100},
				{X: 200, Y: 200},
				{X: 100, Y: 200},
			},
			opts: &PolygonOptions{
				StrokeColor: &Red,
				StrokeWidth: 1.0,
			},
			expectError: false,
		},
		{
			name: "valid polygon with dashed stroke",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
			},
			opts: &PolygonOptions{
				StrokeColor: &Black,
				StrokeWidth: 1.0,
				Dashed:      true,
				DashArray:   []float64{5, 3},
			},
			expectError: false,
		},
		{
			name:        "nil options",
			vertices:    []Point{{X: 100, Y: 100}, {X: 150, Y: 50}, {X: 200, Y: 100}},
			opts:        nil,
			expectError: true,
			errorMsg:    "polygon options cannot be nil",
		},
		{
			name:        "too few vertices (2 vertices)",
			vertices:    []Point{{X: 100, Y: 100}, {X: 150, Y: 50}},
			opts:        &PolygonOptions{StrokeColor: &Black},
			expectError: true,
			errorMsg:    "polygon must have at least 3 vertices",
		},
		{
			name:        "too few vertices (1 vertex)",
			vertices:    []Point{{X: 100, Y: 100}},
			opts:        &PolygonOptions{StrokeColor: &Black},
			expectError: true,
			errorMsg:    "polygon must have at least 3 vertices",
		},
		{
			name:        "empty vertices",
			vertices:    []Point{},
			opts:        &PolygonOptions{StrokeColor: &Black},
			expectError: true,
			errorMsg:    "polygon must have at least 3 vertices",
		},
		{
			name: "neither stroke nor fill",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
			},
			opts:        &PolygonOptions{},
			expectError: true,
			errorMsg:    "polygon must have at least stroke, fill color, or gradient",
		},
		{
			name: "invalid stroke color (R > 1)",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
			},
			opts: &PolygonOptions{
				StrokeColor: &Color{R: 1.5, G: 0, B: 0},
			},
			expectError: true,
			errorMsg:    "stroke color components must be in range [0.0, 1.0]",
		},
		{
			name: "invalid fill color (B < 0)",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
			},
			opts: &PolygonOptions{
				FillColor: &Color{R: 0, G: 0, B: -0.1},
			},
			expectError: true,
			errorMsg:    "fill color components must be in range [0.0, 1.0]",
		},
		{
			name: "negative stroke width",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
			},
			opts: &PolygonOptions{
				StrokeColor: &Black,
				StrokeWidth: -1.0,
			},
			expectError: true,
			errorMsg:    "stroke width must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, err := c.NewPage()
			if err != nil {
				t.Fatalf("failed to create page: %v", err)
			}

			err = page.DrawPolygon(tt.vertices, tt.opts)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// Verify the operation was added
				ops := page.GraphicsOperations()
				if len(ops) != 1 {
					t.Errorf("expected 1 graphics operation, got %d", len(ops))
				} else {
					op := ops[0]
					if op.Type != GraphicsOpPolygon {
						t.Errorf("expected polygon operation, got type %d", op.Type)
					}
					if len(op.Vertices) != len(tt.vertices) {
						t.Errorf("expected %d vertices, got %d", len(tt.vertices), len(op.Vertices))
					}
				}
			}
		})
	}
}

func TestPolygonComplexShapes(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Test star shape (10 vertices)
	starVertices := []Point{
		{X: 150, Y: 50},  // Top
		{X: 165, Y: 90},  // Inner
		{X: 200, Y: 100}, // Right upper
		{X: 170, Y: 120}, // Inner
		{X: 180, Y: 160}, // Right lower
		{X: 150, Y: 135}, // Inner
		{X: 120, Y: 160}, // Left lower
		{X: 130, Y: 120}, // Inner
		{X: 100, Y: 100}, // Left upper
		{X: 135, Y: 90},  // Inner
	}

	opts := &PolygonOptions{
		StrokeColor: &Black,
		StrokeWidth: 1.5,
		FillColor:   &Yellow,
	}

	err = page.DrawPolygon(starVertices, opts)
	if err != nil {
		t.Errorf("failed to draw star polygon: %v", err)
	}

	// Verify 10 vertices
	ops := page.GraphicsOperations()
	if len(ops) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(ops))
	}
	if len(ops[0].Vertices) != 10 {
		t.Errorf("expected 10 vertices, got %d", len(ops[0].Vertices))
	}
}
