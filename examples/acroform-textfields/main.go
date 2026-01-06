// Package main demonstrates creating PDF forms with text fields.
//
// This example shows how to create interactive PDF forms using the
// AcroForm (Interactive Forms) API.
//
// This is Phase 1 implementation: Text Fields only.
// Future phases will add checkboxes, radio buttons, combo boxes, etc.
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
	c.SetTitle("User Registration Form")
	c.SetAuthor("GxPDF Demo")

	// Create page
	page, err := c.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	// Add form title
	err = page.AddText("User Registration Form", 100, 750, creator.HelveticaBold, 18)
	if err != nil {
		log.Fatal(err)
	}

	// Add instruction text
	err = page.AddText("Please fill out the form below:", 100, 720, creator.Helvetica, 12)
	if err != nil {
		log.Fatal(err)
	}

	// Create text fields

	// 1. Name field (required)
	nameField := forms.NewTextField("name", 100, 670, 200, 20)
	nameField.SetPlaceholder("Enter your full name")
	nameField.SetRequired(true)

	err = page.AddField(nameField)
	if err != nil {
		log.Fatal(err)
	}

	// Add label for name field
	err = page.AddText("Name*:", 100, 693, creator.Helvetica, 12)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Email field (required)
	emailField := forms.NewTextField("email", 100, 620, 200, 20)
	emailField.SetPlaceholder("Enter your email")
	emailField.SetRequired(true)

	err = page.AddField(emailField)
	if err != nil {
		log.Fatal(err)
	}

	err = page.AddText("Email*:", 100, 643, creator.Helvetica, 12)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Password field (password type)
	passwordField := forms.NewTextField("password", 100, 570, 200, 20)
	passwordField.SetPassword(true)
	passwordField.SetRequired(true)

	err = page.AddField(passwordField)
	if err != nil {
		log.Fatal(err)
	}

	err = page.AddText("Password*:", 100, 593, creator.Helvetica, 12)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Phone field (with max length)
	phoneField := forms.NewTextField("phone", 100, 520, 150, 20)
	phoneField.SetPlaceholder("10 digits")
	if err := phoneField.SetMaxLength(10); err != nil {
		log.Fatal(err)
	}

	err = page.AddField(phoneField)
	if err != nil {
		log.Fatal(err)
	}

	err = page.AddText("Phone:", 100, 543, creator.Helvetica, 12)
	if err != nil {
		log.Fatal(err)
	}

	// 5. Comments field (multiline)
	commentsField := forms.NewTextField("comments", 100, 400, 300, 100)
	commentsField.SetMultiline(true)
	commentsField.SetPlaceholder("Enter any additional comments")

	err = page.AddField(commentsField)
	if err != nil {
		log.Fatal(err)
	}

	err = page.AddText("Comments:", 100, 503, creator.Helvetica, 12)
	if err != nil {
		log.Fatal(err)
	}

	// 6. Read-only field (pre-filled, cannot be edited)
	readonlyField := forms.NewTextField("userid", 100, 350, 200, 20)
	readonlyField.SetValue("USER-12345")
	readonlyField.SetReadOnly(true)

	err = page.AddField(readonlyField)
	if err != nil {
		log.Fatal(err)
	}

	err = page.AddText("User ID (read-only):", 100, 373, creator.Helvetica, 12)
	if err != nil {
		log.Fatal(err)
	}

	// Add note about required fields
	err = page.AddText("* Required fields", 100, 320, creator.HelveticaOblique, 10)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: PDF writing with AcroForm support is not yet fully implemented
	// This example demonstrates the API design.
	// Full PDF writer integration will be completed in a follow-up task.

	fmt.Println("Form structure created successfully!")
	fmt.Println("\nForm Fields:")
	fmt.Println("- Name (required)")
	fmt.Println("- Email (required)")
	fmt.Println("- Password (required, password field)")
	fmt.Println("- Phone (max 10 digits)")
	fmt.Println("- Comments (multiline)")
	fmt.Println("- User ID (read-only)")

	// When PDF writer is complete, uncomment:
	// err = c.WriteToFile("registration_form.pdf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("\nPDF created: registration_form.pdf")
}
