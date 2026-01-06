// Package writer implements PDF writing infrastructure.
package writer

import "strings"

// EscapePDFString escapes a string for use in PDF literal strings.
//
// PDF literal strings are enclosed in parentheses: (Hello World)
//
// Escapes:
//   - \ → \\
//   - ( → \(
//   - ) → \)
//   - \n → \n (newline)
//   - \r → \r (carriage return)
//   - \t → \t (tab)
//   - \b → \b (backspace)
//   - \f → \f (form feed)
//
// Unicode characters (including Cyrillic) are passed through as-is.
// The caller is responsible for encoding them appropriately (usually UTF-16BE
// for text strings in PDF).
//
// Example:
//
//	EscapePDFString("Hello")              // "Hello"
//	EscapePDFString("Price: $50 (USD)")   // "Price: $50 \\(USD\\)"
//	EscapePDFString("Line1\nLine2")       // "Line1\\nLine2"
//	EscapePDFString("C:\\path")           // "C:\\\\path"
//	EscapePDFString("Привет")             // "Привет" (unchanged)
//
// Reference: PDF 1.7 Spec, Table 3 (Escape sequences in literal strings).
func EscapePDFString(s string) string {
	// Order is critical: backslash must be escaped first!
	// Otherwise we'll double-escape the backslashes we just added.
	s = strings.ReplaceAll(s, "\\", "\\\\") // \ → \\

	// Escape parentheses
	s = strings.ReplaceAll(s, "(", "\\(") // ( → \(
	s = strings.ReplaceAll(s, ")", "\\)") // ) → \)

	// Escape control characters
	// Note: These are Go escape sequences that represent actual control characters.
	// We convert them to PDF escape sequences.
	s = strings.ReplaceAll(s, "\n", "\\n") // newline (0x0A)
	s = strings.ReplaceAll(s, "\r", "\\r") // carriage return (0x0D)
	s = strings.ReplaceAll(s, "\t", "\\t") // tab (0x09)
	s = strings.ReplaceAll(s, "\b", "\\b") // backspace (0x08)
	s = strings.ReplaceAll(s, "\f", "\\f") // form feed (0x0C)

	return s
}
