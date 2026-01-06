package forms_test

import (
	"testing"

	"github.com/coregx/gxpdf/creator/forms"
)

// TestNewCheckbox tests checkbox creation.
func TestNewCheckbox(t *testing.T) {
	checkbox := forms.NewCheckbox("agree", 100, 650, 15, 15)

	if checkbox == nil {
		t.Fatal("NewCheckbox returned nil")
	}

	if checkbox.Name() != "agree" {
		t.Errorf("Expected name 'agree', got '%s'", checkbox.Name())
	}

	if checkbox.Type() != "Btn" {
		t.Errorf("Expected type 'Btn', got '%s'", checkbox.Type())
	}

	rect := checkbox.Rect()
	expectedRect := [4]float64{100, 650, 115, 665} // x, y, x+width, y+height
	if rect != expectedRect {
		t.Errorf("Expected rect %v, got %v", expectedRect, rect)
	}

	// Default values
	if checkbox.IsChecked() {
		t.Error("Checkbox should be unchecked by default")
	}

	if checkbox.Value() != "Off" {
		t.Errorf("Expected value 'Off', got '%v'", checkbox.Value())
	}

	if checkbox.DefaultValue() != "Off" {
		t.Errorf("Expected default value 'Off', got '%v'", checkbox.DefaultValue())
	}

	if checkbox.Flags() != 0 {
		t.Errorf("Expected flags 0, got %d", checkbox.Flags())
	}

	if checkbox.Label() != "" {
		t.Errorf("Expected empty label, got '%s'", checkbox.Label())
	}
}

// TestCheckboxSetChecked tests setting checked state.
func TestCheckboxSetChecked(t *testing.T) {
	checkbox := forms.NewCheckbox("test", 100, 650, 15, 15)

	// Initially unchecked
	if checkbox.IsChecked() {
		t.Error("Checkbox should be unchecked initially")
	}

	// Set checked
	checkbox.SetChecked(true)

	if !checkbox.IsChecked() {
		t.Error("Checkbox should be checked after SetChecked(true)")
	}

	if checkbox.Value() != "Yes" {
		t.Errorf("Expected value 'Yes' when checked, got '%v'", checkbox.Value())
	}

	// Set unchecked
	checkbox.SetChecked(false)

	if checkbox.IsChecked() {
		t.Error("Checkbox should be unchecked after SetChecked(false)")
	}

	if checkbox.Value() != "Off" {
		t.Errorf("Expected value 'Off' when unchecked, got '%v'", checkbox.Value())
	}
}

// TestCheckboxSetLabel tests setting label.
func TestCheckboxSetLabel(t *testing.T) {
	checkbox := forms.NewCheckbox("terms", 100, 650, 15, 15)

	checkbox.SetLabel("I agree to the terms and conditions")

	if checkbox.Label() != "I agree to the terms and conditions" {
		t.Errorf("Expected label 'I agree to the terms and conditions', got '%s'", checkbox.Label())
	}
}

// TestCheckboxSetDefaultChecked tests default checked state.
func TestCheckboxSetDefaultChecked(t *testing.T) {
	checkbox := forms.NewCheckbox("subscribe", 100, 650, 15, 15)

	// Default is unchecked
	if checkbox.DefaultValue() != "Off" {
		t.Errorf("Expected default value 'Off', got '%v'", checkbox.DefaultValue())
	}

	// Set default to checked
	checkbox.SetDefaultChecked(true)

	if checkbox.DefaultValue() != "Yes" {
		t.Errorf("Expected default value 'Yes', got '%v'", checkbox.DefaultValue())
	}

	// Set default to unchecked
	checkbox.SetDefaultChecked(false)

	if checkbox.DefaultValue() != "Off" {
		t.Errorf("Expected default value 'Off', got '%v'", checkbox.DefaultValue())
	}
}

// TestCheckboxSetRequired tests required flag.
func TestCheckboxSetRequired(t *testing.T) {
	checkbox := forms.NewCheckbox("mandatory", 100, 650, 15, 15)

	if checkbox.IsRequired() {
		t.Error("Checkbox should not be required by default")
	}

	checkbox.SetRequired(true)

	if !checkbox.IsRequired() {
		t.Error("Checkbox should be required after SetRequired(true)")
	}

	// Check flags
	if checkbox.Flags()&forms.FlagRequired == 0 {
		t.Error("Required flag should be set in flags bitmask")
	}

	checkbox.SetRequired(false)

	if checkbox.IsRequired() {
		t.Error("Checkbox should not be required after SetRequired(false)")
	}

	if checkbox.Flags()&forms.FlagRequired != 0 {
		t.Error("Required flag should be cleared in flags bitmask")
	}
}

// TestCheckboxSetReadOnly tests readonly flag.
func TestCheckboxSetReadOnly(t *testing.T) {
	checkbox := forms.NewCheckbox("readonly", 100, 650, 15, 15)

	if checkbox.IsReadOnly() {
		t.Error("Checkbox should not be readonly by default")
	}

	checkbox.SetReadOnly(true)

	if !checkbox.IsReadOnly() {
		t.Error("Checkbox should be readonly after SetReadOnly(true)")
	}

	// Check flags
	if checkbox.Flags()&forms.FlagReadOnly == 0 {
		t.Error("ReadOnly flag should be set in flags bitmask")
	}
}

// TestCheckboxSetBorderColor tests border color.
func TestCheckboxSetBorderColor(t *testing.T) {
	checkbox := forms.NewCheckbox("test", 100, 650, 15, 15)

	// Initially no border
	if checkbox.BorderColor() != nil {
		t.Error("Border color should be nil by default")
	}

	// Set border color
	r, g, b := 0.0, 0.0, 0.0
	err := checkbox.SetBorderColor(&r, &g, &b)
	if err != nil {
		t.Fatalf("SetBorderColor failed: %v", err)
	}

	color := checkbox.BorderColor()
	if color == nil {
		t.Fatal("Border color should not be nil after setting")
	}
	if (*color)[0] != 0 || (*color)[1] != 0 || (*color)[2] != 0 {
		t.Errorf("Expected black border [0, 0, 0], got %v", *color)
	}

	// Remove border
	err = checkbox.SetBorderColor(nil, nil, nil)
	if err != nil {
		t.Fatalf("SetBorderColor(nil) failed: %v", err)
	}

	if checkbox.BorderColor() != nil {
		t.Error("Border color should be nil after removing")
	}
}

// TestCheckboxSetBorderColorInvalid tests invalid border color.
func TestCheckboxSetBorderColorInvalid(t *testing.T) {
	checkbox := forms.NewCheckbox("test", 100, 650, 15, 15)

	invalid := 1.5
	r, g := 0.0, 0.0
	err := checkbox.SetBorderColor(&invalid, &r, &g)
	if err == nil {
		t.Error("SetBorderColor(1.5, 0, 0) should return an error")
	}

	invalid = -0.1
	err = checkbox.SetBorderColor(&r, &invalid, &g)
	if err == nil {
		t.Error("SetBorderColor(0, -0.1, 0) should return an error")
	}
}

// TestCheckboxSetFillColor tests fill color.
func TestCheckboxSetFillColor(t *testing.T) {
	checkbox := forms.NewCheckbox("test", 100, 650, 15, 15)

	// Initially no fill
	if checkbox.FillColor() != nil {
		t.Error("Fill color should be nil by default")
	}

	// Set fill color (white)
	r, g, b := 1.0, 1.0, 1.0
	err := checkbox.SetFillColor(&r, &g, &b)
	if err != nil {
		t.Fatalf("SetFillColor failed: %v", err)
	}

	color := checkbox.FillColor()
	if color == nil {
		t.Fatal("Fill color should not be nil after setting")
	}
	if (*color)[0] != 1 || (*color)[1] != 1 || (*color)[2] != 1 {
		t.Errorf("Expected white fill [1, 1, 1], got %v", *color)
	}

	// Remove fill
	err = checkbox.SetFillColor(nil, nil, nil)
	if err != nil {
		t.Fatalf("SetFillColor(nil) failed: %v", err)
	}

	if checkbox.FillColor() != nil {
		t.Error("Fill color should be nil after removing")
	}
}

// TestCheckboxSetFillColorInvalid tests invalid fill color.
func TestCheckboxSetFillColorInvalid(t *testing.T) {
	checkbox := forms.NewCheckbox("test", 100, 650, 15, 15)

	invalid := 1.5
	r, g := 0.0, 0.0
	err := checkbox.SetFillColor(&invalid, &r, &g)
	if err == nil {
		t.Error("SetFillColor(1.5, 0, 0) should return an error")
	}

	invalid = -0.1
	err = checkbox.SetFillColor(&r, &invalid, &g)
	if err == nil {
		t.Error("SetFillColor(0, -0.1, 0) should return an error")
	}
}

// TestCheckboxValidate tests field validation.
func TestCheckboxValidate(t *testing.T) {
	// Valid checkbox
	checkbox := forms.NewCheckbox("valid", 100, 650, 15, 15)
	checkbox.SetChecked(true)

	err := checkbox.Validate()
	if err != nil {
		t.Errorf("Valid checkbox should not return error: %v", err)
	}

	// Empty name
	invalidCheckbox := forms.NewCheckbox("", 100, 650, 15, 15)
	err = invalidCheckbox.Validate()
	if err == nil {
		t.Error("Checkbox with empty name should return error")
	}

	// Invalid rectangle (zero width)
	invalidRect := forms.NewCheckbox("invalid", 100, 650, 0, 15)
	err = invalidRect.Validate()
	if err == nil {
		t.Error("Checkbox with zero width should return error")
	}
}

// TestCheckboxChaining tests method chaining.
func TestCheckboxChaining(t *testing.T) {
	checkbox := forms.NewCheckbox("chained", 100, 650, 15, 15)

	// Chain multiple setters
	r, g, b := 0.0, 0.0, 0.0
	result := checkbox.
		SetChecked(true).
		SetLabel("I agree").
		SetDefaultChecked(false).
		SetRequired(true).
		SetReadOnly(false)

	// SetBorderColor returns error, so test separately
	err := result.SetBorderColor(&r, &g, &b)
	if err != nil {
		t.Fatalf("SetBorderColor failed: %v", err)
	}

	if result != checkbox {
		t.Error("Methods should return the checkbox for chaining")
	}

	if !checkbox.IsChecked() {
		t.Error("Chained SetChecked failed")
	}

	if checkbox.Label() != "I agree" {
		t.Error("Chained SetLabel failed")
	}

	if !checkbox.IsRequired() {
		t.Error("Chained SetRequired failed")
	}

	if checkbox.BorderColor() == nil {
		t.Error("Chained SetBorderColor failed")
	}
}

// TestCheckboxNoRadioFlags tests that radio button flags are NOT set.
func TestCheckboxNoRadioFlags(t *testing.T) {
	checkbox := forms.NewCheckbox("notRadio", 100, 650, 15, 15)

	// Checkbox should NOT have radio button flags
	if checkbox.Flags()&forms.FlagRadio != 0 {
		t.Error("Checkbox should NOT have FlagRadio set")
	}

	if checkbox.Flags()&forms.FlagNoToggleToOff != 0 {
		t.Error("Checkbox should NOT have FlagNoToggleToOff set")
	}

	if checkbox.Flags()&forms.FlagPushbutton != 0 {
		t.Error("Checkbox should NOT have FlagPushbutton set")
	}
}

// TestCheckboxFormFieldInterface tests that Checkbox implements FormField.
func TestCheckboxFormFieldInterface(t *testing.T) {
	checkbox := forms.NewCheckbox("interface", 100, 650, 15, 15)

	// Type assertion to ensure it implements FormField
	var _ forms.FormField = checkbox

	// Test interface methods
	if checkbox.Name() == "" {
		t.Error("Name() should return non-empty string")
	}

	if checkbox.Type() != "Btn" {
		t.Error("Type() should return 'Btn'")
	}

	if checkbox.Value() == nil {
		t.Error("Value() should not return nil")
	}

	if checkbox.DefaultValue() == nil {
		t.Error("DefaultValue() should not return nil")
	}
}
