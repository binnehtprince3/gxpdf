package gxpdf

import "errors"

// Common errors returned by gxpdf functions.
var (
	// ErrInvalidPDF is returned when the file is not a valid PDF.
	ErrInvalidPDF = errors.New("gxpdf: invalid PDF file")

	// ErrEncrypted is returned when the PDF is encrypted and no password was provided.
	ErrEncrypted = errors.New("gxpdf: PDF is encrypted")

	// ErrWrongPassword is returned when the provided password is incorrect.
	ErrWrongPassword = errors.New("gxpdf: wrong password")

	// ErrCorrupted is returned when the PDF structure is corrupted.
	ErrCorrupted = errors.New("gxpdf: PDF file is corrupted")

	// ErrPageNotFound is returned when the requested page does not exist.
	ErrPageNotFound = errors.New("gxpdf: page not found")

	// ErrNoTables is returned when no tables were found on the page.
	ErrNoTables = errors.New("gxpdf: no tables found")

	// ErrUnsupportedFeature is returned for PDF features not yet implemented.
	ErrUnsupportedFeature = errors.New("gxpdf: unsupported PDF feature")
)

// IsEncrypted returns true if the error indicates an encrypted PDF.
func IsEncrypted(err error) bool {
	return errors.Is(err, ErrEncrypted)
}

// IsCorrupted returns true if the error indicates a corrupted PDF.
func IsCorrupted(err error) bool {
	return errors.Is(err, ErrCorrupted)
}
