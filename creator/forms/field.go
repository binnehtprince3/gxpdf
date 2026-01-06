// Package forms provides interactive form field support for PDF documents.
//
// This package implements AcroForm (Interactive Form) support, allowing
// creation of fillable PDF forms with text fields, checkboxes, radio buttons,
// and other interactive elements.
//
// # Thread Safety
//
// Form field instances are NOT safe for concurrent use. Each goroutine
// should create and manipulate its own field instances. Fields should be
// created and configured before adding to a page.
//
// # Example
//
//	c := creator.New()
//	page, _ := c.NewPage()
//
//	// Text field
//	nameField := forms.NewTextField("name", 100, 700, 200, 20)
//	nameField.SetValue("John Doe")
//	page.AddField(nameField)
//
//	c.WriteToFile("form.pdf")
package forms

// FormField represents an interactive form field in a PDF.
//
// All form fields share common properties like name, value, flags,
// and position. Specific field types (TextField, CheckBox, etc.)
// implement this interface.
//
// This follows the DDD pattern where FormField is a polymorphic entity
// within the Form aggregate.
type FormField interface {
	// Name returns the field name (used for form data export).
	Name() string

	// Type returns the PDF field type (/FT value).
	// - "Tx" = Text field
	// - "Btn" = Button (checkbox, radio, pushbutton)
	// - "Ch" = Choice (list box, combo box)
	// - "Sig" = Signature field
	Type() string

	// Rect returns the field's position and size [x, y, width, height].
	Rect() [4]float64

	// Flags returns the field flags (Ff) as a bitmask.
	// Common flags:
	// - Bit 1 (1): ReadOnly
	// - Bit 2 (2): Required
	// - Bit 13 (4096): Multiline (for text fields)
	// - Bit 14 (8192): Password (for text fields)
	Flags() int

	// Value returns the field's current value.
	// For text fields, this is a string.
	// For buttons, this is the button state.
	// For choice fields, this is the selected option(s).
	Value() interface{}

	// DefaultValue returns the field's default value.
	// This is used when the form is reset.
	DefaultValue() interface{}

	// IsReadOnly returns true if the field is read-only.
	IsReadOnly() bool

	// IsRequired returns true if the field is required.
	IsRequired() bool
}

// FieldFlags defines common field flags as constants.
const (
	// FlagReadOnly makes the field read-only (bit 1).
	FlagReadOnly = 1 << 0 // 1

	// FlagRequired marks the field as required (bit 2).
	FlagRequired = 1 << 1 // 2

	// FlagNoExport excludes the field from form data export (bit 3).
	FlagNoExport = 1 << 2 // 4
)

// TextFieldFlags defines text field specific flags.
const (
	// FlagMultiline allows multiple lines in the text field (bit 13).
	FlagMultiline = 1 << 12 // 4096

	// FlagPassword masks text entry as password (bit 14).
	FlagPassword = 1 << 13 // 8192

	// FlagFileSelect treats the field as a file selection (bit 21).
	FlagFileSelect = 1 << 20 // 1048576

	// FlagDoNotSpellCheck disables spell checking (bit 23).
	FlagDoNotSpellCheck = 1 << 22 // 4194304

	// FlagDoNotScroll prevents scrolling in the field (bit 24).
	FlagDoNotScroll = 1 << 23 // 8388608

	// FlagComb creates a comb field (fixed character positions) (bit 25).
	FlagComb = 1 << 24 // 16777216

	// FlagRichText enables rich text formatting (bit 26).
	FlagRichText = 1 << 25 // 33554432
)

// ButtonFieldFlags defines button field specific flags.
const (
	// FlagNoToggleToOff prevents checkboxes/radio from toggling off (bit 15).
	FlagNoToggleToOff = 1 << 14 // 16384

	// FlagRadio makes the button a radio button (bit 16).
	FlagRadio = 1 << 15 // 32768

	// FlagPushbutton makes the button a pushbutton (bit 17).
	FlagPushbutton = 1 << 16 // 65536

	// FlagRadiosInUnison makes all radio buttons with same value toggle together (bit 26).
	FlagRadiosInUnison = 1 << 25 // 33554432
)

// ChoiceFieldFlags defines choice field specific flags.
const (
	// FlagCombo makes the choice field a combo box (editable) (bit 18).
	FlagCombo = 1 << 17 // 131072

	// FlagEdit allows editing in combo box (bit 19).
	FlagEdit = 1 << 18 // 262144

	// FlagSort sorts choice options (bit 20).
	FlagSort = 1 << 19 // 524288

	// FlagMultiSelect allows multiple selections in list (bit 22).
	FlagMultiSelect = 1 << 21 // 2097152

	// FlagCommitOnSelChange commits value on selection change (bit 27).
	FlagCommitOnSelChange = 1 << 26 // 67108864
)
