package forms_test

import (
	"testing"

	"github.com/coregx/gxpdf/creator/forms"
)

// TestNewTextField tests text field creation.
func TestNewTextField(t *testing.T) {
	field := forms.NewTextField("username", 100, 700, 200, 20)

	if field == nil {
		t.Fatal("NewTextField returned nil")
	}

	if field.Name() != "username" {
		t.Errorf("Expected name 'username', got '%s'", field.Name())
	}

	if field.Type() != "Tx" {
		t.Errorf("Expected type 'Tx', got '%s'", field.Type())
	}

	rect := field.Rect()
	expectedRect := [4]float64{100, 700, 300, 720} // x, y, x+width, y+height
	if rect != expectedRect {
		t.Errorf("Expected rect %v, got %v", expectedRect, rect)
	}

	// Default values
	if field.Value() != "" {
		t.Errorf("Expected empty value, got '%v'", field.Value())
	}

	if field.DefaultValue() != "" {
		t.Errorf("Expected empty default value, got '%v'", field.DefaultValue())
	}

	if field.Flags() != 0 {
		t.Errorf("Expected flags 0, got %d", field.Flags())
	}

	if field.FontSize() != 12 {
		t.Errorf("Expected font size 12, got %.2f", field.FontSize())
	}
}

// TestTextFieldSetValue tests setting field value.
func TestTextFieldSetValue(t *testing.T) {
	field := forms.NewTextField("name", 100, 700, 200, 20)

	field.SetValue("John Doe")

	if val := field.Value().(string); val != "John Doe" {
		t.Errorf("Expected value 'John Doe', got '%s'", val)
	}
}

// TestTextFieldSetPlaceholder tests setting placeholder.
func TestTextFieldSetPlaceholder(t *testing.T) {
	field := forms.NewTextField("email", 100, 700, 200, 20)

	field.SetPlaceholder("Enter your email")

	if val := field.DefaultValue().(string); val != "Enter your email" {
		t.Errorf("Expected placeholder 'Enter your email', got '%s'", val)
	}
}

// TestTextFieldSetRequired tests required flag.
func TestTextFieldSetRequired(t *testing.T) {
	field := forms.NewTextField("phone", 100, 700, 200, 20)

	if field.IsRequired() {
		t.Error("Field should not be required by default")
	}

	field.SetRequired(true)

	if !field.IsRequired() {
		t.Error("Field should be required after SetRequired(true)")
	}

	// Check flags
	if field.Flags()&forms.FlagRequired == 0 {
		t.Error("Required flag should be set in flags bitmask")
	}

	field.SetRequired(false)

	if field.IsRequired() {
		t.Error("Field should not be required after SetRequired(false)")
	}

	if field.Flags()&forms.FlagRequired != 0 {
		t.Error("Required flag should be cleared in flags bitmask")
	}
}

// TestTextFieldSetReadOnly tests readonly flag.
func TestTextFieldSetReadOnly(t *testing.T) {
	field := forms.NewTextField("readonly", 100, 700, 200, 20)

	if field.IsReadOnly() {
		t.Error("Field should not be readonly by default")
	}

	field.SetReadOnly(true)

	if !field.IsReadOnly() {
		t.Error("Field should be readonly after SetReadOnly(true)")
	}

	// Check flags
	if field.Flags()&forms.FlagReadOnly == 0 {
		t.Error("ReadOnly flag should be set in flags bitmask")
	}
}

// TestTextFieldSetMultiline tests multiline flag.
func TestTextFieldSetMultiline(t *testing.T) {
	field := forms.NewTextField("comment", 100, 600, 200, 80)

	if field.IsMultiline() {
		t.Error("Field should not be multiline by default")
	}

	field.SetMultiline(true)

	if !field.IsMultiline() {
		t.Error("Field should be multiline after SetMultiline(true)")
	}

	// Check flags
	if field.Flags()&forms.FlagMultiline == 0 {
		t.Error("Multiline flag should be set in flags bitmask")
	}
}

// TestTextFieldSetPassword tests password flag.
func TestTextFieldSetPassword(t *testing.T) {
	field := forms.NewTextField("password", 100, 700, 200, 20)

	if field.IsPassword() {
		t.Error("Field should not be password by default")
	}

	field.SetPassword(true)

	if !field.IsPassword() {
		t.Error("Field should be password after SetPassword(true)")
	}

	// Check flags
	if field.Flags()&forms.FlagPassword == 0 {
		t.Error("Password flag should be set in flags bitmask")
	}
}

// TestTextFieldSetMaxLength tests max length.
func TestTextFieldSetMaxLength(t *testing.T) {
	field := forms.NewTextField("zipcode", 100, 700, 100, 20)

	if field.MaxLength() != 0 {
		t.Errorf("Expected max length 0, got %d", field.MaxLength())
	}

	err := field.SetMaxLength(5)
	if err != nil {
		t.Fatalf("SetMaxLength failed: %v", err)
	}

	if field.MaxLength() != 5 {
		t.Errorf("Expected max length 5, got %d", field.MaxLength())
	}
}

// TestTextFieldSetMaxLengthNegative tests negative max length.
func TestTextFieldSetMaxLengthNegative(t *testing.T) {
	field := forms.NewTextField("test", 100, 700, 200, 20)

	err := field.SetMaxLength(-1)
	if err == nil {
		t.Error("SetMaxLength(-1) should return an error")
	}
}

// TestTextFieldSetFontSize tests font size.
func TestTextFieldSetFontSize(t *testing.T) {
	field := forms.NewTextField("test", 100, 700, 200, 20)

	err := field.SetFontSize(14)
	if err != nil {
		t.Fatalf("SetFontSize failed: %v", err)
	}

	if field.FontSize() != 14 {
		t.Errorf("Expected font size 14, got %.2f", field.FontSize())
	}
}

// TestTextFieldSetFontSizeInvalid tests invalid font size.
func TestTextFieldSetFontSizeInvalid(t *testing.T) {
	field := forms.NewTextField("test", 100, 700, 200, 20)

	err := field.SetFontSize(0)
	if err == nil {
		t.Error("SetFontSize(0) should return an error")
	}

	err = field.SetFontSize(-5)
	if err == nil {
		t.Error("SetFontSize(-5) should return an error")
	}
}

// TestTextFieldSetTextColor tests text color.
func TestTextFieldSetTextColor(t *testing.T) {
	field := forms.NewTextField("test", 100, 700, 200, 20)

	err := field.SetTextColor(1, 0, 0) // Red
	if err != nil {
		t.Fatalf("SetTextColor failed: %v", err)
	}

	color := field.TextColor()
	if color[0] != 1 || color[1] != 0 || color[2] != 0 {
		t.Errorf("Expected red color [1, 0, 0], got %v", color)
	}
}

// TestTextFieldSetTextColorInvalid tests invalid text color.
func TestTextFieldSetTextColorInvalid(t *testing.T) {
	field := forms.NewTextField("test", 100, 700, 200, 20)

	err := field.SetTextColor(1.5, 0, 0)
	if err == nil {
		t.Error("SetTextColor(1.5, 0, 0) should return an error")
	}

	err = field.SetTextColor(0, -0.1, 0)
	if err == nil {
		t.Error("SetTextColor(0, -0.1, 0) should return an error")
	}
}

// TestTextFieldValidate tests field validation.
func TestTextFieldValidate(t *testing.T) {
	// Valid field
	field := forms.NewTextField("valid", 100, 700, 200, 20)
	field.SetValue("Test")

	err := field.Validate()
	if err != nil {
		t.Errorf("Valid field should not return error: %v", err)
	}

	// Empty name
	invalidField := forms.NewTextField("", 100, 700, 200, 20)
	err = invalidField.Validate()
	if err == nil {
		t.Error("Field with empty name should return error")
	}

	// Value exceeds max length
	limitedField := forms.NewTextField("limited", 100, 700, 200, 20)
	if err := limitedField.SetMaxLength(5); err != nil {
		t.Fatalf("SetMaxLength failed: %v", err)
	}
	limitedField.SetValue("Too long value")

	err = limitedField.Validate()
	if err == nil {
		t.Error("Field with value exceeding max length should return error")
	}
}

// TestTextFieldChaining tests method chaining.
func TestTextFieldChaining(t *testing.T) {
	field := forms.NewTextField("chained", 100, 700, 200, 20)

	// Chain multiple setters
	result := field.
		SetValue("John Doe").
		SetPlaceholder("Enter name").
		SetRequired(true).
		SetReadOnly(false).
		SetMultiline(false).
		SetPassword(false).
		SetFontName("Courier")

	if result != field {
		t.Error("Methods should return the field for chaining")
	}

	if field.Value().(string) != "John Doe" {
		t.Error("Chained SetValue failed")
	}

	if !field.IsRequired() {
		t.Error("Chained SetRequired failed")
	}

	if field.FontName() != "Courier" {
		t.Error("Chained SetFontName failed")
	}
}
