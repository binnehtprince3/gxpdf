// Package main demonstrates RC4 encryption support in gxpdf.
//
// This example shows how to create encrypted PDF documents with different
// permission levels.
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Example 1: Basic encryption with user password.
	if err := basicEncryption(); err != nil {
		log.Fatal(err)
	}

	// Example 2: Encryption with different permissions.
	if err := permissionsExample(); err != nil {
		log.Fatal(err)
	}

	// Example 3: Owner and user passwords.
	if err := ownerPasswordExample(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encryption examples completed successfully!")
}

// basicEncryption demonstrates basic PDF encryption.
func basicEncryption() error {
	c := creator.New()
	c.SetTitle("Encrypted Document")

	// Enable encryption with user password.
	err := c.SetEncryption(creator.EncryptionOptions{
		UserPassword: "secret123",
		Permissions:  creator.PermissionPrint | creator.PermissionCopy,
		KeyLength:    128, // 128-bit encryption
	})
	if err != nil {
		return fmt.Errorf("set encryption: %w", err)
	}

	// Add content.
	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	para := creator.NewParagraph("This is an encrypted PDF document.")
	para.SetFont(creator.HelveticaBold, 14)
	if err := page.Draw(para); err != nil {
		return fmt.Errorf("draw paragraph: %w", err)
	}

	// Note: WriteToFile is not yet implemented with encryption support.
	// This is a placeholder for future implementation.
	fmt.Println("Basic encryption example prepared (write not yet implemented)")
	return nil
}

// permissionsExample demonstrates different permission levels.
func permissionsExample() error {
	c := creator.New()
	c.SetTitle("Read-Only Document")

	// Allow only viewing and printing, no modifications.
	err := c.SetEncryption(creator.EncryptionOptions{
		UserPassword: "viewer",
		Permissions:  creator.PermissionPrint, // Only print allowed
		KeyLength:    128,
	})
	if err != nil {
		return fmt.Errorf("set encryption: %w", err)
	}

	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	para := creator.NewParagraph("This PDF is read-only. You can view and print, but not modify.")
	para.SetFont(creator.Helvetica, 12)
	if err := page.Draw(para); err != nil {
		return fmt.Errorf("draw paragraph: %w", err)
	}

	fmt.Println("Permissions example prepared")
	return nil
}

// ownerPasswordExample demonstrates owner and user passwords.
func ownerPasswordExample() error {
	c := creator.New()
	c.SetTitle("Protected Document")

	// Set both user and owner passwords.
	// User can open with restrictions, owner has full access.
	err := c.SetEncryption(creator.EncryptionOptions{
		UserPassword:  "user123",  // Limited access
		OwnerPassword: "owner123", // Full access
		Permissions:   creator.PermissionPrint | creator.PermissionCopy,
		KeyLength:     128,
	})
	if err != nil {
		return fmt.Errorf("set encryption: %w", err)
	}

	page, err := c.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}

	para := creator.NewParagraph("This document has both user and owner passwords.")
	para.SetFont(creator.Helvetica, 12)
	if err := page.Draw(para); err != nil {
		return fmt.Errorf("draw paragraph: %w", err)
	}

	fmt.Println("Owner password example prepared")
	return nil
}
