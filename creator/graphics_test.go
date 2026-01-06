package creator

import (
	"testing"
)

// TestDrawLine_Valid tests valid DrawLine cases.
func TestDrawLine_Valid(t *testing.T) {
	tests := []struct {
		name   string
		x1, y1 float64
		x2, y2 float64
		opts   *LineOptions
	}{
		{"basic line", 100, 700, 500, 700, &LineOptions{Color: Black, Width: 1.0}},
		{"colored line", 100, 600, 500, 600, &LineOptions{Color: Red, Width: 2.0}},
		{"dashed line", 100, 500, 500, 500, &LineOptions{Color: Blue, Width: 1.5, Dashed: true, DashArray: []float64{3, 1}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, _ := c.NewPage()
			if err := page.DrawLine(tt.x1, tt.y1, tt.x2, tt.y2, tt.opts); err != nil {
				t.Errorf("DrawLine() error = %v", err)
			}
			if len(page.graphicsOps) != 1 || page.graphicsOps[0].Type != GraphicsOpLine {
				t.Error("Expected 1 line operation")
			}
		})
	}
}

// TestDrawLine_Invalid tests DrawLine validation.
func TestDrawLine_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		x1, y1 float64
		x2, y2 float64
		opts   *LineOptions
	}{
		{"nil options", 100, 400, 500, 400, nil},
		{"invalid color", 100, 300, 500, 300, &LineOptions{Color: Color{R: 1.5, G: 0, B: 0}, Width: 1.0}},
		{"negative width", 100, 200, 500, 200, &LineOptions{Color: Black, Width: -1.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, _ := c.NewPage()
			if err := page.DrawLine(tt.x1, tt.y1, tt.x2, tt.y2, tt.opts); err == nil {
				t.Error("DrawLine() expected error")
			}
		})
	}
}

// TestDrawRect_Valid tests valid DrawRect cases.
func TestDrawRect_Valid(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		w, h float64
		opts *RectOptions
	}{
		{"stroke only", 100, 600, 200, 100, &RectOptions{StrokeColor: &Black, StrokeWidth: 1.0}},
		{"fill only", 100, 450, 200, 100, &RectOptions{FillColor: &LightGray}},
		{"stroke and fill", 100, 300, 200, 100, &RectOptions{StrokeColor: &Black, StrokeWidth: 2.0, FillColor: &Yellow}},
		{"dashed border", 100, 150, 200, 100, &RectOptions{StrokeColor: &Blue, StrokeWidth: 1.0, FillColor: &Cyan, Dashed: true, DashArray: []float64{5, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, _ := c.NewPage()
			if err := page.DrawRect(tt.x, tt.y, tt.w, tt.h, tt.opts); err != nil {
				t.Errorf("DrawRect() error = %v", err)
			}
			if len(page.graphicsOps) != 1 || page.graphicsOps[0].Type != GraphicsOpRect {
				t.Error("Expected 1 rect operation")
			}
		})
	}
}

// TestDrawRect_Invalid tests DrawRect validation.
func TestDrawRect_Invalid(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		w, h float64
		opts *RectOptions
	}{
		{"nil options", 100, 50, 200, 100, nil},
		{"no stroke or fill", 100, 50, 200, 100, &RectOptions{}},
		{"negative dimensions", 100, 50, -200, 100, &RectOptions{StrokeColor: &Black}},
		{"invalid stroke color", 100, 50, 200, 100, &RectOptions{StrokeColor: &Color{R: 2.0, G: 0, B: 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, _ := c.NewPage()
			if err := page.DrawRect(tt.x, tt.y, tt.w, tt.h, tt.opts); err == nil {
				t.Error("DrawRect() expected error")
			}
		})
	}
}

// TestDrawRectFilled tests the DrawRectFilled convenience method.
func TestDrawRectFilled(t *testing.T) {
	c := New()
	page, _ := c.NewPage()

	if err := page.DrawRectFilled(100, 500, 200, 150, Gray); err != nil {
		t.Errorf("DrawRectFilled() error = %v", err)
	}

	if len(page.graphicsOps) != 1 {
		t.Errorf("Expected 1 graphics operation, got %d", len(page.graphicsOps))
	}

	gop := page.graphicsOps[0]
	if gop.RectOpts.StrokeColor != nil {
		t.Error("Expected no stroke color")
	}
	if gop.RectOpts.FillColor == nil {
		t.Error("Expected fill color to be set")
	}
}

// TestDrawCircle_Valid tests valid DrawCircle cases.
func TestDrawCircle_Valid(t *testing.T) {
	tests := []struct {
		name   string
		cx, cy float64
		radius float64
		opts   *CircleOptions
	}{
		{"stroke only", 300, 400, 50, &CircleOptions{StrokeColor: &Red, StrokeWidth: 2.0}},
		{"fill only", 300, 250, 30, &CircleOptions{FillColor: &Blue}},
		{"stroke and fill", 300, 100, 40, &CircleOptions{StrokeColor: &Black, StrokeWidth: 1.5, FillColor: &Yellow}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, _ := c.NewPage()
			if err := page.DrawCircle(tt.cx, tt.cy, tt.radius, tt.opts); err != nil {
				t.Errorf("DrawCircle() error = %v", err)
			}
			if len(page.graphicsOps) != 1 || page.graphicsOps[0].Type != GraphicsOpCircle {
				t.Error("Expected 1 circle operation")
			}
		})
	}
}

// TestDrawCircle_Invalid tests DrawCircle validation.
func TestDrawCircle_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		cx, cy float64
		radius float64
		opts   *CircleOptions
	}{
		{"nil options", 300, 100, 40, nil},
		{"no stroke or fill", 300, 100, 40, &CircleOptions{}},
		{"negative radius", 300, 100, -40, &CircleOptions{StrokeColor: &Black}},
		{"invalid fill color", 300, 100, 40, &CircleOptions{FillColor: &Color{R: 0, G: -1, B: 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, _ := c.NewPage()
			if err := page.DrawCircle(tt.cx, tt.cy, tt.radius, tt.opts); err == nil {
				t.Error("DrawCircle() expected error")
			}
		})
	}
}

// TestGraphicsOperations tests the GraphicsOperations accessor.
func TestGraphicsOperations(t *testing.T) {
	c := New()
	page, _ := c.NewPage()

	if len(page.GraphicsOperations()) != 0 {
		t.Error("Expected no graphics operations initially")
	}

	_ = page.DrawLine(100, 700, 500, 700, &LineOptions{Color: Black, Width: 1})
	_ = page.DrawRect(100, 600, 200, 100, &RectOptions{StrokeColor: &Black, StrokeWidth: 1})
	_ = page.DrawCircle(300, 300, 50, &CircleOptions{FillColor: &Red})

	if len(page.GraphicsOperations()) != 3 {
		t.Errorf("Expected 3 graphics operations, got %d", len(page.GraphicsOperations()))
	}
}

// TestMixedTextAndGraphics tests using both text and graphics on same page.
func TestMixedTextAndGraphics(t *testing.T) {
	c := New()
	page, _ := c.NewPage()

	_ = page.DrawLine(100, 700, 500, 700, &LineOptions{Color: Black, Width: 2})
	_ = page.DrawRectFilled(100, 600, 200, 80, LightGray)
	_ = page.AddText("Graphics Test", 150, 630, HelveticaBold, 18)

	if len(page.GraphicsOperations()) != 2 {
		t.Errorf("Expected 2 graphics operations, got %d", len(page.GraphicsOperations()))
	}
	if len(page.TextOperations()) != 1 {
		t.Errorf("Expected 1 text operation, got %d", len(page.TextOperations()))
	}
}
