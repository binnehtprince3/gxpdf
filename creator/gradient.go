package creator

import (
	"errors"
	"fmt"
)

// GradientType represents the type of gradient.
type GradientType int

const (
	// GradientTypeLinear represents an axial (linear) gradient.
	// PDF ShadingType 2: gradient along a straight line from start to end point.
	GradientTypeLinear GradientType = 2

	// GradientTypeRadial represents a radial gradient.
	// PDF ShadingType 3: gradient radiating from center point.
	GradientTypeRadial GradientType = 3
)

// ColorStop represents a color transition point in a gradient.
//
// A gradient is defined by multiple color stops along a [0, 1] domain.
// Position 0 is the start, position 1 is the end.
//
// Example:
//
//	stops := []ColorStop{
//	    {Position: 0.0, Color: Red},    // Start: red
//	    {Position: 0.5, Color: Yellow}, // Middle: yellow
//	    {Position: 1.0, Color: Green},  // End: green
//	}
type ColorStop struct {
	// Position is the location of this color stop in the gradient (0.0 to 1.0).
	// 0.0 = start, 1.0 = end.
	Position float64

	// Color is the RGB color at this position.
	Color Color
}

// Gradient represents a color gradient (linear or radial).
//
// Gradients can be used to fill shapes with smooth color transitions.
// Two types are supported:
//   - Linear (axial): Color transitions along a straight line
//   - Radial: Color radiates from a center point
//
// Gradients are defined by color stops at specific positions.
// The PDF renderer interpolates colors between stops.
type Gradient struct {
	// Type is the gradient type (linear or radial).
	Type GradientType

	// ColorStops define the color transitions (minimum 2 required).
	// Stops must be sorted by Position (0.0 to 1.0).
	ColorStops []ColorStop

	// Linear gradient fields (Type == GradientTypeLinear)
	// Coordinates define the gradient axis from (X1, Y1) to (X2, Y2).
	X1, Y1 float64 // Start point
	X2, Y2 float64 // End point

	// Radial gradient fields (Type == GradientTypeRadial)
	// Coordinates define two circles: (X0, Y0, R0) and (X1, Y1, R1).
	// Gradient transitions from inner circle to outer circle.
	X0, Y0, R0 float64 // Starting circle (center + radius)
	// X1, Y1 already defined above (reused for radial end center)
	R1 float64 // Ending radius

	// Extend flags control what happens outside the gradient domain.
	// If true, colors extend beyond the gradient boundaries.
	ExtendStart bool // Extend before the first color stop
	ExtendEnd   bool // Extend after the last color stop
}

// NewLinearGradient creates a new linear (axial) gradient.
//
// The gradient transitions along a line from (x1, y1) to (x2, y2).
// Color stops must be added separately using AddColorStop().
//
// Parameters:
//   - x1, y1: Start point of gradient axis
//   - x2, y2: End point of gradient axis
//
// Example:
//
//	grad := creator.NewLinearGradient(0, 0, 100, 0) // Horizontal gradient
//	grad.AddColorStop(0, creator.Red)
//	grad.AddColorStop(1, creator.Blue)
//
// PDF Reference: ShadingType 2 (Axial Shading).
func NewLinearGradient(x1, y1, x2, y2 float64) *Gradient {
	return &Gradient{
		Type:        GradientTypeLinear,
		X1:          x1,
		Y1:          y1,
		X2:          x2,
		Y2:          y2,
		ExtendStart: true, // Default: extend colors beyond gradient
		ExtendEnd:   true,
		ColorStops:  make([]ColorStop, 0),
	}
}

// NewRadialGradient creates a new radial gradient.
//
// The gradient radiates from an inner circle (x0, y0, r0) to an outer circle (x1, y1, r1).
// For a simple radial gradient from center outward, use same center for both circles
// and set r0 = 0.
//
// Parameters:
//   - x0, y0, r0: Center and radius of starting circle
//   - x1, y1, r1: Center and radius of ending circle
//
// Example:
//
//	// Simple radial gradient from center
//	grad := creator.NewRadialGradient(150, 550, 0, 150, 550, 50)
//	grad.AddColorStop(0, creator.White) // Center: white
//	grad.AddColorStop(1, creator.Blue)  // Edge: blue
//
// PDF Reference: ShadingType 3 (Radial Shading).
func NewRadialGradient(x0, y0, r0, x1, y1, r1 float64) *Gradient {
	return &Gradient{
		Type:        GradientTypeRadial,
		X0:          x0,
		Y0:          y0,
		R0:          r0,
		X1:          x1,
		Y1:          y1,
		R1:          r1,
		ExtendStart: true, // Default: extend colors beyond gradient
		ExtendEnd:   true,
		ColorStops:  make([]ColorStop, 0),
	}
}

// AddColorStop adds a color stop to the gradient.
//
// Color stops define the color at specific positions along the gradient.
// Position must be in range [0.0, 1.0] where 0 is start and 1 is end.
//
// Stops can be added in any order; they will be sorted automatically.
//
// Parameters:
//   - position: Position in gradient (0.0 = start, 1.0 = end)
//   - color: RGB color at this position
//
// Example:
//
//	grad := creator.NewLinearGradient(0, 0, 100, 0)
//	grad.AddColorStop(0.0, creator.Red)
//	grad.AddColorStop(0.5, creator.Yellow)
//	grad.AddColorStop(1.0, creator.Green)
func (g *Gradient) AddColorStop(position float64, color Color) error {
	if position < 0.0 || position > 1.0 {
		return fmt.Errorf("color stop position must be in range [0, 1], got: %f", position)
	}

	if err := validateColor(color); err != nil {
		return fmt.Errorf("color stop color: %w", err)
	}

	g.ColorStops = append(g.ColorStops, ColorStop{
		Position: position,
		Color:    color,
	})

	// Sort color stops by position (required by PDF spec)
	g.sortColorStops()

	return nil
}

// sortColorStops sorts color stops by position (ascending order).
func (g *Gradient) sortColorStops() {
	// Simple insertion sort (efficient for small arrays)
	for i := 1; i < len(g.ColorStops); i++ {
		j := i
		for j > 0 && g.ColorStops[j-1].Position > g.ColorStops[j].Position {
			g.ColorStops[j-1], g.ColorStops[j] = g.ColorStops[j], g.ColorStops[j-1]
			j--
		}
	}
}

// Validate validates the gradient configuration.
//
// Checks:
//   - At least 2 color stops are defined
//   - Color stops are in range [0, 1]
//   - For linear gradients: start and end points are different
//   - For radial gradients: radii are non-negative
//
// Returns an error if validation fails.
func (g *Gradient) Validate() error {
	// Check minimum color stops
	if len(g.ColorStops) < 2 {
		return errors.New("gradient must have at least 2 color stops")
	}

	// Validate color stops
	for i, stop := range g.ColorStops {
		if stop.Position < 0.0 || stop.Position > 1.0 {
			return fmt.Errorf("color stop %d: position must be in range [0, 1], got: %f",
				i, stop.Position)
		}
		if err := validateColor(stop.Color); err != nil {
			return fmt.Errorf("color stop %d: %w", i, err)
		}
	}

	// Type-specific validation
	switch g.Type {
	case GradientTypeLinear:
		return g.validateLinear()
	case GradientTypeRadial:
		return g.validateRadial()
	default:
		return fmt.Errorf("unknown gradient type: %d", g.Type)
	}
}

// validateLinear validates linear gradient configuration.
func (g *Gradient) validateLinear() error {
	// Check that start and end points are different
	if g.X1 == g.X2 && g.Y1 == g.Y2 {
		return errors.New("linear gradient: start and end points must be different")
	}
	return nil
}

// validateRadial validates radial gradient configuration.
func (g *Gradient) validateRadial() error {
	// Check that radii are non-negative
	if g.R0 < 0 {
		return fmt.Errorf("radial gradient: starting radius must be non-negative, got: %f", g.R0)
	}
	if g.R1 < 0 {
		return fmt.Errorf("radial gradient: ending radius must be non-negative, got: %f", g.R1)
	}

	// At least one radius must be positive (otherwise degenerate)
	if g.R0 == 0 && g.R1 == 0 {
		return errors.New("radial gradient: at least one radius must be positive")
	}

	return nil
}
