package creator

import (
	"testing"
)

func TestNewLinearGradient(t *testing.T) {
	grad := NewLinearGradient(0, 0, 100, 0)

	if grad == nil {
		t.Fatal("NewLinearGradient() returned nil")
	}
	if grad.Type != GradientTypeLinear {
		t.Errorf("Type = %d, want %d", grad.Type, GradientTypeLinear)
	}
	if grad.X1 != 0 || grad.Y1 != 0 || grad.X2 != 100 || grad.Y2 != 0 {
		t.Errorf("Coordinates = (%f,%f) to (%f,%f), want (0,0) to (100,0)",
			grad.X1, grad.Y1, grad.X2, grad.Y2)
	}
	if !grad.ExtendStart || !grad.ExtendEnd {
		t.Error("ExtendStart and ExtendEnd should default to true")
	}
	if len(grad.ColorStops) != 0 {
		t.Errorf("ColorStops length = %d, want 0", len(grad.ColorStops))
	}
}

func TestNewRadialGradient(t *testing.T) {
	grad := NewRadialGradient(150, 550, 0, 150, 550, 50)

	if grad == nil {
		t.Fatal("NewRadialGradient() returned nil")
	}
	if grad.Type != GradientTypeRadial {
		t.Errorf("Type = %d, want %d", grad.Type, GradientTypeRadial)
	}
	if grad.X0 != 150 || grad.Y0 != 550 || grad.R0 != 0 {
		t.Errorf("Start circle = (%f,%f,%f), want (150,550,0)", grad.X0, grad.Y0, grad.R0)
	}
	if grad.X1 != 150 || grad.Y1 != 550 || grad.R1 != 50 {
		t.Errorf("End circle = (%f,%f,%f), want (150,550,50)", grad.X1, grad.Y1, grad.R1)
	}
}

func TestGradient_AddColorStop(t *testing.T) {
	grad := NewLinearGradient(0, 0, 100, 0)

	// Add valid color stops
	err := grad.AddColorStop(0, Red)
	if err != nil {
		t.Errorf("AddColorStop(0, Red) failed: %v", err)
	}

	err = grad.AddColorStop(0.5, Yellow)
	if err != nil {
		t.Errorf("AddColorStop(0.5, Yellow) failed: %v", err)
	}

	err = grad.AddColorStop(1, Green)
	if err != nil {
		t.Errorf("AddColorStop(1, Green) failed: %v", err)
	}

	if len(grad.ColorStops) != 3 {
		t.Errorf("ColorStops length = %d, want 3", len(grad.ColorStops))
	}

	// Verify sorting (should be 0, 0.5, 1)
	if grad.ColorStops[0].Position != 0 {
		t.Errorf("ColorStops[0].Position = %f, want 0", grad.ColorStops[0].Position)
	}
	if grad.ColorStops[1].Position != 0.5 {
		t.Errorf("ColorStops[1].Position = %f, want 0.5", grad.ColorStops[1].Position)
	}
	if grad.ColorStops[2].Position != 1 {
		t.Errorf("ColorStops[2].Position = %f, want 1", grad.ColorStops[2].Position)
	}
}

func TestGradient_AddColorStop_Sorting(t *testing.T) {
	grad := NewLinearGradient(0, 0, 100, 0)

	// Add in wrong order
	grad.AddColorStop(1, Green)
	grad.AddColorStop(0, Red)
	grad.AddColorStop(0.5, Yellow)

	// Should be sorted: 0, 0.5, 1
	if grad.ColorStops[0].Position != 0 {
		t.Errorf("After sorting: ColorStops[0].Position = %f, want 0", grad.ColorStops[0].Position)
	}
	if grad.ColorStops[1].Position != 0.5 {
		t.Errorf("After sorting: ColorStops[1].Position = %f, want 0.5", grad.ColorStops[1].Position)
	}
	if grad.ColorStops[2].Position != 1 {
		t.Errorf("After sorting: ColorStops[2].Position = %f, want 1", grad.ColorStops[2].Position)
	}
}

func TestGradient_AddColorStop_InvalidPosition(t *testing.T) {
	grad := NewLinearGradient(0, 0, 100, 0)

	tests := []struct {
		name     string
		position float64
		wantErr  bool
	}{
		{"Negative position", -0.1, true},
		{"Position > 1", 1.1, true},
		{"Position = 0", 0.0, false},
		{"Position = 1", 1.0, false},
		{"Position = 0.5", 0.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := grad.AddColorStop(tt.position, Red)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddColorStop(%f) error = %v, wantErr = %v",
					tt.position, err, tt.wantErr)
			}
		})
	}
}

func TestGradient_Validate_Linear(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Gradient
		wantErr bool
	}{
		{
			name: "Valid linear gradient",
			setup: func() *Gradient {
				g := NewLinearGradient(0, 0, 100, 0)
				g.AddColorStop(0, Red)
				g.AddColorStop(1, Blue)
				return g
			},
			wantErr: false,
		},
		{
			name: "Missing color stops",
			setup: func() *Gradient {
				g := NewLinearGradient(0, 0, 100, 0)
				g.AddColorStop(0, Red)
				return g
			},
			wantErr: true,
		},
		{
			name: "Same start and end points",
			setup: func() *Gradient {
				g := NewLinearGradient(50, 50, 50, 50)
				g.AddColorStop(0, Red)
				g.AddColorStop(1, Blue)
				return g
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grad := tt.setup()
			err := grad.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestGradient_Validate_Radial(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Gradient
		wantErr bool
	}{
		{
			name: "Valid radial gradient",
			setup: func() *Gradient {
				g := NewRadialGradient(150, 550, 0, 150, 550, 50)
				g.AddColorStop(0, White)
				g.AddColorStop(1, Blue)
				return g
			},
			wantErr: false,
		},
		{
			name: "Negative inner radius",
			setup: func() *Gradient {
				g := NewRadialGradient(150, 550, -1, 150, 550, 50)
				g.AddColorStop(0, White)
				g.AddColorStop(1, Blue)
				return g
			},
			wantErr: true,
		},
		{
			name: "Negative outer radius",
			setup: func() *Gradient {
				g := NewRadialGradient(150, 550, 0, 150, 550, -10)
				g.AddColorStop(0, White)
				g.AddColorStop(1, Blue)
				return g
			},
			wantErr: true,
		},
		{
			name: "Both radii zero",
			setup: func() *Gradient {
				g := NewRadialGradient(150, 550, 0, 150, 550, 0)
				g.AddColorStop(0, White)
				g.AddColorStop(1, Blue)
				return g
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grad := tt.setup()
			err := grad.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestRectOptions_Gradient(t *testing.T) {
	grad := NewLinearGradient(0, 0, 100, 0)
	grad.AddColorStop(0, Red)
	grad.AddColorStop(1, Blue)

	opts := &RectOptions{
		FillGradient: grad,
		StrokeColor:  &Black,
		StrokeWidth:  2.0,
	}

	// Test validation - should succeed with gradient
	if opts.FillGradient == nil {
		t.Error("FillGradient should be set")
	}

	// Test mutual exclusivity
	opts.FillColor = &Red
	if opts.FillColor != nil && opts.FillGradient != nil {
		// This should be caught by validation
		err := validateRectOptions(opts)
		if err == nil {
			t.Error("validateRectOptions should reject both fill color and gradient")
		}
	}
}

func TestCircleOptions_Gradient(t *testing.T) {
	grad := NewRadialGradient(0, 0, 0, 0, 0, 50)
	grad.AddColorStop(0, White)
	grad.AddColorStop(1, Blue)

	opts := &CircleOptions{
		FillGradient: grad,
		StrokeColor:  &Black,
		StrokeWidth:  1.0,
	}

	if opts.FillGradient == nil {
		t.Error("FillGradient should be set")
	}
}

func TestPolygonOptions_Gradient(t *testing.T) {
	grad := NewLinearGradient(100, 100, 200, 200)
	grad.AddColorStop(0, Yellow)
	grad.AddColorStop(1, Red)

	opts := &PolygonOptions{
		FillGradient: grad,
		StrokeColor:  &Black,
		StrokeWidth:  1.5,
	}

	if opts.FillGradient == nil {
		t.Error("FillGradient should be set")
	}
}
