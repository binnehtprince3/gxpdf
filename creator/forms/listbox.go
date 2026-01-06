package forms

import (
	"errors"
	"fmt"
)

// ListBox represents a list box field in a PDF form.
//
// List boxes display a scrollable list of options and can optionally
// allow multiple selections.
//
// Example:
//
//	// Single-select list
//	favColor := forms.NewListBox("favorite_color", 100, 450, 150, 80)
//	favColor.AddOptions("Red", "Green", "Blue", "Yellow", "Orange")
//	favColor.SetSelected("Blue")
//
//	// Multi-select list
//	interests := forms.NewListBox("interests", 100, 350, 150, 80)
//	interests.AddOptions("Sports", "Music", "Movies", "Reading", "Travel")
//	interests.SetMultiSelect(true)
//	interests.SetSelectedMultiple("Sports", "Music")
//
// PDF Structure:
//
//	<< /Type /Annot
//	   /Subtype /Widget
//	   /FT /Ch                          % Field Type: Choice
//	   /T (interests)                   % Field name
//	   /V [(Sports) (Music)]            % Selected values (array for multi-select)
//	   /Opt [(Sports) (Music) (Movies) (Reading) (Travel)]  % Options
//	   /Rect [100 350 250 430]          % Position
//	   /F 4                             % Print flag
//	   /Ff 2097152                      % MultiSelect flag (bit 22)
//	>>
type ListBox struct {
	// Required fields
	name    string     // Field name (unique identifier)
	rect    [4]float64 // [x, y, x+width, y+height]
	options []Option   // List of available options

	// Selection
	selected []string // Currently selected values (export values)

	// Optional fields
	defaultValue []string // Default selected values

	// Flags
	flags int // Field flags bitmask (NO Combo flag, optional MultiSelect)

	// Appearance
	fontSize    float64     // Font size for text (default: 12)
	fontName    string      // Font name (default: Helvetica)
	textColor   [3]float64  // RGB text color (default: black)
	borderColor *[3]float64 // RGB border color (nil = no border)
	fillColor   *[3]float64 // RGB fill color (nil = no fill)
}

// NewListBox creates a new list box field at the specified position.
//
// Parameters:
//   - name: Unique field name (used for form data)
//   - x: Left edge position in points
//   - y: Bottom edge position in points
//   - width: Field width in points
//   - height: Field height in points (should be tall enough for multiple items)
//
// Example:
//
//	listbox := forms.NewListBox("colors", 100, 350, 150, 80)
func NewListBox(name string, x, y, width, height float64) *ListBox {
	return &ListBox{
		name:         name,
		rect:         [4]float64{x, y, x + width, y + height},
		options:      make([]Option, 0),
		selected:     make([]string, 0),
		defaultValue: make([]string, 0),
		flags:        0, // NO Combo flag - this is a list box, not a dropdown
		fontSize:     12,
		fontName:     "Helvetica",
		textColor:    [3]float64{0, 0, 0}, // Black
		borderColor:  nil,
		fillColor:    nil,
	}
}

// Name returns the field name.
func (l *ListBox) Name() string {
	return l.name
}

// Type returns the PDF field type (/FT value).
// For list boxes, this is always "Ch" (choice).
func (l *ListBox) Type() string {
	return "Ch"
}

// Rect returns the field's bounding rectangle [x1, y1, x2, y2].
func (l *ListBox) Rect() [4]float64 {
	return l.rect
}

// Flags returns the field flags bitmask.
// For list boxes, Combo flag is NOT set (unlike dropdowns).
func (l *ListBox) Flags() int {
	return l.flags
}

// Value returns the field's current selected value(s).
// For single-select, returns a string.
// For multi-select, returns []string.
func (l *ListBox) Value() interface{} {
	if l.IsMultiSelect() {
		return l.selected // Return array for multi-select
	}
	if len(l.selected) > 0 {
		return l.selected[0] // Return single value
	}
	return ""
}

// DefaultValue returns the field's default value(s).
func (l *ListBox) DefaultValue() interface{} {
	if l.IsMultiSelect() {
		return l.defaultValue
	}
	if len(l.defaultValue) > 0 {
		return l.defaultValue[0]
	}
	return ""
}

// IsReadOnly returns true if the field is read-only.
func (l *ListBox) IsReadOnly() bool {
	return l.flags&FlagReadOnly != 0
}

// IsRequired returns true if the field is required.
func (l *ListBox) IsRequired() bool {
	return l.flags&FlagRequired != 0
}

// IsMultiSelect returns true if multiple selections are allowed.
func (l *ListBox) IsMultiSelect() bool {
	return l.flags&FlagMultiSelect != 0
}

// IsSorted returns true if options are sorted alphabetically.
func (l *ListBox) IsSorted() bool {
	return l.flags&FlagSort != 0
}

// Options returns the list of available options.
func (l *ListBox) Options() []Option {
	return l.options
}

// SelectedValues returns the currently selected values.
func (l *ListBox) SelectedValues() []string {
	return l.selected
}

// AddOption adds a single option to the list box.
//
// Parameters:
//   - exportValue: Value used in form data export
//   - displayValue: Value shown to user in list
//
// Example:
//
//	listbox.AddOption("red", "Red Color")
//	listbox.AddOption("blue", "Blue Color")
func (l *ListBox) AddOption(exportValue, displayValue string) *ListBox {
	l.options = append(l.options, Option{
		ExportValue:  exportValue,
		DisplayValue: displayValue,
	})
	return l
}

// AddOptions adds multiple options where export value equals display value.
//
// This is a convenience method for simple list boxes where the internal
// value is the same as what's displayed.
//
// Example:
//
//	listbox.AddOptions("Red", "Green", "Blue")
func (l *ListBox) AddOptions(values ...string) *ListBox {
	for _, value := range values {
		l.options = append(l.options, Option{
			ExportValue:  value,
			DisplayValue: value,
		})
	}
	return l
}

// SetSelected sets a single selected option by export value.
//
// This is for single-select list boxes (when MultiSelect is false).
//
// Example:
//
//	listbox.SetSelected("blue")
func (l *ListBox) SetSelected(exportValue string) error {
	// Validate that the value exists in options
	if len(l.options) > 0 {
		found := false
		for _, opt := range l.options {
			if opt.ExportValue == exportValue {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value '%s' not found in list box options", exportValue)
		}
	}

	l.selected = []string{exportValue}
	return nil
}

// SetSelectedMultiple sets multiple selected options by export values.
//
// This is for multi-select list boxes (when MultiSelect is true).
//
// Example:
//
//	listbox.SetMultiSelect(true)
//	listbox.SetSelectedMultiple("sports", "music", "movies")
func (l *ListBox) SetSelectedMultiple(exportValues ...string) error {
	// Validate that all values exist in options
	if len(l.options) > 0 {
		for _, value := range exportValues {
			found := false
			for _, opt := range l.options {
				if opt.ExportValue == value {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("value '%s' not found in list box options", value)
			}
		}
	}

	l.selected = exportValues
	return nil
}

// SetDefaultValue sets the default selected value (single selection).
//
// This is used when the form is reset.
//
// Example:
//
//	listbox.SetDefaultValue("blue")
func (l *ListBox) SetDefaultValue(exportValue string) error {
	// Validate that the value exists in options
	if len(l.options) > 0 {
		found := false
		for _, opt := range l.options {
			if opt.ExportValue == exportValue {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("value '%s' not found in list box options", exportValue)
		}
	}

	l.defaultValue = []string{exportValue}
	return nil
}

// SetDefaultValueMultiple sets the default selected values (multi-select).
//
// Example:
//
//	listbox.SetDefaultValueMultiple("sports", "music")
func (l *ListBox) SetDefaultValueMultiple(exportValues ...string) error {
	// Validate that all values exist in options
	if len(l.options) > 0 {
		for _, value := range exportValues {
			found := false
			for _, opt := range l.options {
				if opt.ExportValue == value {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("value '%s' not found in list box options", value)
			}
		}
	}

	l.defaultValue = exportValues
	return nil
}

// SetMultiSelect sets whether multiple selections are allowed.
//
// When enabled, users can select multiple items from the list using
// Ctrl+Click or Shift+Click.
//
// Example:
//
//	listbox.SetMultiSelect(true)  // Allow multiple selections
func (l *ListBox) SetMultiSelect(multiSelect bool) *ListBox {
	if multiSelect {
		l.flags |= FlagMultiSelect // Bit 22
	} else {
		l.flags &^= FlagMultiSelect
		// If disabling multi-select, keep only first selected item
		if len(l.selected) > 1 {
			l.selected = l.selected[:1]
		}
	}
	return l
}

// SetSort sets whether options are sorted alphabetically.
//
// Example:
//
//	listbox.SetSort(true)  // Sort options A-Z
func (l *ListBox) SetSort(sort bool) *ListBox {
	if sort {
		l.flags |= FlagSort // Bit 20
	} else {
		l.flags &^= FlagSort
	}
	return l
}

// SetReadOnly sets whether the field is read-only.
//
// Example:
//
//	listbox.SetReadOnly(true)  // Field cannot be changed
func (l *ListBox) SetReadOnly(readonly bool) *ListBox {
	if readonly {
		l.flags |= FlagReadOnly
	} else {
		l.flags &^= FlagReadOnly
	}
	return l
}

// SetRequired sets whether the field is required.
//
// Example:
//
//	listbox.SetRequired(true)  // User must select at least one option
func (l *ListBox) SetRequired(required bool) *ListBox {
	if required {
		l.flags |= FlagRequired
	} else {
		l.flags &^= FlagRequired
	}
	return l
}

// SetFontSize sets the font size for the text.
//
// Example:
//
//	listbox.SetFontSize(14)
func (l *ListBox) SetFontSize(size float64) error {
	if size <= 0 {
		return errors.New("font size must be positive")
	}
	l.fontSize = size
	return nil
}

// FontSize returns the font size.
func (l *ListBox) FontSize() float64 {
	return l.fontSize
}

// SetFontName sets the font name.
//
// Common values: "Helvetica", "Courier", "Times-Roman"
//
// Example:
//
//	listbox.SetFontName("Courier")
func (l *ListBox) SetFontName(name string) *ListBox {
	l.fontName = name
	return l
}

// FontName returns the font name.
func (l *ListBox) FontName() string {
	return l.fontName
}

// SetTextColor sets the text color (RGB, 0.0-1.0 range).
//
// Example:
//
//	listbox.SetTextColor(0, 0, 1)  // Blue text
func (l *ListBox) SetTextColor(r, g, b float64) error {
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	l.textColor = [3]float64{r, g, b}
	return nil
}

// TextColor returns the text color.
func (l *ListBox) TextColor() [3]float64 {
	return l.textColor
}

// SetBorderColor sets the border color (RGB, 0.0-1.0 range).
//
// Set to nil to remove border.
//
// Example:
//
//	listbox.SetBorderColor(0, 0, 0)  // Black border
func (l *ListBox) SetBorderColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		l.borderColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	l.borderColor = &[3]float64{*r, *g, *b}
	return nil
}

// BorderColor returns the border color (nil if no border).
func (l *ListBox) BorderColor() *[3]float64 {
	return l.borderColor
}

// SetFillColor sets the background fill color (RGB, 0.0-1.0 range).
//
// Set to nil for transparent background.
//
// Example:
//
//	listbox.SetFillColor(1, 1, 1)  // White background
func (l *ListBox) SetFillColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		l.fillColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	l.fillColor = &[3]float64{*r, *g, *b}
	return nil
}

// FillColor returns the fill color (nil if transparent).
func (l *ListBox) FillColor() *[3]float64 {
	return l.fillColor
}

// Validate checks if the field configuration is valid.
//
// Returns an error if:
//   - Name is empty
//   - Rectangle has invalid dimensions
//   - Selected values are not in options
//   - Multiple selections when MultiSelect is false
func (l *ListBox) Validate() error {
	if l.name == "" {
		return errors.New("field name cannot be empty")
	}

	// Validate rectangle
	if l.rect[2] <= l.rect[0] || l.rect[3] <= l.rect[1] {
		return errors.New("invalid rectangle: width and height must be positive")
	}

	// Validate multi-select constraint
	if !l.IsMultiSelect() && len(l.selected) > 1 {
		return errors.New("multiple selections not allowed when MultiSelect is false")
	}

	// Validate selected values exist in options
	return l.validateSelectedValues()
}

// validateSelectedValues checks if all selected values exist in options.
func (l *ListBox) validateSelectedValues() error {
	if len(l.options) == 0 {
		return nil
	}

	for _, selValue := range l.selected {
		if !l.optionExists(selValue) {
			return fmt.Errorf("selected value '%s' not found in options", selValue)
		}
	}
	return nil
}

// optionExists checks if an export value exists in the options list.
func (l *ListBox) optionExists(exportValue string) bool {
	for _, opt := range l.options {
		if opt.ExportValue == exportValue {
			return true
		}
	}
	return false
}
