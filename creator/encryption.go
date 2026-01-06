package creator

import (
	"github.com/coregx/gxpdf/internal/security"
)

// EncryptionAlgorithm specifies the encryption algorithm to use.
type EncryptionAlgorithm int

const (
	// EncryptionRC4_40 uses RC4 with 40-bit keys (PDF 1.1+, legacy).
	EncryptionRC4_40 EncryptionAlgorithm = iota

	// EncryptionRC4_128 uses RC4 with 128-bit keys (PDF 1.4+, legacy).
	EncryptionRC4_128

	// EncryptionAES128 uses AES-128 encryption (PDF 1.5+, recommended).
	EncryptionAES128

	// EncryptionAES256 uses AES-256 encryption (PDF 1.7+, most secure).
	EncryptionAES256
)

// EncryptionOptions holds the encryption settings for PDF creation.
//
// Use SetEncryption to enable password protection and permissions control.
//
// Example:
//
//	c := creator.New()
//	c.SetEncryption(creator.EncryptionOptions{
//	    UserPassword:  "userpass",
//	    OwnerPassword: "ownerpass",
//	    Permissions:   creator.PermissionPrint | creator.PermissionCopy,
//	    Algorithm:     creator.EncryptionAES128, // or EncryptionAES256
//	})
type EncryptionOptions struct {
	// UserPassword allows opening the document with restrictions.
	// If empty, document can be opened without password (but with restrictions).
	UserPassword string

	// OwnerPassword allows opening the document with full access.
	// If empty, defaults to UserPassword.
	OwnerPassword string

	// Permissions specifies what operations are allowed.
	// Use Permission constants (PermissionPrint, PermissionCopy, etc.).
	Permissions Permission

	// Algorithm specifies the encryption algorithm.
	// Valid values: EncryptionRC4_40, EncryptionRC4_128, EncryptionAES128, EncryptionAES256.
	// Default: EncryptionAES128.
	Algorithm EncryptionAlgorithm

	// KeyLength specifies the encryption key length in bits (deprecated, use Algorithm instead).
	// Valid values: 40 (PDF 1.1+), 128 (PDF 1.4+), 256 (PDF 1.7+).
	// If Algorithm is set, KeyLength is ignored.
	// Default: 128 (for backward compatibility).
	KeyLength int
}

// Permission represents PDF document permissions.
//
// Multiple permissions can be combined using the OR operator (|).
//
// Example:
//
//	perms := PermissionPrint | PermissionCopy
type Permission = security.Permission

// Permission constants.
const (
	// PermissionPrint allows printing the document.
	PermissionPrint = security.PermissionPrint

	// PermissionModify allows modifying the document.
	PermissionModify = security.PermissionModify

	// PermissionCopy allows copying text and graphics.
	PermissionCopy = security.PermissionCopy

	// PermissionAnnotate allows adding or modifying annotations.
	PermissionAnnotate = security.PermissionAnnotate

	// PermissionFillForms allows filling form fields.
	PermissionFillForms = security.PermissionFillForms

	// PermissionExtract allows extracting text for accessibility.
	PermissionExtract = security.PermissionExtract

	// PermissionAssemble allows assembling the document.
	PermissionAssemble = security.PermissionAssemble

	// PermissionPrintHighQuality allows high-quality printing.
	PermissionPrintHighQuality = security.PermissionPrintHighQuality

	// PermissionAll grants all permissions.
	PermissionAll = security.PermissionAll

	// PermissionNone grants no permissions.
	PermissionNone = security.PermissionNone
)

// SetEncryption enables encryption for the PDF document.
//
// This must be called BEFORE writing the PDF to file.
//
// Example:
//
//	c := creator.New()
//	c.SetEncryption(creator.EncryptionOptions{
//	    UserPassword:  "userpass",
//	    OwnerPassword: "ownerpass",
//	    Permissions:   creator.PermissionPrint | creator.PermissionCopy,
//	    Algorithm:     creator.EncryptionAES128, // Recommended
//	})
//	c.WriteToFile("protected.pdf")
//
// Note: Encryption is applied during WriteToFile.
func (c *Creator) SetEncryption(opts EncryptionOptions) error {
	// If Algorithm is not set but KeyLength is, map KeyLength to Algorithm for backward compatibility.
	if opts.Algorithm == 0 && opts.KeyLength > 0 {
		opts.Algorithm = mapKeyLengthToAlgorithm(opts.KeyLength)
	}

	// Set default algorithm if not specified.
	if opts.Algorithm == 0 {
		opts.Algorithm = EncryptionAES128 // Default to AES-128 (secure and compatible)
	}

	// Validate algorithm.
	if err := validateAlgorithm(opts.Algorithm); err != nil {
		return err
	}

	// Store encryption options for later use during write.
	// Note: Actual encryption setup happens in WriteToFile when we have the FileID.
	c.encryptionOpts = &opts
	return nil
}

// mapKeyLengthToAlgorithm maps KeyLength to Algorithm for backward compatibility.
func mapKeyLengthToAlgorithm(keyLength int) EncryptionAlgorithm {
	switch keyLength {
	case 40:
		return EncryptionRC4_40
	case 128:
		return EncryptionRC4_128
	case 256:
		return EncryptionAES256
	default:
		return EncryptionAES128 // Default
	}
}

// validateAlgorithm validates the encryption algorithm.
func validateAlgorithm(alg EncryptionAlgorithm) error {
	switch alg {
	case EncryptionRC4_40, EncryptionRC4_128, EncryptionAES128, EncryptionAES256:
		return nil
	default:
		return security.ErrUnsupportedVersion
	}
}

// encryptionOpts stores the encryption options.
// This is added to the Creator struct (see creator.go).
