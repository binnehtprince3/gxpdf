package creator

import (
	"testing"
)

func TestNewTextWatermark(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "simple text",
			text: "CONFIDENTIAL",
			want: "CONFIDENTIAL",
		},
		{
			name: "draft watermark",
			text: "DRAFT",
			want: "DRAFT",
		},
		{
			name: "empty text",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewTextWatermark(tt.text)

			if wm.Text() != tt.want {
				t.Errorf("Text() = %v, want %v", wm.Text(), tt.want)
			}

			// Verify defaults.
			if wm.Font() != HelveticaBold {
				t.Errorf("default font = %v, want %v", wm.Font(), HelveticaBold)
			}
			if wm.FontSize() != 48 {
				t.Errorf("default font size = %v, want 48", wm.FontSize())
			}
			if wm.Opacity() != 0.5 {
				t.Errorf("default opacity = %v, want 0.5", wm.Opacity())
			}
			if wm.Rotation() != 45 {
				t.Errorf("default rotation = %v, want 45", wm.Rotation())
			}
			if wm.Position() != WatermarkCenter {
				t.Errorf("default position = %v, want %v", wm.Position(), WatermarkCenter)
			}
		})
	}
}

func TestTextWatermark_SetFont(t *testing.T) {
	tests := []struct {
		name    string
		font    FontName
		size    float64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid helvetica bold 72pt",
			font:    HelveticaBold,
			size:    72,
			wantErr: false,
		},
		{
			name:    "valid times roman 12pt",
			font:    TimesRoman,
			size:    12,
			wantErr: false,
		},
		{
			name:    "zero font size",
			font:    Helvetica,
			size:    0,
			wantErr: true,
			errMsg:  "font size must be positive",
		},
		{
			name:    "negative font size",
			font:    Helvetica,
			size:    -10,
			wantErr: true,
			errMsg:  "font size must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewTextWatermark("TEST")
			err := wm.SetFont(tt.font, tt.size)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetFont() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("SetFont() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("SetFont() unexpected error: %v", err)
				}
				if wm.Font() != tt.font {
					t.Errorf("Font() = %v, want %v", wm.Font(), tt.font)
				}
				if wm.FontSize() != tt.size {
					t.Errorf("FontSize() = %v, want %v", wm.FontSize(), tt.size)
				}
			}
		})
	}
}

func TestTextWatermark_SetColor(t *testing.T) {
	tests := []struct {
		name    string
		color   Color
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid gray",
			color:   Gray,
			wantErr: false,
		},
		{
			name:    "valid red",
			color:   Red,
			wantErr: false,
		},
		{
			name:    "valid black",
			color:   Black,
			wantErr: false,
		},
		{
			name:    "valid white",
			color:   White,
			wantErr: false,
		},
		{
			name:    "invalid red > 1",
			color:   Color{R: 1.5, G: 0, B: 0},
			wantErr: true,
			errMsg:  "color components must be in range [0.0, 1.0]",
		},
		{
			name:    "invalid green < 0",
			color:   Color{R: 0.5, G: -0.1, B: 0.5},
			wantErr: true,
			errMsg:  "color components must be in range [0.0, 1.0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewTextWatermark("TEST")
			err := wm.SetColor(tt.color)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetColor() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("SetColor() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("SetColor() unexpected error: %v", err)
				}
				if wm.Color() != tt.color {
					t.Errorf("Color() = %v, want %v", wm.Color(), tt.color)
				}
			}
		})
	}
}

func TestTextWatermark_SetOpacity(t *testing.T) {
	tests := []struct {
		name    string
		opacity float64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid 0.3",
			opacity: 0.3,
			wantErr: false,
		},
		{
			name:    "valid 1.0 (fully opaque)",
			opacity: 1.0,
			wantErr: false,
		},
		{
			name:    "valid 0.0 (fully transparent)",
			opacity: 0.0,
			wantErr: false,
		},
		{
			name:    "invalid > 1",
			opacity: 1.5,
			wantErr: true,
			errMsg:  "opacity must be in range [0.0, 1.0]",
		},
		{
			name:    "invalid < 0",
			opacity: -0.1,
			wantErr: true,
			errMsg:  "opacity must be in range [0.0, 1.0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewTextWatermark("TEST")
			err := wm.SetOpacity(tt.opacity)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SetOpacity() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("SetOpacity() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("SetOpacity() unexpected error: %v", err)
				}
				if wm.Opacity() != tt.opacity {
					t.Errorf("Opacity() = %v, want %v", wm.Opacity(), tt.opacity)
				}
			}
		})
	}
}

func TestTextWatermark_SetRotation(t *testing.T) {
	tests := []struct {
		name     string
		rotation float64
	}{
		{
			name:     "horizontal",
			rotation: 0,
		},
		{
			name:     "diagonal 45",
			rotation: 45,
		},
		{
			name:     "vertical",
			rotation: 90,
		},
		{
			name:     "diagonal 135",
			rotation: 135,
		},
		{
			name:     "upside down",
			rotation: 180,
		},
		{
			name:     "negative rotation",
			rotation: -45,
		},
		{
			name:     "full rotation",
			rotation: 360,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewTextWatermark("TEST")
			err := wm.SetRotation(tt.rotation)

			if err != nil {
				t.Errorf("SetRotation() unexpected error: %v", err)
			}
			if wm.Rotation() != tt.rotation {
				t.Errorf("Rotation() = %v, want %v", wm.Rotation(), tt.rotation)
			}
		})
	}
}

func TestTextWatermark_SetPosition(t *testing.T) {
	tests := []struct {
		name     string
		position WatermarkPosition
	}{
		{
			name:     "center",
			position: WatermarkCenter,
		},
		{
			name:     "top left",
			position: WatermarkTopLeft,
		},
		{
			name:     "top right",
			position: WatermarkTopRight,
		},
		{
			name:     "bottom left",
			position: WatermarkBottomLeft,
		},
		{
			name:     "bottom right",
			position: WatermarkBottomRight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewTextWatermark("TEST")
			err := wm.SetPosition(tt.position)

			if err != nil {
				t.Errorf("SetPosition() unexpected error: %v", err)
			}
			if wm.Position() != tt.position {
				t.Errorf("Position() = %v, want %v", wm.Position(), tt.position)
			}
		})
	}
}

func TestPage_DrawWatermark(t *testing.T) {
	tests := []struct {
		name      string
		watermark *TextWatermark
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid watermark",
			watermark: NewTextWatermark("CONFIDENTIAL"),
			wantErr:   false,
		},
		{
			name: "custom watermark",
			watermark: func() *TextWatermark {
				wm := NewTextWatermark("DRAFT")
				_ = wm.SetFont(HelveticaBold, 72)
				_ = wm.SetColor(Red)
				_ = wm.SetOpacity(0.3)
				_ = wm.SetRotation(45)
				return wm
			}(),
			wantErr: false,
		},
		{
			name:      "nil watermark",
			watermark: nil,
			wantErr:   true,
			errMsg:    "watermark cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			page, _ := c.NewPage()
			err := page.DrawWatermark(tt.watermark)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DrawWatermark() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("DrawWatermark() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DrawWatermark() unexpected error: %v", err)
				}

				// Verify watermark was added to graphics operations.
				ops := page.GraphicsOperations()
				if len(ops) != 1 {
					t.Errorf("expected 1 graphics operation, got %d", len(ops))
					return
				}

				op := ops[0]
				if op.Type != GraphicsOpWatermark {
					t.Errorf("operation type = %v, want %v", op.Type, GraphicsOpWatermark)
				}
				if op.WatermarkOp != tt.watermark {
					t.Errorf("operation watermark = %v, want %v", op.WatermarkOp, tt.watermark)
				}
			}
		})
	}
}

func TestCalculateWatermarkPosition(t *testing.T) {
	c := New()
	c.SetPageSize(A4)
	page, _ := c.NewPage() // A4 = 595 × 842 points.

	tests := []struct {
		name     string
		position WatermarkPosition
		fontSize float64
		checkX   bool
		checkY   bool
		wantX    float64
		wantY    float64
	}{
		{
			name:     "center",
			position: WatermarkCenter,
			fontSize: 48,
			checkX:   true,
			checkY:   true,
			wantX:    297.5, // 595 / 2
			wantY:    421,   // 842 / 2
		},
		{
			name:     "top left",
			position: WatermarkTopLeft,
			fontSize: 48,
			checkX:   true,
			checkY:   true,
			wantX:    24,  // fontSize * 0.5
			wantY:    818, // 842 - fontSize * 0.5
		},
		{
			name:     "bottom right",
			position: WatermarkBottomRight,
			fontSize: 48,
			checkX:   false, // Skip X check (text width dependent)
			checkY:   true,
			wantX:    0,  // Not checked
			wantY:    72, // padding + fontSize
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewTextWatermark("TEST")
			_ = wm.SetFont(HelveticaBold, tt.fontSize)
			_ = wm.SetPosition(tt.position)

			x, y := calculateWatermarkPosition(page, wm)

			// Allow some tolerance for text width calculations.
			tolerance := 1.0
			if tt.checkX && !floatNear(x, tt.wantX, tolerance) {
				t.Errorf("x = %v, want %v (tolerance %v)", x, tt.wantX, tolerance)
			}
			if tt.checkY && !floatNear(y, tt.wantY, tolerance) {
				t.Errorf("y = %v, want %v (tolerance %v)", y, tt.wantY, tolerance)
			}
		})
	}
}

func TestRotationMatrix(t *testing.T) {
	tests := []struct {
		name    string
		x       float64
		y       float64
		degrees float64
		wantA   float64 // cos
		wantB   float64 // sin
		wantC   float64 // -sin
		wantD   float64 // cos
	}{
		{
			name:    "no rotation",
			x:       100,
			y:       200,
			degrees: 0,
			wantA:   1, // cos(0) = 1
			wantB:   0, // sin(0) = 0
			wantC:   0, // -sin(0) = 0
			wantD:   1, // cos(0) = 1
		},
		{
			name:    "45 degrees",
			x:       100,
			y:       200,
			degrees: 45,
			wantA:   0.7071,  // cos(45°) ≈ 0.7071
			wantB:   0.7071,  // sin(45°) ≈ 0.7071
			wantC:   -0.7071, // -sin(45°) ≈ -0.7071
			wantD:   0.7071,  // cos(45°) ≈ 0.7071
		},
		{
			name:    "90 degrees",
			x:       100,
			y:       200,
			degrees: 90,
			wantA:   0,  // cos(90°) ≈ 0
			wantB:   1,  // sin(90°) = 1
			wantC:   -1, // -sin(90°) = -1
			wantD:   0,  // cos(90°) ≈ 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrix := rotationMatrix(tt.x, tt.y, tt.degrees)

			tolerance := 0.0001
			if !floatNear(matrix[0], tt.wantA, tolerance) {
				t.Errorf("matrix[0] (a/cos) = %v, want %v", matrix[0], tt.wantA)
			}
			if !floatNear(matrix[1], tt.wantB, tolerance) {
				t.Errorf("matrix[1] (b/sin) = %v, want %v", matrix[1], tt.wantB)
			}
			if !floatNear(matrix[2], tt.wantC, tolerance) {
				t.Errorf("matrix[2] (c/-sin) = %v, want %v", matrix[2], tt.wantC)
			}
			if !floatNear(matrix[3], tt.wantD, tolerance) {
				t.Errorf("matrix[3] (d/cos) = %v, want %v", matrix[3], tt.wantD)
			}
		})
	}
}

// floatNear checks if two floats are within tolerance of each other.
func floatNear(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
