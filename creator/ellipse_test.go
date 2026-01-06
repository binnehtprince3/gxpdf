package creator

import (
	"testing"
)

func TestDrawEllipse(t *testing.T) {
	tests := []struct {
		name        string
		cx, cy      float64
		rx, ry      float64
		opts        *EllipseOptions
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid horizontal ellipse with stroke and fill",
			cx:   150, cy: 200,
			rx: 100, ry: 50,
			opts: &EllipseOptions{
				StrokeColor: &Black,
				StrokeWidth: 2.0,
				FillColor:   &Green,
			},
			expectError: false,
		},
		{
			name: "valid vertical ellipse",
			cx:   150, cy: 200,
			rx: 50, ry: 100,
			opts: &EllipseOptions{
				StrokeColor: &Red,
				StrokeWidth: 1.5,
				FillColor:   &Yellow,
			},
			expectError: false,
		},
		{
			name: "valid circle (rx = ry)",
			cx:   150, cy: 200,
			rx: 75, ry: 75,
			opts: &EllipseOptions{
				StrokeColor: &Blue,
				StrokeWidth: 1.0,
			},
			expectError: false,
		},
		{
			name: "valid ellipse with fill only",
			cx:   150, cy: 200,
			rx: 100, ry: 60,
			opts: &EllipseOptions{
				FillColor: &LightGray,
			},
			expectError: false,
		},
		{
			name: "valid ellipse with stroke only",
			cx:   150, cy: 200,
			rx: 80, ry: 40,
			opts: &EllipseOptions{
				StrokeColor: &Black,
				StrokeWidth: 2.5,
			},
			expectError: false,
		},
		{
			name: "nil options",
			cx:   150, cy: 200,
			rx: 100, ry: 50,
			opts:        nil,
			expectError: true,
			errorMsg:    "ellipse options cannot be nil",
		},
		{
			name: "negative horizontal radius",
			cx:   150, cy: 200,
			rx: -100, ry: 50,
			opts: &EllipseOptions{
				StrokeColor: &Black,
			},
			expectError: true,
			errorMsg:    "horizontal radius must be non-negative",
		},
		{
			name: "negative vertical radius",
			cx:   150, cy: 200,
			rx: 100, ry: -50,
			opts: &EllipseOptions{
				StrokeColor: &Black,
			},
			expectError: true,
			errorMsg:    "vertical radius must be non-negative",
		},
		{
			name: "neither stroke nor fill",
			cx:   150, cy: 200,
			rx: 100, ry: 50,
			opts:        &EllipseOptions{},
			expectError: true,
			errorMsg:    "ellipse must have at least stroke, fill color, or gradient",
		},
		{
			name: "invalid stroke color (R > 1)",
			cx:   150, cy: 200,
			rx: 100, ry: 50,
			opts: &EllipseOptions{
				StrokeColor: &Color{R: 1.5, G: 0, B: 0},
			},
			expectError: true,
			errorMsg:    "stroke color components must be in range [0.0, 1.0]",
		},
		{
			name: "invalid fill color (B < 0)",
			cx:   150, cy: 200,
			rx: 100, ry: 50,
			opts: &EllipseOptions{
				FillColor: &Color{R: 0, G: 0, B: -0.1},
			},
			expectError: true,
			errorMsg:    "fill color components must be in range [0.0, 1.0]",
		},
		{
			name: "negative stroke width",
			cx:   150, cy: 200,
			rx: 100, ry: 50,
			opts: &EllipseOptions{
				StrokeColor: &Black,
				StrokeWidth: -1.0,
			},
			expectError: true,
			errorMsg:    "stroke width must be non-negative",
		},
		{
			name: "zero radii (valid, creates a point)",
			cx:   150, cy: 200,
			rx: 0, ry: 0,
			opts: &EllipseOptions{
				StrokeColor: &Black,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, err := c.NewPage()
			if err != nil {
				t.Fatalf("failed to create page: %v", err)
			}

			err = page.DrawEllipse(tt.cx, tt.cy, tt.rx, tt.ry, tt.opts)

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
					if op.Type != GraphicsOpEllipse {
						t.Errorf("expected ellipse operation, got type %d", op.Type)
					}
					if op.X != tt.cx || op.Y != tt.cy {
						t.Errorf("expected center (%f,%f), got (%f,%f)", tt.cx, tt.cy, op.X, op.Y)
					}
					if op.RX != tt.rx || op.RY != tt.ry {
						t.Errorf("expected radii (%f,%f), got (%f,%f)", tt.rx, tt.ry, op.RX, op.RY)
					}
				}
			}
		})
	}
}

func TestEllipseAspectRatios(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Test various aspect ratios
	ratios := []struct {
		name string
		rx   float64
		ry   float64
	}{
		{"very wide", 200, 25},
		{"wide", 150, 50},
		{"slightly wide", 110, 90},
		{"circle", 100, 100},
		{"slightly tall", 90, 110},
		{"tall", 50, 150},
		{"very tall", 25, 200},
	}

	opts := &EllipseOptions{
		StrokeColor: &Black,
		StrokeWidth: 1.0,
	}

	for i, ratio := range ratios {
		err = page.DrawEllipse(150, 200+float64(i*50), ratio.rx, ratio.ry, opts)
		if err != nil {
			t.Errorf("failed to draw %s ellipse: %v", ratio.name, err)
		}
	}

	// Verify all ellipses were added
	ops := page.GraphicsOperations()
	if len(ops) != len(ratios) {
		t.Errorf("expected %d ellipses, got %d", len(ratios), len(ops))
	}
}
