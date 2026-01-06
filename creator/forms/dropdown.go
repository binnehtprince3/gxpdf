package forms

import (
	"errors"
	"fmt"
)

// Option represents a single option in a dropdown or listbox.
// It consists of an export value (used in form data) and a display value (shown to user).
type Option struct {
	ExportValue  string // Value exported in form data
	DisplayValue string // Value displayed to user
}

// Dropdown represents a dropdown (combo box) field in a PDF form.
//
// Dropdowns allow users to select a single option from a list.
// They can optionally be editable, allowing users to enter custom values.
//
// Example:
//
//	// Simple dropdown
//	country := forms.NewDropdown("country", 100, 550, 150, 20)
//	country.AddOption("us", "United States")
//	country.AddOption("ca", "Canada")
//	country.AddOption("uk", "United Kingdom")
//	country.SetSelected("us")
//
//	// Editable dropdown (user can type custom value)
//	customDropdown := forms.NewDropdown("custom", 100, 500, 150, 20)
//	customDropdown.SetEditable(true)
//	customDropdown.AddOptions("Option 1", "Option 2", "Option 3")
//
// PDF Structure:
//
//	<< /Type /Annot
//	   /Subtype /Widget
//	   /FT /Ch                          % Field Type: Choice
//	   /T (country)                     % Field name
//	   /V (us)                          % Selected value
//	   /Opt [                           % Options array
//	       [(us) (United States)]
//	       [(ca) (Canada)]
//	       [(uk) (United Kingdom)]
//	   ]
//	   /Rect [100 550 250 570]          % Position
//	   /F 4                             % Print flag
//	   /Ff 131072                       % Combo flag (bit 18)
//	>>
type Dropdown struct {
	// Required fields
	name    string     // Field name (unique identifier)
	rect    [4]float64 // [x, y, x+width, y+height]
	options []Option   // List of available options

	// Selection
	selected string // Currently selected value (export value)

	// Optional fields
	defaultValue string // Default selected value

	// Flags
	flags int // Field flags bitmask (includes Combo flag)

	// Appearance
	fontSize    float64     // Font size for text (default: 12)
	fontName    string      // Font name (default: Helvetica)
	textColor   [3]float64  // RGB text color (default: black)
	borderColor *[3]float64 // RGB border color (nil = no border)
	fillColor   *[3]float64 // RGB fill color (nil = no fill)
}

// NewDropdown creates a new dropdown field at the specified position.
//
// Parameters:
//   - name: Unique field name (used for form data)
//   - x: Left edge position in points
//   - y: Bottom edge position in points
//   - width: Field width in points
//   - height: Field height in points (typically 15-25 for single-line dropdown)
//
// Example:
//
//	dropdown := forms.NewDropdown("country", 100, 550, 150, 20)
func NewDropdown(name string, x, y, width, height float64) *Dropdown {
	return &Dropdown{
		name:         name,
		rect:         [4]float64{x, y, x + width, y + height},
		options:      make([]Option, 0),
		selected:     "",
		defaultValue: "",
		flags:        FlagCombo, // Combo flag (bit 18) - makes it a dropdown
		fontSize:     12,
		fontName:     "Helvetica",
		textColor:    [3]float64{0, 0, 0}, // Black
		borderColor:  nil,
		fillColor:    nil,
	}
}

// Name returns the field name.
func (d *Dropdown) Name() string {
	return d.name
}

// Type returns the PDF field type (/FT value).
// For dropdowns, this is always "Ch" (choice).
func (d *Dropdown) Type() string {
	return "Ch"
}

// Rect returns the field's bounding rectangle [x1, y1, x2, y2].
func (d *Dropdown) Rect() [4]float64 {
	return d.rect
}

// Flags returns the field flags bitmask.
// For dropdowns, this includes FlagCombo (bit 18).
func (d *Dropdown) Flags() int {
	return d.flags
}

// Value returns the field's current selected value.
// For dropdowns, this is the export value of the selected option.
func (d *Dropdown) Value() interface{} {
	return d.selected
}

// DefaultValue returns the field's default value.
func (d *Dropdown) DefaultValue() interface{} {
	return d.defaultValue
}

// IsReadOnly returns true if the field is read-only.
func (d *Dropdown) IsReadOnly() bool {
	return d.flags&FlagReadOnly != 0
}

// IsRequired returns true if the field is required.
func (d *Dropdown) IsRequired() bool {
	return d.flags&FlagRequired != 0
}

// IsEditable returns true if users can enter custom values.
func (d *Dropdown) IsEditable() bool {
	return d.flags&FlagEdit != 0
}

// IsSorted returns true if options are sorted alphabetically.
func (d *Dropdown) IsSorted() bool {
	return d.flags&FlagSort != 0
}

// Options returns the list of available options.
func (d *Dropdown) Options() []Option {
	return d.options
}

// SelectedValue returns the currently selected value.
func (d *Dropdown) SelectedValue() string {
	return d.selected
}

// AddOption adds a single option to the dropdown.
//
// Parameters:
//   - exportValue: Value used in form data export
//   - displayValue: Value shown to user in dropdown
//
// Example:
//
//	dropdown.AddOption("us", "United States")
//	dropdown.AddOption("ca", "Canada")
func (d *Dropdown) AddOption(exportValue, displayValue string) *Dropdown {
	d.options = append(d.options, Option{
		ExportValue:  exportValue,
		DisplayValue: displayValue,
	})
	return d
}

// AddOptions adds multiple options where export value equals display value.
//
// This is a convenience method for simple dropdowns where the internal
// value is the same as what's displayed.
//
// Example:
//
//	dropdown.AddOptions("Red", "Green", "Blue")
func (d *Dropdown) AddOptions(values ...string) *Dropdown {
	for _, value := range values {
		d.options = append(d.options, Option{
			ExportValue:  value,
			DisplayValue: value,
		})
	}
	return d
}

// SetSelected sets the currently selected option by export value.
//
// Example:
//
//	dropdown.SetSelected("us")  // Select "United States"
func (d *Dropdown) SetSelected(exportValue string) error {
	// Validate that the value exists in options
	if len(d.options) > 0 {
		found := false
		for _, opt := range d.options {
			if opt.ExportValue == exportValue {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value '%s' not found in dropdown options", exportValue)
		}
	}

	d.selected = exportValue
	return nil
}

// SetDefaultValue sets the default selected value.
//
// This is used when the form is reset.
//
// Example:
//
//	dropdown.SetDefaultValue("us")
func (d *Dropdown) SetDefaultValue(exportValue string) error {
	// Validate that the value exists in options
	if len(d.options) > 0 {
		found := false
		for _, opt := range d.options {
			if opt.ExportValue == exportValue {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value '%s' not found in dropdown options", exportValue)
		}
	}

	d.defaultValue = exportValue
	return nil
}

// SetEditable sets whether users can enter custom values.
//
// When editable, the dropdown becomes a combo box where users can
// either select from the list or type a custom value.
//
// Example:
//
//	dropdown.SetEditable(true)  // Allow custom text entry
func (d *Dropdown) SetEditable(editable bool) *Dropdown {
	if editable {
		d.flags |= FlagEdit // Bit 19
	} else {
		d.flags &^= FlagEdit
	}
	return d
}

// SetSort sets whether options are sorted alphabetically.
//
// Example:
//
//	dropdown.SetSort(true)  // Sort options A-Z
func (d *Dropdown) SetSort(sort bool) *Dropdown {
	if sort {
		d.flags |= FlagSort // Bit 20
	} else {
		d.flags &^= FlagSort
	}
	return d
}

// SetReadOnly sets whether the field is read-only.
//
// Example:
//
//	dropdown.SetReadOnly(true)  // Field cannot be changed
func (d *Dropdown) SetReadOnly(readonly bool) *Dropdown {
	if readonly {
		d.flags |= FlagReadOnly
	} else {
		d.flags &^= FlagReadOnly
	}
	return d
}

// SetRequired sets whether the field is required.
//
// Example:
//
//	dropdown.SetRequired(true)  // User must select an option
func (d *Dropdown) SetRequired(required bool) *Dropdown {
	if required {
		d.flags |= FlagRequired
	} else {
		d.flags &^= FlagRequired
	}
	return d
}

// SetFontSize sets the font size for the text.
//
// Example:
//
//	dropdown.SetFontSize(14)
func (d *Dropdown) SetFontSize(size float64) error {
	if size <= 0 {
		return errors.New("font size must be positive")
	}
	d.fontSize = size
	return nil
}

// FontSize returns the font size.
func (d *Dropdown) FontSize() float64 {
	return d.fontSize
}

// SetFontName sets the font name.
//
// Common values: "Helvetica", "Courier", "Times-Roman"
//
// Example:
//
//	dropdown.SetFontName("Courier")
func (d *Dropdown) SetFontName(name string) *Dropdown {
	d.fontName = name
	return d
}

// FontName returns the font name.
func (d *Dropdown) FontName() string {
	return d.fontName
}

// SetTextColor sets the text color (RGB, 0.0-1.0 range).
//
// Example:
//
//	dropdown.SetTextColor(0, 0, 1)  // Blue text
func (d *Dropdown) SetTextColor(r, g, b float64) error {
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	d.textColor = [3]float64{r, g, b}
	return nil
}

// TextColor returns the text color.
func (d *Dropdown) TextColor() [3]float64 {
	return d.textColor
}

// SetBorderColor sets the border color (RGB, 0.0-1.0 range).
//
// Set to nil to remove border.
//
// Example:
//
//	dropdown.SetBorderColor(0, 0, 0)  // Black border
func (d *Dropdown) SetBorderColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		d.borderColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	d.borderColor = &[3]float64{*r, *g, *b}
	return nil
}

// BorderColor returns the border color (nil if no border).
func (d *Dropdown) BorderColor() *[3]float64 {
	return d.borderColor
}

// SetFillColor sets the background fill color (RGB, 0.0-1.0 range).
//
// Set to nil for transparent background.
//
// Example:
//
//	dropdown.SetFillColor(1, 1, 1)  // White background
func (d *Dropdown) SetFillColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		d.fillColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	d.fillColor = &[3]float64{*r, *g, *b}
	return nil
}

// FillColor returns the fill color (nil if transparent).
func (d *Dropdown) FillColor() *[3]float64 {
	return d.fillColor
}

// Validate checks if the field configuration is valid.
//
// Returns an error if:
//   - Name is empty
//   - Rectangle has invalid dimensions
//   - Selected value is not in options
func (d *Dropdown) Validate() error {
	if d.name == "" {
		return errors.New("field name cannot be empty")
	}

	// Validate rectangle
	if d.rect[2] <= d.rect[0] || d.rect[3] <= d.rect[1] {
		return errors.New("invalid rectangle: width and height must be positive")
	}

	// Validate selected value exists in options (if options are defined)
	if d.selected != "" && len(d.options) > 0 {
		found := false
		for _, opt := range d.options {
			if opt.ExportValue == d.selected {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("selected value '%s' not found in options", d.selected)
		}
	}

	return nil
}
