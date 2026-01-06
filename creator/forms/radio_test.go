package forms_test

import (
	"testing"

	"github.com/coregx/gxpdf/creator/forms"
)

// TestNewRadioGroup tests radio group creation.
func TestNewRadioGroup(t *testing.T) {
	radioGroup := forms.NewRadioGroup("gender")

	if radioGroup == nil {
		t.Fatal("NewRadioGroup returned nil")
	}

	if radioGroup.Name() != "gender" {
		t.Errorf("Expected name 'gender', got '%s'", radioGroup.Name())
	}

	if radioGroup.Type() != "Btn" {
		t.Errorf("Expected type 'Btn', got '%s'", radioGroup.Type())
	}

	// Default values
	if radioGroup.Selected() != "" {
		t.Errorf("Expected empty selection, got '%s'", radioGroup.Selected())
	}

	if radioGroup.Value() != "" {
		t.Errorf("Expected empty value, got '%v'", radioGroup.Value())
	}

	if len(radioGroup.Options()) != 0 {
		t.Errorf("Expected 0 options, got %d", len(radioGroup.Options()))
	}

	// Radio buttons should have Radio + NoToggleToOff flags by default
	expectedFlags := forms.FlagRadio | forms.FlagNoToggleToOff
	if radioGroup.Flags() != expectedFlags {
		t.Errorf("Expected flags %d (Radio + NoToggleToOff), got %d", expectedFlags, radioGroup.Flags())
	}
}

// TestRadioGroupAddOption tests adding options.
func TestRadioGroupAddOption(t *testing.T) {
	radioGroup := forms.NewRadioGroup("payment")

	// Add first option
	radioGroup.AddOption("card", 100, 600, "Credit Card")

	if len(radioGroup.Options()) != 1 {
		t.Errorf("Expected 1 option, got %d", len(radioGroup.Options()))
	}

	opt := radioGroup.Options()[0]
	if opt.Value() != "card" {
		t.Errorf("Expected value 'card', got '%s'", opt.Value())
	}
	if opt.Label() != "Credit Card" {
		t.Errorf("Expected label 'Credit Card', got '%s'", opt.Label())
	}

	expectedRect := [4]float64{100, 600, 115, 615} // x, y, x+15, y+15 (default size)
	if opt.Rect() != expectedRect {
		t.Errorf("Expected rect %v, got %v", expectedRect, opt.Rect())
	}

	// Add second option
	radioGroup.AddOption("cash", 200, 600, "Cash")

	if len(radioGroup.Options()) != 2 {
		t.Errorf("Expected 2 options, got %d", len(radioGroup.Options()))
	}
}

// TestRadioGroupAddOptionCustomSize tests adding options with custom size.
func TestRadioGroupAddOptionCustomSize(t *testing.T) {
	radioGroup := forms.NewRadioGroup("size")

	// Add option with custom size
	radioGroup.AddOption("large", 100, 600, "Large", 20, 20)

	opt := radioGroup.Options()[0]
	expectedRect := [4]float64{100, 600, 120, 620} // x, y, x+20, y+20
	if opt.Rect() != expectedRect {
		t.Errorf("Expected rect %v, got %v", expectedRect, opt.Rect())
	}
}

// TestRadioGroupSetSelected tests setting selected option.
func TestRadioGroupSetSelected(t *testing.T) {
	radioGroup := forms.NewRadioGroup("gender")
	radioGroup.AddOption("male", 100, 600, "Male")
	radioGroup.AddOption("female", 200, 600, "Female")
	radioGroup.AddOption("other", 300, 600, "Other")

	// Initially no selection
	if radioGroup.Selected() != "" {
		t.Error("Initially no option should be selected")
	}

	// Select first option
	err := radioGroup.SetSelected("male")
	if err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	if radioGroup.Selected() != "male" {
		t.Errorf("Expected selected 'male', got '%s'", radioGroup.Selected())
	}

	if radioGroup.Value() != "male" {
		t.Errorf("Expected value 'male', got '%v'", radioGroup.Value())
	}

	// Select different option
	err = radioGroup.SetSelected("female")
	if err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	if radioGroup.Selected() != "female" {
		t.Errorf("Expected selected 'female', got '%s'", radioGroup.Selected())
	}
}

// TestRadioGroupSetSelectedInvalid tests setting invalid selection.
func TestRadioGroupSetSelectedInvalid(t *testing.T) {
	radioGroup := forms.NewRadioGroup("test")
	radioGroup.AddOption("option1", 100, 600, "Option 1")

	// Try to select non-existent option
	err := radioGroup.SetSelected("nonexistent")
	if err == nil {
		t.Error("SetSelected with invalid option should return error")
	}
}

// TestRadioGroupSetDefaultSelected tests default selection.
func TestRadioGroupSetDefaultSelected(t *testing.T) {
	radioGroup := forms.NewRadioGroup("priority")
	radioGroup.AddOption("low", 100, 600, "Low")
	radioGroup.AddOption("medium", 200, 600, "Medium")
	radioGroup.AddOption("high", 300, 600, "High")

	// Set default
	err := radioGroup.SetDefaultSelected("medium")
	if err != nil {
		t.Fatalf("SetDefaultSelected failed: %v", err)
	}

	if radioGroup.DefaultValue() != "medium" {
		t.Errorf("Expected default value 'medium', got '%v'", radioGroup.DefaultValue())
	}
}

// TestRadioGroupSetDefaultSelectedInvalid tests invalid default selection.
func TestRadioGroupSetDefaultSelectedInvalid(t *testing.T) {
	radioGroup := forms.NewRadioGroup("test")
	radioGroup.AddOption("option1", 100, 600, "Option 1")

	// Try to set non-existent default
	err := radioGroup.SetDefaultSelected("nonexistent")
	if err == nil {
		t.Error("SetDefaultSelected with invalid option should return error")
	}
}

// TestRadioGroupSetRequired tests required flag.
func TestRadioGroupSetRequired(t *testing.T) {
	radioGroup := forms.NewRadioGroup("mandatory")

	if radioGroup.IsRequired() {
		t.Error("Radio group should not be required by default")
	}

	radioGroup.SetRequired(true)

	if !radioGroup.IsRequired() {
		t.Error("Radio group should be required after SetRequired(true)")
	}

	// Check flags
	if radioGroup.Flags()&forms.FlagRequired == 0 {
		t.Error("Required flag should be set in flags bitmask")
	}

	radioGroup.SetRequired(false)

	if radioGroup.IsRequired() {
		t.Error("Radio group should not be required after SetRequired(false)")
	}
}

// TestRadioGroupSetReadOnly tests readonly flag.
func TestRadioGroupSetReadOnly(t *testing.T) {
	radioGroup := forms.NewRadioGroup("readonly")

	if radioGroup.IsReadOnly() {
		t.Error("Radio group should not be readonly by default")
	}

	radioGroup.SetReadOnly(true)

	if !radioGroup.IsReadOnly() {
		t.Error("Radio group should be readonly after SetReadOnly(true)")
	}

	// Check flags
	if radioGroup.Flags()&forms.FlagReadOnly == 0 {
		t.Error("ReadOnly flag should be set in flags bitmask")
	}
}

// TestRadioGroupSetAllowToggleOff tests toggle off flag.
func TestRadioGroupSetAllowToggleOff(t *testing.T) {
	radioGroup := forms.NewRadioGroup("toggle")

	// NoToggleToOff should be set by default
	if radioGroup.Flags()&forms.FlagNoToggleToOff == 0 {
		t.Error("NoToggleToOff flag should be set by default")
	}

	// Allow toggle off (remove NoToggleToOff flag)
	radioGroup.SetAllowToggleOff(true)

	if radioGroup.Flags()&forms.FlagNoToggleToOff != 0 {
		t.Error("NoToggleToOff flag should be cleared after SetAllowToggleOff(true)")
	}

	// Disallow toggle off (add NoToggleToOff flag)
	radioGroup.SetAllowToggleOff(false)

	if radioGroup.Flags()&forms.FlagNoToggleToOff == 0 {
		t.Error("NoToggleToOff flag should be set after SetAllowToggleOff(false)")
	}
}

// TestRadioGroupSetBorderColor tests border color.
func TestRadioGroupSetBorderColor(t *testing.T) {
	radioGroup := forms.NewRadioGroup("test")

	// Initially no border
	if radioGroup.BorderColor() != nil {
		t.Error("Border color should be nil by default")
	}

	// Set border color
	r, g, b := 0.0, 0.0, 0.0
	err := radioGroup.SetBorderColor(&r, &g, &b)
	if err != nil {
		t.Fatalf("SetBorderColor failed: %v", err)
	}

	color := radioGroup.BorderColor()
	if color == nil {
		t.Fatal("Border color should not be nil after setting")
	}
	if (*color)[0] != 0 || (*color)[1] != 0 || (*color)[2] != 0 {
		t.Errorf("Expected black border [0, 0, 0], got %v", *color)
	}

	// Remove border
	err = radioGroup.SetBorderColor()
	if err != nil {
		t.Fatalf("SetBorderColor() failed: %v", err)
	}

	if radioGroup.BorderColor() != nil {
		t.Error("Border color should be nil after removing")
	}
}

// TestRadioGroupSetBorderColorInvalid tests invalid border color.
func TestRadioGroupSetBorderColorInvalid(t *testing.T) {
	radioGroup := forms.NewRadioGroup("test")

	invalid := 1.5
	r, g := 0.0, 0.0
	err := radioGroup.SetBorderColor(&invalid, &r, &g)
	if err == nil {
		t.Error("SetBorderColor(1.5, 0, 0) should return an error")
	}

	invalid = -0.1
	err = radioGroup.SetBorderColor(&r, &invalid, &g)
	if err == nil {
		t.Error("SetBorderColor(0, -0.1, 0) should return an error")
	}
}

// TestRadioGroupSetFillColor tests fill color.
func TestRadioGroupSetFillColor(t *testing.T) {
	radioGroup := forms.NewRadioGroup("test")

	// Initially no fill
	if radioGroup.FillColor() != nil {
		t.Error("Fill color should be nil by default")
	}

	// Set fill color (white)
	r, g, b := 1.0, 1.0, 1.0
	err := radioGroup.SetFillColor(&r, &g, &b)
	if err != nil {
		t.Fatalf("SetFillColor failed: %v", err)
	}

	color := radioGroup.FillColor()
	if color == nil {
		t.Fatal("Fill color should not be nil after setting")
	}
	if (*color)[0] != 1 || (*color)[1] != 1 || (*color)[2] != 1 {
		t.Errorf("Expected white fill [1, 1, 1], got %v", *color)
	}

	// Remove fill
	err = radioGroup.SetFillColor()
	if err != nil {
		t.Fatalf("SetFillColor() failed: %v", err)
	}

	if radioGroup.FillColor() != nil {
		t.Error("Fill color should be nil after removing")
	}
}

// TestRadioGroupSetFillColorInvalid tests invalid fill color.
func TestRadioGroupSetFillColorInvalid(t *testing.T) {
	radioGroup := forms.NewRadioGroup("test")

	invalid := 1.5
	r, g := 0.0, 0.0
	err := radioGroup.SetFillColor(&invalid, &r, &g)
	if err == nil {
		t.Error("SetFillColor(1.5, 0, 0) should return an error")
	}

	invalid = -0.1
	err = radioGroup.SetFillColor(&r, &invalid, &g)
	if err == nil {
		t.Error("SetFillColor(0, -0.1, 0) should return an error")
	}
}

// TestRadioGroupValidate tests field validation.
func TestRadioGroupValidate(t *testing.T) {
	// Valid radio group
	radioGroup := forms.NewRadioGroup("valid")
	radioGroup.AddOption("opt1", 100, 600, "Option 1")
	radioGroup.AddOption("opt2", 200, 600, "Option 2")
	if err := radioGroup.SetSelected("opt1"); err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	err := radioGroup.Validate()
	if err != nil {
		t.Errorf("Valid radio group should not return error: %v", err)
	}

	// Empty name
	invalidGroup := forms.NewRadioGroup("")
	invalidGroup.AddOption("opt1", 100, 600, "Option 1")
	err = invalidGroup.Validate()
	if err == nil {
		t.Error("Radio group with empty name should return error")
	}

	// No options
	noOptionsGroup := forms.NewRadioGroup("empty")
	err = noOptionsGroup.Validate()
	if err == nil {
		t.Error("Radio group with no options should return error")
	}

	// Invalid selection - SetSelected returns error for nonexistent values
	// So we can't actually create an invalid state to test Validate() against.
	// This test is removed because SetSelected() already prevents invalid state.
}

// TestRadioGroupChaining tests method chaining.
func TestRadioGroupChaining(t *testing.T) {
	radioGroup := forms.NewRadioGroup("chained")

	// Chain multiple setters
	result := radioGroup.
		AddOption("opt1", 100, 600, "Option 1").
		AddOption("opt2", 200, 600, "Option 2").
		AddOption("opt3", 300, 600, "Option 3").
		SetRequired(true).
		SetReadOnly(false).
		SetAllowToggleOff(false)

	if result != radioGroup {
		t.Error("Methods should return the radio group for chaining")
	}

	if len(radioGroup.Options()) != 3 {
		t.Error("Chained AddOption failed")
	}

	if !radioGroup.IsRequired() {
		t.Error("Chained SetRequired failed")
	}

	// SetSelected returns error, test separately
	err := result.SetSelected("opt2")
	if err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	if radioGroup.Selected() != "opt2" {
		t.Error("SetSelected after chaining failed")
	}
}

// TestRadioGroupHasRadioFlags tests that radio button flags are set.
func TestRadioGroupHasRadioFlags(t *testing.T) {
	radioGroup := forms.NewRadioGroup("hasRadio")

	// Radio group MUST have FlagRadio set
	if radioGroup.Flags()&forms.FlagRadio == 0 {
		t.Error("Radio group MUST have FlagRadio set")
	}

	// Radio group should NOT have FlagPushbutton set
	if radioGroup.Flags()&forms.FlagPushbutton != 0 {
		t.Error("Radio group should NOT have FlagPushbutton set")
	}
}

// TestRadioGroupFormFieldInterface tests that RadioGroup implements FormField.
func TestRadioGroupFormFieldInterface(t *testing.T) {
	radioGroup := forms.NewRadioGroup("interface")
	radioGroup.AddOption("opt1", 100, 600, "Option 1")

	// Type assertion to ensure it implements FormField
	var _ forms.FormField = radioGroup

	// Test interface methods
	if radioGroup.Name() == "" {
		t.Error("Name() should return non-empty string")
	}

	if radioGroup.Type() != "Btn" {
		t.Error("Type() should return 'Btn'")
	}

	if radioGroup.Value() == nil {
		t.Error("Value() should not return nil")
	}

	if radioGroup.DefaultValue() == nil {
		t.Error("DefaultValue() should not return nil")
	}
}

// TestRadioOptionAccessors tests RadioOption accessor methods.
func TestRadioOptionAccessors(t *testing.T) {
	radioGroup := forms.NewRadioGroup("test")
	radioGroup.AddOption("testValue", 100, 600, "Test Label", 20, 20)

	opt := radioGroup.Options()[0]

	if opt.Value() != "testValue" {
		t.Errorf("Expected value 'testValue', got '%s'", opt.Value())
	}

	if opt.Label() != "Test Label" {
		t.Errorf("Expected label 'Test Label', got '%s'", opt.Label())
	}

	expectedRect := [4]float64{100, 600, 120, 620}
	if opt.Rect() != expectedRect {
		t.Errorf("Expected rect %v, got %v", expectedRect, opt.Rect())
	}
}

// TestRadioGroupMultipleSelections tests that only one option can be selected.
func TestRadioGroupMultipleSelections(t *testing.T) {
	radioGroup := forms.NewRadioGroup("single")
	radioGroup.AddOption("opt1", 100, 600, "Option 1")
	radioGroup.AddOption("opt2", 200, 600, "Option 2")
	radioGroup.AddOption("opt3", 300, 600, "Option 3")

	// Select first option
	if err := radioGroup.SetSelected("opt1"); err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	if radioGroup.Selected() != "opt1" {
		t.Error("First option should be selected")
	}

	// Select second option (should replace first)
	if err := radioGroup.SetSelected("opt2"); err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	if radioGroup.Selected() != "opt2" {
		t.Error("Second option should be selected")
	}

	// Only one option should be selected
	if radioGroup.Value() != "opt2" {
		t.Error("Only second option should be selected")
	}
}
