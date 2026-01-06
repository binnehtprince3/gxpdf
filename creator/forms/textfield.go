package forms

import (
	"errors"
)

// TextField represents a text input field in a PDF form.
//
// Text fields allow users to enter single-line or multi-line text.
// They support various options like readonly, required, password masking,
// and multiline input.
//
// Example:
//
//	// Single-line text field
//	nameField := forms.NewTextField("name", 100, 700, 200, 20)
//	nameField.SetValue("John Doe")
//	nameField.SetPlaceholder("Enter your name")
//	nameField.SetRequired(true)
//
//	// Multi-line text area
//	commentField := forms.NewTextField("comment", 100, 600, 200, 80)
//	commentField.SetMultiline(true)
//	commentField.SetValue("Enter your comments here...")
//
// PDF Structure:
//
//	<< /Type /Annot
//	   /Subtype /Widget
//	   /FT /Tx                  % Field Type: Text
//	   /T (name)                % Field name
//	   /V (John Doe)            % Field value
//	   /DV (Enter your name)    % Default value (placeholder)
//	   /Rect [100 700 300 720]  % Position
//	   /F 4                     % Print flag
//	   /Ff 2                    % Field flags (required)
//	   /DA (/Helv 12 Tf 0 g)    % Default appearance
//	>>
type TextField struct {
	// Required fields
	name  string     // Field name (unique identifier)
	rect  [4]float64 // [x, y, x+width, y+height]
	value string     // Current value

	// Optional fields
	defaultValue string // Default value (placeholder)
	maxLength    int    // Maximum character length (0 = unlimited)

	// Flags
	flags int // Field flags bitmask

	// Appearance
	fontSize    float64     // Font size for text (default: 12)
	fontName    string      // Font name (default: Helvetica)
	textColor   [3]float64  // RGB text color (default: black)
	borderColor *[3]float64 // RGB border color (nil = no border)
	fillColor   *[3]float64 // RGB fill color (nil = no fill)
}

// NewTextField creates a new text field at the specified position.
//
// Parameters:
//   - name: Unique field name (used for form data)
//   - x: Left edge position in points
//   - y: Bottom edge position in points
//   - width: Field width in points
//   - height: Field height in points
//
// Example:
//
//	field := forms.NewTextField("email", 100, 700, 200, 20)
func NewTextField(name string, x, y, width, height float64) *TextField {
	return &TextField{
		name:         name,
		rect:         [4]float64{x, y, x + width, y + height},
		value:        "",
		defaultValue: "",
		maxLength:    0,
		flags:        0,
		fontSize:     12,
		fontName:     "Helvetica",
		textColor:    [3]float64{0, 0, 0}, // Black
		borderColor:  nil,
		fillColor:    nil,
	}
}

// Name returns the field name.
func (t *TextField) Name() string {
	return t.name
}

// Type returns the PDF field type (/FT value).
// For text fields, this is always "Tx".
func (t *TextField) Type() string {
	return "Tx"
}

// Rect returns the field's bounding rectangle [x1, y1, x2, y2].
func (t *TextField) Rect() [4]float64 {
	return t.rect
}

// Flags returns the field flags bitmask.
func (t *TextField) Flags() int {
	return t.flags
}

// Value returns the field's current value.
func (t *TextField) Value() interface{} {
	return t.value
}

// DefaultValue returns the field's default value.
func (t *TextField) DefaultValue() interface{} {
	return t.defaultValue
}

// IsReadOnly returns true if the field is read-only.
func (t *TextField) IsReadOnly() bool {
	return t.flags&FlagReadOnly != 0
}

// IsRequired returns true if the field is required.
func (t *TextField) IsRequired() bool {
	return t.flags&FlagRequired != 0
}

// IsMultiline returns true if the field allows multiple lines.
func (t *TextField) IsMultiline() bool {
	return t.flags&FlagMultiline != 0
}

// IsPassword returns true if the field masks text as password.
func (t *TextField) IsPassword() bool {
	return t.flags&FlagPassword != 0
}

// SetValue sets the field's current value.
//
// Example:
//
//	field.SetValue("John Doe")
func (t *TextField) SetValue(value string) *TextField {
	t.value = value
	return t
}

// SetPlaceholder sets the field's placeholder (default value).
//
// The placeholder is shown when the field is empty.
//
// Example:
//
//	field.SetPlaceholder("Enter your name")
func (t *TextField) SetPlaceholder(placeholder string) *TextField {
	t.defaultValue = placeholder
	return t
}

// SetMaxLength sets the maximum character length for the field.
//
// Set to 0 for unlimited length.
//
// Example:
//
//	field.SetMaxLength(50)  // Max 50 characters
func (t *TextField) SetMaxLength(length int) error {
	if length < 0 {
		return errors.New("max length must be non-negative")
	}
	t.maxLength = length
	return nil
}

// MaxLength returns the maximum character length (0 = unlimited).
func (t *TextField) MaxLength() int {
	return t.maxLength
}

// SetReadOnly sets whether the field is read-only.
//
// Example:
//
//	field.SetReadOnly(true)  // Field cannot be edited
func (t *TextField) SetReadOnly(readonly bool) *TextField {
	if readonly {
		t.flags |= FlagReadOnly
	} else {
		t.flags &^= FlagReadOnly
	}
	return t
}

// SetRequired sets whether the field is required.
//
// Example:
//
//	field.SetRequired(true)  // Field must be filled
func (t *TextField) SetRequired(required bool) *TextField {
	if required {
		t.flags |= FlagRequired
	} else {
		t.flags &^= FlagRequired
	}
	return t
}

// SetMultiline sets whether the field allows multiple lines.
//
// Example:
//
//	field.SetMultiline(true)  // Text area (multiple lines)
func (t *TextField) SetMultiline(multiline bool) *TextField {
	if multiline {
		t.flags |= FlagMultiline
	} else {
		t.flags &^= FlagMultiline
	}
	return t
}

// SetPassword sets whether the field masks text as password.
//
// Example:
//
//	field.SetPassword(true)  // Show *** instead of text
func (t *TextField) SetPassword(password bool) *TextField {
	if password {
		t.flags |= FlagPassword
	} else {
		t.flags &^= FlagPassword
	}
	return t
}

// SetFontSize sets the font size for the text.
//
// Example:
//
//	field.SetFontSize(14)
func (t *TextField) SetFontSize(size float64) error {
	if size <= 0 {
		return errors.New("font size must be positive")
	}
	t.fontSize = size
	return nil
}

// FontSize returns the font size.
func (t *TextField) FontSize() float64 {
	return t.fontSize
}

// SetFontName sets the font name.
//
// Common values: "Helvetica", "Courier", "Times-Roman"
//
// Example:
//
//	field.SetFontName("Courier")
func (t *TextField) SetFontName(name string) *TextField {
	t.fontName = name
	return t
}

// FontName returns the font name.
func (t *TextField) FontName() string {
	return t.fontName
}

// SetTextColor sets the text color (RGB, 0.0-1.0 range).
//
// Example:
//
//	field.SetTextColor(0, 0, 1)  // Blue text
func (t *TextField) SetTextColor(r, g, b float64) error {
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	t.textColor = [3]float64{r, g, b}
	return nil
}

// TextColor returns the text color.
func (t *TextField) TextColor() [3]float64 {
	return t.textColor
}

// SetBorderColor sets the border color (RGB, 0.0-1.0 range).
//
// Set to nil to remove border.
//
// Example:
//
//	field.SetBorderColor(0, 0, 0)  // Black border
func (t *TextField) SetBorderColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		t.borderColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	t.borderColor = &[3]float64{*r, *g, *b}
	return nil
}

// BorderColor returns the border color (nil if no border).
func (t *TextField) BorderColor() *[3]float64 {
	return t.borderColor
}

// SetFillColor sets the background fill color (RGB, 0.0-1.0 range).
//
// Set to nil for transparent background.
//
// Example:
//
//	field.SetFillColor(1, 1, 1)  // White background
func (t *TextField) SetFillColor(r, g, b *float64) error {
	if r == nil || g == nil || b == nil {
		t.fillColor = nil
		return nil
	}
	if *r < 0 || *r > 1 || *g < 0 || *g > 1 || *b < 0 || *b > 1 {
		return errors.New("color components must be in range [0.0, 1.0]")
	}
	t.fillColor = &[3]float64{*r, *g, *b}
	return nil
}

// FillColor returns the fill color (nil if transparent).
func (t *TextField) FillColor() *[3]float64 {
	return t.fillColor
}

// Validate checks if the field configuration is valid.
//
// Returns an error if:
//   - Name is empty
//   - Rectangle has invalid dimensions
//   - Value exceeds max length
func (t *TextField) Validate() error {
	if t.name == "" {
		return errors.New("field name cannot be empty")
	}

	// Validate rectangle
	if t.rect[2] <= t.rect[0] || t.rect[3] <= t.rect[1] {
		return errors.New("invalid rectangle: width and height must be positive")
	}

	// Validate value length
	if t.maxLength > 0 && len(t.value) > t.maxLength {
		return errors.New("value exceeds maximum length")
	}

	return nil
}
