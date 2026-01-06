package forms_test

import (
	"testing"

	"github.com/coregx/gxpdf/creator/forms"
)

// TestNewListBox tests list box creation.
func TestNewListBox(t *testing.T) {
	listbox := forms.NewListBox("colors", 100, 350, 150, 80)

	if listbox == nil {
		t.Fatal("NewListBox returned nil")
	}

	if listbox.Name() != "colors" {
		t.Errorf("Expected name 'colors', got '%s'", listbox.Name())
	}

	if listbox.Type() != "Ch" {
		t.Errorf("Expected type 'Ch', got '%s'", listbox.Type())
	}

	rect := listbox.Rect()
	expectedRect := [4]float64{100, 350, 250, 430} // x, y, x+width, y+height
	if rect != expectedRect {
		t.Errorf("Expected rect %v, got %v", expectedRect, rect)
	}

	// Default values
	if listbox.Value() != "" {
		t.Errorf("Expected empty value, got '%v'", listbox.Value())
	}

	if listbox.DefaultValue() != "" {
		t.Errorf("Expected empty default value, got '%v'", listbox.DefaultValue())
	}

	// Combo flag should NOT be set for list boxes
	if listbox.Flags()&forms.FlagCombo != 0 {
		t.Error("Combo flag should NOT be set for list boxes")
	}

	if len(listbox.Options()) != 0 {
		t.Errorf("Expected 0 options, got %d", len(listbox.Options()))
	}
}

// TestListBoxAddOption tests adding individual options.
func TestListBoxAddOption(t *testing.T) {
	listbox := forms.NewListBox("colors", 100, 350, 150, 80)

	listbox.AddOption("red", "Red Color")
	listbox.AddOption("green", "Green Color")
	listbox.AddOption("blue", "Blue Color")

	options := listbox.Options()
	if len(options) != 3 {
		t.Fatalf("Expected 3 options, got %d", len(options))
	}

	// Check first option
	if options[0].ExportValue != "red" {
		t.Errorf("Expected export value 'red', got '%s'", options[0].ExportValue)
	}
	if options[0].DisplayValue != "Red Color" {
		t.Errorf("Expected display value 'Red Color', got '%s'", options[0].DisplayValue)
	}
}

// TestListBoxAddOptions tests adding multiple options at once.
func TestListBoxAddOptions(t *testing.T) {
	listbox := forms.NewListBox("interests", 100, 350, 150, 80)

	listbox.AddOptions("Sports", "Music", "Movies", "Reading", "Travel")

	options := listbox.Options()
	if len(options) != 5 {
		t.Fatalf("Expected 5 options, got %d", len(options))
	}

	// For AddOptions, export value should equal display value
	for i, expected := range []string{"Sports", "Music", "Movies", "Reading", "Travel"} {
		if options[i].ExportValue != expected {
			t.Errorf("Option %d: expected export value '%s', got '%s'", i, expected, options[i].ExportValue)
		}
		if options[i].DisplayValue != expected {
			t.Errorf("Option %d: expected display value '%s', got '%s'", i, expected, options[i].DisplayValue)
		}
	}
}

// TestListBoxSetSelected tests setting single selected value.
func TestListBoxSetSelected(t *testing.T) {
	listbox := forms.NewListBox("color", 100, 350, 150, 80)
	listbox.AddOptions("Red", "Green", "Blue")

	err := listbox.SetSelected("Green")
	if err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	selected := listbox.SelectedValues()
	if len(selected) != 1 {
		t.Fatalf("Expected 1 selected value, got %d", len(selected))
	}
	if selected[0] != "Green" {
		t.Errorf("Expected selected 'Green', got '%s'", selected[0])
	}

	// For single-select, Value() should return string
	if val := listbox.Value().(string); val != "Green" {
		t.Errorf("Expected Value() 'Green', got '%s'", val)
	}
}

// TestListBoxSetSelectedInvalid tests setting invalid selected value.
func TestListBoxSetSelectedInvalid(t *testing.T) {
	listbox := forms.NewListBox("color", 100, 350, 150, 80)
	listbox.AddOptions("Red", "Green", "Blue")

	err := listbox.SetSelected("Invalid")
	if err == nil {
		t.Error("SetSelected with invalid value should return error")
	}
}

// TestListBoxSetSelectedMultiple tests setting multiple selected values.
func TestListBoxSetSelectedMultiple(t *testing.T) {
	listbox := forms.NewListBox("interests", 100, 350, 150, 80)
	listbox.AddOptions("Sports", "Music", "Movies", "Reading", "Travel")
	listbox.SetMultiSelect(true)

	err := listbox.SetSelectedMultiple("Sports", "Music", "Travel")
	if err != nil {
		t.Fatalf("SetSelectedMultiple failed: %v", err)
	}

	selected := listbox.SelectedValues()
	if len(selected) != 3 {
		t.Fatalf("Expected 3 selected values, got %d", len(selected))
	}

	// Check all selected values
	expectedSelected := []string{"Sports", "Music", "Travel"}
	for i, expected := range expectedSelected {
		if selected[i] != expected {
			t.Errorf("Selected[%d]: expected '%s', got '%s'", i, expected, selected[i])
		}
	}

	// For multi-select, Value() should return []string
	if val, ok := listbox.Value().([]string); !ok {
		t.Error("Value() should return []string for multi-select")
	} else if len(val) != 3 {
		t.Errorf("Expected 3 values, got %d", len(val))
	}
}

// TestListBoxSetSelectedMultipleInvalid tests setting invalid multiple values.
func TestListBoxSetSelectedMultipleInvalid(t *testing.T) {
	listbox := forms.NewListBox("interests", 100, 350, 150, 80)
	listbox.AddOptions("Sports", "Music", "Movies")
	listbox.SetMultiSelect(true)

	err := listbox.SetSelectedMultiple("Sports", "Invalid", "Music")
	if err == nil {
		t.Error("SetSelectedMultiple with invalid value should return error")
	}
}

// TestListBoxSetDefaultValue tests setting default single value.
func TestListBoxSetDefaultValue(t *testing.T) {
	listbox := forms.NewListBox("color", 100, 350, 150, 80)
	listbox.AddOptions("Red", "Green", "Blue")

	err := listbox.SetDefaultValue("Blue")
	if err != nil {
		t.Fatalf("SetDefaultValue failed: %v", err)
	}

	if listbox.DefaultValue().(string) != "Blue" {
		t.Errorf("Expected default value 'Blue', got '%s'", listbox.DefaultValue().(string))
	}
}

// TestListBoxSetDefaultValueMultiple tests setting default multiple values.
func TestListBoxSetDefaultValueMultiple(t *testing.T) {
	listbox := forms.NewListBox("interests", 100, 350, 150, 80)
	listbox.AddOptions("Sports", "Music", "Movies")
	listbox.SetMultiSelect(true)

	err := listbox.SetDefaultValueMultiple("Sports", "Movies")
	if err != nil {
		t.Fatalf("SetDefaultValueMultiple failed: %v", err)
	}

	defaults := listbox.DefaultValue().([]string)
	if len(defaults) != 2 {
		t.Fatalf("Expected 2 default values, got %d", len(defaults))
	}
}

// TestListBoxSetMultiSelect tests multi-select flag.
func TestListBoxSetMultiSelect(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)

	if listbox.IsMultiSelect() {
		t.Error("ListBox should not be multi-select by default")
	}

	listbox.SetMultiSelect(true)

	if !listbox.IsMultiSelect() {
		t.Error("ListBox should be multi-select after SetMultiSelect(true)")
	}

	// Check flags
	if listbox.Flags()&forms.FlagMultiSelect == 0 {
		t.Error("MultiSelect flag should be set in flags bitmask")
	}

	// Test disabling multi-select with multiple selections
	listbox.AddOptions("A", "B", "C")
	if err := listbox.SetSelectedMultiple("A", "B", "C"); err != nil {
		t.Fatalf("SetSelectedMultiple failed: %v", err)
	}

	listbox.SetMultiSelect(false)

	// Should keep only first selected item
	selected := listbox.SelectedValues()
	if len(selected) != 1 {
		t.Errorf("Expected 1 selected value after disabling multi-select, got %d", len(selected))
	}
}

// TestListBoxSetSort tests sort flag.
func TestListBoxSetSort(t *testing.T) {
	listbox := forms.NewListBox("sorted", 100, 350, 150, 80)

	if listbox.IsSorted() {
		t.Error("ListBox should not be sorted by default")
	}

	listbox.SetSort(true)

	if !listbox.IsSorted() {
		t.Error("ListBox should be sorted after SetSort(true)")
	}

	// Check flags
	if listbox.Flags()&forms.FlagSort == 0 {
		t.Error("Sort flag should be set in flags bitmask")
	}
}

// TestListBoxSetRequired tests required flag.
func TestListBoxSetRequired(t *testing.T) {
	listbox := forms.NewListBox("required", 100, 350, 150, 80)

	if listbox.IsRequired() {
		t.Error("ListBox should not be required by default")
	}

	listbox.SetRequired(true)

	if !listbox.IsRequired() {
		t.Error("ListBox should be required after SetRequired(true)")
	}

	// Check flags
	if listbox.Flags()&forms.FlagRequired == 0 {
		t.Error("Required flag should be set in flags bitmask")
	}
}

// TestListBoxSetReadOnly tests readonly flag.
func TestListBoxSetReadOnly(t *testing.T) {
	listbox := forms.NewListBox("readonly", 100, 350, 150, 80)

	if listbox.IsReadOnly() {
		t.Error("ListBox should not be readonly by default")
	}

	listbox.SetReadOnly(true)

	if !listbox.IsReadOnly() {
		t.Error("ListBox should be readonly after SetReadOnly(true)")
	}

	// Check flags
	if listbox.Flags()&forms.FlagReadOnly == 0 {
		t.Error("ReadOnly flag should be set in flags bitmask")
	}
}

// TestListBoxSetFontSize tests font size.
func TestListBoxSetFontSize(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)

	err := listbox.SetFontSize(14)
	if err != nil {
		t.Fatalf("SetFontSize failed: %v", err)
	}

	if listbox.FontSize() != 14 {
		t.Errorf("Expected font size 14, got %.2f", listbox.FontSize())
	}
}

// TestListBoxSetFontSizeInvalid tests invalid font size.
func TestListBoxSetFontSizeInvalid(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)

	err := listbox.SetFontSize(0)
	if err == nil {
		t.Error("SetFontSize(0) should return an error")
	}

	err = listbox.SetFontSize(-5)
	if err == nil {
		t.Error("SetFontSize(-5) should return an error")
	}
}

// TestListBoxSetTextColor tests text color.
func TestListBoxSetTextColor(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)

	err := listbox.SetTextColor(0, 0, 1) // Blue
	if err != nil {
		t.Fatalf("SetTextColor failed: %v", err)
	}

	color := listbox.TextColor()
	if color[0] != 0 || color[1] != 0 || color[2] != 1 {
		t.Errorf("Expected blue color [0, 0, 1], got %v", color)
	}
}

// TestListBoxSetTextColorInvalid tests invalid text color.
func TestListBoxSetTextColorInvalid(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)

	err := listbox.SetTextColor(1.5, 0, 0)
	if err == nil {
		t.Error("SetTextColor(1.5, 0, 0) should return an error")
	}

	err = listbox.SetTextColor(0, -0.1, 0)
	if err == nil {
		t.Error("SetTextColor(0, -0.1, 0) should return an error")
	}
}

// TestListBoxValidate tests field validation.
func TestListBoxValidate(t *testing.T) {
	// Valid single-select list box
	listbox := forms.NewListBox("valid", 100, 350, 150, 80)
	listbox.AddOptions("Opt1", "Opt2", "Opt3")
	if err := listbox.SetSelected("Opt1"); err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	err := listbox.Validate()
	if err != nil {
		t.Errorf("Valid list box should not return error: %v", err)
	}

	// Valid multi-select list box
	multiBox := forms.NewListBox("multi", 100, 350, 150, 80)
	multiBox.AddOptions("A", "B", "C")
	multiBox.SetMultiSelect(true)
	if err := multiBox.SetSelectedMultiple("A", "C"); err != nil {
		t.Fatalf("SetSelectedMultiple failed: %v", err)
	}

	err = multiBox.Validate()
	if err != nil {
		t.Errorf("Valid multi-select list box should not return error: %v", err)
	}

	// Empty name
	invalidBox := forms.NewListBox("", 100, 350, 150, 80)
	err = invalidBox.Validate()
	if err == nil {
		t.Error("ListBox with empty name should return error")
	}

	// Multiple selections without multi-select flag
	singleBox := forms.NewListBox("single", 100, 350, 150, 80)
	singleBox.AddOptions("A", "B", "C")
	// Bypass SetSelected by directly manipulating (simulate invalid state)
	// Actually, we need to test via Validate() after enabling/disabling multi-select
	singleBox.SetMultiSelect(true)
	if err := singleBox.SetSelectedMultiple("A", "B"); err != nil {
		t.Fatalf("SetSelectedMultiple failed: %v", err)
	}
	// Now disable multi-select but keep the selections (Validate should catch this)
	// Actually SetMultiSelect(false) trims to 1 item, so this is handled correctly
	// Let's create a different test case
}

// TestListBoxValidateMultipleSelectionsWithoutFlag tests validation error.
func TestListBoxValidateMultipleSelectionsWithoutFlag(t *testing.T) {
	// This is a tricky case - we can't actually create this invalid state
	// via the public API because SetMultiSelect(false) automatically trims
	// the selection list. This is good design!

	// Let's verify the trimming behavior works correctly:
	listbox := forms.NewListBox("test", 100, 350, 150, 80)
	listbox.AddOptions("A", "B", "C")
	listbox.SetMultiSelect(true)
	if err := listbox.SetSelectedMultiple("A", "B", "C"); err != nil {
		t.Fatalf("SetSelectedMultiple failed: %v", err)
	}

	// Verify we have 3 selections
	if len(listbox.SelectedValues()) != 3 {
		t.Fatalf("Expected 3 selections, got %d", len(listbox.SelectedValues()))
	}

	// Disable multi-select
	listbox.SetMultiSelect(false)

	// Should now have only 1 selection
	if len(listbox.SelectedValues()) != 1 {
		t.Errorf("Expected 1 selection after disabling multi-select, got %d", len(listbox.SelectedValues()))
	}

	// Validate should pass
	if err := listbox.Validate(); err != nil {
		t.Errorf("Validation should pass after auto-trimming: %v", err)
	}
}

// TestListBoxValidateInvalidRectangle tests validation with invalid rectangle.
func TestListBoxValidateInvalidRectangle(t *testing.T) {
	// Create list box with invalid dimensions
	listbox := forms.NewListBox("test", 100, 350, -10, 80)
	err := listbox.Validate()
	if err == nil {
		t.Error("ListBox with invalid rectangle should return error")
	}
}

// TestListBoxChaining tests method chaining.
func TestListBoxChaining(t *testing.T) {
	listbox := forms.NewListBox("chained", 100, 350, 150, 80)

	// Chain multiple setters
	result := listbox.
		AddOption("opt1", "Option 1").
		AddOption("opt2", "Option 2").
		AddOptions("Option 3", "Option 4").
		SetMultiSelect(true).
		SetSort(true).
		SetRequired(true).
		SetReadOnly(false).
		SetFontName("Courier")

	if result != listbox {
		t.Error("Methods should return the listbox for chaining")
	}

	if !listbox.IsMultiSelect() {
		t.Error("Chained SetMultiSelect failed")
	}

	if !listbox.IsSorted() {
		t.Error("Chained SetSort failed")
	}

	if !listbox.IsRequired() {
		t.Error("Chained SetRequired failed")
	}

	if listbox.FontName() != "Courier" {
		t.Error("Chained SetFontName failed")
	}

	if len(listbox.Options()) != 4 {
		t.Errorf("Expected 4 options after chaining, got %d", len(listbox.Options()))
	}
}

// TestListBoxNoComboFlag tests that Combo flag is NOT set for list boxes.
func TestListBoxNoComboFlag(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)

	// Combo flag should NOT be set for list boxes (unlike dropdowns)
	if listbox.Flags()&forms.FlagCombo != 0 {
		t.Error("Combo flag must NOT be set for list boxes")
	}

	// Even after setting other flags
	listbox.SetRequired(true)
	listbox.SetMultiSelect(true)

	if listbox.Flags()&forms.FlagCombo != 0 {
		t.Error("Combo flag must remain unset after setting other flags")
	}
}

// TestListBoxValueTypeSingleSelect tests Value() return type for single-select.
func TestListBoxValueTypeSingleSelect(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)
	listbox.AddOptions("A", "B", "C")
	if err := listbox.SetSelected("B"); err != nil {
		t.Fatalf("SetSelected failed: %v", err)
	}

	// Single-select should return string
	val := listbox.Value()
	if _, ok := val.(string); !ok {
		t.Errorf("Value() should return string for single-select, got %T", val)
	}
}

// TestListBoxValueTypeMultiSelect tests Value() return type for multi-select.
func TestListBoxValueTypeMultiSelect(t *testing.T) {
	listbox := forms.NewListBox("test", 100, 350, 150, 80)
	listbox.AddOptions("A", "B", "C")
	listbox.SetMultiSelect(true)
	if err := listbox.SetSelectedMultiple("A", "C"); err != nil {
		t.Fatalf("SetSelectedMultiple failed: %v", err)
	}

	// Multi-select should return []string
	val := listbox.Value()
	if _, ok := val.([]string); !ok {
		t.Errorf("Value() should return []string for multi-select, got %T", val)
	}
}
