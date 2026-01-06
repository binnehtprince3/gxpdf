package creator

import (
	"testing"
)

func TestDrawBezierCurve(t *testing.T) {
	tests := []struct {
		name        string
		segments    []BezierSegment
		opts        *BezierOptions
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid single segment S-curve",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 150, Y: 200},
					C2:    Point{X: 200, Y: 200},
					End:   Point{X: 250, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color: Blue,
				Width: 2.0,
			},
			expectError: false,
		},
		{
			name: "valid multi-segment curve",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 125, Y: 150},
					C2:    Point{X: 150, Y: 150},
					End:   Point{X: 175, Y: 100},
				},
				{
					Start: Point{X: 175, Y: 100},
					C1:    Point{X: 200, Y: 50},
					C2:    Point{X: 225, Y: 50},
					End:   Point{X: 250, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color: Red,
				Width: 1.5,
			},
			expectError: false,
		},
		{
			name: "valid closed curve with fill",
			segments: []BezierSegment{
				{
					Start: Point{X: 150, Y: 50},
					C1:    Point{X: 200, Y: 100},
					C2:    Point{X: 200, Y: 150},
					End:   Point{X: 150, Y: 200},
				},
				{
					Start: Point{X: 150, Y: 200},
					C1:    Point{X: 100, Y: 150},
					C2:    Point{X: 100, Y: 100},
					End:   Point{X: 150, Y: 50},
				},
			},
			opts: &BezierOptions{
				Color:     Black,
				Width:     1.0,
				Closed:    true,
				FillColor: &Yellow,
			},
			expectError: false,
		},
		{
			name: "valid curve with dashed line",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 150, Y: 200},
					C2:    Point{X: 200, Y: 200},
					End:   Point{X: 250, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color:     Green,
				Width:     2.0,
				Dashed:    true,
				DashArray: []float64{10, 5},
			},
			expectError: false,
		},
		{
			name:        "nil options",
			segments:    []BezierSegment{{Start: Point{X: 100, Y: 100}, C1: Point{X: 150, Y: 200}, C2: Point{X: 200, Y: 200}, End: Point{X: 250, Y: 100}}},
			opts:        nil,
			expectError: true,
			errorMsg:    "bezier curve options cannot be nil",
		},
		{
			name:        "empty segments",
			segments:    []BezierSegment{},
			opts:        &BezierOptions{Color: Black, Width: 1.0},
			expectError: true,
			errorMsg:    "bezier curve must have at least 1 segment",
		},
		{
			name: "discontinuous segments",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 125, Y: 150},
					C2:    Point{X: 150, Y: 150},
					End:   Point{X: 175, Y: 100},
				},
				{
					Start: Point{X: 200, Y: 100}, // Does not match previous End point
					C1:    Point{X: 225, Y: 50},
					C2:    Point{X: 250, Y: 50},
					End:   Point{X: 275, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color: Red,
				Width: 1.0,
			},
			expectError: true,
			errorMsg:    "bezier segments must be continuous (segment start point must match previous segment end point)",
		},
		{
			name: "invalid color (R > 1)",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 150, Y: 200},
					C2:    Point{X: 200, Y: 200},
					End:   Point{X: 250, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color: Color{R: 1.5, G: 0, B: 0},
				Width: 1.0,
			},
			expectError: true,
			errorMsg:    "color components must be in range [0.0, 1.0]",
		},
		{
			name: "negative width",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 150, Y: 200},
					C2:    Point{X: 200, Y: 200},
					End:   Point{X: 250, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color: Black,
				Width: -1.0,
			},
			expectError: true,
			errorMsg:    "curve width must be non-negative",
		},
		{
			name: "fill color without closed",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 150, Y: 200},
					C2:    Point{X: 200, Y: 200},
					End:   Point{X: 250, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color:     Black,
				Width:     1.0,
				FillColor: &Yellow,
			},
			expectError: true,
			errorMsg:    "fill color requires closed curve (set Closed: true)",
		},
		{
			name: "invalid fill color",
			segments: []BezierSegment{
				{
					Start: Point{X: 100, Y: 100},
					C1:    Point{X: 150, Y: 200},
					C2:    Point{X: 200, Y: 200},
					End:   Point{X: 250, Y: 100},
				},
			},
			opts: &BezierOptions{
				Color:     Black,
				Width:     1.0,
				Closed:    true,
				FillColor: &Color{R: 0, G: 0, B: -0.5},
			},
			expectError: true,
			errorMsg:    "fill color components must be in range [0.0, 1.0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, err := c.NewPage()
			if err != nil {
				t.Fatalf("failed to create page: %v", err)
			}

			err = page.DrawBezierCurve(tt.segments, tt.opts)

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
					if op.Type != GraphicsOpBezier {
						t.Errorf("expected bezier operation, got type %d", op.Type)
					}
					if len(op.BezierSegs) != len(tt.segments) {
						t.Errorf("expected %d segments, got %d", len(tt.segments), len(op.BezierSegs))
					}
				}
			}
		})
	}
}

func TestBezierComplexCurves(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Test smooth wave pattern (4 segments)
	waveSegments := []BezierSegment{
		{
			Start: Point{X: 50, Y: 100},
			C1:    Point{X: 75, Y: 70},
			C2:    Point{X: 100, Y: 70},
			End:   Point{X: 125, Y: 100},
		},
		{
			Start: Point{X: 125, Y: 100},
			C1:    Point{X: 150, Y: 130},
			C2:    Point{X: 175, Y: 130},
			End:   Point{X: 200, Y: 100},
		},
		{
			Start: Point{X: 200, Y: 100},
			C1:    Point{X: 225, Y: 70},
			C2:    Point{X: 250, Y: 70},
			End:   Point{X: 275, Y: 100},
		},
		{
			Start: Point{X: 275, Y: 100},
			C1:    Point{X: 300, Y: 130},
			C2:    Point{X: 325, Y: 130},
			End:   Point{X: 350, Y: 100},
		},
	}

	opts := &BezierOptions{
		Color: Blue,
		Width: 2.0,
	}

	err = page.DrawBezierCurve(waveSegments, opts)
	if err != nil {
		t.Errorf("failed to draw wave bezier curve: %v", err)
	}

	// Verify 4 segments
	ops := page.GraphicsOperations()
	if len(ops) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(ops))
	}
	if len(ops[0].BezierSegs) != 4 {
		t.Errorf("expected 4 segments, got %d", len(ops[0].BezierSegs))
	}
}

func TestBezierSegmentContinuityValidation(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	// Test: Segments with small floating-point differences should be accepted
	segments := []BezierSegment{
		{
			Start: Point{X: 100, Y: 100},
			C1:    Point{X: 125, Y: 150},
			C2:    Point{X: 150, Y: 150},
			End:   Point{X: 175.0000001, Y: 99.9999999}, // Tiny difference
		},
		{
			Start: Point{X: 175, Y: 100}, // Matches End within epsilon
			C1:    Point{X: 200, Y: 50},
			C2:    Point{X: 225, Y: 50},
			End:   Point{X: 250, Y: 100},
		},
	}

	opts := &BezierOptions{
		Color: Red,
		Width: 1.0,
	}

	err = page.DrawBezierCurve(segments, opts)
	if err != nil {
		t.Errorf("expected continuous segments to be accepted, got error: %v", err)
	}
}
