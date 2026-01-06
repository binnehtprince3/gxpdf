package document

import (
	"errors"
)

// FormField represents an interactive form field widget annotation.
//
// Form fields are special annotations that allow user input.
// They are part of the AcroForm (Interactive Form) system in PDF.
//
// This is a value object that represents the PDF Widget annotation
// combined with the form field dictionary.
//
// PDF Structure:
//
//	<< /Type /Annot
//	   /Subtype /Widget
//	   /FT /Tx                  % Field Type (Tx=Text, Btn=Button, Ch=Choice, Sig=Signature)
//	   /T (fieldName)           % Field name
//	   /V (value)               % Field value
//	   /DV (defaultValue)       % Default value
//	   /Rect [x1 y1 x2 y2]      % Position
//	   /F 4                     % Annotation flags (4 = Print)
//	   /Ff 0                    % Field flags
//	   /DA (/Helv 12 Tf 0 g)    % Default appearance
//	   /MaxLen 100              % Max length (text fields only)
//	>>
type FormField struct {
	// Field identity
	fieldType     string // PDF field type: Tx, Btn, Ch, Sig
	name          string // Field name (unique identifier)
	alternateText string // Alternate text for accessibility

	// Field value
	value        string // Current value
	defaultValue string // Default value

	// Position and appearance
	rect        [4]float64  // [x1, y1, x2, y2]
	flags       int         // Field flags (Ff)
	annotFlags  int         // Annotation flags (F)
	appearance  string      // Default appearance string (/DA)
	borderColor *[3]float64 // Border color RGB
	fillColor   *[3]float64 // Fill color RGB

	// Text field specific
	maxLength int // Maximum text length (0 = unlimited)

	// Choice field specific
	options []string // Choice options
}

// NewFormField creates a new form field.
//
// Parameters:
//   - fieldType: PDF field type (Tx, Btn, Ch, Sig)
//   - name: Unique field name
//   - rect: Position and size [x1, y1, x2, y2]
//
// Example:
//
//	field := NewFormField("Tx", "username", [4]float64{100, 700, 300, 720})
func NewFormField(fieldType, name string, rect [4]float64) *FormField {
	return &FormField{
		fieldType:    fieldType,
		name:         name,
		rect:         rect,
		value:        "",
		defaultValue: "",
		flags:        0,
		annotFlags:   4, // Print flag (bit 3)
		appearance:   "/Helv 12 Tf 0 g",
		borderColor:  nil,
		fillColor:    nil,
		maxLength:    0,
		options:      nil,
	}
}

// FieldType returns the PDF field type.
func (f *FormField) FieldType() string {
	return f.fieldType
}

// Name returns the field name.
func (f *FormField) Name() string {
	return f.name
}

// SetAlternateText sets the alternate text for accessibility.
func (f *FormField) SetAlternateText(text string) {
	f.alternateText = text
}

// AlternateText returns the alternate text.
func (f *FormField) AlternateText() string {
	return f.alternateText
}

// SetValue sets the field value.
func (f *FormField) SetValue(value string) {
	f.value = value
}

// Value returns the field value.
func (f *FormField) Value() string {
	return f.value
}

// SetDefaultValue sets the default value.
func (f *FormField) SetDefaultValue(value string) {
	f.defaultValue = value
}

// DefaultValue returns the default value.
func (f *FormField) DefaultValue() string {
	return f.defaultValue
}

// Rect returns the field rectangle [x1, y1, x2, y2].
func (f *FormField) Rect() [4]float64 {
	return f.rect
}

// SetFlags sets the field flags (Ff).
func (f *FormField) SetFlags(flags int) {
	f.flags = flags
}

// Flags returns the field flags.
func (f *FormField) Flags() int {
	return f.flags
}

// SetAnnotationFlags sets the annotation flags (F).
func (f *FormField) SetAnnotationFlags(flags int) {
	f.annotFlags = flags
}

// AnnotationFlags returns the annotation flags.
func (f *FormField) AnnotationFlags() int {
	return f.annotFlags
}

// SetAppearance sets the default appearance string (/DA).
//
// Example appearance strings:
//   - "/Helv 12 Tf 0 g" - Helvetica 12pt, black
//   - "/Courier 10 Tf 1 0 0 rg" - Courier 10pt, red
func (f *FormField) SetAppearance(appearance string) {
	f.appearance = appearance
}

// Appearance returns the default appearance string.
func (f *FormField) Appearance() string {
	return f.appearance
}

// SetBorderColor sets the border color.
func (f *FormField) SetBorderColor(r, g, b float64) {
	f.borderColor = &[3]float64{r, g, b}
}

// BorderColor returns the border color (nil if no border).
func (f *FormField) BorderColor() *[3]float64 {
	return f.borderColor
}

// SetFillColor sets the fill color.
func (f *FormField) SetFillColor(r, g, b float64) {
	f.fillColor = &[3]float64{r, g, b}
}

// FillColor returns the fill color (nil if transparent).
func (f *FormField) FillColor() *[3]float64 {
	return f.fillColor
}

// SetMaxLength sets the maximum text length (text fields only).
func (f *FormField) SetMaxLength(length int) {
	f.maxLength = length
}

// MaxLength returns the maximum text length.
func (f *FormField) MaxLength() int {
	return f.maxLength
}

// SetOptions sets the choice options (choice fields only).
func (f *FormField) SetOptions(options []string) {
	f.options = make([]string, len(options))
	copy(f.options, options)
}

// Options returns the choice options.
func (f *FormField) Options() []string {
	if f.options == nil {
		return nil
	}
	result := make([]string, len(f.options))
	copy(result, f.options)
	return result
}

// Validate checks if the form field is valid.
//
// Returns an error if:
//   - Field type is invalid
//   - Name is empty
//   - Rectangle has invalid dimensions
//   - Value exceeds max length (for text fields)
func (f *FormField) Validate() error {
	if err := validateFieldType(f.fieldType); err != nil {
		return err
	}

	if f.name == "" {
		return errors.New("field name cannot be empty")
	}

	if err := validateRectangle(f.rect); err != nil {
		return err
	}

	if err := validateTextFieldLength(f.fieldType, f.value, f.maxLength); err != nil {
		return err
	}

	if err := validateColor(f.borderColor, "border"); err != nil {
		return err
	}

	if err := validateColor(f.fillColor, "fill"); err != nil {
		return err
	}

	return nil
}

// validateFieldType validates the PDF field type.
func validateFieldType(fieldType string) error {
	switch fieldType {
	case "Tx", "Btn", "Ch", "Sig":
		return nil
	default:
		return errors.New("invalid field type: must be Tx, Btn, Ch, or Sig")
	}
}

// validateRectangle validates the field rectangle dimensions.
func validateRectangle(rect [4]float64) error {
	if rect[2] <= rect[0] || rect[3] <= rect[1] {
		return errors.New("invalid rectangle: width and height must be positive")
	}
	return nil
}

// validateTextFieldLength validates text field value length.
func validateTextFieldLength(fieldType, value string, maxLength int) error {
	if fieldType == "Tx" && maxLength > 0 {
		if len(value) > maxLength {
			return errors.New("value exceeds maximum length")
		}
	}
	return nil
}

// validateColor validates color component values.
func validateColor(color *[3]float64, colorType string) error {
	if color != nil {
		for _, c := range *color {
			if c < 0 || c > 1 {
				return errors.New(colorType + " color components must be in range [0.0, 1.0]")
			}
		}
	}
	return nil
}
