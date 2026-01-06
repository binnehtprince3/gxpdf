package creator

import (
	"testing"
)

func TestDrawPolyline(t *testing.T) {
	tests := []struct {
		name        string
		vertices    []Point
		opts        *PolylineOptions
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid polyline with 3 vertices",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
			},
			opts: &PolylineOptions{
				Color: Red,
				Width: 2.0,
			},
			expectError: false,
		},
		{
			name: "valid polyline with 2 vertices (minimum)",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 200, Y: 200},
			},
			opts: &PolylineOptions{
				Color: Black,
				Width: 1.0,
			},
			expectError: false,
		},
		{
			name: "valid polyline with dashed line",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 150, Y: 50},
				{X: 200, Y: 100},
				{X: 250, Y: 75},
			},
			opts: &PolylineOptions{
				Color:     Blue,
				Width:     1.5,
				Dashed:    true,
				DashArray: []float64{5, 3, 2, 3},
			},
			expectError: false,
		},
		{
			name: "valid long polyline",
			vertices: []Point{
				{X: 50, Y: 100},
				{X: 100, Y: 80},
				{X: 150, Y: 120},
				{X: 200, Y: 90},
				{X: 250, Y: 110},
				{X: 300, Y: 100},
			},
			opts: &PolylineOptions{
				Color: Green,
				Width: 2.5,
			},
			expectError: false,
		},
		{
			name:        "nil options",
			vertices:    []Point{{X: 100, Y: 100}, {X: 150, Y: 50}},
			opts:        nil,
			expectError: true,
			errorMsg:    "polyline options cannot be nil",
		},
		{
			name:        "too few vertices (1 vertex)",
			vertices:    []Point{{X: 100, Y: 100}},
			opts:        &PolylineOptions{Color: Black, Width: 1.0},
			expectError: true,
			errorMsg:    "polyline must have at least 2 vertices",
		},
		{
			name:        "empty vertices",
			vertices:    []Point{},
			opts:        &PolylineOptions{Color: Black, Width: 1.0},
			expectError: true,
			errorMsg:    "polyline must have at least 2 vertices",
		},
		{
			name: "invalid color (R > 1)",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 200, Y: 200},
			},
			opts: &PolylineOptions{
				Color: Color{R: 2.0, G: 0, B: 0},
				Width: 1.0,
			},
			expectError: true,
			errorMsg:    "color components must be in range [0.0, 1.0]",
		},
		{
			name: "invalid color (G < 0)",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 200, Y: 200},
			},
			opts: &PolylineOptions{
				Color: Color{R: 0, G: -0.5, B: 0},
				Width: 1.0,
			},
			expectError: true,
			errorMsg:    "color components must be in range [0.0, 1.0]",
		},
		{
			name: "negative line width",
			vertices: []Point{
				{X: 100, Y: 100},
				{X: 200, Y: 200},
			},
			opts: &PolylineOptions{
				Color: Black,
				Width: -1.0,
			},
			expectError: true,
			errorMsg:    "line width must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, err := c.NewPage()
			if err != nil {
				t.Fatalf("failed to create page: %v", err)
			}

			err = page.DrawPolyline(tt.vertices, tt.opts)

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
					if op.Type != GraphicsOpPolyline {
						t.Errorf("expected polyline operation, got type %d", op.Type)
					}
					if len(op.Vertices) != len(tt.vertices) {
						t.Errorf("expected %d vertices, got %d", len(tt.vertices), len(op.Vertices))
					}
				}
			}
		})
	}
}

func TestPolylineWavePattern(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Test wave pattern (sine-like curve approximation)
	waveVertices := []Point{
		{X: 50, Y: 100},
		{X: 75, Y: 80},
		{X: 100, Y: 100},
		{X: 125, Y: 120},
		{X: 150, Y: 100},
		{X: 175, Y: 80},
		{X: 200, Y: 100},
	}

	opts := &PolylineOptions{
		Color: Blue,
		Width: 2.0,
	}

	err = page.DrawPolyline(waveVertices, opts)
	if err != nil {
		t.Errorf("failed to draw wave polyline: %v", err)
	}

	// Verify 7 vertices
	ops := page.GraphicsOperations()
	if len(ops) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(ops))
	}
	if len(ops[0].Vertices) != 7 {
		t.Errorf("expected 7 vertices, got %d", len(ops[0].Vertices))
	}
}
