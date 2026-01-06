// Package main demonstrates AcroForm checkbox and radio button creation.
//
// This example shows how to create PDF forms with interactive checkboxes
// and radio button groups using the GxPDF Creator API.
//
// Run: go run main.go
// Output: form_buttons.pdf
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
	"github.com/coregx/gxpdf/creator/forms"
)

func main() {
	// Create new PDF creator
	c := creator.New()

	// Add first page - Checkbox examples
	if err := createCheckboxPage(c); err != nil {
		log.Fatalf("Failed to create checkbox page: %v", err)
	}

	// Add second page - Radio button examples
	if err := createRadioButtonPage(c); err != nil {
		log.Fatalf("Failed to create radio button page: %v", err)
	}

	// Write to file
	if err := c.WriteToFile("form_buttons.pdf"); err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Println("Created form_buttons.pdf with checkboxes and radio buttons")
}

// createCheckboxPage creates a page with checkbox examples.
func createCheckboxPage(c *creator.Creator) error {
	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	// Page title
	page.AddText("AcroForm - Checkbox Examples", 50, 750, creator.HelveticaBold, 24)

	// Section 1: Basic checkbox
	page.AddText("1. Basic Checkbox", 50, 700, creator.HelveticaBold, 14)

	agreeBox := forms.NewCheckbox("agree", 50, 670, 15, 15)
	agreeBox.SetLabel("I agree to the terms and conditions")
	if err := page.AddField(agreeBox); err != nil {
		return fmt.Errorf("add agree checkbox: %w", err)
	}
	page.AddText("I agree to the terms and conditions", 70, 670, creator.Helvetica, 12)

	// Section 2: Pre-checked checkbox
	page.AddText("2. Pre-checked Checkbox", 50, 630, creator.HelveticaBold, 14)

	subscribeBox := forms.NewCheckbox("subscribe", 50, 600, 15, 15)
	subscribeBox.SetLabel("Subscribe to newsletter")
	subscribeBox.SetChecked(true)
	if err := page.AddField(subscribeBox); err != nil {
		return fmt.Errorf("add subscribe checkbox: %w", err)
	}
	page.AddText("Subscribe to newsletter (pre-checked)", 70, 600, creator.Helvetica, 12)

	// Section 3: Required checkbox
	page.AddText("3. Required Checkbox", 50, 560, creator.HelveticaBold, 14)

	requiredBox := forms.NewCheckbox("privacy", 50, 530, 15, 15)
	requiredBox.SetLabel("I have read the privacy policy")
	requiredBox.SetRequired(true)
	if err := page.AddField(requiredBox); err != nil {
		return fmt.Errorf("add privacy checkbox: %w", err)
	}
	page.AddText("I have read the privacy policy (required) *", 70, 530, creator.Helvetica, 12)

	// Section 4: Read-only checkbox
	page.AddText("4. Read-only Checkbox", 50, 490, creator.HelveticaBold, 14)

	readonlyBox := forms.NewCheckbox("verified", 50, 460, 15, 15)
	readonlyBox.SetLabel("Account verified")
	readonlyBox.SetChecked(true)
	readonlyBox.SetReadOnly(true)
	if err := page.AddField(readonlyBox); err != nil {
		return fmt.Errorf("add verified checkbox: %w", err)
	}
	page.AddText("Account verified (read-only)", 70, 460, creator.Helvetica, 12)

	// Section 5: Styled checkboxes
	page.AddText("5. Styled Checkboxes (with colors)", 50, 420, creator.HelveticaBold, 14)

	// Red checkbox
	r, g, b := 1.0, 0.0, 0.0
	redBox := forms.NewCheckbox("option1", 50, 390, 15, 15)
	if err := redBox.SetBorderColor(&r, &g, &b); err != nil {
		return fmt.Errorf("set red border: %w", err)
	}
	if err := page.AddField(redBox); err != nil {
		return fmt.Errorf("add red checkbox: %w", err)
	}
	page.AddText("Option 1 (red border)", 70, 390, creator.Helvetica, 12)

	// Blue checkbox with white fill
	r, g, b = 0.0, 0.0, 1.0
	blueBox := forms.NewCheckbox("option2", 50, 360, 15, 15)
	if err := blueBox.SetBorderColor(&r, &g, &b); err != nil {
		return fmt.Errorf("set blue border: %w", err)
	}
	r, g, b = 1.0, 1.0, 1.0
	if err := blueBox.SetFillColor(&r, &g, &b); err != nil {
		return fmt.Errorf("set white fill: %w", err)
	}
	if err := page.AddField(blueBox); err != nil {
		return fmt.Errorf("add blue checkbox: %w", err)
	}
	page.AddText("Option 2 (blue border, white fill)", 70, 360, creator.Helvetica, 12)

	return nil
}

// createRadioButtonPage creates a page with radio button examples.
func createRadioButtonPage(c *creator.Creator) error {
	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	// Page title
	page.AddText("AcroForm - Radio Button Examples", 50, 750, creator.HelveticaBold, 24)

	// Section 1: Basic radio group
	page.AddText("1. Basic Radio Group", 50, 700, creator.HelveticaBold, 14)
	page.AddText("Gender:", 50, 670, creator.Helvetica, 12)

	gender := forms.NewRadioGroup("gender")
	gender.AddOption("male", 50, 650, "Male")
	gender.AddOption("female", 150, 650, "Female")
	gender.AddOption("other", 250, 650, "Other")
	if err := page.AddField(gender); err != nil {
		return fmt.Errorf("add gender radio: %w", err)
	}
	page.AddText("Male", 70, 650, creator.Helvetica, 12)
	page.AddText("Female", 170, 650, creator.Helvetica, 12)
	page.AddText("Other", 270, 650, creator.Helvetica, 12)

	// Section 2: Pre-selected radio group
	page.AddText("2. Pre-selected Radio Group", 50, 610, creator.HelveticaBold, 14)
	page.AddText("Priority:", 50, 580, creator.Helvetica, 12)

	priority := forms.NewRadioGroup("priority")
	priority.AddOption("low", 50, 560, "Low")
	priority.AddOption("medium", 150, 560, "Medium")
	priority.AddOption("high", 250, 560, "High")
	if err := priority.SetSelected("medium"); err != nil {
		return fmt.Errorf("set priority selection: %w", err)
	}
	if err := page.AddField(priority); err != nil {
		return fmt.Errorf("add priority radio: %w", err)
	}
	page.AddText("Low", 70, 560, creator.Helvetica, 12)
	page.AddText("Medium (selected)", 170, 560, creator.Helvetica, 12)
	page.AddText("High", 270, 560, creator.Helvetica, 12)

	// Section 3: Required radio group
	page.AddText("3. Required Radio Group", 50, 520, creator.HelveticaBold, 14)
	page.AddText("Payment method: *", 50, 490, creator.Helvetica, 12)

	payment := forms.NewRadioGroup("payment")
	payment.AddOption("card", 50, 470, "Credit Card")
	payment.AddOption("cash", 150, 470, "Cash")
	payment.AddOption("transfer", 250, 470, "Bank Transfer")
	payment.SetRequired(true)
	if err := page.AddField(payment); err != nil {
		return fmt.Errorf("add payment radio: %w", err)
	}
	page.AddText("Credit Card", 70, 470, creator.Helvetica, 12)
	page.AddText("Cash", 170, 470, creator.Helvetica, 12)
	page.AddText("Bank Transfer", 270, 470, creator.Helvetica, 12)

	// Section 4: Styled radio buttons
	page.AddText("4. Styled Radio Buttons", 50, 430, creator.HelveticaBold, 14)
	page.AddText("Size:", 50, 400, creator.Helvetica, 12)

	size := forms.NewRadioGroup("size")
	size.AddOption("small", 50, 380, "Small", 20, 20)
	size.AddOption("medium", 150, 380, "Medium", 20, 20)
	size.AddOption("large", 250, 380, "Large", 20, 20)

	r, g, b := 0.0, 0.5, 0.0
	if err := size.SetBorderColor(&r, &g, &b); err != nil {
		return fmt.Errorf("set green border: %w", err)
	}
	r, g, b = 0.9, 1.0, 0.9
	if err := size.SetFillColor(&r, &g, &b); err != nil {
		return fmt.Errorf("set light green fill: %w", err)
	}

	if err := page.AddField(size); err != nil {
		return fmt.Errorf("add size radio: %w", err)
	}
	page.AddText("Small", 75, 380, creator.Helvetica, 12)
	page.AddText("Medium", 175, 380, creator.Helvetica, 12)
	page.AddText("Large", 275, 380, creator.Helvetica, 12)

	return nil
}
