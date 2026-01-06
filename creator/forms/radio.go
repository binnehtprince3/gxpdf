package forms

import (
	"errors"
	"fmt"
)

// RadioGroup represents a group of radio buttons in a PDF form.
//
// Radio buttons allow users to select exactly one option from a group.
// They are button fields with /FT /Btn and radio-specific flags.
//
// Unlike checkboxes, radio buttons are represented as a parent field
// with multiple child widget annotations (one per option).
//
// Example:
//
//	// Create radio group
//	gender := forms.NewRadioGroup("gender")
//	gender.AddOption("male", 100, 600, "Male")
//	gender.AddOption("female", 200, 600, "Female")
//	gender.AddOption("other", 300, 600, "Other")
//	gender.SetSelected("male")
//
//	// Styled radio group
//	priority := forms.NewRadioGroup("priority")
//	priority.AddOption("low", 100, 550, "Low")
//	priority.AddOption("medium", 200, 550, "Medium")
//	priority.AddOption("high", 300, 550, "High")
//	priority.SetBorderColor(0, 0, 0) // Black border
//	priority.SetFillColor(1, 1, 1)   // White background
//
// PDF Structure:
//
//	% Parent radio group
//	<< /FT /Btn
//	   /T (gender)                % Field name
//	   /V /male                   % Selected value
//	   /Ff 49152                  % Flags: Radio (32768) + NoToggleToOff (16384)
//	   /Kids [101 0 R 102 0 R 103 0 R]
//	>>
//
//	% Child widget (one per option)
//	101 0 obj
//	<< /Type /Annot
//	   /Subtype /Widget
//	   /Parent 100 0 R
//	   /T (male)
//	   /Rect [100 600 115 615]
//	   /AS /male                  % Appearance state
//	   /AP << /N << /male 10 0 R /Off 11 0 R >> >>
//	>>
type RadioGroup struct {
	// Required fields
	name    string         // Field name (unique identifier)
	options []*RadioOption // Radio button options

	// Optional fields
	selected        string // Currently selected option value
	defaultSelected string // Default selected option

	// Flags (applied to parent field)
	flags int // Field flags: FlagRadio | FlagNoToggleToOff

	// Appearance (applied to all child widgets)
	borderColor *[3]float64 // RGB border color (nil = no border)
	fillColor   *[3]float64 // RGB fill color (nil = no fill)
}

// RadioOption represents a single radio button option.
type RadioOption struct {
	value string     // Option value (e.g., "male", "female")
	label string     // Display label (e.g., "Male", "Female")
	rect  [4]float64 // [x, y, x+width, y+height]
}

// NewRadioGroup creates a new radio button group.
//
// Parameters:
//   - name: Unique field name (used for form data)
//
// Example:
//
//	radioGroup := forms.NewRadioGroup("payment_method")
func NewRadioGroup(name string) *RadioGroup {
	return &RadioGroup{
		name:            name,
		options:         make([]*RadioOption, 0),
		selected:        "",
		defaultSelected: "",
		// Radio buttons typically have Radio + NoToggleToOff flags
		flags:       FlagRadio | FlagNoToggleToOff,
		borderColor: nil,
		fillColor:   nil,
	}
}

// Name returns the field name.
func (r *RadioGroup) Name() string {
	return r.name
}

// Type returns the PDF field type (/FT value).
// For radio buttons, this is always "Btn" (button).
func (r *RadioGroup) Type() string {
	return "Btn"
}

// Rect returns the bounding rectangle of the first option.
// For radio groups with multiple options, this returns the rect of the first option.
// PDF viewers typically ignore this for parent fields with /Kids.
func (r *RadioGroup) Rect() [4]float64 {
	if len(r.options) == 0 {
		return [4]float64{0, 0, 0, 0}
	}
	return r.options[0].rect
}

// Flags returns the field flags bitmask.
// For radio buttons, this includes FlagRadio and FlagNoToggleToOff.
func (r *RadioGroup) Flags() int {
	return r.flags
}

// Value returns the field's current selected value.
// For radio buttons, this is the value of the selected option (e.g., "male").
// Returns empty string if no option is selected.
func (r *RadioGroup) Value() interface{} {
	return r.selected
}

// DefaultValue returns the field's default selected value.
func (r *RadioGroup) DefaultValue() interface{} {
	return r.defaultSelected
}

// IsReadOnly returns true if the field is read-only.
func (r *RadioGroup) IsReadOnly() bool {
	return r.flags&FlagReadOnly != 0
}

// IsRequired returns true if the field is required.
func (r *RadioGroup) IsRequired() bool {
	return r.flags&FlagRequired != 0
}

// AddOption adds a radio button option to the group.
//
// Parameters:
//   - value: Option value (e.g., "male", "female")
//   - x: Left edge position in points
//   - y: Bottom edge position in points
//   - label: Display label (e.g., "Male", "Female")
//   - width: Button width in points (default: 15)
//   - height: Button height in points (default: 15)
//
// The width and height are optional. If not provided, default to 15x15.
//
// Example:
//
//	gender.AddOption("male", 100, 600, "Male")
//	gender.AddOption("female", 200, 600, "Female")
func (r *RadioGroup) AddOption(value string, x, y float64, label string, dimensions ...float64) *RadioGroup {
	width, height := 15.0, 15.0
	if len(dimensions) >= 1 {
		width = dimensions[0]
	}
	if len(dimensions) >= 2 {
		height = dimensions[1]
	}

	option := &RadioOption{
		value: value,
		label: label,
		rect:  [4]float64{x, y, x + width, y + height},
	}

	r.options = append(r.options, option)
	return r
}

// Options returns all radio button options.
func (r *RadioGroup) Options() []*RadioOption {
	return r.options
}

// SetSelected sets the currently selected option.
//
// The value must match one of the option values added via AddOption.
// If the value doesn't exist, this method does nothing.
//
// Example:
//
//	gender.SetSelected("male")
func (r *RadioGroup) SetSelected(value string) error {
	// Validate that the value exists in options
	found := false
	for _, opt := range r.options {
		if opt.value == value {
			found = true
			break
		}
	}

	if !found && value != "" {
		return fmt.Errorf("option value '%s' not found in radio group '%s'", value, r.name)
	}

	r.selected = value
	return nil
}

// Selected returns the currently selected option value.
func (r *RadioGroup) Selected() string {
	return r.selected
}

// SetDefaultSelected sets the default selected option.
//
// This is used when the form is reset.
//
// Example:
//
//	gender.SetDefaultSelected("male")
func (r *RadioGroup) SetDefaultSelected(value string) error {
	// Validate that the value exists in options
	found := false
	for _, opt := range r.options {
		if opt.value == value {
			found = true
			break
		}
	}

	if !found && value != "" {
		return fmt.Errorf("option value '%s' not found in radio group '%s'", value, r.name)
	}

	r.defaultSelected = value
	return nil
}

// SetReadOnly sets whether the field is read-only.
//
// Example:
//
//	radioGroup.SetReadOnly(true)  // Field cannot be changed
func (r *RadioGroup) SetReadOnly(readonly bool) *RadioGroup {
	if readonly {
		r.flags |= FlagReadOnly
	} else {
		r.flags &^= FlagReadOnly
	}
	return r
}

// SetRequired sets whether the field is required.
//
// Example:
//
//	radioGroup.SetRequired(true)  // User must select an option
func (r *RadioGroup) SetRequired(required bool) *RadioGroup {
	if required {
		r.flags |= FlagRequired
	} else {
		r.flags &^= FlagRequired
	}
	return r
}

// SetAllowToggleOff allows deselecting all radio buttons.
//
// By default (NoToggleToOff flag), once a radio button is selected,
// the user cannot deselect all options. Setting this to true removes
// that restriction.
//
// Example:
//
//	radioGroup.SetAllowToggleOff(true)  // Allow deselecting all
func (r *RadioGroup) SetAllowToggleOff(allow bool) *RadioGroup {
	if allow {
		r.flags &^= FlagNoToggleToOff // Remove NoToggleToOff flag
	} else {
		r.flags |= FlagNoToggleToOff // Add NoToggleToOff flag
	}
	return r
}

// SetBorderColor sets the border color for all radio buttons (RGB, 0.0-1.0 range).
//
// Set to nil to remove border.
//
// Example:
//
//	r, g, b := 0.0, 0.0, 0.0
//	radioGroup.SetBorderColor(&r, &g, &b)  // Black border
func (r *RadioGroup) SetBorderColor(rgbPtr ...*float64) error {
	color, err := parseColorPointers(rgbPtr)
	if err != nil {
		return err
	}
	r.borderColor = color
	return nil
}

// BorderColor returns the border color (nil if no border).
func (r *RadioGroup) BorderColor() *[3]float64 {
	return r.borderColor
}

// SetFillColor sets the background fill color for all radio buttons (RGB, 0.0-1.0 range).
//
// Set to nil for transparent background.
//
// Example:
//
//	r, g, b := 1.0, 1.0, 1.0
//	radioGroup.SetFillColor(&r, &g, &b)  // White background
func (r *RadioGroup) SetFillColor(rgbPtr ...*float64) error {
	color, err := parseColorPointers(rgbPtr)
	if err != nil {
		return err
	}
	r.fillColor = color
	return nil
}

// FillColor returns the fill color (nil if transparent).
func (r *RadioGroup) FillColor() *[3]float64 {
	return r.fillColor
}

// Validate checks if the field configuration is valid.
//
// Returns an error if:
//   - Name is empty
//   - No options have been added
//   - Any option has an invalid rectangle
//   - Selected value doesn't match any option
func (r *RadioGroup) Validate() error {
	if r.name == "" {
		return errors.New("field name cannot be empty")
	}

	if len(r.options) == 0 {
		return errors.New("radio group must have at least one option")
	}

	// Validate all option rectangles
	for i, opt := range r.options {
		if opt.rect[2] <= opt.rect[0] || opt.rect[3] <= opt.rect[1] {
			return fmt.Errorf("option %d has invalid rectangle: width and height must be positive", i)
		}
	}

	// Validate selected value exists in options
	if r.selected != "" {
		found := false
		for _, opt := range r.options {
			if opt.value == r.selected {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("selected value '%s' not found in options", r.selected)
		}
	}

	return nil
}

// Value returns the option value.
func (o *RadioOption) Value() string {
	return o.value
}

// Label returns the option label.
func (o *RadioOption) Label() string {
	return o.label
}

// Rect returns the option rectangle [x1, y1, x2, y2].
func (o *RadioOption) Rect() [4]float64 {
	return o.rect
}

// parseColorPointers parses variadic color pointer arguments.
// Returns nil if no arguments or all nil pointers.
// Returns error if invalid number of arguments or values out of range.
func parseColorPointers(rgbPtr []*float64) (*[3]float64, error) {
	// Handle empty or nil cases
	if shouldReturnNilColor(rgbPtr) {
		return nil, nil //nolint:nilnil // nil color is a valid value (no color)
	}

	// Validate component count
	if len(rgbPtr) < 3 {
		return nil, errors.New("color requires 3 RGB components")
	}

	// Validate component values
	return validateColorComponents(rgbPtr[0], rgbPtr[1], rgbPtr[2])
}

// shouldReturnNilColor checks if color pointers should result in nil.
func shouldReturnNilColor(rgbPtr []*float64) bool {
	if len(rgbPtr) == 0 {
		return true
	}
	if rgbPtr[0] == nil {
		return true
	}
	if len(rgbPtr) >= 3 && (rgbPtr[0] == nil || rgbPtr[1] == nil || rgbPtr[2] == nil) {
		return true
	}
	return false
}

// validateColorComponents validates RGB component values.
func validateColorComponents(r, g, b *float64) (*[3]float64, error) {
	if r == nil || g == nil || b == nil {
		return nil, nil //nolint:nilnil // nil color is a valid value (no color)
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return nil, errors.New("color components must be in range [0.0, 1.0]")
	}
	return &[3]float64{*r, *g, *b}, nil
}
