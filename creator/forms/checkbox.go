package forms

import (
	"errors"
)

// Checkbox represents a checkbox field in a PDF form.
//
// Checkboxes allow users to select or deselect an option.
// They are button fields with /FT /Btn and no special flags (radio/pushbutton).
//
// Example:
//
//	// Create checkbox
//	agreeBox := forms.NewCheckbox("agree", 100, 650, 15, 15)
//	agreeBox.SetLabel("I agree to terms")
//	agreeBox.SetChecked(true)
//
//	// Styled checkbox
//	subscribeBox := forms.NewCheckbox("subscribe", 100, 600, 20, 20)
//	subscribeBox.SetLabel("Subscribe to newsletter")
//	subscribeBox.SetBorderColor(0, 0, 0) // Black border
//	subscribeBox.SetFillColor(1, 1, 1)   // White background
//
// PDF Structure:
//
//	<< /Type /Annot
//	   /Subtype /Widget
//	   /FT /Btn                  % Field Type: Button
//	   /T (agree)                % Field name
//	   /V /Yes                   % Value when checked (/Yes or /Off)
//	   /Rect [100 650 115 665]   % Position
//	   /F 4                      % Print flag
//	   /Ff 0                     % No flags = checkbox (not radio/pushbutton)
//	   /AS /Yes                  % Appearance state (current state)
//	   /AP <<                    % Appearance dictionary (optional)
//	       /N << /Yes 10 0 R /Off 11 0 R >>
//	   >>
//	>>
type Checkbox struct {
	// Required fields
	name    string     // Field name (unique identifier)
	rect    [4]float64 // [x, y, x+width, y+height]
	checked bool       // Current checked state

	// Optional fields
	label          string // Label text (displayed next to checkbox)
	defaultChecked bool   // Default checked state

	// Flags
	flags int // Field flags bitmask (no special flags for checkbox)

	// Appearance
	borderColor *[3]float64 // RGB border color (nil = no border)
	fillColor   *[3]float64 // RGB fill color (nil = no fill)
}

// NewCheckbox creates a new checkbox field at the specified position.
//
// Parameters:
//   - name: Unique field name (used for form data)
//   - x: Left edge position in points
//   - y: Bottom edge position in points
//   - width: Checkbox width in points (typically 12-20)
//   - height: Checkbox height in points (typically 12-20)
//
// Example:
//
//	checkbox := forms.NewCheckbox("terms", 100, 650, 15, 15)
func NewCheckbox(name string, x, y, width, height float64) *Checkbox {
	return &Checkbox{
		name:           name,
		rect:           [4]float64{x, y, x + width, y + height},
		checked:        false,
		label:          "",
		defaultChecked: false,
		flags:          0, // No special flags for checkbox
		borderColor:    nil,
		fillColor:      nil,
	}
}

// Name returns the field name.
func (c *Checkbox) Name() string {
	return c.name
}

// Type returns the PDF field type (/FT value).
// For checkboxes, this is always "Btn" (button).
func (c *Checkbox) Type() string {
	return "Btn"
}

// Rect returns the field's bounding rectangle [x1, y1, x2, y2].
func (c *Checkbox) Rect() [4]float64 {
	return c.rect
}

// Flags returns the field flags bitmask.
// For checkboxes, this is typically 0 (no special flags).
// Radio and pushbutton flags are NOT set.
func (c *Checkbox) Flags() int {
	return c.flags
}

// Value returns the field's current value.
// For checkboxes, this is either "Yes" (checked) or "Off" (unchecked).
func (c *Checkbox) Value() interface{} {
	if c.checked {
		return "Yes"
	}
	return "Off"
}

// DefaultValue returns the field's default value.
func (c *Checkbox) DefaultValue() interface{} {
	if c.defaultChecked {
		return "Yes"
	}
	return "Off"
}

// IsReadOnly returns true if the field is read-only.
func (c *Checkbox) IsReadOnly() bool {
	return c.flags&FlagReadOnly != 0
}

// IsRequired returns true if the field is required.
func (c *Checkbox) IsRequired() bool {
	return c.flags&FlagRequired != 0
}

// IsChecked returns true if the checkbox is checked.
func (c *Checkbox) IsChecked() bool {
	return c.checked
}

// SetChecked sets whether the checkbox is checked.
//
// Example:
//
//	checkbox.SetChecked(true)  // Check the box
func (c *Checkbox) SetChecked(checked bool) *Checkbox {
	c.checked = checked
	return c
}

// SetLabel sets the label text displayed next to the checkbox.
//
// Note: This is for documentation/accessibility purposes.
// The actual label rendering is handled by the application.
//
// Example:
//
//	checkbox.SetLabel("I agree to the terms and conditions")
func (c *Checkbox) SetLabel(label string) *Checkbox {
	c.label = label
	return c
}

// Label returns the label text.
func (c *Checkbox) Label() string {
	return c.label
}

// SetDefaultChecked sets the default checked state.
//
// This is used when the form is reset.
//
// Example:
//
//	checkbox.SetDefaultChecked(false)  // Unchecked by default
func (c *Checkbox) SetDefaultChecked(checked bool) *Checkbox {
	c.defaultChecked = checked
	return c
}

// SetReadOnly sets whether the field is read-only.
//
// Example:
//
//	checkbox.SetReadOnly(true)  // Field cannot be changed
func (c *Checkbox) SetReadOnly(readonly bool) *Checkbox {
	if readonly {
		c.flags |= FlagReadOnly
	} else {
		c.flags &^= FlagReadOnly
	}
	return c
}

// SetRequired sets whether the field is required.
//
// Example:
//
//	checkbox.SetRequired(true)  // User must check this box
func (c *Checkbox) SetRequired(required bool) *Checkbox {
	if required {
		c.flags |= FlagRequired
	} else {
		c.flags &^= FlagRequired
	}
	return c
}

// SetBorderColor sets the border color (RGB, 0.0-1.0 range).
//
// Set to nil to remove border.
//
// Example:
//
//	checkbox.SetBorderColor(0, 0, 0)  // Black border
func (c *Checkbox) SetBorderColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		c.borderColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	c.borderColor = &[3]float64{*r, *g, *b}
	return nil
}

// BorderColor returns the border color (nil if no border).
func (c *Checkbox) BorderColor() *[3]float64 {
	return c.borderColor
}

// SetFillColor sets the background fill color (RGB, 0.0-1.0 range).
//
// Set to nil for transparent background.
//
// Example:
//
//	checkbox.SetFillColor(1, 1, 1)  // White background
func (c *Checkbox) SetFillColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		c.fillColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	c.fillColor = &[3]float64{*r, *g, *b}
	return nil
}

// FillColor returns the fill color (nil if transparent).
func (c *Checkbox) FillColor() *[3]float64 {
	return c.fillColor
}

// Validate checks if the field configuration is valid.
//
// Returns an error if:
//   - Name is empty
//   - Rectangle has invalid dimensions
func (c *Checkbox) Validate() error {
	if c.name == "" {
		return errors.New("field name cannot be empty")
	}

	// Validate rectangle
	if c.rect[2] <= c.rect[0] || c.rect[3] <= c.rect[1] {
		return errors.New("invalid rectangle: width and height must be positive")
	}

	return nil
}
